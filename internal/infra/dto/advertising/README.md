# DTO / Advertising

Estruturas para criar/editar campanhas com assets, segmentação e CTA.

---

## 1. Onde é usado
- handler/advertising.go

## 2. Estruturas principais
| Struct | Campos principais | Direção |
|--------|-------------------|---------|
| CampaignPayload | name, slot, starts_at, ends_at, sponsor_id, category_id, assets | request |
| CampaignResponse | id, status, slots, asset_url, created_at | response |

## 3. Regras de validação
- `starts_at < ends_at`.
- `slot` deve existir na configuração do app.
- `assets` precisa de URL https ou upload prévio.

## 4. Exemplo de request
```json
{
  "name": "Combo sexta",
  "slot": "home-banner",
  "starts_at": "2026-03-10T10:00:00-03:00",
  "ends_at": "2026-03-20T23:59:59-03:00",
  "sponsor_id": "sp-1",
  "category_id": "cat-01",
  "asset_url": "https://s3/banner.png"
}
```

## 5. Exemplo de response
```json
{
  "id": "adv-501",
  "status": "scheduled",
  "slots": [
    "home-banner"
  ],
  "asset_url": "https://s3/banner.png",
  "starts_at": "2026-03-10T10:00:00-03:00"
}
```

## 6. Notas e compatibilidade
- Datas devem ser ISO8601.
- Quando remover campanha, retornar `archived_at`.
