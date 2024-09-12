package fsm

import (
	"fmt"
	"log"
	"p2pbot/internal/db/models"
)

type State int

const (
	Welcome State = iota
	AwaitingExchange
	AwaitingСurrency
	AwaitingExchangeUsername
	AwaitingSide
)

type Event int

const (
	NewTracker Event = iota
	ExchangeFound
	ExchangeNotFound
	UsernameGiven
	CurrencyGiven
	SideGiven
	AdvertisementFound
	AdvertisementNotFound
)

type Action func(args ...any)

type FSM struct {
	current     map[int64]State
	transitions map[State]map[Event]State
	actions     map[State]map[Event]Action
}

func New() *FSM {
	fsm := &FSM{
		current:     make(map[int64]State),
		transitions: make(map[State]map[Event]State),
		actions:     make(map[State]map[Event]Action),
	}

	fsm.transitions[Welcome] = map[Event]State{
		NewTracker: AwaitingExchange,
	}
	fsm.transitions[AwaitingExchange] = map[Event]State{
		ExchangeFound:    AwaitingСurrency,
		ExchangeNotFound: Welcome,
	}
	fsm.transitions[AwaitingСurrency] = map[Event]State{
		CurrencyGiven: AwaitingExchangeUsername,
	}
	fsm.transitions[AwaitingExchangeUsername] = map[Event]State{
		UsernameGiven: AwaitingSide,
	}
	fsm.transitions[AwaitingSide] = map[Event]State{
		AdvertisementFound:    Welcome,
		AdvertisementNotFound: Welcome,
	}
	fsm.actions[Welcome] = map[Event]Action{
		NewTracker: func(args ...any) {

		},
	}

	fsm.actions[AwaitingExchange] = map[Event]Action{
		ExchangeFound: func(args ...any) {
			id, ok := args[0].(int64)
			if !ok {
				log.Println(args)
				log.Fatal("arg[0] must be a int64")
			}

			ex, ok := args[1].(string)
			if !ok {
				log.Fatal("arg[1] must be a string")
			}

			t, ok := args[2].(*models.Tracker)
			if !ok {
				log.Fatal("arg[2] must be a models.Tracker")
			}

			t.Exchange = ex
			t.UserID = id

			log.Printf("New Tracker for %d on %s\n", id, ex)
		},
		ExchangeNotFound: func(args ...any) {

		},
	}

	fsm.actions[AwaitingСurrency] = map[Event]Action{
		CurrencyGiven: func(args ...any) {
			c, ok := args[0].(string)
			if !ok {
				log.Fatal("arg[0] must be a string")
			}

			t, ok := args[1].(*models.Tracker)
			if !ok {
				log.Fatal("arg[1] must be a models.Tracker")
			}

			t.Currency = c
		},
	}
	fsm.actions[AwaitingExchangeUsername] = map[Event]Action{
		UsernameGiven: func(args ...any) {

		},
	}
	fsm.actions[AwaitingSide] = map[Event]Action{
		AdvertisementFound: func(args ...any) {

		},
		AdvertisementNotFound: func(args ...any) {

		},
	}
	return fsm
}

func (fsm *FSM) Transition(chatID int64, event Event, args ...any) (State, error) {

	if newState, ok := fsm.transitions[fsm.current[chatID]][event]; ok {
		if action, ok := fsm.actions[fsm.current[chatID]][event]; ok {
			action(args...)
			fsm.current[chatID] = newState
			return newState, nil
		} else {
			log.Fatalf("Fsm invalid, action not found for state %d, event %d", fsm.current[chatID], event)
		}
	}
	return -1, fmt.Errorf("invalid transition for state %d, event %d", fsm.current[chatID], event)
}

func (fsm *FSM) GetState(id int64) State {
	return fsm.current[id]
}
