package models

import (
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/bwmarrin/discordgo"
)

// EventService is a layer of abstraction leading to the Event interface
type EventService interface {
	Event
}

// Event represents all methods we can use to interact with Event type data
type Event interface {
	AddEvent(User *discordgo.User, MsgID string, limit int)
	EventExists(msgID string) bool
	AddToQueue(UserID *discordgo.User, eventID string) (*discordgo.User, error)
	GetQueue(eventID string) *[]QueueUser
	Close(eventID, role string, user *discordgo.User, roles []string) error
	Clean()
}

// EventData represents an event a user has created
type EventData struct {
	DiscordUser *discordgo.User
	Limit       int
	Queue       []QueueUser
	Expiration  time.Time
}

type eventStore struct {
	eb map[string]*EventData
	m  *sync.Mutex
}

type eventService struct {
	Event
}

// QueueUser represents a user queuing to an event
type QueueUser struct {
	DiscordUser *discordgo.User
}

// internal check to see if interface is implemented correctly
var _ Event = &eventStore{}

// EventExists will check if a requested event exists currently
func (es eventStore) EventExists(msgID string) bool {
	es.m.Lock()
	defer es.m.Unlock()
	if _, ok := es.eb[msgID]; ok {
		return true
	}
	return false
}

// AddEvent creates a new event on the server
func (es eventStore) AddEvent(User *discordgo.User, MsgID string, limit int) {
	es.m.Lock()
	defer es.m.Unlock()
	newQ := make([]QueueUser, 0)
	new := &EventData{
		DiscordUser: User,
		Limit:       limit,
		Queue:       newQ,
		Expiration:  time.Now().Add(2 * time.Hour),
	}
	es.eb[MsgID] = new
}

// AddToQueue will add another user to the queue who registers as long as the
// queue is not full
//
// REVIEW: Check if this is concurrent safe...
func (es eventStore) AddToQueue(User *discordgo.User, eventID string) (*discordgo.User, error) {
	es.m.Lock()
	defer es.m.Unlock()
	val := es.eb[eventID]
	if len(val.Queue) == val.Limit {
		return nil, errors.New("queue limit reached")
	}
	if val.DiscordUser.ID == User.ID {
		return nil, errors.New("you cannot queue for your own event")
	}
	for _, u := range val.Queue {
		if u.DiscordUser.ID == User.ID {
			fmt.Println(u.DiscordUser.ID)
			return nil, errors.New("user already in queue")
		}
	}
	newUser := QueueUser{
		DiscordUser: User,
	}
	val.Queue = append(val.Queue, newUser)
	es.eb[eventID].Queue = val.Queue
	return val.DiscordUser, nil
}

// GetQueue will return the current queue line
func (es eventStore) GetQueue(eventID string) *[]QueueUser {
	es.m.Lock()
	defer es.m.Unlock()
	val := es.eb[eventID]
	return &val.Queue
}

// Close will remove a event listing from the map
func (es eventStore) Close(eventID, role string, user *discordgo.User, roles []string) error {
	es.m.Lock()
	defer es.m.Unlock()
	if !containsRole(role, roles) && es.eb[eventID].DiscordUser.ID != user.ID {
		return errors.New("permission denied")
	}
	delete(es.eb, eventID)
	return nil
}

// containsRole will check if a supplied role ID which controls bot matches the list of role ids
// a member has
func containsRole(item string, container []string) bool {
	for _, r := range container {
		if item == r {
			return true
		}
	}
	return false
}

// Clean will remove event listings from the map that have exceeded time limit
//
// DO NOT CALL THIS RANDOMLY!!
//
// This should only be called in the goroutine in main (ticker to check expiration)
func (es eventStore) Clean() {
	es.m.Lock()
	defer es.m.Unlock()
	for k, v := range es.eb {
		if time.Now().Sub(v.Expiration) > 0 {
			delete(es.eb, k)
		}
	}
}

// NewEventService creates a new Event service
func NewEventService() EventService {
	return eventService{
		Event: eventStore{
			eb: make(map[string]*EventData),
			m:  &sync.Mutex{},
		},
	}
}
