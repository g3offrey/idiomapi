package middleware

import (
	"log/slog"
	"net/http"

	"github.com/g3offrey/idiomapi/internal/model"
	"github.com/gin-gonic/gin"
)

// Recovery returns a gin middleware that recovers from panics and logs them using slog
func Recovery(logger *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				logger.Error("panic recovered",
					"error", err,
					"path", c.Request.URL.Path,
					"method", c.Request.Method,
				)

				c.AbortWithStatusJSON(http.StatusInternalServerError, model.ErrorResponse{
					Error:   "internal_server_error",
					Message: "An unexpected error occurred",
				})
			}
		}()

		c.Next()
	}
}
