package api

import (
	"am2tg/log"
	"am2tg/tg"
	"encoding/json"
	"fmt"
	tgbotapi "gopkg.in/telegram-bot-api.v4"
	"io/ioutil"
	"net/http"
	"sort"
	"strconv"
	"strings"
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
	if err := json.Unmarshal(body, &alerts); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Logger.Print(log.Error, err.Error())
		return
	}

	sli := strings.Split(r.RequestURI, "/")
	chatId, err := strconv.Atoi(sli[len(sli)-1])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Logger.Printf(log.Error, "cannot convert chatId=%s ot int:%s", sli[len(sli)-1], err.Error())
		return
	}
	log.Logger.Printf(log.Debug, "get chat id = %d", chatId)
	bot := tg.GetTGBot()
	msg := tgbotapi.NewMessage(int64(chatId), alerts.format())
	msg.ParseMode = tgbotapi.ModeHTML
	if _, err := bot.Send(msg); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Logger.Println(log.Error, err.Error())
		return
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

func (alerts *Alerts) format() string {
	keys := make([]string, 0, len(alerts.GroupLabels))
	for k := range alerts.GroupLabels {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	groupLabels := make([]string, 0, len(alerts.GroupLabels))
	for _, k := range keys {
		groupLabels = append(groupLabels, fmt.Sprintf("%s=<code>%s</code>", k, alerts.GroupLabels[k]))
	}

	keys = make([]string, 0, len(alerts.CommonLabels))
	for k := range alerts.CommonLabels {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	commonLabels := make([]string, 0, len(alerts.CommonLabels))
	for _, k := range keys {
		if _, ok := alerts.GroupLabels[k]; !ok {
			commonLabels = append(commonLabels, fmt.Sprintf("%s=<code>%s</code>", k, alerts.CommonLabels[k]))
		}
	}

	keys = make([]string, 0, len(alerts.CommonAnnotations))
	for k := range alerts.CommonAnnotations {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	commonAnnotations := make([]string, 0, len(alerts.CommonAnnotations))
	for _, k := range keys {
		commonAnnotations = append(commonAnnotations, fmt.Sprintf("\n%s: <code>%s</code>", k, alerts.CommonAnnotations[k]))
	}

	alertDetails := make([]string, len(alerts.Alerts))
	for i, a := range alerts.Alerts {
		if instance, ok := a.Labels["instance"]; ok {
			instanceString, _ := instance.(string)
			alertDetails[i] += strings.Split(instanceString, ":")[0]
		}
		if job, ok := a.Labels["job"]; ok {
			alertDetails[i] += fmt.Sprintf("[%s]", job)
		}
		if a.GeneratorURL != "" {
			alertDetails[i] = fmt.Sprintf("<a href='%s'>%s</a>", a.GeneratorURL, a.GeneratorURL)
		}
	}
	return fmt.Sprintf(
		"<a href='%s/#/alerts?receiver=%s'>[%s:%d]</a>\ngrouped by: %s\nlabels: %s%s\n%s",
		alerts.ExternalURL,
		alerts.Receiver,
		strings.ToUpper(alerts.Status),
		len(alerts.Alerts),
		strings.Join(groupLabels, ", "),
		strings.Join(commonLabels, ", "),
		strings.Join(commonAnnotations, ""),
		strings.Join(alertDetails, ", "),
	)
}
