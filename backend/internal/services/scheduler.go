package services

import (
	"log"
	"time"
)

func StartRateSync(service *RateService) {
	go func() {
		// Run once immediately
		if err := service.SyncRates("USD"); err != nil {
			log.Println("Initial sync failed:", err)
		} else {
			log.Println("Initial rate sync completed")
		}

		ticker := time.NewTicker(1 * time.Hour)
		defer ticker.Stop()

		for range ticker.C {
			if err := service.SyncRates("USD"); err != nil {
				log.Println("Scheduled sync failed:", err)
			} else {
				log.Println("Rates synchronized successfully")
			}
		}
	}()
}
