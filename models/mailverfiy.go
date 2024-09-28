package models

import "gorm.io/gorm"

type MailVerify struct {
	gorm.Model
	Email        string `json:"email"`
	Mailauthcode string `json:"mailauthcode"`
	Verifycode   string `json:"verifycode"`
}
