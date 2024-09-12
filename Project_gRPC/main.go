package main

import (
	"Project_gRPC/protogen/golang/golang/orders"
	"fmt"
	"log"
	"net"
	"time"

	"google.golang.org/grpc"
)

// server es la estructura que implementa el servicio OrderService
type server struct {
	orders.UnimplementedOrderServiceServer
}

// GetOrders implementa el método GetOrders del servicio OrderService
func (s *server) GetOrders(req *orders.PayloadWithSingleOrder, stream orders.OrderService_GetOrdersServer) error {
	// Para la demostración, vamos a enviar varios mensajes de Order con un retraso entre ellos
	for i := 0; i < 5; i++ {
		// Crear una copia del pedido con valores actualizados para cada mensaje
		msg := &orders.Order{
			Id:   uint64(i + 1),
			Name: fmt.Sprintf("Item %d", i+1),
		}

		// Envía el mensaje al cliente en formato Protobuf
		if err := stream.Send(msg); err != nil {
			log.Println("Error al enviar el mensaje:", err)
			return err
		}

		// Espera un momento antes de enviar el siguiente mensaje
		time.Sleep(3 * time.Second)
	}

	return nil
}

func main() {
	// Configura el listener en el puerto 50051
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Error al intentar escuchar en el puerto: %v", err)
	}

	// Crea un nuevo servidor gRPC
	grpcServer := grpc.NewServer()

	// Registra el servicio OrderService con el servidor
	orders.RegisterOrderServiceServer(grpcServer, &server{})

	log.Println("Servidor gRPC escuchando en el puerto 50051...")
	// Inicia el servidor gRPC
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Error al iniciar el servidor gRPC: %v", err)
	}
}
