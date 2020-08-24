package main

import (
	"encoding/json"
	"fmt"
	m "github.com/davidbigelow25/scaha-entity-model"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
)


var TRUE bool = true
var FALSE bool = false
//
// News
//
func allActivePublishedNewsItems(dao DAO) func(echo.Context) error {
	return func(c echo.Context) error {
		var news []m.NewsItem
		err := dao.FindActivePublishedNewsItems(&news)
		if err != nil {
			log.Error(err.Error())
		}
		return c.JSON(http.StatusOK, news)
	}
}

//
// CLUBS
//
func allClubs(dao DAO) func(echo.Context) error {
	return func(c echo.Context) error {
		var clubs []m.Club
		err := dao.FindAllClubs(&clubs)
		if err != nil {
			log.Error(err.Error())
		}
		return c.JSON(http.StatusOK, clubs)
	}
}

//
// Venues
//
func allVenues(dao DAO) func(echo.Context) error {
	return func(c echo.Context) error {
		venues, err := dao.FindAllVenues()
		if err != nil {
			return c.JSON(http.StatusInternalServerError, venues)
		}
		return c.JSON(http.StatusOK, venues)
	}
}

func allPersons(dao DAO) func(echo.Context) error {
	return func(c echo.Context) error {
		var persons []m.Person
		err := dao.FindAllPersons(&persons)
		if err != nil {
			log.Error(err.Error())
		}
		return c.JSON(http.StatusOK, persons)
	}
}

func getAllLiveGamesByVenueandDate(dao DAO) func(echo.Context) error {
	return func(c echo.Context) error {
		venuetag := c.Param("venuetag")
		livegames, _ := dao.FindAllLiveGamesByVenueandDate(venuetag, "FUTURE TARGET DATE GOES HERE")
		return c.JSON(http.StatusOK, livegames)
	}
}

func getPerson(dao DAO) func(echo.Context) error {
	return func(c echo.Context) error {
		// claim := c.Get("claim").(*m.Claims)  This gets back the claims fished out from the jwt token
		id, _ := strconv.Atoi(c.Param("id"))
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

func getLiveGame(dao DAO) func(echo.Context) error {
	return func(c echo.Context) error {
		id, _ := strconv.Atoi(c.Param("id"))
		livegame, err := dao.FindLiveGame(id)

		if err != nil {
			log.Println(err)
		}
		return c.JSON(http.StatusOK, livegame)
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

// Here we get the live game
// and we update it possibly

func putLiveGame(dao DAO)  func(echo.Context) error {
	return func(c echo.Context) error {

		changeList := map[string]interface{}{}
		if err := c.Bind(&changeList); err != nil {
			return err
		}

		id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
		log.Printf("Here is the live game (id:%d): %+v\n",id,changeList)
		lg := dao.UpdateLiveGame(uint(id),changeList)
		return c.JSON(http.StatusOK, lg)
	}
}

// The real key is LiveGame and Roster
// We simply kick back what the new record looks like
func upsertMiaByLiveGameAndRoster(dao DAO)  func(echo.Context)error {
	return func(c echo.Context) error {
		idlg, _ := strconv.ParseUint(c.Param("idlivegame"), 10, 64)
		idr, _ := strconv.ParseUint(c.Param("idroster"), 10, 64)
		myMia := dao.UpsertMia(uint(idlg), uint(idr),true)
		return c.JSON(http.StatusOK, myMia)
	}
}

func insertScore (dao DAO)  func(echo.Context)error {
	return func(c echo.Context) error {

		idlg, _ := strconv.ParseUint(c.Param("idlivegame"), 10, 64)
		idt, _ := strconv.ParseUint(c.Param("idteam"), 10, 64)
		heshootshescores := m.Scoring{IdLiveGame: uint(idlg), IdTeam: uint(idt)}
		if err := c.Bind(&heshootshescores); err != nil {
			return err
		}

		dao.CreateScoring(&heshootshescores)
		return c.JSON(http.StatusOK, heshootshescores)
	}
}

func updateScore(dao DAO, active bool)  func(echo.Context)error {
	return func(c echo.Context) error {
		idlg, _ := strconv.ParseUint(c.Param("idlivegame"), 10, 64)
		idt, _ := strconv.ParseUint(c.Param("idteam"), 10, 64)
		ids, _ := strconv.ParseUint(c.Param("idscore"), 10, 64)
		heshootshescores := m.Scoring{IdLiveGame: uint(idlg), IdTeam: uint(idt), ID: uint(ids)}
		if err := c.Bind(&heshootshescores); err != nil {
			return err
		}
		dao.UpdateScoring(&heshootshescores,active)
		return c.JSON(http.StatusOK, heshootshescores)
	}
}

func insertPenalty (dao DAO)  func(echo.Context)error {
	return func(c echo.Context) error {
		idlg, _ := strconv.ParseUint(c.Param("idlivegame"), 10, 64)
		idt, _ := strconv.ParseUint(c.Param("idteam"), 10, 64)
		idroster, _ := strconv.ParseUint(c.Param("idroster"), 10, 64)
		pbox := m.Penalty{IdLiveGame: uint(idlg), IdTeam: uint(idt), IdRoster: uint(idroster), IsActive: &TRUE}
		if err := c.Bind(&pbox); err != nil {
			return err
		}
		dao.CreatePenalty(&pbox)
		return c.JSON(http.StatusOK, pbox)
	}
}

func updatePenalty(dao DAO, active bool)  func(echo.Context)error {
	return func(c echo.Context) error {
		idlg, _ := strconv.ParseUint(c.Param("idlivegame"), 10, 64)
		idteam, _ := strconv.ParseUint(c.Param("idteam"), 10, 64)
		idroster, _ := strconv.ParseUint(c.Param("idroster"), 10, 64)
		idpenalty, _ := strconv.ParseUint(c.Param("idpenalty"), 10, 64)
		p := m.Penalty{IdLiveGame: uint(idlg), IdTeam: uint(idteam), IdRoster: uint(idroster), ID: uint(idpenalty)}
		if err := c.Bind(&p); err != nil {
			return err
		}
		dao.UpdatePenalty(&p,active)
		return c.JSON(http.StatusOK, p)
	}
}

func insertSog (dao DAO)  func(echo.Context)error {
	return func(c echo.Context) error {
		idlg, _ := strconv.ParseUint(c.Param("idlivegame"), 10, 64)
		idt, _ := strconv.ParseUint(c.Param("idteam"), 10, 64)
		idroster, _ := strconv.ParseUint(c.Param("idroster"), 10, 64)
		sog := m.Sog{IdLiveGame: uint(idlg), IdTeam: uint(idt), IdRoster: uint(idroster), IsActive: &TRUE}
		if err := c.Bind(&sog); err != nil {
			return err
		}
		dao.CreateSog(&sog)
		return c.JSON(http.StatusOK, sog)
	}
}

func updateSog(dao DAO, active bool)  func(echo.Context)error {
	return func(c echo.Context) error {
		idlg, _ := strconv.ParseUint(c.Param("idlivegame"), 10, 64)
		idteam, _ := strconv.ParseUint(c.Param("idteam"), 10, 64)
		idroster, _ := strconv.ParseUint(c.Param("idroster"), 10, 64)
		idsog, _ := strconv.ParseUint(c.Param("idsog"), 10, 64)
		sog := m.Sog{IdLiveGame: uint(idlg), IdTeam: uint(idteam), IdRoster: uint(idroster), ID: uint(idsog)}
		if err := c.Bind(&sog); err != nil {
			return err
		}
		dao.UpdateSog(&sog,active)
		return c.JSON(http.StatusOK, sog)
	}
}

// The real key is LiveGame and Roster
// We simply kick back what the new record looks like
func deleteMiaByLiveGameAndRoster(dao DAO)  func(echo.Context)error {
	return func(c echo.Context) error {
		idlg, _ := strconv.ParseUint(c.Param("idlivegame"), 10, 64)
		idr, _ := strconv.ParseUint(c.Param("idroster"), 10, 64)
		myMia := dao.UpsertMia(uint(idlg), uint(idr), false)
		return c.JSON(http.StatusOK, myMia)
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
// Thie allows us to add a middleware piece for validating things
//
func asValidate (next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		req, err := http.NewRequest("GET", "http://localhost:4040/validate", nil)
		if err != nil {
			return c.String(http.StatusInternalServerError,"Cannot Authentidate")
		}
		//
		// lets fish out the jwt token
		jwt, err := c.Cookie("jwt")
		if err != nil {
			return c.String(http.StatusForbidden,"Cannot Find Authentication Token")
		}
		req.AddCookie(jwt)
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			return c.String(http.StatusInternalServerError,"Cannot Make the call to the authentication servcie")
		}
		defer resp.Body.Close()
		if resp.StatusCode != 200 {
			return c.String(http.StatusInternalServerError,"Cannot Make the call to the authentication servcie " + strconv.Itoa(resp.StatusCode))
		}
		bodyBytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return c.String(http.StatusInternalServerError,"Cannot Read the Body " + strconv.Itoa(resp.StatusCode))
		}
		// Lets shove it in a the context of the request and pass it on down.
		// it will be used by the rest of the software to help with everything
		claim := m.Claims{}
		json.Unmarshal(bodyBytes, &claim)
		c.Set("claim", &claim)
		return next(c)
	}
}

//
// Lets handle these bad boys
//
func handleRequest(dbgorm *gorm.DB) {

	e := echo.New()
	db := DAO{dbgorm}
	e.Debug = true

	//
	// Restricted group
	// This is an internal call made by all other microservices
	//
	e.Use(middleware.Recover())
	e.Use(middleware.Logger())

	setRoutes(e, &db)
	if Properties.ExternalMS.IsHTTPS {
		e.Logger.Fatal(e.StartTLS(fmt.Sprintf(":%d", Properties.ExternalMS.Port), "./keys/server.crt","./keys/server.key"))
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
	db, err := gorm.Open(Properties.Db.Dialect, connectionString)
	db.LogMode(true)

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