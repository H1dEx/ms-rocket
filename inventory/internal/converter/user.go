package converter

import (
	"time"

	"github.com/H1dEx/ms-rocket/inventory/internal/model"
	common_v1 "github.com/H1dEx/ms-rocket/shared/pkg/proto/common/v1"
)

func NotificationMethodsToModel(methods []*common_v1.NotificationMethod) []model.NotificationMethod {
	res := make([]model.NotificationMethod, 0, len(methods))
	for _, method := range methods {
		res = append(res, model.NotificationMethod{
			ProviderName: method.ProviderName,
			Endpoint:     method.Target,
		})
	}
	return res
}

func UserToModel(user *common_v1.User) model.User {
	if user == nil {
		return model.User{}
	}

	var (
		login               string
		email               string
		createdAt           time.Time
		updated_at          time.Time
		notificationMethods []model.NotificationMethod
	)
	if user.CreatedAt != nil {
		createdAt = user.CreatedAt.AsTime()
	}
	if user.UpdatedAt != nil {
		updated_at = user.UpdatedAt.AsTime()
	}

	if user.Info != nil {
		email = user.Info.Email
		login = user.Info.Login
		notificationMethods = NotificationMethodsToModel(user.Info.NotificationMethods)
	}

	return model.User{
		UUID:                user.Uuid,
		CreatedAt:           createdAt,
		UpdatedAt:           updated_at,
		Login:               login,
		Email:               email,
		NotificationMethods: notificationMethods,
	}
}
