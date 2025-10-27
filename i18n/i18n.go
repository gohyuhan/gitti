package i18n

import (
	"gitti/utils"
	"strings"
)

var SUPPORTED_LANGUAGE_CODE = []string{"EN", "JP", "ZH-TW", "ZH-CN"}

var LANGUAGEMAPPING *LanguageMapping

func InitGittiLanguageMapping(languageCode string) {
	languageCode = strings.ToUpper(languageCode)
	switch languageCode {
	case "EN":
		LANGUAGEMAPPING = &eN
	case "JP":
		LANGUAGEMAPPING = &jP
	case "ZH-TW":
		LANGUAGEMAPPING = &zH_TW
	case "ZH-CN":
		LANGUAGEMAPPING = &zH_CN
	default:
		LANGUAGEMAPPING = &eN
	}
}

func IsLanguageCodeSupported(languageCode string) bool {
	if utils.Contains(SUPPORTED_LANGUAGE_CODE, strings.ToUpper(languageCode)) {
		return true
	}
	return false
}
