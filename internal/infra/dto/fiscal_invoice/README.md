# DTO / Fiscal Invoice

Payloads para emissão/cancelamento/consulta de NF-e.

---

## 1. Onde é usado
- handler/fiscal_invoice.go

## 2. Estruturas principais
| Struct | Campos principais | Direção |
|--------|-------------------|---------|
| FiscalInvoiceRequest | order_id, environment | request |
| FiscalInvoiceResponse | id, status, protocol, xml_url, pdf_url | response |
| FiscalInvoiceCancelRequest | reason | request |

## 3. Regras de validação
- `environment` ∈ {production,sandbox}.
- `reason` mínimo 15 caracteres no cancelamento.

## 4. Exemplo de request
```json
{
  "order_id": "ord-100",
  "environment": "production"
}
```

## 5. Exemplo de response
```json
{
  "id": "nfe-200",
  "status": "authorized",
  "protocol": "1234567890",
  "xml_url": "https://s3/nfe-200.xml"
}
```

## 6. Notas e compatibilidade
- Nunca retornar XML inline; apenas URL.
