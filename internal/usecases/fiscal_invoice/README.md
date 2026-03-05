# Usecase / Fiscal Invoice

Integra com FocusNFe/Transmitenota para emitir, cancelar e consultar NF-e vinculadas a pedidos.

---

## 1. Pontos de entrada
| Método | Rota | Origem | Descrição |
|--------|------|--------|-----------|
| POST | `/fiscal-invoice` | handler/fiscal_invoice.go | Gera NF-e a partir de um pedido. |
| POST | `/fiscal-invoice/{id}/cancel` | handler/fiscal_invoice.go | Cancela nota autorizada. |
| GET | `/fiscal-invoice/{id}` | handler/fiscal_invoice.go | Consulta status e baixa XML/PDF. |

## 2. Dependências
- Repositories: fiscal_invoice, order, company, company_subscription.
- Services: focusnfe, email (envio de DANFE).

## 3. Fluxos e exemplos
### Emitir NF-e
Passos:
- Valida que empresa possui configuração fiscal completa.
- Mapeia itens/pagamentos para layout XML.
- Chama FocusNFe e salva protocolo/arquivos.

Exemplo de request:
```json
{
  "order_id": "ord-100",
  "environment": "production"
}
```
Resposta:
```json
{
  "fiscal_invoice_id": "nfe-200",
  "status": "authorized",
  "protocol": "1234567890"
}
```

### Cancelar
Passos:
- Valida prazo legal.
- Envia justificativa ao provedor.
- Atualiza pedido e dispara email ao cliente.

Exemplo de request:
```json
{
  "reason": "Pedido cancelado pelo cliente"
}
```
Resposta:
```json
{
  "status": "canceled",
  "canceled_at": "2026-03-05T13:00:00Z"
}
```

## 4. Falhas conhecidas
- ErrFiscalConfigMissing
- FocusNFeError
- ErrFiscalInvoiceAlreadyProcessed

## 5. Notas operacionais
- Salvar XML e PDF em S3 para reenvio futuro.
