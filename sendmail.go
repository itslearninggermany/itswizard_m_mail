package itswizard_m_mail

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"github.com/jinzhu/gorm"
	"log"
	"net"
	"net/mail"
	"net/smtp"
	"strconv"
	"text/template"
)

// delete in itswizard_aws
type DbEmailServerData15 struct {
	gorm.Model
	SmtpServer string
	Port       uint
	Password   string
	Username   string
}

func (p *DbEmailServerData15) SendMailWithTemplate(templatefile string, inputTemplateStruct interface{}, receiverMail string, subject string) error {
	from := mail.Address{"", "donotreply@itswizard.de"}
	to := mail.Address{"", receiverMail}
	//	subj := mailsubject
	// Setup headers
	headers := make(map[string]string)
	headers["From"] = from.String()
	headers["To"] = to.String()
	//	headers["Subject"] = subj

	t, _ := template.ParseFiles(templatefile)

	var bo bytes.Buffer

	mimeHeaders := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	bo.Write([]byte(fmt.Sprintf("Subject:"+subject+" \n%s\n\n", mimeHeaders)))

	err := t.Execute(&bo, inputTemplateStruct)
	if err != nil {
		return err
	}

	port := strconv.Itoa(int(p.Port))
	servername := p.SmtpServer + ":" + port
	host, _, _ := net.SplitHostPort(servername)
	fmt.Println(servername)
	auth := smtp.PlainAuth("", p.Username, p.Password, host)

	// TLS config
	tlsconfig := &tls.Config{
		InsecureSkipVerify: true,
		ServerName:         host,
	}

	// Here is the key, you need to call tls.Dial instead of smtp.Dial
	// for smtp servers running on 465 that require an ssl connection
	// from the very beginning (no starttls)
	conn, err := tls.Dial("tcp", servername, tlsconfig)
	if err != nil {
		return err
	}

	c, err := smtp.NewClient(conn, host)
	if err != nil {
		return err
	}

	// Auth
	if err = c.Auth(auth); err != nil {
		return err
	}

	// To && From
	if err = c.Mail(from.Address); err != nil {
		return err
	}

	if err = c.Rcpt(to.Address); err != nil {
		return err
	}

	// Inst
	w, err := c.Data()
	if err != nil {
		return err
	}

	_, err = w.Write(bo.Bytes())
	if err != nil {
		return err
	}

	err = w.Close()
	if err != nil {
		return err
	}

	err = c.Quit()
	log.Println(err)
	return nil
}

// Email
type EmailForCredentials struct {
	gorm.Model
	Name               string `gorm:"unique"`
	Subjekt            string `gorm:"type:TEXT"`
	Logo1              string `gorm:"type:TEXT"`
	Logo2              string `gorm:"type:TEXT"`
	Preheader          string `gorm:"type:TEXT"`
	Welcome            string `gorm:"type:TEXT"`
	First              string `gorm:"type:TEXT"`
	Url                string `gorm:"type:TEXT"`
	Second             string `gorm:"type:TEXT"`
	Greetings          string `gorm:"type:TEXT"`
	Sender             string `gorm:"type:TEXT"`
	ClientDisclaimer   string `gorm:"type:TEXT"`
	ClientInstitution1 string `gorm:"type:TEXT"`
	ClientInstitution2 string `gorm:"type:TEXT"`
	ClientStreet       string `gorm:"type:TEXT"`
	ClientLocal        string `gorm:"type:TEXT"`
	Foot               string `gorm:"type:TEXT"`
}

func (p *EmailForCredentials) SendCredentials(FirstName, LastName, UserName, Password, email string, dbWebserver *gorm.DB, admin bool) error {
	var emailSetup DbEmailServerData15
	err := dbWebserver.Last(&emailSetup).Error
	if err != nil {
		return err
	}

	if admin {
		return emailSetup.SendMailWithTemplate("/home/ubuntu/brooker/emailtemplate/bw_admin.html", struct {
			Logo1              string
			Logo2              string
			Preheader          string
			Welcome            string
			FirstName          string
			LastName           string
			First              string
			Username           string
			Password           string
			Url                string
			Second             string
			Greetings          string
			Sender             string
			ClientDisclaimer   string
			ClientInstitution1 string
			ClientInstitution2 string
			ClientStreet       string
			ClientLocal        string
			Foot               string
		}{
			Logo1:              p.Logo1,
			Logo2:              p.Logo2,
			Preheader:          p.Preheader,
			Welcome:            p.Welcome,
			FirstName:          FirstName,
			LastName:           LastName,
			First:              p.First,
			Username:           UserName,
			Password:           Password,
			Url:                p.Url,
			Second:             p.Second,
			Greetings:          p.Greetings,
			Sender:             p.Sender,
			ClientDisclaimer:   p.ClientDisclaimer,
			ClientInstitution1: p.ClientInstitution1,
			ClientInstitution2: p.ClientInstitution2,
			ClientStreet:       p.ClientStreet,
			ClientLocal:        p.ClientLocal,
			Foot:               p.Foot,
		}, email, p.Subjekt)
	} else {
		return emailSetup.SendMailWithTemplate("/home/ubuntu/brooker/emailtemplate/bw_user_cred.html", struct {
			Logo1              string
			Logo2              string
			Preheader          string
			Welcome            string
			FirstName          string
			LastName           string
			First              string
			Username           string
			Password           string
			Url                string
			Second             string
			Greetings          string
			Sender             string
			ClientDisclaimer   string
			ClientInstitution1 string
			ClientInstitution2 string
			ClientStreet       string
			ClientLocal        string
			Foot               string
		}{
			Logo1:              p.Logo1,
			Logo2:              p.Logo2,
			Preheader:          p.Preheader,
			Welcome:            p.Welcome,
			FirstName:          FirstName,
			LastName:           LastName,
			First:              p.First,
			Username:           UserName,
			Password:           Password,
			Url:                p.Url,
			Second:             p.Second,
			Greetings:          p.Greetings,
			Sender:             p.Sender,
			ClientDisclaimer:   p.ClientDisclaimer,
			ClientInstitution1: p.ClientInstitution1,
			ClientInstitution2: p.ClientInstitution2,
			ClientStreet:       p.ClientStreet,
			ClientLocal:        p.ClientLocal,
			Foot:               p.Foot,
		}, email, p.Subjekt)
	}
}
