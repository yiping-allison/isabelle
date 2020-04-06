package models

import (
	"errors"
	"fmt"
	"sync"

	"github.com/bwmarrin/discordgo"
)

type EventService interface {
	Event
}

type Event interface {
	AddEvent(User *discordgo.User, MsgID string, limit int)
	EventExists(msgID string) bool
	AddToQueue(UserID *discordgo.User, eventID string) error
}

// EventData represents an event a user has created
type EventData struct {
	DiscordUser *discordgo.User
	Limit       int
	Queue       []QueueUser
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
	}
	es.eb[MsgID] = new
}

// AddToQueue will add another user to the queue who registers as long as the
// queue is not full
//
// REVIEW: Check if this is concurrent safe...
func (es eventStore) AddToQueue(User *discordgo.User, eventID string) error {
	es.m.Lock()
	defer es.m.Unlock()
	val := es.eb[eventID]
	if len(val.Queue) == val.Limit {
		return errors.New("limit reached")
	}
	if val.DiscordUser.ID == User.ID {
		return errors.New("you cannot queue for your own event")
	}
	for _, u := range val.Queue {
		if u.DiscordUser.ID == User.ID {
			fmt.Println(u.DiscordUser.ID)
			return errors.New("user already in queue")
		}
	}
	newUser := QueueUser{
		DiscordUser: User,
	}
	val.Queue = append(val.Queue, newUser)
	es.eb[eventID].Queue = val.Queue
	return nil
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
