package controllers

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"twintail/services"

	"github.com/labstack/echo/v5"
)

type mockTailscaleService struct {
	services      []services.ServiceView
	serviceDetail *services.ServiceDetailView
	advertiseErr  error
}

func (m *mockTailscaleService) GetServeStatus() ([]services.ServiceView, error) {
	return m.services, m.advertiseErr
}

func (m *mockTailscaleService) GetServiceByName(name string) (*services.ServiceDetailView, error) {
	return m.serviceDetail, m.advertiseErr
}

func (m *mockTailscaleService) AdvertiseService(params services.AdvertiseServiceParams) error {
	return m.advertiseErr
}

type mockRenderer struct{}

func (m *mockRenderer) Render(ctx *echo.Context, w io.Writer, name string, data any) error {
	return nil
}

func TestIndex_Success(t *testing.T) {
	mockSvc := &mockTailscaleService{
		services: []services.ServiceView{
			{Name: "web-app", HTTPSUrl: "https://example.com", Proxy: "http://localhost:3000"},
		},
	}
	ctrl := NewServiceController(mockSvc)

	e := echo.New()
	e.Renderer = &mockRenderer{}
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := ctrl.Index(c)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if rec.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", rec.Code)
	}
}

func TestIndex_GetServeStatusError(t *testing.T) {
	mockSvc := &mockTailscaleService{
		advertiseErr: &services.AdvertiseError{Message: "Failed to get serve status", Err: nil},
	}
	ctrl := NewServiceController(mockSvc)

	e := echo.New()
	e.Renderer = &mockRenderer{}
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := ctrl.Index(c)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if rec.Code != http.StatusInternalServerError {
		t.Errorf("expected status 500, got %d", rec.Code)
	}
	body := rec.Body.String()
	if !strings.Contains(body, "Failed to get serve status") {
		t.Errorf("expected error message in response, got '%s'", body)
	}
}

func TestCreate(t *testing.T) {
	mockSvc := &mockTailscaleService{}
	ctrl := NewServiceController(mockSvc)

	e := echo.New()
	e.Renderer = &mockRenderer{}
	req := httptest.NewRequest(http.MethodGet, "/services/new", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := ctrl.Create(c)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if rec.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", rec.Code)
	}
}

func TestStore_Success(t *testing.T) {
	mockSvc := &mockTailscaleService{
		advertiseErr: nil,
	}
	ctrl := NewServiceController(mockSvc)

	e := echo.New()
	e.Renderer = &mockRenderer{}
	form := strings.NewReader("service_name=my-service&protocol=https&expose_port=443&destination=http://localhost:8080")
	req := httptest.NewRequest(http.MethodPost, "/services/new", form)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := ctrl.Store(c)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if rec.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", rec.Code)
	}
}

func TestStore_Failure(t *testing.T) {
	mockSvc := &mockTailscaleService{
		advertiseErr: &services.AdvertiseError{Message: "Service already exists", Err: nil},
	}
	ctrl := NewServiceController(mockSvc)

	e := echo.New()
	e.Renderer = &mockRenderer{}
	form := strings.NewReader("service_name=my-service&protocol=https&expose_port=443&destination=http://localhost:8080")
	req := httptest.NewRequest(http.MethodPost, "/services/new", form)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := ctrl.Store(c)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if rec.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", rec.Code)
	}
}

func TestShow_Success(t *testing.T) {
	mockSvc := &mockTailscaleService{
		serviceDetail: &services.ServiceDetailView{
			Name:     "web-app",
			Hostname: "example.com",
			URL:      "https://example.com",
			Ports: []services.PortEntry{
				{Protocol: "https", ExposePort: "443", Destination: "http://localhost:3000"},
			},
		},
	}
	ctrl := NewServiceController(mockSvc)

	e := echo.New()
	e.Renderer = &mockRenderer{}
	e.GET("/services/:name", ctrl.Show)

	req := httptest.NewRequest(http.MethodGet, "/services/web-app", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", rec.Code)
	}
}

func TestShow_NotFound(t *testing.T) {
	mockSvc := &mockTailscaleService{
		serviceDetail: nil,
	}
	ctrl := NewServiceController(mockSvc)

	e := echo.New()
	e.Renderer = &mockRenderer{}
	e.GET("/services/:name", ctrl.Show)

	req := httptest.NewRequest(http.MethodGet, "/services/nonexistent", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Errorf("expected status 404, got %d", rec.Code)
	}
}
