# Guia de Uso do gRPC no Projeto `grpc-file-stream-to-s3`

## 📌 Por que usar gRPC?

gRPC é um framework de comunicação de alto desempenho criado pela Google, baseado em HTTP/2 e Protocol Buffers. Ele permite:
- Comunicação eficiente entre serviços (microservices)
- Streaming bidirecional
- Tipagem forte e contrato bem definido
- Geração automática de código para várias linguagens

No nosso projeto, gRPC é utilizado para **fazer upload de arquivos** via *streaming*, com foco em performance, padronização e escalabilidade.

---

## ⚙️ O que o gRPC faz aqui?

Neste projeto:
- O **client** divide um arquivo em partes (chunks) e envia ao servidor gRPC.
- O **server** recebe cada chunk e armazena o arquivo final diretamente no **Amazon S3**.
- A comunicação entre cliente e servidor é feita com **Protocol Buffers (.proto)**, garantindo performance e consistência entre as partes.

---

## 🔄 Fluxo da aplicação

1. O client se conecta ao servidor gRPC.
2. Cria-se um *stream* com o método `Upload`.
3. O client lê um arquivo local em partes e envia cada pedaço (`chunk`) por esse stream.
4. O server recebe os chunks, reconstrói o arquivo e faz o upload para o S3.
5. Após o envio completo, o server retorna uma resposta final com sucesso ou erro.

---

## 🧱 Estrutura do gRPC no projeto

### 1. Arquivo `.proto`
Define:
- Serviço (`service FileUploader`)
- Métodos (`rpc Upload`)
- Tipos de mensagem (`UploadRequest`, `UploadResponse`)

Exemplo:
```proto
syntax = "proto3";

package uploader;

option go_package = "proto/uploader;uploader";

service FileUploader {
  rpc Upload (stream UploadRequest) returns (UploadResponse);
}

message UploadRequest {
  string filename = 1;
  bytes chunk = 2;
}

message UploadResponse {
  string message = 1;
  bool success = 2;
}
```

## Gerando o código .pb.go automaticamente
protoc --go_out=. --go-grpc_out=. proto/uploader/upload.proto

Isso gera os arquivos .pb.go contendo:

- Interfaces
- Tipos de mensagem

Clientes e servidores prontos para uso

3. Implementação do Servidor
Contém a lógica de:

Receber o stream

Montar o arquivo

Salvar no S3

4. Implementação do Cliente
Responsável por:

Ler o arquivo local

Enviar os chunks via stream

Esperar resposta final

🟢 Vantagens práticas do gRPC

🚀 Performance	Usa HTTP/2 + ProtoBuf, mais leve e rápido que JSON/REST
📦 Contract-first	Garante que cliente e servidor usem o mesmo contrato (definido via .proto)
🔄 Streaming	Suporte nativo a comunicação contínua (ideal para upload de arquivos)

🧩 Quando usar gRPC

- Alta performance e baixo overhead são necessários
- Streaming de dados (ex: arquivos, vídeos, etc)
- Integração entre serviços internos (microservices)
- APIs internas onde controle e eficiência importam mais que compatibilidade com navegadores


🔄 Fluxo de uma aplicação gRPC como esta
🧭 Visão geral do fluxo cliente → servidor → S3

Client                gRPC Server               AWS S3
  |                        |                      |
  | Open file              |                      |
  |----------------------->|                      |
  | Send chunks            |                      |
  |----------------------->|                      |
  |                        | Upload to S3         |
  |                        |--------------------->|
  |                        | Response (success)   |
  |<-----------------------|                      |


  1. O Client:
Abre o arquivo

Conecta-se ao servidor gRPC

Envia o conteúdo do arquivo em chunks via stream

Espera uma resposta com status final

2. O Server:
Recebe os chunks do client

Agrupa e escreve o conteúdo no S3

Retorna uma resposta UploadResponse com status/sucesso

🧱 Estrutura de um projeto gRPC
1. Arquivo .proto
Define o contrato da API:

syntax = "proto3";

package uploader;

option go_package = "proto/uploader;uploader";

service FileUploader {
  rpc Upload (stream UploadRequest) returns (UploadResponse);
}

message UploadRequest {
  string filename = 1;
  bytes chunk = 2;
}

message UploadResponse {
  string message = 1;
  bool success = 2;
}

Esse arquivo define:

O nome do serviço: FileUploader

O método: Upload, que recebe um stream de UploadRequest e retorna um UploadResponse

As mensagens envolvidas

2. Código Gerado (pb.go)
Com protoc, você gera:

upload.pb.go: tipos das mensagens

upload_grpc.pb.go: interface gRPC com métodos Upload, NewFileUploaderClient, etc

3. Servidor (server/main.go)
Implementa a interface FileUploaderServer

Recebe os chunks do cliente

Envia para o bucket S3 (ou armazena no disco)

4. Cliente (client/main.go)
Cria conexão gRPC

Usa Upload() para iniciar o stream

Envia pedaços do arquivo

Fecha e espera o resultado
