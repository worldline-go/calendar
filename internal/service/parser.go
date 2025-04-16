package service

import (
	"context"

	"github.com/rs/zerolog/log"

	"github.com/worldline-go/calendar/pkg/ics"
)

func (s *Service) getRRule(ctx context.Context, rruleStr string) (*ics.RRule, error) {
	s.m.RLock()
	rrule, ok, err := s.cache.Get(ctx, rruleStr)
	s.m.RUnlock()
	if err != nil {
		log.Error().Err(err).Msg("failed to get rrule from cache")
	}

	if ok {
		return rrule, nil
	}

	// Only one goroutine should parse and set for a given rruleStr at a time
	s.m.Lock()
	defer s.m.Unlock()
	// Double-check cache after acquiring write lock (to avoid race)
	rrule, ok, err = s.cache.Get(ctx, rruleStr)
	if err != nil {
		log.Error().Err(err).Msg("failed to get rrule from cache (after lock)")
	}
	if ok {
		return rrule, nil
	}

	rrule, err = ics.ParseRRule(rruleStr)
	if err != nil {
		return nil, err
	}

	if err := s.cache.Set(ctx, rruleStr, rrule); err != nil {
		log.Error().Err(err).Msg("failed to set rrule in cache")
	}

	return rrule, nil
}
