package v1

import (
	"github.com/brianvoe/gofakeit/v7"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/nkolesnikov999/micro2-OK/payment/internal/model"
	paymentV1 "github.com/nkolesnikov999/micro2-OK/shared/pkg/proto/payment/v1"
)

func (s *APISuite) TestPayOrderSuccess() {
	var (
		orderUUID = gofakeit.UUID()
		userUUID  = gofakeit.UUID()
		req       = &paymentV1.PayOrderRequest{
			OrderUuid:     orderUUID,
			UserUuid:      userUUID,
			PaymentMethod: paymentV1.PaymentMethod_PAYMENT_METHOD_CARD,
		}
		expectedTransactionUUID = gofakeit.UUID()
	)

	s.paymentService.On("PayOrder", s.ctx, "CARD").Return(expectedTransactionUUID, nil)

	res, err := s.api.PayOrder(s.ctx, req)
	s.Require().NoError(err)
	s.Require().NotNil(res)
	s.Require().Equal(expectedTransactionUUID, res.TransactionUuid)
}

func (s *APISuite) TestPayOrderWithSBP() {
	var (
		orderUUID = gofakeit.UUID()
		userUUID  = gofakeit.UUID()
		req       = &paymentV1.PayOrderRequest{
			OrderUuid:     orderUUID,
			UserUuid:      userUUID,
			PaymentMethod: paymentV1.PaymentMethod_PAYMENT_METHOD_SBP,
		}
		expectedTransactionUUID = gofakeit.UUID()
	)

	s.paymentService.On("PayOrder", s.ctx, "SBP").Return(expectedTransactionUUID, nil)

	res, err := s.api.PayOrder(s.ctx, req)
	s.Require().NoError(err)
	s.Require().NotNil(res)
	s.Require().Equal(expectedTransactionUUID, res.TransactionUuid)
}

func (s *APISuite) TestPayOrderWithCreditCard() {
	var (
		orderUUID = gofakeit.UUID()
		userUUID  = gofakeit.UUID()
		req       = &paymentV1.PayOrderRequest{
			OrderUuid:     orderUUID,
			UserUuid:      userUUID,
			PaymentMethod: paymentV1.PaymentMethod_PAYMENT_METHOD_CREDIT_CARD,
		}
		expectedTransactionUUID = gofakeit.UUID()
	)

	s.paymentService.On("PayOrder", s.ctx, "CREDIT_CARD").Return(expectedTransactionUUID, nil)

	res, err := s.api.PayOrder(s.ctx, req)
	s.Require().NoError(err)
	s.Require().NotNil(res)
	s.Require().Equal(expectedTransactionUUID, res.TransactionUuid)
}

func (s *APISuite) TestPayOrderWithInvestorMoney() {
	var (
		orderUUID = gofakeit.UUID()
		userUUID  = gofakeit.UUID()
		req       = &paymentV1.PayOrderRequest{
			OrderUuid:     orderUUID,
			UserUuid:      userUUID,
			PaymentMethod: paymentV1.PaymentMethod_PAYMENT_METHOD_INVESTOR_MONEY,
		}
		expectedTransactionUUID = gofakeit.UUID()
	)

	s.paymentService.On("PayOrder", s.ctx, "INVESTOR_MONEY").Return(expectedTransactionUUID, nil)

	res, err := s.api.PayOrder(s.ctx, req)
	s.Require().NoError(err)
	s.Require().NotNil(res)
	s.Require().Equal(expectedTransactionUUID, res.TransactionUuid)
}

func (s *APISuite) TestPayOrderNilRequest() {
	res, err := s.api.PayOrder(s.ctx, nil)
	s.Require().Error(err)
	s.Require().Nil(res)

	st, ok := status.FromError(err)
	s.Require().True(ok)
	s.Require().Equal(codes.Internal, st.Code())
	s.Require().Contains(st.Message(), "internal server error")
}

func (s *APISuite) TestPayOrderInvalidOrderUUID() {
	invalidUUIDs := []string{
		"invalid-uuid",
		"not-a-uuid",
		"123",
		"",
		"not-valid-format",
	}

	for _, invalidUUID := range invalidUUIDs {
		req := &paymentV1.PayOrderRequest{
			OrderUuid:     invalidUUID,
			UserUuid:      gofakeit.UUID(),
			PaymentMethod: paymentV1.PaymentMethod_PAYMENT_METHOD_CARD,
		}

		res, err := s.api.PayOrder(s.ctx, req)
		s.Require().Error(err)
		s.Require().Nil(res)

		st, ok := status.FromError(err)
		s.Require().True(ok)
		s.Require().Equal(codes.InvalidArgument, st.Code())
		s.Require().Contains(st.Message(), "invalid order_uuid format")
	}
}

func (s *APISuite) TestPayOrderUnspecifiedPaymentMethod() {
	var (
		orderUUID = gofakeit.UUID()
		userUUID  = gofakeit.UUID()
		req       = &paymentV1.PayOrderRequest{
			OrderUuid:     orderUUID,
			UserUuid:      userUUID,
			PaymentMethod: paymentV1.PaymentMethod_PAYMENT_METHOD_UNSPECIFIED,
		}
	)

	res, err := s.api.PayOrder(s.ctx, req)
	s.Require().Error(err)
	s.Require().Nil(res)

	st, ok := status.FromError(err)
	s.Require().True(ok)
	s.Require().Equal(codes.InvalidArgument, st.Code())
	s.Require().Contains(st.Message(), "payment_method must be specified")
}

func (s *APISuite) TestPayOrderInvalidPaymentMethod() {
	var (
		orderUUID = gofakeit.UUID()
		userUUID  = gofakeit.UUID()
		req       = &paymentV1.PayOrderRequest{
			OrderUuid:     orderUUID,
			UserUuid:      userUUID,
			PaymentMethod: paymentV1.PaymentMethod(999), // invalid enum value
		}
	)

	res, err := s.api.PayOrder(s.ctx, req)
	s.Require().Error(err)
	s.Require().Nil(res)

	st, ok := status.FromError(err)
	s.Require().True(ok)
	s.Require().Equal(codes.InvalidArgument, st.Code())
	s.Require().Contains(st.Message(), "invalid payment_method")
}

func (s *APISuite) TestPayOrderServiceInvalidPaymentMethod() {
	var (
		orderUUID = gofakeit.UUID()
		userUUID  = gofakeit.UUID()
		req       = &paymentV1.PayOrderRequest{
			OrderUuid:     orderUUID,
			UserUuid:      userUUID,
			PaymentMethod: paymentV1.PaymentMethod_PAYMENT_METHOD_CARD,
		}
	)

	s.paymentService.On("PayOrder", s.ctx, "CARD").Return("", model.ErrInvalidPaymentMethod)

	res, err := s.api.PayOrder(s.ctx, req)
	s.Require().Error(err)
	s.Require().Nil(res)

	st, ok := status.FromError(err)
	s.Require().True(ok)
	s.Require().Equal(codes.InvalidArgument, st.Code())
	s.Require().Contains(st.Message(), "invalid payment method")
}

func (s *APISuite) TestPayOrderServiceError() {
	var (
		orderUUID = gofakeit.UUID()
		userUUID  = gofakeit.UUID()
		req       = &paymentV1.PayOrderRequest{
			OrderUuid:     orderUUID,
			UserUuid:      userUUID,
			PaymentMethod: paymentV1.PaymentMethod_PAYMENT_METHOD_CARD,
		}
		serviceErr = gofakeit.Error()
	)

	s.paymentService.On("PayOrder", s.ctx, "CARD").Return("", serviceErr)

	res, err := s.api.PayOrder(s.ctx, req)
	s.Require().Error(err)
	s.Require().Nil(res)

	st, ok := status.FromError(err)
	s.Require().True(ok)
	s.Require().Equal(codes.Internal, st.Code())
	s.Require().Contains(st.Message(), "internal server error")
}

func (s *APISuite) TestPayOrderEmptyUserUUID() {
	var (
		orderUUID = gofakeit.UUID()
		req       = &paymentV1.PayOrderRequest{
			OrderUuid:     orderUUID,
			UserUuid:      "", // empty user UUID
			PaymentMethod: paymentV1.PaymentMethod_PAYMENT_METHOD_CARD,
		}
		expectedTransactionUUID = gofakeit.UUID()
	)

	// Empty user UUID should not cause validation error in API layer
	s.paymentService.On("PayOrder", s.ctx, "CARD").Return(expectedTransactionUUID, nil)

	res, err := s.api.PayOrder(s.ctx, req)
	s.Require().NoError(err)
	s.Require().NotNil(res)
	s.Require().Equal(expectedTransactionUUID, res.TransactionUuid)
}

func (s *APISuite) TestPayOrderWithValidUUIDs() {
	var (
		orderUUID = uuid.New().String()
		userUUID  = uuid.New().String()
		req       = &paymentV1.PayOrderRequest{
			OrderUuid:     orderUUID,
			UserUuid:      userUUID,
			PaymentMethod: paymentV1.PaymentMethod_PAYMENT_METHOD_SBP,
		}
		expectedTransactionUUID = gofakeit.UUID()
	)

	s.paymentService.On("PayOrder", s.ctx, "SBP").Return(expectedTransactionUUID, nil)

	res, err := s.api.PayOrder(s.ctx, req)
	s.Require().NoError(err)
	s.Require().NotNil(res)
	s.Require().Equal(expectedTransactionUUID, res.TransactionUuid)
}

func (s *APISuite) TestPayOrderWithLongTransactionUUID() {
	var (
		orderUUID = gofakeit.UUID()
		userUUID  = gofakeit.UUID()
		req       = &paymentV1.PayOrderRequest{
			OrderUuid:     orderUUID,
			UserUuid:      userUUID,
			PaymentMethod: paymentV1.PaymentMethod_PAYMENT_METHOD_CREDIT_CARD,
		}
		expectedTransactionUUID = gofakeit.UUID() + gofakeit.UUID() // long UUID
	)

	s.paymentService.On("PayOrder", s.ctx, "CREDIT_CARD").Return(expectedTransactionUUID, nil)

	res, err := s.api.PayOrder(s.ctx, req)
	s.Require().NoError(err)
	s.Require().NotNil(res)
	s.Require().Equal(expectedTransactionUUID, res.TransactionUuid)
}

func (s *APISuite) TestPayOrderWithSameOrderAndUserUUID() {
	var (
		sharedUUID = gofakeit.UUID()
		req        = &paymentV1.PayOrderRequest{
			OrderUuid:     sharedUUID,
			UserUuid:      sharedUUID, // same UUID for order and user
			PaymentMethod: paymentV1.PaymentMethod_PAYMENT_METHOD_INVESTOR_MONEY,
		}
		expectedTransactionUUID = gofakeit.UUID()
	)

	s.paymentService.On("PayOrder", s.ctx, "INVESTOR_MONEY").Return(expectedTransactionUUID, nil)

	res, err := s.api.PayOrder(s.ctx, req)
	s.Require().NoError(err)
	s.Require().NotNil(res)
	s.Require().Equal(expectedTransactionUUID, res.TransactionUuid)
}
