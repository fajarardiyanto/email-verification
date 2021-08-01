package service

import (
	"github.com/d-vignesh/go-jwt-auth/utils"
	"github.com/hashicorp/go-hclog"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

type MailService interface {
	CreateMail(mr *Mail) []byte
	SendMail(mr *Mail) error
	NewMail(from string, to []string, subject string, mailType MailType, data *MailData) *Mail
}

type MailType int

const (
	MailConfirmation MailType = iota + 1
	PassReset
)

type MailData struct {
	Username string
	Code     string
}

type Mail struct {
	from    string
	to      []string
	subject string
	body    string
	mtype   MailType
	data    *MailData
}

type SGMailService struct {
	logger  hclog.Logger
	configs *utils.Configurations
}

func NewSGMailService(logger hclog.Logger, configs *utils.Configurations) *SGMailService {
	return &SGMailService{logger, configs}
}

func (ms *SGMailService) CreateMail(mr *Mail) []byte {

	m := mail.NewV3Mail()

	from := mail.NewEmail("bookite", mr.from)
	m.SetFrom(from)

	if mr.mtype == MailConfirmation {
		m.SetTemplateID(ms.configs.MailVerifTemplateID)
	} else if mr.mtype == PassReset {
		m.SetTemplateID(ms.configs.PassResetTemplateID)
	}

	p := mail.NewPersonalization()

	tos := make([]*mail.Email, 0)
	for _, to := range mr.to {
		tos = append(tos, mail.NewEmail("user", to))
	}

	p.AddTos(tos...)

	p.SetDynamicTemplateData("Username", mr.data.Username)
	p.SetDynamicTemplateData("Code", mr.data.Code)

	m.AddPersonalizations(p)
	return mail.GetRequestBody(m)
}

func (ms *SGMailService) SendMail(mr *Mail) error {

	request := sendgrid.GetRequest(ms.configs.SendGridApiKey, "/v3/mail/send", "https://api.sendgrid.com")
	request.Method = "POST"
	var Body = ms.CreateMail(mr)
	request.Body = Body
	response, err := sendgrid.API(request)
	if err != nil {
		ms.logger.Error("unable to send mail", "error", err)
		return err
	}
	ms.logger.Info("mail sent successfully", "sent status code", response.StatusCode)
	return nil
}

func (ms *SGMailService) NewMail(from string, to []string, subject string, mailType MailType, data *MailData) *Mail {
	return &Mail{
		from:    from,
		to:      to,
		subject: subject,
		mtype:   mailType,
		data:    data,
	}
}
