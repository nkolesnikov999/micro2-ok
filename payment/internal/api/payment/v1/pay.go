package v1

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/nkolesnikov999/micro2-OK/payment/internal/model"
	"github.com/nkolesnikov999/micro2-OK/platform/pkg/logger"
	paymentV1 "github.com/nkolesnikov999/micro2-OK/shared/pkg/proto/payment/v1"
)

func (a *api) PayOrder(ctx context.Context, req *paymentV1.PayOrderRequest) (*paymentV1.PayOrderResponse, error) {
	if req == nil {
		return nil, status.Error(codes.Internal, "internal server error")
	}

	if _, err := uuid.Parse(req.GetOrderUuid()); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid order_uuid format: %v", err)
	}

	pm := req.GetPaymentMethod()
	if pm == paymentV1.PaymentMethod_PAYMENT_METHOD_UNSPECIFIED {
		return nil, status.Error(codes.InvalidArgument, "payment_method must be specified")
	}

	var method model.PaymentMethod
	switch pm {
	case paymentV1.PaymentMethod_PAYMENT_METHOD_CARD:
		method = model.PaymentMethodCard
	case paymentV1.PaymentMethod_PAYMENT_METHOD_SBP:
		method = model.PaymentMethodSBP
	case paymentV1.PaymentMethod_PAYMENT_METHOD_CREDIT_CARD:
		method = model.PaymentMethodCreditCard
	case paymentV1.PaymentMethod_PAYMENT_METHOD_INVESTOR_MONEY:
		method = model.PaymentMethodInvestorMoney
	default:
		return nil, status.Errorf(codes.InvalidArgument, "invalid payment_method: %v", pm)
	}

	txn, err := a.paymentService.PayOrder(ctx, string(method))
	if err != nil {
		if errors.Is(err, model.ErrInvalidPaymentMethod) {
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}
		return nil, status.Error(codes.Internal, "internal server error")
	}

	logger.Info(ctx, "payment successful", zap.String("transaction_uuid", txn),
		zap.String("order_uuid", req.GetOrderUuid()),
		zap.String("user_uuid", req.GetUserUuid()),
		zap.String("payment_method", pm.String()),
		zap.Error(err),
	)

	return &paymentV1.PayOrderResponse{TransactionUuid: txn}, nil
}
