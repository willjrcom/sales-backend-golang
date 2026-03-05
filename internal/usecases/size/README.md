# Usecase / Size

Mantém tamanhos e preços associados às variações de produto.

---

## 1. Pontos de entrada
| Método | Rota | Origem | Descrição |
|--------|------|--------|-----------|
| POST | `/products/{id}/sizes` | handler/size.go | Cria tamanho para produto. |
| PUT | `/products/{id}/sizes/{size_id}` | handler/size.go | Atualiza preço/flag. |

## 2. Dependências
- Repositories: size, product, stock.

## 3. Fluxos e exemplos
### Adicionar tamanho
Passos:
- Valida combinação tamanho+variação única.
- Calcula preço final baseado em markup.
- Atualiza estoque se controlado por tamanho.

Exemplo de request:
```json
{
  "name": "Grande",
  "price": 32,
  "sku_suffix": "G"
}
```
Resposta:
```json
{
  "size_id": "size-1",
  "status": "active"
}
```

## 4. Falhas conhecidas
- ErrSizeDuplicate
- ErrSizeReferenced

## 5. Notas operacionais
- Quando remover tamanho, inative ao invés de deletar para preservar históricos.
