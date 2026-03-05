# DTO / Product Category

DTOs para categorias do cardápio.

---

## 1. Onde é usado
- handler/product_category.go

## 2. Estruturas principais
| Struct | Campos principais | Direção |
|--------|-------------------|---------|
| CategoryRequest | name, parent_id, process_rule_id | request |
| CategoryResponse | id, name, parent_id, process_rule_id | response |

## 3. Regras de validação
- Sem loops hierárquicos.

## 4. Exemplo de request
```json
{
  "name": "Burgers",
  "parent_id": null,
  "process_rule_id": "pr-10"
}
```

## 5. Exemplo de response
```json
{
  "id": "cat-burgers",
  "name": "Burgers",
  "parent_id": null,
  "process_rule_id": "pr-10"
}
```

## 6. Notas e compatibilidade
- `label_key` usado nos relatórios deve ser estável.
