# DTO / Sponsor

DTOs para patrocinadores e assets.

---

## 1. Onde é usado
- handler/sponsor.go

## 2. Estruturas principais
| Struct | Campos principais | Direção |
|--------|-------------------|---------|
| SponsorRequest | name, contract_start, contract_end, benefits | request |
| SponsorResponse | id, status, contract_end | response |

## 3. Regras de validação
- Datas no formato ISO.

## 4. Exemplo de request
```json
{
  "name": "Bebidas XYZ",
  "contract_start": "2026-01-01",
  "contract_end": "2026-12-31"
}
```

## 5. Exemplo de response
```json
{
  "id": "sp-1",
  "status": "active",
  "contract_end": "2026-12-31"
}
```

## 6. Notas e compatibilidade
- Armazenar logos via S3 DTO.
