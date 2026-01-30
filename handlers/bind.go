package handlers

import "github.com/labstack/echo/v5"

// BindAndValidate はバインドとバリデーションを一括で行うヘルパー関数
// エラーがあった場合はerrorを返す。ハンドラー側でエラーハンドリングを行う。
func BindAndValidate(ctx *echo.Context, req any) error {
	if err := ctx.Bind(req); err != nil {
		return err
	}
	if err := ctx.Validate(req); err != nil {
		return err
	}
	return nil
}
