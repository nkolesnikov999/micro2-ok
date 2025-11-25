package telegram

import (
	"bytes"
	"context"
	"embed"
	"text/template"

	"go.uber.org/zap"

	"github.com/nkolesnikov999/micro2-OK/notification/internal/client/http"
	"github.com/nkolesnikov999/micro2-OK/notification/internal/model"
	"github.com/nkolesnikov999/micro2-OK/platform/pkg/logger"
)

//go:embed templates/order_paid.tmpl
//go:embed templates/ship_assembled.tmpl
var templateFS embed.FS

type orderPaidTemplateData struct {
	EventUUID       string
	OrderUUID       string
	UserUUID        string
	PaymentMethod   string
	TransactionUUID string
}

type shipAssembledTemplateData struct {
	EventUUID    string
	OrderUUID    string
	UserUUID     string
	BuildTimeSec int64
}

var (
	orderPaidTemplate     = template.Must(template.ParseFS(templateFS, "templates/order_paid.tmpl"))
	shipAssembledTemplate = template.Must(template.ParseFS(templateFS, "templates/ship_assembled.tmpl"))
)

type service struct {
	telegramClient http.TelegramClient
	chatID         int64
}

func NewService(telegramClient http.TelegramClient, chatID int64) *service {
	return &service{
		telegramClient: telegramClient,
		chatID:         chatID,
	}
}

func (s *service) SendOrderPaidNotification(ctx context.Context, orderPaidEvent model.OrderPaidEvent) error {
	message, err := s.buildOrderPaidMessage(orderPaidEvent)
	if err != nil {
		return err
	}

	err = s.telegramClient.SendMessage(ctx, s.chatID, message)
	if err != nil {
		return err
	}

	logger.Info(ctx, "Telegram message sent to chat", zap.Int64("chat_id", s.chatID), zap.String("message", message))
	return nil
}

func (s *service) buildOrderPaidMessage(orderPaidEvent model.OrderPaidEvent) (string, error) {
	data := orderPaidTemplateData{
		EventUUID:       orderPaidEvent.EventUUID,
		OrderUUID:       orderPaidEvent.OrderUUID,
		UserUUID:        orderPaidEvent.UserUUID,
		PaymentMethod:   orderPaidEvent.PaymentMethod,
		TransactionUUID: orderPaidEvent.TransactionUUID,
	}

	var buf bytes.Buffer
	err := orderPaidTemplate.Execute(&buf, data)
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}

func (s *service) SendOrderAssembledNotification(ctx context.Context, shipAssembledEvent model.ShipAssembledEvent) error {
	message, err := s.buildShipAssembledMessage(shipAssembledEvent)
	if err != nil {
		return err
	}

	err = s.telegramClient.SendMessage(ctx, s.chatID, message)
	if err != nil {
		return err
	}

	logger.Info(ctx, "Telegram message sent to chat", zap.Int64("chat_id", s.chatID), zap.String("message", message))
	return nil
}

func (s *service) buildShipAssembledMessage(shipAssembledEvent model.ShipAssembledEvent) (string, error) {
	data := shipAssembledTemplateData{
		EventUUID:    shipAssembledEvent.EventUUID,
		OrderUUID:    shipAssembledEvent.OrderUUID,
		UserUUID:     shipAssembledEvent.UserUUID,
		BuildTimeSec: shipAssembledEvent.BuildTimeSec,
	}

	var buf bytes.Buffer
	err := shipAssembledTemplate.Execute(&buf, data)
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}
