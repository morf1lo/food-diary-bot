package fooddiarybot

import "net/http"

type Server struct {
	httpServer *http.Server
}

func (s *Server) Start(port string) error {
	s.httpServer = &http.Server{
		Handler: s,
		Addr: ":" + port,
	}

	return s.httpServer.ListenAndServe()
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {}
