# Usecase / Contact

Centraliza contatos (telefone/email) independentes de clientes ou funcionários, evitando duplicidade e permitindo reuso.

---

## 1. Pontos de entrada
| Método | Rota | Origem | Descrição |
|--------|------|--------|-----------|
| POST | `/contacts` | handler/contact.go | Cria contato e associa a uma pessoa. |
| PUT | `/contacts/{id}` | handler/contact.go | Atualiza telefone/email e flags de principal. |
| DELETE | `/contacts/{id}` | handler/contact.go | Remove contato preservando histórico em entidades dependentes. |

## 2. Dependências
- Repositories: contact, person.

## 3. Fluxos e exemplos
### Criar contato com validação
Passos:
- Normaliza formatos (E.164 para telefone, lowercase para email).
- Verifica se já existe contato igual e reusa ID.
- Associa a entidade (client/employee).

Exemplo de request:
```json
{
  "person_id": "per-1",
  "type": "phone",
  "value": "+5511988889999",
  "is_primary": true
}
```
Resposta:
```json
{
  "id": "con-500",
  "person_id": "per-1",
  "type": "phone",
  "value": "+5511988889999"
}
```

## 4. Falhas conhecidas
- ErrContactDuplicate
- ErrContactInUse

## 5. Notas operacionais
- Nunca apague contatos que estejam marcados como principal sem antes promover outro.
