package utils

import (
	"math/rand"
	"time"
)

func SleepRandomTime(min, max int) {
	rand.Seed(time.Now().UnixNano())
	duration := rand.Intn(max-min+1) + min
	time.Sleep(time.Duration(duration) * time.Second)
}
