package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"

	inventoryV1 "github.com/H1dEx/ms-rocket/shared/pkg/proto/inventory/v1"
)

const grpcPort = 50052

type inventoryService struct {
	inventoryV1.UnimplementedInventoryServiceServer

	mu    sync.RWMutex
	parts map[string]*inventoryV1.Part
}

func (s *inventoryService) GetPart(ctx context.Context, req *inventoryV1.GetPartRequest) (*inventoryV1.GetPartResponse, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	part, ok := s.parts[req.GetUuid()]
	if !ok {
		return nil, status.Errorf(codes.NotFound, "part with UUID %s not found", req.GetUuid())
	}

	return &inventoryV1.GetPartResponse{Part: part}, nil
}

type FilterSets struct {
	Uuids                 map[string]struct{}
	Names                 map[string]struct{}
	Categories            map[string]struct{}
	ManufacturerCountries map[string]struct{}
	Tags                  map[string]struct{}
}

func BuildFilterSets(f *inventoryV1.PartsFilter) FilterSets {
	if f == nil {
		return FilterSets{}
	}
	return FilterSets{
		Uuids:                 filterToMap(f.Uuids),
		Names:                 filterToMap(f.Names),
		Categories:            filterToMap(f.Categories),
		ManufacturerCountries: filterToMap(f.ManufacturerCountries),
		Tags:                  filterToMap(f.Tags),
	}
}

func filterToMap(f []string) map[string]struct{} {
	if f == nil {
		return nil
	}
	acc := make(map[string]struct{}, len(f))
	for _, v := range f {
		acc[v] = struct{}{}
	}
	return acc
}

func hasAnyTag(tags []string, f map[string]struct{}) bool {
	for _, t := range tags {
		if _, ok := f[t]; ok {
			return true
		}
	}
	return false
}

func filterParts(parts map[string]*inventoryV1.Part, f FilterSets) []*inventoryV1.Part {
	result := []*inventoryV1.Part{}
	for _, p := range parts {
		if f.Uuids != nil {
			if _, ok := f.Uuids[p.Uuid]; !ok {
				continue
			}
		}
		if f.Names != nil {
			if _, ok := f.Names[p.Uuid]; !ok {
				continue
			}
		}
		if f.Categories != nil {
			if _, ok := f.Categories[p.Uuid]; !ok {
				continue
			}
		}
		if f.ManufacturerCountries != nil {
			if _, ok := f.ManufacturerCountries[p.Uuid]; !ok {
				continue
			}
		}
		if f.Tags != nil {
			if !hasAnyTag(p.Tags, f.Tags) {
				continue
			}
		}
		result = append(result, p)
	}
	return result
}

func (s *inventoryService) ListParts(ctx context.Context, req *inventoryV1.ListPartsRequest) (*inventoryV1.ListPartsResponse, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	f := BuildFilterSets(req.GetFilter())
	res := filterParts(s.parts, f)
	return &inventoryV1.ListPartsResponse{Parts: res}, nil
}

func main() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", grpcPort))
	if err != nil {
		log.Printf("failed to listen: %v\n", err)
		return
	}

	defer func() {
		if cerr := lis.Close(); cerr != nil {
			log.Printf("failed to close listener: %v\n", cerr)
		}
	}()

	s := grpc.NewServer()
	firstMock := &inventoryV1.Part{
		Uuid:  "111",
		Name:  "First detail",
		Price: 100,
	}
	secondMock := &inventoryV1.Part{
		Uuid:  "222",
		Name:  "Second detail",
		Price: 200,
	}
	service := &inventoryService{
		parts: map[string]*inventoryV1.Part{},
	}

	service.parts[firstMock.Uuid] = firstMock
	service.parts[secondMock.Uuid] = secondMock

	inventoryV1.RegisterInventoryServiceServer(s, service)
	reflection.Register(s)

	go func() {
		log.Printf("🚀 gRPC server listening on %d\n", grpcPort)
		err = s.Serve(lis)
		if err != nil {
			log.Printf("failed to serve: %v\n", err)
			return
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("🛑 Shutting down gRPC server...")
	s.GracefulStop()
	log.Println("✅ Server stopped")
}
