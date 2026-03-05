# DTO / Company

DTOs para onboarding e atualização de empresa/tenant.

---

## 1. Onde é usado
- handler/company.go

## 2. Estruturas principais
| Struct | Campos principais | Direção |
|--------|-------------------|---------|
| CompanyCreateRequest | name, document, owner_email, preferences | request |
| CompanyResponse | id, schema, status, preferences | response |

## 3. Regras de validação
- Documento no formato CNPJ normalizado.
- `owner_email` precisa ser único.
- Preferências aceitam apenas chaves conhecidas.

## 4. Exemplo de request
```json
{
  "name": "Loja Central",
  "document": "12345678000199",
  "owner_email": "admin@loja.com",
  "preferences": {
    "allow_negative_stock": true
  }
}
```

## 5. Exemplo de response
```json
{
  "id": "cmp-100",
  "schema": "tenant_abc",
  "status": "trial",
  "preferences": {
    "allow_negative_stock": true
  }
}
```

## 6. Notas e compatibilidade
- Nunca retornar tokens ou segredos neste DTO.
