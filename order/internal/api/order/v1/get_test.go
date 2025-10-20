package v1

import (
	"net/http"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/google/uuid"

	"github.com/nkolesnikov999/micro2-OK/order/internal/model"
	orderV1 "github.com/nkolesnikov999/micro2-OK/shared/pkg/openapi/order/v1"
)

func (s *APISuite) TestGetOrderByUuidSuccess() {
	var (
		orderUUID = uuid.MustParse(gofakeit.UUID())
		userUUID  = uuid.MustParse(gofakeit.UUID())
		partUUID1 = uuid.MustParse(gofakeit.UUID())
		partUUID2 = uuid.MustParse(gofakeit.UUID())
		order     = model.Order{
			OrderUUID:       orderUUID,
			UserUUID:        userUUID,
			PartUuids:       []uuid.UUID{partUUID1, partUUID2},
			TotalPrice:      150.50,
			TransactionUUID: gofakeit.UUID(),
			PaymentMethod:   "CARD",
			Status:          "PAID",
		}
		params = orderV1.GetOrderByUuidParams{
			OrderUUID: orderUUID,
		}
	)

	s.orderService.On("GetOrder", s.ctx, orderUUID).Return(order, nil)

	res, err := s.api.GetOrderByUuid(s.ctx, params)
	s.Require().NoError(err)
	s.Require().NotNil(res)

	orderDto, ok := res.(*orderV1.OrderDto)
	s.Require().True(ok)
	s.Require().Equal(order.OrderUUID, orderDto.OrderUUID)
	s.Require().Equal(order.UserUUID, orderDto.UserUUID)
	s.Require().Len(orderDto.PartUuids, 2)
	s.Require().Equal(partUUID1, orderDto.PartUuids[0])
	s.Require().Equal(partUUID2, orderDto.PartUuids[1])
	s.Require().Equal(float32(150.50), orderDto.TotalPrice)
	s.Require().Equal(order.TransactionUUID, orderDto.TransactionUUID.Value)
	s.Require().Equal(order.PaymentMethod, string(orderDto.PaymentMethod.Value))
	s.Require().Equal(order.Status, string(orderDto.Status))
}

func (s *APISuite) TestGetOrderByUuidNotFound() {
	var (
		orderUUID = uuid.MustParse(gofakeit.UUID())
		params    = orderV1.GetOrderByUuidParams{
			OrderUUID: orderUUID,
		}
	)

	s.orderService.On("GetOrder", s.ctx, orderUUID).Return(model.Order{}, model.ErrOrderNotFound)

	res, err := s.api.GetOrderByUuid(s.ctx, params)
	s.Require().NoError(err)
	s.Require().NotNil(res)

	notFoundErr, ok := res.(*orderV1.NotFoundError)
	s.Require().True(ok)
	s.Require().Equal(http.StatusNotFound, notFoundErr.Code)
	s.Require().Contains(notFoundErr.Message, "order not found")
}

func (s *APISuite) TestGetOrderByUuidServiceError() {
	var (
		orderUUID  = uuid.MustParse(gofakeit.UUID())
		serviceErr = gofakeit.Error()
		params     = orderV1.GetOrderByUuidParams{
			OrderUUID: orderUUID,
		}
	)

	s.orderService.On("GetOrder", s.ctx, orderUUID).Return(model.Order{}, serviceErr)

	res, err := s.api.GetOrderByUuid(s.ctx, params)
	s.Require().NoError(err)
	s.Require().NotNil(res)

	internalErr, ok := res.(*orderV1.InternalServerError)
	s.Require().True(ok)
	s.Require().Equal(http.StatusInternalServerError, internalErr.Code)
	s.Require().Equal("internal server error", internalErr.Message)
}

func (s *APISuite) TestGetOrderByUuidWithEmptyTransactionUUID() {
	var (
		orderUUID = uuid.MustParse(gofakeit.UUID())
		userUUID  = uuid.MustParse(gofakeit.UUID())
		partUUID  = uuid.MustParse(gofakeit.UUID())
		order     = model.Order{
			OrderUUID:       orderUUID,
			UserUUID:        userUUID,
			PartUuids:       []uuid.UUID{partUUID},
			TotalPrice:      99.99,
			TransactionUUID: "", // empty transaction UUID
			PaymentMethod:   "",
			Status:          "PENDING",
		}
		params = orderV1.GetOrderByUuidParams{
			OrderUUID: orderUUID,
		}
	)

	s.orderService.On("GetOrder", s.ctx, orderUUID).Return(order, nil)

	res, err := s.api.GetOrderByUuid(s.ctx, params)
	s.Require().NoError(err)
	s.Require().NotNil(res)

	orderDto, ok := res.(*orderV1.OrderDto)
	s.Require().True(ok)
	s.Require().Equal(order.OrderUUID, orderDto.OrderUUID)
	s.Require().Equal(order.UserUUID, orderDto.UserUUID)
	s.Require().Len(orderDto.PartUuids, 1)
	s.Require().Equal(partUUID, orderDto.PartUuids[0])
	s.Require().Equal(float32(99.99), orderDto.TotalPrice)
	s.Require().Equal("", orderDto.TransactionUUID.Value)
	s.Require().Equal("", string(orderDto.PaymentMethod.Value))
	s.Require().Equal(order.Status, string(orderDto.Status))
}

func (s *APISuite) TestGetOrderByUuidWithManyParts() {
	var (
		orderUUID = uuid.MustParse(gofakeit.UUID())
		userUUID  = uuid.MustParse(gofakeit.UUID())
		partUUIDs = make([]uuid.UUID, 10)
		order     = model.Order{
			OrderUUID:       orderUUID,
			UserUUID:        userUUID,
			TotalPrice:      500.75,
			TransactionUUID: gofakeit.UUID(),
			PaymentMethod:   "SBP",
			Status:          "PAID",
		}
		params = orderV1.GetOrderByUuidParams{
			OrderUUID: orderUUID,
		}
	)

	// Generate multiple part UUIDs
	for i := 0; i < 10; i++ {
		partUUIDs[i] = uuid.MustParse(gofakeit.UUID())
	}
	order.PartUuids = partUUIDs

	s.orderService.On("GetOrder", s.ctx, orderUUID).Return(order, nil)

	res, err := s.api.GetOrderByUuid(s.ctx, params)
	s.Require().NoError(err)
	s.Require().NotNil(res)

	orderDto, ok := res.(*orderV1.OrderDto)
	s.Require().True(ok)
	s.Require().Equal(order.OrderUUID, orderDto.OrderUUID)
	s.Require().Equal(order.UserUUID, orderDto.UserUUID)
	s.Require().Len(orderDto.PartUuids, 10)
	s.Require().Equal(float32(500.75), orderDto.TotalPrice)
	s.Require().Equal(order.TransactionUUID, orderDto.TransactionUUID.Value)
	s.Require().Equal(order.PaymentMethod, string(orderDto.PaymentMethod.Value))
	s.Require().Equal(order.Status, string(orderDto.Status))
}

func (s *APISuite) TestGetOrderByUuidWithZeroPrice() {
	var (
		orderUUID = uuid.MustParse(gofakeit.UUID())
		userUUID  = uuid.MustParse(gofakeit.UUID())
		partUUID  = uuid.MustParse(gofakeit.UUID())
		order     = model.Order{
			OrderUUID:       orderUUID,
			UserUUID:        userUUID,
			PartUuids:       []uuid.UUID{partUUID},
			TotalPrice:      0.0, // zero price
			TransactionUUID: gofakeit.UUID(),
			PaymentMethod:   "INVESTOR_MONEY",
			Status:          "PAID",
		}
		params = orderV1.GetOrderByUuidParams{
			OrderUUID: orderUUID,
		}
	)

	s.orderService.On("GetOrder", s.ctx, orderUUID).Return(order, nil)

	res, err := s.api.GetOrderByUuid(s.ctx, params)
	s.Require().NoError(err)
	s.Require().NotNil(res)

	orderDto, ok := res.(*orderV1.OrderDto)
	s.Require().True(ok)
	s.Require().Equal(order.OrderUUID, orderDto.OrderUUID)
	s.Require().Equal(order.UserUUID, orderDto.UserUUID)
	s.Require().Len(orderDto.PartUuids, 1)
	s.Require().Equal(float32(0.0), orderDto.TotalPrice)
	s.Require().Equal(order.TransactionUUID, orderDto.TransactionUUID.Value)
	s.Require().Equal(order.PaymentMethod, string(orderDto.PaymentMethod.Value))
	s.Require().Equal(order.Status, string(orderDto.Status))
}

func (s *APISuite) TestGetOrderByUuidWithNegativePrice() {
	var (
		orderUUID = uuid.MustParse(gofakeit.UUID())
		userUUID  = uuid.MustParse(gofakeit.UUID())
		partUUID  = uuid.MustParse(gofakeit.UUID())
		order     = model.Order{
			OrderUUID:       orderUUID,
			UserUUID:        userUUID,
			PartUuids:       []uuid.UUID{partUUID},
			TotalPrice:      -50.25, // negative price
			TransactionUUID: gofakeit.UUID(),
			PaymentMethod:   "CREDIT_CARD",
			Status:          "PAID",
		}
		params = orderV1.GetOrderByUuidParams{
			OrderUUID: orderUUID,
		}
	)

	s.orderService.On("GetOrder", s.ctx, orderUUID).Return(order, nil)

	res, err := s.api.GetOrderByUuid(s.ctx, params)
	s.Require().NoError(err)
	s.Require().NotNil(res)

	orderDto, ok := res.(*orderV1.OrderDto)
	s.Require().True(ok)
	s.Require().Equal(order.OrderUUID, orderDto.OrderUUID)
	s.Require().Equal(order.UserUUID, orderDto.UserUUID)
	s.Require().Len(orderDto.PartUuids, 1)
	s.Require().Equal(float32(-50.25), orderDto.TotalPrice)
	s.Require().Equal(order.TransactionUUID, orderDto.TransactionUUID.Value)
	s.Require().Equal(order.PaymentMethod, string(orderDto.PaymentMethod.Value))
	s.Require().Equal(order.Status, string(orderDto.Status))
}

func (s *APISuite) TestGetOrderByUuidWithDifferentStatuses() {
	statuses := []string{"PENDING", "PAID", "CANCELLED", "PROCESSING"}

	for _, status := range statuses {
		var (
			orderUUID = uuid.MustParse(gofakeit.UUID())
			userUUID  = uuid.MustParse(gofakeit.UUID())
			partUUID  = uuid.MustParse(gofakeit.UUID())
			order     = model.Order{
				OrderUUID:       orderUUID,
				UserUUID:        userUUID,
				PartUuids:       []uuid.UUID{partUUID},
				TotalPrice:      gofakeit.Price(10, 1000),
				TransactionUUID: gofakeit.UUID(),
				PaymentMethod:   "CARD",
				Status:          status,
			}
			params = orderV1.GetOrderByUuidParams{
				OrderUUID: orderUUID,
			}
		)

		s.orderService.On("GetOrder", s.ctx, orderUUID).Return(order, nil)

		res, err := s.api.GetOrderByUuid(s.ctx, params)
		s.Require().NoError(err)
		s.Require().NotNil(res)

		orderDto, ok := res.(*orderV1.OrderDto)
		s.Require().True(ok)
		s.Require().Equal(order.OrderUUID, orderDto.OrderUUID)
		s.Require().Equal(order.UserUUID, orderDto.UserUUID)
		s.Require().Equal(status, string(orderDto.Status))
	}
}

func (s *APISuite) TestGetOrderByUuidWithDifferentPaymentMethods() {
	paymentMethods := []string{"CARD", "SBP", "CREDIT_CARD", "INVESTOR_MONEY"}

	for _, paymentMethod := range paymentMethods {
		var (
			orderUUID = uuid.MustParse(gofakeit.UUID())
			userUUID  = uuid.MustParse(gofakeit.UUID())
			partUUID  = uuid.MustParse(gofakeit.UUID())
			order     = model.Order{
				OrderUUID:       orderUUID,
				UserUUID:        userUUID,
				PartUuids:       []uuid.UUID{partUUID},
				TotalPrice:      gofakeit.Price(10, 1000),
				TransactionUUID: gofakeit.UUID(),
				PaymentMethod:   paymentMethod,
				Status:          "PAID",
			}
			params = orderV1.GetOrderByUuidParams{
				OrderUUID: orderUUID,
			}
		)

		s.orderService.On("GetOrder", s.ctx, orderUUID).Return(order, nil)

		res, err := s.api.GetOrderByUuid(s.ctx, params)
		s.Require().NoError(err)
		s.Require().NotNil(res)

		orderDto, ok := res.(*orderV1.OrderDto)
		s.Require().True(ok)
		s.Require().Equal(order.OrderUUID, orderDto.OrderUUID)
		s.Require().Equal(order.UserUUID, orderDto.UserUUID)
		s.Require().Equal(paymentMethod, string(orderDto.PaymentMethod.Value))
	}
}

func (s *APISuite) TestGetOrderByUuidWithSameUserAndOrderUUID() {
	var (
		sharedUUID = uuid.MustParse(gofakeit.UUID())
		partUUID   = uuid.MustParse(gofakeit.UUID())
		order      = model.Order{
			OrderUUID:       sharedUUID,
			UserUUID:        sharedUUID, // same UUID for order and user
			PartUuids:       []uuid.UUID{partUUID},
			TotalPrice:      75.50,
			TransactionUUID: gofakeit.UUID(),
			PaymentMethod:   "CARD",
			Status:          "PAID",
		}
		params = orderV1.GetOrderByUuidParams{
			OrderUUID: sharedUUID,
		}
	)

	s.orderService.On("GetOrder", s.ctx, sharedUUID).Return(order, nil)

	res, err := s.api.GetOrderByUuid(s.ctx, params)
	s.Require().NoError(err)
	s.Require().NotNil(res)

	orderDto, ok := res.(*orderV1.OrderDto)
	s.Require().True(ok)
	s.Require().Equal(order.OrderUUID, orderDto.OrderUUID)
	s.Require().Equal(order.UserUUID, orderDto.UserUUID)
	s.Require().Equal(sharedUUID, orderDto.OrderUUID)
	s.Require().Equal(sharedUUID, orderDto.UserUUID)
	s.Require().Len(orderDto.PartUuids, 1)
	s.Require().Equal(float32(75.50), orderDto.TotalPrice)
	s.Require().Equal(order.TransactionUUID, orderDto.TransactionUUID.Value)
	s.Require().Equal(order.PaymentMethod, string(orderDto.PaymentMethod.Value))
	s.Require().Equal(order.Status, string(orderDto.Status))
}
