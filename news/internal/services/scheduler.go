package services

import (
	"log"
	"time"
)

func StartNewsSync(service *NewsService) {
	go func() {
		service.SyncAll()

		ticker := time.NewTicker(1 * time.Hour)
		defer ticker.Stop()

		for range ticker.C {
			service.SyncAll()
		}
	}()
	log.Println("News sync scheduler started")
}
