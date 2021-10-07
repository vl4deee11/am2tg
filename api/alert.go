package api

import (
	"am2tg/log"
	"am2tg/tg"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sort"
	"strconv"
	"strings"

	tgbotapi "gopkg.in/telegram-bot-api.v4"
)

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
	if err = json.Unmarshal(body, &alerts); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Logger.Print(log.Error, err.Error())
		return
	}

	sli := strings.Split(r.RequestURI, "/")
	chatID, err := strconv.Atoi(sli[len(sli)-1])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Logger.Printf(log.Error, "cannot convert chatId=%s ot int:%s", sli[len(sli)-1], err.Error())
		return
	}
	log.Logger.Printf(log.Debug, "get chat id = %d", chatID)
	bot := tg.GetTGBot()

	chunkedMsg := alerts.format()
	for i := range chunkedMsg {
		msg := tgbotapi.NewMessage(int64(chatID), chunkedMsg[i])
		if _, err := bot.Send(msg); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Logger.Println(log.Error, err.Error())
			_, _ = bot.Send(tgbotapi.NewMessage(int64(chatID), "Error sending message, checkout logs"))
			return
		}
	}
	log.Logger.Println(log.Info, "send alerts successfully")
	w.WriteHeader(http.StatusOK)
}

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

func (alerts *Alerts) format() []string {
	keys := make([]string, 0, len(alerts.GroupLabels))
	for k := range alerts.GroupLabels {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	groupLabels := make([]string, 0, len(alerts.GroupLabels))
	for _, k := range keys {
		groupLabels = append(groupLabels, fmt.Sprintf("%s=%s", k, alerts.GroupLabels[k]))
	}

	keys = make([]string, 0, len(alerts.CommonLabels))
	for k := range alerts.CommonLabels {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	commonLabels := make([]string, 0, len(alerts.CommonLabels))
	for _, k := range keys {
		if _, ok := alerts.GroupLabels[k]; !ok {
			commonLabels = append(commonLabels, fmt.Sprintf("%s=%s", k, alerts.CommonLabels[k]))
		}
	}

	keys = make([]string, 0, len(alerts.CommonAnnotations))
	for k := range alerts.CommonAnnotations {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	commonAnnotations := make([]string, 0, len(alerts.CommonAnnotations))
	for _, k := range keys {
		commonAnnotations = append(commonAnnotations, fmt.Sprintf("\n%s: %s", k, alerts.CommonAnnotations[k]))
	}

	alertDetails := make([]string, len(alerts.Alerts))
	for i := range alerts.Alerts {
		if alerts.Alerts[i].Status == "firing" {
			alertDetails[i] = fmt.Sprintf(
				"Alert[%d]: \n starts_at=%s",
				i,
				alerts.Alerts[i].StartsAt,
			)
		} else {
			alertDetails[i] = fmt.Sprintf(
				"Alert[%d]: \n starts_at= %s \n ends_at=%s",
				i+1,
				alerts.Alerts[i].StartsAt,
				alerts.Alerts[i].EndsAt,
			)
		}
	}
	return chunkMsg(fmt.Sprintf(
		"[%s:%d]\nGrouped by: %s\nLabels: %s%s\n%s",
		strings.ToUpper(alerts.Status),
		len(alerts.Alerts),
		strings.Join(groupLabels, ", "),
		strings.Join(commonLabels, ", "),
		strings.Join(commonAnnotations, ""),
		strings.Join(alertDetails, "\n"),
	))
}

func chunkMsg(s string) []string {
	// TG Api max msg size
	max := 4000

	var sb strings.Builder
	var chunks []string

	runes := bytes.Runes([]byte(s))
	l := len(runes) - 1
	for i := range runes {
		sb.WriteRune(runes[i])
		if sb.Len() == max {
			chunks = append(chunks, sb.String())
			sb.Reset()
		} else if i == l {
			chunks = append(chunks, sb.String())
		}
	}

	return chunks
}
