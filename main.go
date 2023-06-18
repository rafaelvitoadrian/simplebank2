package main

import (
	"database/sql"
	"log"
	"net"

	_ "github.com/lib/pq"
	"github.com/rafaelvitoadrian/simplebank2/api"
	db "github.com/rafaelvitoadrian/simplebank2/db/sqlc"
	"github.com/rafaelvitoadrian/simplebank2/gapi"
	"github.com/rafaelvitoadrian/simplebank2/pb"
	"github.com/rafaelvitoadrian/simplebank2/utils"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	config, err := utils.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load config: ", err)
	}

	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("cannot connect db: ", err)
	}

	store := db.NewStore(conn)
	runGRPCServer(config, store)
}

func runGRPCServer(config utils.Config, store db.Store) {
	server, err := gapi.NewServer(config, store)
	if err != nil {
		log.Fatal("cannot create server : ", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterSimpleBankServer(grpcServer, server)
	reflection.Register(grpcServer)

	listener, err := net.Listen("tcp", config.GRPCServerAddress)
	if err != nil {
		log.Fatal("cannot create listener : ", err)
	}

	log.Printf("Start gRPC Server at %s", listener.Addr().String())

	err = grpcServer.Serve(listener)
	if err != nil {
		log.Fatal("Cannot crreate gRPC server : ", err)
	}

}

func runGinServer(config utils.Config, store db.Store) {
	server, err := api.NewServer(config, store)
	if err != nil {
		log.Fatal(err)
	}

	err = server.Start(config.HTTPServerAddress)
	if err != nil {
		log.Fatal(err)
	}
}
