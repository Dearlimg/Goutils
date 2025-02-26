package email

import (
	"crypto/tls"
	"gopkg.in/gomail.v2"
)

/*
Gomail 是一个简单高效的发送电子邮件 go 语言的模块包，其使用 SMTP 服务器发送电子邮件含附件。
目前只支持使用 SMTP 服务器发送电子邮件，但是其 API 较为灵活，如果有其它的定制需求也可以轻易地借助其实现，这恰恰好符合我们的需求，
因为目前我们只需要一个小而美的发送电子邮件的库就可以了。
*/

// SMTPInfo 发送邮箱所必须的信息
type SMTPInfo struct {
	Port     int
	IsSSL    bool
	Host     string
	UserName string
	Password string
	From     string
}

type Email struct {
	*SMTPInfo
}

func NewEmail(info *SMTPInfo) *Email {
	return &Email{SMTPInfo: info}
}

func (e *Email) SendMail(to []string, subject, body string) error {
	//1.首先构建一个 Message 对象，也就是邮件对象
	msg := gomail.NewMessage()
	//2.填充 From，注意第一个字母要大写
	msg.SetHeader("From", e.From)
	//3.填充 To
	msg.SetHeader("To", to...)
	//4.设置邮件主题
	msg.SetHeader("Subject", subject)
	//5.设置要发送邮件的正文
	//第一个参数是类型，第二个参数是内容
	//如果是 html，第一个参数是 `text/html`，如果是文本则是 `text/plain`
	msg.SetBody("text/html", body)

	//6.创建 SMTP 实例
	//SMTP（Simple Mail Transfer Protocol）是一种提供可靠且有效的电子邮件传输的协议。
	//简单来说，SMTP 可以用来发送电子邮件，将电子邮件从发送方传递到接受方的邮箱服务器。这样收件人就可以收到发件人发送的邮箱了
	dialer := gomail.NewDialer(e.Host, e.Port, e.UserName, e.Password)
	//7.设置 TLS 设置，InsecureSkipVerify 字段用于控制是否跳过对服务器证书的验证。
	//如果 e.IsSSL 为 true，则表示跳过对服务器证书的验证如果为false，则表示不跳过验证，即需要验证服务器证书的有效性。跳过验证可能会存在安全风险，因为无法确保连接的目标服务器的身份可信。
	//建议在生产环境中谨慎使用InsecureSkipVerify，最好保持默认值为false，以确保安全性。如果需要在开发或测试环境中使用自签名证书等情况，才考虑设置为true跳过证书验证。
	dialer.TLSConfig = &tls.Config{InsecureSkipVerify: e.IsSSL}

	//8.发送邮件，打开与 SMTP 服务器的连接并发送电子邮件，发送完就关闭
	return dialer.DialAndSend(msg)
}
