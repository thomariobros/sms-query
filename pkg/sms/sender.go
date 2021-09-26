package sms

import (
	"encoding/json"
	"fmt"
	corehttp "net/http"
	"net/url"
	"time"

	"github.com/avast/retry-go"
	"github.com/golang/glog"

	"main/pkg/config"
	"main/pkg/http"
)

type Sender struct {
	phoneNumber string
	message     string
}

func Send(to string, message string) error {
	// https://help.nexmo.com/hc/en-us/articles/204076866-How-Long-is-a-Single-SMS-body-
	if len(message) > 800 {
		// truncate too long message
		glog.Warningf("truncate sms message %s", message)
		message = message[:800] + "..."
	}
	sender := config.GetInstance().FindCustomSender(to)
	if sender != nil {
		var err error
		switch sender.Provider.ID {
		case config.FREE_MOBILE:
			err = sendFreeMobile(to, message)
		default:
			glog.Warningf("no sms sender for %s so fallback", to)
			err = sendNexmo(to, message)
		}
		if err != nil {
			glog.Errorf("error while sending sms to %s: %s", to, err)
			return err
		}
	} else {
		err := sendNexmo(to, message)
		if err != nil {
			glog.Errorf("error while sending sms to %s: %s", to, err)
			return err
		}
	}
	return nil
}

func sendFreeMobile(to string, message string) error {
	provider := config.GetInstance().FindCustomSender(to).Provider
	params := map[string]string{"user": provider.Key, "pass": provider.Secret, "msg": message}
	body, err := json.Marshal(params)
	if err != nil {
		return err
	}
	var resp *corehttp.Response = nil
	err = retry.Do(
		func() error {
			resp, err = http.Post("https://smsapi.free-mobile.fr/sendmsg", "application/json", string(body))
			if err == nil {
				defer resp.Body.Close()
			}
			return err
		},
		retry.Attempts(5),
		retry.Delay(5),
		retry.Units(time.Second),
		retry.OnRetry(func(n uint, err error) {
			glog.Errorf("attempt #%d: %s\n", n, err)
		}),
	)
	if err != nil {
		return err
	}
	if resp.StatusCode != 200 {
		return fmt.Errorf("error while sending sms through free mobile, http code: %d", resp.StatusCode)
	}
	return nil
}

func sendNexmo(to string, message string) error {
	provider := config.GetInstance().SMSGateway.Provider
	params := "api_key=" + provider.Key +
		"&api_secret=" + provider.Secret +
		"&type=text" +
		"&from=" + provider.PhoneNumber +
		"&to=" + to +
		"&text=" + url.QueryEscape(message)
	var err error = nil
	var resp *corehttp.Response = nil
	err = retry.Do(
		func() error {
			resp, err = http.Post("https://rest.nexmo.com/sms/json", "application/x-www-form-urlencoded", params)
			if err == nil {
				defer resp.Body.Close()
			}
			return err
		},
		retry.Attempts(5),
		retry.Delay(5),
		retry.Units(time.Second),
		retry.OnRetry(func(n uint, err error) {
			glog.Errorf("attempt #%d: %s\n", n, err)
		}),
	)
	if err != nil {
		return err
	}
	if resp.StatusCode != 200 {
		return fmt.Errorf("error while sending sms through nexmo, http code: %d", resp.StatusCode)
	}
	return nil
}
