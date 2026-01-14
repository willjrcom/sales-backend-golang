# Como Executar as Migrations

## Pré-requisitos

Certifique-se de que o banco de dados PostgreSQL está rodando e acessível.

## Executar Migrations

Execute as migrations na ordem correta:

```bash
# 1. Adicionar campos fiscais à tabela companies
psql -U your_user -d your_database -f scripts/migrations/001_add_fiscal_fields_to_companies.sql

# 2. Criar tabela company_usage_costs
psql -U your_user -d your_database -f scripts/migrations/002_create_company_usage_costs_table.sql

# 3. Criar tabela fiscal_invoices
psql -U your_user -d your_database -f scripts/migrations/003_create_fiscal_invoices_table.sql
```

## Verificar Migrations

```bash
# Verificar se as tabelas foram criadas
psql -U your_user -d your_database -c "\dt"

# Verificar estrutura da tabela companies
psql -U your_user -d your_database -c "\d companies"

# Verificar estrutura da tabela company_usage_costs
psql -U your_user -d your_database -c "\d company_usage_costs"

# Verificar estrutura da tabela fiscal_invoices
psql -U your_user -d your_database -c "\d fiscal_invoices"
```

## Rollback (se necessário)

```bash
# Remover tabela fiscal_invoices
psql -U your_user -d your_database -c "DROP TABLE IF EXISTS fiscal_invoices CASCADE;"

# Remover tabela company_usage_costs
psql -U your_user -d your_database -c "DROP TABLE IF EXISTS company_usage_costs CASCADE;"

# Remover campos fiscais da tabela companies
psql -U your_user -d your_database -c "
ALTER TABLE companies 
DROP COLUMN IF EXISTS fiscal_enabled,
DROP COLUMN IF EXISTS inscricao_estadual,
DROP COLUMN IF EXISTS regime_tributario,
DROP COLUMN IF EXISTS cnae,
DROP COLUMN IF EXISTS crt,
DROP COLUMN IF EXISTS transmitenota_usuario,
DROP COLUMN IF EXISTS transmitenota_senha;
"
```

## Configuração Pós-Migration

Após executar as migrations:

1. **Adicione as variáveis de ambiente** ao seu `.env.local`:
   ```bash
   cat .env.transmitenota.example >> .env.local
   ```

2. **Configure os dados fiscais da empresa** via API ou diretamente no banco:
   ```sql
   UPDATE companies 
   SET 
     fiscal_enabled = true,
     inscricao_estadual = '123456789',
     regime_tributario = 1,
     cnae = '4711302',
     crt = 1,
     transmitenota_usuario = 'seu_usuario',
     transmitenota_senha = 'sua_senha'
   WHERE id = 'company_uuid_here';
   ```

3. **Teste a emissão** usando o sandbox da Transmitenota (`TRANSMITENOTA_SANDBOX=true`)
