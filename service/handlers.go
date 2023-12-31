package service

import (
	"fmt"

	"report-bot/doc"

	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Handler struct {
	*tg.BotAPI
	data *doc.Data
	doc  *doc.Doc
}

func New(botToken string) (*Handler, error) {
	bot, err := tg.NewBotAPI(botToken)
	if err != nil {
		return nil, fmt.Errorf("failed to create bot API: %s", err.Error())
	}

	return &Handler{
		bot,
		&doc.Data{},
		&doc.Doc{},
	}, nil
}

func (h *Handler) Run() error {
	// if err := h.DeleteDocument(); err != nil {
	// 	return fmt.Errorf("failed to delete document: %s", err.Error())
	// }

	upd := tg.NewUpdate(0)
	upd.Timeout = 60
	//h.Debug = true
	upds := h.GetUpdatesChan(upd)

	for u := range upds {
		switch true {
		case u.Message != nil && u.Message.Text == "/start":
			h.data = &doc.Data{}
			if err := h.Start(u.Message.Chat.ID); err != nil {
				return fmt.Errorf("error in 'start' func: %s", err.Error())
			}

		case u.Message == nil && u.CallbackQuery != nil:
			if u.CallbackData() == "/create" {
				Mode = "/create"
			} else if u.CallbackData() == "/list" {
				Mode = "/list"
			}

			if err := h.Next(u.CallbackQuery.Message.Chat.ID, u.CallbackData()); err != nil {
				return fmt.Errorf("error in 'next' func: %s", err.Error())
			}

			_, _ = h.Request(tg.NewDeleteMessage(u.CallbackQuery.Message.Chat.ID, u.CallbackQuery.Message.MessageID))

		case u.Message != nil && u.Message.Text != "/start":
			if err := h.Next(u.Message.Chat.ID, u.Message.Text); err != nil {
				return fmt.Errorf("error in 'next' func: %s", err.Error())
			}
			_, _ = h.Request(tg.NewDeleteMessage(u.Message.Chat.ID, u.Message.MessageID-1))
		}
	}

	return nil
}
