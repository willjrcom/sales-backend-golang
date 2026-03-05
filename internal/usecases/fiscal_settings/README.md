# Usecase / Fiscal Settings

Permite configurar credenciais e parâmetros fiscais exigidos para emissão (CRT, CNAE, certificados).

---

## 1. Pontos de entrada
| Método | Rota | Origem | Descrição |
|--------|------|--------|-----------|
| GET | `/fiscal/settings` | handler/fiscal_settings.go | Retorna configurações atuais. |
| POST | `/fiscal/settings` | handler/fiscal_settings.go | Atualiza campos fiscais e validações. |
| POST | `/fiscal/settings/test` | handler/fiscal_settings.go | Executa teste de comunicação com provedor. |

## 2. Dependências
- Repositories: fiscal_settings, company.
- Services: focusnfe.

## 3. Fluxos e exemplos
### Salvar configuração
Passos:
- Valida campos obrigatórios conforme regime tributário.
- Criptografa credenciais antes de salvar.
- Se `test_immediately=true`, dispara teste assíncrono.

Exemplo de request:
```json
{
  "company_id": "cmp-100",
  "crt": 1,
  "cnae": "5611201",
  "certificate_base64": "..."
}
```
Resposta:
```json
{
  "company_id": "cmp-100",
  "status": "ready"
}
```

## 4. Falhas conhecidas
- ErrCertificateInvalid
- ErrMissingFields

## 5. Notas operacionais
- Sensível: logs nunca devem incluir credenciais ou certificados.
