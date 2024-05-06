package handler

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/morf1lo/food-diary-bot/configs"
)

type Bot struct {
	handler *Handler
}

func NewBot(handler *Handler) *Bot {
	return &Bot{handler: handler}
}

func (b *Bot) Start(cfg *configs.TgBotConfig) error {
	tgbot, err := tgbotapi.NewBotAPI(cfg.Token)
	if err != nil {
		return err
	}

	tgbot.Debug = cfg.Debug

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := tgbot.GetUpdatesChan(u)
	for update := range updates {
		if update.Message != nil && update.Message.IsCommand() {
			go b.handler.Command(tgbot, update)
		}

		if update.Message != nil && !update.Message.IsCommand() {
			go b.handler.Message(tgbot, update)
		}
	}

	return nil
}
