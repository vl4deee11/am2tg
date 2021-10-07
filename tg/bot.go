package tg

import (
	"golang.org/x/net/proxy"
	"gopkg.in/telegram-bot-api.v4"
	"net/http"
)

var bot *tgbotapi.BotAPI

func GetTGBot() *tgbotapi.BotAPI {
	if bot != nil {
		return bot
	}
	panic("Bot not initialized")
}

func MakeBot(t string, proxyURL string) error {
	if proxyURL != "" {
		dialSocksProxy, err := proxy.SOCKS5("tcp", proxyURL, nil, proxy.Direct)
		if err != nil {
			return err
		}
		cli := &http.Client{Transport: &http.Transport{Dial: dialSocksProxy.Dial}}
		_bot, err := tgbotapi.NewBotAPIWithClient(t, cli)
		if err != nil {
			return err
		}
		bot = _bot
		return nil
	}
	_bot, err := tgbotapi.NewBotAPI(t)
	if err != nil {
		return err
	}
	bot = _bot
	return nil
}
