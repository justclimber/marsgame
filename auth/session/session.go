package session

import (
	"crypto/rand"
	"encoding/base64"
	"io"
	"sync"
	"time"
)

const defaultTTL = time.Hour * 24 * 14
const defaultCookieName = "sid"

type Data struct {
	UserId    uint32
	updatedAt time.Time
}

func (d *Data) isExpired(at time.Time) bool {
	return d.updatedAt.Before(at)
}

type Session struct {
	sync.Mutex
	data       map[string]*Data
	Ttl        time.Duration
	CookieName string
}

func NewSession() *Session {
	return &Session{
		data:       make(map[string]*Data),
		Ttl:        defaultTTL,
		CookieName: defaultCookieName,
	}
}

func (s *Session) Start(userId uint32) string {
	sessionId := generateSessionId()
	s.Lock()
	defer s.Unlock()

	s.data[sessionId] = &Data{UserId: userId, updatedAt: time.Now()}
	return sessionId
}

func (s *Session) Get(sessionId string) (*Data, bool) {
	data, ok := s.data[sessionId]
	if data.isExpired(time.Now().Add(-s.Ttl)) {
		delete(s.data, sessionId)
		return nil, false
	}
	return data, ok
}

func generateSessionId() string {
	randbytes := make([]byte, 40)
	if _, err := io.ReadFull(rand.Reader, randbytes); err != nil {
		return ""
	}
	return base64.StdEncoding.EncodeToString(randbytes)
}
