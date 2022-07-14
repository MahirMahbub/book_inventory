package services

import (
	"fmt"
	"os"
)

func SendVerifyEmail(resetToken string) {

	link := os.Getenv("SERVICE_URL") + ":" + os.Getenv("USER_SERVICE_PORT") + os.Getenv("VERSION_GROUP") + "/user/verify/?verify_token=" + resetToken
	body := "Here is your reset <a href='" + link + "'>link</a>"
	html := "<strong>" + body + "</strong>"
	fmt.Println(html)
}
