# `internal/usecases`

Camada responsável por orquestrar regras de negócio usando as entidades do domínio e as interfaces providas pela infraestrutura. Cada subpasta mantém um README com detalhes do fluxo específico.

## Módulos disponíveis

advertising · checkout · client · company · company_category · contact · delivery_driver · employee · fiscal_invoice · fiscal_settings · order · order_queue · order_table · place · print_manager · process_rule · product · product_category · report · shift · size · sponsor · stock · table · user

## Convenção

1. Usecases expõem construtores que recebem interfaces, não implementações concretas.
2. Métodos `AddDependencies` são usados para injetar dependências circulares sem quebrar testes.
3. Toda validação crítica (quantidades, permissões, status) deve acontecer aqui antes de chamar infraestrutura.
4. Testes unitários residem na mesma pasta (`*_test.go`) e usam fakes/mocks leves.

