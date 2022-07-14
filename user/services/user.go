package services

import (
	"fmt"
	"os"
)

func SendVerifyEmail(resetToken string) {

	link := os.Getenv("SERVICE_URL") + ":" + os.Getenv("FRONTEND_PORT") + "/verify-user/?verify_token=" + resetToken
	body := "Here is your account verify <a href='" + link + "'>link</a>"
	html := "<strong>" + body + "</strong>"
	fmt.Println(html)
}

func SendPasswordChangeEmail(resetToken string) {

	link := os.Getenv("SERVICE_URL") + ":" + os.Getenv("FRONTEND_PORT") + "/change-password/?verify_token=" + resetToken
	body := "Here is your password reset <a href='" + link + "'>link</a>"
	html := "<strong>" + body + "</strong>"
	fmt.Println(html)
}
