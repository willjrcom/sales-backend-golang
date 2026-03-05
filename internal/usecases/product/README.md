# Usecase / Product

Administra catálogo de produtos, variações, tamanhos e flags de estoque.

---

## 1. Pontos de entrada
| Método | Rota | Origem | Descrição |
|--------|------|--------|-----------|
| GET | `/products` | handler/product.go | Lista produtos com filtros. |
| POST | `/products` | handler/product.go | Cria produto completo. |
| PUT | `/products/{id}` | handler/product.go | Atualiza atributos, variações, estoque. |

## 2. Dependências
- Repositories: product, product_variation, size, stock.
- Usecases: stock, process_rule, order.

## 3. Fluxos e exemplos
### Criar produto
Passos:
- Valida categoria, preço mínimo e integrações.
- Cria produto + variações + tamanhos num único fluxo.
- Se `track_stock=true`, cria registros em stock.

Exemplo de request:
```json
{
  "name": "Combo X",
  "category_id": "cat-burguer",
  "base_price": 25,
  "variations": [
    {
      "name": "Carne",
      "sku": "CBX-C"
    }
  ]
}
```
Resposta:
```json
{
  "product_id": "prod-77",
  "status": "active"
}
```

## 4. Falhas conhecidas
- ErrProductDuplicateSKU
- ErrInvalidVariation

## 5. Notas operacionais
- Mudanças de preço devem disparar evento para sincronizar PDV offline.
