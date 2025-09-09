package utils

import (
	"fmt"

	"gopkg.in/gomail.v2"
)

func SendOtp(to string, otp string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", "tejeswar.raju357@gmail.com")
	m.SetHeader("To", to)
	m.SetHeader("Subject", "Your Login OTP")
	m.SetBody("text/plain", fmt.Sprintf("Your OTP is %s", otp))

	d := gomail.NewDialer("smtp.gmail.com", 587, "tejeswar.raju357@gmail.com", "xcvk phwz hqxb karr")
	err := d.DialAndSend(m)
	if err != nil {
		fmt.Println("SendOTP error", err)
	}
	return err
}
