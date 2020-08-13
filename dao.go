package main

import (
	m "github.com/davidbigelow25/scaha-entity-model"
	"github.com/jinzhu/gorm"
	log "github.com/sirupsen/logrus"
)

type DAO struct {
	DB *gorm.DB
}


//
// Eventually we will have microservices for everythign
// for for now to get thingsg going - lets put things in this
// file and group them for easy split-off down the road
//

//-------------------------------------------
// NewsItem Related Pulls
//-------------------------------------------
func (d DAO) FindActivePublishedNewsItems(news *[]m.NewsItem) error {
	return d.DB.Where("is_active = 1 AND state='publish'").Find(news).Error
}


//-------------------------------------------
// Update A Live game for the
//-------------------------------------------
func (d DAO) UpdateLiveGame(idp uint, changeMap map[string]interface{}) int {
	lg := m.Livegame{ID:idp}
	db := d.DB.Model(&lg).Assign(&lg).Updates(changeMap).First(&lg)
	return int(db.RowsAffected)
}

//-------------------------------------------
// Update A Live game for the
//-------------------------------------------
func (d DAO) UpsertMia(idlg uint, idr uint, active bool)  *m.Mia {
	mymia := m.Mia{IdLiveGame: idlg, IdRoster: idr, IsActive: &active}
	d.DB.Assign(mymia).FirstOrCreate(&mymia)
	return &mymia
}

func (d DAO) CreateScoring(scoring *m.Scoring)  *m.Scoring {
	d.DB.Create(scoring).Assign(scoring)
	return scoring
}

func (d DAO) UpdateScoring(scoring *m.Scoring, isactive bool)  *m.Scoring {
	scoring.IsActive = &isactive
	d.DB.Save(scoring).Assign(scoring)
	return scoring
}


//-------------------------------------------
// Venue Related Pulls
//-------------------------------------------
func (d DAO) FindAllVenues() (*[]m.Venue, error) {
	var venues = []m.Venue{}
	err := d.DB.Where("isactive = 1").Find(&venues).Error
	return &venues,err
}


//-------------------------------------------
// Pull all games within a time range and venue
//
// This is a very Shallow pull to give a top level
// view of information to display in a list
//-------------------------------------------
func (d DAO) FindAllLiveGamesByVenueandDate(venuetag string, date string) (*[]m.Livegame, error) {
	var livegame = []m.Livegame{}
	// We will hardcode the date for now to gernerate a small list of games for the given venue
	//
	err := d.DB.Where("venuetag = ? AND actdate between '2020-01-10' AND '2020-01-20'", venuetag).
		Preload("HomeTeam").
		Preload("AwayTeam").
		Preload("Schedule").
		Preload("Venue").
		Find(&livegame).Error
	return &livegame,err
}

//-------------------------------------------
// Lets grab a live game in all of its glory
//-------------------------------------------
func (d DAO) FindLiveGame(id int) (*m.Livegame, error) {
	var livegame = m.Livegame{}
	err := d.DB.Debug().
		Preload("HomeTeam").
		Preload("AwayTeam").
		Preload("Schedule").
		Preload("Venue").
		Preload("Mia", "isActive = 1").
		Preload("HomeShotsOnGoal","idlivegame= ? ", id).
		Preload("AwayShotsOnGoal","idlivegame= ? ", id).
		Preload("HomeScoring","idlivegame= ? ", id).
		Preload("AwayScoring","idlivegame= ? ", id).
		Preload("HomePenalties","idlivegame= ? ", id).
		Preload("AwayPenalties","idlivegame= ? ", id).
		Preload("HomeTeam.Suspensions", "NOT isServed ").
		Preload("AwayTeam.Suspensions"," NOT isServed ").
		Preload("HomeTeam.RosterSpots").
		Preload("AwayTeam.RosterSpots").
		Preload("AwayTeam.RosterSpots.Person").
		Preload("HomeTeam.RosterSpots.Person").
		Preload("HomeTeam.RosterSpots.Player").
		Preload("AwayTeam.RosterSpots.Player").
		Where("isActive = 1 AND idlivegame = ?", id).
		First(&livegame).Error

	return &livegame, err
}

//-------------------------------------------
// CLUB Related Pulls
//-------------------------------------------

func (d DAO) FindAllClubs(clubs *[]m.Club) error {
	return d.DB.Find(clubs).Error
}

//-------------------------------------------
// Person Related Pulls
//-------------------------------------------

func (d DAO) FindAllPersons(persons *[]m.Person) error {
	return d.DB.Find(persons).Error
}

// Lets find a person
func (d DAO) FindPerson(id int) (*m.Person, error) {
	var person = m.Person{}
	err := d.DB.Debug().Preload("Profile").Where("id = ?", id).First(&person).Error
	return &person, err
}

//
// Here are some simple DAO routines that will
// Lets find a person
func (d DAO) FindPersonByProfile(profile m.Profile) (*m.Person, error) {
	var person = m.Person{}
	err := d.DB.Debug().Preload("Profile").Where("id_profile = ?", profile.ID).First(&person).Error

	log.Debug(err)

	return &person, err
}

//
// Here are some simple DAO routines that will
// Lets find a profile and everone underneath it
func (d DAO) FindProfile(usercode string, pwd string) (*m.Profile, error) {
	var profile = m.Profile{}
	err := d.DB.Debug().Where("user_code = ? AND pwd = ?", usercode, pwd).
		Preload("Person").
		Preload("Roles").
		Preload("Roles.InheritedRoles").
		Preload("Roles.InheritedRoles.InheritedRoles").
		First(&profile).Error
	r := profile.Roles
	profile.Roles = *r.Flatten()
	return &profile, err
}

//
// Here are some simple DAO routines that will
// Lets find a person
// We want to control exactly how these structures get loaded because of the recursive nature of this.
//
func (d DAO) FindFamily(id int) (*m.Family, error) {
	var family = m.Family{}
	err := d.DB.Debug().
		Where("id = ?", id).
		Preload("Person").
		Preload("Person.Profile").
		Preload("Person.Profile.Roles").
		Preload("FamilyMembers").
		Preload("FamilyMembers.Person").
		First(&family).Error
	return &family, err
}

func (d DAO) FindFamilyByPerson(person m.Person) (*m.Family, error) {
	var family = m.Family{}
	err := d.DB.Debug().
		Where("id_person = ?", person.ID).
		Preload("Person").
		Preload("Person.UsaHockeys").
		Preload("Person.Profile").
		Preload("Person.Profile.Roles").
		Preload("Person.Profile.Roles.InheritedRoles").
		Preload("Person.Profile.Roles.InheritedRoles.InheritedRoles").
		Preload("Person.Profile.Roles.InheritedRoles.InheritedRoles.InheritedRoles").
		Preload("FamilyMembers").
		Preload("FamilyMembers.Person").
		Preload("FamilyMembers.Person.UsaHockeys").
		First(&family).Error
		r := family.Person.Profile.Roles
		family.Person.Profile.Roles = *r.Flatten()
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

