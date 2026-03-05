# Usecase / Client

Administra clientes finais, histórico, preferências e vínculos de endereço/contato.

---

## 1. Pontos de entrada
| Método | Rota | Origem | Descrição |
|--------|------|--------|-----------|
| GET | `/clients` | handler/client.go | Listagem paginada com filtros por documento, telefone e status. |
| POST | `/clients` | handler/client.go | Cria cliente sincronizando person/contact/address. |
| PUT | `/clients/{id}` | handler/client.go | Atualiza dados pessoais e flags de bloqueio. |
| GET | `/clients/{id}/history` | handler/client.go | Retorna pedidos e tickets médios. |

## 2. Dependências
- Repositories: client, person, contact, address, order.
- Services: geocode (endereços), email.

## 3. Fluxos e exemplos
### Cadastro completo
Passos:
- Normaliza CPF/telefone e verifica duplicidade.
- Cria/atualiza registros de `person`, `contact` e `address` numa transação.
- Salva preferências e retorna cliente completo.

Exemplo de request:
```json
{
  "name": "Maria Alves",
  "document": "12345678909",
  "phones": [
    "+5511988887777"
  ],
  "address": {
    "zip_code": "04000-000",
    "street": "Rua Exemplo",
    "number": "123"
  }
}
```
Resposta:
```json
{
  "id": "cli-200",
  "status": "active",
  "loyalty_score": 4
}
```

### Consulta histórico
Passos:
- Busca pedidos em `order` filtrando por cliente e período.
- Agrega métricas (ticket médio, últimos status).
- Retorna timeline usada pelo atendimento.

Exemplo de request:
```json
{}
```
Resposta:
```json
{
  "orders": [
    {
      "order_id": "ord-10",
      "status": "finished",
      "total": 85
    },
    {
      "order_id": "ord-11",
      "status": "canceled",
      "total": 50
    }
  ],
  "metrics": {
    "avg_ticket": 67.5
  }
}
```

## 4. Falhas conhecidas
- ErrClientDuplicate: documento já cadastrado.
- ErrInvalidAddress: geocode não encontrou o CEP informado.

## 5. Notas operacionais
- Flag `blocked_reason` impede checkout até revisão manual.
- Quando apagar cliente, manter soft delete para preservar histórico.
