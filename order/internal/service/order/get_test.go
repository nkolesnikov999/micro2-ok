package order

import (
	"github.com/brianvoe/gofakeit/v7"
	"github.com/google/uuid"

	"github.com/nkolesnikov999/micro2-OK/order/internal/model"
)

func (s *ServiceSuite) TestGetOrderSuccess() {
	order := model.Order{
		OrderUUID:       uuid.New(),
		UserUUID:        uuid.New(),
		PartUuids:       []uuid.UUID{uuid.New(), uuid.New()},
		TotalPrice:      gofakeit.Price(100, 1000),
		TransactionUUID: uuid.New().String(),
		PaymentMethod:   "CARD",
		Status:          "PAID",
	}

	s.orderRepository.On("GetOrder", s.ctx, order.OrderUUID.String()).Return(order, nil)

	res, err := s.service.GetOrder(s.ctx, order.OrderUUID)
	s.NoError(err)
	s.Equal(order, res)
}

func (s *ServiceSuite) TestGetOrderNotFound() {
	orderUUID := uuid.New()

	s.orderRepository.On("GetOrder", s.ctx, orderUUID.String()).Return(model.Order{}, model.ErrOrderNotFound)

	res, err := s.service.GetOrder(s.ctx, orderUUID)
	s.Error(err)
	s.ErrorIs(err, model.ErrOrderNotFound)
	s.Empty(res)
}

func (s *ServiceSuite) TestGetOrderRepositoryError() {
	orderUUID := uuid.New()
	repoErr := gofakeit.Error()

	s.orderRepository.On("GetOrder", s.ctx, orderUUID.String()).Return(model.Order{}, repoErr)

	res, err := s.service.GetOrder(s.ctx, orderUUID)
	s.Error(err)
	s.ErrorIs(err, model.ErrOrderGetFailed)
	s.Empty(res)
}

func (s *ServiceSuite) TestGetOrderWithPendingStatus() {
	order := model.Order{
		OrderUUID:       uuid.New(),
		UserUUID:        uuid.New(),
		PartUuids:       []uuid.UUID{uuid.New()},
		TotalPrice:      gofakeit.Price(100, 1000),
		TransactionUUID: "",
		PaymentMethod:   "",
		Status:          "PENDING_PAYMENT",
	}

	s.orderRepository.On("GetOrder", s.ctx, order.OrderUUID.String()).Return(order, nil)

	res, err := s.service.GetOrder(s.ctx, order.OrderUUID)
	s.NoError(err)
	s.Equal(order, res)
	s.Equal("PENDING_PAYMENT", res.Status)
	s.Empty(res.TransactionUUID)
	s.Empty(res.PaymentMethod)
}

func (s *ServiceSuite) TestGetOrderWithCancelledStatus() {
	order := model.Order{
		OrderUUID:       uuid.New(),
		UserUUID:        uuid.New(),
		PartUuids:       []uuid.UUID{uuid.New()},
		TotalPrice:      gofakeit.Price(100, 1000),
		TransactionUUID: "",
		PaymentMethod:   "",
		Status:          "CANCELLED",
	}

	s.orderRepository.On("GetOrder", s.ctx, order.OrderUUID.String()).Return(order, nil)

	res, err := s.service.GetOrder(s.ctx, order.OrderUUID)
	s.NoError(err)
	s.Equal(order, res)
	s.Equal("CANCELLED", res.Status)
}

func (s *ServiceSuite) TestGetOrderWithEmptyPartUUIDs() {
	order := model.Order{
		OrderUUID:       uuid.New(),
		UserUUID:        uuid.New(),
		PartUuids:       []uuid.UUID{}, // empty parts
		TotalPrice:      0.0,
		TransactionUUID: "",
		PaymentMethod:   "",
		Status:          "PENDING_PAYMENT",
	}

	s.orderRepository.On("GetOrder", s.ctx, order.OrderUUID.String()).Return(order, nil)

	res, err := s.service.GetOrder(s.ctx, order.OrderUUID)
	s.NoError(err)
	s.Equal(order, res)
	s.Empty(res.PartUuids)
}

func (s *ServiceSuite) TestGetOrderWithManyParts() {
	partUUIDs := make([]uuid.UUID, 10)
	for i := range partUUIDs {
		partUUIDs[i] = uuid.New()
	}

	order := model.Order{
		OrderUUID:       uuid.New(),
		UserUUID:        uuid.New(),
		PartUuids:       partUUIDs,
		TotalPrice:      gofakeit.Price(1000, 10000),
		TransactionUUID: uuid.New().String(),
		PaymentMethod:   "SBP",
		Status:          "PAID",
	}

	s.orderRepository.On("GetOrder", s.ctx, order.OrderUUID.String()).Return(order, nil)

	res, err := s.service.GetOrder(s.ctx, order.OrderUUID)
	s.NoError(err)
	s.Equal(order, res)
	s.Len(res.PartUuids, 10)
}

func (s *ServiceSuite) TestGetOrderWithZeroPrice() {
	order := model.Order{
		OrderUUID:       uuid.New(),
		UserUUID:        uuid.New(),
		PartUuids:       []uuid.UUID{uuid.New()},
		TotalPrice:      0.0, // zero price
		TransactionUUID: "",
		PaymentMethod:   "",
		Status:          "PENDING_PAYMENT",
	}

	s.orderRepository.On("GetOrder", s.ctx, order.OrderUUID.String()).Return(order, nil)

	res, err := s.service.GetOrder(s.ctx, order.OrderUUID)
	s.NoError(err)
	s.Equal(order, res)
	s.Equal(0.0, res.TotalPrice)
}

func (s *ServiceSuite) TestGetOrderWithNegativePrice() {
	order := model.Order{
		OrderUUID:       uuid.New(),
		UserUUID:        uuid.New(),
		PartUuids:       []uuid.UUID{uuid.New()},
		TotalPrice:      -100.0, // negative price
		TransactionUUID: "",
		PaymentMethod:   "",
		Status:          "PENDING_PAYMENT",
	}

	s.orderRepository.On("GetOrder", s.ctx, order.OrderUUID.String()).Return(order, nil)

	res, err := s.service.GetOrder(s.ctx, order.OrderUUID)
	s.NoError(err)
	s.Equal(order, res)
	s.Equal(-100.0, res.TotalPrice)
}

func (s *ServiceSuite) TestGetOrderWithVeryHighPrice() {
	order := model.Order{
		OrderUUID:       uuid.New(),
		UserUUID:        uuid.New(),
		PartUuids:       []uuid.UUID{uuid.New()},
		TotalPrice:      999999.99, // very high price
		TransactionUUID: uuid.New().String(),
		PaymentMethod:   "CREDIT_CARD",
		Status:          "PAID",
	}

	s.orderRepository.On("GetOrder", s.ctx, order.OrderUUID.String()).Return(order, nil)

	res, err := s.service.GetOrder(s.ctx, order.OrderUUID)
	s.NoError(err)
	s.Equal(order, res)
	s.Equal(999999.99, res.TotalPrice)
}

func (s *ServiceSuite) TestGetOrderWithDifferentPaymentMethods() {
	paymentMethods := []string{"CARD", "SBP", "CREDIT_CARD", "INVESTOR_MONEY"}

	for _, method := range paymentMethods {
		order := model.Order{
			OrderUUID:       uuid.New(),
			UserUUID:        uuid.New(),
			PartUuids:       []uuid.UUID{uuid.New()},
			TotalPrice:      gofakeit.Price(100, 1000),
			TransactionUUID: uuid.New().String(),
			PaymentMethod:   method,
			Status:          "PAID",
		}

		s.orderRepository.On("GetOrder", s.ctx, order.OrderUUID.String()).Return(order, nil)

		res, err := s.service.GetOrder(s.ctx, order.OrderUUID)
		s.NoError(err)
		s.Equal(order, res)
		s.Equal(method, res.PaymentMethod)
	}
}

func (s *ServiceSuite) TestGetOrderWithEmptyTransactionUUID() {
	order := model.Order{
		OrderUUID:       uuid.New(),
		UserUUID:        uuid.New(),
		PartUuids:       []uuid.UUID{uuid.New()},
		TotalPrice:      gofakeit.Price(100, 1000),
		TransactionUUID: "", // empty transaction UUID
		PaymentMethod:   "CARD",
		Status:          "PENDING_PAYMENT",
	}

	s.orderRepository.On("GetOrder", s.ctx, order.OrderUUID.String()).Return(order, nil)

	res, err := s.service.GetOrder(s.ctx, order.OrderUUID)
	s.NoError(err)
	s.Equal(order, res)
	s.Empty(res.TransactionUUID)
}

func (s *ServiceSuite) TestGetOrderWithEmptyPaymentMethod() {
	order := model.Order{
		OrderUUID:       uuid.New(),
		UserUUID:        uuid.New(),
		PartUuids:       []uuid.UUID{uuid.New()},
		TotalPrice:      gofakeit.Price(100, 1000),
		TransactionUUID: uuid.New().String(),
		PaymentMethod:   "", // empty payment method
		Status:          "PENDING_PAYMENT",
	}

	s.orderRepository.On("GetOrder", s.ctx, order.OrderUUID.String()).Return(order, nil)

	res, err := s.service.GetOrder(s.ctx, order.OrderUUID)
	s.NoError(err)
	s.Equal(order, res)
	s.Empty(res.PaymentMethod)
}

func (s *ServiceSuite) TestGetOrderWithDifferentStatuses() {
	statuses := []string{"PENDING_PAYMENT", "PAID", "CANCELLED"}

	for _, status := range statuses {
		order := model.Order{
			OrderUUID:       uuid.New(),
			UserUUID:        uuid.New(),
			PartUuids:       []uuid.UUID{uuid.New()},
			TotalPrice:      gofakeit.Price(100, 1000),
			TransactionUUID: uuid.New().String(),
			PaymentMethod:   "CARD",
			Status:          status,
		}

		s.orderRepository.On("GetOrder", s.ctx, order.OrderUUID.String()).Return(order, nil)

		res, err := s.service.GetOrder(s.ctx, order.OrderUUID)
		s.NoError(err)
		s.Equal(order, res)
		s.Equal(status, res.Status)
	}
}

func (s *ServiceSuite) TestGetOrderWithSameUserAndOrderUUID() {
	sharedUUID := uuid.New()
	order := model.Order{
		OrderUUID:       sharedUUID,
		UserUUID:        sharedUUID, // same UUID for user and order
		PartUuids:       []uuid.UUID{uuid.New()},
		TotalPrice:      gofakeit.Price(100, 1000),
		TransactionUUID: uuid.New().String(),
		PaymentMethod:   "CARD",
		Status:          "PAID",
	}

	s.orderRepository.On("GetOrder", s.ctx, order.OrderUUID.String()).Return(order, nil)

	res, err := s.service.GetOrder(s.ctx, order.OrderUUID)
	s.NoError(err)
	s.Equal(order, res)
	s.Equal(sharedUUID, res.OrderUUID)
	s.Equal(sharedUUID, res.UserUUID)
}

func (s *ServiceSuite) TestGetOrderWithNilPartUUIDs() {
	order := model.Order{
		OrderUUID:       uuid.New(),
		UserUUID:        uuid.New(),
		PartUuids:       nil, // nil parts
		TotalPrice:      gofakeit.Price(100, 1000),
		TransactionUUID: uuid.New().String(),
		PaymentMethod:   "CARD",
		Status:          "PAID",
	}

	s.orderRepository.On("GetOrder", s.ctx, order.OrderUUID.String()).Return(order, nil)

	res, err := s.service.GetOrder(s.ctx, order.OrderUUID)
	s.NoError(err)
	s.Equal(order, res)
	s.Nil(res.PartUuids)
}
