# Usecase / Place

Administra locais físicos (lojas, dark kitchens, pontos pickup) e integra com mesas e entregas.

---

## 1. Pontos de entrada
| Método | Rota | Origem | Descrição |
|--------|------|--------|-----------|
| GET | `/places` | handler/place.go | Lista locais com horários. |
| POST | `/places` | handler/place.go | Cria novo local. |
| PUT | `/places/{id}` | handler/place.go | Atualiza horas, capacidade e integrações de entrega. |

## 2. Dependências
- Repositories: place, table, address.

## 3. Fluxos e exemplos
### Cadastrar local
Passos:
- Valida endereço/CEP.
- Define capacidades (mesas, slots pickup).
- Ativa canais disponíveis (delivery/pickup).

Exemplo de request:
```json
{
  "name": "Loja Norte",
  "address_id": "addr-9",
  "channels": [
    "delivery",
    "pickup"
  ]
}
```
Resposta:
```json
{
  "place_id": "plc-2",
  "status": "active"
}
```

## 4. Falhas conhecidas
- ErrPlaceDuplicate
- ErrInvalidChannel

## 5. Notas operacionais
- Mudanças afetam cálculo de SLA de entrega; invalidar caches.
