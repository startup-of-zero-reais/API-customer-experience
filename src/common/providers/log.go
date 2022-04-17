package providers

import (
	"os"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/sirupsen/logrus"
	"github.com/startup-of-zero-reais/API-customer-experience/src/common/domain"
)

type (
	LogProvider struct {
		client *logrus.Logger
		*logrus.Entry
	}
)

func NewLogProvider() *LogProvider {
	log := logrus.New()

	if os.Getenv("ENVIRONMENT") != "development" {
		log.SetFormatter(&logrus.JSONFormatter{})
	} else {
		log.SetFormatter(&logrus.TextFormatter{})
	}
	log.SetOutput(os.Stdout)

	l := &LogProvider{
		client: log,
	}

	l.Entry = l.client.WithTime(time.Now())

	return l
}

func (l *LogProvider) LoggerConfig(event events.APIGatewayV2HTTPRequest) {
	http := event.RequestContext.HTTP

	l.Entry = l.WithFields(logrus.Fields{
		"method":    http.Method,
		"path":      http.Path,
		"protocol":  http.Protocol,
		"sourceIp":  http.SourceIP,
		"userAgent": http.UserAgent,
	})

}

func (l *LogProvider) LogResponse(res domain.Response) {
	l.Entry = l.Entry.WithFields(logrus.Fields{
		"status":  res.StatusCode,
		"headers": res.Headers,
		"cookies": res.Cookies,
	})

	l.Infoln(res.Body.ToJson())
}
