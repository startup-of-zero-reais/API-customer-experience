package data

import (
	"context"
	"log"
	"math/rand"
	"time"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/startup-of-zero-reais/dynamo-for-lambda/expressions"

	"github.com/startup-of-zero-reais/API-customer-experience/src/common/validation"
	"github.com/startup-of-zero-reais/API-customer-experience/src/session/domain"

	domayn "github.com/startup-of-zero-reais/dynamo-for-lambda/domain"
	"github.com/startup-of-zero-reais/dynamo-for-lambda/drivers"
	"github.com/startup-of-zero-reais/dynamo-for-lambda/table"
)

type (
	OTPRepositoryImpl struct {
		dynamo *drivers.DynamoClient
	}
)

func NewOTPRepository() domain.OTPRepository {
	var awsConfs []func(*config.LoadOptions) error
	awsConf := append(awsConfs, config.WithRegion("us-east-1"))

	cfg, err := config.LoadDefaultConfig(context.TODO(), awsConf...)
	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
	}

	dynamo := drivers.NewDynamoClient(
		context.TODO(),
		&domayn.Config{
			TableName: "PassTokens",
			Table: table.NewTable(
				"PassTokens",
				PassTokens{},
			),
			Client: dynamodb.NewFromConfig(cfg),
		},
	)

	return &OTPRepositoryImpl{
		dynamo: dynamo,
	}
}

func (o *OTPRepositoryImpl) getOtp(email string, otp int) (*domain.PassTokens, error) {
	sql := o.dynamo.NewExpressionBuilder().Where(
		expressions.NewKeyCondition("Email", email),
	).AndWhere(
		expressions.NewSortKeyCondition("OTP").Equal(otp),
	)

	var passToken domain.PassTokens
	err := o.dynamo.Perform(drivers.GET, sql, &passToken)
	if err != nil {
		return nil, err
	}

	return &passToken, nil
}

func (o *OTPRepositoryImpl) getOtps(email string) ([]domain.PassTokens, error) {
	sql := o.dynamo.NewExpressionBuilder().Where(
		expressions.NewKeyCondition("Email", email),
	)

	var passToken []domain.PassTokens
	err := o.dynamo.Perform(drivers.QUERY, sql, &passToken)
	if err != nil {
		return nil, err
	}

	return passToken, nil
}

func (o *OTPRepositoryImpl) deleteOtp(email string, otp int) error {
	sql := o.dynamo.NewExpressionBuilder().Where(
		expressions.NewKeyCondition("Email", email),
	).AndWhere(
		expressions.NewSortKeyCondition("OTP").Equal(otp),
	)

	err := o.dynamo.Perform(drivers.DELETE, sql, &PassTokens{})
	if err != nil {
		return err
	}

	return nil
}

func (o *OTPRepositoryImpl) New(email string) (*domain.PassTokens, error) {
	sql := o.dynamo.NewExpressionBuilder()

	rand.Seed(time.Now().Unix())
	otp := rand.Intn(999999)

	passToken := &domain.PassTokens{
		Email:     email,
		ExpiresIn: time.Now().Add(time.Minute * 5).Unix(),
		OTP:       otp,
	}
	sql.SetItem(*passToken)

	err := o.dynamo.Perform(drivers.PUT, sql, &domain.PassTokens{})
	if err != nil {
		return nil, err
	}

	return passToken, nil
}

func (o *OTPRepositoryImpl) Invalidate(email string) error {
	otps, err := o.getOtps(email)
	if err != nil {
		return err
	}

	if len(otps) == 0 {
		return validation.NotFoundError("EDO001: código de recuperação inválido")
	}

	for _, otp := range otps {
		err = o.deleteOtp(email, otp.OTP)
		if err != nil {
			return err
		}
	}

	return nil
}

func (o *OTPRepositoryImpl) IsValid(email string, OTP int) bool {
	otp, err := o.getOtp(email, OTP)
	if err != nil {
		log.Println("err", err.Error())
		return false
	}

	if otp.IsExpired() {
		log.Println("\n\tIsExpired", otp.IsExpired())
		return false
	}

	return true
}

func (o *OTPRepositoryImpl) SearchOtp(otp int) ([]domain.PassTokens, error) {
	sql := o.dynamo.NewExpressionBuilder().Where(
		expressions.NewKeyCondition("OTP", otp),
	).SetIndex("OtpIndex")

	var passTokens []domain.PassTokens
	err := o.dynamo.Perform(drivers.QUERY, sql, &passTokens)
	if err != nil {
		return nil, err
	}

	if len(passTokens) == 0 {
		return nil, validation.NotFoundError("EDO002: código de recuperação inválido")
	}

	return passTokens, nil
}
