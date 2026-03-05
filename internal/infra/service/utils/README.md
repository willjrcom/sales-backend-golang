# Service / Utils

Funções de apoio (conversão, tempo, formatação) usadas por múltiplas camadas.

---

## 1. Responsabilidades
- Conversões de decimal, datas, máscaras.
- Gerar IDs amigáveis para relatórios.

## 2. Métodos principais
| Assinatura | Descrição |
|------------|-----------|
| `ParseDecimal(string) decimal.Decimal` | Garante precisão para valores monetários. |
| `TruncateWithEllipsis(s string, limit int) string` | Formata strings longas para impressão. |

## 3. Fluxo típico
- Chamada direta, sem dependências externas.

## 4. Configuração / Env Vars
- (sem configuração específica)

## 5. Exemplo de uso
```go
go
value := utils.ParseDecimal("12.34")
```

## 6. Falhas comuns
- ErrInvalidDecimal

## 7. Notas operacionais
- Evitar adicionar regras de negócio aqui para não criar dependências cíclicas.
