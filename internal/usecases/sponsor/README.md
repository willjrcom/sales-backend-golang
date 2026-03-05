# Usecase / Sponsor

Gerencia patrocinadores, contratos e incentivos ligados a categorias/ads.

---

## 1. Pontos de entrada
| Método | Rota | Origem | Descrição |
|--------|------|--------|-----------|
| POST | `/sponsors` | handler/sponsor.go | Cria patrocinador. |
| PUT | `/sponsors/{id}` | handler/sponsor.go | Atualiza contrato. |
| POST | `/sponsors/{id}/assets` | handler/sponsor.go | Cadastra logos/benefícios. |

## 2. Dependências
- Repositories: sponsor, company_category, advertising.
- Services: S3 para assets.

## 3. Fluxos e exemplos
### Cadastrar patrocinador
Passos:
- Salva dados legais e vigência.
- Relaciona benefícios (descontos, combos).

Exemplo de request:
```json
{
  "name": "Bebidas XYZ",
  "contract_start": "2026-01-01",
  "contract_end": "2026-12-31"
}
```
Resposta:
```json
{
  "sponsor_id": "sp-1",
  "status": "active"
}
```

## 4. Falhas conhecidas
- ErrSponsorOverlap

## 5. Notas operacionais
- Renovações devem ser criadas como novo registro para manter histórico.
