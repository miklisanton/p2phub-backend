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
	trStaging map[int64]*models.Tracker
}

func NewTrackerService(repo *repository.TrackerRepository) *TrackerService {
	exchanges := make(map[string]bool)
	exchanges["binance"] = true
	exchanges["bybit"] = true
	return &TrackerService{repo: repo, Exchanges: exchanges, trStaging: make(map[int64]*models.Tracker)}
}

func (s *TrackerService) CreateTracker(tracker *models.Tracker) error {
	if tracker == nil {
		return fmt.Errorf("Tracker is nil")
	}

	tracker.Side = strings.ToUpper(tracker.Side)
	if tracker.Side != "BUY" && tracker.Side != "SELL" {
		return fmt.Errorf("Side must be BUY/SELL")
	}

	tracker.Currency = strings.ToUpper(tracker.Currency)
	if len(tracker.Currency) != 3 {
		return fmt.Errorf("Currency must be 3 symbols <EUR>")
	}

	tracker.Exchange = strings.ToLower(tracker.Exchange)
	if _, ok := s.Exchanges[tracker.Exchange]; !ok {
		return fmt.Errorf("exchange %s not supported", tracker.Exchange)
	}

	// Remove tracker from staging area
	s.DeleteTrackerStaging(tracker.ChatID)

	return s.repo.Save(tracker)
}

func (s *TrackerService) SetWaitingFlag(id int64, flag bool) error {
	return s.repo.UpdateWaitingFlag(id, flag)
}

func (s *TrackerService) GetAllTrackers() ([]*models.UserTracker, error) {
	return s.repo.GetAllTrackers()
}

func (s *TrackerService) GetTrackerStaging(id int64) *models.Tracker {
	tr, ok := s.trStaging[id]
	if !ok {
		tr := &models.Tracker{}
		s.trStaging[id] = tr
		return tr
	}
	return tr
}

func (s *TrackerService) DeleteTrackerStaging(id int64) {
	delete(s.trStaging, id)
}
