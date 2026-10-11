package api

import (
	"crypto/subtle"
	"net/http"

	"github.com/gin-gonic/gin"
)

// 業務錯誤碼
const (
	CodeOK           = 0
	CodeInvalidParam = 40001
	CodeUnauthorized = 40101
	// 付費憑證無效或過期
	CodePremiumTokenInvalid = 40102
	// 超出免費範圍, 需付費
	CodePremiumRequired = 40301
	CodeNotFound        = 40401
	CodeDataNotReady    = 50301
	HeaderAPIKey        = "X-API-Key"
	MessageOK           = "ok"
	MessageUnauthorize  = "invalid or missing api key"
)

// Response 統一回傳格式
type Response struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

// OK 回傳成功
func OK(c *gin.Context, data any) {
	c.JSON(http.StatusOK, Response{Code: CodeOK, Message: MessageOK, Data: data})
}

// Fail 回傳失敗
func Fail(c *gin.Context, status int, code int, msg string) {
	c.AbortWithStatusJSON(status, Response{Code: code, Message: msg})
}

// APIKeyAuth 檢查 X-API-Key, key 為空時不檢查
func APIKeyAuth(key string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if key == "" {
			c.Next()
			return
		}
		got := c.GetHeader(HeaderAPIKey)
		if subtle.ConstantTimeCompare([]byte(got), []byte(key)) != 1 {
			Fail(c, http.StatusUnauthorized, CodeUnauthorized, MessageUnauthorize)
			return
		}
		c.Next()
	}
}
