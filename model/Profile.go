package model

import (
	"github.com/jinzhu/gorm"
)

type Profile struct {
	gorm.Model
	UserCode  string    `gorm:"size:50;not null;`
	Pwd  string          `gorm:"size:50;not null;`
	NickName string    	 `gorm:"size:50;not null;`
}
