# gRPC File Upload (Client Streaming) to S3 using Go

Projeto em Go que realiza upload de arquivos para a AWS S3 utilizando gRPC com streaming de arquivos.

## 📦 Requisitos

- Go 1.20 ou superior
- `protoc` (Protocol Buffers compiler)
- Plugins do `protoc` para Go:
  - `protoc-gen-go`
  - `protoc-gen-go-grpc`
- Conta e bucket S3 configurados
- Variáveis de ambiente AWS (`AWS_ACCESS_KEY_ID`, `AWS_SECRET_ACCESS_KEY`, `AWS_REGION`)

## 🧱 Instalação dos plugins `protoc` (caso ainda não tenha)

```bash
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

grpc-file-stream-to-s3/
├── proto/
│   └── uploader/
│       └── upload.proto
├── server/
│   └── main.go
├── client/
│   └── main.go
├── go.mod

Certifique-se de que o .proto possui a diretiva correta:

option go_package = "proto/uploader;uploader";
Em seguida, execute o comando:

protoc --go_out=. --go-grpc_out=. proto/uploader/upload.proto
Isso vai gerar os arquivos:

proto/uploader/upload.pb.go
proto/uploader/upload_grpc.pb.go

▶️ Executando o servidor
go run server/main.go

📤 Executando o client (para enviar arquivos)

Edite o caminho do arquivo no client (caso necessário) e execute:
go run client/main.go

🐘 Conectando à AWS S3
Defina as variáveis de ambiente (pode usar um .env se usar algo como godotenv):

export AWS_ACCESS_KEY_ID="sua-chave"
export AWS_SECRET_ACCESS_KEY="sua-chave-secreta"
export AWS_REGION="sa-east-1"

documentação do projeto em:
docs/grpc-guia-projeto.md