package service

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/startup-of-zero-reais/API-customer-experience/src/common/providers"

	"github.com/google/uuid"

	d "github.com/startup-of-zero-reais/API-customer-experience/src/common/domain"
	"github.com/startup-of-zero-reais/API-customer-experience/src/common/validation"

	"github.com/startup-of-zero-reais/API-customer-experience/src/session/domain"

	log "github.com/sirupsen/logrus"
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
		logger        *providers.LogProvider
	}
)

func NewRecoverPassword(userRepository domain.UserRepository, otpRepository domain.OTPRepository, sender domain.Sender, evtRepository d.EventRepository, logger *providers.LogProvider) RecoverPassword {
	return &RecoverPasswordImpl{
		userRepository: userRepository,
		otpRepository:  otpRepository,

		sender: sender,

		evtRepository: evtRepository,
		logger:        logger,
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

	bytes, err := json.Marshal(passToken)
	if err != nil {
		return err
	}

	r.logger.WithFields(log.Fields{
		"user_id": user.ID,
		"event":   d.RequestPasswordRecover,
	}).Infoln(string(bytes))

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
		r.logger.Debugln("passToken:", passToken)
		r.logger.Debugln("isExpired:", passToken.IsExpired())

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

	bytes, err := json.Marshal(validToken)
	if err != nil {
		return err
	}

	r.logger.WithFields(log.Fields{
		"email": validToken.Email,
		"event": d.PasswordRecovered,
	}).Infoln(string(bytes))

	return nil
}
