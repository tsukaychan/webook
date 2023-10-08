package captcha

import (
	"context"
	_ "embed"
	"errors"
)

var (
	ErrSetCaptchaTooManyTimes    = errors.New("send captcha too many times")
	ErrInternal                  = errors.New("internal error")
	ErrCaptchaVerifyTooManyTimes = errors.New("verify captcha too many times")
	ErrUnknownForCode            = errors.New("unknown error for code")
)

//go:embed lua/set_code.lua
var luaSetCaptcha string

//go:embed lua/verify_code.lua
var luaVerifyCode string

var _ CaptchaCache = (*CaptchaRedisCache)(nil)

type CaptchaCache interface {
	Set(ctx context.Context, biz string, phone string, code string) error
	Verify(ctx context.Context, biz string, phone string, inputCaptcha string) (bool, error)
}