package i18n

import (
	"strings"

	"github.com/gohyuhan/gitti/utils"
)

// we are using ISO 639-1 as our default for the language,
// but when it comes to language like chinese which have simplified and traditional we plus a script tag (ISO 15924) along with ISO 639-1
// eg, english will be en ( doesn't matter if it was uk or us ) but chinese will be categorize to ZH-HANS and ZH-HANT (simplified and traditional)
var SUPPORTED_LANGUAGE_CODE = []string{"EN", "JA", "ZH-HANT", "ZH-HANS"}

var LANGUAGEMAPPING *LanguageMapping

func InitGittiLanguageMapping(languageCode string) {
	languageCode = strings.ToUpper(languageCode)
	switch languageCode {
	case "EN":
		LANGUAGEMAPPING = &eN
	case "JA":
		LANGUAGEMAPPING = &jA
	case "ZH-HANT":
		LANGUAGEMAPPING = &zH_HANT
	case "ZH-HANS":
		LANGUAGEMAPPING = &zH_HANS
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
