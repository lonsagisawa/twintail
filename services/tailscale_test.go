package services

import (
	"errors"
	"strings"
	"testing"
)

type mockCmd struct {
	output []byte
	err    error
}

func (m *mockCmd) Output() ([]byte, error) {
	return m.output, m.err
}

func (m *mockCmd) CombinedOutput() ([]byte, error) {
	return m.output, m.err
}

var mockServeOutput []byte
var mockAdvertiseError error

func TestGetServeStatus_Success(t *testing.T) {
	jsonData := `{
		"Services": {
			"svc:web-app": {
				"Web": {
					"example.com:443": {
						"Handlers": {
							"/": {"Proxy": "http://localhost:3000"}
						}
					}
				}
			}
		}
	}`
	mockServeOutput = []byte(jsonData)
	mockAdvertiseError = nil
	defer func() {
		mockServeOutput = nil
		mockAdvertiseError = nil
	}()

	oldExecCommand := execCommand
	defer func() { execCommand = oldExecCommand }()
	execCommand = func(name string, args ...string) interface {
		Output() ([]byte, error)
		CombinedOutput() ([]byte, error)
	} {
		switch {
		case strings.Contains(strings.Join(args, " "), "serve status"):
			if mockServeOutput != nil {
				return &mockCmd{output: mockServeOutput}
			}
			return &mockCmd{output: []byte(`{"Services": {}}`)}
		case strings.Contains(strings.Join(args, " "), "serve --service="):
			if mockAdvertiseError != nil {
				return &mockCmd{err: mockAdvertiseError, output: []byte("command failed")}
			}
			return &mockCmd{output: []byte("success")}
		}
		return &mockCmd{err: errors.New("unexpected command")}
	}

	svc := NewTailscaleService()
	services, err := svc.GetServeStatus()

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(services) != 1 {
		t.Fatalf("expected 1 service, got %d", len(services))
	}
	if services[0].Name != "web-app" {
		t.Errorf("expected name 'web-app', got '%s'", services[0].Name)
	}
	if services[0].HTTPSUrl != "https://example.com" {
		t.Errorf("expected URL 'https://example.com', got '%s'", services[0].HTTPSUrl)
	}
	if services[0].Proxy != "http://localhost:3000" {
		t.Errorf("expected proxy 'http://localhost:3000', got '%s'", services[0].Proxy)
	}
}

func TestGetServeStatus_JSONParseError(t *testing.T) {
	mockServeOutput = []byte("invalid json")
	mockAdvertiseError = nil
	defer func() {
		mockServeOutput = nil
		mockAdvertiseError = nil
	}()

	oldExecCommand := execCommand
	defer func() { execCommand = oldExecCommand }()
	execCommand = func(name string, args ...string) interface {
		Output() ([]byte, error)
		CombinedOutput() ([]byte, error)
	} {
		switch {
		case strings.Contains(strings.Join(args, " "), "serve status"):
			if mockServeOutput != nil {
				return &mockCmd{output: mockServeOutput}
			}
			return &mockCmd{output: []byte(`{"Services": {}}`)}
		case strings.Contains(strings.Join(args, " "), "serve --service="):
			if mockAdvertiseError != nil {
				return &mockCmd{err: mockAdvertiseError, output: []byte("command failed")}
			}
			return &mockCmd{output: []byte("success")}
		}
		return &mockCmd{err: errors.New("unexpected command")}
	}

	svc := NewTailscaleService()
	_, err := svc.GetServeStatus()

	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestAdvertiseService_Success(t *testing.T) {
	mockServeOutput = nil
	mockAdvertiseError = nil
	defer func() {
		mockServeOutput = nil
		mockAdvertiseError = nil
	}()

	oldExecCommand := execCommand
	defer func() { execCommand = oldExecCommand }()
	execCommand = func(name string, args ...string) interface {
		Output() ([]byte, error)
		CombinedOutput() ([]byte, error)
	} {
		switch {
		case strings.Contains(strings.Join(args, " "), "serve status"):
			if mockServeOutput != nil {
				return &mockCmd{output: mockServeOutput}
			}
			return &mockCmd{output: []byte(`{"Services": {}}`)}
		case strings.Contains(strings.Join(args, " "), "serve --service="):
			if mockAdvertiseError != nil {
				return &mockCmd{err: mockAdvertiseError, output: []byte("command failed")}
			}
			return &mockCmd{output: []byte("success")}
		}
		return &mockCmd{err: errors.New("unexpected command")}
	}

	svc := NewTailscaleService()
	params := AdvertiseServiceParams{
		ServiceName: "my-service",
		Protocol:    "https",
		ExposePort:  "443",
		Destination: "http://localhost:8080",
	}
	err := svc.AdvertiseService(params)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestAdvertiseService_Failure(t *testing.T) {
	mockServeOutput = nil
	mockAdvertiseError = errors.New("command failed")
	defer func() {
		mockServeOutput = nil
		mockAdvertiseError = nil
	}()

	oldExecCommand := execCommand
	defer func() { execCommand = oldExecCommand }()
	execCommand = func(name string, args ...string) interface {
		Output() ([]byte, error)
		CombinedOutput() ([]byte, error)
	} {
		switch {
		case strings.Contains(strings.Join(args, " "), "serve status"):
			if mockServeOutput != nil {
				return &mockCmd{output: mockServeOutput}
			}
			return &mockCmd{output: []byte(`{"Services": {}}`)}
		case strings.Contains(strings.Join(args, " "), "serve --service="):
			if mockAdvertiseError != nil {
				return &mockCmd{err: mockAdvertiseError, output: []byte("command failed")}
			}
			return &mockCmd{output: []byte("success")}
		}
		return &mockCmd{err: errors.New("unexpected command")}
	}

	svc := NewTailscaleService()
	params := AdvertiseServiceParams{
		ServiceName: "my-service",
		Protocol:    "https",
		ExposePort:  "443",
		Destination: "http://localhost:8080",
	}
	err := svc.AdvertiseService(params)

	if err == nil {
		t.Fatal("expected error, got nil")
	}

	advErr, ok := err.(*AdvertiseError)
	if !ok {
		t.Fatalf("expected AdvertiseError, got %T", err)
	}
	if advErr.Message != "command failed" {
		t.Errorf("expected message 'command failed', got '%s'", advErr.Message)
	}
}

func TestGetServeStatus_EmptyServices(t *testing.T) {
	mockServeOutput = []byte(`{"Services": {}}`)
	mockAdvertiseError = nil
	defer func() {
		mockServeOutput = nil
		mockAdvertiseError = nil
	}()

	oldExecCommand := execCommand
	defer func() { execCommand = oldExecCommand }()
	execCommand = func(name string, args ...string) interface {
		Output() ([]byte, error)
		CombinedOutput() ([]byte, error)
	} {
		switch {
		case strings.Contains(strings.Join(args, " "), "serve status"):
			if mockServeOutput != nil {
				return &mockCmd{output: mockServeOutput}
			}
			return &mockCmd{output: []byte(`{"Services": {}}`)}
		case strings.Contains(strings.Join(args, " "), "serve --service="):
			if mockAdvertiseError != nil {
				return &mockCmd{err: mockAdvertiseError, output: []byte("command failed")}
			}
			return &mockCmd{output: []byte("success")}
		}
		return &mockCmd{err: errors.New("unexpected command")}
	}

	svc := NewTailscaleService()
	services, err := svc.GetServeStatus()

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(services) != 0 {
		t.Fatalf("expected 0 services, got %d", len(services))
	}
}

func setupMockExecCommand() func() {
	oldExecCommand := execCommand
	execCommand = func(name string, args ...string) interface {
		Output() ([]byte, error)
		CombinedOutput() ([]byte, error)
	} {
		switch {
		case strings.Contains(strings.Join(args, " "), "serve status"):
			if mockServeOutput != nil {
				return &mockCmd{output: mockServeOutput}
			}
			return &mockCmd{output: []byte(`{"Services": {}}`)}
		case strings.Contains(strings.Join(args, " "), "serve --service="):
			if mockAdvertiseError != nil {
				return &mockCmd{err: mockAdvertiseError, output: []byte("command failed")}
			}
			return &mockCmd{output: []byte("success")}
		}
		return &mockCmd{err: errors.New("unexpected command")}
	}
	return func() { execCommand = oldExecCommand }
}

func TestGetServiceByName_Success(t *testing.T) {
	jsonData := `{
		"Services": {
			"svc:web-app": {
				"Web": {
					"example.com:443": {
						"Handlers": {
							"/": {"Proxy": "http://localhost:3000"}
						}
					}
				}
			}
		}
	}`
	mockServeOutput = []byte(jsonData)
	mockAdvertiseError = nil
	defer func() {
		mockServeOutput = nil
		mockAdvertiseError = nil
	}()
	defer setupMockExecCommand()()

	svc := NewTailscaleService()
	detail, err := svc.GetServiceByName("web-app")

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if detail == nil {
		t.Fatal("expected detail, got nil")
	}
	if detail.Name != "web-app" {
		t.Errorf("expected name 'web-app', got '%s'", detail.Name)
	}
	if detail.Hostname != "example.com" {
		t.Errorf("expected hostname 'example.com', got '%s'", detail.Hostname)
	}
	if detail.URL != "https://example.com" {
		t.Errorf("expected URL 'https://example.com', got '%s'", detail.URL)
	}
	if len(detail.Ports) != 1 {
		t.Fatalf("expected 1 port, got %d", len(detail.Ports))
	}
	if detail.Ports[0].Protocol != "https" {
		t.Errorf("expected protocol 'https', got '%s'", detail.Ports[0].Protocol)
	}
	if detail.Ports[0].ExposePort != "443" {
		t.Errorf("expected port '443', got '%s'", detail.Ports[0].ExposePort)
	}
	if detail.Ports[0].Destination != "http://localhost:3000" {
		t.Errorf("expected destination 'http://localhost:3000', got '%s'", detail.Ports[0].Destination)
	}
}

func TestGetServiceByName_NotFound(t *testing.T) {
	mockServeOutput = []byte(`{"Services": {}}`)
	mockAdvertiseError = nil
	defer func() {
		mockServeOutput = nil
		mockAdvertiseError = nil
	}()
	defer setupMockExecCommand()()

	svc := NewTailscaleService()
	detail, err := svc.GetServiceByName("nonexistent")

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if detail != nil {
		t.Errorf("expected nil, got %+v", detail)
	}
}

func TestGetServiceByName_MultiplePorts(t *testing.T) {
	jsonData := `{
		"Services": {
			"svc:multi-port": {
				"Web": {
					"example.com:443": {
						"Handlers": {
							"/": {"Proxy": "http://localhost:3000"}
						}
					},
					"example.com:8080": {
						"Handlers": {
							"/": {"Proxy": "http://localhost:8080"}
						}
					}
				}
			}
		}
	}`
	mockServeOutput = []byte(jsonData)
	mockAdvertiseError = nil
	defer func() {
		mockServeOutput = nil
		mockAdvertiseError = nil
	}()
	defer setupMockExecCommand()()

	svc := NewTailscaleService()
	detail, err := svc.GetServiceByName("multi-port")

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if detail == nil {
		t.Fatal("expected detail, got nil")
	}
	if len(detail.Ports) != 2 {
		t.Fatalf("expected 2 ports, got %d", len(detail.Ports))
	}
	if detail.URL != "https://example.com" {
		t.Errorf("expected URL 'https://example.com', got '%s'", detail.URL)
	}
}

func TestGetServiceByName_HTTPOnly(t *testing.T) {
	jsonData := `{
		"Services": {
			"svc:http-only": {
				"Web": {
					"example.com:80": {
						"Handlers": {
							"/": {"Proxy": "http://localhost:3000"}
						}
					}
				}
			}
		}
	}`
	mockServeOutput = []byte(jsonData)
	mockAdvertiseError = nil
	defer func() {
		mockServeOutput = nil
		mockAdvertiseError = nil
	}()
	defer setupMockExecCommand()()

	svc := NewTailscaleService()
	detail, err := svc.GetServiceByName("http-only")

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if detail == nil {
		t.Fatal("expected detail, got nil")
	}
	if detail.URL != "http://example.com" {
		t.Errorf("expected URL 'http://example.com', got '%s'", detail.URL)
	}
	if detail.Ports[0].Protocol != "http" {
		t.Errorf("expected protocol 'http', got '%s'", detail.Ports[0].Protocol)
	}
}

func TestGetServiceByName_HTTPNonStandardPort(t *testing.T) {
	jsonData := `{
		"Services": {
			"svc:custom-port": {
				"Web": {
					"example.com:8080": {
						"Handlers": {
							"/": {"Proxy": "http://localhost:3000"}
						}
					}
				}
			}
		}
	}`
	mockServeOutput = []byte(jsonData)
	mockAdvertiseError = nil
	defer func() {
		mockServeOutput = nil
		mockAdvertiseError = nil
	}()
	defer setupMockExecCommand()()

	svc := NewTailscaleService()
	detail, err := svc.GetServiceByName("custom-port")

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if detail == nil {
		t.Fatal("expected detail, got nil")
	}
	if detail.URL != "http://example.com:8080" {
		t.Errorf("expected URL 'http://example.com:8080', got '%s'", detail.URL)
	}
}
