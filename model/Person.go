package model

import (
	"errors"
	"github.com/jinzhu/gorm"
	"html"
	"strings"
	"time"
)

type Person struct {
	gorm.Model
	FirstName  string    `gorm:"size:45;not null;`
	LastName  string     `gorm:"size:45;not null;`
	Email string    	 `gorm:"size:45;not null;`
	Phone  string    	 `gorm:"size:14"`
	Adress string 		 `gorm:"size:75"`
	City string 		 `gorm:"size:45"`
	State string 		 `gorm:"size:2"`
	Zipcode int
	Country string  	 `gorm:"size:45"`
	Gender string    	 `gorm:"size:1"`
	Dob time.Time
	Citizenship string 	  `gorm:"size:45"`
	BirthCertificate BitBool
	IdProfile int
	Profile *Profile		 `gorm:"foreignkey:IdProfile;PRELOAD:false"`
	IdFamily int
	Family *Family        `gorm:"foreignkey:IdPerson"`
	UsaHockeys []UsaHockey `gorm:"foreignkey:IdPerson"`

}
//
// Hey, lets be smart and filter out all the garbaded that can come it
// make it html safeish
func (p *Person) Prepare() {
	p.ID = 0
	p.FirstName = html.EscapeString(strings.TrimSpace(p.FirstName))
	p.LastName = html.EscapeString(strings.TrimSpace(p.LastName))
	p.CreatedAt = time.Now()
	p.UpdatedAt = time.Now()
}

//
// Here are all the validation rules we can apply here for the given object
//
func (p *Person) Validate() error {

	if p.FirstName == "" {
		return errors.New("First Name Required")
	}
	if p.LastName == "" {
		return errors.New("Last Name")
	}
	return nil
}