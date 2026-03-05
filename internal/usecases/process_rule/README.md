# Usecase / Process Rule

Define passos automáticos de produção associados a categorias/produtos, influenciando fila e impressão.

---

## 1. Pontos de entrada
| Método | Rota | Origem | Descrição |
|--------|------|--------|-----------|
| POST | `/process-rules` | handler/process_rule.go | Cria regra com etapas. |
| PUT | `/process-rules/{id}` | handler/process_rule.go | Atualiza sequência. |
| DELETE | `/process-rules/{id}` | handler/process_rule.go | Desativa regra mantendo histórico. |

## 2. Dependências
- Repositories: process_rule, product_category, order_process.

## 3. Fluxos e exemplos
### Criar regra
Passos:
- Valida que cada etapa referencia estação/tempo.
- Salva regra e relacionamentos com categorias.
- Atualiza cache usado pelo order usecase.

Exemplo de request:
```json
{
  "name": "Burguer padrão",
  "category_id": "cat-burguer",
  "steps": [
    {
      "station": "grill",
      "duration_sec": 300
    }
  ]
}
```
Resposta:
```json
{
  "process_rule_id": "pr-10",
  "status": "active"
}
```

## 4. Falhas conhecidas
- ErrProcessRuleConflict

## 5. Notas operacionais
- Alterações impactam fila imediatamente; considerar janela de manutenção.
