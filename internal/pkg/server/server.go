package server

import (
	"github.com/dinowar/gateway-service/internal/pkg/service"
	"net/http"
)

type Server struct {
	rep    *service.RepositoryService
	logger *service.LogService
}

func NewAppServer(rep *service.RepositoryService, logger *service.LogService) *Server {
	return &Server{
		rep:    rep,
		logger: logger,
	}
}

func (server *Server) HandleDeposit(w http.ResponseWriter, r *http.Request) {

}

func (server *Server) HandleWithdraw(w http.ResponseWriter, r *http.Request) {

}

func (server *Server) HandleCallback(w http.ResponseWriter, r *http.Request) {

}

func (server *Server) HandleGetTransaction(w http.ResponseWriter, r *http.Request) {

}
