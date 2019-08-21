package service

import (
	"context"
	"github.com/turfaa/order-dinner/dinner"
	"log"
	"time"
)

type dinnerService struct {
	ctx           context.Context
	client        dinner.Client
	interval      int
	visited       map[int64]bool
	startDateHash int64
}

func NewDinnerService(ctx context.Context, client dinner.Client, interval int, startDate time.Time) (Service, error) {
	s := &dinnerService{
		ctx:           ctx,
		client:        client,
		interval:      interval,
		visited:       make(map[int64]bool),
		startDateHash: getDateHash(startDate),
	}

	for t := time.Now().UTC(); getDateHash(t) < s.startDateHash; t = t.AddDate(0, 0, 1) {
		s.visited[getDateHash(t)] = true
	}

	return s, nil
}

func (s *dinnerService) Serve() error {
	for {
		select {
		case <-time.Tick(time.Duration(s.interval) * time.Millisecond):
			v := false
			now := time.Now().UTC()
			h, m, _ := now.Clock()
			hash := getDateHash(now)

			if val, ok := s.visited[hash]; ok {
				v = val
			} else {
				s.visited[hash] = false
			}

			if !v && (h > 4 || (h == 4 && m >= 30)) {
				if err := s.client.Order(); err == nil {
					s.visited[hash] = true
					log.Print("Order success")
				} else {
					log.Printf("Order error: %s", err)
				}
			} else {
				if err := s.client.HealthCheck(); err == nil {
					log.Print("Health check success")
				} else {
					log.Printf("Health check error: %s", err)
				}

			}

		case <-s.ctx.Done():
			log.Print("Context closed")
			return s.ctx.Err()
		}
	}
}
