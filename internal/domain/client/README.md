# Domain / Client

Representa cliente final com vínculo a person, contatos e preferências de fidelidade.

---

## 1. Entidades principais
| Nome | Descrição |
|------|-----------|
| Client | Campos de fidelidade, bloqueios e métricas. |
| ClientHistory | Snapshots de pedidos e tickets. |

## 2. Regras de negócio
- Cada cliente pertence a uma empresa (schema) e usa soft delete.
- `blocked_reason` impede novos pedidos até revisão.
- Relaciona último endereço/contato preferido para checkout rápido.

## 3. Interações e consumidores
- Usecases: client, order, checkout.
- Repositories: client, order.

## 4. Exemplo de estrutura
```json
{
  "id": "cli-200",
  "person_id": "per-50",
  "loyalty_score": 4,
  "blocked_reason": null,
  "last_order_at": "2026-03-03T12:00:00Z"
}
```
