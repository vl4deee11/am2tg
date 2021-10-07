package api

import (
	"am2tg/log"
	"am2tg/tg"
	"encoding/json"
	"fmt"
	tgbotapi "gopkg.in/telegram-bot-api.v4"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

type Alerts struct {
	Alerts            []Alert                `json:"alerts"`
	CommonAnnotations map[string]interface{} `json:"commonAnnotations"`
	CommonLabels      map[string]interface{} `json:"commonLabels"`
	ExternalURL       string                 `json:"externalURL"`
	GroupKey          string                 `json:"groupKey"`
	GroupLabels       map[string]interface{} `json:"groupLabels"`
	Receiver          string                 `json:"receiver"`
	Status            string                 `json:"status"`
	Version           string                 `json:"version"`
}

type Alert struct {
	Annotations  map[string]interface{} `json:"annotations"`
	EndsAt       string                 `json:"endsAt"`
	GeneratorURL string                 `json:"generatorURL"`
	Labels       map[string]interface{} `json:"labels"`
	StartsAt     string                 `json:"startsAt"`
	Status       string                 `json:"status"`
}

func AlertsPOST(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprintf(w, "Sorry, only POST methods are supported")
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Logger.Print(log.Error, err.Error())
		return
	}

	var alerts Alerts
	if err := json.Unmarshal(body, &alerts); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Logger.Print(log.Error, err.Error())
		return
	}

	sli := strings.Split(r.RequestURI, "/")
	chatId, err := strconv.Atoi(sli[len(sli)-1])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Logger.Printf(log.Error, "cannot convert chatId=%s ot int:", sli[len(sli)-1], err.Error())
		return
	}

	log.Logger.Printf(log.Debug, "get chat id = %s", chatId)
	msg := tgbotapi.NewMessage(int64(chatId), string(body))
	if _, err := tg.GetTGBot().Send(msg); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Logger.Println(log.Error, err.Error())
		return
	}
	log.Logger.Println(log.Info, "send alerts successfully")
	w.WriteHeader(http.StatusOK)
}
