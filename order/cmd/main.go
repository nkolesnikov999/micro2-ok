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
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"

	orderV1 "github.com/nkolesnikov999/micro2-OK/shared/pkg/openapi/order/v1"
	inventoryV1 "github.com/nkolesnikov999/micro2-OK/shared/pkg/proto/inventory/v1"
	paymentV1 "github.com/nkolesnikov999/micro2-OK/shared/pkg/proto/payment/v1"
)

const (
	httpPort = "8080"
	// адрес inventory gRPC‑сервера
	inventoryAddr = "127.0.0.1:50051"
	paymentAddr   = "127.0.0.1:50050"
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
	storage         *orderStorage
	inventoryClient inventoryV1.InventoryServiceClient
	paymentClient   paymentV1.PaymentServiceClient
}

// NewOrderHandler создает новый обработчик запросов к API заказах
func NewOrderHandler(storage *orderStorage, inventoryClient inventoryV1.InventoryServiceClient, paymentClient paymentV1.PaymentServiceClient) *OrderHandler {
	return &OrderHandler{
		storage:         storage,
		inventoryClient: inventoryClient,
		paymentClient:   paymentClient,
	}
}

// validateUUID проверяет корректность формата UUID
func validateUUID(uuidStr, fieldName string) error {
	if _, err := uuid.Parse(uuidStr); err != nil {
		return status.Errorf(codes.InvalidArgument, "invalid %s format: %v", fieldName, err)
	}
	return nil
}

func (h *OrderHandler) CancelOrder(ctx context.Context, params orderV1.CancelOrderParams) (orderV1.CancelOrderRes, error) {
	// Validate UUID format
	orderUUIDStr := params.OrderUUID.String()
	if err := validateUUID(orderUUIDStr, "order_uuid"); err != nil {
		return &orderV1.BadRequestError{Code: http.StatusBadRequest, Message: err.Error()}, nil
	}

	order := h.storage.GetOrder(orderUUIDStr)
	if order == nil {
		return &orderV1.NotFoundError{Code: http.StatusNotFound, Message: "order not found"}, nil
	}

	// If already paid, cannot be cancelled
	if order.Status == orderV1.OrderStatusPAID {
		return &orderV1.ConflictError{Code: http.StatusConflict, Message: "order already paid and cannot be cancelled"}, nil
	}

	// If waiting for payment, cancel it
	if order.Status == orderV1.OrderStatusPENDINGPAYMENT {
		order.Status = orderV1.OrderStatusCANCELLED
		if err := h.storage.SaveOrder(order); err != nil {
			return &orderV1.InternalServerError{Code: http.StatusInternalServerError, Message: err.Error()}, nil
		}
	}

	// Return 204 No Content on successful cancellation (or if already cancelled)
	return nil, &orderV1.GenericErrorStatusCode{
		StatusCode: http.StatusNoContent,
		Response:   orderV1.GenericError{},
	}
}

func (h *OrderHandler) CreateOrder(ctx context.Context, req *orderV1.CreateOrderRequest) (orderV1.CreateOrderRes, error) {
	if req == nil {
		log.Printf("CRITICAL: received nil request in CreateOrder - potential infrastructure issue")
		return &orderV1.InternalServerError{Code: http.StatusInternalServerError, Message: "internal server error"}, nil
	}
	if len(req.PartUuids) == 0 {
		return &orderV1.BadRequestError{Code: http.StatusBadRequest, Message: "part_uuids must not be empty"}, nil
	}

	// Validate user UUID format
	if err := validateUUID(req.UserUUID.String(), "user_uuid"); err != nil {
		return &orderV1.BadRequestError{Code: http.StatusBadRequest, Message: err.Error()}, nil
	}

	// Build filter for inventory service by requested UUIDs and validate them
	uuids := make([]string, 0, len(req.PartUuids))
	for _, id := range req.PartUuids {
		uuidStr := id.String()
		if err := validateUUID(uuidStr, "part_uuid"); err != nil {
			return &orderV1.BadRequestError{Code: http.StatusBadRequest, Message: err.Error()}, nil
		}
		uuids = append(uuids, uuidStr)
	}

	inventoryResp, err := h.callInventoryListParts(ctx, uuids)
	if err != nil {
		return &orderV1.ServiceUnavailableError{Code: http.StatusServiceUnavailable, Message: err.Error()}, nil
	}

	// Ensure all requested parts exist
	found := make(map[string]struct{}, len(inventoryResp.GetParts()))
	var total float64
	for _, part := range inventoryResp.GetParts() {
		found[part.GetUuid()] = struct{}{}
		total += part.GetPrice()
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

// callInventoryListParts выполняет ListParts через gRPC клиент
func (h *OrderHandler) callInventoryListParts(ctx context.Context, uuids []string) (*inventoryV1.ListPartsResponse, error) {
	// Добавляем явный таймаут для gRPC вызова (например, 5 секунд)
	grpcCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	return h.inventoryClient.ListParts(grpcCtx, &inventoryV1.ListPartsRequest{
		Filter: &inventoryV1.PartsFilter{Uuids: uuids},
	})
}

// callPaymentService выполняет PayOrder через gRPC клиент
func (h *OrderHandler) callPaymentService(ctx context.Context, order *orderV1.OrderDto, paymentMethod paymentV1.PaymentMethod) (*paymentV1.PayOrderResponse, error) {
	// Добавляем явный таймаут для gRPC вызова
	grpcCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	return h.paymentClient.PayOrder(grpcCtx, &paymentV1.PayOrderRequest{
		OrderUuid:     order.OrderUUID.String(),
		UserUuid:      order.UserUUID.String(),
		PaymentMethod: paymentMethod,
	})
}

func (h *OrderHandler) GetOrderByUuid(ctx context.Context, params orderV1.GetOrderByUuidParams) (orderV1.GetOrderByUuidRes, error) {
	// Validate UUID format
	orderUUIDStr := params.OrderUUID.String()
	if err := validateUUID(orderUUIDStr, "order_uuid"); err != nil {
		return &orderV1.BadRequestError{Code: http.StatusBadRequest, Message: err.Error()}, nil
	}

	order := h.storage.GetOrder(orderUUIDStr)
	if order == nil {
		return &orderV1.NotFoundError{Code: http.StatusNotFound, Message: "order not found"}, nil
	}
	return order, nil
}

func (h *OrderHandler) PayOrder(ctx context.Context, req *orderV1.PayOrderRequest, params orderV1.PayOrderParams) (orderV1.PayOrderRes, error) {
	if req == nil {
		log.Printf("CRITICAL: received nil request in CreateOrder - potential infrastructure issue")
		return &orderV1.InternalServerError{Code: http.StatusInternalServerError, Message: "internal server error"}, nil
	}
	// Validate UUID format
	orderUUIDStr := params.OrderUUID.String()
	if err := validateUUID(orderUUIDStr, "order_uuid"); err != nil {
		return &orderV1.BadRequestError{Code: http.StatusBadRequest, Message: err.Error()}, nil
	}

	order := h.storage.GetOrder(orderUUIDStr)
	if order == nil {
		return &orderV1.NotFoundError{Code: http.StatusNotFound, Message: "order not found"}, nil
	}

	// Validate order status before payment attempt
	if order.Status == orderV1.OrderStatusPAID {
		return &orderV1.ConflictError{Code: http.StatusConflict, Message: "order already paid"}, nil
	}
	if order.Status == orderV1.OrderStatusCANCELLED {
		return &orderV1.ConflictError{Code: http.StatusConflict, Message: "cannot pay cancelled order"}, nil
	}

	paymentMethod := randomPaymentMethod()
	paymentResp, err := h.callPaymentService(ctx, order, paymentMethod)
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

func randomPaymentMethod() paymentV1.PaymentMethod {
	// Values from proto: 0=UNSPECIFIED, 1=CARD, 2=SBP, 3=CREDIT_CARD, 4=INVESTOR_MONEY
	vals := []paymentV1.PaymentMethod{
		paymentV1.PaymentMethod_PAYMENT_METHOD_CARD,
		paymentV1.PaymentMethod_PAYMENT_METHOD_SBP,
		paymentV1.PaymentMethod_PAYMENT_METHOD_CREDIT_CARD,
		paymentV1.PaymentMethod_PAYMENT_METHOD_INVESTOR_MONEY,
	}
	return vals[gofakeit.IntRange(0, len(vals)-1)]
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
	// Создаем gRPC соединение с inventory сервисом
	inventoryConn, err := grpc.NewClient(
		inventoryAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalf("не удалось подключиться к inventory сервису: %v", err)
	}
	defer func() {
		if cerr := inventoryConn.Close(); cerr != nil {
			log.Printf("ошибка при закрытии соединения с inventory: %v", cerr)
		}
	}()

	// Создаем gRPC соединение с payment сервисом
	paymentConn, err := grpc.NewClient(
		paymentAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalf("не удалось подключиться к payment сервису: %v", err)
	}
	defer func() {
		if cerr := paymentConn.Close(); cerr != nil {
			log.Printf("ошибка при закрытии соединения с payment: %v", cerr)
		}
	}()

	// Создаем gRPC клиенты
	inventoryClient := inventoryV1.NewInventoryServiceClient(inventoryConn)
	paymentClient := paymentV1.NewPaymentServiceClient(paymentConn)

	// Создаем хранилище для данных о заказах
	storage := NewOrderStorage()

	// Создаем обработчик API заказах с gRPC клиентами
	orderHandler := NewOrderHandler(storage, inventoryClient, paymentClient)

	// Создаем OpenAPI сервер
	orderServer, err := orderV1.NewServer(
		orderHandler,
		orderV1.WithPathPrefix("/api/v1"),
	)
	if err != nil {
		log.Fatalf("ошибка создания сервера OpenAPI: %v", err)
	}

	// Инициализируем роутер Chi
	router := chi.NewRouter()

	// Добавляем middleware
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(middleware.Timeout(10 * time.Second))

	// Монтируем обработчики OpenAPI
	router.Mount("/", orderServer)

	// Запускаем HTTP-сервер
	server := &http.Server{
		Addr:              net.JoinHostPort("localhost", httpPort),
		Handler:           router,
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
