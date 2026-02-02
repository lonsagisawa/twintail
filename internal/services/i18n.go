package services

import (
	"encoding/json"
	"io/fs"
	"strings"
	"sync"
)

type I18n struct {
	translations map[string]map[string]string
	defaultLang  string
	mu           sync.RWMutex
}

func NewI18n(localesFS fs.FS, defaultLang string) (*I18n, error) {
	i18n := &I18n{
		translations: make(map[string]map[string]string),
		defaultLang:  defaultLang,
	}

	entries, err := fs.ReadDir(localesFS, ".")
	if err != nil {
		return nil, err
	}

	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".json") {
			continue
		}

		lang := strings.TrimSuffix(entry.Name(), ".json")
		data, err := fs.ReadFile(localesFS, entry.Name())
		if err != nil {
			return nil, err
		}

		var messages map[string]string
		if err := json.Unmarshal(data, &messages); err != nil {
			return nil, err
		}

		i18n.translations[lang] = messages
	}

	return i18n, nil
}

func (i *I18n) T(lang, key string) string {
	i.mu.RLock()
	defer i.mu.RUnlock()

	if messages, ok := i.translations[lang]; ok {
		if msg, ok := messages[key]; ok {
			return msg
		}
	}

	if messages, ok := i.translations[i.defaultLang]; ok {
		if msg, ok := messages[key]; ok {
			return msg
		}
	}

	return key
}

func (i *I18n) GetTranslator(lang string) func(key string) string {
	return func(key string) string {
		return i.T(lang, key)
	}
}

func ParseAcceptLanguage(header string) string {
	if header == "" {
		return "en"
	}

	parts := strings.Split(header, ",")
	for _, part := range parts {
		lang := strings.TrimSpace(strings.Split(part, ";")[0])

		if strings.HasPrefix(lang, "ja") {
			return "ja"
		}
		if strings.HasPrefix(lang, "en") {
			return "en"
		}
	}

	return "en"
}

func NormalizeLang(lang string) string {
	switch lang {
	case "ja", "en":
		return lang
	default:
		return "en"
	}
}

func GetSupportedLanguages() []string {
	return []string{"en", "ja"}
}
