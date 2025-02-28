package email

import (
	"fmt"
	"github.com/Dearlimg/Goutils/pkg/times"
	"testing"
	"time"
)

var (
	Host     = "smtp.qq.com"
	Port     = 465
	UserName = "xyy@qq.com"
	Password = ""
	IsSSL    = true
	From     = "xyy@qq.com"
	To       = []string{"xyy@qq.com"}
)

func TestEmail_SendMail(t *testing.T) {
	defaultMailer := NewEmail(&SMTPInfo{
		Port:     Port,
		IsSSL:    IsSSL,
		Host:     Host,
		UserName: UserName,
		Password: Password,
		From:     From,
	})
	err := defaultMailer.SendMail( //短信通知
		To,
		fmt.Sprintf("异常抛出，发生时间：%s,%d", times.GetNowDateTimeStr(), time.Now().Unix()),
		fmt.Sprintf("错误信息：%v", "NO"),
	)
	if err != nil {
		fmt.Println(err)
	}
}
