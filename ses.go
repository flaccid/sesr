package sesr

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ses"

	log "github.com/Sirupsen/logrus"
)

func Send(awsRegion string, awsAccessKey string, awsSecretAccessKey string, sender string, recipients []string, subject string, textBody string, charset string) (err error) {
	// create a new aws session
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(awsRegion)},
	)

	// create a new session
	svc := ses.New(sess)

	// assemble the email
	input := &ses.SendEmailInput{
		Destination: &ses.Destination{
			CcAddresses: []*string{},
			ToAddresses: aws.StringSlice(recipients),
		},
		Message: &ses.Message{
			Body: &ses.Body{
				// TODO: support html
				// Html: &ses.Content{
				// 	Charset: aws.String(c.String("charset")),
				// 	Data:    aws.String(c.String("body")),
				// },
				Text: &ses.Content{
					Charset: aws.String(charset),
					Data:    aws.String(textBody),
				},
			},
			Subject: &ses.Content{
				Charset: aws.String(charset),
				Data:    aws.String(subject),
			},
		},
		Source: aws.String(sender),
	}

	// attempt to send the email
	result, err := svc.SendEmail(input)

	log.WithFields(log.Fields{
		"sender":     sender,
		"recipients": fmt.Sprint(recipients),
		"subject":    subject,
		"charset":    charset,
		"result":     result,
		"body":       "[redacted]",
	}).Debug("ses mail send")

	if err != nil {
		return err.(awserr.Error)
	}

	return nil
}
