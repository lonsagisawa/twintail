package handlers

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"twintail/services"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v5"
)

type endpointTestValidator struct {
	validator *validator.Validate
}

func (tv *endpointTestValidator) Validate(i any) error {
	return tv.validator.Struct(i)
}

func newEndpointTestValidator() *endpointTestValidator {
	return &endpointTestValidator{validator: validator.New()}
}

type mockEndpointService struct {
	serviceDetail *services.ServiceDetailView
	endpointErr   error
}

func (m *mockEndpointService) GetServiceByName(name string) (*services.ServiceDetailView, error) {
	return m.serviceDetail, nil
}

func (m *mockEndpointService) AddEndpoint(params services.EndpointParams) error {
	return m.endpointErr
}

func (m *mockEndpointService) RemoveEndpoint(params services.EndpointParams) error {
	return m.endpointErr
}

func (m *mockEndpointService) UpdateEndpoint(params services.UpdateEndpointParams) error {
	return m.endpointErr
}

func TestEndpointCreate(t *testing.T) {
	mockSvc := &mockEndpointService{}
	ctrl := NewEndpointHandler(mockSvc)

	e := echo.New()
	e.Renderer = &mockRenderer{}
	e.GET("/services/:name/endpoints/new", ctrl.Create)

	req := httptest.NewRequest(http.MethodGet, "/services/my-service/endpoints/new", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", rec.Code)
	}
}

func TestEndpointStore_Success(t *testing.T) {
	mockSvc := &mockEndpointService{
		endpointErr: nil,
	}
	ctrl := NewEndpointHandler(mockSvc)

	e := echo.New()
	e.Renderer = &mockRenderer{}
	e.Validator = newEndpointTestValidator()
	e.POST("/services/:name/endpoints/new", ctrl.Store)

	form := strings.NewReader("protocol=https&expose_port=443&destination=http://localhost:8080")
	req := httptest.NewRequest(http.MethodPost, "/services/my-service/endpoints/new", form)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code != http.StatusSeeOther {
		t.Errorf("expected status 303, got %d", rec.Code)
	}
	if rec.Header().Get("Location") != "/services/my-service" {
		t.Errorf("expected redirect to /services/my-service, got %s", rec.Header().Get("Location"))
	}
}

func TestEndpointStore_Failure(t *testing.T) {
	mockSvc := &mockEndpointService{
		endpointErr: &services.CommandError{Message: "Failed to add endpoint", Err: nil},
	}
	ctrl := NewEndpointHandler(mockSvc)

	e := echo.New()
	e.Renderer = &mockRenderer{}
	e.POST("/services/:name/endpoints/new", ctrl.Store)

	form := strings.NewReader("protocol=https&expose_port=443&destination=http://localhost:8080")
	req := httptest.NewRequest(http.MethodPost, "/services/my-service/endpoints/new", form)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", rec.Code)
	}
}

func TestEndpointDelete(t *testing.T) {
	mockSvc := &mockEndpointService{}
	ctrl := NewEndpointHandler(mockSvc)

	e := echo.New()
	e.Renderer = &mockRenderer{}
	e.GET("/services/:name/endpoints/delete", ctrl.Delete)

	req := httptest.NewRequest(http.MethodGet, "/services/my-service/endpoints/delete?protocol=https&port=443&destination=http://localhost:8080", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", rec.Code)
	}
}

func TestEndpointDestroy_Success_ServiceExists(t *testing.T) {
	mockSvc := &mockEndpointService{
		endpointErr: nil,
		serviceDetail: &services.ServiceDetailView{
			Name:     "my-service",
			Hostname: "example.com",
			Ports: []services.PortEntry{
				{Protocol: "http", ExposePort: "80", Destination: "http://localhost:8080"},
			},
		},
	}
	ctrl := NewEndpointHandler(mockSvc)

	e := echo.New()
	e.Renderer = &mockRenderer{}
	e.Validator = newEndpointTestValidator()
	e.POST("/services/:name/endpoints/delete", ctrl.Destroy)

	form := strings.NewReader("protocol=https&expose_port=443&destination=http://localhost:8080")
	req := httptest.NewRequest(http.MethodPost, "/services/my-service/endpoints/delete", form)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code != http.StatusSeeOther {
		t.Errorf("expected status 303, got %d", rec.Code)
	}
	if rec.Header().Get("Location") != "/services/my-service" {
		t.Errorf("expected redirect to /services/my-service, got %s", rec.Header().Get("Location"))
	}
}

func TestEndpointDestroy_Success_ServiceGone(t *testing.T) {
	mockSvc := &mockEndpointService{
		endpointErr:   nil,
		serviceDetail: nil,
	}
	ctrl := NewEndpointHandler(mockSvc)

	e := echo.New()
	e.Renderer = &mockRenderer{}
	e.Validator = newEndpointTestValidator()
	e.POST("/services/:name/endpoints/delete", ctrl.Destroy)

	form := strings.NewReader("protocol=https&expose_port=443&destination=http://localhost:8080")
	req := httptest.NewRequest(http.MethodPost, "/services/my-service/endpoints/delete", form)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code != http.StatusSeeOther {
		t.Errorf("expected status 303, got %d", rec.Code)
	}
	if rec.Header().Get("Location") != "/" {
		t.Errorf("expected redirect to /, got %s", rec.Header().Get("Location"))
	}
}

func TestEndpointDestroy_Failure(t *testing.T) {
	mockSvc := &mockEndpointService{
		endpointErr: &services.CommandError{Message: "Failed to remove endpoint", Err: nil},
	}
	ctrl := NewEndpointHandler(mockSvc)

	e := echo.New()
	e.Renderer = &mockRenderer{}
	e.POST("/services/:name/endpoints/delete", ctrl.Destroy)

	form := strings.NewReader("protocol=https&expose_port=443&destination=http://localhost:8080")
	req := httptest.NewRequest(http.MethodPost, "/services/my-service/endpoints/delete", form)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code != http.StatusInternalServerError {
		t.Errorf("expected status 500, got %d", rec.Code)
	}
}
