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
		*logrus.Logger
	}
)

func NewLogProvider() *LogProvider {
	log := logrus.New()

	log.SetFormatter(&logrus.JSONFormatter{})
	log.SetOutput(os.Stdout)

	return &LogProvider{
		Logger: log,
	}
}

func (l *LogProvider) LoggerConfig(event events.APIGatewayV2HTTPRequest) {
	http := event.RequestContext.HTTP

	l.WithFields(logrus.Fields{
		"method":    http.Method,
		"path":      http.Path,
		"protocol":  http.Protocol,
		"sourceIp":  http.SourceIP,
		"userAgent": http.UserAgent,
	})

	l.WithTime(time.Now())
}

func (l *LogProvider) LogResponse(res domain.Response) {
	l.WithFields(logrus.Fields{
		"status":  res.StatusCode,
		"headers": res.Headers,
		"cookies": res.Cookies,
	})

	l.Infoln(res.Body.ToJson())
}
