package handler

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/morf1lo/food-diary-bot/model"
	"github.com/morf1lo/food-diary-bot/service"
)

func (h *Handler) Message(tgbot *tgbotapi.BotAPI, update tgbotapi.Update) {
	record := &model.Record{
		UserID: update.Message.From.ID,
		Body: update.Message.Text,
	}

	err := h.services.Record.Create(record)
	if err == nil {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Record added successfully!")
		tgbot.Send(msg)
		return
	}

	if err != service.ErrDailyLimitReached {
		return
	}

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, err.Error())
	tgbot.Send(msg)
}
