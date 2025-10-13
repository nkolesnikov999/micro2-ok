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
	httpPort          = "8080"
	inventoryAddr     = "127.0.0.1:50051"
	paymentAddr       = "127.0.0.1:50050"
	readHeaderTimeout = 5 * time.Second
	shutdownTimeout   = 10 * time.Second
)

var (
	ErrOrderNotFound      = errors.New("order not found")
	ErrOrderAlreadyExists = errors.New("order already exists")
)

// orderStorage представляет потокобезопасное хранилище данных о заказах
type orderStorage struct {
	mu     sync.RWMutex
	orders map[string]*orderV1.OrderDto
}

func NewOrderStorage() *orderStorage {
	return &orderStorage{
		orders: make(map[string]*orderV1.OrderDto),
	}
}

func (s *orderStorage) CreateOrder(order *orderV1.OrderDto) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	orderUUID := order.OrderUUID.String()
	if _, exists := s.orders[orderUUID]; exists {
		return ErrOrderAlreadyExists
	}

	s.orders[orderUUID] = order
	return nil
}

func (s *orderStorage) UpdateOrder(order *orderV1.OrderDto) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	orderUUID := order.OrderUUID.String()
	if _, exists := s.orders[orderUUID]; !exists {
		return ErrOrderNotFound
	}

	s.orders[orderUUID] = order
	return nil
}

func (s *orderStorage) GetOrder(uuid string) (*orderV1.OrderDto, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	order, ok := s.orders[uuid]
	if !ok {
		return nil, ErrOrderNotFound
	}

	return order, nil
}

type OrderHandler struct {
	storage         *orderStorage
	inventoryClient inventoryV1.InventoryServiceClient
	paymentClient   paymentV1.PaymentServiceClient
}

func NewOrderHandler(storage *orderStorage, inventoryClient inventoryV1.InventoryServiceClient, paymentClient paymentV1.PaymentServiceClient) *OrderHandler {
	return &OrderHandler{
		storage:         storage,
		inventoryClient: inventoryClient,
		paymentClient:   paymentClient,
	}
}

func validateUUID(uuidStr, fieldName string) error {
	if _, err := uuid.Parse(uuidStr); err != nil {
		return status.Errorf(codes.InvalidArgument, "invalid %s format: %v", fieldName, err)
	}
	return nil
}

func (h *OrderHandler) CancelOrder(ctx context.Context, params orderV1.CancelOrderParams) (orderV1.CancelOrderRes, error) {
	orderUUIDStr := params.OrderUUID.String()
	if err := validateUUID(orderUUIDStr, "order_uuid"); err != nil {
		return &orderV1.BadRequestError{Code: http.StatusBadRequest, Message: err.Error()}, nil
	}

	order, err := h.storage.GetOrder(orderUUIDStr)
	if err != nil {
		if errors.Is(err, ErrOrderNotFound) {
			return &orderV1.NotFoundError{Code: http.StatusNotFound, Message: "order not found"}, nil
		}
		return &orderV1.InternalServerError{Code: http.StatusInternalServerError, Message: err.Error()}, nil
	}

	if order.Status == orderV1.OrderStatusPAID {
		return &orderV1.ConflictError{Code: http.StatusConflict, Message: "order already paid and cannot be cancelled"}, nil
	}

	if order.Status == orderV1.OrderStatusPENDINGPAYMENT {
		order.Status = orderV1.OrderStatusCANCELLED
		if err := h.storage.UpdateOrder(order); err != nil {
			if errors.Is(err, ErrOrderNotFound) {
				return &orderV1.NotFoundError{Code: http.StatusNotFound, Message: "order not found"}, nil
			}
			return &orderV1.InternalServerError{Code: http.StatusInternalServerError, Message: err.Error()}, nil
		}
	}

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

	if err := validateUUID(req.UserUUID.String(), "user_uuid"); err != nil {
		return &orderV1.BadRequestError{Code: http.StatusBadRequest, Message: err.Error()}, nil
	}

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
	if err := h.storage.CreateOrder(order); err != nil {
		if errors.Is(err, ErrOrderAlreadyExists) {
			// UUID коллизия крайне маловероятна, но теоретически возможна
			log.Printf("CRITICAL: UUID collision detected for order %s", orderID)
		}
		return &orderV1.InternalServerError{Code: http.StatusInternalServerError, Message: err.Error()}, nil
	}

	return &orderV1.CreateOrderResponse{
		OrderUUID:  orderID,
		TotalPrice: float32(total),
	}, nil
}

func (h *OrderHandler) callInventoryListParts(ctx context.Context, uuids []string) (*inventoryV1.ListPartsResponse, error) {
	// Явный таймаут для защиты от зависания при проблемах с внешним сервисом
	grpcCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	return h.inventoryClient.ListParts(grpcCtx, &inventoryV1.ListPartsRequest{
		Filter: &inventoryV1.PartsFilter{Uuids: uuids},
	})
}

func (h *OrderHandler) callPaymentService(ctx context.Context, order *orderV1.OrderDto, paymentMethod paymentV1.PaymentMethod) (*paymentV1.PayOrderResponse, error) {
	// Явный таймаут для защиты от зависания при проблемах с внешним сервисом
	grpcCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	return h.paymentClient.PayOrder(grpcCtx, &paymentV1.PayOrderRequest{
		OrderUuid:     order.OrderUUID.String(),
		UserUuid:      order.UserUUID.String(),
		PaymentMethod: paymentMethod,
	})
}

func (h *OrderHandler) GetOrderByUuid(ctx context.Context, params orderV1.GetOrderByUuidParams) (orderV1.GetOrderByUuidRes, error) {
	orderUUIDStr := params.OrderUUID.String()
	if err := validateUUID(orderUUIDStr, "order_uuid"); err != nil {
		return &orderV1.BadRequestError{Code: http.StatusBadRequest, Message: err.Error()}, nil
	}

	order, err := h.storage.GetOrder(orderUUIDStr)
	if err != nil {
		if errors.Is(err, ErrOrderNotFound) {
			return &orderV1.NotFoundError{Code: http.StatusNotFound, Message: "order not found"}, nil
		}
		return &orderV1.InternalServerError{Code: http.StatusInternalServerError, Message: err.Error()}, nil
	}
	return order, nil
}

func (h *OrderHandler) PayOrder(ctx context.Context, req *orderV1.PayOrderRequest, params orderV1.PayOrderParams) (orderV1.PayOrderRes, error) {
	if req == nil {
		log.Printf("CRITICAL: received nil request in CreateOrder - potential infrastructure issue")
		return &orderV1.InternalServerError{Code: http.StatusInternalServerError, Message: "internal server error"}, nil
	}
	orderUUIDStr := params.OrderUUID.String()
	if err := validateUUID(orderUUIDStr, "order_uuid"); err != nil {
		return &orderV1.BadRequestError{Code: http.StatusBadRequest, Message: err.Error()}, nil
	}

	order, err := h.storage.GetOrder(orderUUIDStr)
	if err != nil {
		if errors.Is(err, ErrOrderNotFound) {
			return &orderV1.NotFoundError{Code: http.StatusNotFound, Message: "order not found"}, nil
		}
		return &orderV1.InternalServerError{Code: http.StatusInternalServerError, Message: err.Error()}, nil
	}

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
	if err := h.storage.UpdateOrder(order); err != nil {
		if errors.Is(err, ErrOrderNotFound) {
			return &orderV1.NotFoundError{Code: http.StatusNotFound, Message: "order not found"}, nil
		}
		return &orderV1.InternalServerError{Code: http.StatusInternalServerError, Message: err.Error()}, nil
	}
	return &orderV1.PayOrderResponse{
		TransactionUUID: paymentResp.TransactionUuid,
	}, nil
}

func randomPaymentMethod() paymentV1.PaymentMethod {
	// Генерируем случайный метод оплаты, исключая UNSPECIFIED (значение 0 из proto)
	vals := []paymentV1.PaymentMethod{
		paymentV1.PaymentMethod_PAYMENT_METHOD_CARD,
		paymentV1.PaymentMethod_PAYMENT_METHOD_SBP,
		paymentV1.PaymentMethod_PAYMENT_METHOD_CREDIT_CARD,
		paymentV1.PaymentMethod_PAYMENT_METHOD_INVESTOR_MONEY,
	}
	return vals[gofakeit.IntRange(0, len(vals)-1)]
}

// convertPaymentMethod маппит enum из payment сервиса в OpenAPI enum
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

func initGRPCConnections() (*grpc.ClientConn, *grpc.ClientConn, error) {
	inventoryConn, err := grpc.NewClient(
		inventoryAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, nil, fmt.Errorf("не удалось подключиться к inventory сервису: %w", err)
	}

	paymentConn, err := grpc.NewClient(
		paymentAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		// Cleanup: закрываем уже открытое соединение при ошибке
		if cerr := inventoryConn.Close(); cerr != nil {
			log.Printf("ошибка при закрытии соединения с inventory: %v", cerr)
		}
		return nil, nil, fmt.Errorf("не удалось подключиться к payment сервису: %w", err)
	}

	return inventoryConn, paymentConn, nil
}

func initApplication() (*grpc.ClientConn, *grpc.ClientConn, *orderV1.Server, error) {
	inventoryConn, paymentConn, err := initGRPCConnections()
	if err != nil {
		return nil, nil, nil, fmt.Errorf("ошибка инициализации gRPC соединений: %w", err)
	}

	inventoryClient := inventoryV1.NewInventoryServiceClient(inventoryConn)
	paymentClient := paymentV1.NewPaymentServiceClient(paymentConn)

	storage := NewOrderStorage()
	orderHandler := NewOrderHandler(storage, inventoryClient, paymentClient)

	orderServer, err := orderV1.NewServer(
		orderHandler,
		orderV1.WithPathPrefix("/api/v1"),
	)
	if err != nil {
		// Cleanup: закрываем соединения при ошибке создания сервера
		if cerr := inventoryConn.Close(); cerr != nil {
			log.Printf("ошибка при закрытии соединения с inventory: %v", cerr)
		}
		if cerr := paymentConn.Close(); cerr != nil {
			log.Printf("ошибка при закрытии соединения с payment: %v", cerr)
		}
		return nil, nil, nil, fmt.Errorf("ошибка создания сервера OpenAPI: %w", err)
	}

	return inventoryConn, paymentConn, orderServer, nil
}

func main() {
	inventoryConn, paymentConn, orderServer, err := initApplication()
	if err != nil {
		log.Fatalf("ошибка инициализации приложения: %v", err)
	}
	defer func() {
		if cerr := inventoryConn.Close(); cerr != nil {
			log.Printf("ошибка при закрытии соединения с inventory: %v", cerr)
		}
	}()
	defer func() {
		if cerr := paymentConn.Close(); cerr != nil {
			log.Printf("ошибка при закрытии соединения с payment: %v", cerr)
		}
	}()

	router := chi.NewRouter()
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(middleware.Timeout(10 * time.Second))
	router.Mount("/", orderServer)

	server := &http.Server{
		Addr:    net.JoinHostPort("localhost", httpPort),
		Handler: router,
		// Защита от Slowloris атак: принудительно закрывает соединение, если клиент
		// не успел отправить все заголовки за отведенное время
		ReadHeaderTimeout: readHeaderTimeout,
	}

	go func() {
		log.Printf("🚀 HTTP-сервер запущен на порту %s\n", httpPort)
		err = server.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Printf("❌ Ошибка запуска сервера: %v\n", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("🛑 Завершение работы сервера...")

	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	err = server.Shutdown(ctx)
	if err != nil {
		log.Printf("❌ Ошибка при остановке сервера: %v\n", err)
	}

	log.Println("✅ Сервер остановлен")
}
