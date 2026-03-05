# Domain / Product

Produtos, variações, tamanhos, categorias e regras de processo.

---

## 1. Entidades principais
| Nome | Descrição |
|------|-----------|
| Product | Dados principais e flags. |
| ProductVariation | SKUs/estoque individual. |
| ProductSize | Tamanhos com preços. |
| Category linkage | Referência a product_category. |
| ProcessRule link | Define etapas de produção. |

## 2. Regras de negócio
- Produto pode exigir controle de estoque FIFO (variations).
- Preços armazenados com currency decimal fixo.
- Flag `allow_additionals` habilita complementos no order.

## 3. Interações e consumidores
- Usecases: product, product_category, size, stock.

## 4. Exemplo de estrutura
```json
{
  "id": "prod-77",
  "name": "Combo X",
  "track_stock": true,
  "variations": [
    {
      "id": "var-1",
      "sku": "CBX-C"
    }
  ],
  "sizes": [
    {
      "id": "size-1",
      "name": "Grande",
      "price": 32
    }
  ]
}
```
