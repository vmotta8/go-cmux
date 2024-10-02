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
	"google.golang.org/protobuf/types/known/emptypb" // Import necessário
)

// Renomeie a estrutura para evitar conflitos
type grpcSrv struct {
	pb.UnimplementedInfoServiceServer
}

// Atualize a assinatura para usar *emptypb.Empty
func (s *grpcSrv) GetInfoGrpc(ctx context.Context, in *emptypb.Empty) (*pb.InfoResponse, error) {
	return &pb.InfoResponse{Message: "Olá do gRPC!"}, nil
}

func main() {
	// Cria o listener principal na porta 8081
	lis, err := net.Listen("tcp", ":8081")
	if err != nil {
		log.Fatalf("Falha ao escutar na porta 8081: %v", err)
	}

	// Cria um cmux a partir do listener
	m := cmux.New(lis)

	// Matchers para gRPC e HTTP
	grpcL := m.Match(cmux.HTTP2HeaderField("content-type", "application/grpc"))
	httpL := m.Match(cmux.Any())

	// Configura o servidor gRPC
	grpcServer := grpc.NewServer()
	pb.RegisterInfoServiceServer(grpcServer, &grpcSrv{}) // Use a nova estrutura
	// Registrar reflection para facilitar testes
	reflection.Register(grpcServer)

	// Configura o servidor HTTP com Gin
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()

	// Endpoint REST em /get-info-rest
	router.GET("/get-info-rest", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Olá do REST!",
		})
	})

	// Inicia o servidor gRPC em uma goroutine
	go func() {
		if err := grpcServer.Serve(grpcL); err != nil {
			log.Fatalf("Falha ao servir gRPC: %v", err)
		}
	}()

	// Inicia o servidor HTTP em outra goroutine
	go func() {
		if err := http.Serve(httpL, router); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Falha ao servir HTTP: %v", err)
		}
	}()

	// Inicia o cmux para servir ambos os servidores
	log.Println("Servidor escutando na porta 8081")
	if err := m.Serve(); err != nil {
		log.Fatalf("cmux falhou: %v", err)
	}
}
