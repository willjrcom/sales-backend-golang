# Domain / Employee

Modela funcionários, cargos e pagamentos.

---

## 1. Entidades principais
| Nome | Descrição |
|------|-----------|
| Employee | Dados funcionais e status. |
| EmployeePayment | Histórico de pagamentos e comissões. |

## 2. Regras de negócio
- Relaciona-se a `Person` e opcionalmente `User`.
- Permite bloqueio temporário mantendo histórico.
- Pagamentos ligados a turno/ordem para conciliação.

## 3. Interações e consumidores
- Usecases: employee, shift, delivery_driver, user.

## 4. Exemplo de estrutura
```json
{
  "id": "emp-10",
  "person_id": "per-20",
  "role": "cashier",
  "status": "active"
}
```
