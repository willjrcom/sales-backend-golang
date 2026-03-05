# Domain / Company

Entidade do tenant com status operacional, assinaturas e preferências.

---

## 1. Entidades principais
| Nome | Descrição |
|------|-----------|
| Company | Dados principais e flags. |
| CompanyPreference | Configurações do PDV/estoque. |
| CompanySubscription | Plano vigente e billing. |
| CompanyUsageCost | Medições de custo variável. |

## 2. Regras de negócio
- Status (trial, active, suspended) habilita/desabilita módulos.
- Preferências versionadas para auditoria.
- Uso excedente gera cobrança automática registrada em `CompanyUsageCost`.

## 3. Interações e consumidores
- Usecases: company, checkout, fiscal_settings, report.
- Infra: repository/postgres/company.*

## 4. Exemplo de estrutura
```json
{
  "id": "cmp-100",
  "schema": "tenant_abc",
  "status": "active",
  "preferences": {
    "allow_negative_stock": true
  }
}
```
