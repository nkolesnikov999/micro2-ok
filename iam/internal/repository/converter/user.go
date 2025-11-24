package converter

import (
	"encoding/json"

	"github.com/nkolesnikov999/micro2-OK/iam/internal/model"
	repoModel "github.com/nkolesnikov999/micro2-OK/iam/internal/repository/model"
)

// ToRepoUser конвертирует доменную модель пользователя в репозиторную модель для PostgreSQL.
func ToRepoUser(user model.User) (repoModel.User, error) {
	var notificationMethodsJSON []byte
	if len(user.Info.NotificationMethods) > 0 {
		var err error
		notificationMethodsJSON, err = json.Marshal(user.Info.NotificationMethods)
		if err != nil {
			return repoModel.User{}, err
		}
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

// ToRepoUserWithPasswordHash конвертирует доменную модель пользователя в репозиторную модель с установленным passwordHash.
func ToRepoUserWithPasswordHash(user model.User, passwordHash string) (repoModel.User, error) {
	repoUser, err := ToRepoUser(user)
	if err != nil {
		return repoModel.User{}, err
	}
	repoUser.PasswordHash = passwordHash
	return repoUser, nil
}

// ToModelUser конвертирует репозиторную модель пользователя из PostgreSQL в доменную модель.
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
