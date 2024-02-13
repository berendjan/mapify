package main

import (
	"com/mapify/pagerouter"
	"log"
)

func main() {

	pageRouter, err := pagerouter.NewPageRouter("pages")
	if err != nil {
		println(err.Error())
	}
	pageRouter.Print()
	err = pageRouter.Run(3000)
	if err != nil {
		log.Printf("Error running server: %s", err.Error())
	}
}
