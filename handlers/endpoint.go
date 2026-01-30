package handlers

import (
	"twintail/services"

	"github.com/labstack/echo/v5"
)

type EndpointService interface {
	GetServiceByName(name string) (*services.ServiceDetailView, error)
	AddEndpoint(params services.EndpointParams) error
	RemoveEndpoint(params services.EndpointParams) error
	UpdateEndpoint(params services.UpdateEndpointParams) error
}

type EndpointHandler struct {
	tailscale EndpointService
}

func NewEndpointHandler(tailscale EndpointService) *EndpointHandler {
	return &EndpointHandler{
		tailscale: tailscale,
	}
}

type EndpointFormData struct {
	Protocol    string `form:"protocol" validate:"required,oneof=https http tcp+tls tcp"`
	ExposePort  string `form:"expose_port" validate:"required,numeric"`
	Destination string `form:"destination" validate:"required"`
}

func (h *EndpointHandler) Create(ctx *echo.Context) error {
	name := ctx.Param("name")
	return ctx.Render(200, "new_endpoint.html", map[string]any{
		"ServiceName": name,
		"FormData":    EndpointFormData{Protocol: "https", ExposePort: "443"},
	})
}

func (h *EndpointHandler) Store(ctx *echo.Context) error {
	name := ctx.Param("name")
	var formData EndpointFormData
	if err := ctx.Bind(&formData); err != nil {
		return ctx.Render(200, "new_endpoint.html", map[string]any{
			"ServiceName": name,
			"Error":       err.Error(),
			"FormData":    formData,
		})
	}
	if err := ctx.Validate(&formData); err != nil {
		return ctx.Render(200, "new_endpoint.html", map[string]any{
			"ServiceName": name,
			"Error":       err.Error(),
			"FormData":    formData,
		})
	}

	params := services.EndpointParams{
		ServiceName: name,
		Protocol:    formData.Protocol,
		ExposePort:  formData.ExposePort,
		Destination: formData.Destination,
	}

	if err := h.tailscale.AddEndpoint(params); err != nil {
		return ctx.Render(200, "new_endpoint.html", map[string]any{
			"ServiceName": name,
			"Error":       err.Error(),
			"FormData":    formData,
		})
	}

	return ctx.Redirect(303, "/services/"+name)
}

func (h *EndpointHandler) Delete(ctx *echo.Context) error {
	name := ctx.Param("name")
	protocol := ctx.QueryParam("protocol")
	exposePort := ctx.QueryParam("port")
	destination := ctx.QueryParam("destination")

	return ctx.Render(200, "confirm_delete_endpoint.html", map[string]any{
		"ServiceName": name,
		"Protocol":    protocol,
		"ExposePort":  exposePort,
		"Destination": destination,
	})
}

func (h *EndpointHandler) Destroy(ctx *echo.Context) error {
	name := ctx.Param("name")
	params := services.EndpointParams{
		ServiceName: name,
		Protocol:    ctx.FormValue("protocol"),
		ExposePort:  ctx.FormValue("expose_port"),
		Destination: ctx.FormValue("destination"),
	}

	if err := h.tailscale.RemoveEndpoint(params); err != nil {
		return ctx.String(500, "Failed to delete endpoint: "+err.Error())
	}

	svc, _ := h.tailscale.GetServiceByName(name)
	if svc == nil {
		return ctx.Redirect(303, "/")
	}

	return ctx.Redirect(303, "/services/"+name)
}

type EditEndpointFormData struct {
	Protocol       string `form:"protocol" validate:"required,oneof=https http tcp+tls tcp"`
	ExposePort     string `form:"expose_port" validate:"required,numeric"`
	OldDestination string `form:"old_destination" validate:"required"`
	NewDestination string `form:"new_destination" validate:"required"`
}

func (h *EndpointHandler) Edit(ctx *echo.Context) error {
	name := ctx.Param("name")
	protocol := ctx.QueryParam("protocol")
	exposePort := ctx.QueryParam("port")
	destination := ctx.QueryParam("destination")

	return ctx.Render(200, "edit_endpoint.html", map[string]any{
		"ServiceName": name,
		"FormData": EditEndpointFormData{
			Protocol:       protocol,
			ExposePort:     exposePort,
			OldDestination: destination,
			NewDestination: destination,
		},
	})
}

func (h *EndpointHandler) Update(ctx *echo.Context) error {
	name := ctx.Param("name")
	var formData EditEndpointFormData
	if err := ctx.Bind(&formData); err != nil {
		return ctx.Render(200, "edit_endpoint.html", map[string]any{
			"ServiceName": name,
			"Error":       err.Error(),
			"FormData":    formData,
		})
	}
	if err := ctx.Validate(&formData); err != nil {
		return ctx.Render(200, "edit_endpoint.html", map[string]any{
			"ServiceName": name,
			"Error":       err.Error(),
			"FormData":    formData,
		})
	}

	params := services.UpdateEndpointParams{
		ServiceName:    name,
		Protocol:       formData.Protocol,
		ExposePort:     formData.ExposePort,
		OldDestination: formData.OldDestination,
		NewDestination: formData.NewDestination,
	}

	if err := h.tailscale.UpdateEndpoint(params); err != nil {
		return ctx.Render(200, "edit_endpoint.html", map[string]any{
			"ServiceName": name,
			"Error":       err.Error(),
			"FormData":    formData,
		})
	}

	return ctx.Redirect(303, "/services/"+name)
}
