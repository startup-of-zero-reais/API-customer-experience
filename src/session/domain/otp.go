package domain

import (
	"log"
	"time"
)

type (
	PassTokens struct {
		Email     string `json:"email"`
		OTP       int    `json:"otp"`
		ExpiresIn int64  `json:"expires_in"`
	}

	OTPRepository interface {
		New(email string) (*PassTokens, error)
		Invalidate(email string) error
		IsValid(email string, OTP int) bool
		SearchOtp(otp int) ([]PassTokens, error)
	}

	ResetPassInput struct {
		Password        string `json:"password"`
		ConfirmPassword string `json:"confirm_password"`
		OTP             int    `json:"otp"`
	}
)

func (p *PassTokens) IsExpired() bool {
	log.Println("ExpiresIn:\t", p.ExpiresIn)
	log.Println("Now time:\t", time.Now().Unix())
	return p.ExpiresIn < time.Now().Unix()
}
