package domain

import "time"

type (
	// UserSession struct represents a user session
	UserSession struct {
		UserID       string `json:"userId"`
		SessionID    string `json:"sessionId"`
		CreatedAt    int64  `json:"createdAt"`
		ExpiresIn    int64  `json:"expiresIn"`
		SessionToken string `json:"sessionToken"`
	}

	// SessionRepository is a repository to access the session data
	SessionRepository interface {
		NewSession(usrId, email string) (*UserSession, error)
		UserSessions(usrId string) ([]UserSession, error)
		DeleteSession(sessionID string, createdAt int64) error
	}
)

func NewUserSession(usrId, sessId, sessToken string, createdAt, expiresIn int64) *UserSession {
	return &UserSession{
		UserID:       usrId,
		SessionID:    sessId,
		CreatedAt:    createdAt,
		ExpiresIn:    expiresIn,
		SessionToken: sessToken,
	}
}

func (u *UserSession) IsExpired() bool {
	return u.ExpiresIn < time.Now().Unix()
}

func (u *UserSession) Token() string {
	return u.SessionToken
}
