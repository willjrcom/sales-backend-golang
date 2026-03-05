# DTO / Entity

DTOs genéricos usados em selects (id/label).

---

## 1. Onde é usado
- handlers diversos

## 2. Estruturas principais
| Struct | Campos principais | Direção |
|--------|-------------------|---------|
| SelectOption | id, label | response |

## 3. Regras de validação
- Ideal para combos no frontend.

## 4. Exemplo de request
```json
{}
```

## 5. Exemplo de response
```json
{
  "id": "prod-1",
  "label": "Burger Clássico"
}
```

## 6. Notas e compatibilidade
- Não adicionar informações sensíveis aqui.
