# Domain / Preference

Estruturas auxiliares para armazenar preferências configuráveis (tema, impressão, flags).

---

## 1. Entidades principais
| Nome | Descrição |
|------|-----------|
| PreferenceValue | Chave, valor e escopo. |

## 2. Regras de negócio
- Valores tipados (bool, int, string, json).
- Suporta fallback para defaults globais.

## 3. Interações e consumidores
- Usecases: company, order, print_manager.

## 4. Exemplo de estrutura
```json
{
  "key": "receipt.template",
  "value": "kitchen_v2",
  "scope": "company"
}
```
