package service

import (
	"encoding/json"

	"github.com/startup-of-zero-reais/API-customer-experience/src/common/validation"

	log "github.com/sirupsen/logrus"
	d "github.com/startup-of-zero-reais/API-customer-experience/src/common/domain"
	"github.com/startup-of-zero-reais/API-customer-experience/src/common/providers"
	"github.com/startup-of-zero-reais/API-customer-experience/src/session/domain"
)

type (
	SignOut interface {
		ClearSession(usrId string) error
	}

	SignOutImpl struct {
		userRepository    domain.UserRepository
		sessionRepository domain.SessionRepository

		evtRepository d.EventRepository
		logger        *providers.LogProvider
	}
)

func NewSignOut(userRepository domain.UserRepository, sessionRepository domain.SessionRepository, evtRepository d.EventRepository, logger *providers.LogProvider) SignOut {
	return &SignOutImpl{
		userRepository:    userRepository,
		sessionRepository: sessionRepository,
		evtRepository:     evtRepository,
		logger:            logger,
	}
}

func (s *SignOutImpl) ClearSession(usrId string) error {
	sessions, err := s.sessionRepository.UserSessions(usrId)
	if err != nil {
		return err
	}

	if len(sessions) == 0 {
		return validation.UnauthorizedError("o usuário não está logado")
	}

	for _, session := range sessions {
		s.logger.Debugln("SESSION:", session)
		err = s.sessionRepository.DeleteSession(session.UserID, session.CreatedAt)
		if err != nil {
			return err
		}
	}
	sessionBytes, err := json.Marshal(sessions[0])
	if err != nil {
		return err
	}

	s.logger.WithFields(log.Fields{
		"user_id": sessions[0].UserID,
		"event":   d.SessionEnded,
	}).Info(string(sessionBytes))

	return nil
}
