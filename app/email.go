package app

import (
	"bytes"
	"fmt"
	"html/template"
	"log"

	"os"

	"github.com/sendgrid/sendgrid-go"
)

var tmpls *template.Template

// var sender = sendgrid.NewSendGridClient(config.Config.Email.User, config.Config.Email.Key)
var FromEmail = "testing@theplant.jp"
var FromName = "Auto Testing"

var emailAccount = ""
var emailKey = ""

func init() {
	var err error
	gopath := os.Getenv("GOPATH")
	if tmpls, err = template.New("emails").ParseGlob(gopath + "/src/github.com/azumads/selenium/app/*.tmpl"); err != nil {
		panic(err)
	}
}

func NewHtmlEmail(to, toName, subject, content string) (email *sendgrid.SGMail) {
	email = sendgrid.NewMail()
	email.SetFrom(FromEmail)
	email.SetFromName(FromName)
	email.AddTo(to)
	email.AddToName(toName)
	email.SetSubject(subject)
	email.SetHTML(content)
	return
}

func Send(email *sendgrid.SGMail, emailType string) (err error) {
	sender := sendgrid.NewSendGridClient(emailAccount, emailKey)
	if err = sender.Send(email); err != nil {
		log.Printf("send %s email to %s, error : %s", emailType, email.To, err.Error())
	} else {
		log.Printf("send %s email to %s, success", emailType, email.To)
	}
	return
}

type EmailData struct {
	Message string
	Link    string
}

func SendNotifyErrorEmail(to, project, testcase, jobId string) (err error) {
	data := EmailData{
		Message: "Project: " + project + "    TestCase: " + testcase,
		Link:    HostUrl() + fmt.Sprintf("/admin/run_tests/%s", jobId),
	}
	var buf bytes.Buffer
	var content string
	if err = tmpls.ExecuteTemplate(&buf, "notify_error_email", data); err != nil {
		return
	}
	content = buf.String()
	email := NewHtmlEmail(to, "", "Test fail", content)
	err = Send(email, "register email confirm")
	return
}
