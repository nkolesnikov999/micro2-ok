package v1

import (
	"net/http"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/google/uuid"

	"github.com/nkolesnikov999/micro2-OK/order/internal/model"
	orderV1 "github.com/nkolesnikov999/micro2-OK/shared/pkg/openapi/order/v1"
)

func (s *APISuite) TestCreateOrderSuccess() {
	var (
		userUUID  = uuid.MustParse(gofakeit.UUID())
		partUUID1 = uuid.MustParse(gofakeit.UUID())
		partUUID2 = uuid.MustParse(gofakeit.UUID())
		req       = &orderV1.CreateOrderRequest{
			UserUUID:  userUUID,
			PartUuids: []uuid.UUID{partUUID1, partUUID2},
		}
		expectedOrder = model.Order{
			OrderUUID:  uuid.MustParse(gofakeit.UUID()),
			UserUUID:   userUUID,
			PartUuids:  []uuid.UUID{partUUID1, partUUID2},
			TotalPrice: 150.50,
			Status:     "PENDING",
		}
	)

	s.orderService.On("CreateOrder", s.ctx, userUUID, []uuid.UUID{partUUID1, partUUID2}).Return(expectedOrder, nil)

	res, err := s.api.CreateOrder(s.ctx, req)
	s.Require().NoError(err)
	s.Require().NotNil(res)

	createOrderResp, ok := res.(*orderV1.CreateOrderResponse)
	s.Require().True(ok)
	s.Require().Equal(expectedOrder.OrderUUID, createOrderResp.OrderUUID)
	s.Require().Equal(float32(150.50), createOrderResp.TotalPrice)
}

func (s *APISuite) TestCreateOrderWithSinglePart() {
	var (
		userUUID = uuid.MustParse(gofakeit.UUID())
		partUUID = uuid.MustParse(gofakeit.UUID())
		req      = &orderV1.CreateOrderRequest{
			UserUUID:  userUUID,
			PartUuids: []uuid.UUID{partUUID},
		}
		expectedOrder = model.Order{
			OrderUUID:  uuid.MustParse(gofakeit.UUID()),
			UserUUID:   userUUID,
			PartUuids:  []uuid.UUID{partUUID},
			TotalPrice: 99.99,
			Status:     "PENDING",
		}
	)

	s.orderService.On("CreateOrder", s.ctx, userUUID, []uuid.UUID{partUUID}).Return(expectedOrder, nil)

	res, err := s.api.CreateOrder(s.ctx, req)
	s.Require().NoError(err)
	s.Require().NotNil(res)

	createOrderResp, ok := res.(*orderV1.CreateOrderResponse)
	s.Require().True(ok)
	s.Require().Equal(expectedOrder.OrderUUID, createOrderResp.OrderUUID)
	s.Require().Equal(float32(99.99), createOrderResp.TotalPrice)
}

func (s *APISuite) TestCreateOrderWithManyParts() {
	var (
		userUUID  = uuid.MustParse(gofakeit.UUID())
		numParts  = gofakeit.IntRange(5, 10)
		partUUIDs = make([]uuid.UUID, numParts)
	)

	for i := 0; i < numParts; i++ {
		partUUIDs[i] = uuid.MustParse(gofakeit.UUID())
	}

	req := &orderV1.CreateOrderRequest{
		UserUUID:  userUUID,
		PartUuids: partUUIDs,
	}

	expectedOrder := model.Order{
		OrderUUID:  uuid.MustParse(gofakeit.UUID()),
		UserUUID:   userUUID,
		PartUuids:  partUUIDs,
		TotalPrice: gofakeit.Price(100, 1000),
		Status:     "PENDING",
	}

	s.orderService.On("CreateOrder", s.ctx, userUUID, partUUIDs).Return(expectedOrder, nil)

	res, err := s.api.CreateOrder(s.ctx, req)
	s.Require().NoError(err)
	s.Require().NotNil(res)

	createOrderResp, ok := res.(*orderV1.CreateOrderResponse)
	s.Require().True(ok)
	s.Require().Equal(expectedOrder.OrderUUID, createOrderResp.OrderUUID)
	s.Require().Equal(float32(expectedOrder.TotalPrice), createOrderResp.TotalPrice)
}

func (s *APISuite) TestCreateOrderNilRequest() {
	res, err := s.api.CreateOrder(s.ctx, nil)
	s.Require().NoError(err)
	s.Require().NotNil(res)

	internalErr, ok := res.(*orderV1.InternalServerError)
	s.Require().True(ok)
	s.Require().Equal(http.StatusInternalServerError, internalErr.Code)
	s.Require().Contains(internalErr.Message, "internal server error")
}

func (s *APISuite) TestCreateOrderEmptyPartUUIDs() {
	var (
		userUUID = uuid.MustParse(gofakeit.UUID())
		req      = &orderV1.CreateOrderRequest{
			UserUUID:  userUUID,
			PartUuids: []uuid.UUID{}, // Empty part UUIDs
		}
	)

	res, err := s.api.CreateOrder(s.ctx, req)
	s.Require().NoError(err)
	s.Require().NotNil(res)

	badRequestErr, ok := res.(*orderV1.BadRequestError)
	s.Require().True(ok)
	s.Require().Equal(http.StatusBadRequest, badRequestErr.Code)
	s.Require().Contains(badRequestErr.Message, "part_uuids must not be empty")
}

func (s *APISuite) TestCreateOrderServiceEmptyPartUUIDs() {
	var (
		userUUID = uuid.MustParse(gofakeit.UUID())
		partUUID = uuid.MustParse(gofakeit.UUID())
		req      = &orderV1.CreateOrderRequest{
			UserUUID:  userUUID,
			PartUuids: []uuid.UUID{partUUID},
		}
	)

	s.orderService.On("CreateOrder", s.ctx, userUUID, []uuid.UUID{partUUID}).Return(model.Order{}, model.ErrEmptyPartUUIDs)

	res, err := s.api.CreateOrder(s.ctx, req)
	s.Require().NoError(err)
	s.Require().NotNil(res)

	badRequestErr, ok := res.(*orderV1.BadRequestError)
	s.Require().True(ok)
	s.Require().Equal(http.StatusBadRequest, badRequestErr.Code)
	s.Require().Contains(badRequestErr.Message, "invalid request")
}

func (s *APISuite) TestCreateOrderPartsNotFound() {
	var (
		userUUID  = uuid.MustParse(gofakeit.UUID())
		partUUID1 = uuid.MustParse(gofakeit.UUID())
		partUUID2 = uuid.MustParse(gofakeit.UUID())
		req       = &orderV1.CreateOrderRequest{
			UserUUID:  userUUID,
			PartUuids: []uuid.UUID{partUUID1, partUUID2},
		}
	)

	s.orderService.On("CreateOrder", s.ctx, userUUID, []uuid.UUID{partUUID1, partUUID2}).Return(model.Order{}, model.ErrPartsNotFound)

	res, err := s.api.CreateOrder(s.ctx, req)
	s.Require().NoError(err)
	s.Require().NotNil(res)

	notFoundErr, ok := res.(*orderV1.NotFoundError)
	s.Require().True(ok)
	s.Require().Equal(http.StatusNotFound, notFoundErr.Code)
	s.Require().Contains(notFoundErr.Message, "parts not found")
}

func (s *APISuite) TestCreateOrderInventoryUnavailable() {
	var (
		userUUID  = uuid.MustParse(gofakeit.UUID())
		partUUID1 = uuid.MustParse(gofakeit.UUID())
		partUUID2 = uuid.MustParse(gofakeit.UUID())
		req       = &orderV1.CreateOrderRequest{
			UserUUID:  userUUID,
			PartUuids: []uuid.UUID{partUUID1, partUUID2},
		}
	)

	s.orderService.On("CreateOrder", s.ctx, userUUID, []uuid.UUID{partUUID1, partUUID2}).Return(model.Order{}, model.ErrInventoryUnavailable)

	res, err := s.api.CreateOrder(s.ctx, req)
	s.Require().NoError(err)
	s.Require().NotNil(res)

	serviceUnavailableErr, ok := res.(*orderV1.ServiceUnavailableError)
	s.Require().True(ok)
	s.Require().Equal(http.StatusServiceUnavailable, serviceUnavailableErr.Code)
	s.Require().Contains(serviceUnavailableErr.Message, "inventory service unavailable")
}

func (s *APISuite) TestCreateOrderServiceError() {
	var (
		userUUID   = uuid.MustParse(gofakeit.UUID())
		partUUID   = uuid.MustParse(gofakeit.UUID())
		serviceErr = gofakeit.Error()
		req        = &orderV1.CreateOrderRequest{
			UserUUID:  userUUID,
			PartUuids: []uuid.UUID{partUUID},
		}
	)

	s.orderService.On("CreateOrder", s.ctx, userUUID, []uuid.UUID{partUUID}).Return(model.Order{}, serviceErr)

	res, err := s.api.CreateOrder(s.ctx, req)
	s.Require().NoError(err)
	s.Require().NotNil(res)

	internalErr, ok := res.(*orderV1.InternalServerError)
	s.Require().True(ok)
	s.Require().Equal(http.StatusInternalServerError, internalErr.Code)
	s.Require().Equal("internal server error", internalErr.Message)
}

func (s *APISuite) TestCreateOrderWithZeroPrice() {
	var (
		userUUID = uuid.MustParse(gofakeit.UUID())
		partUUID = uuid.MustParse(gofakeit.UUID())
		req      = &orderV1.CreateOrderRequest{
			UserUUID:  userUUID,
			PartUuids: []uuid.UUID{partUUID},
		}
		expectedOrder = model.Order{
			OrderUUID:  uuid.MustParse(gofakeit.UUID()),
			UserUUID:   userUUID,
			PartUuids:  []uuid.UUID{partUUID},
			TotalPrice: 0.0, // Zero price
			Status:     "PENDING",
		}
	)

	s.orderService.On("CreateOrder", s.ctx, userUUID, []uuid.UUID{partUUID}).Return(expectedOrder, nil)

	res, err := s.api.CreateOrder(s.ctx, req)
	s.Require().NoError(err)
	s.Require().NotNil(res)

	createOrderResp, ok := res.(*orderV1.CreateOrderResponse)
	s.Require().True(ok)
	s.Require().Equal(expectedOrder.OrderUUID, createOrderResp.OrderUUID)
	s.Require().Equal(float32(0.0), createOrderResp.TotalPrice)
}

func (s *APISuite) TestCreateOrderWithNegativePrice() {
	var (
		userUUID = uuid.MustParse(gofakeit.UUID())
		partUUID = uuid.MustParse(gofakeit.UUID())
		req      = &orderV1.CreateOrderRequest{
			UserUUID:  userUUID,
			PartUuids: []uuid.UUID{partUUID},
		}
		expectedOrder = model.Order{
			OrderUUID:  uuid.MustParse(gofakeit.UUID()),
			UserUUID:   userUUID,
			PartUuids:  []uuid.UUID{partUUID},
			TotalPrice: -50.0, // Negative price
			Status:     "PENDING",
		}
	)

	s.orderService.On("CreateOrder", s.ctx, userUUID, []uuid.UUID{partUUID}).Return(expectedOrder, nil)

	res, err := s.api.CreateOrder(s.ctx, req)
	s.Require().NoError(err)
	s.Require().NotNil(res)

	createOrderResp, ok := res.(*orderV1.CreateOrderResponse)
	s.Require().True(ok)
	s.Require().Equal(expectedOrder.OrderUUID, createOrderResp.OrderUUID)
	s.Require().Equal(float32(-50.0), createOrderResp.TotalPrice)
}

func (s *APISuite) TestCreateOrderWithHighPrice() {
	var (
		userUUID = uuid.MustParse(gofakeit.UUID())
		partUUID = uuid.MustParse(gofakeit.UUID())
		req      = &orderV1.CreateOrderRequest{
			UserUUID:  userUUID,
			PartUuids: []uuid.UUID{partUUID},
		}
		expectedOrder = model.Order{
			OrderUUID:  uuid.MustParse(gofakeit.UUID()),
			UserUUID:   userUUID,
			PartUuids:  []uuid.UUID{partUUID},
			TotalPrice: 99999.99, // High price
			Status:     "PENDING",
		}
	)

	s.orderService.On("CreateOrder", s.ctx, userUUID, []uuid.UUID{partUUID}).Return(expectedOrder, nil)

	res, err := s.api.CreateOrder(s.ctx, req)
	s.Require().NoError(err)
	s.Require().NotNil(res)

	createOrderResp, ok := res.(*orderV1.CreateOrderResponse)
	s.Require().True(ok)
	s.Require().Equal(expectedOrder.OrderUUID, createOrderResp.OrderUUID)
	s.Require().Equal(float32(99999.99), createOrderResp.TotalPrice)
}

func (s *APISuite) TestCreateOrderWithSameUserAndPartUUIDs() {
	var (
		sharedUUID = uuid.MustParse(gofakeit.UUID())
		req        = &orderV1.CreateOrderRequest{
			UserUUID:  sharedUUID,
			PartUuids: []uuid.UUID{sharedUUID}, // Same UUID for user and part
		}
		expectedOrder = model.Order{
			OrderUUID:  uuid.MustParse(gofakeit.UUID()),
			UserUUID:   sharedUUID,
			PartUuids:  []uuid.UUID{sharedUUID},
			TotalPrice: 75.50,
			Status:     "PENDING",
		}
	)

	s.orderService.On("CreateOrder", s.ctx, sharedUUID, []uuid.UUID{sharedUUID}).Return(expectedOrder, nil)

	res, err := s.api.CreateOrder(s.ctx, req)
	s.Require().NoError(err)
	s.Require().NotNil(res)

	createOrderResp, ok := res.(*orderV1.CreateOrderResponse)
	s.Require().True(ok)
	s.Require().Equal(expectedOrder.OrderUUID, createOrderResp.OrderUUID)
	s.Require().Equal(float32(75.50), createOrderResp.TotalPrice)
}

func (s *APISuite) TestCreateOrderWithDuplicatePartUUIDs() {
	var (
		userUUID = uuid.MustParse(gofakeit.UUID())
		partUUID = uuid.MustParse(gofakeit.UUID())
		req      = &orderV1.CreateOrderRequest{
			UserUUID:  userUUID,
			PartUuids: []uuid.UUID{partUUID, partUUID}, // Duplicate part UUIDs
		}
		expectedOrder = model.Order{
			OrderUUID:  uuid.MustParse(gofakeit.UUID()),
			UserUUID:   userUUID,
			PartUuids:  []uuid.UUID{partUUID, partUUID},
			TotalPrice: 150.0,
			Status:     "PENDING",
		}
	)

	s.orderService.On("CreateOrder", s.ctx, userUUID, []uuid.UUID{partUUID, partUUID}).Return(expectedOrder, nil)

	res, err := s.api.CreateOrder(s.ctx, req)
	s.Require().NoError(err)
	s.Require().NotNil(res)

	createOrderResp, ok := res.(*orderV1.CreateOrderResponse)
	s.Require().True(ok)
	s.Require().Equal(expectedOrder.OrderUUID, createOrderResp.OrderUUID)
	s.Require().Equal(float32(150.0), createOrderResp.TotalPrice)
}
