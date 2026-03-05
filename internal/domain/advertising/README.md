# Domain / Advertising

Entidades para campanhas promocionais/ads com assets, vigência e segmentação por categoria/patrocinador.

---

## 1. Entidades principais
| Nome | Descrição |
|------|-----------|
| Advertising | Dados principais (nome, slot, assets, CTA). |
| AdvertisingTarget | Liga campanhas a categorias, patrocinadores ou locais. |

## 2. Regras de negócio
- Um slot/categoria só pode ter uma campanha ativa por janela.
- Status calculado automaticamente (scheduled, running, expired).
- Assets armazenam versões para web e PDV.

## 3. Interações e consumidores
- Usecases: advertising, company_category, sponsor.
- Infra: handler/advertising, report dashboards.

## 4. Exemplo de estrutura
```json
{
  "id": "adv-10",
  "slot": "home-banner",
  "category_id": "cat-burgers",
  "sponsor_id": "sp-1",
  "starts_at": "2026-03-10T10:00:00Z",
  "ends_at": "2026-03-20T23:59:00Z",
  "asset_url": "https://s3/banner.png"
}
```
