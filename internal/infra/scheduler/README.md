# Infra / Scheduler

Jobs recorrentes e workers que rodam dentro do processo HTTP. Exemplo atual: `daily_scheduler.go`, responsável por validar alertas de estoque, expirações e disparar notificações.

## Boas práticas

- Sempre execute jobs com `context.Context` e respeite multi-tenant (iterando schemas).
- Use canais/cron lightweight; evitar bloqueios longos no thread principal.
- Logs devem indicar o schema e a tarefa executada.

Novos agendamentos devem ser registrados aqui descrevendo periodicidade e dependências.

