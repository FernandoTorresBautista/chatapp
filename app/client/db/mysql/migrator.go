package mysql

import (
	"embed"
	"fmt"
	"log"
	"net/http"
	"time"

	"chatapp/app/config"

	// migrate documentation: https://github.com/golang-migrate/migrate

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/source/httpfs"

	_ "github.com/golang-migrate/migrate/v4/database/mysql"
)

//go:embed scripts
var migrationScripts embed.FS

type Migrator struct {
	cfg    *config.Configuration
	logger *log.Logger
}

func NewMigrtor(logger *log.Logger, cfg *config.Configuration) *Migrator {
	return &Migrator{
		cfg:    cfg,
		logger: logger,
	}
}

func (m *Migrator) Run() error {
	// get the files to migrate
	sf, err := httpfs.New(http.FS(migrationScripts), "scripts")
	if err != nil {
		return err
	}
	m.logger.Println("Migration connecting to db")
	dbURL := fmt.Sprintf("mysql://%s:%s@tcp(%s)/%s", m.cfg.MySQL.MigrateDBUser, m.cfg.MySQL.MigrateDBPass, m.cfg.MySQL.DBIP, m.cfg.MySQL.DBName)
	var mig *migrate.Migrate
	delay := 3
	for count := 0; mig == nil && count < m.cfg.MySQL.DBRetryCount; count++ {
		// sleep 3 seconds in every try
		if count > 0 {
			time.Sleep(time.Duration(delay) * time.Second)
		}
		mig, err = migrate.NewWithSourceInstance("httpfs", sf, dbURL)
		if err != nil {
			m.logger.Printf("error in migration %+v\n", err)
			// increment the delay
			delay += 3
		}
	}
	if mig != nil {
		defer mig.Close()
	}
	if err != nil {
		return err
	} else if mig == nil {
		return fmt.Errorf("error in migration, ran our of range attemps")
	}
	m.logger.Println("starting the migration")
	if err := mig.Up(); err != nil && err != migrate.ErrNoChange {
		m.logger.Printf("error migration db: %+v", err)
		return err
	}
	v, _, err := mig.Version()
	if err != nil {
		return err
	}
	m.logger.Printf("migrate active version: %d", v)
	return nil
}
