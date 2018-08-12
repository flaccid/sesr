package main

import (
	"os"
	"strings"

	"github.com/flaccid/sesr"
	"github.com/urfave/cli"

	log "github.com/Sirupsen/logrus"
)

var (
	VERSION = "v0.0.0-dev"
)

func beforeApp(c *cli.Context) error {
	if c.GlobalBool("debug") {
		log.SetLevel(log.DebugLevel)
	}

	// check if needed aws variables are available
	if len(c.String("aws-region")) < 1 || len(c.String("aws-access-key-id")) < 1 || len(c.String("aws-secret-access-key")) < 1 {
		log.Fatal("insufficient aws credentials")
	}

	if len(c.String("sender")) < 1 && !c.Bool("daemon") {
		log.Fatal("sender email required")
	}

	if len(c.String("recipients")) < 1 && !c.Bool("daemon") {
		log.Fatal("recipient email(s) required")
	}

	if len(c.String("subject")) < 1 && !c.Bool("daemon") {
		log.Fatal("email subject required")
	}

	if c.NArg() < 1 && !c.Bool("daemon") {
		log.Fatal("email body required")
	}

	return nil
}

func main() {
	app := cli.NewApp()
	app.Name = "sesr"
	app.Version = VERSION
	app.Usage = "send email using ses"
	app.Action = start
	app.Before = beforeApp
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "sender,s",
			Usage:  "email sender",
			EnvVar: "EMAIL_SENDER",
		},
		cli.StringFlag{
			Name:   "recipients,r",
			Usage:  "email recipients",
			EnvVar: "EMAIL_RECIPIENTS",
		},
		cli.StringFlag{
			Name:   "subject",
			Usage:  "email subject",
			EnvVar: "EMAIL_SUBJECT",
		},
		cli.StringFlag{
			Name:   "body",
			Usage:  "email body",
			EnvVar: "EMAIL_BODY",
		},
		cli.StringFlag{
			Name:   "charset",
			Usage:  "email character set",
			EnvVar: "EMAIL_CHARSET",
			Value:  "UTF-8",
		},
		cli.StringFlag{
			Name:   "aws-region",
			Usage:  "aws region for the sns service",
			EnvVar: "AWS_REGION",
		},
		cli.StringFlag{
			Name:   "aws-access-key-id",
			Usage:  "aws access key id",
			EnvVar: "AWS_ACCESS_KEY_ID",
		},
		cli.StringFlag{
			Name:   "aws-secret-access-key",
			Usage:  "aws secret access key",
			EnvVar: "AWS_SECRET_ACCESS_KEY",
		},
		cli.BoolFlag{
			Name:  "daemon",
			Usage: "run in daemon mode (web service)",
		},
		// not yet implemented
		cli.BoolFlag{
			Name:  "dry",
			Usage: "run in dry mode",
		},
		cli.BoolFlag{
			Name:  "debug,d",
			Usage: "run in debug mode",
		},
	}
	app.Run(os.Args)
}

func start(c *cli.Context) error {
	if c.Bool("daemon") {
		// pass mandatory cli params to the api service (aws credentials)
		sesr.Serve(c.String("aws-region"), c.String("aws-access-key-id"), c.String("aws-secret-access-key"))
	} else {
		log.WithFields(log.Fields{
			"recipients": c.String("recipients"),
			"message":    c.Args().Get(0),
		}).Info("send email to")

		sesr.Send(c.String("aws-region"), c.String("aws-access-key-id"), c.String("aws-secret-access-key"), c.String("sender"), strings.Split(c.String("recipients"), ","), c.String("subject"), c.Args().Get(0), c.String("charset"))
		//log.Info(err)
	}

	return nil
}
