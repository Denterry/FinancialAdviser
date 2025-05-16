package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/Denterry/FinancialAdviser/Backend/gateway-service/config"
	"github.com/Denterry/FinancialAdviser/Backend/gateway-service/internal/service"
	"github.com/Denterry/FinancialAdviser/Backend/gateway-service/pkg/logger"
	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

// Logger middleware for request logging
func Logger(l *logger.Logger) gin.HandlerFunc {
	return gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		l.Info("%s - [%s] \"%s %s %s %d %s \"%s\" %s\"",
			param.ClientIP,
			param.TimeStamp.Format(time.RFC1123),
			param.Method,
			param.Path,
			param.Request.Proto,
			param.StatusCode,
			param.Latency,
			param.Request.UserAgent(),
			param.ErrorMessage,
		)
		return ""
	})
}

// Recovery middleware for panic recovery
func Recovery() gin.HandlerFunc {
	return gin.CustomRecovery(func(c *gin.Context, recovered interface{}) {
		c.AbortWithStatus(http.StatusInternalServerError)
		fmt.Printf("[PANIC] %s | %v\n", time.Now().Format(time.RFC3339), recovered)
	})
}

// CORS middleware for handling Cross-Origin Resource Sharing
func CORS() gin.HandlerFunc {
	allowOrigin := "*"

	return func(c *gin.Context) {
		h := c.Writer.Header()
		h.Set("Access-Control-Allow-Origin", allowOrigin)
		h.Set("Access-Control-Allow-Credentials", "true")

		// Echo requested headers/methods if present.
		if reqHeaders := c.GetHeader("Access-Control-Request-Headers"); reqHeaders != "" {
			h.Set("Access-Control-Allow-Headers", reqHeaders)
		} else {
			h.Set("Access-Control-Allow-Headers", "Authorization, Content-Type")
		}

		if reqMethods := c.GetHeader("Access-Control-Request-Methods"); reqMethods != "" {
			h.Set("Access-Control-Allow-Methods", reqMethods)
		} else {
			h.Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		}

		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		c.Next()
	}
}

// RateLimit middleware for request rate limiting
func RateLimit(cfg *config.Config) gin.HandlerFunc {
	rps := rate.Limit(cfg.HTTP.GIN.RateLimit) // tokens / second
	burst := int(cfg.HTTP.GIN.RateLimit)      // same as rps for now

	limiters := make(map[string]*rate.Limiter)
	var lastPurge time.Time

	return func(c *gin.Context) {
		ip := c.ClientIP()

		// purge old entries once a minute
		if time.Since(lastPurge) > time.Minute {
			for k, l := range limiters {
				if l.Allow() { // rough heuristic
					delete(limiters, k)
				}
			}
			lastPurge = time.Now()
		}

		l, ok := limiters[ip]
		if !ok {
			l = rate.NewLimiter(rps, burst)
			limiters[ip] = l
		}

		if !l.Allow() {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"error": "rate limit exceeded",
			})
			return
		}

		c.Next()
	}
}

// Auth middleware for JWT authentication
func Auth(authService *service.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		raw := c.GetHeader("Authorization")
		if raw == "" {
			unauth(c, "missing Authorization header")
			return
		}

		const bearer = "Bearer "
		if strings.HasPrefix(raw, bearer) {
			raw = strings.TrimPrefix(raw, bearer)
		}

		claims, err := authService.ValidateToken(context.Background(), raw)
		if err != nil {
			unauth(c, err.Error())
			return
		}

		// Make claims available to downstream handlers.
		c.Set("claims", claims)
		c.Next()
	}
}

func unauth(c *gin.Context, msg string) {
	c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": msg})
}
