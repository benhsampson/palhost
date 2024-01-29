package services

import (
	"database/sql"
	"fmt"
	"log"
	env "palhost/env"
	"testing"

	"github.com/golang-migrate/migrate"
	"github.com/golang-migrate/migrate/database/postgres"
	_ "github.com/golang-migrate/migrate/source/file"
	"github.com/stretchr/testify/suite"
)

const MIGRATE_DOWN bool = false

type DBTestSuite struct {
	suite.Suite

	db *sql.DB
	m  *migrate.Migrate
}

func (s *DBTestSuite) SetupSuite() {
	config := env.NewEnvConfig("../.testing.env")
	s.db = ConnectToDB(config)
	driver, err := postgres.WithInstance(s.db, &postgres.Config{
		DatabaseName: config.DBName,
	})
	if err != nil {
		log.Fatal(err)
	}
	s.m, err = migrate.NewWithDatabaseInstance("file://migrations", config.DBName, driver)
	if err != nil {
		log.Fatal(err)
	}
	if err = s.m.Up(); err != nil {
		if err == migrate.ErrNoChange {
			fmt.Println("Migration UP skipped")
		} else {
			log.Fatal(err)
		}
	}
}

func (s *DBTestSuite) TearDownTest() {
	for _, table := range []string{"users"} {
		if _, err := s.db.Exec(fmt.Sprintf("DELETE FROM %s;", table)); err != nil {
			log.Fatal(err)
		}
	}
}

func (s *DBTestSuite) TearDownSuite() {
	if MIGRATE_DOWN {
		if err := s.m.Down(); err != nil {
			if err == migrate.ErrNoChange {
				fmt.Println("Migration DOWN skipped")
			} else {
				log.Fatal(err)
			}
		}
	}
}

func TestDBSuite(t *testing.T) {
	suite.Run(t, new(DBTestSuite))
}
