package main

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/labstack/echo"
	log "github.com/sirupsen/logrus"
	"net/http"
	"os"
	m "scaha-entity-model"
	"strconv"
)

func allPersons(dao DAO) func(echo.Context) error {
	return func(c echo.Context) error {
		var persons []m.Person
		err := dao.FindAll(&persons)
		if err != nil {
			log.Error(err.Error())
		}
		return c.JSON(http.StatusOK, persons)
	}
}

func getPerson(dao DAO) func(echo.Context) error {
	return func(c echo.Context) error {
		id, _ := strconv.Atoi(c.Param("id"))
		log.Debug(id)
		person, _ := dao.FindPerson(id)
		return c.JSON(http.StatusOK, person)
	}
}

func getProfile(dao DAO) func(echo.Context) error {
	return func(c echo.Context) error {
		usercode := c.Param("usercode")
		pwd := c.Param("pwd")
		profile,err := dao.FindProfile(usercode, pwd)
		if err != nil {
			return c.JSON(http.StatusNotFound,"Cannot find Profile")
		}
		person,err2 := dao.FindPersonByProfile(*profile)
		if err2 != nil {
			return c.JSON(http.StatusNotFound,"Cannot find Person For profile")
		}
		family,err3 := dao.FindFamilyByPerson(*person)
		if err3 != nil {
			return c.JSON(http.StatusNotFound,"Cannot find Family for Profile")
		}
		return c.JSON(http.StatusOK, family)
	}
}

func getProfileAndRoles(dao DAO) func(echo.Context) error {
	return func(c echo.Context) error {
		usercode := c.Param("usercode")
		pwd := c.Param("pwd")
		profile,err := dao.FindProfile(usercode, pwd)
		if err != nil {
			return c.JSON(http.StatusNotFound,"Cannot find Profile")
		}
		profile.Person.Profile = nil
		return c.JSON(http.StatusOK, profile)
	}
}

func getFamily(dao DAO) func(echo.Context) error {
	return func(c echo.Context) error {
		id, _ := strconv.Atoi(c.Param("id"))
		log.Debug(id)
		family, _ := dao.FindFamily(id)
		return c.JSON(http.StatusOK, family)
	}
}

func getFamilyMemberByFamily(dao DAO) func(echo.Context) error {
	return func(c echo.Context) error {
		id, _ := strconv.Atoi(c.Param("id"))
		log.Debug(id)
		family, _ := dao.FindFamilyMemberByFamilyId(id)
		return c.JSON(http.StatusOK, family)
	}
}


/*func newPerson(dao repo.DAO) func(echo.Context) error {
	return func(c echo.Context) error {
		name := c.Param("name")
		email := c.Param("email")
		dao.create(&m.Person{Name: name, Email: email})
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
*/

//
// Lets handle these bad boys
//
func handleRequest(dbgorm *gorm.DB) {

	e := echo.New()
	db := DAO{dbgorm}

	e.GET("/person", allPersons(db))
	e.GET("/person/:id", getPerson(db))
	e.GET("/family/:id", getFamily(db))
	e.GET("/familymember/family/:id", getFamilyMemberByFamily(db))
	e.GET("/profile/:usercode/:pwd", getProfile(db))


	if Properties.ExternalMS.IsHTTPS {
		e.Logger.Fatal(e.StartTLS(fmt.Sprintf(":%d", Properties.ExternalMS.Port), "./server.crt","./server.key"))
	} else {
		e.Logger.Fatal(e.Start(fmt.Sprintf(":%d", Properties.ExternalMS.Port)))
	}

}

func init() {

	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.TraceLevel)
	InitConfiguration("./")

}

func main() {

	// Lets hook up the database and launch the microservice

	connectionString := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8%sparseTime=true", Properties.Db.User, Properties.Db.Pass, Properties.Db.Host, Properties.Db.Port, Properties.Db.Dbname,"&")
	db, err := gorm.Open("mysql", connectionString)
	if err != nil || db == nil {
		log.Error(err.Error())
		log.Panic("failed to connect database")
	} else {
		log.Info("Connected to the database with the following String: %s", connectionString)
		db.SingularTable(true)
		defer db.Close()
		handleRequest(db)
	}
}