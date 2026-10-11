package api

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

/*
付費憑證 : v1.<userID>.<到期 unix 秒>.<HMAC-SHA256 hex>
由伺服器以 premium_secret 簽發, 手機放在 X-Premium-Token 標頭
*/

const (
	HeaderPremiumToken = "X-Premium-Token"
	premiumTokenPrefix = "v1"
	ctxPremium         = "premium"
	ctxPremiumUser     = "premium_user"
	ctxPremiumExpires  = "premium_expires"
)

var premiumUserPattern = regexp.MustCompile(`^[A-Za-z0-9_-]{1,64}$`)

var (
	ErrPremiumTokenInvalid = errors.New("premium token invalid")
	ErrPremiumTokenExpired = errors.New("premium token expired")
)

func signPremium(secret, payload string) string {
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(payload))
	return hex.EncodeToString(mac.Sum(nil))
}

// IssuePremiumToken 簽發付費憑證
func IssuePremiumToken(secret, user string, expires time.Time) (string, error) {
	if secret == "" {
		return "", errors.New("premium_secret is not configured")
	}
	if !premiumUserPattern.MatchString(user) {
		return "", errors.New("user must be 1-64 characters of A-Z a-z 0-9 _ -")
	}
	payload := fmt.Sprintf("%s.%s.%d", premiumTokenPrefix, user, expires.Unix())
	return payload + "." + signPremium(secret, payload), nil
}

// VerifyPremiumToken 驗證付費憑證, 回傳使用者與到期時間
func VerifyPremiumToken(secret, token string, now time.Time) (string, time.Time, error) {
	parts := strings.Split(token, ".")
	if secret == "" || len(parts) != 4 || parts[0] != premiumTokenPrefix || !premiumUserPattern.MatchString(parts[1]) {
		return "", time.Time{}, ErrPremiumTokenInvalid
	}
	payload := strings.Join(parts[:3], ".")
	if !hmac.Equal([]byte(parts[3]), []byte(signPremium(secret, payload))) {
		return "", time.Time{}, ErrPremiumTokenInvalid
	}
	exp, err := strconv.ParseInt(parts[2], 10, 64)
	if err != nil {
		return "", time.Time{}, ErrPremiumTokenInvalid
	}
	expires := time.Unix(exp, 0)
	if !now.Before(expires) {
		return parts[1], expires, ErrPremiumTokenExpired
	}
	return parts[1], expires, nil
}

// PremiumAuth 解析 X-Premium-Token. 沒帶視為免費用戶; 帶了但無效或過期回 401, 讓 App 知道要更新憑證
func PremiumAuth(secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader(HeaderPremiumToken)
		if token == "" {
			c.Next()
			return
		}
		user, expires, err := VerifyPremiumToken(secret, token, time.Now())
		if err != nil {
			Fail(c, http.StatusUnauthorized, CodePremiumTokenInvalid, err.Error())
			return
		}
		c.Set(ctxPremium, true)
		c.Set(ctxPremiumUser, user)
		c.Set(ctxPremiumExpires, expires)
		c.Next()
	}
}

// IsPremium 目前請求是否為付費用戶
func IsPremium(c *gin.Context) bool {
	return c.GetBool(ctxPremium)
}

// Entitlement 目前用戶的權限
type Entitlement struct {
	Premium          bool   `json:"premium"`
	User             string `json:"user,omitempty"`
	PremiumExpiresAt string `json:"premiumExpiresAt,omitempty"`
	FreeDrawLimit    int    `json:"freeDrawLimit"` // 免費用戶可看的最新期數, 0 表示不限
}

// MeHandler GET /me
func MeHandler(freeDrawLimit int) gin.HandlerFunc {
	return func(c *gin.Context) {
		e := Entitlement{Premium: IsPremium(c), FreeDrawLimit: freeDrawLimit}
		if e.Premium {
			e.User = c.GetString(ctxPremiumUser)
			if t, ok := c.Get(ctxPremiumExpires); ok {
				e.PremiumExpiresAt = t.(time.Time).Format(time.RFC3339)
			}
		}
		OK(c, e)
	}
}
