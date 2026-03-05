# Usecase / Company

Executa onboarding de novas empresas (tenants), gerencia preferências operacionais e sincroniza billing (subscription + usage cost).

---

## 1. Pontos de entrada
| Método | Rota | Origem | Descrição |
|--------|------|--------|-----------|
| POST | `/company` | handler/company.go | Cria empresa, schema e usuário administrador. |
| PUT | `/company/{id}` | handler/company.go | Atualiza dados cadastrais e status. |
| POST | `/company/{id}/preferences` | handler/company.go | Atualiza preferências operacionais e flags do PDV. |
| POST | `/company/{id}/subscription/activate` | handler/company.go | Liga/desliga assinatura e calcula billing inicial. |

## 2. Dependências
- Repositories: company, schema, user, address, company_subscription, company_usage_cost.
- Services: rabbitmq (eventos de onboard), email (convite admin), schema service (criação física).

## 3. Fluxos e exemplos
### Onboarding completo
Passos:
- Valida CNPJ e evita duplicidade.
- Cria schema dedicado via módulo `schema` e executa migrations básicas.
- Registra empresa, endereço, preferências padrão e usuário master.
- Dispara email de boas-vindas e evento `company.created`.

Exemplo de request:
```json
{
  "name": "Loja Central",
  "document": "12.345.678/0001-99",
  "owner_email": "admin@loja.com",
  "preferences": {
    "allow_negative_stock": true
  }
}
```
Resposta:
```json
{
  "id": "cmp-100",
  "schema": "tenant_abc",
  "status": "trial",
  "admin_user_id": "usr-1"
}
```

### Atualizar preferências
Passos:
- Recebe payload parcial e valida regras (ex.: não permitir desativar estoque se há produtos controlados).
- Persiste preferências versionadas.
- Invalida caches de configuração e notifica apps conectados.

Exemplo de request:
```json
{
  "allow_negative_stock": false,
  "default_checkout_methods": [
    "cash",
    "pix"
  ],
  "stock_alert_webhook": "https://hooks/alert"
}
```
Resposta:
```json
{
  "company_id": "cmp-100",
  "updated_fields": [
    "allow_negative_stock",
    "default_checkout_methods"
  ]
}
```

## 4. Falhas conhecidas
- ErrDuplicateCompany
- ErrSchemaProvision
- ErrPreferenceConflict

## 5. Notas operacionais
- Alterações sensíveis devem ser logadas com `company_id` e `user_id`.
- Ao suspender empresa, também revogue tokens ativos.
