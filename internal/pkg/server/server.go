package server

import (
	. "github.com/dinowar/gateway-service/internal/pkg/domain"
	"github.com/dinowar/gateway-service/internal/pkg/service"
	"log"
	"net/http"
)

var ()

type Server struct {
	rep      *service.RepositoryService
	logger   *service.LogService
	gateways map[string]PaymentGateway
}

func NewAppServer(rep *service.RepositoryService, logger *service.LogService, workerCount int) *Server {
	return &Server{
		rep:      rep,
		logger:   logger,
		gateways: make(map[string]PaymentGateway),
	}
}

func (server *Server) RegisterGateway(name string, gateway PaymentGateway) {
	server.gateways[name] = gateway
	log.Printf("registered gateway: %s", name)
}

func (server *Server) HandleDeposit(w http.ResponseWriter, r *http.Request) {

}

func (server *Server) HandleWithdraw(w http.ResponseWriter, r *http.Request) {

}

func (server *Server) HandleCallback(w http.ResponseWriter, r *http.Request) {

}

func (server *Server) HandleGetTransaction(w http.ResponseWriter, r *http.Request) {

}
