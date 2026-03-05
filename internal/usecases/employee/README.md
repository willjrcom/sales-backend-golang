# Usecase / Employee

Mantém funcionários, vínculos com usuários do sistema e controle de pagamentos/comissões.

---

## 1. Pontos de entrada
| Método | Rota | Origem | Descrição |
|--------|------|--------|-----------|
| POST | `/employees` | handler/employee.go | Cria funcionário e opcionalmente usuário de acesso. |
| PUT | `/employees/{id}` | handler/employee.go | Atualiza cargo, jornada, permissões. |
| POST | `/employees/{id}/payment` | handler/employee.go | Registra pagamento ou comissão. |

## 2. Dependências
- Repositories: employee, person, contact, user, employee_payment.
- Services: email (convite de acesso).

## 3. Fluxos e exemplos
### Onboard funcionário
Passos:
- Valida CPF e dados pessoais.
- Cria employee + contato + vínculo com company.
- Se `create_user=true`, delega ao usecase user.

Exemplo de request:
```json
{
  "name": "João",
  "document": "11122233344",
  "role": "cashier",
  "create_user": true
}
```
Resposta:
```json
{
  "id": "emp-10",
  "user_id": "usr-55",
  "status": "active"
}
```

### Registrar pagamento
Passos:
- Calcula valor líquido considerando adiantamentos.
- Cria `employee_payment` com referência ao shift.
- Atualiza saldo devedor/credor.

Exemplo de request:
```json
{
  "amount": 300,
  "shift_id": "shf-1",
  "type": "commission"
}
```
Resposta:
```json
{
  "payment_id": "pay-emp-1",
  "status": "recorded"
}
```

## 4. Falhas conhecidas
- ErrEmployeeDuplicate
- ErrMissingShift

## 5. Notas operacionais
- Desativar funcionário não remove vínculos com pedidos passados; apenas bloqueia escalação futura.
