package v1

import (
	"net/http"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/google/uuid"

	"github.com/nkolesnikov999/micro2-OK/order/internal/model"
	orderV1 "github.com/nkolesnikov999/micro2-OK/shared/pkg/openapi/order/v1"
)

func (s *APISuite) TestCancelOrderSuccess() {
	var (
		orderUUID = uuid.MustParse(gofakeit.UUID())
		params    = orderV1.CancelOrderParams{
			OrderUUID: orderUUID,
		}
	)

	s.orderService.On("CancelOrder", s.ctx, orderUUID).Return(nil)

	res, err := s.api.CancelOrder(s.ctx, params)
	s.Require().Error(err)
	s.Require().Nil(res)

	var statusCode *orderV1.GenericErrorStatusCode
	s.Require().ErrorAs(err, &statusCode)
	s.Require().Equal(http.StatusNoContent, statusCode.StatusCode)
}

func (s *APISuite) TestCancelOrderNotFound() {
	var (
		orderUUID = uuid.MustParse(gofakeit.UUID())
		params    = orderV1.CancelOrderParams{
			OrderUUID: orderUUID,
		}
	)

	s.orderService.On("CancelOrder", s.ctx, orderUUID).Return(model.ErrOrderNotFound)

	res, err := s.api.CancelOrder(s.ctx, params)
	s.Require().NoError(err)
	s.Require().NotNil(res)

	notFoundErr, ok := res.(*orderV1.NotFoundError)
	s.Require().True(ok)
	s.Require().Equal(http.StatusNotFound, notFoundErr.Code)
	s.Require().Contains(notFoundErr.Message, "order not found")
}

func (s *APISuite) TestCancelOrderCannotCancelPaidOrder() {
	var (
		orderUUID = uuid.MustParse(gofakeit.UUID())
		params    = orderV1.CancelOrderParams{
			OrderUUID: orderUUID,
		}
	)

	s.orderService.On("CancelOrder", s.ctx, orderUUID).Return(model.ErrCannotCancelPaidOrder)

	res, err := s.api.CancelOrder(s.ctx, params)
	s.Require().NoError(err)
	s.Require().NotNil(res)

	conflictErr, ok := res.(*orderV1.ConflictError)
	s.Require().True(ok)
	s.Require().Equal(http.StatusConflict, conflictErr.Code)
	s.Require().Contains(conflictErr.Message, "order cannot be cancelled")
}

func (s *APISuite) TestCancelOrderServiceError() {
	var (
		orderUUID  = uuid.MustParse(gofakeit.UUID())
		serviceErr = gofakeit.Error()
		params     = orderV1.CancelOrderParams{
			OrderUUID: orderUUID,
		}
	)

	s.orderService.On("CancelOrder", s.ctx, orderUUID).Return(serviceErr)

	res, err := s.api.CancelOrder(s.ctx, params)
	s.Require().NoError(err)
	s.Require().NotNil(res)

	internalErr, ok := res.(*orderV1.InternalServerError)
	s.Require().True(ok)
	s.Require().Equal(http.StatusInternalServerError, internalErr.Code)
	s.Require().Equal("internal server error", internalErr.Message)
}

func (s *APISuite) TestCancelOrderWithDifferentUUIDs() {
	orderUUIDs := []uuid.UUID{
		uuid.MustParse(gofakeit.UUID()),
		uuid.MustParse(gofakeit.UUID()),
		uuid.MustParse(gofakeit.UUID()),
	}

	for _, orderUUID := range orderUUIDs {
		params := orderV1.CancelOrderParams{
			OrderUUID: orderUUID,
		}

		s.orderService.On("CancelOrder", s.ctx, orderUUID).Return(nil)

		res, err := s.api.CancelOrder(s.ctx, params)
		s.Require().Error(err)
		s.Require().Nil(res)

		var statusCode *orderV1.GenericErrorStatusCode
		s.Require().ErrorAs(err, &statusCode)
		s.Require().Equal(http.StatusNoContent, statusCode.StatusCode)
	}
}

func (s *APISuite) TestCancelOrderWithSameUUID() {
	var (
		orderUUID = uuid.MustParse(gofakeit.UUID())
		params    = orderV1.CancelOrderParams{
			OrderUUID: orderUUID,
		}
	)

	s.orderService.On("CancelOrder", s.ctx, orderUUID).Return(nil)

	res, err := s.api.CancelOrder(s.ctx, params)
	s.Require().Error(err)
	s.Require().Nil(res)

	var statusCode *orderV1.GenericErrorStatusCode
	s.Require().ErrorAs(err, &statusCode)
	s.Require().Equal(http.StatusNoContent, statusCode.StatusCode)
}

func (s *APISuite) TestCancelOrderWithZeroUUID() {
	var (
		orderUUID = uuid.Nil // Zero UUID
		params    = orderV1.CancelOrderParams{
			OrderUUID: orderUUID,
		}
	)

	s.orderService.On("CancelOrder", s.ctx, orderUUID).Return(model.ErrOrderNotFound)

	res, err := s.api.CancelOrder(s.ctx, params)
	s.Require().NoError(err)
	s.Require().NotNil(res)

	notFoundErr, ok := res.(*orderV1.NotFoundError)
	s.Require().True(ok)
	s.Require().Equal(http.StatusNotFound, notFoundErr.Code)
	s.Require().Contains(notFoundErr.Message, "order not found")
}

func (s *APISuite) TestCancelOrderWithMaxUUID() {
	var (
		orderUUID = uuid.MustParse("ffffffff-ffff-ffff-ffff-ffffffffffff") // Max UUID
		params    = orderV1.CancelOrderParams{
			OrderUUID: orderUUID,
		}
	)

	s.orderService.On("CancelOrder", s.ctx, orderUUID).Return(nil)

	res, err := s.api.CancelOrder(s.ctx, params)
	s.Require().Error(err)
	s.Require().Nil(res)

	var statusCode *orderV1.GenericErrorStatusCode
	s.Require().ErrorAs(err, &statusCode)
	s.Require().Equal(http.StatusNoContent, statusCode.StatusCode)
}

func (s *APISuite) TestCancelOrderMultipleTimes() {
	var (
		orderUUID = uuid.MustParse(gofakeit.UUID())
		params    = orderV1.CancelOrderParams{
			OrderUUID: orderUUID,
		}
	)

	// First cancellation - success
	s.orderService.On("CancelOrder", s.ctx, orderUUID).Return(nil).Once()

	res, err := s.api.CancelOrder(s.ctx, params)
	s.Require().Error(err)
	s.Require().Nil(res)

	var statusCode *orderV1.GenericErrorStatusCode
	s.Require().ErrorAs(err, &statusCode)
	s.Require().Equal(http.StatusNoContent, statusCode.StatusCode)

	// Second cancellation - already cancelled (assuming this returns an error)
	s.orderService.On("CancelOrder", s.ctx, orderUUID).Return(model.ErrOrderNotFound).Once()

	res, err = s.api.CancelOrder(s.ctx, params)
	s.Require().NoError(err)
	s.Require().NotNil(res)

	notFoundErr, ok := res.(*orderV1.NotFoundError)
	s.Require().True(ok)
	s.Require().Equal(http.StatusNotFound, notFoundErr.Code)
}

func (s *APISuite) TestCancelOrderWithRandomUUIDs() {
	// Test with multiple random UUIDs
	for i := 0; i < 5; i++ {
		var (
			orderUUID = uuid.MustParse(gofakeit.UUID())
			params    = orderV1.CancelOrderParams{
				OrderUUID: orderUUID,
			}
		)

		s.orderService.On("CancelOrder", s.ctx, orderUUID).Return(nil)

		res, err := s.api.CancelOrder(s.ctx, params)
		s.Require().Error(err)
		s.Require().Nil(res)

		var statusCode *orderV1.GenericErrorStatusCode
		s.Require().ErrorAs(err, &statusCode)
		s.Require().Equal(http.StatusNoContent, statusCode.StatusCode)
	}
}

func (s *APISuite) TestCancelOrderWithSpecialUUID() {
	specialUUIDs := []string{
		"00000000-0000-0000-0000-000000000001",
		"12345678-1234-1234-1234-123456789012",
		"abcdefab-cdef-abcd-efab-cdefabcdefab",
	}

	for _, uuidStr := range specialUUIDs {
		var (
			orderUUID = uuid.MustParse(uuidStr)
			params    = orderV1.CancelOrderParams{
				OrderUUID: orderUUID,
			}
		)

		s.orderService.On("CancelOrder", s.ctx, orderUUID).Return(nil)

		res, err := s.api.CancelOrder(s.ctx, params)
		s.Require().Error(err)
		s.Require().Nil(res)

		var statusCode *orderV1.GenericErrorStatusCode
		s.Require().ErrorAs(err, &statusCode)
		s.Require().Equal(http.StatusNoContent, statusCode.StatusCode)
	}
}
