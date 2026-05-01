# OpenTelemetry Go Collector & Instrumentation

Este projeto é um módulo Go projetado para padronizar e facilitar a implementação de observabilidade (Traces, Metrics e Logs) em serviços Go, seguindo os princípios de **12-Factor App** e utilizando o padrão **OpenTelemetry (OTel)**.

## 🚀 Funcionalidades

- **Tracing**: Inicialização automática de TracerProvider com suporte a exportação via OTLP (gRPC/HTTP).
- **Metrics**: Implementação de MeterProvider para coleta de métricas de performance e saúde do sistema.
- **Logging Estruturado**: Integração com `logrus` utilizando um formatador JSON rígido que inclui automaticamente `trace_id` e `span_id`.
- **Middlewares**: Suporte nativo para **Gin**, **Fiber**, **gRPC** e `net/http` standard.
- **Configuração via Env**: Totalmente configurável através de variáveis de ambiente.

---

## 📦 Instalação

```bash
go get github.com/alanfranciscos/otel-collector
```

---

## ⚙️ Configuração (Variáveis de Ambiente)

O módulo utiliza as seguintes variáveis para se auto-configurar:

| Variável | Descrição | Exemplo/Padrão |
| :--- | :--- | :--- |
| `ENVIRONMENT` | Ambiente de execução (`PRODUCTION`, `STAGING`, `LOCAL`) | `LOCAL` |
| `OTEL_SERVICE_NAME` | Nome do serviço para identificação nos traces e logs | `my-api-service` |
| `OTEL_EXPORTER_OTLP_PROTOCOL` | Protocolo de exportação OTLP | `http` ou `grpc` |
| `OTEL_EXPORTER_OTLP_ENDPOINT` | Endpoint do OTel Collector ou Backend | `localhost:4317` |

> **Nota**: Em modo `LOCAL`, os exportadores OTLP são desativados para evitar erros de conexão, mas a geração de IDs e o log estruturado continuam funcionando normalmente.

---

## 🛠️ Como Utilizar

### 1. Inicialização

Chame o `Initialize` o mais cedo possível no seu `main.go`.

```go
func main() {
    ctx := context.Background()
    serviceName := "MY-APP"

    // Inicializa Logs, Traces e Metrics
    shutdown, err := telemetry.NewTelemetry(ctx, &serviceName).Initialize()
    if err != nil {
        log.Fatal(err)
    }
    defer shutdown(ctx) // Garante o flush dos dados antes de sair
}
```

### 2. Traces (Rastreamento)

Utilize o `tracer` global para criar spans manuais em operações críticas.

```go
var tracer = otel.Tracer("my-handler")

func processData(ctx context.Context) {
    ctx, span := tracer.Start(ctx, "process.data.operation")
    defer span.End()

    // ... lógica de negócio
}
```

### 3. Metrics (Métricas)

Crie contadores ou histogramas para monitorar o comportamento da aplicação.

```go
var meter = otel.Meter("my-service-meter")
var counter, _ = meter.Int64Counter("orders_processed_total")

func handleOrder(ctx context.Context) {
    counter.Add(ctx, 1)
}
```

### 4. Logging Estruturado

Use sempre o `logrus.WithContext(ctx)` para garantir que o `trace_id` atual seja injetado no log. **Importante**: No Gin, prefira `ctx.Request.Context()`.

```go
logrus.WithContext(ctx).Info("Operação realizada com sucesso")
```

### 5. Middlewares (Exemplo Gin)

```go
app := gin.New()

// Adiciona middleware que gerencia Spans e Logs de requisição automaticamente
ginMdw := ginmdw.NewGinMiddlewareConfig()
app.Use(ginMdw.Middleware()...)
```

---

## 🏗️ Estrutura do Projeto

- `/pkg/telemetry`: Ponto de entrada público para inicialização.
- `/middleware`: Implementações de middlewares para diversos frameworks.
- `/internal/pkg/telemetry/provider`: Lógica interna de configuração de cada sinal (Log, Trace, Meter).
- `/internal/pkg/telemetry/schema`: Definição do JSON rígido utilizado nos logs.

---

## 🧪 Desenvolvimento

Utilize o `Makefile` para tarefas comuns:

- `make test`: Executa os testes unitários.
- `make lint`: Verifica estilo e erros estáticos.
- `make run-example`: Executa a API de exemplo em `examples/gin`.
- `make tidy`: Organiza as dependências do módulo e exemplos.
