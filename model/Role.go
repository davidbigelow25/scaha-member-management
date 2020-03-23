package model

import "github.com/jinzhu/gorm"



type Role struct {
	gorm.Model
	RoleName string
	RoleDescription string
	IdParentRole int
	DefaultRole int
	InheritedRoles Roles  `gorm:"foreignkey:IdParentRole;"`
}


//
// We have to flatten out the roles so the user gets ALL direct and interited roles
// we do not care how the app uses them
//
func (r Roles) Flatten() *Roles {

	var flatten  Roles

	for _, role := range r {
		fr := Role{role.Model, role.RoleName, role.RoleDescription, role.IdParentRole,role.DefaultRole,nil}
		flatten = append(flatten,fr)
		if role.InheritedRoles != nil {
			flatten = append(flatten,*role.InheritedRoles.Flatten()...)
		}
	}
	return &flatten
}