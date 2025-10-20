package model

// PaymentMethod represents supported payment methods in a type-safe way.
type PaymentMethod string

const (
	PaymentMethodUnspecified   PaymentMethod = ""
	PaymentMethodCard          PaymentMethod = "CARD"
	PaymentMethodSBP           PaymentMethod = "SBP"
	PaymentMethodCreditCard    PaymentMethod = "CREDIT_CARD"
	PaymentMethodInvestorMoney PaymentMethod = "INVESTOR_MONEY"
)
