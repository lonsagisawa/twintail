package services

import (
	"encoding/json"
	"os/exec"
	"strings"
)

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
	TCP map[string]TCPEntry  `json:"TCP"`
	Web map[string]WebEntry  `json:"Web"`
}

type ServeStatus struct {
	Services map[string]Service `json:"Services"`
}

type ServiceView struct {
	Name     string
	HTTPSUrl string
	Proxy    string
}

type TailscaleService struct{}

func NewTailscaleService() *TailscaleService {
	return &TailscaleService{}
}

func (s *TailscaleService) GetServeStatus() ([]ServiceView, error) {
	cmd := exec.Command("tailscale", "serve", "status", "--json")
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
		var httpsUrl, proxy string

		for host, web := range svc.Web {
			if strings.Contains(host, ":443") {
				httpsUrl = "https://" + strings.TrimSuffix(host, ":443")
				for _, handler := range web.Handlers {
					if handler.Proxy != "" {
						proxy = handler.Proxy
						break
					}
				}
				break
			}
		}

		services = append(services, ServiceView{
			Name:     displayName,
			HTTPSUrl: httpsUrl,
			Proxy:    proxy,
		})
	}

	return services, nil
}
