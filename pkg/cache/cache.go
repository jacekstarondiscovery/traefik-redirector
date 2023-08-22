package cache

import (
	"github.com/jacekstarondiscovery/traefik-redirector/pkg/model"
	"sync"
	"time"
)

type TimeNow func() time.Time

type RedirectCache struct {
	Redirects  []model.Redirect `json:"redirects"`
	UpdateDate time.Time        `json:"updateDate"`
	ExpireDate time.Time        `json:"expireDate"`
	MaxAge     time.Duration    `json:"maxAge"`
	mutex      sync.Mutex       `json:"-"`
	timeNow    TimeNow          `json:"-"`
}

func New(maxAge int64, timeNow TimeNow) *RedirectCache {
	return &RedirectCache{
		ExpireDate: time.Time{},
		UpdateDate: time.Time{},
		Redirects:  []model.Redirect{},
		MaxAge:     time.Duration(maxAge),
		mutex:      sync.Mutex{},
		timeNow:    timeNow,
	}
}

func (rc *RedirectCache) IsFresh() bool {
	if rc.ExpireDate.Equal(time.Time{}) {
		return false
	}

	now := rc.timeNow()
	return now.Before(rc.ExpireDate)
}

func (rc *RedirectCache) Update(redirects []model.Redirect) {
	now := rc.timeNow()

	rc.Redirects = redirects
	rc.UpdateDate = now
	rc.ExpireDate = now.Add(rc.MaxAge * time.Second)
}

func (rc *RedirectCache) Lock() bool {
	return rc.mutex.TryLock()
}

func (rc *RedirectCache) Unlock() {
	rc.mutex.Unlock()
}
