package service

import (
	"time"

	"github.com/google/uuid"
	"github.com/startup-of-zero-reais/API-customer-experience/src/common/validation"

	"github.com/startup-of-zero-reais/API-customer-experience/src/session/domain"

	d "github.com/startup-of-zero-reais/API-customer-experience/src/common/domain"
)

type (
	SignIn interface {
		SignIn(email, password string) (*domain.UserSession, error)
	}

	SignInImpl struct {
		userRepository    domain.UserRepository
		sessionRepository domain.SessionRepository

		evtRepository d.EventRepository
	}
)

func NewSignIn(userRepository domain.UserRepository, sessionRepository domain.SessionRepository, evtRepository d.EventRepository) SignIn {
	return &SignInImpl{
		userRepository:    userRepository,
		sessionRepository: sessionRepository,

		evtRepository: evtRepository,
	}
}

func (s *SignInImpl) SignIn(email, password string) (*domain.UserSession, error) {
	user, err := s.userRepository.Find(email)
	if err != nil {
		return nil, err
	}

	passwordChecker := domain.NewPassword(user.Password)
	err = passwordChecker.Validate(password)
	if err != nil {
		return nil, err
	}

	sessions, err := s.sessionRepository.UserSessions(user.ID)
	if err != nil {
		return nil, err
	}

	if len(sessions) > 0 {
		lastSession := sessions[0]
		if lastSession.IsExpired() {
			return nil, validation.BadRequestError("sessão expirada")
		} else {
			return &lastSession, nil
		}
	}

	session, err := s.sessionRepository.NewSession(user.ID, user.Email)
	if err != nil {
		return nil, err
	}

	err = s.evtRepository.Emit(
		user.ID,
		uuid.NewString(),
		d.SessionStarted,
		session,
		time.Now().Unix(),
	)
	if err != nil {
		return nil, err
	}

	return session, nil
}
