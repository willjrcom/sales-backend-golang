# DTO / Contact

DTOs para contatos (telefone/email) independentes.

---

## 1. Onde é usado
- handler/contact.go

## 2. Estruturas principais
| Struct | Campos principais | Direção |
|--------|-------------------|---------|
| ContactRequest | person_id, type, value, is_primary | request |
| ContactResponse | id, person_id, type, value, is_primary | response |

## 3. Regras de validação
- type ∈ {phone,email}.
- phone formato E.164; email validado com regex.

## 4. Exemplo de request
```json
{
  "person_id": "per-1",
  "type": "phone",
  "value": "+5511988887777",
  "is_primary": true
}
```

## 5. Exemplo de response
```json
{
  "id": "con-1",
  "person_id": "per-1",
  "type": "phone",
  "value": "+5511988887777",
  "is_primary": true
}
```

## 6. Notas e compatibilidade
- `is_primary` deve ser único por tipo/pessoa.
