package main

import (
	"context"
	"log"
	"net"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/soheilhy/cmux"
	pb "github.com/vmotta8/go-cmux/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/types/known/emptypb"
)

type grpcSrv struct {
	pb.UnimplementedInfoServiceServer
}

func (s *grpcSrv) GetInfoGrpc(ctx context.Context, in *emptypb.Empty) (*pb.InfoResponse, error) {
	return &pb.InfoResponse{Message: "Olá do gRPC!"}, nil
}

func main() {

	lis, err := net.Listen("tcp", ":8081")
	if err != nil {
		log.Fatalf("Falha ao escutar na porta 8081: %v", err)
	}

	m := cmux.New(lis)

	grpcL := m.Match(cmux.HTTP2HeaderField("content-type", "application/grpc"))
	httpL := m.Match(cmux.Any())

	grpcServer := grpc.NewServer()
	pb.RegisterInfoServiceServer(grpcServer, &grpcSrv{})

	reflection.Register(grpcServer)

	gin.SetMode(gin.ReleaseMode)
	router := gin.New()

	router.GET("/get-info-rest", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Olá do REST!",
		})
	})

	go func() {
		if err := grpcServer.Serve(grpcL); err != nil {
			log.Fatalf("Falha ao servir gRPC: %v", err)
		}
	}()

	go func() {
		if err := http.Serve(httpL, router); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Falha ao servir HTTP: %v", err)
		}
	}()

	log.Println("Servidor escutando na porta 8081")
	if err := m.Serve(); err != nil {
		log.Fatalf("cmux falhou: %v", err)
	}
}
