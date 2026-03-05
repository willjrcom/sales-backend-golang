# Domain / Entity

Estruturas base com IDs, timestamps e audit log reutilizadas pelo domínio.

---

## 1. Entidades principais
| Nome | Descrição |
|------|-----------|
| Entity | ID, CreatedAt, UpdatedAt, DeletedAt. |
| EntityMetadata | Pares chave/valor auxiliares. |

## 2. Regras de negócio
- Soft delete padronizado.
- Metadados sempre opcional e versionado.

## 3. Interações e consumidores
- Incorporado em quase todos os domínios.

## 4. Exemplo de estrutura
```json
{
  "id": "uuid",
  "created_at": "2026-03-05T12:00:00Z",
  "updated_at": "2026-03-05T12:00:00Z",
  "deleted_at": null
}
```
