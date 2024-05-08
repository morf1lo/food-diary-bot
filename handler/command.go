package handler

import (
	"strconv"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/morf1lo/food-diary-bot/model"
	"github.com/sirupsen/logrus"
)

func (h *Handler) Command(tgbot *tgbotapi.BotAPI, update tgbotapi.Update) {
	if update.Message.Command() == "start" {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Hello!")
		tgbot.Send(msg)
	}

	if update.Message.Command() == "help" {
		msgText := `
Commands:
/add - Create a record

/month - All records made within 1 month
/week  - All records made within 1 week
/day All records made on this day
/all - All your records
/last - Your last meal

/search <query> - Search records that contains query
		`
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, msgText)
		tgbot.Send(msg)
	}

	if update.Message.Command() == "add" {
		if err := h.services.Record.RequestToAdd(strconv.Itoa(int(update.Message.From.ID))); err != nil {
			logrus.Fatalf("error creating redis record: %s", err.Error())
		}

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Send a message")
		tgbot.Send(msg)
	}

	if update.Message.Command() == "month" {
		records, err := h.services.Record.FindByMonth(update.Message.From.ID)
		if err != nil {
			return
		}

		text := formatText(records)
	
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, text)
		tgbot.Send(msg)
	}

	if update.Message.Command() == "week" {
		records, err := h.services.Record.FindByWeek(update.Message.From.ID)
		if err != nil {
			return
		}

		text := formatText(records)
	
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, text)
		tgbot.Send(msg)
	}

	if update.Message.Command() == "day" {
		records, err := h.services.Record.FindByDay(update.Message.From.ID)
		if err != nil {
			return
		}

		text := formatText(records)
	
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, text)
		tgbot.Send(msg)
	}

	if update.Message.Command() == "all" {
		records, err := h.services.Record.FindAll(update.Message.From.ID)
		if err != nil {
			return
		}

		text := formatText(records)

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, text)
		tgbot.Send(msg)
	}

	if update.Message.Command() == "last" {
		record, err := h.services.Record.FindLast(update.Message.From.ID)
		if err != nil {
			return
		}

		text := record.DateAdded.Format("02.01.2006") + "\n" + record.Body + " - " + record.DateAdded.Format("15:04")

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, text)
		tgbot.Send(msg)
	}

	if update.Message.Command() == "search" {
		query := strings.TrimSpace(strings.Split(update.Message.Text, "/search")[1])
		if query == "" {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Please type query")
			tgbot.Send(msg)
			return
		}

		records, err := h.services.Record.Search(update.Message.From.ID, query)
		if err != nil {
			return
		}
		if len(records) == 0 {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Records not found")
			tgbot.Send(msg)
			return
		}

		text := formatText(records)

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, text)
		tgbot.Send(msg)
	}
}

func formatText(records []*model.Record) string {
	var text string
	var currentDate time.Time
	for _, record := range records {
		if !record.DateAdded.Truncate(time.Hour * 24).Equal(currentDate) {
			text += "\n" + record.DateAdded.Format("02.01.2006") + "\n"
			currentDate = record.DateAdded.Truncate(time.Hour * 24)
		}
		text += record.Body + " - " + record.DateAdded.Format("15:04") + "\n"
	}

	return text
}
