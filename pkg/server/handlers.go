package server

import (
	"context"
	"github.com/whereiskurt/tiogo/pkg/ui"
	"net/http"
	"time"
)

func (s *Server) Shutdown(w http.ResponseWriter, r *http.Request) {
	s.Log.Debugf("/Shutdown called - beginning s Shutdown")

	_, _ = w.Write([]byte(ui.Gopher()))
	_, _ = w.Write([]byte("\n...bye felicia\n"))

	timeout, cancel := context.WithTimeout(s.Context, 5*time.Second)
	err := s.HTTP.Shutdown(timeout)
	if err != nil {
		s.Log.Errorf("server error during Shutdown: %+v", err)
	}
	s.Finished()
	cancel()
}
