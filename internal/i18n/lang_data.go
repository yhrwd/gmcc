package i18n

type langEntry struct {
	Key   string
	Value string
}

func (i *I18n) loadLanguage(lang string) {
	data := getLangData(lang)
	for _, entry := range data {
		i.messages[entry.Key] = entry.Value
	}
}
