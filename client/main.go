package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"

	pb "github.com/luisteixeira74/grpc-file-stream-to-s3/proto/uploader" // Importa o pacote gerado pelo protoc
	"google.golang.org/grpc"
)

const (
	serverAddress = "localhost:50051" // Endereço do servidor gRPC
)

func main() {
	// Conectando-se ao servidor gRPC
	conn, err := grpc.Dial(serverAddress, grpc.WithInsecure()) // Usa uma conexão sem TLS
	if err != nil {
		log.Fatalf("Erro ao conectar ao servidor: %v", err)
	}
	defer conn.Close()

	client := pb.NewFileUploaderClient(conn)

	// Abre o arquivo que será enviado
	file, err := os.Open("client/files/teste.jpg") // Substitua pelo caminho do seu arquivo
	if err != nil {
		log.Fatalf("Erro ao abrir o arquivo: %v", err)
	}
	defer file.Close()

	// Criação do stream de upload
	stream, err := client.Upload(context.Background())
	if err != nil {
		log.Fatalf("Erro ao iniciar o upload: %v", err)
	}

	// Envia o arquivo em chunks
	buffer := make([]byte, 1024) // Buffer de 1KB
	for {
		n, err := file.Read(buffer)
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("Erro ao ler o arquivo: %v", err)
		}

		// Envia o chunk para o servidor
		req := &pb.UploadRequest{
			Filename: "teste.jpg", // Nome do arquivo (pode ser alterado conforme necessário)
			Chunk:    buffer[:n],
		}

		if err := stream.Send(req); err != nil {
			log.Fatalf("Erro ao enviar chunk: %v", err)
		}
		fmt.Printf("Enviado chunk de %d bytes\n", n)
	}

	// Finaliza o upload e aguarda a resposta do servidor
	resp, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("Erro ao receber resposta do servidor: %v", err)
	}

	// Exibe a resposta do servidor
	fmt.Printf("Resposta do servidor: %s (sucesso: %v)\n", resp.GetMessage(), resp.GetSuccess())
}
