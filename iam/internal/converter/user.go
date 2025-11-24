package converter

import (
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/nkolesnikov999/micro2-OK/iam/internal/model"
	commonV1 "github.com/nkolesnikov999/micro2-OK/shared/pkg/proto/common/v1"
)

// ToProtoUser конвертирует доменную модель User в proto сообщение.
func ToProtoUser(user model.User) *commonV1.User {
	protoNotificationMethods := make([]*commonV1.NotificationMethod, 0, len(user.Info.NotificationMethods))
	for _, nm := range user.Info.NotificationMethods {
		protoNotificationMethods = append(protoNotificationMethods, ToProtoNotificationMethod(nm))
	}

	return &commonV1.User{
		Uuid: user.UUID.String(),
		Info: &commonV1.UserInfo{
			Login:               user.Info.Login,
			Email:               user.Info.Email,
			NotificationMethods: protoNotificationMethods,
		},
		CreatedAt: timestamppb.New(user.CreatedAt),
		UpdatedAt: timestamppb.New(user.UpdatedAt),
	}
}

// ToProtoNotificationMethod конвертирует доменную модель NotificationMethod в proto сообщение.
func ToProtoNotificationMethod(nm model.NotificationMethod) *commonV1.NotificationMethod {
	return &commonV1.NotificationMethod{
		ProviderName: nm.ProviderName,
		Target:       nm.Target,
	}
}
