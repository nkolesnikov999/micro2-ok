package order

import (
	"github.com/brianvoe/gofakeit/v7"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"

	"github.com/nkolesnikov999/micro2-OK/order/internal/model"
)

func (s *ServiceSuite) TestCreateOrderSuccess() {
	userUUID := uuid.New()
	partUUIDs := []uuid.UUID{uuid.New(), uuid.New()}
	parts := []model.Part{
		{
			Uuid:  partUUIDs[0],
			Price: 100.0,
		},
		{
			Uuid:  partUUIDs[1],
			Price: 200.0,
		},
	}

	s.inventoryClient.On("ListParts", s.ctx, model.PartsFilter{Uuids: partUUIDs}).Return(parts, nil)
	s.orderRepository.On("CreateOrder", s.ctx, mock.MatchedBy(func(order model.Order) bool {
		return order.UserUUID == userUUID &&
			len(order.PartUuids) == len(partUUIDs) &&
			order.TotalPrice == 300.0 &&
			order.Status == "PENDING_PAYMENT" &&
			order.OrderUUID != uuid.Nil
	})).Return(nil)

	order, err := s.service.CreateOrder(s.ctx, userUUID, partUUIDs)
	s.NoError(err)
	s.Equal(userUUID, order.UserUUID)
	s.Equal(partUUIDs, order.PartUuids)
	s.Equal(300.0, order.TotalPrice)
	s.Equal("PENDING_PAYMENT", order.Status)
	s.NotEmpty(order.OrderUUID)
}

func (s *ServiceSuite) TestCreateOrderEmptyPartUUIDs() {
	userUUID := uuid.New()
	partUUIDs := []uuid.UUID{}

	order, err := s.service.CreateOrder(s.ctx, userUUID, partUUIDs)
	s.Error(err)
	s.ErrorIs(err, model.ErrEmptyPartUUIDs)
	s.Empty(order)
}

func (s *ServiceSuite) TestCreateOrderNilPartUUIDs() {
	userUUID := uuid.New()
	partUUIDs := []uuid.UUID(nil)

	order, err := s.service.CreateOrder(s.ctx, userUUID, partUUIDs)
	s.Error(err)
	s.ErrorIs(err, model.ErrEmptyPartUUIDs)
	s.Empty(order)
}

func (s *ServiceSuite) TestCreateOrderInventoryUnavailable() {
	userUUID := uuid.New()
	partUUIDs := []uuid.UUID{uuid.New()}
	inventoryErr := gofakeit.Error()

	s.inventoryClient.On("ListParts", s.ctx, model.PartsFilter{Uuids: partUUIDs}).Return([]model.Part{}, inventoryErr)

	order, err := s.service.CreateOrder(s.ctx, userUUID, partUUIDs)
	s.Error(err)
	s.ErrorIs(err, model.ErrInventoryUnavailable)
	s.Empty(order)
}

func (s *ServiceSuite) TestCreateOrderPartsNotFound() {
	userUUID := uuid.New()
	partUUIDs := []uuid.UUID{uuid.New(), uuid.New()}
	parts := []model.Part{
		{
			Uuid:  partUUIDs[0],
			Price: 100.0,
		},
		// partUUIDs[1] is missing
	}

	s.inventoryClient.On("ListParts", s.ctx, model.PartsFilter{Uuids: partUUIDs}).Return(parts, nil)

	order, err := s.service.CreateOrder(s.ctx, userUUID, partUUIDs)
	s.Error(err)

	var partsNotFoundErr *model.PartsNotFoundError
	s.ErrorAs(err, &partsNotFoundErr)
	s.Len(partsNotFoundErr.MissingUUIDs, 1)
	s.Equal(partUUIDs[1].String(), partsNotFoundErr.MissingUUIDs[0])
	s.Empty(order)
}

func (s *ServiceSuite) TestCreateOrderMultiplePartsNotFound() {
	userUUID := uuid.New()
	partUUIDs := []uuid.UUID{uuid.New(), uuid.New(), uuid.New()}
	parts := []model.Part{
		{
			Uuid:  partUUIDs[0],
			Price: 100.0,
		},
		// partUUIDs[1] and partUUIDs[2] are missing
	}

	s.inventoryClient.On("ListParts", s.ctx, model.PartsFilter{Uuids: partUUIDs}).Return(parts, nil)

	order, err := s.service.CreateOrder(s.ctx, userUUID, partUUIDs)
	s.Error(err)

	var partsNotFoundErr *model.PartsNotFoundError
	s.ErrorAs(err, &partsNotFoundErr)
	s.Len(partsNotFoundErr.MissingUUIDs, 2)
	s.Contains(partsNotFoundErr.MissingUUIDs, partUUIDs[1].String())
	s.Contains(partsNotFoundErr.MissingUUIDs, partUUIDs[2].String())
	s.Empty(order)
}

func (s *ServiceSuite) TestCreateOrderAllPartsNotFound() {
	userUUID := uuid.New()
	partUUIDs := []uuid.UUID{uuid.New(), uuid.New()}
	parts := []model.Part{} // no parts found

	s.inventoryClient.On("ListParts", s.ctx, model.PartsFilter{Uuids: partUUIDs}).Return(parts, nil)

	order, err := s.service.CreateOrder(s.ctx, userUUID, partUUIDs)
	s.Error(err)

	var partsNotFoundErr *model.PartsNotFoundError
	s.ErrorAs(err, &partsNotFoundErr)
	s.Len(partsNotFoundErr.MissingUUIDs, 2)
	s.Contains(partsNotFoundErr.MissingUUIDs, partUUIDs[0].String())
	s.Contains(partsNotFoundErr.MissingUUIDs, partUUIDs[1].String())
	s.Empty(order)
}

func (s *ServiceSuite) TestCreateOrderRepositoryError() {
	userUUID := uuid.New()
	partUUIDs := []uuid.UUID{uuid.New()}
	parts := []model.Part{
		{
			Uuid:  partUUIDs[0],
			Price: 100.0,
		},
	}
	repoErr := gofakeit.Error()

	s.inventoryClient.On("ListParts", s.ctx, model.PartsFilter{Uuids: partUUIDs}).Return(parts, nil)
	s.orderRepository.On("CreateOrder", s.ctx, mock.MatchedBy(func(order model.Order) bool {
		return order.UserUUID == userUUID &&
			len(order.PartUuids) == len(partUUIDs) &&
			order.TotalPrice == 100.0 &&
			order.Status == "PENDING_PAYMENT" &&
			order.OrderUUID != uuid.Nil
	})).Return(repoErr)

	order, err := s.service.CreateOrder(s.ctx, userUUID, partUUIDs)
	s.Error(err)
	s.ErrorIs(err, model.ErrOrderCreateFailed)
	s.Empty(order)
}

func (s *ServiceSuite) TestCreateOrderAlreadyExists() {
	userUUID := uuid.New()
	partUUIDs := []uuid.UUID{uuid.New()}
	parts := []model.Part{
		{
			Uuid:  partUUIDs[0],
			Price: 100.0,
		},
	}

	s.inventoryClient.On("ListParts", s.ctx, model.PartsFilter{Uuids: partUUIDs}).Return(parts, nil)
	s.orderRepository.On("CreateOrder", s.ctx, mock.MatchedBy(func(order model.Order) bool {
		return order.UserUUID == userUUID &&
			len(order.PartUuids) == len(partUUIDs) &&
			order.TotalPrice == 100.0 &&
			order.Status == "PENDING_PAYMENT" &&
			order.OrderUUID != uuid.Nil
	})).Return(model.ErrOrderAlreadyExists)

	order, err := s.service.CreateOrder(s.ctx, userUUID, partUUIDs)
	s.Error(err)
	s.ErrorIs(err, model.ErrOrderAlreadyExists)
	s.Empty(order)
}

func (s *ServiceSuite) TestCreateOrderWithZeroPrice() {
	userUUID := uuid.New()
	partUUIDs := []uuid.UUID{uuid.New()}
	parts := []model.Part{
		{
			Uuid:  partUUIDs[0],
			Price: 0.0, // zero price
		},
	}

	s.inventoryClient.On("ListParts", s.ctx, model.PartsFilter{Uuids: partUUIDs}).Return(parts, nil)
	s.orderRepository.On("CreateOrder", s.ctx, mock.MatchedBy(func(order model.Order) bool {
		return order.UserUUID == userUUID &&
			len(order.PartUuids) == len(partUUIDs) &&
			order.TotalPrice == 0.0 &&
			order.Status == "PENDING_PAYMENT" &&
			order.OrderUUID != uuid.Nil
	})).Return(nil)

	order, err := s.service.CreateOrder(s.ctx, userUUID, partUUIDs)
	s.NoError(err)
	s.Equal(0.0, order.TotalPrice)
}

func (s *ServiceSuite) TestCreateOrderWithNegativePrice() {
	userUUID := uuid.New()
	partUUIDs := []uuid.UUID{uuid.New()}
	parts := []model.Part{
		{
			Uuid:  partUUIDs[0],
			Price: -100.0, // negative price
		},
	}

	s.inventoryClient.On("ListParts", s.ctx, model.PartsFilter{Uuids: partUUIDs}).Return(parts, nil)
	s.orderRepository.On("CreateOrder", s.ctx, mock.MatchedBy(func(order model.Order) bool {
		return order.UserUUID == userUUID &&
			len(order.PartUuids) == len(partUUIDs) &&
			order.TotalPrice == -100.0 &&
			order.Status == "PENDING_PAYMENT" &&
			order.OrderUUID != uuid.Nil
	})).Return(nil)

	order, err := s.service.CreateOrder(s.ctx, userUUID, partUUIDs)
	s.NoError(err)
	s.Equal(-100.0, order.TotalPrice)
}

func (s *ServiceSuite) TestCreateOrderWithVeryHighPrice() {
	userUUID := uuid.New()
	partUUIDs := []uuid.UUID{uuid.New()}
	parts := []model.Part{
		{
			Uuid:  partUUIDs[0],
			Price: 999999.99, // very high price
		},
	}

	s.inventoryClient.On("ListParts", s.ctx, model.PartsFilter{Uuids: partUUIDs}).Return(parts, nil)
	s.orderRepository.On("CreateOrder", s.ctx, mock.MatchedBy(func(order model.Order) bool {
		return order.UserUUID == userUUID &&
			len(order.PartUuids) == len(partUUIDs) &&
			order.TotalPrice == 999999.99 &&
			order.Status == "PENDING_PAYMENT" &&
			order.OrderUUID != uuid.Nil
	})).Return(nil)

	order, err := s.service.CreateOrder(s.ctx, userUUID, partUUIDs)
	s.NoError(err)
	s.Equal(999999.99, order.TotalPrice)
}

func (s *ServiceSuite) TestCreateOrderWithManyParts() {
	userUUID := uuid.New()
	partUUIDs := make([]uuid.UUID, 10)
	parts := make([]model.Part, 10)
	var totalPrice float64

	for i := 0; i < 10; i++ {
		partUUIDs[i] = uuid.New()
		price := gofakeit.Price(10, 100)
		parts[i] = model.Part{
			Uuid:  partUUIDs[i],
			Price: price,
		}
		totalPrice += price
	}

	s.inventoryClient.On("ListParts", s.ctx, model.PartsFilter{Uuids: partUUIDs}).Return(parts, nil)
	s.orderRepository.On("CreateOrder", s.ctx, mock.MatchedBy(func(order model.Order) bool {
		return order.UserUUID == userUUID &&
			len(order.PartUuids) == len(partUUIDs) &&
			order.TotalPrice == totalPrice &&
			order.Status == "PENDING_PAYMENT" &&
			order.OrderUUID != uuid.Nil
	})).Return(nil)

	order, err := s.service.CreateOrder(s.ctx, userUUID, partUUIDs)
	s.NoError(err)
	s.Equal(userUUID, order.UserUUID)
	s.Equal(partUUIDs, order.PartUuids)
	s.Equal(totalPrice, order.TotalPrice)
	s.Equal("PENDING_PAYMENT", order.Status)
}

func (s *ServiceSuite) TestCreateOrderWithSameUserAndOrderUUID() {
	sharedUUID := uuid.New()
	partUUIDs := []uuid.UUID{uuid.New()}
	parts := []model.Part{
		{
			Uuid:  partUUIDs[0],
			Price: 100.0,
		},
	}

	s.inventoryClient.On("ListParts", s.ctx, model.PartsFilter{Uuids: partUUIDs}).Return(parts, nil)
	s.orderRepository.On("CreateOrder", s.ctx, mock.MatchedBy(func(order model.Order) bool {
		return order.UserUUID == sharedUUID &&
			len(order.PartUuids) == len(partUUIDs) &&
			order.TotalPrice == 100.0 &&
			order.Status == "PENDING_PAYMENT" &&
			order.OrderUUID != uuid.Nil
	})).Return(nil)

	order, err := s.service.CreateOrder(s.ctx, sharedUUID, partUUIDs)
	s.NoError(err)
	s.Equal(sharedUUID, order.UserUUID)
	s.NotEqual(sharedUUID, order.OrderUUID) // OrderUUID should be different
}

func (s *ServiceSuite) TestCreateOrderWithDuplicatePartUUIDs() {
	userUUID := uuid.New()
	duplicateUUID := uuid.New()
	partUUIDs := []uuid.UUID{duplicateUUID, duplicateUUID} // duplicate UUIDs
	parts := []model.Part{
		{
			Uuid:  duplicateUUID,
			Price: 100.0,
		},
		// Only one part returned for duplicate UUID
	}

	s.inventoryClient.On("ListParts", s.ctx, model.PartsFilter{Uuids: partUUIDs}).Return(parts, nil)
	s.orderRepository.On("CreateOrder", s.ctx, mock.MatchedBy(func(order model.Order) bool {
		return order.UserUUID == userUUID &&
			len(order.PartUuids) == len(partUUIDs) &&
			order.TotalPrice == 100.0 &&
			order.Status == "PENDING_PAYMENT" &&
			order.OrderUUID != uuid.Nil
	})).Return(nil)

	order, err := s.service.CreateOrder(s.ctx, userUUID, partUUIDs)
	s.NoError(err)
	s.Equal(100.0, order.TotalPrice)
}

func (s *ServiceSuite) TestCreateOrderWithMixedPrices() {
	userUUID := uuid.New()
	partUUIDs := []uuid.UUID{uuid.New(), uuid.New(), uuid.New()}
	parts := []model.Part{
		{
			Uuid:  partUUIDs[0],
			Price: 50.0,
		},
		{
			Uuid:  partUUIDs[1],
			Price: -25.0, // negative price
		},
		{
			Uuid:  partUUIDs[2],
			Price: 0.0, // zero price
		},
	}

	s.inventoryClient.On("ListParts", s.ctx, model.PartsFilter{Uuids: partUUIDs}).Return(parts, nil)
	s.orderRepository.On("CreateOrder", s.ctx, mock.MatchedBy(func(order model.Order) bool {
		return order.UserUUID == userUUID &&
			len(order.PartUuids) == len(partUUIDs) &&
			order.TotalPrice == 25.0 &&
			order.Status == "PENDING_PAYMENT" &&
			order.OrderUUID != uuid.Nil
	})).Return(nil)

	order, err := s.service.CreateOrder(s.ctx, userUUID, partUUIDs)
	s.NoError(err)
	s.Equal(25.0, order.TotalPrice)
}

func (s *ServiceSuite) TestCreateOrderGeneratesUniqueOrderUUID() {
	userUUID := uuid.New()
	partUUIDs := []uuid.UUID{uuid.New()}
	parts := []model.Part{
		{
			Uuid:  partUUIDs[0],
			Price: 100.0,
		},
	}

	s.inventoryClient.On("ListParts", s.ctx, model.PartsFilter{Uuids: partUUIDs}).Return(parts, nil)
	s.orderRepository.On("CreateOrder", s.ctx, mock.MatchedBy(func(order model.Order) bool {
		return order.UserUUID == userUUID &&
			len(order.PartUuids) == len(partUUIDs) &&
			order.TotalPrice == 100.0 &&
			order.Status == "PENDING_PAYMENT" &&
			order.OrderUUID != uuid.Nil
	})).Return(nil)

	order1, err1 := s.service.CreateOrder(s.ctx, userUUID, partUUIDs)
	s.NoError(err1)

	order2, err2 := s.service.CreateOrder(s.ctx, userUUID, partUUIDs)
	s.NoError(err2)

	s.NotEqual(order1.OrderUUID, order2.OrderUUID)
}
