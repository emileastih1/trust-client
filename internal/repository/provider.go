package repository

import (
	"bulletin-board-api/internal/constants"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/google/wire"
	"github.com/lib/pq"
	"github.com/spf13/viper"
	sqltrace "gopkg.in/DataDog/dd-trace-go.v1/contrib/database/sql"
	gormtrace "gopkg.in/DataDog/dd-trace-go.v1/contrib/gorm.io/gorm.v1"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/plugin/dbresolver"
)

const (
	maxIdleConnectionCount = 50
	maxOpenConnectionCount = 200
	maxConnectionLifetime  = time.Hour
)

func InitializeDB() (*gorm.DB, error) {
	sqltrace.Register("postgres", &pq.Driver{}, sqltrace.WithServiceName(constants.BulletinBoardServiceName))

	dbURL := viper.GetString("DB_URL")
	if len(dbURL) == 0 {
		dbURL = viper.GetString("database.URL")
	}

	dbWriterConn, err := sqltrace.Open("postgres", dbURL)
	if err != nil {
		return nil, fmt.Errorf("sqltrace open failed: %w", err)
	}

	logLevel := logger.Info
	if viper.GetString("stage") == constants.StageProduction ||
		viper.GetString("stage") == constants.StagePreprod {
		logLevel = logger.Warn
	}

	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second, // Slow SQL threshold
			LogLevel:                  logLevel,    // Log level
			IgnoreRecordNotFoundError: false,       // Ignore ErrRecordNotFound error for logger
			Colorful:                  true,        // Enable color
		},
	)

	conn, err := gormtrace.Open(
		postgres.New(postgres.Config{Conn: dbWriterConn}),
		&gorm.Config{
			AllowGlobalUpdate: false,
			CreateBatchSize:   25,
			Logger:            newLogger,
			PrepareStmt:       true,
			NowFunc: func() time.Time {
				return time.Now().UTC()
			},
		},
		gormtrace.WithServiceName(constants.BulletinBoardServiceName),
		gormtrace.WithAnalytics(true),
	)
	if err != nil {
		return nil, fmt.Errorf("init gorm failed: %w", err)
	}

	// Configure read replica
	DbReadURL := viper.GetString("DB_READ_URL")
	if len(DbReadURL) > 0 {
		dbReadConn, err := sqltrace.Open("postgres", DbReadURL)
		if err != nil {
			return nil, fmt.Errorf("sqltrace open failed: %w", err)
		}
		dbReadReplica := postgres.New(postgres.Config{Conn: dbReadConn})
		// create db resolver
		err = conn.Use(dbresolver.Register(dbresolver.Config{
			Sources:  []gorm.Dialector{postgres.New(postgres.Config{Conn: dbWriterConn})},
			Replicas: []gorm.Dialector{dbReadReplica},
			Policy:   dbresolver.RandomPolicy{},
		}))
		if err != nil {
			return nil, err
		}
	}

	genericDB, err := conn.DB()
	if err != nil {
		log.Fatal("init gorm failed", err)
		return nil, fmt.Errorf("init gorm failed: %w", err)
	}

	if err = genericDB.Ping(); err != nil {
		return nil, fmt.Errorf("can not connect to database, ping failure: %w", err)
	}

	genericDB.SetMaxIdleConns(maxIdleConnectionCount)
	genericDB.SetMaxOpenConns(maxOpenConnectionCount)
	genericDB.SetConnMaxLifetime(maxConnectionLifetime)
	return conn, nil
}

var DBProvider = wire.NewSet(
	InitializeDB, NewVaspRepository, NewAddressRepository,
)
