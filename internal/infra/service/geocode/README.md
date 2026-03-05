# Service / Geocode

Integra com provider externo para obter latitude/longitude e validar zonas de entrega.

---

## 1. Responsabilidades
- Converter CEP/endereço em coordenadas.
- Calcular precisão e armazenar status (`ROOFTOP`, `APPROXIMATE`).
- Enriquecer Address antes de salvar.

## 2. Métodos principais
| Assinatura | Descrição |
|------------|-----------|
| `Lookup(ctx, address Address) (GeoCode, error)` | Consulta provider (Google, Here, ViaCEP). |
| `WithinDeliveryZone(lat, lng float64, placeID string) bool` | Valida se ponto está dentro da zona configurada. |

## 3. Fluxo típico
- Quando um endereço novo é salvo, usecase chama `Lookup`.
- Serviço chama API externa e retorna coordenadas + precisão.
- Caso `WithinDeliveryZone` seja falso, flag é salva para alerta no checkout.

## 4. Configuração / Env Vars
- `GEOCODE_PROVIDER`
- `GEOCODE_API_KEY`

## 5. Exemplo de uso
```go
go
geo, err := geocode.Lookup(ctx, addr)
if err == nil {
    addr.GeoCode = geo
}
```

## 6. Falhas comuns
- ProviderTimeout
- ErrNoResults

## 7. Notas operacionais
- Implemente cache (ex.: Redis) para CEPs populares e reduzir custo.
