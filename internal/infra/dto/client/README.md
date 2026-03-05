# DTO / Client

Payloads para cadastro/edição de clientes finais e retorno resumido/detalhado.

---

## 1. Onde é usado
- handler/client.go

## 2. Estruturas principais
| Struct | Campos principais | Direção |
|--------|-------------------|---------|
| ClientRequest | name, document, phones[], address, preferences | request |
| ClientResponse | id, name, document, loyalty_score, blocked_reason, last_order_at | response |

## 3. Regras de validação
- Documento deve ter 11 dígitos (CPF).
- `phones` no padrão E.164.
- Endereço pode ser null para clientes apenas com pickup.

## 4. Exemplo de request
```json
{
  "name": "Maria",
  "document": "12345678909",
  "phones": [
    "+5511988887777"
  ],
  "address": {
    "zip_code": "04000-000",
    "street": "Rua X",
    "number": "123"
  }
}
```

## 5. Exemplo de response
```json
{
  "id": "cli-200",
  "name": "Maria",
  "document": "12345678909",
  "loyalty_score": 4,
  "blocked_reason": null,
  "last_order_at": "2026-03-01T10:00:00Z"
}
```

## 6. Notas e compatibilidade
- Sempre mascarar documento na UI (***).
- `blocked_reason` null significa liberado.
