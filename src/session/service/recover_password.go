package service

import (
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"

	d "github.com/startup-of-zero-reais/API-customer-experience/src/common/domain"
	"github.com/startup-of-zero-reais/API-customer-experience/src/common/validation"

	"github.com/startup-of-zero-reais/API-customer-experience/src/session/domain"
)

type (
	RecoverPassword interface {
		SendOTP(email string) error
		ResetPassword(otp int, password string) error
	}

	RecoverPasswordImpl struct {
		userRepository domain.UserRepository
		otpRepository  domain.OTPRepository

		sender domain.Sender

		evtRepository d.EventRepository
	}
)

func NewRecoverPassword(userRepository domain.UserRepository, otpRepository domain.OTPRepository, sender domain.Sender, evtRepository d.EventRepository) RecoverPassword {
	return &RecoverPasswordImpl{
		userRepository: userRepository,
		otpRepository:  otpRepository,

		sender: sender,

		evtRepository: evtRepository,
	}
}

func (r *RecoverPasswordImpl) SendOTP(email string) error {
	user, err := r.userRepository.Find(email)
	if err != nil {
		return err
	}

	passToken, err := r.otpRepository.New(email)
	if err != nil {
		return err
	}

	err = r.sender.SendSMS(
		fmt.Sprintf("Use o código %d para recuperar sua senha", passToken.OTP),
		user.Phone,
	)
	if err != nil {
		return err
	}

	err = r.evtRepository.Emit(
		user.ID,
		uuid.NewString(),
		d.RequestPasswordRecover,
		passToken,
		time.Now().Unix(),
	)
	if err != nil {
		return err
	}

	return nil
}

func (r *RecoverPasswordImpl) ResetPassword(otp int, password string) error {
	passTokens, err := r.otpRepository.SearchOtp(otp)
	if err != nil {
		return err
	}

	if len(passTokens) == 0 {
		return validation.UnauthorizedError("ERP001: código de recuperação inválido")
	}

	var validToken *domain.PassTokens
	allExpired := true
	for _, passToken := range passTokens {
		log.Printf("\npassToken: %+v\n", passToken)
		log.Printf("\nisExpired: %+v\n", passToken.IsExpired())

		if !passToken.IsExpired() {
			allExpired = false
			validToken = &passToken
			break
		}
	}

	if allExpired {
		return validation.UnauthorizedError("ERP002: código de recuperação inválido")
	}

	if !r.otpRepository.IsValid(validToken.Email, otp) {
		return validation.UnauthorizedError("ERP003: código de recuperação inválido")
	}

	err = r.userRepository.UpdatePassword(validToken.Email, password)
	if err != nil {
		return err
	}

	err = r.otpRepository.Invalidate(validToken.Email)
	if err != nil {
		return err
	}

	err = r.evtRepository.Emit(
		validToken.Email,
		uuid.NewString(),
		d.PasswordRecovered,
		validToken,
		time.Now().Unix(),
	)
	if err != nil {
		return err
	}

	return nil
}
