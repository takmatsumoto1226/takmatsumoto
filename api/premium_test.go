package api

import (
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestPremiumToken(t *testing.T) {
	a := assert.New(t)
	now := time.Unix(1_800_000_000, 0)
	token, err := IssuePremiumToken("s1", "user_1", now.Add(time.Hour))
	a.NoError(err)

	user, exp, err := VerifyPremiumToken("s1", token, now)
	a.NoError(err)
	a.Equal("user_1", user)
	a.Equal(now.Add(time.Hour).Unix(), exp.Unix())

	_, _, err = VerifyPremiumToken("s1", token, now.Add(time.Hour))
	a.ErrorIs(err, ErrPremiumTokenExpired)

	_, _, err = VerifyPremiumToken("s2", token, now)
	a.ErrorIs(err, ErrPremiumTokenInvalid)

	// 竄改使用者或到期時間
	forged := strings.Replace(token, "user_1", "user_2", 1)
	_, _, err = VerifyPremiumToken("s1", forged, now)
	a.ErrorIs(err, ErrPremiumTokenInvalid)

	_, _, err = VerifyPremiumToken("", token, now)
	a.ErrorIs(err, ErrPremiumTokenInvalid)

	_, err = IssuePremiumToken("", "user", now)
	a.Error(err)
	_, err = IssuePremiumToken("s1", "bad.user", now)
	a.Error(err)
}
