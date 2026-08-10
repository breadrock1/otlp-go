package main

import (
	"github.com/gofiber/fiber/v2"
)

const (
	DefaultUserIDHeader = "unknown"
	userIDHeaderKey     = "x_user_id"
	userIDHeader        = "X-User-Id"
)

func UserContext() fiber.Handler {
	return func(eCtx *fiber.Ctx) error {
		authHeaderValue := eCtx.Get(userIDHeader)

		if authHeaderValue == "" {
			authHeaderValue = DefaultUserIDHeader
		}

		eCtx.Locals(userIDHeaderKey, authHeaderValue)

		return eCtx.Next()
	}
}
