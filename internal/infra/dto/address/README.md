# DTO / Address

DTOs usados em cadastros de endereço (empresa, cliente, delivery) incluindo suporte a geocode e validação de CEP.

---

## 1. Onde é usado
- handler/company.go
- handler/client.go
- handler/order_delivery.go

## 2. Estruturas principais
| Struct | Campos principais | Direção |
|--------|-------------------|---------|
| AddressRequest | zip_code, street, number, complement, neighborhood, city, state, reference | request |
| AddressResponse | id, zip_code, street, number, city, state, geocode | response |

## 3. Regras de validação
- `zip_code` deve ter 8 dígitos numéricos.
- `state` sempre em duas letras maiúsculas.
- Quando `geocode_required=true`, validar lat/lng no payload.

## 4. Exemplo de request
```json
{
  "zip_code": "04000-000",
  "street": "Rua Exemplo",
  "number": "123",
  "city": "São Paulo",
  "state": "SP",
  "geocode_required": true
}
```

## 5. Exemplo de response
```json
{
  "id": "addr-1",
  "zip_code": "04000-000",
  "street": "Rua Exemplo",
  "number": "123",
  "city": "São Paulo",
  "state": "SP",
  "geocode": {
    "lat": -23.5,
    "lng": -46.6,
    "precision": "ROOFTOP"
  }
}
```

## 6. Notas e compatibilidade
- Sempre retornar `id` para permitir atualização incremental no frontend.
- Campos opcionais devem aparecer como null para manter consistência.
