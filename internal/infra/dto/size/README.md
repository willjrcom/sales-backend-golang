# DTO / Size

DTOs para tamanhos vinculados a produtos.

---

## 1. Onde é usado
- handler/size.go

## 2. Estruturas principais
| Struct | Campos principais | Direção |
|--------|-------------------|---------|
| SizeRequest | name, sku_suffix, price | request |
| SizeResponse | id, name, price, status | response |

## 3. Regras de validação
- `price`>0.
- `sku_suffix` max 3 chars.

## 4. Exemplo de request
```json
{
  "name": "Grande",
  "sku_suffix": "G",
  "price": 32
}
```

## 5. Exemplo de response
```json
{
  "id": "size-1",
  "name": "Grande",
  "price": 32,
  "status": "active"
}
```

## 6. Notas e compatibilidade
- Quando inativado, manter `status=inactive`.
