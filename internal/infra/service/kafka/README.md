# Service / Kafka

Wrapper básico para publicar/consumir eventos Kafka quando habilitado.

---

## 1. Responsabilidades
- Criar producers/consumers com config única.
- Serializar mensagens JSON.
- Fornecer helpers para tópicos padronizados (order.events, stock.alerts).

## 2. Métodos principais
| Assinatura | Descrição |
|------------|-----------|
| `Publish(topic string, payload any) error` | Serializa e envia com chave opcional. |
| `Consume(topic string, handler func(Message))` | Loop de consumo com commit manual. |

## 3. Fluxo típico
- Usecase chama `Publish` quando precisa notificar outros sistemas.
- Consumidores são configurados em módulos específicos (ex.: integrations).

## 4. Configuração / Env Vars
- `KAFKA_BROKERS`
- `KAFKA_CLIENT_ID`

## 5. Exemplo de uso
```go
go
kafka.Publish("order.events", map[string]any{"order_id": order.ID, "status": order.Status})
```

## 6. Falhas comuns
- ErrKafkaUnavailable

## 7. Notas operacionais
- Atualmente opcional; degrade para RabbitMQ se brokers indisponíveis.
