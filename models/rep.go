package models

import (
	"errors"

	"github.com/jinzhu/gorm"
)

const (
	ErrDiscordIDRequired string = "need discord user ID"
)

type Rep struct {
	gorm.Model
	DiscordID string `gorm:"not_null;unique_index"`
	RepNum    int    `gorm:"not_null"`
}

type RepService interface {
	RepDB
}

type RepDB interface {
	Create(rep *Rep) error
	Exists(userID string) bool
	GetRep(userID string) int
}

type repGorm struct {
	db *gorm.DB
}

type repService struct {
	RepDB
}

type repValidator struct {
	RepDB
}

var _ RepDB = &repGorm{}

// NewRepService creates the rep service object
func NewRepService(db *gorm.DB) RepService {
	return &repService{
		RepDB: &repValidator{
			RepDB: &repGorm{
				db: db,
			},
		},
	}
}

type repValFn func(*Rep) error

func runRepValFns(rep *Rep, fns ...repValFn) error {
	for _, fn := range fns {
		if err := fn(rep); err != nil {
			return err
		}
	}
	return nil
}

// Helper func to check if provided rep item contains a discord ID
//
// This is required.
func (rv *repValidator) discordIDRequired(rep *Rep) error {
	if rep.DiscordID == "" {
		return errors.New(ErrDiscordIDRequired)
	}
	return nil
}

// Create inserts a new value into the database
func (rv *repValidator) Create(rep *Rep) error {
	err := runRepValFns(rep, rv.discordIDRequired)
	if err != nil {
		return err
	}
	return rv.RepDB.Create(rep)
}

// Create inserts a new value into the database
func (rg *repGorm) Create(rep *Rep) error {
	return rg.db.Create(rep).Error
}

// Exists will check if an user is in the database
func (rg *repGorm) Exists(userID string) bool {
	var user Rep
	db := rg.db.Where("discord_id = ?", userID)
	err := first(db, &user)
	if err != nil {
		return false
	}
	return true
}

// GetRep returns the rep number of an individual
func (rg *repGorm) GetRep(userID string) int {
	var user Rep
	db := rg.db.Where("discord_id = ?", userID)
	err := first(db, &user)
	if err != nil {
		return -1
	}
	return user.RepNum
}
