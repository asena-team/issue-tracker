package main

import (
	"golang.org/x/time/rate"
	"sync"
	"time"
)

type Visitor struct {
	Limiter  *rate.Limiter
	LastSeen time.Time
}

var visitors = make(map[string]*Visitor)
var mutex sync.Mutex

func GetVisitor(ip string) *rate.Limiter {
	mutex.Lock()
	defer mutex.Unlock()

	v, exists := visitors[ip]
	if !exists {
		r := rate.Every(time.Minute * 30)
		limiter := rate.NewLimiter(r, 2)
		visitors[ip] = &Visitor{
			limiter,
			time.Now(),
		}

		return limiter
	}

	v.LastSeen = time.Now()
	return v.Limiter
}

func CleanupVisitors() {
	for {
		time.Sleep(3 * time.Second)

		mutex.Lock()
		for ip, v := range visitors {
			if time.Since(v.LastSeen) > 30*time.Minute {
				delete(visitors, ip)
			}
		}
		mutex.Unlock()
	}
}
