package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/watchlist-kata/auth/cmd/config"
	crl "github.com/watchlist-kata/auth/internal/controller/grpc"
	pb "github.com/watchlist-kata/auth/internal/controller/grpc/proto"
	repo "github.com/watchlist-kata/auth/internal/repository/postgres"
	sv "github.com/watchlist-kata/auth/internal/service"
	"github.com/watchlist-kata/auth/pkg/tokenJWT"

	"google.golang.org/grpc"
)

func main() {
	cf := config.LoadConfig()
	jwt := tokenJWT.NewTokenJWT(cf.SecretKey, time.Minute*15)

	db := config.NewPostgres().InitDB(cf)
	defer db.Close()

	serverGRPC := grpc.NewServer()

	pb.RegisterAuthServiceServer(serverGRPC,
		crl.NewAuthServiceGRPC(sv.NewAuthSeviceImpl(jwt,
			repo.NewPostgesRepositoryImpl(db.GetDB()))))

	l, err := net.Listen("tcp", fmt.Sprintf(":%v", cf.GRPCServerPort))
	defer l.Close()

	if err != nil {
		log.Fatalf("listen error %w", err)
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		log.Println(fmt.Sprintf("The auth server gRPC is running on port %v", cf.GRPCServerPort))
		if err = serverGRPC.Serve(l); err != nil {
			log.Fatalf("Accept error %w", err)
		}
	}()

	<-stop
	serverGRPC.GracefulStop()
	log.Println("Server stopped")

}
