package converter

import (
	"github.com/H1dEx/ms-rocket/iam/internal/model"
	repoModel "github.com/H1dEx/ms-rocket/iam/internal/repository/model"
)

func notificationMethodToModel(repoNotificationMethod repoModel.NotificationMethod) model.NotificationMethod {
	return model.NotificationMethod{
		ProviderName: repoNotificationMethod.ProviderName,
		Endpoint:     repoNotificationMethod.Endpoint,
	}
}

func notificationMethodsToModel(repoNotificationMethods []repoModel.NotificationMethod) []model.NotificationMethod {
	notificationMethods := make([]model.NotificationMethod, 0, len(repoNotificationMethods))
	for _, repoNotificationMethod := range repoNotificationMethods {
		notificationMethods = append(notificationMethods, notificationMethodToModel(repoNotificationMethod))
	}
	return notificationMethods
}

func UserToModel(repoUser repoModel.User) model.User {
	return model.User{
		UUID:                repoUser.UUID,
		Login:               repoUser.Login,
		Email:               repoUser.Email,
		NotificationMethods: notificationMethodsToModel(repoUser.NotificationMethods),
		PasswordHash:        repoUser.Password,
		CreatedAt:           repoUser.CreatedAt,
		UpdatedAt:           repoUser.UpdatedAt,
	}
}

func notificationMethodToRepo(notificationMethod model.NotificationMethod) repoModel.NotificationMethod {
	return repoModel.NotificationMethod{
		ProviderName: notificationMethod.ProviderName,
		Endpoint:     notificationMethod.Endpoint,
	}
}

func notificationMethodsToRepo(notificationMethods []model.NotificationMethod) []repoModel.NotificationMethod {
	repoNotificationMethods := make([]repoModel.NotificationMethod, 0, len(notificationMethods))
	for _, notificationMethod := range notificationMethods {
		repoNotificationMethods = append(repoNotificationMethods, notificationMethodToRepo(notificationMethod))
	}
	return repoNotificationMethods
}

func UserToRepo(user model.User) repoModel.User {
	return repoModel.User{
		UUID:                user.UUID,
		Login:               user.Login,
		Email:               user.Email,
		NotificationMethods: notificationMethodsToRepo(user.NotificationMethods),
		Password:            user.PasswordHash,
		CreatedAt:           user.CreatedAt,
		UpdatedAt:           user.UpdatedAt,
	}
}
