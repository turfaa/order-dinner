package service

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/turfaa/order-dinner/dinner"
)

type dinnerService struct {
	ctx context.Context

	client        dinner.Client
	interval      int
	visited       map[int64]bool
	startDateHash int64

	mutex *sync.Mutex
}

func NewDinnerService(ctx context.Context, client dinner.Client, interval int, startDate time.Time) (Service, error) {
	s := &dinnerService{
		ctx: ctx,

		client:        client,
		interval:      interval,
		visited:       make(map[int64]bool),
		startDateHash: getDateHash(startDate),

		mutex: &sync.Mutex{},
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
			go func() {
				s.mutex.Lock()
				defer s.mutex.Unlock()

				v := false
				now := time.Now().UTC()
				h, m, _ := now.Clock()
				hash := getDateHash(now)

				if val, ok := s.visited[hash]; ok {
					v = val
				} else {
					s.visited[hash] = false
				}

				if !v && (h > 4 || (h == 4 && m >= 30)) && s.client.IsReady() {
					if err := s.client.Order(); err == nil {
						s.visited[hash] = true
						log.Print("Order success")
					} else {
						log.Printf("Order error: %s", err)
					}
				} else {
					if err := s.client.UpdateMenu(); err == nil {
						log.Print("Update menu success")
					} else {
						log.Printf("Update menu error: %s", err)
					}
				}
			}()

		case <-s.ctx.Done():
			log.Print("Context closed")
			return s.ctx.Err()
		}
	}
}
