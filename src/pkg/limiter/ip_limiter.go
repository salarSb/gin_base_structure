package limiter

import (
	"golang.org/x/time/rate"
	"sync"
)

type IpRateLimiter struct {
	ips map[string]*rate.Limiter
	mu  *sync.RWMutex
	r   rate.Limit
	b   int
}

func NewIpRateLimiter(r rate.Limit, b int) *IpRateLimiter {
	return &IpRateLimiter{
		ips: make(map[string]*rate.Limiter),
		mu:  &sync.RWMutex{},
		r:   r,
		b:   b,
	}
}

func (l *IpRateLimiter) AddIp(ip string) *rate.Limiter {
	l.mu.Lock()
	defer l.mu.Unlock()
	limiter := rate.NewLimiter(l.r, l.b)
	l.ips[ip] = limiter
	return limiter
}

func (l *IpRateLimiter) GetLimiter(ip string) *rate.Limiter {
	l.mu.Lock()
	limiter, exists := l.ips[ip]
	if !exists {
		l.mu.Unlock()
		return l.AddIp(ip)
	}
	l.mu.Unlock()
	return limiter
}
