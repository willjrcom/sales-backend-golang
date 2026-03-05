# Bootstrap / RBAC

Implementação mínima de RBAC (Role Based Access Control) usada como utilitário pelos handlers para definir permissões de rota.

## Conceitos

- **Role**: representa o papel atrelado ao usuário (ex.: `admin`, `cashier`, `viewer`).
- **Resource**: string abstrata normalmente alinhada ao módulo (`order`, `stock`, `report`).
- **AccessLevel**: enum `NoAccess`, `Read`, `Write`, `Admin` que facilita comparações.
- **RBAC**: mapa em memória que relaciona role → resource → access level.

## Como aplicar

1. Configure as permissões em tempo de boot (ex.: dentro de `internal/infra/modules` ou de um middleware dedicado).
2. Ao autenticar o usuário, associe os roles correspondentes (vindos da empresa) e injete no contexto.
3. Chame `rbac.CanAccess(user, resource)` antes de executar o handler para bloquear acessos.

Embora simples, este pacote centraliza a semântica de acesso e evita duplicação de lógica entre handlers.

