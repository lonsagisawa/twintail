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
var mockCommandError error

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
	mockCommandError = nil
	defer func() {
		mockServeOutput = nil
		mockCommandError = nil
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
			if mockCommandError != nil {
				return &mockCmd{err: mockCommandError, output: []byte("command failed")}
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
	mockCommandError = nil
	defer func() {
		mockServeOutput = nil
		mockCommandError = nil
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
			if mockCommandError != nil {
				return &mockCmd{err: mockCommandError, output: []byte("command failed")}
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
	mockCommandError = nil
	defer func() {
		mockServeOutput = nil
		mockCommandError = nil
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
			if mockCommandError != nil {
				return &mockCmd{err: mockCommandError, output: []byte("command failed")}
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
	mockCommandError = errors.New("command failed")
	defer func() {
		mockServeOutput = nil
		mockCommandError = nil
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
			if mockCommandError != nil {
				return &mockCmd{err: mockCommandError, output: []byte("command failed")}
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

	advErr, ok := err.(*CommandError)
	if !ok {
		t.Fatalf("expected CommandError, got %T", err)
	}
	if advErr.Message != "command failed" {
		t.Errorf("expected message 'command failed', got '%s'", advErr.Message)
	}
}

func TestGetServeStatus_EmptyServices(t *testing.T) {
	mockServeOutput = []byte(`{"Services": {}}`)
	mockCommandError = nil
	defer func() {
		mockServeOutput = nil
		mockCommandError = nil
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
			if mockCommandError != nil {
				return &mockCmd{err: mockCommandError, output: []byte("command failed")}
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
			if mockCommandError != nil {
				return &mockCmd{err: mockCommandError, output: []byte("command failed")}
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
	mockCommandError = nil
	defer func() {
		mockServeOutput = nil
		mockCommandError = nil
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
	mockCommandError = nil
	defer func() {
		mockServeOutput = nil
		mockCommandError = nil
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
	mockCommandError = nil
	defer func() {
		mockServeOutput = nil
		mockCommandError = nil
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
	mockCommandError = nil
	defer func() {
		mockServeOutput = nil
		mockCommandError = nil
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
	mockCommandError = nil
	defer func() {
		mockServeOutput = nil
		mockCommandError = nil
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

func setupMockExecCommandWithClear() func() {
	oldExecCommand := execCommand
	execCommand = func(name string, args ...string) interface {
		Output() ([]byte, error)
		CombinedOutput() ([]byte, error)
	} {
		argsStr := strings.Join(args, " ")
		switch {
		case strings.Contains(argsStr, "serve status"):
			if mockServeOutput != nil {
				return &mockCmd{output: mockServeOutput}
			}
			return &mockCmd{output: []byte(`{"Services": {}}`)}
		case strings.Contains(argsStr, "serve --service="):
			if mockCommandError != nil {
				return &mockCmd{err: mockCommandError, output: []byte("command failed")}
			}
			return &mockCmd{output: []byte("success")}
		case strings.Contains(argsStr, "serve clear"):
			if mockCommandError != nil {
				return &mockCmd{err: mockCommandError, output: []byte("clear failed")}
			}
			return &mockCmd{output: []byte("success")}
		}
		return &mockCmd{err: errors.New("unexpected command")}
	}
	return func() { execCommand = oldExecCommand }
}

func TestClearService_Success(t *testing.T) {
	mockCommandError = nil
	defer func() {
		mockCommandError = nil
	}()
	defer setupMockExecCommandWithClear()()

	svc := NewTailscaleService()
	err := svc.ClearService("my-service")

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestClearService_Failure(t *testing.T) {
	mockCommandError = errors.New("clear failed")
	defer func() {
		mockCommandError = nil
	}()
	defer setupMockExecCommandWithClear()()

	svc := NewTailscaleService()
	err := svc.ClearService("my-service")

	if err == nil {
		t.Fatal("expected error, got nil")
	}

	clearErr, ok := err.(*CommandError)
	if !ok {
		t.Fatalf("expected CommandError, got %T", err)
	}
	if clearErr.Message != "clear failed" {
		t.Errorf("expected message 'clear failed', got '%s'", clearErr.Message)
	}
}

var mockRemoveCommandError error
var capturedRemoveArgs []string

func setupMockExecCommandWithEndpoint() func() {
	oldExecCommand := execCommand
	execCommand = func(name string, args ...string) interface {
		Output() ([]byte, error)
		CombinedOutput() ([]byte, error)
	} {
		argsStr := strings.Join(args, " ")
		switch {
		case strings.Contains(argsStr, "serve status"):
			if mockServeOutput != nil {
				return &mockCmd{output: mockServeOutput}
			}
			return &mockCmd{output: []byte(`{"Services": {}}`)}
		case strings.Contains(argsStr, "serve --service=") && strings.Contains(argsStr, " off"):
			capturedRemoveArgs = args
			if mockRemoveCommandError != nil {
				return &mockCmd{err: mockRemoveCommandError, output: []byte("remove failed")}
			}
			return &mockCmd{output: []byte("success")}
		case strings.Contains(argsStr, "serve --service="):
			if mockCommandError != nil {
				return &mockCmd{err: mockCommandError, output: []byte("command failed")}
			}
			return &mockCmd{output: []byte("success")}
		case strings.Contains(argsStr, "serve clear"):
			if mockCommandError != nil {
				return &mockCmd{err: mockCommandError, output: []byte("clear failed")}
			}
			return &mockCmd{output: []byte("success")}
		}
		return &mockCmd{err: errors.New("unexpected command")}
	}
	return func() { execCommand = oldExecCommand }
}

func TestAddEndpoint_Success(t *testing.T) {
	mockCommandError = nil
	defer func() {
		mockCommandError = nil
	}()
	defer setupMockExecCommandWithEndpoint()()

	svc := NewTailscaleService()
	params := EndpointParams{
		ServiceName: "my-service",
		Protocol:    "https",
		ExposePort:  "443",
		Destination: "http://localhost:8080",
	}
	err := svc.AddEndpoint(params)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestAddEndpoint_Failure(t *testing.T) {
	mockCommandError = errors.New("command failed")
	defer func() {
		mockCommandError = nil
	}()
	defer setupMockExecCommandWithEndpoint()()

	svc := NewTailscaleService()
	params := EndpointParams{
		ServiceName: "my-service",
		Protocol:    "https",
		ExposePort:  "443",
		Destination: "http://localhost:8080",
	}
	err := svc.AddEndpoint(params)

	if err == nil {
		t.Fatal("expected error, got nil")
	}

	endpointErr, ok := err.(*CommandError)
	if !ok {
		t.Fatalf("expected CommandError, got %T", err)
	}
	if endpointErr.Message != "command failed" {
		t.Errorf("expected message 'command failed', got '%s'", endpointErr.Message)
	}
}

func TestRemoveEndpoint_Success(t *testing.T) {
	mockRemoveCommandError = nil
	capturedRemoveArgs = nil
	defer func() {
		mockRemoveCommandError = nil
		capturedRemoveArgs = nil
	}()
	defer setupMockExecCommandWithEndpoint()()

	svc := NewTailscaleService()
	params := EndpointParams{
		ServiceName: "my-service",
		Protocol:    "https",
		ExposePort:  "443",
		Destination: "http://localhost:8080",
	}
	err := svc.RemoveEndpoint(params)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if capturedRemoveArgs == nil {
		t.Fatal("expected command to be called")
	}
	argsStr := strings.Join(capturedRemoveArgs, " ")
	if !strings.Contains(argsStr, "--service=svc:my-service") {
		t.Errorf("expected args to contain '--service=svc:my-service', got '%s'", argsStr)
	}
	if !strings.Contains(argsStr, "--https=443") {
		t.Errorf("expected args to contain '--https=443', got '%s'", argsStr)
	}
	if !strings.Contains(argsStr, "off") {
		t.Errorf("expected args to contain 'off', got '%s'", argsStr)
	}
}

func TestRemoveEndpoint_Failure(t *testing.T) {
	mockRemoveCommandError = errors.New("remove failed")
	defer func() {
		mockRemoveCommandError = nil
	}()
	defer setupMockExecCommandWithEndpoint()()

	svc := NewTailscaleService()
	params := EndpointParams{
		ServiceName: "my-service",
		Protocol:    "https",
		ExposePort:  "443",
		Destination: "http://localhost:8080",
	}
	err := svc.RemoveEndpoint(params)

	if err == nil {
		t.Fatal("expected error, got nil")
	}

	endpointErr, ok := err.(*CommandError)
	if !ok {
		t.Fatalf("expected CommandError, got %T", err)
	}
	if endpointErr.Message != "remove failed" {
		t.Errorf("expected message 'remove failed', got '%s'", endpointErr.Message)
	}
}

func TestUpdateEndpoint_Success(t *testing.T) {
	mockCommandError = nil
	mockRemoveCommandError = nil
	defer func() {
		mockCommandError = nil
		mockRemoveCommandError = nil
	}()
	defer setupMockExecCommandWithEndpoint()()

	svc := NewTailscaleService()
	params := UpdateEndpointParams{
		ServiceName:    "my-service",
		Protocol:       "https",
		ExposePort:     "443",
		OldDestination: "http://localhost:8080",
		NewDestination: "http://localhost:9000",
	}
	err := svc.UpdateEndpoint(params)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestUpdateEndpoint_RemoveFailure(t *testing.T) {
	mockRemoveCommandError = errors.New("remove failed")
	defer func() {
		mockRemoveCommandError = nil
	}()
	defer setupMockExecCommandWithEndpoint()()

	svc := NewTailscaleService()
	params := UpdateEndpointParams{
		ServiceName:    "my-service",
		Protocol:       "https",
		ExposePort:     "443",
		OldDestination: "http://localhost:8080",
		NewDestination: "http://localhost:9000",
	}
	err := svc.UpdateEndpoint(params)

	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestUpdateEndpoint_AddFailure(t *testing.T) {
	mockRemoveCommandError = nil
	mockCommandError = errors.New("add failed")
	defer func() {
		mockRemoveCommandError = nil
		mockCommandError = nil
	}()
	defer setupMockExecCommandWithEndpoint()()

	svc := NewTailscaleService()
	params := UpdateEndpointParams{
		ServiceName:    "my-service",
		Protocol:       "https",
		ExposePort:     "443",
		OldDestination: "http://localhost:8080",
		NewDestination: "http://localhost:9000",
	}
	err := svc.UpdateEndpoint(params)

	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestIsTailscaleNotInstalledError_NilError(t *testing.T) {
	result := IsTailscaleNotInstalledError(nil)
	if result {
		t.Error("expected false for nil error")
	}
}

func TestIsTailscaleNotInstalledError_DirectError(t *testing.T) {
	result := IsTailscaleNotInstalledError(ErrTailscaleNotInstalled)
	if !result {
		t.Error("expected true for ErrTailscaleNotInstalled")
	}
}

func TestIsTailscaleNotInstalledError_OtherError(t *testing.T) {
	result := IsTailscaleNotInstalledError(errors.New("some other error"))
	if result {
		t.Error("expected false for unrelated error")
	}
}

func TestCheckInstalled_Success(t *testing.T) {
	oldExecCommand := execCommand
	defer func() { execCommand = oldExecCommand }()
	execCommand = func(name string, args ...string) interface {
		Output() ([]byte, error)
		CombinedOutput() ([]byte, error)
	} {
		return &mockCmd{output: []byte("1.50.0")}
	}

	svc := NewTailscaleService()
	err := svc.CheckInstalled()

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestCheckInstalled_CommandFailed(t *testing.T) {
	oldExecCommand := execCommand
	defer func() { execCommand = oldExecCommand }()
	execCommand = func(name string, args ...string) interface {
		Output() ([]byte, error)
		CombinedOutput() ([]byte, error)
	} {
		return &mockCmd{err: errors.New("command failed")}
	}

	svc := NewTailscaleService()
	err := svc.CheckInstalled()

	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestCommandError_ErrorWithMessage(t *testing.T) {
	err := &CommandError{
		Message: "custom error message",
		Err:     errors.New("underlying error"),
	}

	if err.Error() != "custom error message" {
		t.Errorf("expected 'custom error message', got '%s'", err.Error())
	}
}

func TestCommandError_ErrorWithoutMessage(t *testing.T) {
	err := &CommandError{
		Message: "",
		Err:     errors.New("underlying error"),
	}

	if err.Error() != "underlying error" {
		t.Errorf("expected 'underlying error', got '%s'", err.Error())
	}
}

func TestGetServeStatus_MultipleServices(t *testing.T) {
	jsonData := `{
		"Services": {
			"svc:web-app": {
				"Web": {
					"web.example.com:443": {
						"Handlers": {
							"/": {"Proxy": "http://localhost:3000"}
						}
					}
				}
			},
			"svc:api-server": {
				"Web": {
					"api.example.com:443": {
						"Handlers": {
							"/": {"Proxy": "http://localhost:8080"}
						}
					}
				}
			},
			"svc:db-proxy": {
				"Web": {
					"db.example.com:5432": {
						"Handlers": {
							"/": {"Proxy": "localhost:5432"}
						}
					}
				}
			}
		}
	}`
	mockServeOutput = []byte(jsonData)
	mockCommandError = nil
	defer func() {
		mockServeOutput = nil
		mockCommandError = nil
	}()
	defer setupMockExecCommand()()

	svc := NewTailscaleService()
	services, err := svc.GetServeStatus()

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(services) != 3 {
		t.Fatalf("expected 3 services, got %d", len(services))
	}
	expectedOrder := []string{"api-server", "db-proxy", "web-app"}
	for i, expected := range expectedOrder {
		if services[i].Name != expected {
			t.Errorf("expected services[%d].Name = '%s', got '%s'", i, expected, services[i].Name)
		}
	}
}

func TestGetServeStatus_NullServices(t *testing.T) {
	jsonData := `{"Services": null}`
	mockServeOutput = []byte(jsonData)
	mockCommandError = nil
	defer func() {
		mockServeOutput = nil
		mockCommandError = nil
	}()
	defer setupMockExecCommand()()

	svc := NewTailscaleService()
	services, err := svc.GetServeStatus()

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(services) != 0 {
		t.Fatalf("expected 0 services, got %d", len(services))
	}
}

func TestGetServeStatus_EmptyWebHandlers(t *testing.T) {
	jsonData := `{
		"Services": {
			"svc:empty-web": {
				"Web": {}
			}
		}
	}`
	mockServeOutput = []byte(jsonData)
	mockCommandError = nil
	defer func() {
		mockServeOutput = nil
		mockCommandError = nil
	}()
	defer setupMockExecCommand()()

	svc := NewTailscaleService()
	services, err := svc.GetServeStatus()

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(services) != 1 {
		t.Fatalf("expected 1 service, got %d", len(services))
	}
	if services[0].Name != "empty-web" {
		t.Errorf("expected name 'empty-web', got '%s'", services[0].Name)
	}
	if services[0].HTTPSUrl != "" {
		t.Errorf("expected empty HTTPSUrl, got '%s'", services[0].HTTPSUrl)
	}
	if services[0].Proxy != "" {
		t.Errorf("expected empty Proxy, got '%s'", services[0].Proxy)
	}
}

func TestGetServeStatus_EmptyHandlers(t *testing.T) {
	jsonData := `{
		"Services": {
			"svc:no-handlers": {
				"Web": {
					"example.com:443": {
						"Handlers": {}
					}
				}
			}
		}
	}`
	mockServeOutput = []byte(jsonData)
	mockCommandError = nil
	defer func() {
		mockServeOutput = nil
		mockCommandError = nil
	}()
	defer setupMockExecCommand()()

	svc := NewTailscaleService()
	services, err := svc.GetServeStatus()

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(services) != 1 {
		t.Fatalf("expected 1 service, got %d", len(services))
	}
	if services[0].Name != "no-handlers" {
		t.Errorf("expected name 'no-handlers', got '%s'", services[0].Name)
	}
	if services[0].Proxy != "" {
		t.Errorf("expected empty Proxy, got '%s'", services[0].Proxy)
	}
}

func TestGetServeStatus_InvalidHostFormat(t *testing.T) {
	jsonData := `{
		"Services": {
			"svc:invalid-host": {
				"Web": {
					"example.com": {
						"Handlers": {
							"/": {"Proxy": "http://localhost:3000"}
						}
					}
				}
			}
		}
	}`
	mockServeOutput = []byte(jsonData)
	mockCommandError = nil
	defer func() {
		mockServeOutput = nil
		mockCommandError = nil
	}()
	defer setupMockExecCommand()()

	svc := NewTailscaleService()
	services, err := svc.GetServeStatus()

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(services) != 1 {
		t.Fatalf("expected 1 service, got %d", len(services))
	}
	if services[0].Name != "invalid-host" {
		t.Errorf("expected name 'invalid-host', got '%s'", services[0].Name)
	}
	if services[0].HTTPSUrl != "" {
		t.Errorf("expected empty HTTPSUrl, got '%s'", services[0].HTTPSUrl)
	}
}
