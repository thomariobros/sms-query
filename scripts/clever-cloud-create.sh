#!/bin/bash
set -e

clever create --type go sms-query --region par
clever scale --flavor pico --min-instances 1 --max-instances 1 --build-flavor disabled
clever config update --enable-zero-downtime --enable-force-https

# go build env variables https://www.clever-cloud.com/doc/deploy/application/golang/go/#build-the-application
clever env set CC_GO_BUILD_TOOL gomod
clever env set CC_GO_PKG ./cmd/sms-query/server.go
# sms-query env variables
clever env set SMS_GATEWAY_ALLOWED_PHONE_NUMBERS_1 "$PHONE_NUMBER"
clever env set NEXMO_KEY "$NEXMO_KEY"
clever env set NEXMO_SECRET "$NEXMO_SECRET"
clever env set NEXMO_SIGNATURE_SECRET "$NEXMO_SIGNATURE_SECRET"
clever env set NEXMO_PHONE_NUMBER "$NEXMO_PHONE_NUMBER"
clever env set CYCLOCITY_JCDECAUX_API_KEY "$CYCLOCITY_JCDECAUX_API_KEY"

clever deploy
