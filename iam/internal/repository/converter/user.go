package converter

import (
	"encoding/json"

	"github.com/nkolesnikov999/micro2-OK/iam/internal/model"
	repoModel "github.com/nkolesnikov999/micro2-OK/iam/internal/repository/model"
)

func ToRepoUser(user model.User) (repoModel.User, error) {
	// Всегда маршалим notificationMethods, даже если массив пустой, чтобы получить [] вместо null
	notificationMethodsJSON, err := json.Marshal(user.Info.NotificationMethods)
	if err != nil {
		return repoModel.User{}, err
	}

	return repoModel.User{
		UUID:                user.UUID,
		Login:               user.Info.Login,
		Email:               user.Info.Email,
		PasswordHash:        "", // PasswordHash устанавливается отдельно при создании/обновлении
		NotificationMethods: notificationMethodsJSON,
		CreatedAt:           user.CreatedAt,
		UpdatedAt:           user.UpdatedAt,
	}, nil
}

func ToRepoUserWithPasswordHash(user model.User, passwordHash string) (repoModel.User, error) {
	repoUser, err := ToRepoUser(user)
	if err != nil {
		return repoModel.User{}, err
	}
	repoUser.PasswordHash = passwordHash
	return repoUser, nil
}

func ToModelUser(repoUser repoModel.User) (model.User, error) {
	var notificationMethods []model.NotificationMethod
	if len(repoUser.NotificationMethods) > 0 {
		if err := json.Unmarshal(repoUser.NotificationMethods, &notificationMethods); err != nil {
			return model.User{}, err
		}
	}

	return model.User{
		UUID:      repoUser.UUID,
		CreatedAt: repoUser.CreatedAt,
		UpdatedAt: repoUser.UpdatedAt,
		Info: model.UserInfo{
			Login:               repoUser.Login,
			Email:               repoUser.Email,
			NotificationMethods: notificationMethods,
		},
	}, nil
}
