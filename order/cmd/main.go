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

	"github.com/brianvoe/gofakeit/v7"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/google/uuid"
	ht "github.com/ogen-go/ogen/http"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	orderV1 "github.com/nkolesnikov999/micro2-OK/shared/pkg/openapi/order/v1"
	inventoryV1 "github.com/nkolesnikov999/micro2-OK/shared/pkg/proto/inventory/v1"
	paymentV1 "github.com/nkolesnikov999/micro2-OK/shared/pkg/proto/payment/v1"
)

const (
	httpPort = "8080"
	// адрес inventory gRPC‑сервера
	inventoryAddr = "localhost:50051"
	paymentAddr   = "localhost:50050"
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

// GetOrder возвращает информацию о заказе по uuid
func (s *orderStorage) GetOrder(uuid string) *orderV1.OrderDto {
	s.mu.RLock()
	defer s.mu.RUnlock()

	order, ok := s.orders[uuid]
	if !ok {
		return nil
	}

	return order
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

func randomPaymentMethod() paymentV1.PaymentMethod {
	// Values from proto: 0=UNSPECIFIED, 1=CARD, 2=SBP, 3=CREDIT_CARD, 4=INVESTOR_MONEY
	vals := []paymentV1.PaymentMethod{
		paymentV1.PaymentMethod_PAYMENT_METHOD_UNSPECIFIED,
		paymentV1.PaymentMethod_PAYMENT_METHOD_CARD,
		paymentV1.PaymentMethod_PAYMENT_METHOD_SBP,
		paymentV1.PaymentMethod_PAYMENT_METHOD_CREDIT_CARD,
		paymentV1.PaymentMethod_PAYMENT_METHOD_INVESTOR_MONEY,
	}
	return vals[gofakeit.IntRange(0, len(vals)-1)]
}

func callPaymentService(ctx context.Context, addr string, order *orderV1.OrderDto, paymentMethod paymentV1.PaymentMethod) (*paymentV1.PayOrderResponse, error) {
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

	client := paymentV1.NewPaymentServiceClient(conn)
	return client.PayOrder(ctx, &paymentV1.PayOrderRequest{
		OrderUuid:     order.OrderUUID.String(),
		UserUuid:      order.UserUUID.String(),
		PaymentMethod: paymentMethod,
	})
}

func (h *OrderHandler) GetOrderByUuid(ctx context.Context, params orderV1.GetOrderByUuidParams) (orderV1.GetOrderByUuidRes, error) {
	order := h.storage.GetOrder(params.OrderUUID.String())
	if order == nil {
		return &orderV1.NotFoundError{Code: http.StatusNotFound, Message: "order not found"}, nil
	}
	return order, nil
}

func (h *OrderHandler) PayOrder(ctx context.Context, req *orderV1.PayOrderRequest, params orderV1.PayOrderParams) (orderV1.PayOrderRes, error) {
	order := h.storage.GetOrder(params.OrderUUID.String())
	if order == nil {
		return &orderV1.NotFoundError{Code: http.StatusNotFound, Message: "order not found"}, nil
	}
	paymentMethod := randomPaymentMethod()
	paymentResp, err := callPaymentService(ctx, paymentAddr, order, paymentMethod)
	if err != nil {
		return &orderV1.InternalServerError{Code: http.StatusInternalServerError, Message: err.Error()}, nil
	}
	order.Status = orderV1.OrderStatusPAID
	order.TransactionUUID = orderV1.NewOptNilString(paymentResp.TransactionUuid)
	order.PaymentMethod = orderV1.NewOptPaymentMethod(convertPaymentMethod(paymentMethod))
	if err := h.storage.SaveOrder(order); err != nil {
		return &orderV1.InternalServerError{Code: http.StatusInternalServerError, Message: err.Error()}, nil
	}
	return &orderV1.PayOrderResponse{
		TransactionUUID: paymentResp.TransactionUuid,
	}, nil
}

// convertPaymentMethod maps payment service enum to OpenAPI enum.
func convertPaymentMethod(pm paymentV1.PaymentMethod) orderV1.PaymentMethod {
	switch pm {
	case paymentV1.PaymentMethod_PAYMENT_METHOD_CARD:
		return orderV1.PaymentMethodPAYMENTMETHODCARD
	case paymentV1.PaymentMethod_PAYMENT_METHOD_SBP:
		return orderV1.PaymentMethodPAYMENTMETHODSBP
	case paymentV1.PaymentMethod_PAYMENT_METHOD_CREDIT_CARD:
		return orderV1.PaymentMethodPAYMENTMETHODCREDITCARD
	case paymentV1.PaymentMethod_PAYMENT_METHOD_INVESTOR_MONEY:
		return orderV1.PaymentMethodPAYMENTMETHODINVESTORMONEY
	case paymentV1.PaymentMethod_PAYMENT_METHOD_UNSPECIFIED:
		fallthrough
	default:
		return orderV1.PaymentMethodPAYMENTMETHODUNKNOWN
	}
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
