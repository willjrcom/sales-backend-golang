# Usecase / Table

CRUD de mesas físicas/virtuais e integração com layout de salão.

---

## 1. Pontos de entrada
| Método | Rota | Origem | Descrição |
|--------|------|--------|-----------|
| GET | `/tables` | handler/table.go | Lista mesas e status. |
| POST | `/tables` | handler/table.go | Cria mesa. |
| PUT | `/tables/{id}` | handler/table.go | Atualiza capacidade e posição. |

## 2. Dependências
- Repositories: table, place.

## 3. Fluxos e exemplos
### Criar mesa
Passos:
- Valida que identificador é único por local.
- Associa place/zona.

Exemplo de request:
```json
{
  "name": "M10",
  "place_id": "plc-2",
  "capacity": 4
}
```
Resposta:
```json
{
  "table_id": "tbl-10",
  "status": "available"
}
```

## 4. Falhas conhecidas
- ErrTableDuplicate

## 5. Notas operacionais
- Mesas virtuais (delivery hubs) devem ter flag `virtual=true`.
