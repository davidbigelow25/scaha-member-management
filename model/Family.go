package model

import (
	"github.com/jinzhu/gorm"
)

type Family struct {
	gorm.Model
	FamilyName  string    `gorm:"size:75;not null;`
	IdPerson int
	Person Person		 `gorm:"foreignkey:IdPerson;"`
	FamilyMembers []FamilyMember  `gorm:"foreignkey:IdFamily;"`
}
