# Usecase / Report

Executa consultas analíticas (vendas, estoque, marketing) retornando séries para dashboards web e mobile.

---

## 1. Pontos de entrada
| Método | Rota | Origem | Descrição |
|--------|------|--------|-----------|
| POST | `/report/sales-summary` | handler/report.go | Resumo de vendas por período. |
| POST | `/report/additional-items-sold` | handler/report.go | Top adicionais. |
| POST | `/report/complements-sold` | handler/report.go | Top complementos. |

## 2. Dependências
- Repositories: report (consultas SQL customizadas), order, stock.
- Services: cache (opcional).

## 3. Fluxos e exemplos
### Gerar resumo de vendas
Passos:
- Valida range de datas (máx 92 dias).
- Constrói SQL com filtros por canal/empresa.
- Retorna série agregada e totais.

Exemplo de request:
```json
{
  "from": "2026-03-01",
  "to": "2026-03-05",
  "group_by": "day"
}
```
Resposta:
```json
{
  "series": [
    {
      "date": "2026-03-01",
      "total": 1230
    }
  ],
  "totals": {
    "orders": 45,
    "gross": 3200
  }
}
```

## 4. Falhas conhecidas
- ErrReportTooLarge
- ErrUnknownMetric

## 5. Notas operacionais
- Relatórios intensivos devem usar cache em Redis para repetidas consultas.
