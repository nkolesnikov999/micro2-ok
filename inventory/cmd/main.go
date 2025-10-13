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

	"github.com/brianvoe/gofakeit/v7"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	inventoryV1 "github.com/nkolesnikov999/micro2-OK/shared/pkg/proto/inventory/v1"
)

// grpcPort — порт, на котором слушает gRPC‑сервер inventory.
const grpcPort = 50051

// inventoryService реализует gRPC‑сервис inventory.
type inventoryService struct {
	inventoryV1.UnimplementedInventoryServiceServer

	mu    sync.RWMutex
	parts map[string]*inventoryV1.Part
}

func (s *inventoryService) GetPart(ctx context.Context, req *inventoryV1.GetPartRequest) (*inventoryV1.GetPartResponse, error) {
	part_uuid := req.GetUuid()
	// Validate UUID format
	if _, err := uuid.Parse(part_uuid); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid uuid format: %v", err)
	}

	s.mu.RLock()
	defer s.mu.RUnlock()
	part, ok := s.parts[part_uuid]
	if !ok {
		return nil, status.Errorf(codes.NotFound, "part not found")
	}
	return &inventoryV1.GetPartResponse{Part: part}, nil
}

func (s *inventoryService) ListParts(ctx context.Context, req *inventoryV1.ListPartsRequest) (*inventoryV1.ListPartsResponse, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	filter := req.GetFilter()

	// If filter is nil or all its fields are empty, return all parts
	if filter == nil || (len(filter.GetUuids()) == 0 && len(filter.GetNames()) == 0 && len(filter.GetCategories()) == 0 && len(filter.GetManufacturerCountries()) == 0 && len(filter.GetTags()) == 0) {
		parts := make([]*inventoryV1.Part, 0, len(s.parts))
		for _, part := range s.parts {
			parts = append(parts, part)
		}
		return &inventoryV1.ListPartsResponse{Parts: parts}, nil
	}

	// Build sets for O(1) membership checks (OR within a single field)
	uuidsSet := makeStringSet(filter.GetUuids())
	namesSet := makeStringSet(filter.GetNames())
	categoriesSet := makeCategorySet(filter.GetCategories())
	countriesSet := makeStringSet(filter.GetManufacturerCountries())
	tagsSet := makeStringSet(filter.GetTags())

	// AND across different fields
	parts := make([]*inventoryV1.Part, 0, len(s.parts))
	for _, part := range s.parts {
		if uuidsSet != nil {
			if _, ok := uuidsSet[part.GetUuid()]; !ok {
				continue
			}
		}
		if namesSet != nil {
			if _, ok := namesSet[part.GetName()]; !ok {
				continue
			}
		}
		if categoriesSet != nil {
			if _, ok := categoriesSet[part.GetCategory()]; !ok {
				continue
			}
		}
		if countriesSet != nil {
			country := ""
			if part.GetManufacturer() != nil {
				country = part.GetManufacturer().GetCountry()
			}
			if _, ok := countriesSet[country]; !ok {
				continue
			}
		}
		if tagsSet != nil {
			if !hasAnyTag(part.GetTags(), tagsSet) {
				continue
			}
		}
		parts = append(parts, part)
	}

	return &inventoryV1.ListPartsResponse{Parts: parts}, nil
}

// makeStringSet creates a set from a slice of strings. Returns nil for empty input.
func makeStringSet(values []string) map[string]struct{} {
	if len(values) == 0 {
		return nil
	}
	set := make(map[string]struct{}, len(values))
	for _, v := range values {
		set[v] = struct{}{}
	}
	return set
}

// makeCategorySet creates a set from a slice of categories. Returns nil for empty input.
func makeCategorySet(values []inventoryV1.Category) map[inventoryV1.Category]struct{} {
	if len(values) == 0 {
		return nil
	}
	set := make(map[inventoryV1.Category]struct{}, len(values))
	for _, v := range values {
		set[v] = struct{}{}
	}
	return set
}

// hasAnyTag returns true if partTags contains at least one tag from wanted.
func hasAnyTag(partTags []string, wanted map[string]struct{}) bool {
	for _, t := range partTags {
		if _, ok := wanted[t]; ok {
			return true
		}
	}
	return false
}

// fakeDimensions generates random dimensions with realistic ranges.
func fakeDimensions() *inventoryV1.Dimensions {
	return &inventoryV1.Dimensions{
		Length: gofakeit.Float64Range(1.0, 300.0),
		Width:  gofakeit.Float64Range(1.0, 300.0),
		Height: gofakeit.Float64Range(0.5, 150.0),
		Weight: gofakeit.Float64Range(0.1, 500.0),
	}
}

// fakeManufacturer generates a random manufacturer.
func fakeManufacturer() *inventoryV1.Manufacturer {
	return &inventoryV1.Manufacturer{
		Name:    gofakeit.Company(),
		Country: gofakeit.Country(),
		Website: gofakeit.URL(),
	}
}

// fakeTags returns a small set of random tags.
func fakeTags() []string {
	n := gofakeit.IntRange(1, 5)
	tags := make([]string, 0, n)
	for i := 0; i < n; i++ {
		tags = append(tags, gofakeit.Word())
	}
	return tags
}

// randomCategory returns a random inventory category (excluding unspecified most of the time).
func randomCategory() inventoryV1.Category {
	// Values from proto: 0=UNSPECIFIED, 1=ENGINE, 2=FUEL, 3=PORTHOLE, 4=WING
	vals := []inventoryV1.Category{
		inventoryV1.Category_CATEGORY_ENGINE,
		inventoryV1.Category_CATEGORY_FUEL,
		inventoryV1.Category_CATEGORY_PORTHOLE,
		inventoryV1.Category_CATEGORY_WING,
	}
	return vals[gofakeit.IntRange(0, len(vals)-1)]
}

func createParts(count int) []*inventoryV1.Part {
	parts := make([]*inventoryV1.Part, 0, count)
	for range count {
		parts = append(parts, &inventoryV1.Part{
			Uuid:          gofakeit.UUID(),
			Name:          gofakeit.Name(),
			Description:   gofakeit.Sentence(10),
			Price:         gofakeit.Price(100, 1000),
			StockQuantity: int64(gofakeit.IntRange(1, 100)),
			Category:      randomCategory(),
			Dimensions:    fakeDimensions(),
			Manufacturer:  fakeManufacturer(),
			Tags:          fakeTags(),
			CreatedAt:     timestamppb.New(gofakeit.Date()),
			UpdatedAt:     timestamppb.New(gofakeit.Date()),
		})
	}
	return parts
}

// initParts creates and initializes the parts inventory with sample data.
func initParts() map[string]*inventoryV1.Part {
	parts := createParts(100)
	partsMap := make(map[string]*inventoryV1.Part)
	for _, part := range parts {
		partsMap[part.GetUuid()] = part
	}
	return partsMap
}

func main() {
	partsMap := initParts()

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

	// Создаем gRPC сервер
	s := grpc.NewServer()

	// Регистрируем наш сервис
	service := &inventoryService{
		parts: partsMap,
	}

	inventoryV1.RegisterInventoryServiceServer(s, service)

	// Включаем рефлексию для отладки
	reflection.Register(s)

	go func() {
		log.Printf("🚀 gRPC server listening on %d\n", grpcPort)
		err = s.Serve(lis)
		if err != nil {
			log.Printf("failed to serve: %v\n", err)
			return
		}
	}()

	// Корректное завершение (graceful shutdown): ждём сигнал и останавливаем сервер.
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("🛑 Shutting down gRPC server...")
	s.GracefulStop()
	log.Println("✅ Server stopped")
}
