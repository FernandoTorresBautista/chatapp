package mysql

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

const (
	mySQLURL = "%s:%s@tcp(%s)/"
)

// MySQL Base States
const (
	stateMySQLPreInit = iota
	stateMySQLPostInit
)

// Base represents a generalized MySQL connection.  Meant to be composed into
// an implementation-specific repository struct
type Base struct {
	SQLDB          *sqlx.DB
	DBName         string
	mysqlBaseState int // see "MySQL Base States" above
}

// Ping issues a quick query to MySQL to see if the connection
// is alive.
func (m *Base) Ping() error {
	if m.SQLDB == nil {
		return fmt.Errorf("db not connected")
	}
	if m.mysqlBaseState < stateMySQLPostInit {
		return fmt.Errorf("db initialization has not run")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err := m.SQLDB.PingContext(ctx)
	if err != nil {
		return fmt.Errorf("error pinging db: %v", err)
	}
	return nil
}

// MySQLConnecCheck
func (m *Base) MySQLConnecCheck() error {
	if m.mysqlBaseState < stateMySQLPostInit {
		return errors.New("mysql not started")
	} else if err := m.Ping(); err != nil {
		return err
	}
	return nil
}

// Close will close the db connection
func (m *Base) Close() (e error) {
	if m.SQLDB != nil {
		e = m.SQLDB.Close()
	}
	return
}

// ConnectInit will connect to the given MySQL db - if the 'dbName' schema
// does not exist, it will attempt to create it.
func (m *Base) ConnectInit(logger *log.Logger, dbIP, dbUser, dbPass, dbName string, retryCount int) error {
	if len(dbIP) == 0 {
		return errors.New("missing database ip address configuration")
	}
	url, displayURL := buildURLs(dbIP, dbUser, dbPass)
	dbc, err := connect(logger, url, displayURL, dbName, retryCount)
	if err != nil {
		return fmt.Errorf("db connect failure to: %s  Reason: %w", displayURL, err)
	}
	m.SQLDB = dbc
	m.DBName = dbName
	m.SQLDB.SetMaxIdleConns(10)
	m.SQLDB.SetMaxOpenConns(100)
	m.mysqlBaseState = stateMySQLPostInit
	return nil
}

// return the url to display and use
func buildURLs(mySQLIP string, user string, sec string) (url string, displayURL string) {
	url = fmt.Sprintf(mySQLURL, user, sec, mySQLIP)
	displayURL = fmt.Sprintf(mySQLURL, "****", "****", mySQLIP) // for logging
	return
}

// connect tries ot connect to the given displayURL.  If retryCount > 0, it
// will retry that many times.  A retryCount < 0 indicates try forever.
// retryCount == 0 means only try once.
func connect(logger *log.Logger, url string, displayURL string, dbName string, retryCount int) (*sqlx.DB, error) {
	delay := time.Second * 2
	logger.Printf("connect start url: %s, db: %s", displayURL, dbName)
	var dbc *sqlx.DB
	var err error
	limit := retryCount >= 0
	for ; !limit || retryCount >= 0; retryCount-- {
		dbc, err = sqlx.Connect("mysql", url+dbName+"?parseTime=true")
		if err != nil {
			logger.Printf("database connect failed, error: %v", err)
			// If the database was not created, return right away (retrying
			//   will not resolve this)
			if dbName != "" && strings.Contains(err.Error(), "Unknown database") {
				err = fmt.Errorf("database '%s' does not exist", dbName)
				logger.Printf("database does not exist: %s", dbName)
				return nil, err
			}
			time.Sleep(delay)
			// backoff delay
			if delay < 30 {
				delay *= 2
			}
			continue
		}
		if dbc != nil {
			break
		}
	}
	if err != nil {
		return nil, err
	}
	// Catch the case where we failed to connect before we ran out of retries.
	if dbc == nil {
		return nil, fmt.Errorf("error connecting to db [" + dbName + "]")
	}
	logger.Printf("connect complete url: %s, db: %s", displayURL, dbName)
	return dbc, nil
}
