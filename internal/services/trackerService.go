package services

import (
	"fmt"
	"p2pbot/internal/db/models"
	"p2pbot/internal/db/repository"
	"strings"
)

type TrackerService struct {
	repo      *repository.TrackerRepository
	Exchanges map[string]bool
	trStaging map[int]*models.Tracker
}

func NewTrackerService(repo *repository.TrackerRepository) *TrackerService {
	exchanges := make(map[string]bool)
	exchanges["binance"] = true
	exchanges["bybit"] = true
	return &TrackerService{repo: repo, Exchanges: exchanges, trStaging: make(map[int]*models.Tracker)}
}



/*
ValidateTracker checks if tracker fields are correct and creates new tracker

tracker - tracker to create
staging - if true, tracker will be removed from staging area after creation
(use true if added to staging before)
return error if tracker is nil, side is not BUY/SELL,
currency length is not 3 or exchange is not supported
*/

func (s *TrackerService) ValidateTracker(tracker *models.Tracker, staging bool) error {
	if tracker == nil {
		return fmt.Errorf("Tracker is nil")
	}

	tracker.Side = strings.ToUpper(tracker.Side)
	if tracker.Side != "BUY" && tracker.Side != "SELL" {
		return fmt.Errorf("Side must be BUY/SELL")
	}

	tracker.Currency = strings.ToUpper(tracker.Currency)
	if len(tracker.Currency) != 3 {
		return fmt.Errorf("Currency ticker must be 3 symbols long, EUR for example")
	}

	tracker.Exchange = strings.ToLower(tracker.Exchange)
	if _, ok := s.Exchanges[tracker.Exchange]; !ok {
		return fmt.Errorf("exchange %s not supported", tracker.Exchange)
	}

	// Remove tracker from staging area
    if (staging) {
        s.DeleteTrackerStaging(tracker.UserID)
    }

	return nil
}

/*
CreateTracker creates new tracker

tracker - tracker to create
return nil if tracker created, error otherwise 
*/
func (s *TrackerService) CreateTracker(tracker *models.Tracker) error {
	return s.repo.Save(tracker)
}

func (s *TrackerService) SetWaitingFlag(id int64, flag bool) error { 
	return s.repo.UpdateWaitingFlag(id, flag)
}

func (s *TrackerService) GetAllTrackers() ([]*models.UserTracker, error) {
	return s.repo.GetAllTrackers()
}

func (s *TrackerService) GetTrackersByUserId(id int) ([]*models.UserTracker, error) {
	return s.repo.GetTrackersByUserId(id)
}

func (s *TrackerService) GetTrackerStaging(id int) *models.Tracker {
	tr, ok := s.trStaging[id]
	if !ok {
		tr := &models.Tracker{}
		s.trStaging[id] = tr
		return tr
	}
	return tr
}

func (s *TrackerService) DeleteTrackerStaging(id int) {
	delete(s.trStaging, id)
}
