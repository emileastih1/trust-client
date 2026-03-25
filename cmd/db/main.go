package main

import (
	"bulletin-board-api/cmd"
	"bulletin-board-api/internal/repository"
	"flag"
	"log"
	"os"

	"github.com/spf13/viper"

	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
)

const (
	Postgres = "postgres"
)

var (
	flags   = flag.NewFlagSet("goose", flag.ExitOnError)
	verbose = flags.Bool("v", false, "Verbose mode")
	help    = flags.Bool("h", false, "Print help")
)

func main() {
	cmd.InitConfig()

	_ = flags.Parse(os.Args[1:])
	goose.SetVerbose(*verbose)
	db, err := repository.InitializeDB()
	if err != nil {
		panic(err)
	}

	sqlMigrationDir := viper.GetString("migration_dir")
	if len(sqlMigrationDir) == 0 {
		panic("migration dir is missing")
	}

	args := flags.Args()
	if len(args) == 0 || *help {
		flags.Usage()
		return
	}
	command := args[0]
	if err := goose.SetDialect(Postgres); err != nil {
		log.Fatal(err)
	}

	arguments := []string{}
	if len(args) > 3 {
		arguments = append(arguments, args[3:]...)
	}

	sqlDB, err := db.DB()
	if err != nil {
		panic(err)
	}

	if err := goose.Run(command, sqlDB, sqlMigrationDir, arguments...); err != nil {
		log.Fatalf("goose run: %v", err)
	}
}
