# Usecase / Company Category

Mantém categorias operacionais da empresa e vínculos com patrocinadores/ads para habilitar fluxos específicos.

---

## 1. Pontos de entrada
| Método | Rota | Origem | Descrição |
|--------|------|--------|-----------|
| GET | `/company/categories` | handler/company_category.go | Lista categorias disponíveis com flags. |
| POST | `/company/{id}/categories` | handler/company_category.go | Atribui conjunto de categorias à empresa. |
| POST | `/company/{id}/categories/{category_id}/sponsors` | handler/company_category.go | Liga patrocinadores a uma categoria da empresa. |

## 2. Dependências
- Repositories: company_category, sponsor, advertising.
- Usecases relacionados: advertising, product.

## 3. Fluxos e exemplos
### Atribuir categorias
Passos:
- Valida se a categoria existe e é compatível com o plano.
- Remove vínculos antigos e cria novos registros na tabela pivot.
- Comunica módulos de produto/processo para liberar fluxos.

Exemplo de request:
```json
{
  "categories": [
    "restaurant",
    "bar"
  ]
}
```
Resposta:
```json
{
  "company_id": "cmp-100",
  "categories": [
    "restaurant",
    "bar"
  ],
  "effective_at": "2026-03-05T10:00:00Z"
}
```

### Vincular patrocinador
Passos:
- Checa exclusividade do patrocinador para a categoria.
- Cria vínculo com data início/fim.
- Dispara evento `category.sponsor.updated`.

Exemplo de request:
```json
{
  "sponsor_id": "sp-10",
  "starts_at": "2026-03-10",
  "ends_at": "2026-06-10"
}
```
Resposta:
```json
{
  "company_id": "cmp-100",
  "category": "restaurant",
  "sponsor_id": "sp-10",
  "status": "active"
}
```

## 4. Falhas conhecidas
- ErrCategoryNotAllowed
- ErrSponsorConflict

## 5. Notas operacionais
- Categorias controlam quais módulos ficam visíveis no app do operador.
