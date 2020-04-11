package models

import (
	"errors"
	"sync"

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
	// AddRep adds a repID linked with a user ID to be repped
	// into a temp map
	AddRep(userID, repID string)

	// Clean will delete a repID event from the tmp map
	Clean(repID string)

	// Create inserts a new value into the database
	Create(rep *Rep) error

	// Exists will check if an user is in the database
	Exists(userID string) bool

	// GetRep returns the rep number of an individual
	GetRep(userID string) int

	// RepIDExists returns true if a given RepID event exists in
	// tmp map
	RepIDExists(repID string) bool

	// GetUser will return the userID behind the repID event
	// saved in tmp map
	GetUser(repID string) string

	// Increase will increase the rep number in the database for a given
	// user by one
	Increase(userID string) error
}

type repGorm struct {
	db      *gorm.DB
	tmpReps map[string]string
	m       *sync.RWMutex
}

type repService struct {
	RepDB
}

type repValidator struct {
	RepDB
}

var _ RepDB = &repGorm{}

func (rg repGorm) Clean(repID string) {
	rg.m.Lock()
	defer rg.m.Unlock()
	delete(rg.tmpReps, repID)
}

// Increase will increase the rep number in the database for a given
// user by one
func (rg repGorm) Increase(userID string) error {
	var rep Rep
	db := rg.db.Where("discord_id = ?", userID)
	err := first(db, &rep)
	if err != nil {
		return nil
	}
	rep.RepNum++
	return rg.db.Save(rep).Error
}

// GetUser will return the userID behind the repID event
// saved in tmp map
func (rg repGorm) GetUser(repID string) string {
	rg.m.RLock()
	defer rg.m.RUnlock()
	return rg.tmpReps[repID]
}

// RepIDExists returns true if a given RepID event exists in
// tmp map
func (rg repGorm) RepIDExists(repID string) bool {
	rg.m.RLock()
	defer rg.m.RUnlock()
	if _, ok := rg.tmpReps[repID]; ok {
		return true
	}
	return false
}

// AddRep adds a repID linked with a user ID to be repped
// into a temp map
func (rg repGorm) AddRep(userID, repID string) {
	rg.m.Lock()
	defer rg.m.Unlock()
	rg.tmpReps[repID] = userID
}

// NewRepService creates the rep service object
func NewRepService(db *gorm.DB) RepService {
	return &repService{
		RepDB: &repValidator{
			RepDB: &repGorm{
				db:      db,
				tmpReps: make(map[string]string),
				m:       &sync.RWMutex{},
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
