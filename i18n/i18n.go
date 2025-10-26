package i18n

import "strings"

var LANGUAGEMAPPING *LanguageMapping

func InitGittiLanguageMapping(languageCode string) {
	languageCode = strings.ToUpper(languageCode)
	switch languageCode {
	case "EN":
		LANGUAGEMAPPING = &EN
	case "JP":
		LANGUAGEMAPPING = &JP
	case "ZH-TW":
		LANGUAGEMAPPING = &ZH_TW
	case "ZH-CN":
		LANGUAGEMAPPING = &ZH_CN
	default:
		LANGUAGEMAPPING = &EN
	}
}
