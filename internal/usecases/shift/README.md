# Usecase / Shift

Controla abertura/fechamento de turnos, consolida vendas por atendente e calcula comissões.

---

## 1. Pontos de entrada
| Método | Rota | Origem | Descrição |
|--------|------|--------|-----------|
| POST | `/shift/open` | handler/shift.go | Abre turno. |
| POST | `/shift/close` | handler/shift.go | Fecha turno e gera resumo. |
| GET | `/shift/{id}` | handler/shift.go | Consulta detalhada. |

## 2. Dependências
- Repositories: shift, order, employee, delivery_driver_tax.
- Services: print_manager (relatórios).

## 3. Fluxos e exemplos
### Abrir turno
Passos:
- Valida que não existe turno aberto para o funcionário.
- Registra caixa inicial e disponibiliza turno para pedidos.

Exemplo de request:
```json
{
  "employee_id": "emp-1",
  "cash_float": 200
}
```
Resposta:
```json
{
  "shift_id": "shf-1",
  "status": "open"
}
```

### Fechar turno
Passos:
- Gera resumo de pagamentos, gorjetas e taxas.
- Registra divergências e dispara impressão.

Exemplo de request:
```json
{
  "shift_id": "shf-1",
  "cash_counted": 380
}
```
Resposta:
```json
{
  "shift_id": "shf-1",
  "status": "closed",
  "cash_diff": -20
}
```

## 4. Falhas conhecidas
- ErrShiftAlreadyOpen
- ErrShiftMismatchTotals

## 5. Notas operacionais
- Fechamento deve bloquear novos pedidos para o atendente até reabrir turno.
