# Domain / Fiscal Settings

Configuração fiscal da empresa (IE, CNAE, CRT, credenciais).

---

## 1. Entidades principais
| Nome | Descrição |
|------|-----------|
| FiscalSettings | Campos fiscais e certificados. |

## 2. Regras de negócio
- Credenciais armazenadas criptografadas.
- Validação muda conforme regime (Simples, Lucro Real).

## 3. Interações e consumidores
- Usecases: fiscal_settings, company, fiscal_invoice.

## 4. Exemplo de estrutura
```json
{
  "company_id": "cmp-100",
  "crt": 1,
  "cnae": "5611201",
  "certificate_hash": "sha256..."
}
```
