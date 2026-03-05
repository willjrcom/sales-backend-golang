# Infra / Handler

Responsável por expor os casos de uso via HTTP (Chi). Cada arquivo monta um sub-router específico, converte DTOs e delega as regras de negócio para a camada `usecases`.

---

## 1. Pipeline padrão

1. Instanciar `chi.Router`.
2. Aplicar middlewares locais (validação de schema, RBAC, métricas).
3. Fazer bind/validação do payload (`dto.<Context>`).
4. Chamar o usecase correspondente e traduzir o resultado em resposta HTTP.
5. Registrar o router com `handlerimpl.NewHandler(path, router, unprotectedRoutes...)`.

---

## 2. Principais handlers

| Arquivo | Prefixo base | Usecases chamados | Endpoints principais |
|---------|--------------|-------------------|----------------------|
| `order.go` | `/order` | `order`, `stock`, `checkout` | `POST /order`, `POST /order/{id}/items`, `POST /order/{id}/status`, `DELETE /order/{id}` |
| `item.go` | `/order/{id}/items` | `item`, `stock` | cria/edita/remove itens vinculados ao pedido |
| `group_item.go` | `/group-item` | `group_item`, `order_process` | abre/fecha grupos e sincroniza status da cozinha |
| `stock.go` | `/stock` | `stock` | movimentos manuais, alertas e relatórios de estoque |
| `checkout.go` | `/checkout` | `checkout`, `order`, `company` | cálculo, confirmação e webhooks de pagamento |
| `report.go` | `/report` | `report` | dashboards de vendas/estoque/marketing |
| `product.go` | `/product` | `product`, `stock` | CRUD produtos/variações/tamanhos |
| `product_category.go` | `/product/categories` | `product_category`, `process_rule` | hierarquia + regras de processo |
| `company.go` | `/company` | `company`, `schema`, `user` | onboarding, preferências, assinatura |
| `client.go` | `/clients` | `client`, `contact`, `address` | CRUD clientes + histórico |
| `delivery_driver.go` | `/delivery-drivers` | `delivery_driver`, `employee` | cadastro, status, redistribuição |
| `order_queue.go` | `/order/queue` | `order_queue`, `order_process` | avanço/reatribuição nas filas |
| `order_table.go` | `/tables` | `order_table`, `table`, `order` | abrir/mover/fechar mesas |
| `place.go` | `/places` | `place`, `table` | locais físicos, canais e horários |
| `fiscal_invoice.go` | `/fiscal-invoice` | `fiscal_invoice`, `fiscal_settings` | emitir/cancelar/consultar NF-e |
| `user.go` | `/users` | `user`, `employee`, `auth` | CRUD usuários, reset senha, roles |
| `s3.go` | `/storage` | `s3` | geração de URLs pré-assinadas |
| `public.go` | `/public` | `company`, `user` | endpoints sem autenticação (lookup, forgot password) |

> Dica: mantenha o nome do arquivo alinhado com o prefixo base; isso facilita localizar o handler correto.

---

## 3. Tratamento de erros

- Utilize `httperrors.From(err)` para padronizar status e payload.
- Sempre logue IDs críticos (`order_id`, `company_id`, `user_id`) antes de retornar erro.
- Para validações complexas retorne 422 com detalhes (`dto.ValidationError`).
- Panics são capturados por `RecoverMiddleware`, mas os handlers devem evitar panics previsíveis validando entrada.

---

## 4. Rotas públicas

- Webhooks (`/checkout/webhook/mercadopago`, `/public/users/forgot-password`) precisam ser listados em `UnprotectedRoutes`.
- Mesmo em rotas públicas, valide assinaturas HMAC e sanitize os headers antes de logar.

---

## 5. Boas práticas

- Handler **não** contém regra de negócio; mantenha a camada fina.
- Reutilize os DTOs da pasta `internal/infra/dto` e mantenha-os sincronizados com o frontend.
- Propague `context.Context` do Chi para o usecase para preservar tracing/schema.
- Para streams longos (fila/cozinha) considere SSE/WebSocket dedicado para não bloquear o worker HTTP.

---

## 3. Detalhes por handler
### `order.go` — prefixo `/order`
Usecases: order, stock, checkout, group_item
Endpoints:
- `POST /order` — Cria pedido em modo draft e retorna queue_number.
- `POST /order/{id}/items` — Adiciona itens/adicionais e dispara reserva de estoque.
- `POST /order/{id}/status` — Transiciona status (pending → in_progress → finished).
- `DELETE /order/{id}` — Cancela pedido, restaura estoque e estorna pagamentos.
Notas:
- Propaga `context.Context` com schema e usuário logado para auditoria.
- Utiliza DTOs `order`, `item`, `group_item`; não aceita payload fora desses contratos.

### `item.go` — prefixo `/order/{id}/items`
Usecases: item, stock
Endpoints:
- `POST /order/{id}/items` — Cria item isolado (usado quando front envia um item por vez).
- `PUT /order/{id}/items/{item_id}` — Atualiza quantidade/notas de um item.
- `DELETE /order/{id}/items/{item_id}` — Remove item e chama `RestoreFromItem`.
Notas:
- Sempre valida se o pedido permite edição (status < finished).

### `stock.go` — prefixo `/stock`
Usecases: stock
Endpoints:
- `GET /stock/{id}` — Consulta estoque e lotes.
- `POST /stock/{id}/movement/add` — Entrada manual (IN) com custo opcional.
- `POST /stock/{id}/movement/remove` — Saída manual (OUT) FIFO.
- `POST /stock/{id}/movement/adjust` — Ajuste inventário.
- `GET /stock/{id}/movement` — Histórico paginado.
- `GET /stock/alerts` — Alertas ativos.
- `POST /stock/report` — Relatório consolidado por produto.
Notas:
- Todos os endpoints exigem schema e role com permissão `stock:write` ou `stock:read`.

### `checkout.go` — prefixo `/checkout`
Usecases: checkout, order, company
Endpoints:
- `POST /checkout/calculate` — Calcula totais, troco e taxas antes de confirmar.
- `POST /checkout/confirm` — Confirma pagamento local/online.
- `POST /checkout/webhook/mercadopago` — Recebe webhooks de pagamento (rota pública).
Notas:
- Webhook está em `UnprotectedRoutes` e valida assinatura HMAC.
- Quando confirmar POS, o handler aguarda resposta do terminal e retorna 202 em caso de timeout.

### `report.go` — prefixo `/report`
Usecases: report
Endpoints:
- `POST /report/sales-summary` — Série temporal de vendas.
- `POST /report/additional-items-sold` — Top adicionais vendidos.
- `POST /report/complements-sold` — Top complementos (limit 10).
- `POST /report/top-products` — Produtos mais vendidos (limit configurável).
Notas:
- Para filtros grandes, aplica limite de 92 dias; retorna 422 se exceder.

### `product.go` — prefixo `/product`
Usecases: product, stock
Endpoints:
- `GET /product` — Listagem paginada com filtros por categoria/canal.
- `POST /product` — Cria produto completo com variações e tamanhos.
- `PUT /product/{id}` — Atualiza atributos e disponibilidade.
- `DELETE /product/{id}` — Soft delete (marca como inactive).
Notas:
- Quando `track_stock=true`, handler garante criação de estoque chamando `StockUsecase`.

### `product_category.go` — prefixo `/product/categories`
Usecases: product_category, process_rule
Endpoints:
- `GET /product/categories` — Retorna árvore de categorias.
- `POST /product/categories` — Cria categoria.
- `PUT /product/categories/{id}` — Atualiza hierarquia/process rules.
- `DELETE /product/categories/{id}` — Desativa categoria mantendo histórico.
Notas:
- Aplica validação de loops antes de chamar o usecase.

### `company.go` — prefixo `/company`
Usecases: company, schema, user
Endpoints:
- `POST /company` — Onboarding completo (empresa + schema + usuário admin).
- `PUT /company/{id}` — Atualiza dados cadastrais.
- `POST /company/{id}/preferences` — Atualiza preferências operacionais.
- `POST /company/{id}/subscription/activate` — Ativa/desativa assinatura.
Notas:
- Alguns endpoints executam no schema público antes de configurar search_path.

### `client.go` — prefixo `/clients`
Usecases: client, contact, address
Endpoints:
- `GET /clients` — Listagem com filtros por documento/telefone/status.
- `POST /clients` — Cria cliente + person/contact/address em uma tx.
- `PUT /clients/{id}` — Atualiza dados e flags de bloqueio.
- `GET /clients/{id}/history` — Retorna pedidos/ticket médio.
Notas:
- Documentos e telefones são normalizados antes de chamar o usecase.

### `delivery_driver.go` — prefixo `/delivery-drivers`
Usecases: delivery_driver, employee
Endpoints:
- `GET /delivery-drivers` — Lista drivers + status em tempo real.
- `POST /delivery-drivers` — Cadastra entregador vinculado a employee.
- `PATCH /delivery-drivers/{id}/status` — Atualiza disponibilidade/veículo/zona.
Notas:
- Retorna 409 se tentar marcar como `available` enquanto possui pedido em andamento.

### `user.go` — prefixo `/users`
Usecases: user, employee, auth
Endpoints:
- `POST /users` — Cria usuário com roles iniciais.
- `GET /users` — Listagem com filtros por role/status.
- `POST /users/{id}/reset-password` — Dispara reset (webhook/email).
- `POST /users/{id}/roles` — Atualiza papéis RBAC.
Notas:
- Reset password é unprotected, mas exige token temporário.

---

## 3. Detalhes por handler
### `advertising.go` — prefixo `/advertising`
Usecases: advertising, s3
Endpoints:
- `GET /advertising` — Lista campanhas ativas com filtros por categoria/patrocinador.
- `POST /advertising` — Cria campanha e agenda publicação.
- `PUT /advertising/{id}` — Atualiza vigência, assets e segmentação.
- `DELETE /advertising/{id}` — Arquiva campanha e dispara evento de remoção.
Notas:
- Uploads grandes usam URLs pré-assinadas geradas pelo serviço S3.
- Handler valida exclusividade do slot antes de chamar o usecase.

### `company.go` — prefixo `/company`
Usecases: company, schema, user
Endpoints:
- `POST /company` — Onboarding completo (empresa + schema + usuário admin).
- `PUT /company/{id}` — Atualiza dados cadastrais e status.
- `POST /company/{id}/preferences` — Atualiza preferências operacionais.
- `POST /company/{id}/subscription/activate` — Controla assinatura e billing.
Notas:
- Operações sensíveis acontecem no schema público antes de configurar search_path.

### `company_category.go` — prefixo `/company/{id}/categories`
Usecases: company_category, sponsor, advertising
Endpoints:
- `GET /company/{id}/categories` — Lista categorias atribuídas e patrocinadores ativos.
- `POST /company/{id}/categories` — Substitui conjunto de categorias da empresa.
- `POST /company/{id}/categories/{category_id}/sponsors` — Vincula patrocinador à categoria.
- `DELETE /company/{id}/categories/{category_id}/sponsors/{sponsor_id}` — Remove vínculo de patrocínio.
Notas:
- Aplica bloqueio otimista para evitar conflitos entre múltiplos operadores.

### `contact.go` — prefixo `/contacts`
Usecases: contact
Endpoints:
- `POST /contacts` — Cria contato (telefone/email) ligado a uma pessoa.
- `PUT /contacts/{id}` — Atualiza valor ou flag `is_primary`.
- `DELETE /contacts/{id}` — Remove contato (mantém registro auditável).
Notas:
- Normaliza telefone/email antes de persisitir.
- Retorna 409 se já existir contato primário para o tipo.

### `client.go` — prefixo `/clients`
Usecases: client, contact, address
Endpoints:
- `GET /clients` — Listagem com filtros por documento/telefone/status.
- `POST /clients` — Cria cliente completo (person/contact/address).
- `PUT /clients/{id}` — Atualiza dados, bloqueios e preferências.
- `GET /clients/{id}/history` — Retorna histórico de pedidos e métricas.
Notas:
- Documento e telefone são mascarados nas respostas para proteger dados.

### `delivery_driver.go` — prefixo `/delivery-drivers`
Usecases: delivery_driver, employee
Endpoints:
- `GET /delivery-drivers` — Lista entregadores com status em tempo real.
- `POST /delivery-drivers` — Cadastra entregador vinculado a employee.
- `PATCH /delivery-drivers/{id}/status` — Atualiza disponibilidade, veículo ou zona.
Notas:
- Retorna 409 se driver tentar ficar disponível com pedido em andamento.

### `employee.go` — prefixo `/employees`
Usecases: employee, user
Endpoints:
- `POST /employees` — Cria funcionário (opcionalmente cria usuário).
- `PUT /employees/{id}` — Atualiza cargo, jornada, permissões.
- `POST /employees/{id}/payment` — Registra pagamento/comissão.
Notas:
- Quando `create_user=true`, handler delega criação ao usecase user.

### `fiscal_invoice.go` — prefixo `/fiscal-invoice`
Usecases: fiscal_invoice, fiscal_settings
Endpoints:
- `POST /fiscal-invoice` — Gera NF-e para pedido.
- `POST /fiscal-invoice/{id}/cancel` — Cancela NF-e com justificativa.
- `GET /fiscal-invoice/{id}` — Consulta status, XML/PDF assinados.
Notas:
- Sempre valida se empresa tem configuração fiscal completa antes de emitir.

### `fiscal_settings.go` — prefixo `/fiscal/settings`
Usecases: fiscal_settings, company
Endpoints:
- `GET /fiscal/settings` — Retorna configuração atual.
- `POST /fiscal/settings` — Atualiza CRT, CNAE, certificado.
- `POST /fiscal/settings/test` — Executa teste de comunicação com provedor.
Notas:
- Payloads sensíveis são mascarados nos logs.

### `group_item.go` — prefixo `/group-item`
Usecases: group_item, order_process
Endpoints:
- `POST /group-item` — Cria grupo para itens (ex.: combos).
- `POST /group-item/{id}/status` — Atualiza status do grupo (em preparo, pronto).
- `DELETE /group-item/{id}` — Cancela grupo e restaura estoque dos itens.
Notas:
- Dispara eventos para a cozinha imprimir novamente quando status muda.

### `order_delivery.go` — prefixo `/order/{id}/delivery`
Usecases: order, order_delivery, delivery_driver
Endpoints:
- `POST /order/{id}/delivery` — Configura endereço e driver.
- `PATCH /order/{id}/delivery/status` — Atualiza status e ETA.
- `POST /order/{id}/delivery/assign` — Reatribui driver manualmente.
Notas:
- Valida CEP e zona de entrega antes de confirmar.

### `order_pickup.go` — prefixo `/order/{id}/pickup`
Usecases: order, order_pickup
Endpoints:
- `POST /order/{id}/pickup` — Define horário e contato de retirada.
- `PATCH /order/{id}/pickup/status` — Atualiza status (waiting, ready, delivered).
Notas:
- Gera `pickup_code` e devolve no payload para conferência.

### `order_print.go` — prefixo `/print`
Usecases: print_manager, order, shift
Endpoints:
- `POST /print/order` — Imprime pedido completo.
- `POST /print/kitchen` — Imprime tickets por estação.
- `POST /print/shift` — Imprime fechamento de turno.
Notas:
- Executa fora do request principal usando goroutines para não bloquear.

### `order_process.go` — prefixo `/order/process`
Usecases: order_process, order_queue
Endpoints:
- `GET /order/process` — Lista processos e status.
- `POST /order/process/{id}/status` — Atualiza etapa/funcionário responsável.
- `POST /order/process/{id}/queue` — Enfileira itens manualmente.
Notas:
- Requer permissão específica `process:write`.

### `order_queue.go` — prefixo `/order/queue`
Usecases: order_queue, order_process
Endpoints:
- `GET /order/queue` — Lista filas por processo com posição.
- `POST /order/{id}/queue/advance` — Avança pedido para próxima etapa.
- `POST /order/{id}/queue/reassign` — Reatribui responsável/processo.
Notas:
- Eventos são transmitidos por SSE/WebSocket para painéis.

### `order_table.go` — prefixo `/tables`
Usecases: order_table, table, order
Endpoints:
- `POST /tables/{id}/open` — Abre mesa e cria pedido dine-in.
- `POST /tables/{id}/move` — Transfere pedido para outra mesa.
- `POST /tables/{id}/close` — Fecha mesa e consolida pagamentos.
Notas:
- Mantém lock na mesa durante abertura para evitar doble booking.

### `place.go` — prefixo `/places`
Usecases: place, table
Endpoints:
- `GET /places` — Lista locais com horários e canais.
- `POST /places` — Cria novo local.
- `PUT /places/{id}` — Atualiza canais/horários/capacidade.
Notas:
- Quando canais mudam, invalida caches de SLA de entrega.

### `process_rule.go` — prefixo `/process-rules`
Usecases: process_rule
Endpoints:
- `GET /process-rules` — Lista regras ativas.
- `POST /process-rules` — Cria regra com etapas.
- `PUT /process-rules/{id}` — Atualiza sequência de etapas.
- `DELETE /process-rules/{id}` — Desativa regra mantendo histórico.
Notas:
- Atualizações impactam a cozinha em tempo real; peça confirmação ao usuário.

### `product.go` — prefixo `/product`
Usecases: product, stock
Endpoints:
- `GET /product` — Listagem paginada com filtros por categoria/canal.
- `POST /product` — Cria produto completo com variações e tamanhos.
- `PUT /product/{id}` — Atualiza atributos e disponibilidade.
- `DELETE /product/{id}` — Soft delete (status inactive).
Notas:
- Quando `track_stock=true`, cria registro em stock automaticamente.

### `product_category.go` — prefixo `/product/categories`
Usecases: product_category, process_rule
Endpoints:
- `GET /product/categories` — Retorna árvore completa.
- `POST /product/categories` — Cria categoria.
- `PUT /product/categories/{id}` — Atualiza hierarquia/process rule.
- `DELETE /product/categories/{id}` — Desativa categoria.
Notas:
- Valida loops antes de salvar.

### `public.go` — prefixo `/public`
Usecases: company, user
Endpoints:
- `GET /public/companies` — Lista empresas disponíveis para login.
- `POST /public/users/forgot-password` — Inicia fluxo de reset (sem auth).
Notas:
- Precisa tratar rate-limit/abuso pois não exige token.

### `report.go` — prefixo `/report`
Usecases: report
Endpoints:
- `POST /report/sales-summary` — Série temporal de vendas.
- `POST /report/additional-items-sold` — Top adicionais.
- `POST /report/complements-sold` — Top complementos (limit 10).
- `POST /report/top-products` — Produtos mais vendidos.
Notas:
- Limita range a 92 dias; retorna 422 se exceder.

### `s3.go` — prefixo `/storage`
Usecases: s3
Endpoints:
- `POST /storage/presign` — Gera URL pré-assinada para upload.
- `DELETE /storage/object` — Remove asset por chave (quando permitido).
Notas:
- Retorna apenas URLs HTTPS; chave inclui company_id para isolamento.

### `shift.go` — prefixo `/shift`
Usecases: shift, order, print_manager
Endpoints:
- `POST /shift/open` — Abre turno e registra caixa inicial.
- `POST /shift/close` — Fecha turno e gera resumo.
- `GET /shift/{id}` — Detalhes do turno.
Notas:
- Ao fechar, dispara impressão automática via `print_manager`.

### `size.go` — prefixo `/products/{product_id}/sizes`
Usecases: size, product
Endpoints:
- `POST /products/{product_id}/sizes` — Cria tamanho.
- `PUT /products/{product_id}/sizes/{size_id}` — Atualiza preço/status.
- `DELETE /products/{product_id}/sizes/{size_id}` — Inativa tamanho.
Notas:
- Valida SKU suffix único por produto.

### `sponsor.go` — prefixo `/sponsors`
Usecases: sponsor, advertising
Endpoints:
- `GET /sponsors` — Lista patrocinadores ativos.
- `POST /sponsors` — Cadastra patrocinador e vigência.
- `PUT /sponsors/{id}` — Atualiza contrato/benefícios.
- `POST /sponsors/{id}/assets` — Cadastra logos/benefícios extras.
Notas:
- Uploads usam serviço S3; remover ativos órfãos ao deletar.

### `stock.go` — prefixo `/stock`
Usecases: stock
Endpoints:
- `GET /stock/{id}` — Consulta estoque e lotes.
- `POST /stock/{id}/movement/add` — Entrada manual (IN) com custo opcional.
- `POST /stock/{id}/movement/remove` — Saída manual (OUT) FIFO.
- `POST /stock/{id}/movement/adjust` — Ajuste inventário.
- `GET /stock/{id}/movement` — Histórico paginado.
- `GET /stock/alerts` — Alertas ativos.
- `POST /stock/report` — Relatório consolidado por produto.
Notas:
- Requer permissão específica (`stock:read` ou `stock:write`).

### `table.go` — prefixo `/tables`
Usecases: table, place
Endpoints:
- `GET /tables` — Lista mesas por local.
- `POST /tables` — Cria mesa física/virtual.
- `PUT /tables/{id}` — Atualiza capacidade/status.
- `DELETE /tables/{id}` — Inativa mesa (se não houver pedidos abertos).
Notas:
- Valida se a mesa está livre antes de excluir/inativar.

### `order.go` — prefixo `/order`
Usecases: order, stock, checkout, group_item
Endpoints:
- `POST /order` — Cria pedido em modo draft e retorna queue_number.
- `POST /order/{id}/items` — Adiciona itens e dispara reserva de estoque.
- `POST /order/{id}/status` — Transiciona status (pending → in_progress → finished).
- `DELETE /order/{id}` — Cancela pedido, restaura estoque e estorna pagamentos.
Notas:
- Propaga `context.Context` com schema/usuário para auditoria.
- Utiliza DTOs `order`, `item`, `group_item`.

### `item.go` — prefixo `/order/{id}/items`
Usecases: item, stock
Endpoints:
- `POST /order/{id}/items` — Cria item isolado quando o front envia um por vez.
- `PUT /order/{id}/items/{item_id}` — Atualiza quantidade e observações.
- `DELETE /order/{id}/items/{item_id}` — Remove item e chama `RestoreFromItem`.
Notas:
- Valida se o pedido permite edição (status < finished).

