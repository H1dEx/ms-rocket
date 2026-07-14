package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	orderV1 "github.com/H1dEx/ms-rocket/shared/pkg/openapi/order/v1"
	inventoryV1 "github.com/H1dEx/ms-rocket/shared/pkg/proto/inventory/v1"
	paymentV1 "github.com/H1dEx/ms-rocket/shared/pkg/proto/payment/v1"
)

const (
	httpPort = "8080"
	// Таймауты для HTTP-сервера
	readHeaderTimeout = 5 * time.Second
	shutdownTimeout   = 10 * time.Second
	paymentPort       = "localhost:50051"
	inventoryPort     = "localhost:50052"
)

var (
	ErrNotFound = errors.New("not found")
	ErrConflict = errors.New("Conflict")
)

type OrderStorage struct {
	mu              sync.RWMutex
	orders          map[string]*orderV1.OrderDto
	paymentClient   paymentV1.PaymentServiceClient
	inventoryClient inventoryV1.InventoryServiceClient
}

func NewOrderStorage(p paymentV1.PaymentServiceClient, i inventoryV1.InventoryServiceClient) *OrderStorage {
	return &OrderStorage{
		orders:          make(map[string]*orderV1.OrderDto),
		paymentClient:   p,
		inventoryClient: i,
	}
}

func (o *OrderStorage) AddOrder(details *orderV1.OrderDto) {
	o.mu.Lock()
	defer o.mu.Unlock()

	o.orders[details.OrderUUID] = details
}

func (o *OrderStorage) PayOrder(orderUUID string, method orderV1.PaymentMethod, transactionUUID string) error {
	o.mu.Lock()
	defer o.mu.Unlock()
	order, ok := o.orders[orderUUID]
	if !ok {
		return ErrNotFound
	}
	if order.Status != orderV1.OrderStatusPENDINGPAYMENT {
		return ErrConflict
	}
	o.orders[orderUUID].PaymentMethod = orderV1.NewOptPaymentMethod(method)
	o.orders[orderUUID].TransactionUUID = orderV1.NewOptString(transactionUUID)
	o.orders[orderUUID].Status = orderV1.OrderStatusPAID

	return nil
}

func (o *OrderStorage) GetOrder(orderUUID string) *orderV1.OrderDto {
	o.mu.RLock()
	defer o.mu.RUnlock()

	order, ok := o.orders[orderUUID]
	if !ok {
		return nil
	}

	return order
}

func (o *OrderStorage) CancelOrder(orderUUID string) error {
	o.mu.Lock()
	defer o.mu.Unlock()

	order, ok := o.orders[orderUUID]

	if !ok {
		return ErrNotFound
	}

	if order.Status != orderV1.OrderStatusPENDINGPAYMENT {
		return ErrConflict
	}

	o.orders[orderUUID].Status = orderV1.OrderStatusCANCELLED
	return nil
}

type OrderHandler struct {
	storage *OrderStorage
}

func NewOrderHandler(storage *OrderStorage) *OrderHandler {
	return &OrderHandler{
		storage: storage,
	}
}

func (h *OrderHandler) CreateOrder(ctx context.Context, req *orderV1.CreateOrderRequest) (orderV1.CreateOrderRes, error) {
	filters := &inventoryV1.PartsFilter{Uuids: req.PartUuids}
	res, err := h.storage.inventoryClient.ListParts(ctx, &inventoryV1.ListPartsRequest{Filter: filters})
	if err != nil {
		return &orderV1.InternalServerError{
			Code:    500,
			Message: fmt.Sprintf("list getting error: %s", err.Error()),
		}, nil
	}

	if len(res.Parts) != len(req.PartUuids) {
		return &orderV1.NotFoundError{Code: 404, Message: "Some parts not exist"}, nil
	}

	var sum float32
	for _, p := range res.GetParts() {
		sum += float32(p.Price)
	}

	uuid := uuid.New().String()
	order := &orderV1.OrderDto{
		OrderUUID:  uuid,
		UserUUID:   req.UserUUID,
		PartUuids:  req.PartUuids,
		TotalPrice: sum,
		Status:     orderV1.OrderStatusPENDINGPAYMENT,
	}
	h.storage.AddOrder(order)
	return &orderV1.CreateOrderResponse{OrderUUID: uuid, TotalPrice: sum}, nil
}

func (h *OrderHandler) GetOrderByID(ctx context.Context, params orderV1.GetOrderByIDParams) (orderV1.GetOrderByIDRes, error) {
	order := h.storage.GetOrder(params.OrderUUID)
	if order == nil {
		return &orderV1.NotFoundError{
			Code:    404,
			Message: fmt.Sprintf("Order with id %s is not found", params.OrderUUID),
		}, nil
	}

	return &orderV1.GetOrderResponse{Order: *order}, nil
}

func (h *OrderHandler) OrderCancelById(ctx context.Context, params orderV1.OrderCancelByIdParams) (orderV1.OrderCancelByIdRes, error) {
	if err := h.storage.CancelOrder(params.OrderUUID); err != nil {
		if errors.Is(err, ErrNotFound) {
			return &orderV1.NotFoundError{
				Code:    404,
				Message: fmt.Sprintf("Order with id %s is not found", params.OrderUUID),
			}, nil
		}
		if errors.Is(err, ErrConflict) {
			return &orderV1.ConflictError{
				Code:    409,
				Message: fmt.Sprintf("Conflict error with id %s", params.OrderUUID),
			}, nil
		}
		return &orderV1.InternalServerError{Code: 500, Message: "Service error"}, nil
	}
	return &orderV1.OrderCancelByIdNoContent{}, nil
}

func ToPaymentPaymentMethod(m orderV1.PaymentMethod) paymentV1.PaymentMethod {
	switch m {
	case orderV1.PaymentMethodCARD:
		return paymentV1.PaymentMethod_PAYMENT_METHOD_CARD
	case orderV1.PaymentMethodSBP:
		return paymentV1.PaymentMethod_PAYMENT_METHOD_SBP
	case orderV1.PaymentMethodCREDITCARD:
		return paymentV1.PaymentMethod_PAYMENT_METHOD_CREDIT_CARD
	case orderV1.PaymentMethodINVESTORMONEY:
		return paymentV1.PaymentMethod_PAYMENT_METHOD_INVESTOR_MONEY
	default:
		return paymentV1.PaymentMethod_PAYMENT_METHOD_UNKNOWN
	}
}

func (h *OrderHandler) PayOrderById(ctx context.Context, req *orderV1.PayOrderRequest, params orderV1.PayOrderByIdParams) (orderV1.PayOrderByIdRes, error) {
	paymentMethod := ToPaymentPaymentMethod(req.GetPaymentMethod())
	if paymentMethod == paymentV1.PaymentMethod_PAYMENT_METHOD_UNKNOWN || req.PaymentMethod == orderV1.PaymentMethodUNKNOWN {
		return &orderV1.BadRequestError{
			Code:    400,
			Message: fmt.Sprintf("Unknown payment method %s", paymentMethod),
		}, nil
	}
	o := h.storage.GetOrder(params.OrderUUID)
	if o == nil {
		return &orderV1.NotFoundError{
			Code:    404,
			Message: fmt.Sprintf("Order with id %s is not found", params.OrderUUID),
		}, nil
	}

	p := &paymentV1.PayOrderRequest{PaymentMethod: paymentMethod}
	res, err := h.storage.paymentClient.PayOrder(ctx, p)
	if err != nil {
		return &orderV1.InternalServerError{
			Code:    500,
			Message: fmt.Sprintf("payment error %s", err.Error()),
		}, nil
	}
	transactionID := res.GetTransactionUuid()
	if err := h.storage.PayOrder(params.OrderUUID, req.PaymentMethod, res.GetTransactionUuid()); err != nil {
		if errors.Is(err, ErrConflict) {
			return &orderV1.ConflictError{
				Code:    409,
				Message: fmt.Sprintf("Conflict error with id %s", params.OrderUUID),
			}, nil
		}
		return &orderV1.InternalServerError{Code: 500, Message: "Service error"}, nil
	}
	return &orderV1.PayOrderResponse{TransactionUUID: transactionID}, nil
}

func (h *OrderHandler) NewError(ctx context.Context, err error) *orderV1.GenericErrorStatusCode {
	return &orderV1.GenericErrorStatusCode{
		StatusCode: http.StatusInternalServerError,
		Response: orderV1.GenericError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		},
	}
}

func main() {
	paymentConn, err := grpc.NewClient(
		paymentPort,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Printf("failed to connect: %v\n", err)
		return
	}
	defer func() {
		if cerr := paymentConn.Close(); cerr != nil {
			log.Printf("failed to close connect: %v", cerr)
		}
	}()

	paymentClient := paymentV1.NewPaymentServiceClient(paymentConn)

	inventoryConn, err := grpc.NewClient(
		inventoryPort,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Printf("failed to connect: %v\n", err)
		return
	}
	defer func() {
		if cerr := inventoryConn.Close(); cerr != nil {
			log.Printf("failed to close connect: %v", cerr)
		}
	}()

	inventoryClient := inventoryV1.NewInventoryServiceClient(inventoryConn)

	storage := NewOrderStorage(paymentClient, inventoryClient)
	orderHandler := NewOrderHandler(storage)

	orderServer, err := orderV1.NewServer(orderHandler)
	if err != nil {
		log.Printf("ошибка создания сервера OpenAPI: %v", err)
		return
	}
	// Инициализируем роутер Chi
	r := chi.NewRouter()

	// Добавляем middleware
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(10 * time.Second))
	// r.Use(customMiddleware.RequestLogger)

	r.Mount("/", orderServer)

	server := &http.Server{
		Addr:              net.JoinHostPort("localhost", httpPort),
		Handler:           r,
		ReadHeaderTimeout: readHeaderTimeout, // Защита от Slowloris атак - тип DDoS-атаки, при которой
		// атакующий умышленно медленно отправляет HTTP-заголовки, удерживая соединения открытыми и истощая
		// пул доступных соединений на сервере. ReadHeaderTimeout принудительно закрывает соединение,
		// если клиент не успел отправить все заголовки за отведенное время.
	}

	// Запускаем сервер в отдельной горутине
	go func() {
		log.Printf("🚀 HTTP-сервер запущен на порту %s\n", httpPort)
		err = server.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Printf("❌ Ошибка запуска сервера: %v\n", err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("🛑 Завершение работы сервера...")

	// Создаем контекст с таймаутом для остановки сервера
	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	err = server.Shutdown(ctx)
	if err != nil {
		log.Printf("❌ Ошибка при остановке сервера: %v\n", err)
	}

	log.Println("✅ Сервер остановлен")
}
