# Domain / Company Category

Associa empresas a categorias operacionais e seus patrocinadores.

---

## 1. Entidades principais
| Nome | Descrição |
|------|-----------|
| CompanyCategory | Categoria atribuída à empresa. |
| CategoryAdvertising | Ligação com campanhas ativas. |
| CategorySponsor | Patrocinadores exclusivos por categoria. |

## 2. Regras de negócio
- Uma categoria pode exigir patrocinador exclusivo.
- Categorias controlam disponibilidade de produtos/processos.

## 3. Interações e consumidores
- Usecases: company_category, advertising, product.

## 4. Exemplo de estrutura
```json
{
  "company_id": "cmp-100",
  "category_id": "restaurant",
  "sponsors": [
    "sp-1"
  ]
}
```
