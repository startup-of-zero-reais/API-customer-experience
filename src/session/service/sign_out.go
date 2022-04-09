package service

import (
	"errors"
	"log"

	"github.com/startup-of-zero-reais/API-customer-experience/src/session/domain"
)

type (
	SignOut interface {
		ClearSession(usrId string) error
	}

	SignOutImpl struct {
		userRepository    domain.UserRepository
		sessionRepository domain.SessionRepository
	}
)

func NewSignOut(userRepository domain.UserRepository, sessionRepository domain.SessionRepository) SignOut {
	return &SignOutImpl{
		userRepository:    userRepository,
		sessionRepository: sessionRepository,
	}
}

func (s *SignOutImpl) ClearSession(usrId string) error {
	sessions, err := s.sessionRepository.UserSessions(usrId)
	if err != nil {
		return err
	}

	if len(sessions) == 0 {
		return errors.New("o usuário não está logado")
	}

	for _, session := range sessions {
		log.Printf("\nSESSION: %+v\n\n", session)
		err = s.sessionRepository.DeleteSession(session.UserID, session.CreatedAt)
		if err != nil {
			return err
		}
	}

	return nil
}
