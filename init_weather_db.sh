#!/bin/bash

# Set variables
DB_USER="weather_user"
DB_PASSWORD="pa$$w0rd"
DB_NAME="weather"

echo "🔐 Creating PostgreSQL user and database..."

# Execute SQL commands
psql -U postgres <<EOF
DO \$\$
BEGIN
    IF NOT EXISTS (
        SELECT FROM pg_catalog.pg_user WHERE usename = '${DB_USER}'
    ) THEN
        CREATE USER ${DB_USER} WITH PASSWORD '${DB_PASSWORD}';
    END IF;
END
\$\$;

CREATE DATABASE ${DB_NAME} OWNER ${DB_USER};
GRANT ALL PRIVILEGES ON DATABASE ${DB_NAME} TO ${DB_USER};
EOF

echo "✅ Database '${DB_NAME}' and user '${DB_USER}' created."
