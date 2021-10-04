[![Go](https://github.com/thomariobros/sms-query/actions/workflows/go.yml/badge.svg)](https://github.com/thomariobros/sms-query/actions/workflows/go.yml)

# Vonage APIs configuration

1. https://dashboard.nexmo.com/your-numbers => buy virtual phone number
2. https://dashboard.nexmo.com/settings => activate *Messages API*
3. https://dashboard.nexmo.com/applications => new application, activate *Messages* capability and set *Inbound URL*
4. link virtual phone number to the new application

# local dev

## Build

```bash
make
```

## cli to test query and print response or send response as SMS

```bash
make build-cli
make run-cli from=... query=help send=false # print
make run-cli from=... query=help send=true # send SMS
```

## Unit tests

```bash
make test
```

## pre-commit

see config file [.pre-commit-config.yaml](.pre-commit-config.yaml)

### Install

```bash
pip3 install pre-commit
pre-commit install
```

## Vonage APIs webhooks

- https://developer.nexmo.com/messages/code-snippets/configure-webhooks
- https://developer.nexmo.com/messages/concepts/signed-webhooks

# Deploy

## Clever Cloud https://www.clever-cloud.com/

### Create
```bash
export PHONE_NUMBER=...
export NEXMO_KEY=...
export NEXMO_SECRET=...
export NEXMO_SIGNATURE_SECRET=...
export NEXMO_PHONE_NUMBER=...
export CLEVER_TOKEN=...
export CLEVER_SECRET=...
./scripts/clever-cloud-create.sh
```

### Update
```bash
export CLEVER_TOKEN=...
export CLEVER_SECRET=...
clever deploy
```

# External data APIs

## JCDecaux’s self-service bicycles (bicloo in Nantes, France)

https://developer.jcdecaux.com/

## Transport

http://data.nantes.fr/donnees/detail/info-trafic-temps-reel-de-la-tan/

https://opendata.stif.info/page/home/

https://www.vianavigo.com/accueil

http://doc.navitia.io/#next-departures-and-arrivals

https://canaltp.github.io/navitia-playground/

## Weather

https://darksky.net/dev/login

https://darksky.net/forecast/47.2186,-1.5541/us12/en
