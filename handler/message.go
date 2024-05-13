package handler

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/morf1lo/food-diary-bot/model"
)

func (h *Handler) Message(tgbot *tgbotapi.BotAPI, update tgbotapi.Update) {
	record := &model.Record{
		UserID: update.Message.From.ID,
		Body: update.Message.Text,
	}

	if err := h.services.Record.Create(record); err != nil {
		return
	}

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Record added successfully!")
	tgbot.Send(msg)
}
