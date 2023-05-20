package i18n

var DefaultLang = "en"

var translations map[string]map[string]string = map[string]map[string]string{} //map[lang]map[en]otherlang

func RegisterTranslation(lang string, translation map[string]string) {
	translations[lang] = translation
}

func RegisterTranslate(lang string, text string, translate string) {
	if translations[lang] == nil {
		translations[lang] = map[string]string{}
	}

	translations[lang][text] = translate
}

func Translate(lang string, text string) string {
	return translations[lang][text]
}
