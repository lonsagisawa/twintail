package handlers

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v5"
)

func TestSettingsHandler_Show(t *testing.T) {
	handler := NewSettingsHandler()

	e := echo.New()
	e.Renderer = &mockRenderer{}
	req := httptest.NewRequest(http.MethodGet, "/settings", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Set("lang", "en")

	err := handler.Show(c)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if rec.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", rec.Code)
	}
}

func TestSettingsHandler_Update_Success(t *testing.T) {
	handler := NewSettingsHandler()

	e := echo.New()
	e.Renderer = &mockRenderer{}
	e.Validator = newTestValidator()
	form := strings.NewReader("lang=ja")
	req := httptest.NewRequest(http.MethodPost, "/settings", form)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := handler.Update(c)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if rec.Code != http.StatusSeeOther {
		t.Errorf("expected status 303, got %d", rec.Code)
	}

	cookies := rec.Result().Cookies()
	var langCookie *http.Cookie
	for _, cookie := range cookies {
		if cookie.Name == "lang" {
			langCookie = cookie
			break
		}
	}

	if langCookie == nil {
		t.Fatal("expected lang cookie to be set")
	}
	if langCookie.Value != "ja" {
		t.Errorf("expected cookie value 'ja', got '%s'", langCookie.Value)
	}
	if langCookie.Path != "/" {
		t.Errorf("expected cookie path '/', got '%s'", langCookie.Path)
	}
	if !langCookie.HttpOnly {
		t.Error("expected cookie to be HttpOnly")
	}
}

func TestSettingsHandler_Update_EnglishLanguage(t *testing.T) {
	handler := NewSettingsHandler()

	e := echo.New()
	e.Renderer = &mockRenderer{}
	e.Validator = newTestValidator()
	form := strings.NewReader("lang=en")
	req := httptest.NewRequest(http.MethodPost, "/settings", form)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := handler.Update(c)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	cookies := rec.Result().Cookies()
	var langCookie *http.Cookie
	for _, cookie := range cookies {
		if cookie.Name == "lang" {
			langCookie = cookie
			break
		}
	}

	if langCookie == nil {
		t.Fatal("expected lang cookie to be set")
	}
	if langCookie.Value != "en" {
		t.Errorf("expected cookie value 'en', got '%s'", langCookie.Value)
	}
}

func TestSettingsHandler_Update_InvalidLanguage(t *testing.T) {
	handler := NewSettingsHandler()

	e := echo.New()
	e.Renderer = &mockRenderer{}
	e.Validator = newTestValidator()
	form := strings.NewReader("lang=fr")
	req := httptest.NewRequest(http.MethodPost, "/settings", form)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := handler.Update(c)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if rec.Code != http.StatusSeeOther {
		t.Errorf("expected status 303, got %d", rec.Code)
	}

	cookies := rec.Result().Cookies()
	for _, cookie := range cookies {
		if cookie.Name == "lang" {
			t.Error("expected no lang cookie for invalid language")
		}
	}
}

func TestSettingsHandler_Update_EmptyLanguage(t *testing.T) {
	handler := NewSettingsHandler()

	e := echo.New()
	e.Renderer = &mockRenderer{}
	e.Validator = newTestValidator()
	form := strings.NewReader("lang=")
	req := httptest.NewRequest(http.MethodPost, "/settings", form)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := handler.Update(c)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if rec.Code != http.StatusSeeOther {
		t.Errorf("expected status 303, got %d", rec.Code)
	}

	cookies := rec.Result().Cookies()
	for _, cookie := range cookies {
		if cookie.Name == "lang" {
			t.Error("expected no lang cookie for empty language")
		}
	}
}

func TestSettingsHandler_Update_MissingField(t *testing.T) {
	handler := NewSettingsHandler()

	e := echo.New()
	e.Renderer = &mockRenderer{}
	e.Validator = newTestValidator()
	form := strings.NewReader("")
	req := httptest.NewRequest(http.MethodPost, "/settings", form)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := handler.Update(c)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if rec.Code != http.StatusSeeOther {
		t.Errorf("expected status 303, got %d", rec.Code)
	}
}
