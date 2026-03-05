package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/scmbr/device-tsv-processor/internal/errs"
	"github.com/scmbr/device-tsv-processor/pkg/logger"
)

type APIError struct {
	Error  string            `json:"error"`
	Kind   errs.Kind         `json:"kind"`
	Code   string            `json:"code,omitempty"`
	Fields map[string]string `json:"fields,omitempty"`
}

func Error() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Next()

		if len(ctx.Errors) == 0 {
			return
		}

		raw := ctx.Errors.Last().Err
		if raw == nil || ctx.Writer.Written() {
			return
		}

		status, resp, logLevel := mapErr(raw)

		switch logLevel {
		case "error":
			logger.Error(resp.Error, raw, map[string]interface{}{
				"kind": resp.Kind,
				"code": resp.Code,
			})
		case "warn":
			logger.Warn(resp.Error, map[string]interface{}{
				"kind": resp.Kind,
				"code": resp.Code,
			})
		}

		ctx.AbortWithStatusJSON(status, resp)
	}
}

func mapErr(err error) (status int, resp APIError, logLevel string) {
	if e, ok := errs.As(err); ok {
		resp.Kind = e.Kind
		resp.Code = e.Code

		switch e.Kind {
		case errs.KindInvalid:
			resp.Error = "invalid request"
			resp.Fields = e.Fields
			return http.StatusBadRequest, resp, "warn"

		case errs.KindNotFound:
			resp.Error = "not found"
			return http.StatusNotFound, resp, "warn"

		case errs.KindForbidden:
			resp.Error = "forbidden"
			return http.StatusForbidden, resp, "warn"

		case errs.KindConflict:
			resp.Error = "conflict"
			return http.StatusConflict, resp, "warn"

		default:
			resp.Error = "internal server error"
			return http.StatusInternalServerError, resp, "error"
		}
	}

	resp.Kind = errs.KindInternal
	resp.Error = "internal server error"
	return http.StatusInternalServerError, resp, "error"
}
