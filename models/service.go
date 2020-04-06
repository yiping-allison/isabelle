package models

import (
	"github.com/jinzhu/gorm"
)

// Code in this file adapted from Jon Calhoun

// Services handles services for bot
type Services struct {
	db     *gorm.DB
	events map[string]EventData
	Entry  EntryService
	Event  EventService
}

// ServicesConfig represents functions that are meant to be running configurations
// for services
type ServicesConfig func(*Services) error

// NewServices initializes the configuration for every service
func NewServices(cfgs ...ServicesConfig) (*Services, error) {
	var s Services
	for _, cfg := range cfgs {
		if err := cfg(&s); err != nil {
			return nil, err
		}
	}
	return &s, nil
}

// WithGorm opens a database connection using the Gorm package and
// sets the database
func WithGorm(dialect, connectionInfo string) ServicesConfig {
	return func(s *Services) error {
		db, err := gorm.Open(dialect, connectionInfo)
		if err != nil {
			return err
		}
		s.db = db
		return nil
	}
}

// WithEntries will initialize entry database
func WithEntries() ServicesConfig {
	return func(s *Services) error {
		s.Entry = NewEntryService(s.db)
		return nil
	}
}

// WithEvents will initialize an events server 'database'
func WithEvents() ServicesConfig {
	return func(s *Services) error {
		s.Event = NewEventService()
		return nil
	}
}

// WithLogMode makes sure that every database interaction in logged whether
// for debugging or other logging purposes
func WithLogMode(mode bool) ServicesConfig {
	return func(s *Services) error {
		s.db.LogMode(mode)
		return nil
	}
}

// Close will close the database connection
func (s Services) Close() error {
	return s.db.Close()
}

// AutoMigrate attempts to automigrate sql tables
func (s Services) AutoMigrate() error {
	return s.db.AutoMigrate(&Entry{}).Error
}
