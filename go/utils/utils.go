package main

import (
	"log"
	"os"
	"time"
)

func WaitUntilFileEmpty(path string) bool {
	delay := 50 * time.Millisecond

	for {
		info, err := os.Stat(path)
		if err != nil {
			log.Printf("❌ File stat error: %v", err)
			time.Sleep(delay)
			continue
		}

		if info.Size() == 0 {
			log.Printf("✅ File is empty")
			return true
		}

		log.Printf("⏳ Waiting for file to be empty (current size: %d bytes)", info.Size())
		time.Sleep(delay)

		if delay < 200*time.Millisecond {
			delay += 10 * time.Millisecond
		}
	}
}
