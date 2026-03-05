# DTO / Company Category

Estruturas para atribuir categorias e patrocinadores a uma empresa.

---

## 1. Onde é usado
- handler/company_category.go

## 2. Estruturas principais
| Struct | Campos principais | Direção |
|--------|-------------------|---------|
| AssignCategoriesRequest | categories[] | request |
| CompanyCategoryResponse | company_id, categories[], sponsors[] | response |

## 3. Regras de validação
- Lista não pode ser vazia.
- Categorias devem existir na tabela de referência.

## 4. Exemplo de request
```json
{
  "categories": [
    "restaurant",
    "bar"
  ]
}
```

## 5. Exemplo de response
```json
{
  "company_id": "cmp-100",
  "categories": [
    "restaurant",
    "bar"
  ],
  "sponsors": [
    "sp-1"
  ]
}
```

## 6. Notas e compatibilidade
- Quando categories muda, invalidar caches de menu.
