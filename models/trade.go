package models

import (
	"errors"
	"sync"
	"time"

	"github.com/bwmarrin/discordgo"
)

// TradeService wraps to the Trade interface
type TradeService interface {
	Trade
}

// Trade contains all the methods we can use to interact with
// trade data
type Trade interface {
	// Clean will remove an expired item from the map
	Clean()

	// GetExpiration returns the expiration time of the trade event
	GetExpiration(tradeID string) time.Time

	// AddOffer will track an offer to a tradeID
	//
	// This func will return an err if the user is already in trade, else nil
	AddOffer(tradeID, offer string, user *discordgo.User) error

	// AddTrade will add a new trade event to tracking
	AddTrade(tradeID string, user *discordgo.User)

	// Exists returns true if an event with the trade ID exists
	Exists(tradeID string) bool

	// GetHost returns the creator of the trade
	GetHost(tradeID string) *discordgo.User

	// Close will close a trade event. If the user does not have permission to close the event, the func
	// will return an error
	Close(tradeID string, user *discordgo.User, userRoles []string, adminID string) error
}

// TradeData represents all data needed to keep
// track of a trade event
type TradeData struct {
	DiscordUser *discordgo.User
	Expiration  time.Time
	Offers      []TradeOfferer
}

// TradeOfferer defines someone offering a response to a trade
type TradeOfferer struct {
	user  *discordgo.User
	offer string
}

type tradeStore struct {
	ts map[string]*TradeData
	m  *sync.RWMutex
}

type tradeService struct {
	Trade
}

var _ Trade = &tradeStore{}

// Clean will remove an expired item from the map
func (ts tradeStore) Clean() {
	ts.m.Lock()
	defer ts.m.Unlock()
	for k, v := range ts.ts {
		if time.Now().Sub(v.Expiration) > 0 {
			delete(ts.ts, k)
		}
	}
}

// GetExpiration returns the expiration time of the trade event
func (ts tradeStore) GetExpiration(tradeID string) time.Time {
	ts.m.RLock()
	defer ts.m.RUnlock()
	return ts.ts[tradeID].Expiration
}

// AddOffer will track an offer to a tradeID
//
// This func will return an err if the user is already in trade, else nil
func (ts tradeStore) AddOffer(tradeID, tradeOffer string, user *discordgo.User) error {
	ts.m.Lock()
	defer ts.m.Unlock()
	new := TradeOfferer{
		user:  user,
		offer: tradeOffer,
	}
	val := ts.ts[tradeID]
	if containsUser(user, val.Offers) {
		return errors.New("user already in trade")
	}
	if val.DiscordUser.ID == user.ID {
		return errors.New("you cannot offer for your own trade")
	}
	ts.ts[tradeID].Offers = append(ts.ts[tradeID].Offers, new)
	return nil
}

func containsUser(user *discordgo.User, offers []TradeOfferer) bool {
	for _, u := range offers {
		if user.ID == u.user.ID {
			return true
		}
	}
	return false
}

// Close will close a trade event. If the user does not have permission to close the event, the func
// will return an error
func (ts tradeStore) Close(tradeID string, user *discordgo.User, userRoles []string, adminID string) error {
	ts.m.Lock()
	defer ts.m.Unlock()
	val := ts.ts[tradeID]
	if val.DiscordUser.ID != user.ID && !containsRole(adminID, userRoles) {
		return errors.New("you do not have permission to close this event")
	}
	delete(ts.ts, tradeID)
	return nil
}

// GetHost returns the creator of the trade
func (ts tradeStore) GetHost(tradeID string) *discordgo.User {
	ts.m.RLock()
	defer ts.m.RUnlock()
	val := ts.ts[tradeID]
	return val.DiscordUser
}

// AddTrade will add a new trade event to tracking
func (ts tradeStore) AddTrade(tradeID string, user *discordgo.User) {
	ts.m.Lock()
	defer ts.m.Unlock()
	o := make([]TradeOfferer, 0)
	new := TradeData{
		DiscordUser: user,
		Expiration:  time.Now().Add(4 * time.Hour),
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
