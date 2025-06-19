package main

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"net"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	pb "github.com/luisteixeira74/grpc-file-stream-to-s3/proto/uploader"
	"google.golang.org/grpc"
)

const (
	port   = ":50051"
	bucket = "seu-bucket-s3" // Substitua pelo nome do seu bucket
	region = "us-east-1"     // Ou a região que você escolheu (ex: sa-east-1)
)

type server struct {
	pb.UnimplementedFileUploaderServer
	s3Client *s3.Client
}

func (s *server) Upload(stream pb.FileUploader_UploadServer) error {
	var buffer bytes.Buffer
	var filename string

	for {
		req, err := stream.Recv()
		if err == io.EOF {
			break // fim do envio
		}
		if err != nil {
			return fmt.Errorf("erro ao receber chunk: %v", err)
		}

		filename = req.GetFilename()
		_, err = buffer.Write(req.GetChunk())
		if err != nil {
			return fmt.Errorf("erro ao gravar chunk: %v", err)
		}
	}

	// Envia para o S3
	_, err := s.s3Client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket: &bucket,
		Key:    &filename,
		Body:   bytes.NewReader(buffer.Bytes()),
	})
	if err != nil {
		return stream.SendAndClose(&pb.UploadResponse{
			Success: false,
			Message: fmt.Sprintf("erro ao enviar para o S3: %v", err),
		})
	}

	return stream.SendAndClose(&pb.UploadResponse{
		Success: true,
		Message: "Arquivo enviado com sucesso para o S3",
	})
}

func main() {
	// Carrega config da AWS
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(region))
	if err != nil {
		log.Fatalf("Erro ao carregar config AWS: %v", err)
	}

	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("Erro ao ouvir na porta %s: %v", port, err)
	}

	s := grpc.NewServer()
	pb.RegisterFileUploaderServer(s, &server{
		s3Client: s3.NewFromConfig(cfg),
	})

	log.Printf("Servidor gRPC ouvindo em %s...\n", port)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Erro ao iniciar servidor: %v", err)
	}
}
