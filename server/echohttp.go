package server

import "goserver/router"

func Run() {

	e := router.New()

	e.Start(":8001")
}
