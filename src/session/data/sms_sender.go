package data

import (
	"context"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sns"
	"github.com/startup-of-zero-reais/API-customer-experience/src/session/domain"
)

type (
	SenderImpl struct {
		ctx context.Context
		sns *sns.Client
	}
)

func NewSender() domain.Sender {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		panic("configuration error, " + err.Error())
	}

	client := sns.NewFromConfig(cfg)

	return &SenderImpl{
		ctx: context.TODO(),
		sns: client,
	}
}

func (s *SenderImpl) SendSMS(message, phone string) error {
	output, err := s.sns.Publish(s.ctx, &sns.PublishInput{
		PhoneNumber: aws.String(phone),
		Message:     aws.String(message),
	})
	if err != nil {
		return err
	}

	log.Printf("\nMessage ID: %s\n\n", *output.MessageId)

	return nil
}
