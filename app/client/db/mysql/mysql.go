package mysql

import (
	"log"

	"chatapp/pkg/clients/mysql"
)

type repoStateValue = int

const (
	stateRepoPreInit repoStateValue = iota
	stateRepoReady
)

// Repo repository
type Repo struct {
	mysql.Base
	logger *log.Logger
	opts   Options
	state  repoStateValue
}

// Options defines the config options for the mysql repo
type Options struct {
	DBIP           string
	User           string
	Pass           string
	DBName         string
	ConnRetryCount int
}

// New return a new instance of db
func New(logger *log.Logger, opts Options) *Repo {
	return &Repo{
		logger: logger,
		opts:   opts,
		state:  stateRepoPreInit,
	}
}

// isReady returns true if the repo is connected and initialized
func (r *Repo) isReady() bool {
	return !(r.SQLDB == nil || r.state < stateRepoReady)
}

// Start will connect to the db and get it ready to go.
func (r *Repo) Start() error {
	err := r.Base.ConnectInit(r.logger, r.opts.DBIP, r.opts.User,
		r.opts.Pass, r.opts.DBName, r.opts.ConnRetryCount)
	if err != nil {
		return err
	}
	r.state = stateRepoReady
	return nil
}

// Stop the mysql client
func (r *Repo) Stop() error {
	return r.Close()
}
