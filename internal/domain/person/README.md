# Domain / Person

Dados pessoais compartilhados (nome, documento, contatos).

---

## 1. Entidades principais
| Nome | Descrição |
|------|-----------|
| Person | Nome, documento, data nascimento. |
| Contact | Telefone/email principal. |

## 2. Regras de negócio
- Documento único por schema.
- Campos sensíveis (CPF) criptografados em repouso.

## 3. Interações e consumidores
- Usecases: client, employee, user.

## 4. Exemplo de estrutura
```json
{
  "id": "per-1",
  "name": "Maria",
  "document": "12345678909",
  "contacts": [
    {
      "type": "phone",
      "value": "+5511988887777"
    }
  ]
}
```
