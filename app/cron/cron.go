package cron

import (
	"log"

	"chatapp/app/config"
	// "github.com/robfig/cron/v3"
)

// simple example
// this is used to setup a cron if needed to check anythig you want...
type Cron struct {
	// cr     *cron.Cron
	cfg    *config.Configuration
	logger *log.Logger
	// you can add here bridges to the db repository functions or
	// others clients like kafka, etc
	//  this example add the db with just one function
	db Handle
}

// New return a new instance of the cron
func New(cfg *config.Configuration, logger *log.Logger) *Cron {
	return &Cron{
		// cr:     cron.New(),
		cfg:    cfg,
		logger: logger,
	}
}

// Handle a bridge to use db functions...
// thisis just an example on how to use it
type Handle interface {
	// functions to check
	GetExample() (interface{}, error)
}

// Start the cron if is enable
func (c *Cron) Start() error {
	if c.cfg.Cron.Enable {
		c.logger.Printf("adding jobs to the cron")
		// // adding jobs
		// if _, err := c.cr.AddFunc(c.cfg.Cron.CheckExample, c.Example); err != nil {
		// 	return err
		// }
		// // start the cron
		// c.cr.Start()
		c.logger.Printf("cron startted with schedule")
	}
	return nil
}

// Stop the cron
func (c *Cron) Stop() error {
	// c.cr.Stop()
	return nil
}

// Example job
func (c *Cron) Example() {
	c.logger.Printf("Example jobs started")
	// using the function from the db
	resp, err := c.db.GetExample()
	if err != nil {
		c.logger.Printf("CRON Error db.GetExample: %#v", err)
		return
	}
	c.logger.Printf("CRON OK db.GetExample: %#v", resp)
}
