package data

import (
	"context"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/startup-of-zero-reais/API-customer-experience/src/session/domain"
	domayn "github.com/startup-of-zero-reais/dynamo-for-lambda/domain"
	"github.com/startup-of-zero-reais/dynamo-for-lambda/drivers"
	"github.com/startup-of-zero-reais/dynamo-for-lambda/expressions"
	"github.com/startup-of-zero-reais/dynamo-for-lambda/table"

	s "github.com/startup-of-zero-reais/API-customer-experience/src/common/service"
)

type (
	// SessionRepository is a repository to access the session data
	SessionRepositoryImpl struct {
		dynamo     *drivers.DynamoClient
		jwtService s.JwtService
	}
)

// NewSessionRepository creates a new session repository
func NewSessionRepository(jwtService s.JwtService) domain.SessionRepository {
	dynamo := drivers.NewDynamoClient(
		context.TODO(),
		&domayn.Config{
			TableName: "UserSession",
			Table: table.NewTable(
				"UserSession",
				UserSession{},
			),
			Environment: domayn.Environment(os.Getenv("ENVIRONMENT")),
			Endpoint:    os.Getenv("ENDPOINT"),
		},
	)

	err := dynamo.Migrate()
	if err != nil {
		panic(err)
	}

	return &SessionRepositoryImpl{
		dynamo:     dynamo,
		jwtService: jwtService,
	}
}

// NewSession creates a new session for a user
func (s *SessionRepositoryImpl) NewSession(usrId, email string) (*domain.UserSession, error) {

	token, err := s.jwtService.GenerateToken(usrId, email)
	if err != nil {
		return nil, err
	}

	session := domain.UserSession{
		UserID:       usrId,
		SessionID:    uuid.NewString(),
		CreatedAt:    time.Now().Unix(),
		ExpiresIn:    time.Now().Add(time.Hour * 24).Unix(),
		SessionToken: token,
	}

	sql := s.dynamo.NewExpressionBuilder().SetItem(session)

	err = s.dynamo.Perform(drivers.PUT, sql, &domain.UserSession{})
	if err != nil {
		return nil, err
	}

	return &session, nil
}

func (s *SessionRepositoryImpl) UserSessions(usrId string) ([]domain.UserSession, error) {
	sql := s.dynamo.NewExpressionBuilder().Where(
		expressions.NewKeyCondition("UserID", usrId),
	)

	var sessions []domain.UserSession
	err := s.dynamo.Perform(drivers.QUERY, sql, &sessions)
	if err != nil {
		return nil, err
	}

	if len(sessions) > 0 {
		sessions = sort(sessions)
	}

	return sessions, nil
}

func (s *SessionRepositoryImpl) DeleteSession(userID string, createdAt int64) error {
	sql := s.dynamo.NewExpressionBuilder().Where(
		expressions.NewKeyCondition("UserID", userID),
	).AndWhere(
		expressions.NewSortKeyCondition("CreatedAt").Equal(createdAt),
	)

	err := s.dynamo.Perform(drivers.DELETE, sql, &domain.UserSession{})
	if err != nil {
		return err
	}

	return nil
}
