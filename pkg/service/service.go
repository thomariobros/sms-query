package service

import (
	"sms-query/pkg/service/maps/nominatim"
	"sms-query/pkg/service/search/qwant"
	"sms-query/pkg/service/translate/deepl"
)

func NewMapsService() *nominatim.NominatimService {
	return nominatim.NewNominatimService()
}

func NewSearchService() *qwant.QwantSearchService {
	return qwant.NewQwantSearchService()
}

func NewTranslateService() *deepl.DeepLTranslateService {
	return deepl.NewDeepLTranslateService()
}
