package tg

import (
	"gopkg.in/telegram-bot-api.v4"
)

var bot *tgbotapi.BotAPI

func GetTGBot() *tgbotapi.BotAPI {
	if bot != nil {
		return bot
	}
	panic("Bot not initialized")
}

func MakeBot(t string) error {
	_bot, err := tgbotapi.NewBotAPI(t)
	if err != nil {
		return err
	}
	bot = _bot
	return nil
}
