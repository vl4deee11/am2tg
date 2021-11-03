package tg

import (
	"net/http"

	"golang.org/x/net/proxy"
)

var bot *Bot

func GetTGBot() *Bot {
	if bot != nil {
		return bot
	}
	panic("Bot not initialized")
}

func MakeBot(t, proxyURL string) error {
	if proxyURL != "" {
		dialSocksProxy, err := proxy.SOCKS5("tcp", proxyURL, nil, proxy.Direct)
		if err != nil {
			return err
		}
		bot = &Bot{
			token: t,
			cli: &http.Client{
				Transport: &http.Transport{
					Dial: dialSocksProxy.Dial,
				},
			},
		}

		if err := bot.ping(); err != nil {
			return err
		}

		return nil
	}
	bot = &Bot{token: t, cli: new(http.Client)}
	if err := bot.ping(); err != nil {
		return err
	}
	return nil
}
