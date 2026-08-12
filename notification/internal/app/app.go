package app

import (
	"context"

	"github.com/go-faster/errors"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"go.uber.org/zap"

	"github.com/H1dEx/ms-rocket/notification/internal/config"
	"github.com/H1dEx/ms-rocket/platform/pkg/closer"
	"github.com/H1dEx/ms-rocket/platform/pkg/logger"
)

type App struct {
	diContainer *diContainer
}

func New(ctx context.Context) (*App, error) {
	a := &App{}
	if err := a.initDeps(ctx); err != nil {
		return nil, err
	}
	return a, nil
}

func (a *App) initDeps(ctx context.Context) error {
	inits := []func(context.Context) error{
		a.initDI,
		a.InitLogger,
		a.initCloser,
		a.initTelegramBot,
	}
	for _, init := range inits {
		if err := init(ctx); err != nil {
			return err
		}
	}
	return nil
}

func (a *App) Run(ctx context.Context) error {
	errCh := make(chan error, 2)

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	go func() {
		if err := a.runAssembledConsumer(ctx); err != nil {
			errCh <- errors.Errorf("Assembled consumer crashed: %v", err)
		}
	}()
	go func() {
		if err := a.runPaidConsumer(ctx); err != nil {
			errCh <- errors.Errorf("Paid consumer crashed: %v", err)
		}
	}()

	select {
	case err := <-errCh:
		logger.Error(ctx, "💥 Consumer crashed", zap.Error(err))
		cancel()
		return err
	case <-ctx.Done():
		logger.Info(ctx, "🛑 Consumer context cancelled, shutting down")
		return ctx.Err()
	}
}

func (a *App) runAssembledConsumer(ctx context.Context) error {
	logger.Info(ctx, "Starting assembled consumer")
	err := a.diContainer.OrderAssembledConsumer(ctx).RunConsumer(ctx)
	if err != nil {
		logger.Error(ctx, "failed to run assembled consumer", zap.Error(err))
		return err
	}
	return nil
}

func (a *App) runPaidConsumer(ctx context.Context) error {
	logger.Info(ctx, "Starting paid consumer")
	err := a.diContainer.OrderPaidConsumer(ctx).RunConsumer(ctx)
	if err != nil {
		logger.Error(ctx, "failed to run paid consumer", zap.Error(err))
		return err
	}
	return nil
}

func (a *App) initDI(_ context.Context) error {
	a.diContainer = NewDIContainer()
	return nil
}

func (a *App) InitLogger(_ context.Context) error {
	err := logger.Init(config.GetConfig().Logger.Level(), config.GetConfig().Logger.AsJSON())
	if err != nil {
		return err
	}
	return nil
}

func (a *App) initCloser(_ context.Context) error {
	closer.SetLogger(logger.Logger())
	return nil
}

func (a *App) initTelegramBot(ctx context.Context) error {
	telegramBot := a.diContainer.TelegramBot(ctx)
	telegramBot.RegisterHandler(bot.HandlerTypeMessageText, "/start", bot.MatchTypeExact, func(ctx context.Context, b *bot.Bot, update *models.Update) {
		logger.Info(ctx, "chat id", zap.Int64("chat_id", update.Message.Chat.ID))
		_, err := b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   "👋 Привет! Я бот для уведомлений о заказах. Напиши мне /help для получения помощи.",
		})
		if err != nil {
			logger.Error(ctx, "failed to send activation message", zap.Error(err))
		}
	})

	go func() {
		logger.Info(ctx, "Starting telegram bot")
		telegramBot.Start(ctx)
	}()

	return nil
}
