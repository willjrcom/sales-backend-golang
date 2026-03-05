# DTO / Report

Filtros e respostas para dashboards (sales, estoque, marketing).

---

## 1. Onde é usado
- handler/report.go

## 2. Estruturas principais
| Struct | Campos principais | Direção |
|--------|-------------------|---------|
| ReportFilter | from, to, group_by, channel, category_ids[] | request |
| ReportSeriesResponse | label, value, meta | response |

## 3. Regras de validação
- `to - from` ≤ 92 dias.
- `group_by` ∈ {hour,day,week,month}.

## 4. Exemplo de request
```json
{
  "from": "2026-03-01",
  "to": "2026-03-05",
  "group_by": "day",
  "channel": "delivery"
}
```

## 5. Exemplo de response
```json
{
  "series": [
    {
      "label": "2026-03-01",
      "value": 1230
    }
  ],
  "totals": {
    "orders": 45
  }
}
```

## 6. Notas e compatibilidade
- Para dados grandes, habilitar paginação/stream.
