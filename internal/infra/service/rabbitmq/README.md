# Service / RabbitMQ

Wrapper para conexão/canais RabbitMQ usado por email, checkout, notificações. Fornece reconexão automática.

---

## 1. Responsabilidades
- Abrir conexão única por processo.
- Criar canais publicados/consumidores.
- Expor métodos para declarar filas e ligar exchanges.

## 2. Métodos principais
| Assinatura | Descrição |
|------------|-----------|
| `Publish(queue string, body []byte) error` | Envia mensagem com confirmação. |
| `Consume(queue string, handler func(Message))` | Registra consumidor com QoS configurável. |
| `Close()` | Fecha conexão/canais. |

## 3. Fluxo típico
- Módulos chamam `NewInstance` com `RABBITMQ_URL`.
- Cada serviço (email, checkout) cria filas necessárias.
- Consume executa handler e faz ack/nack manual.

## 4. Configuração / Env Vars
- `RABBITMQ_URL`

## 5. Exemplo de uso
```go
go
rabbitmq.Publish("email.send", payload)
rabbitmq.Consume("email.send", func(msg Message) { /* process */ })
```

## 6. Falhas comuns
- ErrConnectionClosed

## 7. Notas operacionais
- Aplicar retry exponencial ao publicar quando conexão cair.
