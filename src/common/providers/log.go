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
		log.SetFormatter(&logrus.TextFormatter{
			ForceColors: true,
		})
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
	fields := logrus.Fields{
		"status":  res.StatusCode,
		"headers": res.Headers,
	}

	if len(res.Cookies) > 0 {
		fields["cookies"] = res.Cookies
	}

	l.Entry = l.Entry.WithFields(fields)

	switch res.StatusCode {
	case 200, 201, 202, 204:
		l.Infoln(res.Body.ToJson())
	case 300, 301, 302, 303, 304, 305, 306, 307:
		l.Warnln(res.Body.ToJson())
	default:
		l.Errorln(res.Body.ToJson())
	}
}
