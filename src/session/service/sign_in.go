package service

import (
	"encoding/json"

	log "github.com/sirupsen/logrus"

	"github.com/startup-of-zero-reais/API-customer-experience/src/session/domain"

	d "github.com/startup-of-zero-reais/API-customer-experience/src/common/domain"
	"github.com/startup-of-zero-reais/API-customer-experience/src/common/providers"
)

type (
	SignIn interface {
		SignIn(email, password string) (*domain.UserSession, error)
	}

	SignInImpl struct {
		userRepository    domain.UserRepository
		sessionRepository domain.SessionRepository

		evtRepository d.EventRepository
		logger        *providers.LogProvider
	}
)

func NewSignIn(userRepository domain.UserRepository, sessionRepository domain.SessionRepository, evtRepository d.EventRepository, logger *providers.LogProvider) SignIn {
	return &SignInImpl{
		userRepository:    userRepository,
		sessionRepository: sessionRepository,

		evtRepository: evtRepository,
		logger:        logger,
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
		if !lastSession.IsExpired() {
			err = s.logEvent(&lastSession, true)
			if err != nil {
				return nil, err
			}
			return &lastSession, nil
		}
	}

	session, err := s.sessionRepository.NewSession(user.ID, user.Email)
	if err != nil {
		return nil, err
	}

	err = s.logEvent(session, false)
	if err != nil {
		return nil, err
	}

	return session, nil
}

func (s *SignInImpl) logEvent(session *domain.UserSession, recovered_session bool) error {
	sessionBytes, err := json.Marshal(session)
	if err != nil {
		return err
	}

	s.logger.WithFields(log.Fields{
		"user_id":   session.UserID,
		"event":     d.SessionStarted,
		"recovered": recovered_session,
	}).Info(string(sessionBytes))

	return nil
}
