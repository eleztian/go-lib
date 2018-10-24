package email

import (
	"bytes"
	"fmt"
	"html/template"
	"path/filepath"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ses"
	"go.uber.org/zap"

	"fangcun.vesync.com/vdmp/DeveloperPlatform/pkg/log"
)

type Email struct {
	To      string `json:"to"`
	Subject string `json:"subject"`
	Msg     string `json:"msg"`
}

var (
	Sender       string
	emailSession *session.Session
)

var TemplateRootPath = "template"

func Init(region string, sender string, accessId string, accessKey string, templateRootPath string) {
	emailSession = session.Must(session.NewSession(&aws.Config{
		Region:                        aws.String(region),
		CredentialsChainVerboseErrors: aws.Bool(true),
		Credentials:                   credentials.NewStaticCredentials(accessId, accessKey, ""),
	}))
	Sender = sender
	if templateRootPath != "" {
		TemplateRootPath = templateRootPath
	}
}

func sendEmail(to string, topic string, content string) error {
	svc := ses.New(emailSession)

	input := &ses.SendEmailInput{
		Destination: &ses.Destination{
			ToAddresses: []*string{
				aws.String(to),
			},
		},
		Message: &ses.Message{
			Body: &ses.Body{
				Html: &ses.Content{
					Charset: aws.String("UTF-8"),
					Data:    aws.String(content),
				},
				Text: &ses.Content{
					Charset: aws.String("UTF-8"),
					Data:    aws.String(content),
				},
			},
			Subject: &ses.Content{
				Charset: aws.String("UTF-8"),
				Data:    aws.String(topic),
			},
		},
		Source: aws.String(Sender),
		//SourceArn: aws.String("arn:aws:ses:us-east-1:704832300401:identity/vesync.com"),
	}

	result, err := svc.SendEmail(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case ses.ErrCodeMessageRejected:
				return fmt.Errorf("%s: %s", ses.ErrCodeMessageRejected, aerr.Error())
			case ses.ErrCodeMailFromDomainNotVerifiedException:
				return fmt.Errorf("%s: %s", ses.ErrCodeMailFromDomainNotVerifiedException, aerr.Error())
			case ses.ErrCodeConfigurationSetDoesNotExistException:
				return fmt.Errorf("%s: %s", ses.ErrCodeConfigurationSetDoesNotExistException, aerr.Error())
			default:
				return fmt.Errorf("%s", aerr.Error())
			}
		}
		return fmt.Errorf("%s", err.Error())
	}
	log.Debug("send email", zap.String("messageId", result.String()))
	return nil
}

const (
	BaseFile = "base.html"
)

func SendEmailWithTemplateFile(to string, topic string, fileTemplatePath string, data interface{}) {
	if to == "" {
		return
	}
	fileTemplatePath = filepath.Join(TemplateRootPath, fileTemplatePath)
	bPath := filepath.Join(TemplateRootPath, BaseFile)
	t, err := template.ParseFiles(bPath, fileTemplatePath)
	if err != nil {
		log.Error("send email parse file", err, zap.String("to", to), zap.String("topic", topic))
		return
	}
	buff := bytes.NewBuffer([]byte(""))
	err = t.Execute(buff, data)
	if err != nil {
		log.Error("send email execute file", err, zap.String("to", to), zap.String("topic", topic))
		return
	}
	err = sendEmail(to, topic, buff.String())
	if err != nil {
		log.Error("send email", err, zap.String("to", to), zap.String("topic", topic))
		return
	}
}
func SendEmailWithTemplate(to string, topic string, templateContent string, template2 string, data interface{}) error {
	t, err := template.New("test").Parse(templateContent)
	if err != nil {
		return err
	}
	t, err = t.Parse(template2)
	if err != nil {
		return err
	}
	buff := bytes.NewBuffer([]byte(""))
	err = t.Execute(buff, data)
	if err != nil {
		return err
	}
	return sendEmail(to, topic, buff.String())
}
