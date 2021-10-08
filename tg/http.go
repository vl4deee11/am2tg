package tg

import (
	"am2tg/log"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"golang.org/x/net/context"
)

const (
	APIReqTimeout = 30 * time.Second
	getMe         = "getMe"
	sendMsg       = "sendMessage"
	fmtAPI        = "https://api.telegram.org/bot%s/%s"
)

type Bot struct {
	token string
	cli   *http.Client
}

type Resp struct {
	Ok          bool            `json:"ok"`
	Result      json.RawMessage `json:"result"`
	ErrorCode   int             `json:"error_code"`
	Description string          `json:"description"`
	Parameters  struct {
		MigrateToChatID int64 `json:"migrate_to_chat_id"`
		RetryAfter      int   `json:"retry_after"`
	} `json:"parameters"`
}

func (b *Bot) makeHTTPReq(method string, params url.Values) error {
	ctx, _ := context.WithTimeout(context.Background(), APIReqTimeout)
	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		fmt.Sprintf(fmtAPI, b.token, method),
		strings.NewReader(params.Encode()),
	)
	if err != nil {
		return err
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	resp, err := b.cli.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Logger.Warn("cannot read response body")
		return nil
	}

	log.Logger.Trace(data)

	var tgResp Resp
	err = json.Unmarshal(data, &tgResp)
	if err != nil {
		return err
	}

	if !tgResp.Ok {
		return fmt.Errorf("code:[%d]:%s", tgResp.ErrorCode, tgResp.Description)
	}

	return nil
}

func (b *Bot) ping() error {
	return b.makeHTTPReq(getMe, nil)
}

func (b *Bot) SendMsg(chatID int64, txt string) error {
	v := url.Values{}
	v.Add("chat_id", strconv.FormatInt(chatID, 10))
	v.Add("text", txt)
	return b.makeHTTPReq(sendMsg, v)
}
