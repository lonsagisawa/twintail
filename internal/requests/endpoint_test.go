package requests

import (
	"testing"

	"github.com/go-playground/validator/v10"
)

func TestStoreEndpointRequest_Validation(t *testing.T) {
	v := validator.New()

	tests := []struct {
		name    string
		req     StoreEndpointRequest
		wantErr bool
	}{
		{
			name: "valid request",
			req: StoreEndpointRequest{
				Protocol:    "https",
				ExposePort:  "443",
				Destination: "http://localhost:8080",
			},
			wantErr: false,
		},
		{
			name: "invalid protocol",
			req: StoreEndpointRequest{
				Protocol:    "ftp",
				ExposePort:  "443",
				Destination: "http://localhost:8080",
			},
			wantErr: true,
		},
		{
			name: "non-numeric port",
			req: StoreEndpointRequest{
				Protocol:    "https",
				ExposePort:  "abc",
				Destination: "http://localhost:8080",
			},
			wantErr: true,
		},
		{
			name: "missing destination",
			req: StoreEndpointRequest{
				Protocol:   "https",
				ExposePort: "443",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := v.Struct(tt.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("Validation() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestStoreEndpointRequest_Default(t *testing.T) {
	var req StoreEndpointRequest
	defaults := req.Default()

	if defaults.Protocol != "https" {
		t.Errorf("expected Protocol 'https', got '%s'", defaults.Protocol)
	}
	if defaults.ExposePort != "443" {
		t.Errorf("expected ExposePort '443', got '%s'", defaults.ExposePort)
	}
}

func TestStoreEndpointRequest_ToParams(t *testing.T) {
	req := StoreEndpointRequest{
		Protocol:    "https",
		ExposePort:  "443",
		Destination: "http://localhost:8080",
	}

	params := req.ToParams("my-service")

	if params.ServiceName != "my-service" {
		t.Errorf("expected ServiceName 'my-service', got '%s'", params.ServiceName)
	}
	if params.Protocol != req.Protocol {
		t.Errorf("expected Protocol '%s', got '%s'", req.Protocol, params.Protocol)
	}
}

func TestUpdateEndpointRequest_Validation(t *testing.T) {
	v := validator.New()

	tests := []struct {
		name    string
		req     UpdateEndpointRequest
		wantErr bool
	}{
		{
			name: "valid request",
			req: UpdateEndpointRequest{
				Protocol:       "https",
				ExposePort:     "443",
				OldDestination: "http://localhost:8080",
				NewDestination: "http://localhost:9090",
			},
			wantErr: false,
		},
		{
			name: "missing old destination",
			req: UpdateEndpointRequest{
				Protocol:       "https",
				ExposePort:     "443",
				NewDestination: "http://localhost:9090",
			},
			wantErr: true,
		},
		{
			name: "missing new destination",
			req: UpdateEndpointRequest{
				Protocol:       "https",
				ExposePort:     "443",
				OldDestination: "http://localhost:8080",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := v.Struct(tt.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("Validation() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestUpdateEndpointRequest_ToParams(t *testing.T) {
	req := UpdateEndpointRequest{
		Protocol:       "https",
		ExposePort:     "443",
		OldDestination: "http://localhost:8080",
		NewDestination: "http://localhost:9090",
	}

	params := req.ToParams("my-service")

	if params.ServiceName != "my-service" {
		t.Errorf("expected ServiceName 'my-service', got '%s'", params.ServiceName)
	}
	if params.OldDestination != req.OldDestination {
		t.Errorf("expected OldDestination '%s', got '%s'", req.OldDestination, params.OldDestination)
	}
	if params.NewDestination != req.NewDestination {
		t.Errorf("expected NewDestination '%s', got '%s'", req.NewDestination, params.NewDestination)
	}
}

func TestDestroyEndpointRequest_Validation(t *testing.T) {
	v := validator.New()

	tests := []struct {
		name    string
		req     DestroyEndpointRequest
		wantErr bool
	}{
		{
			name: "valid request",
			req: DestroyEndpointRequest{
				Protocol:    "https",
				ExposePort:  "443",
				Destination: "http://localhost:8080",
			},
			wantErr: false,
		},
		{
			name: "missing protocol",
			req: DestroyEndpointRequest{
				ExposePort:  "443",
				Destination: "http://localhost:8080",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := v.Struct(tt.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("Validation() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
