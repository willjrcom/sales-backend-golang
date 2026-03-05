# Repository / Postgres / Fiscal Invoice

Armazena NF-e emitidas, protocolos e links para XML/PDF.

---

## 1. Principais consultas/métodos
| Método | Descrição |
|--------|-----------|
| `Create(ctx, invoice)` | Insere nota com payload do provedor. |
| `GetByOrder(ctx, orderID)` | Busca notas associadas a um pedido. |
| `UpdateStatus(ctx, id)` | Atualiza status/protocolo e timestamps. |

## 2. Transações e locking
- Operações de emissão/cancelamento rodam em tx com order/fiscal_settings.

## 3. Exemplo de SQL
```sql
SELECT id, status, protocol, xml_url
FROM fiscal_invoices
WHERE order_id=@order;
```

## 4. Notas operacionais
- Persistir payload enviado/recebido para auditoria fiscal.
- Não retornar XML inline; apenas URL segura.
