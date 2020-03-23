package model

import (
	"github.com/jinzhu/gorm"
	"time"
)

type UsaHockey struct {
	gorm.Model
	UsaHockeyNumber string `gorm:"size:14;not null;`
	UsaYear int
	LastName  string     `gorm:"size:45;not null;`
	FirstName  string    `gorm:"size:45;not null;`
	MiddleInitial  string    `gorm:"size:1;not null;`
	Email string    	 `gorm:"size:45;not null;`
	HomePhone  string  	 `gorm:"size:14"`
	WorkPhone  string  	 `gorm:"size:14"`
	Adress string 		 `gorm:"size:75"`
	City string 		 `gorm:"size:45"`
	State string 		 `gorm:"size:2"`
	Zipcode string       `gorm:"size:10"`
	Country string  	 `gorm:"size:45"`
	Gender string    	 `gorm:"size:1"`
	Dob time.Time
	Forzip	string		  `gorm:"size:16"`
	Citizenship string 	  `gorm:"size:45"`
	SeasonTag  string 	  `gorm:"size:16"`
	PgLastName  string     `gorm:"size:45;not null;`
	PgFirstName  string    `gorm:"size:45;not null;`
	PgMiddleInitial  string    `gorm:"size:1;not null;`
	PgEmail string    	 `gorm:"size:45;not null;`
	IdPerson int
	Person *Person		 `gorm:"foreignkey:IdPerson"`
}
