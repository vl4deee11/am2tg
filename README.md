# am2tg

Service for sending alerts from prometheus alert manager to telegram channel

## AlertManager config

Add to `alertmanager.yml`
More docs [here](https://prometheus.io/docs/alerting/latest/configuration/)

```yaml
receivers:
  - name: 'web.hook'
    webhook_configs:
      - send_resolved: True
        # Get 'chat_id' -> https://web.telegram.org/z/#-{number} -> -100{number}
        url: 'http://0.0.0.0:80/alerts/chat_id'
```

## Config

Available environment variables

`AM2TG_API_HOST` - service hostname

`AM2TG_API_PORT` - service port

`AM2TG_SOCKS5PROXY` - proxy address

`AM2TG_TOKEN` - telegram bot token

`AM2TG_LOGLVL` - log level (`DEBUG` || `INFO` || `WARN` || `ERROR` || `TRACE`)

## API
### GET /health
Health probe

### POST /alerts/:chatid
Push new alert to telegram
