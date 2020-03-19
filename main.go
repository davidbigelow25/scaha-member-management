package main

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"github.com/labstack/echo"
	"net/http"
	"strconv"
)

type dbops struct {
	db *gorm.DB
}


//
// Here is all the DAO shit
// Not sure it belongs here
func (d dbops) findAll(persons *[]Person) error {
	return d.db.Find(persons).Error
}

func (d dbops) create(person *Person) error {
	return d.db.Create(person).Error
}

func (d dbops) findByPage(person *[]Person, page, view int) error {
	return d.db.Limit(view).Offset(view * (page - 1)).Find(&person).Error

}

func (d dbops) updateByName(name, email string) error {
	var person Person
	d.db.Where("name=?", name).Find(&person)
	person.Email = email
	return d.db.Save(&person).Error
}

func (d dbops) deleteByName(name string) error {
	var person Person
	d.db.Where("name=?", name).Find(&person)
	return d.db.Delete(&person).Error
}

func handlerFunc(msg string) func(echo.Context) error {
	return func(c echo.Context) error {
		return c.String(http.StatusOK, msg)
	}
}

func allPersons(dbobj dbops) func(echo.Context) error {
	return func(c echo.Context) error {
		var persons []Person
		dbobj.findAll(&persons)
		fmt.Println("{}", persons)

		return c.JSON(http.StatusOK, persons)
	}
}

func newPerson(dbobj dbops) func(echo.Context) error {
	return func(c echo.Context) error {
		name := c.Param("name")
		email := c.Param("email")
		dbobj.create(&Person{Name: name, Email: email})
		return c.String(http.StatusOK, name+" user successfully created")
	}
}

func deletePerson(dbobj dbops) func(echo.Context) error {
	return func(c echo.Context) error {
		name := c.Param("name")

		dbobj.(name)

		return c.String(http.StatusOK, name+" user successfully deleted")
	}
}

func updateUser(dbobj dbops) func(echo.Context) error {
	return func(c echo.Context) error {
		name := c.Param("name")
		email := c.Param("email")
		dbobj.updateByName(name, email)
		return c.String(http.StatusOK, name+" user successfully updated")
	}
}

func usersByPage(dbobj dbops) func(echo.Context) error {
	return func(c echo.Context) error {
		limit, _ := strconv.Atoi(c.QueryParam("limit"))
		page, _ := strconv.Atoi(c.QueryParam("page"))
		var result []User
		dbobj.findByPage(&result, page, limit)
		return c.JSON(http.StatusOK, result)
	}
}
func handleRequest(dbgorm *gorm.DB) {
	e := echo.New()
	db := dbops{dbgorm}

	e.GET("/person", allPersons(db))
	e.POST("/person/:id", newPerson(db))
	e.PUT("/user/:name/:email", updatePerson(db))
	e.DELETE("/person/:id", deletePerson(db))

	e.Logger.Fatal(e.Start(":3000"))
}


func main() {
	fmt.Println("Go ORM tutorial")
	db, err := gorm.Open("sqlite3", "sqlite3gorm.db")
	if err != nil {
		fmt.Println(err.Error())
		panic("failed to connect database")
	}
	defer db.Close()
	handleRequest(db)
}