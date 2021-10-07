package main

import (
	"am2tg/api"
	"am2tg/log"
	"am2tg/tg"
	"fmt"
	"net/http"
	"regexp"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	API struct {
		Host string `default:""`
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
		log.Logger.Fatal(err)
	}

	http.HandleFunc("/", route)
	log.Logger.Fatal(http.ListenAndServe(fmt.Sprintf("%s:%d", c.API.Host, c.API.Port), nil))
}

var rAlerts = regexp.MustCompile(`/alerts/.*`)

func route(w http.ResponseWriter, r *http.Request) {
	lw := log.NewLogWriter(w)
	defer log.Logger.Printf(
		log.Info,
		"%s %s%s [%d]",
		r.Method,
		r.Host,
		r.URL.String(),
		lw.StatusCode,
	)
	switch {
	case r.URL.Path == "/health" || r.URL.Path == "/ready":
		w.WriteHeader(http.StatusOK)
		_, _ = fmt.Fprintf(w, "Ok")
		return
	case rAlerts.MatchString(r.URL.Path):
		api.AlertsPOST(w, r)
		return
	default:
		w.WriteHeader(http.StatusNotFound)
		_, _ = fmt.Fprintf(w, "404 Not found")
	}
}
