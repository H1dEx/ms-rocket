package telegram

import (
	"bytes"
	"context"
	"embed"
	"text/template"

	"go.uber.org/zap"

	"github.com/H1dEx/ms-rocket/notification/internal/client/http"
	"github.com/H1dEx/ms-rocket/notification/internal/model"
	def "github.com/H1dEx/ms-rocket/notification/internal/service"
	"github.com/H1dEx/ms-rocket/platform/pkg/logger"
)

var _ def.TelegramService = (*service)(nil)

const chatID = -5345253226

//go:embed templates/paid_notification.tmpl
var paidFS embed.FS

//go:embed templates/assembled_notification.tmpl
var assembledFS embed.FS

type assembledTemplateData struct {
	EventUUID        string
	OrderUUID        string
	UserUUID         string
	BuildTimeSeconds int64
}

type paidTemplateData struct {
	EventUUID       string
	OrderUUID       string
	UserUUID        string
	PaymentMethod   string
	TransactionUUID string
}

type service struct {
	telegramClient http.TelegramClient
}

var (
	assembledTemplate = template.Must(template.ParseFS(assembledFS, "templates/assembled_notification.tmpl"))
	paidTemplate      = template.Must(template.ParseFS(paidFS, "templates/paid_notification.tmpl"))
)

func NewService(telegramClient http.TelegramClient) *service {
	return &service{
		telegramClient: telegramClient,
	}
}

func (s *service) buildPaidTemplate(data model.OrderPaidEvent) (string, error) {
	paidTemplateData := paidTemplateData{
		EventUUID:       data.EventUUID,
		OrderUUID:       data.OrderUUID,
		UserUUID:        data.UserUUID,
		PaymentMethod:   string(data.PaymentMethod),
		TransactionUUID: data.TransactionUUID,
	}
	var buf bytes.Buffer
	if err := paidTemplate.Execute(&buf, paidTemplateData); err != nil {
		return "", err
	}
	return buf.String(), nil
}

func (s *service) buildAssembledTemplate(data model.ShipAssembledEvent) (string, error) {
	assembledTemplateData := assembledTemplateData{
		EventUUID:        data.EventUUID,
		OrderUUID:        data.OrderUUID,
		UserUUID:         data.UserUUID,
		BuildTimeSeconds: data.BuildTimeSeconds,
	}

	var buf bytes.Buffer
	if err := assembledTemplate.Execute(&buf, assembledTemplateData); err != nil {
		return "", err
	}
	return buf.String(), nil
}

func (s *service) SendPaidNotification(ctx context.Context, data model.OrderPaidEvent) error {
	message, err := s.buildPaidTemplate(data)
	if err != nil {
		return err
	}
	err = s.telegramClient.SendMessage(ctx, chatID, message)
	if err != nil {
		return err
	}
	logger.Info(ctx, "paid notification sent", zap.String("message", message))
	return nil
}

func (s *service) SendAssembledNotification(ctx context.Context, data model.ShipAssembledEvent) error {
	message, err := s.buildAssembledTemplate(data)
	if err != nil {
		return err
	}
	err = s.telegramClient.SendMessage(ctx, chatID, message)
	if err != nil {
		return err
	}
	logger.Info(ctx, "assembled notification sent", zap.String("message", message))
	return nil
}
