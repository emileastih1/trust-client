package main

import (
	"bulletin-board-api/cmd"
)

func main() {
	cmd.InitConfig()
	app, err := cmd.InitApp()
	checkErr(err)
	checkErr(app.Start())
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
