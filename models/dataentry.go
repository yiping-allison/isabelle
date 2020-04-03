package models

import (
	"errors"

	"github.com/jinzhu/gorm"
)

const (
	ErrNameRequired = "models: name is required"
	ErrNotFound     = "models: entry not found"
)

// Entry represents a database entry of either an insect or
// fish in the postgres database
type Entry struct {
	Name      string `gorm:"not null"`
	SellPrice int    `gorm:"column:sell_price"`
	NorthSt   string `gorm:"type:varchar(255);column:north_start"`
	NorthEnd  string `gorm:"type:varchar(255);column:north_end"`
	SouthSt   string `gorm:"type:varchar(255);column:south_start"`
	SouthEnd  string `gorm:"type:varchar(255);column:south_end"`
	Time      string `gorm:"type:varchar(255);column:time_of_day"`
	Location  string `gorm:"type:varchar(255);column:location"`
	Image     string `gorm:"type:varchar(255);column:image"`
	Type      string `gorm:"type:varchar(255);column:type"`
}

type EntryService interface {
	EntryDB
}

// EntryDB is used to interact with data entry database
//
// If the entry is found, we will return the entry and nil error
//
// If the entry is not found, we will return ErrNotFound
//
// Lastly, we will return any other errors not generated by
// this package
type EntryDB interface {
	ByName(name, tableName string) (*Entry, error)
}

type entryGorm struct {
	db *gorm.DB
}

type entryService struct {
	EntryDB
}

type entryValidator struct {
	EntryDB
}

// Internal check if we're correctly implementing interface
var _ EntryDB = &entryGorm{}

// NewEntryService creates a new service to data entry database
func NewEntryService(db *gorm.DB) EntryService {
	return &entryService{
		EntryDB: &entryValidator{
			EntryDB: &entryGorm{
				db: db,
			},
		},
	}
}

type entryValFn func(*Entry) error

func runEntryValFns(entry *Entry, fns ...entryValFn) error {
	for _, fn := range fns {
		if err := fn(entry); err != nil {
			return err
		}
	}
	return nil
}

func (ev *entryValidator) hasName(e *Entry) error {
	if e.Name == "" {
		return errors.New(ErrNameRequired)
	}
	return nil
}

// ByName looks up an entry by its name and returns the entry
// if it exists.
//
// If not, we will return an error
func (eg *entryGorm) ByName(name, tableName string) (*Entry, error) {
	var entry Entry
	db := eg.db.Table(tableName).Where("name = ?", name)
	err := first(db, &entry)
	if err != nil {
		return nil, err
	}
	return &entry, nil
}

// first will query using the provided gorm.DB and it will
// get the first item returned and place it into dst. If nothing
// is found in the query, it will return ErrNotFound
func first(db *gorm.DB, dst interface{}) error {
	err := db.First(dst).Error
	if err == gorm.ErrRecordNotFound {
		return errors.New(ErrNotFound)
	}
	return err
}
