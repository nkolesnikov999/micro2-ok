//go:build integration

package integration

import (
	"context"

	"github.com/brianvoe/gofakeit/v7"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	authV1 "github.com/nkolesnikov999/micro2-OK/shared/pkg/proto/auth/v1"
	commonV1 "github.com/nkolesnikov999/micro2-OK/shared/pkg/proto/common/v1"
	userV1 "github.com/nkolesnikov999/micro2-OK/shared/pkg/proto/user/v1"
)

// CreateTestSession создает тестового пользователя и сессию через IAM сервис
func (env *TestEnvironment) CreateTestSession(ctx context.Context) (string, error) {
	// Генерируем уникальный логин/емейл, чтобы избежать конфликта "user already exists"
	login := gofakeit.Email()

	// Подключаемся к IAM сервису
	iamConn, err := grpc.DialContext(
		ctx,
		env.IAM.Address(),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return "", err
	}
	defer iamConn.Close()

	// Создаем клиенты
	userClient := userV1.NewUserServiceClient(iamConn)
	authClient := authV1.NewAuthServiceClient(iamConn)

	// Регистрируем тестового пользователя
	_, err = userClient.Register(ctx, &userV1.RegisterRequest{
		Info: &userV1.UserRegistrationInfo{
			Info: &commonV1.UserInfo{
				Login: login,
				Email: login,
			},
			Password: "testpassword123",
		},
	})
	if err != nil {
		return "", err
	}

	// Создаем сессию для пользователя
	loginResp, err := authClient.Login(ctx, &authV1.LoginRequest{
		Login:    login,
		Password: "testpassword123",
	})
	if err != nil {
		return "", err
	}

	return loginResp.SessionUuid, nil
}
