package server

import (
	"goserver/database"
	"goserver/router"
	"goserver/templates"
	"log"
)

func Run() {

	client, err := database.New()
	if err != nil {
		log.Fatal(err)
	}
	templates.Init()

	e := router.New(client)

	e.Logger.Fatal(e.Start(":1323"))

}
