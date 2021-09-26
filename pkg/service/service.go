package service

import (
	"main/pkg/service/maps/nominatim"
	"main/pkg/service/search/qwant"
	"main/pkg/service/translate/deepl"
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
