package payment

import (
	"strings"

	"github.com/brianvoe/gofakeit/v7"

	"github.com/nkolesnikov999/micro2-OK/payment/internal/model"
)

func (s *ServiceSuite) TestPayOrderSuccess() {
	paymentMethod := "CARD"

	transactionUUID, err := s.service.PayOrder(s.ctx, paymentMethod)
	s.NoError(err)
	s.NotEmpty(transactionUUID)
	s.Len(transactionUUID, 36) // UUID length
}

func (s *ServiceSuite) TestPayOrderWithCard() {
	paymentMethod := "CARD"

	transactionUUID, err := s.service.PayOrder(s.ctx, paymentMethod)
	s.NoError(err)
	s.NotEmpty(transactionUUID)
	s.Len(transactionUUID, 36)
}

func (s *ServiceSuite) TestPayOrderWithSBP() {
	paymentMethod := "SBP"

	transactionUUID, err := s.service.PayOrder(s.ctx, paymentMethod)
	s.NoError(err)
	s.NotEmpty(transactionUUID)
	s.Len(transactionUUID, 36)
}

func (s *ServiceSuite) TestPayOrderWithCreditCard() {
	paymentMethod := "CREDIT_CARD"

	transactionUUID, err := s.service.PayOrder(s.ctx, paymentMethod)
	s.NoError(err)
	s.NotEmpty(transactionUUID)
	s.Len(transactionUUID, 36)
}

func (s *ServiceSuite) TestPayOrderWithInvestorMoney() {
	paymentMethod := "INVESTOR_MONEY"

	transactionUUID, err := s.service.PayOrder(s.ctx, paymentMethod)
	s.NoError(err)
	s.NotEmpty(transactionUUID)
	s.Len(transactionUUID, 36)
}

func (s *ServiceSuite) TestPayOrderWithLowercaseMethod() {
	paymentMethod := "card" // lowercase

	transactionUUID, err := s.service.PayOrder(s.ctx, paymentMethod)
	s.NoError(err)
	s.NotEmpty(transactionUUID)
	s.Len(transactionUUID, 36)
}

func (s *ServiceSuite) TestPayOrderWithMixedCaseMethod() {
	paymentMethod := "Card" // mixed case

	transactionUUID, err := s.service.PayOrder(s.ctx, paymentMethod)
	s.NoError(err)
	s.NotEmpty(transactionUUID)
	s.Len(transactionUUID, 36)
}

func (s *ServiceSuite) TestPayOrderWithWhitespace() {
	paymentMethod := "  CARD  " // with whitespace

	transactionUUID, err := s.service.PayOrder(s.ctx, paymentMethod)
	s.NoError(err)
	s.NotEmpty(transactionUUID)
	s.Len(transactionUUID, 36)
}

func (s *ServiceSuite) TestPayOrderWithTabWhitespace() {
	paymentMethod := "\tCARD\t" // with tab whitespace

	transactionUUID, err := s.service.PayOrder(s.ctx, paymentMethod)
	s.NoError(err)
	s.NotEmpty(transactionUUID)
	s.Len(transactionUUID, 36)
}

func (s *ServiceSuite) TestPayOrderWithNewlineWhitespace() {
	paymentMethod := "\nCARD\n" // with newline whitespace

	transactionUUID, err := s.service.PayOrder(s.ctx, paymentMethod)
	s.NoError(err)
	s.NotEmpty(transactionUUID)
	s.Len(transactionUUID, 36)
}

func (s *ServiceSuite) TestPayOrderWithInvalidMethod() {
	paymentMethod := "INVALID_METHOD"

	transactionUUID, err := s.service.PayOrder(s.ctx, paymentMethod)
	s.Error(err)
	s.ErrorIs(err, model.ErrInvalidPaymentMethod)
	s.Empty(transactionUUID)
}

func (s *ServiceSuite) TestPayOrderWithEmptyMethod() {
	paymentMethod := ""

	transactionUUID, err := s.service.PayOrder(s.ctx, paymentMethod)
	s.Error(err)
	s.ErrorIs(err, model.ErrInvalidPaymentMethod)
	s.Empty(transactionUUID)
}

func (s *ServiceSuite) TestPayOrderWithWhitespaceOnly() {
	paymentMethod := "   " // only whitespace

	transactionUUID, err := s.service.PayOrder(s.ctx, paymentMethod)
	s.Error(err)
	s.ErrorIs(err, model.ErrInvalidPaymentMethod)
	s.Empty(transactionUUID)
}

func (s *ServiceSuite) TestPayOrderWithRandomString() {
	paymentMethod := gofakeit.Word()

	transactionUUID, err := s.service.PayOrder(s.ctx, paymentMethod)
	s.Error(err)
	s.ErrorIs(err, model.ErrInvalidPaymentMethod)
	s.Empty(transactionUUID)
}

func (s *ServiceSuite) TestPayOrderWithNumberString() {
	paymentMethod := "12345"

	transactionUUID, err := s.service.PayOrder(s.ctx, paymentMethod)
	s.Error(err)
	s.ErrorIs(err, model.ErrInvalidPaymentMethod)
	s.Empty(transactionUUID)
}

func (s *ServiceSuite) TestPayOrderWithSpecialCharacters() {
	paymentMethod := "CARD@#$%"

	transactionUUID, err := s.service.PayOrder(s.ctx, paymentMethod)
	s.Error(err)
	s.ErrorIs(err, model.ErrInvalidPaymentMethod)
	s.Empty(transactionUUID)
}

func (s *ServiceSuite) TestPayOrderWithPartialMatch() {
	paymentMethod := "CAR" // partial match of CARD

	transactionUUID, err := s.service.PayOrder(s.ctx, paymentMethod)
	s.Error(err)
	s.ErrorIs(err, model.ErrInvalidPaymentMethod)
	s.Empty(transactionUUID)
}

func (s *ServiceSuite) TestPayOrderWithMultipleWords() {
	paymentMethod := "CREDIT CARD" // space in method name

	transactionUUID, err := s.service.PayOrder(s.ctx, paymentMethod)
	s.Error(err)
	s.ErrorIs(err, model.ErrInvalidPaymentMethod)
	s.Empty(transactionUUID)
}

func (s *ServiceSuite) TestPayOrderWithVeryLongString() {
	paymentMethod := strings.Repeat("CARD", 100) // very long string

	transactionUUID, err := s.service.PayOrder(s.ctx, paymentMethod)
	s.Error(err)
	s.ErrorIs(err, model.ErrInvalidPaymentMethod)
	s.Empty(transactionUUID)
}

func (s *ServiceSuite) TestPayOrderWithUnicodeCharacters() {
	paymentMethod := "КАРТА" // unicode characters

	transactionUUID, err := s.service.PayOrder(s.ctx, paymentMethod)
	s.Error(err)
	s.ErrorIs(err, model.ErrInvalidPaymentMethod)
	s.Empty(transactionUUID)
}

func (s *ServiceSuite) TestPayOrderWithNullCharacter() {
	paymentMethod := "CARD\x00" // null character

	transactionUUID, err := s.service.PayOrder(s.ctx, paymentMethod)
	s.Error(err)
	s.ErrorIs(err, model.ErrInvalidPaymentMethod)
	s.Empty(transactionUUID)
}

func (s *ServiceSuite) TestPayOrderGeneratesUniqueUUIDs() {
	paymentMethod := "CARD"

	// Generate multiple UUIDs and ensure they're unique
	uuids := make(map[string]bool)
	for i := 0; i < 100; i++ {
		transactionUUID, err := s.service.PayOrder(s.ctx, paymentMethod)
		s.NoError(err)
		s.NotEmpty(transactionUUID)

		// Check for uniqueness
		s.False(uuids[transactionUUID], "UUID should be unique: %s", transactionUUID)
		uuids[transactionUUID] = true
	}

	s.Len(uuids, 100, "Should generate 100 unique UUIDs")
}

func (s *ServiceSuite) TestPayOrderWithAllValidMethods() {
	validMethods := []string{"CARD", "SBP", "CREDIT_CARD", "INVESTOR_MONEY"}

	for _, method := range validMethods {
		transactionUUID, err := s.service.PayOrder(s.ctx, method)
		s.NoError(err, "Method %s should be valid", method)
		s.NotEmpty(transactionUUID)
		s.Len(transactionUUID, 36)
	}
}

func (s *ServiceSuite) TestPayOrderWithAllValidMethodsLowercase() {
	validMethods := []string{"card", "sbp", "credit_card", "investor_money"}

	for _, method := range validMethods {
		transactionUUID, err := s.service.PayOrder(s.ctx, method)
		s.NoError(err, "Method %s should be valid", method)
		s.NotEmpty(transactionUUID)
		s.Len(transactionUUID, 36)
	}
}
