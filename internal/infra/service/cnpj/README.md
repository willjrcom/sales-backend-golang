# Service / CNPJ

Valida e normaliza CNPJ para cadastros de empresa.

---

## 1. Responsabilidades
- Remover máscara e validar dígitos verificadores.
- Formatar para exibição (##.###.###/####-##).
- Fornecer helpers usados em company onboarding e integrações fiscais.

## 2. Métodos principais
| Assinatura | Descrição |
|------------|-----------|
| `Normalize(raw string) (string, error)` | Remove caracteres não numéricos e valida tamanho. |
| `IsValid(cnpj string) bool` | Executa cálculo de dígito verificador. |
| `Format(cnpj string) string` | Aplica máscara padrão para UI. |

## 3. Fluxo típico
- Recebe documento do formulário.
- Normalize → valida tamanho → calcula DV.
- Se válido, retorna string com 14 dígitos para persistência.
- Usecases chamam `Format` apenas para resposta ao frontend.

## 4. Configuração / Env Vars
- (sem configuração específica)

## 5. Exemplo de uso
```go
go
normalized, err := cnpj.Normalize(input)
if err != nil {
    return ErrInvalidCNPJ
}
if !cnpj.IsValid(normalized) {
    return ErrInvalidCNPJ
}
company.Document = normalized
```

## 6. Falhas comuns
- ErrInvalidLength
- ErrInvalidDigit

## 7. Notas operacionais
- Não faça chamadas externas aqui; validação é puramente local.
