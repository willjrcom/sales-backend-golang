# Domain / Address

Modela endereços completos utilizados por empresas, clientes e entregas, suportando geocodificação e múltiplos contatos.

---

## 1. Entidades principais
| Nome | Descrição |
|------|-----------|
| Address | Logradouro, número, complemento, bairro, cidade, UF, CEP. |
| GeoCode | Latitude/longitude e status da precisão da busca. |

## 2. Regras de negócio
- Campos normalizados (CEP numérico, UF em duas letras).
- Permite marcar endereço principal por entidade/grupo.
- Se `GeoCodeStatus=partial`, bloqueia entregas fora da zona configurada.

## 3. Interações e consumidores
- Usecases: company, client, order_delivery, place.
- Serviços: geocode, delivery SLA.

## 4. Exemplo de estrutura
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
