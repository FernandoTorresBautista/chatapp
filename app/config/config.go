package config

import (
	"github.com/jinzhu/configor"
)

// Configuration ...
type Configuration struct {
	AppName string   `default:"example" env:"APP_NAME"`
	Port    uint     `default:"8000" env:"PORT"`
	DBType  string   `default:"mysql" env:"DB_TYPE"`
	MySQL   struct { // mysql
		DBIP          string `default:"127.0.0.1:3306" env:"DB_MYSQL_IP"`
		DBName        string `default:"test" env:"DB_MYSQL_NAME"`
		DBUser        string `default:"testuser" env:"DB_MYSQL_USER"`
		DBPass        string `default:"*****" env:"DB_MYSQL_PASS"`
		DBRetryCount  int    `default:"10" env:"DB_MYSQL_RETRY"`
		MigrateDBUser string `default:"testadmin" env:"MIGRATE_DB_USER"`
		MigrateDBPass string `default:"*****" env:"MIGRATE_DB_PASS"`
	}
	MigrateDB            bool   `default:"false" env:"MIGRATE_DB"`
	ContinueAfterMigrate bool   `default:"false" env:"CONTINUE_AFTER_MIGRATE"`
	Fail                 bool   `defualt:"false"`
	FailMessage          string `defualt:""`
	Cron                 struct {
		Enable       bool   `default:"false" env:"CRON_ENABLE"`
		CheckExample string `default:"@daily" env:"CRON_SCHEDULE_EXAMPLE"` // we can use also "'@every [Xh XhYm Ym, ...]'"
	}

	RabbitMQ struct {
		Host     string `default:"host" env:"RABBITMQ_HOST"`       // defualt user for rabbit
		User     string `default:"user" env:"RABBIT_USER"`         // defualt user for rabbit
		Password string `default:"password" env:"RABBIT_PASSWORD"` // default password for rabbit
	}
}

// Cfg global configuration
var Cfg = Configuration{}

func init() {
	err := configor.Load(&Cfg, "app/config/config.yml")
	if err != nil {
		Cfg.Fail = true
		Cfg.FailMessage = err.Error()
	}
}
