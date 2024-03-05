package utils

import (
	"fmt"
	"math/rand"
	"mime/multipart"
	"net/smtp"
	"os"
	"path"
	"regexp"
	"time"

	"github.com/jordan-wright/email"
)

// 生成随机字符串
func RandomString(l int, Inner string) string {
	var letters = []byte(Inner)
	var result = make([]byte, l)
	rand.NewSource(time.Now().UnixNano())
	for i := range l {
		result[i] = letters[rand.Intn(len(letters))]
	}
	return string(result)
}

// 重命名文件
func Rename(file *multipart.FileHeader, path_file string, ext_allowed map[string]bool) (string, error) {
	var dst string
	oldname := file.Filename
	ext := path.Ext(oldname)
	if !ext_allowed[ext] {
		return "", fmt.Errorf("error extention name") // 文件扩展名不允许
	}
	for {
		times := 0
		dst = "./resource/template/" + fmt.Sprintf("%v(%v)%v", oldname, times, ext)
		_, err := os.Stat(dst)
		if err != nil {
			times++
			continue
		}
		break
	}
	return dst, nil
}

func SendMail(Toemail, smtpUser, smtpPassword, title string) (string, error) {
	smtpHost := "smtp.qq.com"             // SMTP服务器地址
	smtpPort := "587"                     // SMTP服务器端口
	toUserEmail := Toemail                // 接收者邮箱地址
	code := RandomString(6, "0123456789") // 验证码
	if !IsEmailLegal(toUserEmail) {
		return "", fmt.Errorf("邮箱格式错误")
	}
	e := email.NewEmail()
	e.From = smtpUser                                                                                                                                                 // 发件人邮箱账号
	e.To = append(e.To, toUserEmail)                                                                                                                                  // 收件人邮箱地址                                                                               // 收件人邮箱地址
	e.Subject = title                                                                                                                                                 // 邮件主题
	e.Text = []byte("验证码:" + code)                                                                                                                                    // 邮件正文内容（纯文本）
	e.HTML = []byte("<strong>" + string(e.Text) + "</strong><br><p>有效时长5分钟</p><br><br><br><p>  本项目由mahaonan001在GitHub上开源的问卷系统go项目,如果有兴趣参加,欢迎联系1649801526@qq.com</p>") // 邮件正文内容（HTML格式）
	err := e.Send(smtpHost+":"+smtpPort, smtp.PlainAuth("", smtpUser, smtpPassword, smtpHost))                                                                        // 发送邮件
	if err != nil {
		return "", err
	}
	return code, nil

}
func IsEmailLegal(email string) bool {
	var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,4}$`)
	// 使用MatchString()函数来判断电子邮件地址是否匹配正则表达式
	return emailRegex.MatchString(email)
}
