package model

import (
	"github.com/jinzhu/gorm"
)

type Profile struct {
	gorm.Model
	UserCode  string    `gorm:"size:50;not null;`
	Pwd  string          `gorm:"size:50;not null;`
	NickName string    	 `gorm:"size:50;not null;`
	Person *Person        `gorm:"foreignkey:IdProfile"`
	Roles Roles         `gorm:"many2many:role_set;association_jointable_foreignkey:id_role;jointable_foreignkey:id_profile;"`
}
