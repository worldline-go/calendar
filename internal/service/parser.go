package service

import (
	"context"
	"time"

	"github.com/rs/zerolog/log"

	"github.com/worldline-go/calendar/pkg/ical"
)

func (s *Service) getRRule(ctx context.Context, repeatStr string) (*ical.Repeat, error) {
	s.m.RLock()
	rrule, ok, err := s.cacheRule.Get(ctx, repeatStr)
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
	rrule, ok, err = s.cacheRule.Get(ctx, repeatStr)
	if err != nil {
		log.Error().Err(err).Msg("failed to get rrule from cache (after lock)")
	}
	if ok {
		return rrule, nil
	}

	rrule, err = ical.ParseRepeat(repeatStr)
	if err != nil {
		return nil, err
	}

	if err := s.cacheRule.Set(ctx, repeatStr, rrule); err != nil {
		log.Error().Err(err).Msg("failed to set rrule in cache")
	}

	return rrule, nil
}

func (s *Service) TZLocation(tz string) (*time.Location, error) {
	s.m.RLock()
	loc, ok, err := s.cacheTZ.Get(context.Background(), tz)
	s.m.RUnlock()
	if err != nil {
		log.Error().Err(err).Msg("failed to get location from cache")
	}

	if ok {
		return loc, nil
	}

	// Only one goroutine should parse and set for a given tz at a time
	s.m.Lock()
	defer s.m.Unlock()
	// Double-check cache after acquiring write lock (to avoid race)
	loc, ok, err = s.cacheTZ.Get(context.Background(), tz)
	if err != nil {
		log.Error().Err(err).Msg("failed to get location from cache (after lock)")
	}
	if ok {
		return loc, nil
	}

	loc, err = time.LoadLocation(tz)
	if err != nil {
		return nil, err
	}
	if err := s.cacheTZ.Set(context.Background(), tz, loc); err != nil {
		log.Error().Err(err).Msg("failed to set location in cache")
	}

	return loc, nil
}
