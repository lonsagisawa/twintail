package services

import (
	"encoding/json"
	"os/exec"
	"strings"
)

var execCommand = func(name string, arg ...string) interface {
	Output() ([]byte, error)
	CombinedOutput() ([]byte, error)
} {
	return exec.Command(name, arg...)
}

type Handler struct {
	Proxy string `json:"Proxy,omitempty"`
}

type WebEntry struct {
	Handlers map[string]Handler `json:"Handlers"`
}

type TCPEntry struct {
	HTTP  bool `json:"HTTP,omitempty"`
	HTTPS bool `json:"HTTPS,omitempty"`
}

type Service struct {
	TCP map[string]TCPEntry `json:"TCP"`
	Web map[string]WebEntry `json:"Web"`
}

type ServeStatus struct {
	Services map[string]Service `json:"Services"`
}

type ServiceView struct {
	Name     string
	HTTPSUrl string
	HTTPUrl  string
	Proxy    string
}

type PortEntry struct {
	Protocol    string
	ExposePort  string
	Destination string
}

type ServiceDetailView struct {
	Name     string
	Hostname string
	URL      string
	Ports    []PortEntry
}

type TailscaleService struct{}

func NewTailscaleService() *TailscaleService {
	return &TailscaleService{}
}

func (s *TailscaleService) GetServeStatus() ([]ServiceView, error) {
	cmd := execCommand("tailscale", "serve", "status", "--json")
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	var status ServeStatus
	if err := json.Unmarshal(output, &status); err != nil {
		return nil, err
	}

	var services []ServiceView
	for name, svc := range status.Services {
		displayName := strings.TrimPrefix(name, "svc:")
		var httpsUrl, httpUrl, proxy string

		for host, web := range svc.Web {
			parts := strings.Split(host, ":")
			if len(parts) != 2 {
				continue
			}
			hostname := parts[0]
			port := parts[1]

			for _, handler := range web.Handlers {
				if handler.Proxy != "" && proxy == "" {
					proxy = handler.Proxy
				}
			}

			if port == "443" {
				httpsUrl = "https://" + hostname
			} else if port == "80" {
				httpUrl = "http://" + hostname
			} else {
				httpUrl = "http://" + hostname + ":" + port
			}
		}

		services = append(services, ServiceView{
			Name:     displayName,
			HTTPSUrl: httpsUrl,
			HTTPUrl:  httpUrl,
			Proxy:    proxy,
		})
	}

	return services, nil
}

func (s *TailscaleService) GetServiceByName(name string) (*ServiceDetailView, error) {
	cmd := execCommand("tailscale", "serve", "status", "--json")
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	var status ServeStatus
	if err := json.Unmarshal(output, &status); err != nil {
		return nil, err
	}

	svcKey := "svc:" + name
	svc, ok := status.Services[svcKey]
	if !ok {
		return nil, nil
	}

	detail := &ServiceDetailView{
		Name: name,
	}

	var hasHTTPS, hasHTTP bool
	var httpPort string

	for host, web := range svc.Web {
		parts := strings.Split(host, ":")
		if len(parts) == 2 && detail.Hostname == "" {
			detail.Hostname = parts[0]
		}

		port := ""
		protocol := "http"
		if len(parts) == 2 {
			port = parts[1]
			if port == "443" {
				protocol = "https"
				hasHTTPS = true
			} else {
				hasHTTP = true
				httpPort = port
			}
		}

		for _, handler := range web.Handlers {
			if handler.Proxy != "" {
				detail.Ports = append(detail.Ports, PortEntry{
					Protocol:    protocol,
					ExposePort:  port,
					Destination: handler.Proxy,
				})
			}
		}
	}

	if hasHTTPS {
		detail.URL = "https://" + detail.Hostname
	} else if hasHTTP {
		if httpPort == "80" {
			detail.URL = "http://" + detail.Hostname
		} else {
			detail.URL = "http://" + detail.Hostname + ":" + httpPort
		}
	}

	return detail, nil
}

type AdvertiseServiceParams struct {
	ServiceName string
	Protocol    string
	ExposePort  string
	Destination string
}

func (s *TailscaleService) AdvertiseService(params AdvertiseServiceParams) error {
	args := []string{
		"serve",
		"--service=svc:" + params.ServiceName,
		"--" + params.Protocol + "=" + params.ExposePort,
		params.Destination,
	}

	cmd := execCommand("tailscale", args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return &AdvertiseError{
			Message: string(output),
			Err:     err,
		}
	}
	return nil
}

type AdvertiseError struct {
	Message string
	Err     error
}

func (e *AdvertiseError) Error() string {
	if e.Message != "" {
		return e.Message
	}
	return e.Err.Error()
}

func (s *TailscaleService) ClearService(name string) error {
	cmd := execCommand("tailscale", "serve", "clear", "svc:"+name)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return &ClearError{
			Message: string(output),
			Err:     err,
		}
	}
	return nil
}

type ClearError struct {
	Message string
	Err     error
}

func (e *ClearError) Error() string {
	if e.Message != "" {
		return e.Message
	}
	return e.Err.Error()
}

type EndpointParams struct {
	ServiceName string
	Protocol    string
	ExposePort  string
	Destination string
}

type EndpointError struct {
	Message string
	Err     error
}

func (e *EndpointError) Error() string {
	if e.Message != "" {
		return e.Message
	}
	return e.Err.Error()
}

func (s *TailscaleService) AddEndpoint(params EndpointParams) error {
	args := []string{
		"serve",
		"--service=svc:" + params.ServiceName,
		"--" + params.Protocol + "=" + params.ExposePort,
		params.Destination,
	}

	cmd := execCommand("tailscale", args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return &EndpointError{
			Message: string(output),
			Err:     err,
		}
	}
	return nil
}

func (s *TailscaleService) RemoveEndpoint(params EndpointParams) error {
	args := []string{
		"serve",
		"--service=svc:" + params.ServiceName,
		"--" + params.Protocol + "=" + params.ExposePort,
		params.Destination,
		"off",
	}

	cmd := execCommand("tailscale", args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return &EndpointError{
			Message: string(output),
			Err:     err,
		}
	}
	return nil
}

type UpdateEndpointParams struct {
	ServiceName    string
	Protocol       string
	ExposePort     string
	OldDestination string
	NewDestination string
}

func (s *TailscaleService) UpdateEndpoint(params UpdateEndpointParams) error {
	removeParams := EndpointParams{
		ServiceName: params.ServiceName,
		Protocol:    params.Protocol,
		ExposePort:  params.ExposePort,
		Destination: params.OldDestination,
	}
	if err := s.RemoveEndpoint(removeParams); err != nil {
		return err
	}

	addParams := EndpointParams{
		ServiceName: params.ServiceName,
		Protocol:    params.Protocol,
		ExposePort:  params.ExposePort,
		Destination: params.NewDestination,
	}
	return s.AddEndpoint(addParams)
}
