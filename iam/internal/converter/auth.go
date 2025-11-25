package converter

import (
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/nkolesnikov999/micro2-OK/iam/internal/model"
	commonV1 "github.com/nkolesnikov999/micro2-OK/shared/pkg/proto/common/v1"
)

func ToProtoSession(session model.Session) *commonV1.Session {
	return &commonV1.Session{
		Uuid:      session.UUID.String(),
		CreatedAt: timestamppb.New(session.CreatedAt),
		UpdatedAt: timestamppb.New(session.UpdatedAt),
		ExpiresAt: timestamppb.New(session.ExpiresAt),
	}
}
