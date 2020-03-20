package model

import "github.com/jinzhu/gorm"

type FamilyMember struct {
	gorm.Model
	RelationType string   `gorm:"size:45;not null;`
	IdPerson int
	Person Person        `gorm:"foreignkey:IdPerson;"`
	IdFamily int
	Family *Family
}

