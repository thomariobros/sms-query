package i18n

import (
	"github.com/goodsign/monday"
	"github.com/nicksnyder/go-i18n/i18n"

	"sms-query/pkg/config"
)

func LoadWithRootPath(configRootPath string) {
	i18n.MustLoadTranslationFile(configRootPath + "/i18n/en_US.json")
	i18n.MustLoadTranslationFile(configRootPath + "/i18n/fr_FR.json")
}

func LoadWithContent(content map[string][]byte) {
	i18n.ParseTranslationFileBytes("config/i18n/en_US.json", content["i18n-en_US"])
	i18n.ParseTranslationFileBytes("config/i18n/fr_FR.json", content["i18n-fr_FR"])
}

func GetTranslationForPhoneNumber(phoneNumber string, key string) string {
	locale := ""
	defaults := config.GetInstance().GetDefaults(phoneNumber)
	if defaults != nil {
		locale = defaults.Locale
	}
	return GetTranslation(locale, key)
}

func GetTranslation(locale string, key string) string {
	T, _ := i18n.Tfunc(locale, "en_US")
	return T(key)
}

func GetMondayLocale(locale string) monday.Locale {
	switch locale {
	case "fr_FR":
		return monday.LocaleFrFR
	default:
		return monday.LocaleEnUS
	}
}
