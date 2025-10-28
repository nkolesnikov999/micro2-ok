package order

import (
	"github.com/brianvoe/gofakeit/v7"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"

	"github.com/nkolesnikov999/micro2-OK/order/internal/model"
)

// Helper function to create a matcher for updated order
func (s *ServiceSuite) createUpdatedOrderMatcher(originalOrder model.Order) interface{} {
	return mock.MatchedBy(func(updatedOrder model.Order) bool {
		return updatedOrder.Status == "CANCELLED" &&
			updatedOrder.OrderUUID == originalOrder.OrderUUID &&
			updatedOrder.UserUUID == originalOrder.UserUUID &&
			updatedOrder.TotalPrice == originalOrder.TotalPrice &&
			updatedOrder.TransactionUUID == originalOrder.TransactionUUID &&
			updatedOrder.PaymentMethod == originalOrder.PaymentMethod &&
			len(updatedOrder.PartUuids) == len(originalOrder.PartUuids) &&
			!updatedOrder.UpdatedAt.IsZero() // Проверяем, что UpdatedAt установлен
	})
}

func (s *ServiceSuite) TestCancelOrderSuccess() {
	order := model.Order{
		OrderUUID:       uuid.New(),
		UserUUID:        uuid.New(),
		PartUuids:       []uuid.UUID{uuid.New()},
		TotalPrice:      gofakeit.Price(100, 1000),
		TransactionUUID: "",
		PaymentMethod:   "",
		Status:          "PENDING_PAYMENT",
	}

	s.orderRepository.On("GetOrder", s.ctx, order.OrderUUID).Return(order, nil)
	s.orderRepository.On("UpdateOrder", s.ctx, order.OrderUUID, s.createUpdatedOrderMatcher(order)).Return(nil)

	err := s.service.CancelOrder(s.ctx, order.OrderUUID)
	s.NoError(err)
}

func (s *ServiceSuite) TestCancelOrderNotFound() {
	orderUUID := uuid.New()

	s.orderRepository.On("GetOrder", s.ctx, orderUUID).Return(model.Order{}, model.ErrOrderNotFound)

	err := s.service.CancelOrder(s.ctx, orderUUID)
	s.Error(err)
	s.ErrorIs(err, model.ErrOrderNotFound)
}

func (s *ServiceSuite) TestCancelOrderGetFailed() {
	orderUUID := uuid.New()
	repoErr := gofakeit.Error()

	s.orderRepository.On("GetOrder", s.ctx, orderUUID).Return(model.Order{}, repoErr)

	err := s.service.CancelOrder(s.ctx, orderUUID)
	s.Error(err)
	s.ErrorIs(err, model.ErrOrderGetFailed)
}

func (s *ServiceSuite) TestCancelOrderAlreadyPaid() {
	order := model.Order{
		OrderUUID:       uuid.New(),
		UserUUID:        uuid.New(),
		PartUuids:       []uuid.UUID{uuid.New()},
		TotalPrice:      gofakeit.Price(100, 1000),
		TransactionUUID: uuid.New().String(),
		PaymentMethod:   "CARD",
		Status:          "PAID", // already paid
	}

	s.orderRepository.On("GetOrder", s.ctx, order.OrderUUID).Return(order, nil)

	err := s.service.CancelOrder(s.ctx, order.OrderUUID)
	s.Error(err)
	s.ErrorIs(err, model.ErrCannotCancelPaidOrder)
}

func (s *ServiceSuite) TestCancelOrderAlreadyCancelled() {
	order := model.Order{
		OrderUUID:       uuid.New(),
		UserUUID:        uuid.New(),
		PartUuids:       []uuid.UUID{uuid.New()},
		TotalPrice:      gofakeit.Price(100, 1000),
		TransactionUUID: "",
		PaymentMethod:   "",
		Status:          "CANCELLED", // already cancelled
	}

	s.orderRepository.On("GetOrder", s.ctx, order.OrderUUID).Return(order, nil)

	err := s.service.CancelOrder(s.ctx, order.OrderUUID)
	s.NoError(err) // Should succeed without updating
}

func (s *ServiceSuite) TestCancelOrderUpdateFailed() {
	order := model.Order{
		OrderUUID:       uuid.New(),
		UserUUID:        uuid.New(),
		PartUuids:       []uuid.UUID{uuid.New()},
		TotalPrice:      gofakeit.Price(100, 1000),
		TransactionUUID: "",
		PaymentMethod:   "",
		Status:          "PENDING_PAYMENT",
	}
	updateErr := gofakeit.Error()

	s.orderRepository.On("GetOrder", s.ctx, order.OrderUUID).Return(order, nil)
	s.orderRepository.On("UpdateOrder", s.ctx, order.OrderUUID, s.createUpdatedOrderMatcher(order)).Return(updateErr)

	err := s.service.CancelOrder(s.ctx, order.OrderUUID)
	s.Error(err)
	s.ErrorIs(err, model.ErrOrderUpdateFailed)
}

func (s *ServiceSuite) TestCancelOrderUpdateNotFound() {
	order := model.Order{
		OrderUUID:       uuid.New(),
		UserUUID:        uuid.New(),
		PartUuids:       []uuid.UUID{uuid.New()},
		TotalPrice:      gofakeit.Price(100, 1000),
		TransactionUUID: "",
		PaymentMethod:   "",
		Status:          "PENDING_PAYMENT",
	}

	s.orderRepository.On("GetOrder", s.ctx, order.OrderUUID).Return(order, nil)
	s.orderRepository.On("UpdateOrder", s.ctx, order.OrderUUID, s.createUpdatedOrderMatcher(order)).Return(model.ErrOrderNotFound)

	err := s.service.CancelOrder(s.ctx, order.OrderUUID)
	s.Error(err)
	s.ErrorIs(err, model.ErrOrderNotFound)
}

func (s *ServiceSuite) TestCancelOrderWithDifferentStatuses() {
	statuses := []string{"PENDING_PAYMENT", "CANCELLED"}

	for _, status := range statuses {
		order := model.Order{
			OrderUUID:       uuid.New(),
			UserUUID:        uuid.New(),
			PartUuids:       []uuid.UUID{uuid.New()},
			TotalPrice:      gofakeit.Price(100, 1000),
			TransactionUUID: "",
			PaymentMethod:   "",
			Status:          status,
		}

		if status == "PENDING_PAYMENT" {
			s.orderRepository.On("GetOrder", s.ctx, order.OrderUUID).Return(order, nil)
			s.orderRepository.On("UpdateOrder", s.ctx, order.OrderUUID, s.createUpdatedOrderMatcher(order)).Return(nil)
		} else {
			s.orderRepository.On("GetOrder", s.ctx, order.OrderUUID).Return(order, nil)
		}

		err := s.service.CancelOrder(s.ctx, order.OrderUUID)
		s.NoError(err)
	}
}

func (s *ServiceSuite) TestCancelOrderWithEmptyPartUUIDs() {
	order := model.Order{
		OrderUUID:       uuid.New(),
		UserUUID:        uuid.New(),
		PartUuids:       []uuid.UUID{}, // empty parts
		TotalPrice:      gofakeit.Price(100, 1000),
		TransactionUUID: "",
		PaymentMethod:   "",
		Status:          "PENDING_PAYMENT",
	}

	s.orderRepository.On("GetOrder", s.ctx, order.OrderUUID).Return(order, nil)
	s.orderRepository.On("UpdateOrder", s.ctx, order.OrderUUID, s.createUpdatedOrderMatcher(order)).Return(nil)

	err := s.service.CancelOrder(s.ctx, order.OrderUUID)
	s.NoError(err)
}

func (s *ServiceSuite) TestCancelOrderWithNilPartUUIDs() {
	order := model.Order{
		OrderUUID:       uuid.New(),
		UserUUID:        uuid.New(),
		PartUuids:       nil, // nil parts
		TotalPrice:      gofakeit.Price(100, 1000),
		TransactionUUID: "",
		PaymentMethod:   "",
		Status:          "PENDING_PAYMENT",
	}

	s.orderRepository.On("GetOrder", s.ctx, order.OrderUUID).Return(order, nil)
	s.orderRepository.On("UpdateOrder", s.ctx, order.OrderUUID, s.createUpdatedOrderMatcher(order)).Return(nil)

	err := s.service.CancelOrder(s.ctx, order.OrderUUID)
	s.NoError(err)
}

func (s *ServiceSuite) TestCancelOrderWithManyParts() {
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

	s.orderRepository.On("GetOrder", s.ctx, order.OrderUUID).Return(order, nil)
	s.orderRepository.On("UpdateOrder", s.ctx, order.OrderUUID, s.createUpdatedOrderMatcher(order)).Return(nil)

	err := s.service.CancelOrder(s.ctx, order.OrderUUID)
	s.NoError(err)
}

func (s *ServiceSuite) TestCancelOrderWithZeroPrice() {
	order := model.Order{
		OrderUUID:       uuid.New(),
		UserUUID:        uuid.New(),
		PartUuids:       []uuid.UUID{uuid.New()},
		TotalPrice:      0.0, // zero price
		TransactionUUID: "",
		PaymentMethod:   "",
		Status:          "PENDING_PAYMENT",
	}

	s.orderRepository.On("GetOrder", s.ctx, order.OrderUUID).Return(order, nil)
	s.orderRepository.On("UpdateOrder", s.ctx, order.OrderUUID, s.createUpdatedOrderMatcher(order)).Return(nil)

	err := s.service.CancelOrder(s.ctx, order.OrderUUID)
	s.NoError(err)
}

func (s *ServiceSuite) TestCancelOrderWithNegativePrice() {
	order := model.Order{
		OrderUUID:       uuid.New(),
		UserUUID:        uuid.New(),
		PartUuids:       []uuid.UUID{uuid.New()},
		TotalPrice:      -100.0, // negative price
		TransactionUUID: "",
		PaymentMethod:   "",
		Status:          "PENDING_PAYMENT",
	}

	s.orderRepository.On("GetOrder", s.ctx, order.OrderUUID).Return(order, nil)
	s.orderRepository.On("UpdateOrder", s.ctx, order.OrderUUID, s.createUpdatedOrderMatcher(order)).Return(nil)

	err := s.service.CancelOrder(s.ctx, order.OrderUUID)
	s.NoError(err)
}

func (s *ServiceSuite) TestCancelOrderWithVeryHighPrice() {
	order := model.Order{
		OrderUUID:       uuid.New(),
		UserUUID:        uuid.New(),
		PartUuids:       []uuid.UUID{uuid.New()},
		TotalPrice:      999999.99, // very high price
		TransactionUUID: "",
		PaymentMethod:   "",
		Status:          "PENDING_PAYMENT",
	}

	s.orderRepository.On("GetOrder", s.ctx, order.OrderUUID).Return(order, nil)
	s.orderRepository.On("UpdateOrder", s.ctx, order.OrderUUID, s.createUpdatedOrderMatcher(order)).Return(nil)

	err := s.service.CancelOrder(s.ctx, order.OrderUUID)
	s.NoError(err)
}

func (s *ServiceSuite) TestCancelOrderWithSameUserAndOrderUUID() {
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

	s.orderRepository.On("GetOrder", s.ctx, order.OrderUUID).Return(order, nil)
	s.orderRepository.On("UpdateOrder", s.ctx, order.OrderUUID, s.createUpdatedOrderMatcher(order)).Return(nil)

	err := s.service.CancelOrder(s.ctx, order.OrderUUID)
	s.NoError(err)
}

func (s *ServiceSuite) TestCancelOrderWithEmptyTransactionUUID() {
	order := model.Order{
		OrderUUID:       uuid.New(),
		UserUUID:        uuid.New(),
		PartUuids:       []uuid.UUID{uuid.New()},
		TotalPrice:      gofakeit.Price(100, 1000),
		TransactionUUID: "", // empty transaction UUID
		PaymentMethod:   "",
		Status:          "PENDING_PAYMENT",
	}

	s.orderRepository.On("GetOrder", s.ctx, order.OrderUUID).Return(order, nil)
	s.orderRepository.On("UpdateOrder", s.ctx, order.OrderUUID, s.createUpdatedOrderMatcher(order)).Return(nil)

	err := s.service.CancelOrder(s.ctx, order.OrderUUID)
	s.NoError(err)
}

func (s *ServiceSuite) TestCancelOrderWithEmptyPaymentMethod() {
	order := model.Order{
		OrderUUID:       uuid.New(),
		UserUUID:        uuid.New(),
		PartUuids:       []uuid.UUID{uuid.New()},
		TotalPrice:      gofakeit.Price(100, 1000),
		TransactionUUID: "",
		PaymentMethod:   "", // empty payment method
		Status:          "PENDING_PAYMENT",
	}

	s.orderRepository.On("GetOrder", s.ctx, order.OrderUUID).Return(order, nil)
	s.orderRepository.On("UpdateOrder", s.ctx, order.OrderUUID, s.createUpdatedOrderMatcher(order)).Return(nil)

	err := s.service.CancelOrder(s.ctx, order.OrderUUID)
	s.NoError(err)
}
