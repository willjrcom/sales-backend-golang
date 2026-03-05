# Infra / Modules

Responsável por montar cada módulo do sistema (repositórios, serviços e handlers) e registrá-los no servidor. Arquivo principal: `main.go`.

## Fluxo

1. Construtores `New<Context>Module` instanciam repositórios (Bun), serviços auxiliares e handlers.
2. Dependências cruzadas entre usecases são resolvidas via métodos `AddDependencies`.
3. `MainModules` executa todos os construtores em ordem para garantir que serviços críticos (estoque, pedidos, fila) recebam suas dependências antes de iniciar o servidor.

## Orientações

- Ao adicionar um novo contexto, crie `NewXModule` que retorna `(repository, service, handler)`.
- Registre quaisquer consumidores assíncronos (ex.: email worker) aqui para que iniciem junto ao servidor.
- Evite lógica de negócio; mantenha este arquivo apenas como composição.

