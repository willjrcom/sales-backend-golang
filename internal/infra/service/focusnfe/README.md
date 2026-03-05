# Service / FocusNFe

Cliente HTTP para emissão/consulta/cancelamento de NF-e via FocusNFe/Transmitenota.

---

## 1. Responsabilidades
- Gerar payloads JSON a partir de entidades fiscais.
- Assinar requests com token e lidar com polling de status.
- Tratar erros específicos (ex.: rejeição SEFAZ).

## 2. Métodos principais
| Assinatura | Descrição |
|------------|-----------|
| `CreateNF(ctx, invoice FiscalInvoice) (Response, error)` | Envia NF-e e retorna protocolo. |
| `GetStatus(ctx, id string) (Response, error)` | Consulta status pelo ID Focus. |
| `Cancel(ctx, id string, reason string) error` | Envia cancelamento com justificativa. |

## 3. Fluxo típico
- Usecase fiscal_invoice monta DTO e chama `CreateNF`.
- Serviço converte para JSON FocusNFe, envia via POST `/v2/nfe`.
- Processa resposta, salvando protocolo/numero lote.
- Para acompanhamento, usecase chama `GetStatus` até `authorized`/`rejected`.

## 4. Configuração / Env Vars
- `FOCUSNFE_API_URL`
- `FOCUSNFE_TOKEN`
- `FOCUSNFE_TIMEOUT_MS`

## 5. Exemplo de uso
```go
go
resp, err := focusnfe.CreateNF(ctx, invoice)
if err != nil {
    return fmt.Errorf("nfe error: %w", err)
}
logger.Info("NF-e enviada", "protocol", resp.Protocol)
```

## 6. Falhas comuns
- ErrUnauthorized (token inválido)
- ErrRejected (SEFAZ retornou rejeição)

## 7. Notas operacionais
- Sempre persistir payload enviado/recebido para auditoria fiscal.
