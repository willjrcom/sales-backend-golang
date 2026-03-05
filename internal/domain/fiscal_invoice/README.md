# Domain / Fiscal Invoice

Representa NF-e emitidas para pedidos.

---

## 1. Entidades principais
| Nome | Descrição |
|------|-----------|
| FiscalInvoice | Status, protocolo, xml/pdf. |
| FiscalConstants | CST/CFOP tabelados. |

## 2. Regras de negócio
- Somente empresas com `fiscal_enabled` podem criar.
- Armazena XML/PDF para reprocessamento.
- Controla prazos de cancelamento.

## 3. Interações e consumidores
- Usecase fiscal_invoice, repositories.*, service focusnfe.

## 4. Exemplo de estrutura
```json
{
  "id": "nfe-200",
  "order_id": "ord-500",
  "status": "authorized",
  "protocol": "1234567890"
}
```
