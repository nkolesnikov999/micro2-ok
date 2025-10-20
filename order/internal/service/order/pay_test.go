package order

import (
	"github.com/brianvoe/gofakeit/v7"
	"github.com/google/uuid"

	"github.com/nkolesnikov999/micro2-OK/order/internal/model"
)

// createExpectedOrder creates an order with updated payment fields
func createExpectedOrder(order model.Order, transactionUUID, paymentMethod string) model.Order {
	expectedOrder := order
	expectedOrder.Status = "PAID"
	expectedOrder.TransactionUUID = transactionUUID
	expectedOrder.PaymentMethod = paymentMethod
	return expectedOrder
}

func (s *ServiceSuite) TestPayOrderSuccess() {
	order := model.Order{
		OrderUUID:       uuid.New(),
		UserUUID:        uuid.New(),
		PartUuids:       []uuid.UUID{uuid.New()},
		TotalPrice:      gofakeit.Price(100, 1000),
		TransactionUUID: "",
		PaymentMethod:   "",
		Status:          "PENDING_PAYMENT",
	}
	paymentMethod := "CARD"
	transactionUUID := uuid.New().String()

	// Create expected order after payment
	expectedOrder := order
	expectedOrder.Status = "PAID"
	expectedOrder.TransactionUUID = transactionUUID
	expectedOrder.PaymentMethod = paymentMethod

	s.orderRepository.On("GetOrder", s.ctx, order.OrderUUID).Return(order, nil)
	s.paymentClient.On("PayOrder", s.ctx, order.OrderUUID.String(), order.UserUUID.String(), paymentMethod).Return(transactionUUID, nil)
	s.orderRepository.On("UpdateOrder", s.ctx, order.OrderUUID, expectedOrder).Return(nil)

	res, err := s.service.PayOrder(s.ctx, order.OrderUUID, paymentMethod)
	s.NoError(err)
	s.Equal(transactionUUID, res)
}

func (s *ServiceSuite) TestPayOrderNotFound() {
	orderUUID := uuid.New()
	paymentMethod := "CARD"

	s.orderRepository.On("GetOrder", s.ctx, orderUUID).Return(model.Order{}, model.ErrOrderNotFound)

	res, err := s.service.PayOrder(s.ctx, orderUUID, paymentMethod)
	s.Error(err)
	s.ErrorIs(err, model.ErrOrderNotFound)
	s.Empty(res)
}

func (s *ServiceSuite) TestPayOrderGetFailed() {
	orderUUID := uuid.New()
	paymentMethod := "CARD"
	repoErr := gofakeit.Error()

	s.orderRepository.On("GetOrder", s.ctx, orderUUID).Return(model.Order{}, repoErr)

	res, err := s.service.PayOrder(s.ctx, orderUUID, paymentMethod)
	s.Error(err)
	s.ErrorIs(err, model.ErrOrderGetFailed)
	s.Empty(res)
}

func (s *ServiceSuite) TestPayOrderAlreadyPaid() {
	order := model.Order{
		OrderUUID:       uuid.New(),
		UserUUID:        uuid.New(),
		PartUuids:       []uuid.UUID{uuid.New()},
		TotalPrice:      gofakeit.Price(100, 1000),
		TransactionUUID: uuid.New().String(),
		PaymentMethod:   "CARD",
		Status:          "PAID", // already paid
	}
	paymentMethod := "CARD"

	s.orderRepository.On("GetOrder", s.ctx, order.OrderUUID).Return(order, nil)

	res, err := s.service.PayOrder(s.ctx, order.OrderUUID, paymentMethod)
	s.Error(err)
	s.ErrorIs(err, model.ErrOrderNotPayable)
	s.Empty(res)
}

func (s *ServiceSuite) TestPayOrderCancelled() {
	order := model.Order{
		OrderUUID:       uuid.New(),
		UserUUID:        uuid.New(),
		PartUuids:       []uuid.UUID{uuid.New()},
		TotalPrice:      gofakeit.Price(100, 1000),
		TransactionUUID: "",
		PaymentMethod:   "",
		Status:          "CANCELLED", // cancelled order
	}
	paymentMethod := "CARD"

	s.orderRepository.On("GetOrder", s.ctx, order.OrderUUID).Return(order, nil)

	res, err := s.service.PayOrder(s.ctx, order.OrderUUID, paymentMethod)
	s.Error(err)
	s.ErrorIs(err, model.ErrOrderNotPayable)
	s.Empty(res)
}

func (s *ServiceSuite) TestPayOrderPaymentFailed() {
	order := model.Order{
		OrderUUID:       uuid.New(),
		UserUUID:        uuid.New(),
		PartUuids:       []uuid.UUID{uuid.New()},
		TotalPrice:      gofakeit.Price(100, 1000),
		TransactionUUID: "",
		PaymentMethod:   "",
		Status:          "PENDING_PAYMENT",
	}
	paymentMethod := "CARD"
	paymentErr := gofakeit.Error()

	s.orderRepository.On("GetOrder", s.ctx, order.OrderUUID).Return(order, nil)
	s.paymentClient.On("PayOrder", s.ctx, order.OrderUUID.String(), order.UserUUID.String(), paymentMethod).Return("", paymentErr)

	res, err := s.service.PayOrder(s.ctx, order.OrderUUID, paymentMethod)
	s.Error(err)
	s.ErrorIs(err, model.ErrPaymentFailed)
	s.Empty(res)
}

func (s *ServiceSuite) TestPayOrderUpdateFailed() {
	order := model.Order{
		OrderUUID:       uuid.New(),
		UserUUID:        uuid.New(),
		PartUuids:       []uuid.UUID{uuid.New()},
		TotalPrice:      gofakeit.Price(100, 1000),
		TransactionUUID: "",
		PaymentMethod:   "",
		Status:          "PENDING_PAYMENT",
	}
	paymentMethod := "CARD"
	transactionUUID := uuid.New().String()
	updateErr := gofakeit.Error()

	s.orderRepository.On("GetOrder", s.ctx, order.OrderUUID).Return(order, nil)
	s.paymentClient.On("PayOrder", s.ctx, order.OrderUUID.String(), order.UserUUID.String(), paymentMethod).Return(transactionUUID, nil)
	s.orderRepository.On("UpdateOrder", s.ctx, order.OrderUUID, createExpectedOrder(order, transactionUUID, paymentMethod)).Return(updateErr)

	res, err := s.service.PayOrder(s.ctx, order.OrderUUID, paymentMethod)
	s.Error(err)
	s.ErrorIs(err, model.ErrOrderUpdateFailed)
	s.Empty(res)
}

func (s *ServiceSuite) TestPayOrderUpdateNotFound() {
	order := model.Order{
		OrderUUID:       uuid.New(),
		UserUUID:        uuid.New(),
		PartUuids:       []uuid.UUID{uuid.New()},
		TotalPrice:      gofakeit.Price(100, 1000),
		TransactionUUID: "",
		PaymentMethod:   "",
		Status:          "PENDING_PAYMENT",
	}
	paymentMethod := "CARD"
	transactionUUID := uuid.New().String()

	s.orderRepository.On("GetOrder", s.ctx, order.OrderUUID).Return(order, nil)
	s.paymentClient.On("PayOrder", s.ctx, order.OrderUUID.String(), order.UserUUID.String(), paymentMethod).Return(transactionUUID, nil)
	s.orderRepository.On("UpdateOrder", s.ctx, order.OrderUUID, createExpectedOrder(order, transactionUUID, paymentMethod)).Return(model.ErrOrderNotFound)

	res, err := s.service.PayOrder(s.ctx, order.OrderUUID, paymentMethod)
	s.Error(err)
	s.ErrorIs(err, model.ErrOrderNotFound)
	s.Empty(res)
}

func (s *ServiceSuite) TestPayOrderWithDifferentPaymentMethods() {
	paymentMethods := []string{"CARD", "SBP", "CREDIT_CARD", "INVESTOR_MONEY"}

	for _, method := range paymentMethods {
		order := model.Order{
			OrderUUID:       uuid.New(),
			UserUUID:        uuid.New(),
			PartUuids:       []uuid.UUID{uuid.New()},
			TotalPrice:      gofakeit.Price(100, 1000),
			TransactionUUID: "",
			PaymentMethod:   "",
			Status:          "PENDING_PAYMENT",
		}
		transactionUUID := uuid.New().String()

		s.orderRepository.On("GetOrder", s.ctx, order.OrderUUID).Return(order, nil)
		s.paymentClient.On("PayOrder", s.ctx, order.OrderUUID.String(), order.UserUUID.String(), method).Return(transactionUUID, nil)
		s.orderRepository.On("UpdateOrder", s.ctx, order.OrderUUID, createExpectedOrder(order, transactionUUID, method)).Return(nil)

		res, err := s.service.PayOrder(s.ctx, order.OrderUUID, method)
		s.NoError(err)
		s.Equal(transactionUUID, res)
	}
}

func (s *ServiceSuite) TestPayOrderWithEmptyPaymentMethod() {
	order := model.Order{
		OrderUUID:       uuid.New(),
		UserUUID:        uuid.New(),
		PartUuids:       []uuid.UUID{uuid.New()},
		TotalPrice:      gofakeit.Price(100, 1000),
		TransactionUUID: "",
		PaymentMethod:   "",
		Status:          "PENDING_PAYMENT",
	}
	paymentMethod := "" // empty payment method
	transactionUUID := uuid.New().String()

	s.orderRepository.On("GetOrder", s.ctx, order.OrderUUID).Return(order, nil)
	s.paymentClient.On("PayOrder", s.ctx, order.OrderUUID.String(), order.UserUUID.String(), paymentMethod).Return(transactionUUID, nil)
	s.orderRepository.On("UpdateOrder", s.ctx, order.OrderUUID, createExpectedOrder(order, transactionUUID, paymentMethod)).Return(nil)

	res, err := s.service.PayOrder(s.ctx, order.OrderUUID, paymentMethod)
	s.NoError(err)
	s.Equal(transactionUUID, res)
}

func (s *ServiceSuite) TestPayOrderWithZeroPrice() {
	order := model.Order{
		OrderUUID:       uuid.New(),
		UserUUID:        uuid.New(),
		PartUuids:       []uuid.UUID{uuid.New()},
		TotalPrice:      0.0, // zero price
		TransactionUUID: "",
		PaymentMethod:   "",
		Status:          "PENDING_PAYMENT",
	}
	paymentMethod := "CARD"
	transactionUUID := uuid.New().String()

	s.orderRepository.On("GetOrder", s.ctx, order.OrderUUID).Return(order, nil)
	s.paymentClient.On("PayOrder", s.ctx, order.OrderUUID.String(), order.UserUUID.String(), paymentMethod).Return(transactionUUID, nil)
	s.orderRepository.On("UpdateOrder", s.ctx, order.OrderUUID, createExpectedOrder(order, transactionUUID, paymentMethod)).Return(nil)

	res, err := s.service.PayOrder(s.ctx, order.OrderUUID, paymentMethod)
	s.NoError(err)
	s.Equal(transactionUUID, res)
}

func (s *ServiceSuite) TestPayOrderWithNegativePrice() {
	order := model.Order{
		OrderUUID:       uuid.New(),
		UserUUID:        uuid.New(),
		PartUuids:       []uuid.UUID{uuid.New()},
		TotalPrice:      -100.0, // negative price
		TransactionUUID: "",
		PaymentMethod:   "",
		Status:          "PENDING_PAYMENT",
	}
	paymentMethod := "CARD"
	transactionUUID := uuid.New().String()

	s.orderRepository.On("GetOrder", s.ctx, order.OrderUUID).Return(order, nil)
	s.paymentClient.On("PayOrder", s.ctx, order.OrderUUID.String(), order.UserUUID.String(), paymentMethod).Return(transactionUUID, nil)
	s.orderRepository.On("UpdateOrder", s.ctx, order.OrderUUID, createExpectedOrder(order, transactionUUID, paymentMethod)).Return(nil)

	res, err := s.service.PayOrder(s.ctx, order.OrderUUID, paymentMethod)
	s.NoError(err)
	s.Equal(transactionUUID, res)
}

func (s *ServiceSuite) TestPayOrderWithVeryHighPrice() {
	order := model.Order{
		OrderUUID:       uuid.New(),
		UserUUID:        uuid.New(),
		PartUuids:       []uuid.UUID{uuid.New()},
		TotalPrice:      999999.99, // very high price
		TransactionUUID: "",
		PaymentMethod:   "",
		Status:          "PENDING_PAYMENT",
	}
	paymentMethod := "CARD"
	transactionUUID := uuid.New().String()

	s.orderRepository.On("GetOrder", s.ctx, order.OrderUUID).Return(order, nil)
	s.paymentClient.On("PayOrder", s.ctx, order.OrderUUID.String(), order.UserUUID.String(), paymentMethod).Return(transactionUUID, nil)
	s.orderRepository.On("UpdateOrder", s.ctx, order.OrderUUID, createExpectedOrder(order, transactionUUID, paymentMethod)).Return(nil)

	res, err := s.service.PayOrder(s.ctx, order.OrderUUID, paymentMethod)
	s.NoError(err)
	s.Equal(transactionUUID, res)
}

func (s *ServiceSuite) TestPayOrderWithEmptyPartUUIDs() {
	order := model.Order{
		OrderUUID:       uuid.New(),
		UserUUID:        uuid.New(),
		PartUuids:       []uuid.UUID{}, // empty parts
		TotalPrice:      gofakeit.Price(100, 1000),
		TransactionUUID: "",
		PaymentMethod:   "",
		Status:          "PENDING_PAYMENT",
	}
	paymentMethod := "CARD"
	transactionUUID := uuid.New().String()

	s.orderRepository.On("GetOrder", s.ctx, order.OrderUUID).Return(order, nil)
	s.paymentClient.On("PayOrder", s.ctx, order.OrderUUID.String(), order.UserUUID.String(), paymentMethod).Return(transactionUUID, nil)
	s.orderRepository.On("UpdateOrder", s.ctx, order.OrderUUID, createExpectedOrder(order, transactionUUID, paymentMethod)).Return(nil)

	res, err := s.service.PayOrder(s.ctx, order.OrderUUID, paymentMethod)
	s.NoError(err)
	s.Equal(transactionUUID, res)
}

func (s *ServiceSuite) TestPayOrderWithNilPartUUIDs() {
	order := model.Order{
		OrderUUID:       uuid.New(),
		UserUUID:        uuid.New(),
		PartUuids:       nil, // nil parts
		TotalPrice:      gofakeit.Price(100, 1000),
		TransactionUUID: "",
		PaymentMethod:   "",
		Status:          "PENDING_PAYMENT",
	}
	paymentMethod := "CARD"
	transactionUUID := uuid.New().String()

	s.orderRepository.On("GetOrder", s.ctx, order.OrderUUID).Return(order, nil)
	s.paymentClient.On("PayOrder", s.ctx, order.OrderUUID.String(), order.UserUUID.String(), paymentMethod).Return(transactionUUID, nil)
	s.orderRepository.On("UpdateOrder", s.ctx, order.OrderUUID, createExpectedOrder(order, transactionUUID, paymentMethod)).Return(nil)

	res, err := s.service.PayOrder(s.ctx, order.OrderUUID, paymentMethod)
	s.NoError(err)
	s.Equal(transactionUUID, res)
}

func (s *ServiceSuite) TestPayOrderWithManyParts() {
	partUUIDs := make([]uuid.UUID, 10)
	for i := range partUUIDs {
		partUUIDs[i] = uuid.New()
	}

	order := model.Order{
		OrderUUID:       uuid.New(),
		UserUUID:        uuid.New(),
		PartUuids:       partUUIDs,
		TotalPrice:      gofakeit.Price(1000, 10000),
		TransactionUUID: "",
		PaymentMethod:   "",
		Status:          "PENDING_PAYMENT",
	}
	paymentMethod := "CARD"
	transactionUUID := uuid.New().String()

	s.orderRepository.On("GetOrder", s.ctx, order.OrderUUID).Return(order, nil)
	s.paymentClient.On("PayOrder", s.ctx, order.OrderUUID.String(), order.UserUUID.String(), paymentMethod).Return(transactionUUID, nil)
	s.orderRepository.On("UpdateOrder", s.ctx, order.OrderUUID, createExpectedOrder(order, transactionUUID, paymentMethod)).Return(nil)

	res, err := s.service.PayOrder(s.ctx, order.OrderUUID, paymentMethod)
	s.NoError(err)
	s.Equal(transactionUUID, res)
}

func (s *ServiceSuite) TestPayOrderWithSameUserAndOrderUUID() {
	sharedUUID := uuid.New()
	order := model.Order{
		OrderUUID:       sharedUUID,
		UserUUID:        sharedUUID, // same UUID for user and order
		PartUuids:       []uuid.UUID{uuid.New()},
		TotalPrice:      gofakeit.Price(100, 1000),
		TransactionUUID: "",
		PaymentMethod:   "",
		Status:          "PENDING_PAYMENT",
	}
	paymentMethod := "CARD"
	transactionUUID := uuid.New().String()

	s.orderRepository.On("GetOrder", s.ctx, order.OrderUUID).Return(order, nil)
	s.paymentClient.On("PayOrder", s.ctx, order.OrderUUID.String(), order.UserUUID.String(), paymentMethod).Return(transactionUUID, nil)
	s.orderRepository.On("UpdateOrder", s.ctx, order.OrderUUID, createExpectedOrder(order, transactionUUID, paymentMethod)).Return(nil)

	res, err := s.service.PayOrder(s.ctx, order.OrderUUID, paymentMethod)
	s.NoError(err)
	s.Equal(transactionUUID, res)
}

func (s *ServiceSuite) TestPayOrderUpdatesOrderStatus() {
	order := model.Order{
		OrderUUID:       uuid.New(),
		UserUUID:        uuid.New(),
		PartUuids:       []uuid.UUID{uuid.New()},
		TotalPrice:      gofakeit.Price(100, 1000),
		TransactionUUID: "",
		PaymentMethod:   "",
		Status:          "PENDING_PAYMENT",
	}
	paymentMethod := "CARD"
	transactionUUID := uuid.New().String()

	// Create expected order after payment
	expectedOrder := order
	expectedOrder.Status = "PAID"
	expectedOrder.TransactionUUID = transactionUUID
	expectedOrder.PaymentMethod = paymentMethod

	s.orderRepository.On("GetOrder", s.ctx, order.OrderUUID).Return(order, nil)
	s.paymentClient.On("PayOrder", s.ctx, order.OrderUUID.String(), order.UserUUID.String(), paymentMethod).Return(transactionUUID, nil)
	s.orderRepository.On("UpdateOrder", s.ctx, order.OrderUUID, expectedOrder).Return(nil)

	res, err := s.service.PayOrder(s.ctx, order.OrderUUID, paymentMethod)
	s.NoError(err)
	s.Equal(transactionUUID, res)
}
