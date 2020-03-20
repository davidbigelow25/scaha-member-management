package repository

import (
	"github.com/jinzhu/gorm"
	m "scaha_micro_member/model"
)

type DAO struct {
	DB *gorm.DB
}

//
// Here are some simple DAO routines that will
func (d DAO) FindAll(persons *[]m.Person) error {
	return d.DB.Find(persons).Error
}

//
// Here are some simple DAO routines that will
// Lets find a person
func (d DAO) FindPerson(id int) (*m.Person, error) {
	var person = m.Person{}
	err := d.DB.Debug().Preload("Profile").Where("id = ?", id).First(&person).Error
	return &person, err
}

//
// Here are some simple DAO routines that will
// Lets find a person
func (d DAO) FindFamily(id int) (*m.Family, error) {
	var family = m.Family{}
	err := d.DB.Debug().
		Where("id = ?", id).
		Preload("Person").
		Preload("Person.Profile").
		Preload("FamilyMembers").
		Preload("FamilyMembers.Person").
		Preload("FamilyMembers.Person.Profile").
		First(&family).Error
	return &family, err
}

//
// Here are some simple DAO routines that will
// Lets find a person
func (d DAO) FindFamilyMemberByFamilyId(id int) (*m.FamilyMember, error) {
	var familymember = m.FamilyMember{}
	err := d.DB.Debug().Where("id_family = ?", id).
		First(&familymember).Error
	return &familymember, err
}


func (d DAO) Create(person *m.Person) error {
	return d.DB.Create(person).Error
}

func (d DAO) FindByPage(person *[]m.Person, page, view int) error {
	return d.DB.Limit(view).Offset(view * (page - 1)).Find(&person).Error

}

func (d DAO) UpdateByName(name, email string) error {
	var person m.Person
	d.DB.Where("name=?", name).Find(&person)
	person.Email = email
	return d.DB.Save(&person).Error
}

func (d DAO) DeleteByName(name string) error {
	var person m.Person
	d.DB.Where("name=?", name).Find(&person)
	return d.DB.Delete(&person).Error
}

