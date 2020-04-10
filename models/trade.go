package models

import (
	"sync"

	"github.com/bwmarrin/discordgo"
)

// TradeService wraps to the Trade interface
type TradeService interface {
	Trade
}

// Trade contains all the methods we cna use to interact with
// trade data
type Trade interface {
	AddTrade(tradeID string, user *discordgo.User)
	Exists(tradeID string) bool
}

// TradeData represents all data needed to keep
// track of a trade event
type TradeData struct {
	DiscordUser *discordgo.User
	Offers      []TradeOfferer
}

// TradeOfferer defines someone offering a response to a trade
type TradeOfferer struct {
	user *discordgo.User
}

type tradeStore struct {
	ts map[string]*TradeData
	m  *sync.RWMutex
}

type tradeService struct {
	Trade
}

var _ Trade = &tradeStore{}

// AddTrade will add a new trade event to tracking
func (ts tradeStore) AddTrade(tradeID string, user *discordgo.User) {
	ts.m.Lock()
	defer ts.m.Unlock()
	o := make([]TradeOfferer, 0)
	new := TradeData{
		DiscordUser: user,
		Offers:      o,
	}
	ts.ts[tradeID] = &new
}

// Exists returns true if an event with the trade ID exists
func (ts tradeStore) Exists(tradeID string) bool {
	ts.m.RLock()
	defer ts.m.RUnlock()
	if _, ok := ts.ts[tradeID]; ok {
		return true
	}
	return false
}

// NewTradeService initializes a new Trade Service
func NewTradeService() TradeService {
	return tradeService{
		Trade: tradeStore{
			ts: make(map[string]*TradeData),
			m:  &sync.RWMutex{},
		},
	}
}
