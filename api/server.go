package api

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	db "github.com/rafaelvitoadrian/simplebank2/db/sqlc"
	"github.com/rafaelvitoadrian/simplebank2/token"
	"github.com/rafaelvitoadrian/simplebank2/utils"
)

type Server struct {
	config     utils.Config
	store      db.Store
	router     *gin.Engine
	tokenMaker token.Maker
}

func NewServer(config utils.Config, store db.Store) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(config.TokenSymetricKey)
	if err != nil {
		return nil, fmt.Errorf("Cannot Create Token maker : %w", err)
	}

	server := &Server{
		config:     config,
		store:      store,
		tokenMaker: tokenMaker,
	}
	// router := gin.Default()

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("currency", validCurrency)
	}

	server.setupRouter()

	return server, nil
}

func (server *Server) setupRouter() {
	router := gin.Default()
	router.POST("/accounts", server.createAccount)
	router.GET("/accounts/:id", server.getAccount)
	router.GET("/accounts", server.listAccount)

	router.POST("/transfer", server.createTransfer)

	router.POST("/user", server.createUser)
	router.POST("/user/login", server.loginUser)

	router.POST("/token/renew_access", server.renewAccessToken)
	server.router = router
}

func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
