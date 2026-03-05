# Usecase / Product Category

Gerencia categorias do cardápio, hierarquias e vínculos com processos/patrocinadores.

---

## 1. Pontos de entrada
| Método | Rota | Origem | Descrição |
|--------|------|--------|-----------|
| GET | `/product/categories` | handler/product_category.go | Lista categorias. |
| POST | `/product/categories` | handler/product_category.go | Cria categoria. |
| PUT | `/product/categories/{id}` | handler/product_category.go | Atualiza hierarquia e processos. |

## 2. Dependências
- Repositories: product_category, process_rule, sponsor.

## 3. Fluxos e exemplos
### Criar categoria com processo
Passos:
- Valida ausência de loops hierárquicos.
- Relaciona process_rule quando informado.
- Atualiza permissões de estoque/adicionais.

Exemplo de request:
```json
{
  "name": "Burgers",
  "parent_id": null,
  "process_rule_id": "pr-10"
}
```
Resposta:
```json
{
  "category_id": "cat-burgers",
  "status": "active"
}
```

## 4. Falhas conhecidas
- ErrCategoryLoop

## 5. Notas operacionais
- Categorias são usadas para dashboards; mantenha `label_key` amigável.
