package v1

import (
	"net/http"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/google/uuid"

	"github.com/nkolesnikov999/micro2-OK/order/internal/model"
	orderV1 "github.com/nkolesnikov999/micro2-OK/shared/pkg/openapi/order/v1"
)

func (s *APISuite) TestPayOrderSuccess() {
	var (
		orderUUID = uuid.MustParse(gofakeit.UUID())
		req       = &orderV1.PayOrderRequest{
			PaymentMethod: orderV1.PaymentMethodPAYMENTMETHODCARD,
		}
		params = orderV1.PayOrderParams{
			OrderUUID: orderUUID,
		}
		expectedTransactionUUID = gofakeit.UUID()
	)

	s.orderService.On("PayOrder", s.ctx, orderUUID, "CARD").Return(expectedTransactionUUID, nil)

	res, err := s.api.PayOrder(s.ctx, req, params)
	s.Require().NoError(err)
	s.Require().NotNil(res)

	payOrderResp, ok := res.(*orderV1.PayOrderResponse)
	s.Require().True(ok)
	s.Require().Equal(expectedTransactionUUID, payOrderResp.TransactionUUID)
}

func (s *APISuite) TestPayOrderWithSBP() {
	var (
		orderUUID = uuid.MustParse(gofakeit.UUID())
		req       = &orderV1.PayOrderRequest{
			PaymentMethod: orderV1.PaymentMethodPAYMENTMETHODSBP,
		}
		params = orderV1.PayOrderParams{
			OrderUUID: orderUUID,
		}
		expectedTransactionUUID = gofakeit.UUID()
	)

	s.orderService.On("PayOrder", s.ctx, orderUUID, "SBP").Return(expectedTransactionUUID, nil)

	res, err := s.api.PayOrder(s.ctx, req, params)
	s.Require().NoError(err)
	s.Require().NotNil(res)

	payOrderResp, ok := res.(*orderV1.PayOrderResponse)
	s.Require().True(ok)
	s.Require().Equal(expectedTransactionUUID, payOrderResp.TransactionUUID)
}

func (s *APISuite) TestPayOrderWithCreditCard() {
	var (
		orderUUID = uuid.MustParse(gofakeit.UUID())
		req       = &orderV1.PayOrderRequest{
			PaymentMethod: orderV1.PaymentMethodPAYMENTMETHODCREDITCARD,
		}
		params = orderV1.PayOrderParams{
			OrderUUID: orderUUID,
		}
		expectedTransactionUUID = gofakeit.UUID()
	)

	s.orderService.On("PayOrder", s.ctx, orderUUID, "CREDIT_CARD").Return(expectedTransactionUUID, nil)

	res, err := s.api.PayOrder(s.ctx, req, params)
	s.Require().NoError(err)
	s.Require().NotNil(res)

	payOrderResp, ok := res.(*orderV1.PayOrderResponse)
	s.Require().True(ok)
	s.Require().Equal(expectedTransactionUUID, payOrderResp.TransactionUUID)
}

func (s *APISuite) TestPayOrderWithInvestorMoney() {
	var (
		orderUUID = uuid.MustParse(gofakeit.UUID())
		req       = &orderV1.PayOrderRequest{
			PaymentMethod: orderV1.PaymentMethodPAYMENTMETHODINVESTORMONEY,
		}
		params = orderV1.PayOrderParams{
			OrderUUID: orderUUID,
		}
		expectedTransactionUUID = gofakeit.UUID()
	)

	s.orderService.On("PayOrder", s.ctx, orderUUID, "INVESTOR_MONEY").Return(expectedTransactionUUID, nil)

	res, err := s.api.PayOrder(s.ctx, req, params)
	s.Require().NoError(err)
	s.Require().NotNil(res)

	payOrderResp, ok := res.(*orderV1.PayOrderResponse)
	s.Require().True(ok)
	s.Require().Equal(expectedTransactionUUID, payOrderResp.TransactionUUID)
}

func (s *APISuite) TestPayOrderNilRequest() {
	params := orderV1.PayOrderParams{
		OrderUUID: uuid.MustParse(gofakeit.UUID()),
	}

	res, err := s.api.PayOrder(s.ctx, nil, params)
	s.Require().NoError(err)
	s.Require().NotNil(res)

	internalErr, ok := res.(*orderV1.InternalServerError)
	s.Require().True(ok)
	s.Require().Equal(http.StatusInternalServerError, internalErr.Code)
	s.Require().Contains(internalErr.Message, "internal server error")
}

func (s *APISuite) TestPayOrderNotFound() {
	var (
		orderUUID = uuid.MustParse(gofakeit.UUID())
		req       = &orderV1.PayOrderRequest{
			PaymentMethod: orderV1.PaymentMethodPAYMENTMETHODCARD,
		}
		params = orderV1.PayOrderParams{
			OrderUUID: orderUUID,
		}
	)

	s.orderService.On("PayOrder", s.ctx, orderUUID, "CARD").Return("", model.ErrOrderNotFound)

	res, err := s.api.PayOrder(s.ctx, req, params)
	s.Require().NoError(err)
	s.Require().NotNil(res)

	notFoundErr, ok := res.(*orderV1.NotFoundError)
	s.Require().True(ok)
	s.Require().Equal(http.StatusNotFound, notFoundErr.Code)
	s.Require().Contains(notFoundErr.Message, "order not found")
}

func (s *APISuite) TestPayOrderAlreadyPaid() {
	var (
		orderUUID = uuid.MustParse(gofakeit.UUID())
		req       = &orderV1.PayOrderRequest{
			PaymentMethod: orderV1.PaymentMethodPAYMENTMETHODCARD,
		}
		params = orderV1.PayOrderParams{
			OrderUUID: orderUUID,
		}
	)

	s.orderService.On("PayOrder", s.ctx, orderUUID, "CARD").Return("", model.ErrOrderAlreadyPaid)

	res, err := s.api.PayOrder(s.ctx, req, params)
	s.Require().NoError(err)
	s.Require().NotNil(res)

	conflictErr, ok := res.(*orderV1.ConflictError)
	s.Require().True(ok)
	s.Require().Equal(http.StatusConflict, conflictErr.Code)
	s.Require().Contains(conflictErr.Message, "order already paid")
}

func (s *APISuite) TestPayOrderCannotPayCancelledOrder() {
	var (
		orderUUID = uuid.MustParse(gofakeit.UUID())
		req       = &orderV1.PayOrderRequest{
			PaymentMethod: orderV1.PaymentMethodPAYMENTMETHODCARD,
		}
		params = orderV1.PayOrderParams{
			OrderUUID: orderUUID,
		}
	)

	s.orderService.On("PayOrder", s.ctx, orderUUID, "CARD").Return("", model.ErrCannotPayCancelledOrder)

	res, err := s.api.PayOrder(s.ctx, req, params)
	s.Require().NoError(err)
	s.Require().NotNil(res)

	conflictErr, ok := res.(*orderV1.ConflictError)
	s.Require().True(ok)
	s.Require().Equal(http.StatusConflict, conflictErr.Code)
	s.Require().Contains(conflictErr.Message, "cannot pay cancelled order")
}

func (s *APISuite) TestPayOrderPaymentFailed() {
	var (
		orderUUID = uuid.MustParse(gofakeit.UUID())
		req       = &orderV1.PayOrderRequest{
			PaymentMethod: orderV1.PaymentMethodPAYMENTMETHODCARD,
		}
		params = orderV1.PayOrderParams{
			OrderUUID: orderUUID,
		}
	)

	s.orderService.On("PayOrder", s.ctx, orderUUID, "CARD").Return("", model.ErrPaymentFailed)

	res, err := s.api.PayOrder(s.ctx, req, params)
	s.Require().NoError(err)
	s.Require().NotNil(res)

	internalErr, ok := res.(*orderV1.InternalServerError)
	s.Require().True(ok)
	s.Require().Equal(http.StatusInternalServerError, internalErr.Code)
	s.Require().Contains(internalErr.Message, "payment failed")
}

func (s *APISuite) TestPayOrderServiceError() {
	var (
		orderUUID  = uuid.MustParse(gofakeit.UUID())
		serviceErr = gofakeit.Error()
		req        = &orderV1.PayOrderRequest{
			PaymentMethod: orderV1.PaymentMethodPAYMENTMETHODCARD,
		}
		params = orderV1.PayOrderParams{
			OrderUUID: orderUUID,
		}
	)

	s.orderService.On("PayOrder", s.ctx, orderUUID, "CARD").Return("", serviceErr)

	res, err := s.api.PayOrder(s.ctx, req, params)
	s.Require().NoError(err)
	s.Require().NotNil(res)

	internalErr, ok := res.(*orderV1.InternalServerError)
	s.Require().True(ok)
	s.Require().Equal(http.StatusInternalServerError, internalErr.Code)
	s.Require().Equal(serviceErr.Error(), internalErr.Message)
}

func (s *APISuite) TestPayOrderWithLongTransactionUUID() {
	var (
		orderUUID = uuid.MustParse(gofakeit.UUID())
		req       = &orderV1.PayOrderRequest{
			PaymentMethod: orderV1.PaymentMethodPAYMENTMETHODSBP,
		}
		params = orderV1.PayOrderParams{
			OrderUUID: orderUUID,
		}
		expectedTransactionUUID = gofakeit.UUID() + gofakeit.UUID() // long UUID
	)

	s.orderService.On("PayOrder", s.ctx, orderUUID, "SBP").Return(expectedTransactionUUID, nil)

	res, err := s.api.PayOrder(s.ctx, req, params)
	s.Require().NoError(err)
	s.Require().NotNil(res)

	payOrderResp, ok := res.(*orderV1.PayOrderResponse)
	s.Require().True(ok)
	s.Require().Equal(expectedTransactionUUID, payOrderResp.TransactionUUID)
}

func (s *APISuite) TestPayOrderWithEmptyTransactionUUID() {
	var (
		orderUUID = uuid.MustParse(gofakeit.UUID())
		req       = &orderV1.PayOrderRequest{
			PaymentMethod: orderV1.PaymentMethodPAYMENTMETHODCARD,
		}
		params = orderV1.PayOrderParams{
			OrderUUID: orderUUID,
		}
		expectedTransactionUUID = "" // empty transaction UUID
	)

	s.orderService.On("PayOrder", s.ctx, orderUUID, "CARD").Return(expectedTransactionUUID, nil)

	res, err := s.api.PayOrder(s.ctx, req, params)
	s.Require().NoError(err)
	s.Require().NotNil(res)

	payOrderResp, ok := res.(*orderV1.PayOrderResponse)
	s.Require().True(ok)
	s.Require().Equal(expectedTransactionUUID, payOrderResp.TransactionUUID)
}

func (s *APISuite) TestPayOrderWithDifferentPaymentMethods() {
	paymentMethods := []struct {
		apiMethod     orderV1.PaymentMethod
		serviceMethod string
	}{
		{orderV1.PaymentMethodPAYMENTMETHODCARD, "CARD"},
		{orderV1.PaymentMethodPAYMENTMETHODSBP, "SBP"},
		{orderV1.PaymentMethodPAYMENTMETHODCREDITCARD, "CREDIT_CARD"},
		{orderV1.PaymentMethodPAYMENTMETHODINVESTORMONEY, "INVESTOR_MONEY"},
	}

	for _, pm := range paymentMethods {
		var (
			orderUUID = uuid.MustParse(gofakeit.UUID())
			req       = &orderV1.PayOrderRequest{
				PaymentMethod: pm.apiMethod,
			}
			params = orderV1.PayOrderParams{
				OrderUUID: orderUUID,
			}
			expectedTransactionUUID = gofakeit.UUID()
		)

		s.orderService.On("PayOrder", s.ctx, orderUUID, pm.serviceMethod).Return(expectedTransactionUUID, nil)

		res, err := s.api.PayOrder(s.ctx, req, params)
		s.Require().NoError(err)
		s.Require().NotNil(res)

		payOrderResp, ok := res.(*orderV1.PayOrderResponse)
		s.Require().True(ok)
		s.Require().Equal(expectedTransactionUUID, payOrderResp.TransactionUUID)
	}
}

func (s *APISuite) TestPayOrderWithSameOrderUUID() {
	var (
		orderUUID = uuid.MustParse(gofakeit.UUID())
		req       = &orderV1.PayOrderRequest{
			PaymentMethod: orderV1.PaymentMethodPAYMENTMETHODCARD,
		}
		params = orderV1.PayOrderParams{
			OrderUUID: orderUUID,
		}
		expectedTransactionUUID = gofakeit.UUID()
	)

	s.orderService.On("PayOrder", s.ctx, orderUUID, "CARD").Return(expectedTransactionUUID, nil)

	res, err := s.api.PayOrder(s.ctx, req, params)
	s.Require().NoError(err)
	s.Require().NotNil(res)

	payOrderResp, ok := res.(*orderV1.PayOrderResponse)
	s.Require().True(ok)
	s.Require().Equal(expectedTransactionUUID, payOrderResp.TransactionUUID)
}
