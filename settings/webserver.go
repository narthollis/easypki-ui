package settings

import (
	"os"
	"time"
)

type WebServerSettings struct {
	Address string
	// "the duration for which the server gracefully wait for existing connections to finish - e.g. 15s or 1m"
	GracefulTimeout time.Duration
}

func (s *WebServerSettings) Create() {
	s.Address = os.Getenv("HTTP_LISTEN")
	s.GracefulTimeout = time.Second * 15
}
