package utils

import (
	"fmt"
	"net/smtp"
	"net/url"
	"posServer/internal/config"
	"strconv"
)

var (
	conf         = config.NewConfig()
	configMailer = conf.SMTP
)

func parseSmtpUrl(smtpURL string) (host string, port int, username, password, sender string, useTLS bool, err error) {
	u, err := url.Parse(smtpURL)
	if err != nil {
		return
	}

	host = u.Hostname()
	port, _ = strconv.Atoi(u.Port())
	username = u.User.Username()
	password, _ = u.User.Password()

	q := u.Query()
	useTLS, _ = strconv.ParseBool(q.Get("useTLS"))
	sender = q.Get("sender")

	return
}

func SendOTP(otp, mailAddr string) error {
	host, port, username, password, sender, useTLS, err := parseSmtpUrl(configMailer)
	if err != nil {
		return fmt.Errorf("Error parsing SMTP url: %v", err)
	}

	addr := fmt.Sprintf("%s:%d", host, port)
	subject := fmt.Sprintf("SQL POS Terminal Binding Verification Code\n")
	body := fmt.Sprintf(`
		<html>
			<body>
				<p>Hi!<p/>
				<p>You are completing terminal binding verification. Your verification code is: <strong style="font-size:18px;color:blue;">%s</strong></p>
				<p>Please complete this process within 5 minutes.</p>
				<p>SQL POS</p>
				<br>
				<p style="font-size:12px;color:grey;">This is an automated email. Please do not reply to this email.</p>
			</body>
		</html>	
	`, otp)

	msg := "From: " + sender + "\n" +
		"To: " + mailAddr + "\n" +
		"Subject: " + subject +
		"MIME-version: 1.0;\n" +
		"Content-Type: text/html; charset=\"UTF-8\";\n\n" +
		body

	auth := smtp.PlainAuth("", username, password, host)

	if useTLS {
		err = smtp.SendMail(addr, auth, sender, []string{mailAddr}, []byte(msg))
	} else {
		err = smtp.SendMail(addr, nil, sender, []string{mailAddr}, []byte(msg))
	}
	if err != nil {
		return fmt.Errorf("SMTP error: %s\n", err)
	}
	return nil
}
