package server

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"twintail/internal/handlers"
	"twintail/internal/services"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v5"
)

type mockTailscaleService struct {
	services          []services.ServiceView
	serviceDetail     *services.ServiceDetailView
	advertiseErr      error
	clearErr          error
	checkInstalledErr error
}

func (m *mockTailscaleService) CheckInstalled() error {
	return m.checkInstalledErr
}

func (m *mockTailscaleService) GetServeStatus() ([]services.ServiceView, error) {
	if m.advertiseErr != nil {
		return nil, m.advertiseErr
	}
	return m.services, nil
}

func (m *mockTailscaleService) GetServiceByName(name string) (*services.ServiceDetailView, error) {
	return m.serviceDetail, nil
}

func (m *mockTailscaleService) AdvertiseService(params services.AdvertiseServiceParams) error {
	return m.advertiseErr
}

func (m *mockTailscaleService) ClearService(name string) error {
	return m.clearErr
}

func (m *mockTailscaleService) AddEndpoint(params services.EndpointParams) error {
	return m.advertiseErr
}

func (m *mockTailscaleService) RemoveEndpoint(params services.EndpointParams) error {
	return m.advertiseErr
}

func (m *mockTailscaleService) UpdateEndpoint(params services.UpdateEndpointParams) error {
	return m.advertiseErr
}

func setupTestServer(tailscaleSvc *mockTailscaleService) *echo.Echo {
	e := echo.New()
	e.Use(I18nMiddleware())

	e.Renderer = &testRenderer{}
	e.Validator = &testValidator{}

	container := handlers.NewContainer(tailscaleSvc)
	RegisterRoutes(e, container)

	return e
}

type testRenderer struct{}

func (r *testRenderer) Render(ctx *echo.Context, w io.Writer, name string, data any) error {
	fmt.Fprintf(w, "template:%s", name)
	return nil
}

type testValidator struct{}

func (v *testValidator) Validate(i any) error {
	return validator.New().Struct(i)
}

func TestIntegration_IndexRoute(t *testing.T) {
	mockSvc := &mockTailscaleService{
		services: []services.ServiceView{
			{Name: "test-service", HTTPSUrl: "https://test.example.com"},
		},
	}
	e := setupTestServer(mockSvc)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", rec.Code)
	}
}

func TestIntegration_LanguageCookie(t *testing.T) {
	mockSvc := &mockTailscaleService{}
	e := setupTestServer(mockSvc)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.AddCookie(&http.Cookie{Name: "lang", Value: "ja"})
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", rec.Code)
	}
}

func TestIntegration_AcceptLanguageHeader(t *testing.T) {
	mockSvc := &mockTailscaleService{}
	e := setupTestServer(mockSvc)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Accept-Language", "ja-JP,ja;q=0.9,en;q=0.8")
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", rec.Code)
	}
}

func TestIntegration_ServiceCreateRedirect(t *testing.T) {
	mockSvc := &mockTailscaleService{}
	e := setupTestServer(mockSvc)

	form := strings.NewReader("service_name=new-service&protocol=https&expose_port=443&destination=http://localhost:8080")
	req := httptest.NewRequest(http.MethodPost, "/services/new", form)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code != http.StatusSeeOther {
		t.Errorf("expected status 303, got %d", rec.Code)
	}
	if rec.Header().Get("Location") != "/services/new-service" {
		t.Errorf("expected redirect to /services/new-service, got %s", rec.Header().Get("Location"))
	}
}

func TestIntegration_ServiceShowNotFound(t *testing.T) {
	mockSvc := &mockTailscaleService{
		serviceDetail: nil,
	}
	e := setupTestServer(mockSvc)

	req := httptest.NewRequest(http.MethodGet, "/services/nonexistent", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Errorf("expected status 404, got %d", rec.Code)
	}
}

func TestIntegration_SettingsRouteExists(t *testing.T) {
	mockSvc := &mockTailscaleService{}
	e := setupTestServer(mockSvc)

	req := httptest.NewRequest(http.MethodGet, "/settings", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", rec.Code)
	}
}
