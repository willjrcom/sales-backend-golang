# Domain / Schema

Metadados do schema PostgreSQL por empresa.

---

## 1. Entidades principais
| Nome | Descrição |
|------|-----------|
| Schema | Nome, status, provisioned_at. |
| SchemaHistory | Migrations aplicadas. |

## 2. Regras de negócio
- Status determina se requests podem usar o schema.
- Histórico usado pelos comandos `migrate`.

## 3. Interações e consumidores
- Bootstrap/migrations, Company usecase.

## 4. Exemplo de estrutura
```json
{
  "name": "tenant_abc",
  "status": "active",
  "migrations": [
    "001_base.sql"
  ]
}
```
