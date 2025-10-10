package main

import (
	"context"
	"errors"
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
	ht "github.com/ogen-go/ogen/http"

	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	orderV1 "github.com/nkolesnikov999/micro2-OK/shared/pkg/openapi/order/v1"
	inventoryV1 "github.com/nkolesnikov999/micro2-OK/shared/pkg/proto/inventory/v1"
)

const (
	httpPort = "8080"
	// адрес inventory gRPC‑сервера
	inventoryAddr = "localhost:50051"
	// Таймауты для HTTP-сервера
	readHeaderTimeout = 5 * time.Second
	shutdownTimeout   = 10 * time.Second
)

// orderStorage представляет потокобезопасное хранилище данных о заказах
type orderStorage struct {
	mu     sync.RWMutex
	orders map[string]*orderV1.OrderDto
}

// NeworderStorage создает новое хранилище данных о заказах
func NewOrderStorage() *orderStorage {
	return &orderStorage{
		orders: make(map[string]*orderV1.OrderDto),
	}
}

func (s *orderStorage) SaveOrder(order *orderV1.OrderDto) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.orders[order.OrderUUID.String()] = order
	return nil
}

// OrderHandler реализует интерфейс orderV1.Handler для обработки запросов к API заказах
type OrderHandler struct {
	storage *orderStorage
}

// NewOrderHandler создает новый обработчик запросов к API заказах
func NewOrderHandler(storage *orderStorage) *OrderHandler {
	return &OrderHandler{
		storage: storage,
	}
}

func (h *OrderHandler) CancelOrder(ctx context.Context, params orderV1.CancelOrderParams) (orderV1.CancelOrderRes, error) {
	return nil, ht.ErrNotImplemented
}

func (h *OrderHandler) CreateOrder(ctx context.Context, req *orderV1.CreateOrderRequest) (orderV1.CreateOrderRes, error) {
	if req == nil {
		return &orderV1.BadRequestError{Code: http.StatusBadRequest, Message: "empty request"}, nil
	}
	if len(req.PartUuids) == 0 {
		return &orderV1.BadRequestError{Code: http.StatusBadRequest, Message: "part_uuids must not be empty"}, nil
	}

	// Build filter for inventory service by requested UUIDs
	uuids := make([]string, 0, len(req.PartUuids))
	for _, id := range req.PartUuids {
		uuids = append(uuids, id.String())
	}

	resp, err := callInventoryListParts(ctx, inventoryAddr, uuids)
	if err != nil {
		return &orderV1.ServiceUnavailableError{Code: http.StatusServiceUnavailable, Message: err.Error()}, nil
	}

	// Ensure all requested parts exist
	found := make(map[string]struct{}, len(resp.GetParts()))
	var total float64
	for _, p := range resp.GetParts() {
		found[p.GetUuid()] = struct{}{}
		total += p.GetPrice()
	}
	missing := make([]string, 0)
	for _, id := range uuids {
		if _, ok := found[id]; !ok {
			missing = append(missing, id)
		}
	}
	if len(missing) > 0 {
		return &orderV1.NotFoundError{Code: http.StatusNotFound, Message: "one or more parts not found"}, nil
	}

	orderID := uuid.New()
	order := &orderV1.OrderDto{
		OrderUUID:  orderID,
		UserUUID:   req.UserUUID,
		PartUuids:  req.PartUuids,
		TotalPrice: float32(total),
		Status:     orderV1.OrderStatusPENDINGPAYMENT,
	}
	if err := h.storage.SaveOrder(order); err != nil {
		return &orderV1.InternalServerError{Code: http.StatusInternalServerError, Message: err.Error()}, nil
	}

	return &orderV1.CreateOrderResponse{
		OrderUUID:  orderID,
		TotalPrice: float32(total),
	}, nil
}

// callInventoryListParts создает gRPC‑клиента, выполняет ListParts и возвращает ответ.
func callInventoryListParts(ctx context.Context, addr string, uuids []string) (*inventoryV1.ListPartsResponse, error) {

	conn, err := grpc.NewClient(
		addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Printf("failed to connect: %v\n", err)
		return nil, err
	}
	defer func() {
		if cerr := conn.Close(); cerr != nil {
			log.Printf("failed to close connect: %v", cerr)
		}
	}()

	client := inventoryV1.NewInventoryServiceClient(conn)
	return client.ListParts(ctx, &inventoryV1.ListPartsRequest{
		Filter: &inventoryV1.PartsFilter{Uuids: uuids},
	})
}

func (h *OrderHandler) GetOrderByUuid(ctx context.Context, params orderV1.GetOrderByUuidParams) (orderV1.GetOrderByUuidRes, error) {
	return nil, ht.ErrNotImplemented
}

func (h *OrderHandler) PayOrder(ctx context.Context, req *orderV1.PayOrderRequest, params orderV1.PayOrderParams) (orderV1.PayOrderRes, error) {
	return nil, ht.ErrNotImplemented
}

func (h *OrderHandler) NewError(ctx context.Context, err error) *orderV1.GenericErrorStatusCode {
	return &orderV1.GenericErrorStatusCode{
		StatusCode: http.StatusInternalServerError,
		Response: orderV1.GenericError{
			Code:    orderV1.NewOptInt(http.StatusInternalServerError),
			Message: orderV1.NewOptString(err.Error()),
		},
	}
}

func main() {
	// Создаем хранилище для данных о заказах
	storage := NewOrderStorage()

	// Создаем обработчик API заказах
	orderHandler := NewOrderHandler(storage)

	// Создаем OpenAPI сервер
	orderServer, err := orderV1.NewServer(
		orderHandler,
		orderV1.WithPathPrefix("/api/v1"),
	)
	if err != nil {
		log.Fatalf("ошибка создания сервера OpenAPI: %v", err)
	}

	// Инициализируем роутер Chi
	r := chi.NewRouter()

	// Добавляем middleware
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(10 * time.Second))

	// Монтируем обработчики OpenAPI
	r.Mount("/", orderServer)

	// Запускаем HTTP-сервер
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
