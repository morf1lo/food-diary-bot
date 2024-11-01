package handler

import (
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/morf1lo/food-diary-bot/model"
	"github.com/sirupsen/logrus"
)

const maxMessageLength = 4096

var replyMarkup = tgbotapi.NewReplyKeyboard(
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("/add"),
	),
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("/month"),
		tgbotapi.NewKeyboardButton("/week"),
		tgbotapi.NewKeyboardButton("/day"),
	),
)

func (h *Handler) Command(tgbot *tgbotapi.BotAPI, update tgbotapi.Update) {
	command := update.Message.Command()

	if _, ok := h.lastCommands[update.Message.From.ID][command]; !ok {
		h.lastCommands[update.Message.From.ID] = make(map[string]time.Time)
	}

	if command == "day" || command == "week" || command == "month" || command == "search" {
		if lastTime, ok := h.lastCommands[update.Message.From.ID][command]; ok {
			if time.Since(lastTime) < cooldownDuration {
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "You are using commands too quickly. Please wait before using this command again")
				tgbot.Send(msg)
				return
			}
		}
		h.lastCommands[update.Message.From.ID][command] = time.Now()
	}

	if command == "start" {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Hello!")
		msg.ReplyMarkup = replyMarkup
		tgbot.Send(msg)
	}

	if command == "help" {
		msgText := `
Commands:
/add - Create a record

/month - All records made within 1 month
/week  - All records made within 1 week
/day All records made on this day
/last - Your last meal

/search <query> - Search records that contains query
		`
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, msgText)
		tgbot.Send(msg)
	}

	if command == "add" {
		if err := h.services.Record.RequestToAdd(update.Message.From.ID); err != nil {
			logrus.Fatalf("error creating redis record: %s", err.Error())
		}

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Send a message")
		tgbot.Send(msg)
	}

	if command == "month" {
		records, err := h.services.Record.FindWithinMonth(update.Message.From.ID)
		if err != nil {
			return
		}

		text := formatText(records)
		sendLongMessage(update.Message.Chat.ID, text, tgbot)
	}

	if command == "week" {
		records, err := h.services.Record.FindWithinWeek(update.Message.From.ID)
		if err != nil {
			return
		}

		text := formatText(records)
		sendLongMessage(update.Message.Chat.ID, text, tgbot)
	}

	if command == "day" {
		records, err := h.services.Record.FindWithinDay(update.Message.From.ID)
		if err != nil {
			return
		}

		text := formatText(records)
	
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, text)
		tgbot.Send(msg)
	}

	if command == "last" {
		record, err := h.services.Record.FindLast(update.Message.From.ID)
		if err != nil {
			return
		}

		text := record.DateAdded.Format("02.01.2006") + "\n" + record.Body + " - " + record.DateAdded.Format("15:04")

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, text)
		tgbot.Send(msg)
	}

	if command == "search" {
		query := strings.TrimSpace(strings.Split(update.Message.Text, "/search")[1])
		if query == "" {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Please type query. Command should look like: /search <query>")
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

func sendLongMessage(chatID int64, text string, bot *tgbotapi.BotAPI) {
	for len(text) > maxMessageLength {
		part := text[:maxMessageLength]
		
		lastSpace := strings.LastIndex(part, "\n")
		if lastSpace == -1 {
			lastSpace = strings.LastIndex(part, " ")
		}
		if lastSpace != -1 {
			part = text[:lastSpace]
		}

		msg := tgbotapi.NewMessage(chatID, part)
		bot.Send(msg)

		text = text[len(part):]
	}

	if len(text) > 0 {
		msg := tgbotapi.NewMessage(chatID, text)
		bot.Send(msg)
	}
}
