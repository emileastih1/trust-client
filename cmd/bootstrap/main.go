package main

import (
	"bulletin-board-api/cmd"
	"log"
)

func main() {
	cmd.InitConfig()
	app, err := cmd.InitBootstrapApp()
	checkErr(err)
	checkErr(app.BootstrapSecrets())
}

func checkErr(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}
