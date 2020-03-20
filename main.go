package main

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"github.com/labstack/echo"
	"net/http"
	"os"
	"strconv"

	_ "github.com/jinzhu/gorm/dialects/mysql"
	log "github.com/sirupsen/logrus"
	conf "scaha_micro_member/config"
	m "scaha_micro_member/model"
	repo "scaha_micro_member/repository"
)



func handlerFunc(msg string) func(echo.Context) error {
	return func(c echo.Context) error {
		return c.String(http.StatusOK, msg)
	}
}


func allPersons(dao repo.DAO) func(echo.Context) error {
	return func(c echo.Context) error {
		var persons []m.Person
		dao.FindAll(&persons)
		return c.JSON(http.StatusOK, persons)
	}
}

func getPerson(dao repo.DAO) func(echo.Context) error {
	return func(c echo.Context) error {
		id, _ := strconv.Atoi(c.Param("id"))
		log.Debug(id)
		person, _ := dao.FindPerson(id)
		return c.JSON(http.StatusOK, person)
	}
}

func getProfile(dao repo.DAO) func(echo.Context) error {
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


func getFamily(dao repo.DAO) func(echo.Context) error {
	return func(c echo.Context) error {
		id, _ := strconv.Atoi(c.Param("id"))
		log.Debug(id)
		family, _ := dao.FindFamily(id)
		return c.JSON(http.StatusOK, family)
	}
}

func getFamilyMemberByFamily(dao repo.DAO) func(echo.Context) error {
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
	db := repo.DAO{dbgorm}

	e.GET("/person", allPersons(db))
	e.GET("/person/:id", getPerson(db))
	e.GET("/family/:id", getFamily(db))
	e.GET("/familymember/family/:id", getFamilyMemberByFamily(db))
	e.GET("/profile/:usercode/:pwd", getProfile(db))
//	e.POST("/person/:id", newPerson(db)
//	e.PUT("/user/:name/:email", updatePerson(db))
//	e.DELETE("/person/:id", deletePerson(db))
	e.Logger.Fatal(e.Start(":4000"))
}

func initLogging() {
	// Log as JSON instead of the default ASCII formatter.
	log.SetFormatter(&log.JSONFormatter{})
	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	log.SetOutput(os.Stdout)
	// Only log the info severity or above.
	log.SetLevel(log.InfoLevel)
	log.WithFields(log.Fields{
		"prefix":      "sensor",
		"temperature": -4,
	}).Info("Temperature changes")
}

func main() {


	initLogging()
	conf.InitConfiguration("./")
	connectionString := fmt.Sprintf("%s:%s@/%s?charset=utf8%sparseTime=true", conf.Properties.Db.User, conf.Properties.Db.Pass, conf.Properties.Db.Dbname,"&")
	log.Info(connectionString)

	db, err := gorm.Open("mysql", connectionString)
	db.SingularTable(true)
	db.LogMode(true)
//	db.Model(m.Family).Related(m.FamilyMember)
//	db.Model(&m.FamilyMember{}).AddForeignKey("id_family", "customers(customer_id)", "CASCADE", "CASCADE") // Foreign key need to define manually
	defer db.Close()

	if err != nil {
		fmt.Println(err.Error())
		log.Panic("failed to connect database")
	}
	handleRequest(db)
}