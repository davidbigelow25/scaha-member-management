package main

import "github.com/labstack/echo"

func setRoutes(e *echo.Echo, db *DAO) {
	// All the general pulls (non volitale)
	e.GET("/person", allPersons(*db))
	e.GET("/venue", allVenues(*db))
	e.GET("/news", allActivePublishedNewsItems(*db))
	e.GET("/club", allClubs(*db))
	e.GET("/person/:id", getPerson(*db), asValidate)
	e.GET("/family/:id", getFamily(*db))
	e.GET("/livegame/:id", getLiveGame(*db))
	e.GET("/livegame/byvenue/:venuetag", getAllLiveGamesByVenueandDate(*db))
	e.GET("/familymember/family/:id", getFamilyMemberByFamily(*db))
	e.GET("/profile/:usercode/:pwd", getProfile(*db))

	//
	// All the update sections
	//

	// Live Game.  Here we do a single level update (we do not worry about relationships and nested updates
	// we are handling all the using the path of the user and inserting keys in the path itself
	e.PUT("/livegame/:id",putLiveGame(*db))

	// MIA work, since the primary key is not going to be usefull and is not the ID here
	// we will do a basic upsert.  There is no body here to worry about.
	// a delete simply deactivates the data with
	e.PUT("/livegame/:idlivegame/roster/:idroster/mia",upsertMiaByLiveGameAndRoster(*db))
	e.DELETE("/livegame/:idlivegame/roster/:idroster/mia",deleteMiaByLiveGameAndRoster(*db))

	// Scoring Section  - This is how we do it when we have a key of sorts
	e.POST("/livegame/:idlivegame/team/:idteam/scoring",insertScore(*db))  // Create a scoring situation
	e.PUT("/livegame/:idlivegame/team/:idteam/scoring/:idscore",updateScore(*db,true)) // Update change something about a score for the given team
	e.DELETE("/livegame/:idlivegame/team/:idteam/scoring/:idscore",updateScore(*db, false)) // deactivate something about a score for the given team

	// Penalty Section  - This is how we do it when we have a key of sorts
	e.POST("/livegame/:idlivegame/team/:idteam/roster/:idroster/penalty",insertPenalty(*db))  // Create a scoring situation
	e.PUT("/livegame/:idlivegame/team/:idteam/roster/:idroster/penalty/:idpenalty",updatePenalty(*db,true)) // Update change something about a score for the given team
	e.DELETE("/livegame/:idlivegame/team/:idteam/roster/:idroster/penalty/:idpenalty",updatePenalty(*db, false)) // deactivate something about a score for the given team

	// Game Suspensions Section  - This is how we do it when we have a key of sorts
	e.POST("/livegame/:idlivegame/team/:idteam/person/:idperson/suspensions",insertScore(*db))  // Create a suspension
	e.PUT("/livegame/:idlivegame/team/:idteam/person/:idperson/suspensions/:idsuspensions",updateScore(*db,true)) // Update change something about a score for the given team
	e.DELETE("/livegame/:idlivegame/team/:idteam/person/:idperson/suspensions/:idsuspensions",updateScore(*db, false)) // deactivate something about a score for the given team

	// Shots On Goal Section  Section  - This is how we do it when we have a key of sorts
	e.POST("/livegame/:idlivegame/team/:idteam/roster/:idroster/sog",insertScore(*db))  // Create a suspension
	e.PUT("/livegame/:idlivegame/team/:idteam/roster/:idroster/sog/:idsog",updateScore(*db,true)) // Update change something about a score for the given team
	e.DELETE("/livegame/:idlivegame/team/:idteam/roster/:idroster/sog/:idsog",updateScore(*db, false)) // deactivate something about a score for the given team


}
