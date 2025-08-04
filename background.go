package main

import (
	"context"
	"log"
	"time"
)

func StartURLCleaner(appConfig *AppConfig, ctx context.Context) {
	ticker := time.NewTicker(time.Minute)
	log.Printf("Cleaner set to delete expired records every %v \n", time.Minute)

	for {
		select {
		case <-ticker.C:
			err := appConfig.DB.DeleteExpiredURLEntries(ctx, time.Now().UTC())
			if err != nil {
				log.Println("Failed clearing db, err:", err)
			}

			log.Println("Cleaned expired urls, time:", time.Now())
		case <-ctx.Done():
			ticker.Stop()
			return
		}
	}
}
