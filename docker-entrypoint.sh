#!/bin/sh
set -e

# Executa todas as migrações pendentes automaticamente
# A menos que AUTO_MIGRATE=false seja definido
if [ "$AUTO_MIGRATE" != "false" ]; then
    echo "🔄 Running pending migrations..."
    /app/sales-backend migrate-all || {
        echo "⚠️  Migration failed, but continuing..."
    }
    echo "✅ Migrations complete"
fi

# Executa o comando passado (ou o CMD padrão)
exec /app/sales-backend "$@"
