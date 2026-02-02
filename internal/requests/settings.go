package requests

import "github.com/labstack/echo/v5"

type UpdateSettingsRequest struct {
	Lang string `form:"lang" validate:"required,oneof=en ja"`
}

func (r *UpdateSettingsRequest) FromContext(ctx *echo.Context) error {
	if err := ctx.Bind(r); err != nil {
		return err
	}
	return ctx.Validate(r)
}
