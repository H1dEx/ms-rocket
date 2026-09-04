package converter

import (
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/H1dEx/ms-rocket/iam/internal/model"
	v1 "github.com/H1dEx/ms-rocket/shared/pkg/proto/common/v1"
)

func NotificationMethodsToProto(notificationMethods []model.NotificationMethod) []*v1.NotificationMethod {
	methods := make([]*v1.NotificationMethod, 0, len(notificationMethods))

	for _, method := range notificationMethods {
		methods = append(methods, &v1.NotificationMethod{
			ProviderName: method.ProviderName,
			Target:       method.Endpoint,
		})
	}

	return methods
}

func UserInfoToProto(userInfo model.User) *v1.UserInfo {
	return &v1.UserInfo{
		Login:               userInfo.Login,
		Email:               userInfo.Email,
		NotificationMethods: NotificationMethodsToProto(userInfo.NotificationMethods),
	}
}

func UserToProto(user model.User) *v1.User {
	return &v1.User{
		Uuid:      user.UUID,
		Info:      UserInfoToProto(user),
		CreatedAt: timestamppb.New(user.CreatedAt),
		UpdatedAt: timestamppb.New(user.UpdatedAt),
	}
}

func NotificationMethodsToModel(notificationMethods []*v1.NotificationMethod) []model.NotificationMethod {
	methods := make([]model.NotificationMethod, 0, len(notificationMethods))

	for _, method := range notificationMethods {
		methods = append(methods, model.NotificationMethod{
			ProviderName: method.ProviderName,
			Endpoint:     method.Target,
		})
	}
	return methods
}
