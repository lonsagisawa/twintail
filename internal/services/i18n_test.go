package services

import (
	"testing"
	"testing/fstest"
)

func TestNewI18n_Success(t *testing.T) {
	fs := fstest.MapFS{
		"en.json": &fstest.MapFile{Data: []byte(`{"hello": "Hello", "world": "World"}`)},
		"ja.json": &fstest.MapFile{Data: []byte(`{"hello": "こんにちは", "world": "世界"}`)},
	}

	i18n, err := NewI18n(fs, "en")

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if i18n == nil {
		t.Fatal("expected i18n instance, got nil")
	}
	if len(i18n.translations) != 2 {
		t.Errorf("expected 2 languages, got %d", len(i18n.translations))
	}
}

func TestNewI18n_InvalidJSON(t *testing.T) {
	fs := fstest.MapFS{
		"en.json": &fstest.MapFile{Data: []byte(`invalid json`)},
	}

	_, err := NewI18n(fs, "en")

	if err == nil {
		t.Fatal("expected error for invalid JSON, got nil")
	}
}

func TestNewI18n_SkipsNonJSONFiles(t *testing.T) {
	fs := fstest.MapFS{
		"en.json":  &fstest.MapFile{Data: []byte(`{"hello": "Hello"}`)},
		"readme":   &fstest.MapFile{Data: []byte(`not json`)},
		"test.txt": &fstest.MapFile{Data: []byte(`not json`)},
	}

	i18n, err := NewI18n(fs, "en")

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(i18n.translations) != 1 {
		t.Errorf("expected 1 language, got %d", len(i18n.translations))
	}
}

func TestT_ReturnsTranslation(t *testing.T) {
	fs := fstest.MapFS{
		"en.json": &fstest.MapFile{Data: []byte(`{"hello": "Hello"}`)},
		"ja.json": &fstest.MapFile{Data: []byte(`{"hello": "こんにちは"}`)},
	}

	i18n, _ := NewI18n(fs, "en")

	if got := i18n.T("en", "hello"); got != "Hello" {
		t.Errorf("expected 'Hello', got '%s'", got)
	}
	if got := i18n.T("ja", "hello"); got != "こんにちは" {
		t.Errorf("expected 'こんにちは', got '%s'", got)
	}
}

func TestT_FallbackToDefaultLang(t *testing.T) {
	fs := fstest.MapFS{
		"en.json": &fstest.MapFile{Data: []byte(`{"hello": "Hello", "world": "World"}`)},
		"ja.json": &fstest.MapFile{Data: []byte(`{"hello": "こんにちは"}`)},
	}

	i18n, _ := NewI18n(fs, "en")

	// "world" key exists only in en, should fallback
	if got := i18n.T("ja", "world"); got != "World" {
		t.Errorf("expected 'World' (fallback), got '%s'", got)
	}
}

func TestT_ReturnsKeyIfNotFound(t *testing.T) {
	fs := fstest.MapFS{
		"en.json": &fstest.MapFile{Data: []byte(`{"hello": "Hello"}`)},
	}

	i18n, _ := NewI18n(fs, "en")

	if got := i18n.T("en", "nonexistent"); got != "nonexistent" {
		t.Errorf("expected 'nonexistent', got '%s'", got)
	}
}

func TestT_UnknownLanguageFallsBackToDefault(t *testing.T) {
	fs := fstest.MapFS{
		"en.json": &fstest.MapFile{Data: []byte(`{"hello": "Hello"}`)},
	}

	i18n, _ := NewI18n(fs, "en")

	if got := i18n.T("fr", "hello"); got != "Hello" {
		t.Errorf("expected 'Hello' (fallback), got '%s'", got)
	}
}

func TestGetTranslator(t *testing.T) {
	fs := fstest.MapFS{
		"en.json": &fstest.MapFile{Data: []byte(`{"hello": "Hello"}`)},
		"ja.json": &fstest.MapFile{Data: []byte(`{"hello": "こんにちは"}`)},
	}

	i18n, _ := NewI18n(fs, "en")

	translator := i18n.GetTranslator("ja")
	if got := translator("hello"); got != "こんにちは" {
		t.Errorf("expected 'こんにちは', got '%s'", got)
	}
}

func TestParseAcceptLanguage(t *testing.T) {
	tests := []struct {
		header   string
		expected string
	}{
		{"", "en"},
		{"en", "en"},
		{"en-US", "en"},
		{"ja", "ja"},
		{"ja-JP", "ja"},
		{"ja,en;q=0.9", "ja"},
		{"en,ja;q=0.9", "en"},
		{"fr,de", "en"}, // unsupported languages fallback to en
		{"ja-JP,ja;q=0.9,en;q=0.8", "ja"},
	}

	for _, tt := range tests {
		t.Run(tt.header, func(t *testing.T) {
			got := ParseAcceptLanguage(tt.header)
			if got != tt.expected {
				t.Errorf("ParseAcceptLanguage(%q) = %q, want %q", tt.header, got, tt.expected)
			}
		})
	}
}

func TestNormalizeLang(t *testing.T) {
	tests := []struct {
		lang     string
		expected string
	}{
		{"en", "en"},
		{"ja", "ja"},
		{"fr", "en"},
		{"", "en"},
		{"zh", "en"},
	}

	for _, tt := range tests {
		t.Run(tt.lang, func(t *testing.T) {
			got := NormalizeLang(tt.lang)
			if got != tt.expected {
				t.Errorf("NormalizeLang(%q) = %q, want %q", tt.lang, got, tt.expected)
			}
		})
	}
}

func TestGetSupportedLanguages(t *testing.T) {
	langs := GetSupportedLanguages()

	if len(langs) != 2 {
		t.Errorf("expected 2 languages, got %d", len(langs))
	}
	if langs[0] != "en" || langs[1] != "ja" {
		t.Errorf("expected [en, ja], got %v", langs)
	}
}
