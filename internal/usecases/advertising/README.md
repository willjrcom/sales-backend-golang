# Usecase / Advertising

Gerencia campanhas promocionais exibidas no PDV e app do cliente, garantindo que apenas anúncios válidos e dentro de vigência sejam entregues.

---

## 1. Pontos de entrada
| Método | Rota | Origem | Descrição |
|--------|------|--------|-----------|
| GET | `/advertising` | handler/advertising.go | Lista campanhas ativas, com filtros por categoria e patrocinador. |
| POST | `/advertising` | handler/advertising.go | Cria campanha e agenda publicação. |
| PUT | `/advertising/{id}` | handler/advertising.go | Atualiza mídias, vigência e segmentação. |
| DELETE | `/advertising/{id}` | handler/advertising.go | Arquiva a campanha e emite evento para remover do app. |

## 2. Dependências
- Repositories: advertising, sponsor, company_category.
- Services: S3 (upload de assets), RabbitMQ (evento advertising.changed).
- DTOs: infra/dto/advertising.

## 3. Fluxos e exemplos
### Publicar campanha
Passos:
- Valida janela inicial/final e exclusividade por slot/categoria.
- Upload da arte principal para S3 e salva URL assinada.
- Cria registro no repo e dispara evento `advertising.created` para o app atualizar o carrossel.

Exemplo de request:
```json
{
  "name": "Combo de sexta",
  "category_id": "cat-01",
  "sponsor_id": "spon-9",
  "starts_at": "2026-03-10T10:00:00-03:00",
  "ends_at": "2026-03-20T23:59:59-03:00",
  "cta_url": "https://promo.exemplo"
}
```
Resposta:
```json
{
  "id": "adv-501",
  "status": "scheduled",
  "asset_url": "https://s3/promo.png",
  "slots": [
    "home-banner"
  ]
}
```

### Expirar automaticamente
Passos:
- Job diário consulta campanhas cujo `ends_at < now`.
- Atualiza status para `expired` e remove dos caches.
- Envia evento `advertising.expired` para o frontend invalidar cards.

Tarefa agendada:
```json
{
  "task": "expire-advertising",
  "batch_size": 200
}
```
Resposta:
```json
{
  "processed": 27,
  "expired_ids": [
    "adv-1",
    "adv-2"
  ]
}
```

## 4. Falhas conhecidas
- ErrOverlappingSlot: existe outra campanha com o mesmo slot e janela.
- ErrInvalidSponsor: patrocinador inativo ou fora da categoria.
- S3UploadError: falha ao armazenar ativo gráfico.

## 5. Notas operacionais
- Sempre remover assets órfãos no S3 quando campanha é deletada.
- Logs devem incluir `company_id` e `slot` para auditoria.
