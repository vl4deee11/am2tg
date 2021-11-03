package main

import (
	"am2tg/pkg/api"
	"am2tg/pkg/log"
	"am2tg/pkg/tg"
	"fmt"
	"github.com/kelseyhightower/envconfig"
	"net/http"
	"regexp"
)

type Config struct {
	API struct {
		Host string `default:"0.0.0.0"`
		Port uint16 `default:"80"`
	} `split_words:"true"`
	Socks5Proxy string
	Token       string `required:"true"`
	LogLvL      string `default:"INFO"`
}

func main() {
	var c Config
	if err := envconfig.Process("AM2TG", &c); err != nil {
		fmt.Println(err)
		return
	}
	log.MakeLogger(c.LogLvL)

	if err := tg.MakeBot(c.Token, c.Socks5Proxy); err != nil {
		log.Logger.Error(err)
		return
	}

	http.HandleFunc("/", route)
	log.Logger.Info("service start")
	log.Logger.Fatal(http.ListenAndServe(fmt.Sprintf("%s:%d", c.API.Host, c.API.Port), nil))
}

var rAlerts = regexp.MustCompile(`/alerts/.*`)

type WriterWithStatusCode struct {
	http.ResponseWriter
	StatusCode int
}

func NewWriter(w http.ResponseWriter) *WriterWithStatusCode {
	return &WriterWithStatusCode{w, http.StatusOK}
}

func (w *WriterWithStatusCode) WriteHeader(code int) {
	w.StatusCode = code
	w.ResponseWriter.WriteHeader(code)
}

func route(w http.ResponseWriter, r *http.Request) {
	lw := NewWriter(w)
	defer log.Logger.Infof(
		"%s %s%s [%d]",
		r.Method,
		r.Host,
		r.URL.String(),
		lw.StatusCode,
	)
	switch {
	case r.URL.Path == "/health" || r.URL.Path == "/ready":
		w.WriteHeader(http.StatusOK)
		_, _ = fmt.Fprintf(w, "ok")
		return
	case rAlerts.MatchString(r.URL.Path):
		api.AlertsPOST(w, r)
		return
	default:
		w.WriteHeader(http.StatusNotFound)
		_, _ = fmt.Fprintf(w, "404 not found")
	}
}
