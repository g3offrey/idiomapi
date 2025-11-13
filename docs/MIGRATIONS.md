# Database Migrations with Goose

This project uses [goose](https://github.com/pressly/goose) for database migration management.

## Why Goose?

✅ **Simple and reliable**: SQL-first approach  
✅ **Versioned migrations**: Track migration history  
✅ **Rollback support**: Easily revert changes  
✅ **Go embedded**: Can also use Go migrations  
✅ **Database agnostic**: Works with PostgreSQL, MySQL, SQLite, etc.  

## Installation

Goose is automatically installed when you run:

```bash
make install-tools
```

Or manually:

```bash
go install github.com/pressly/goose/v3/cmd/goose@latest
```

## Migration Commands

### Run All Pending Migrations

```bash
make migrate-up
```

This applies all pending migrations to bring the database to the latest version.

### Rollback Last Migration

```bash
make migrate-down
```

This reverts the most recently applied migration.

### Reset All Migrations

```bash
make migrate-reset
```

⚠️ **Warning**: This drops all tables and resets the database. Use with caution!

### Check Migration Status

```bash
make migrate-status
```

Shows which migrations have been applied:

```
Applied At                  Migration
=======================================
2024-01-15 10:30:00 UTC    00001_create_todos_table.sql
Pending                     00002_add_tags_to_todos.sql
```

### Create New Migration

```bash
make migrate-create NAME=add_users_table
```

This creates a new migration file in the `migrations/` directory.

## Migration File Format

Goose uses special comments to separate up and down migrations:

```sql
-- +goose Up
-- SQL for applying the migration
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    email VARCHAR(255) NOT NULL UNIQUE
);

-- +goose Down
-- SQL for reverting the migration
DROP TABLE IF EXISTS users;
```

### Key Points

1. **Up Section**: `-- +goose Up` defines forward migration
2. **Down Section**: `-- +goose Down` defines rollback
3. **Statement Blocks**: Use `-- +goose StatementBegin` and `-- +goose StatementEnd` for complex SQL
4. **Naming**: Files are named `XXXXX_description.sql` where XXXXX is a sequential number

## Example Migrations

### Simple Table Creation

**File**: `00002_create_users.sql`

```sql
-- +goose Up
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    email VARCHAR(255) NOT NULL UNIQUE,
    name VARCHAR(255) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

-- +goose Down
DROP TABLE IF EXISTS users;
```

### Complex Migration with Functions

**File**: `00003_add_audit_trigger.sql`

```sql
-- +goose Up
-- +goose StatementBegin
CREATE OR REPLACE FUNCTION audit_trigger()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER users_audit
    BEFORE UPDATE ON users
    FOR EACH ROW
    EXECUTE FUNCTION audit_trigger();
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TRIGGER IF EXISTS users_audit ON users;
DROP FUNCTION IF EXISTS audit_trigger();
-- +goose StatementEnd
```

### Adding Columns

**File**: `00004_add_user_role.sql`

```sql
-- +goose Up
ALTER TABLE users ADD COLUMN role VARCHAR(50) DEFAULT 'user';
CREATE INDEX idx_users_role ON users(role);

-- +goose Down
DROP INDEX IF EXISTS idx_users_role;
ALTER TABLE users DROP COLUMN IF EXISTS role;
```

### Data Migration

**File**: `00005_populate_default_data.sql`

```sql
-- +goose Up
INSERT INTO users (email, name, role) VALUES
    ('admin@example.com', 'Admin User', 'admin'),
    ('user@example.com', 'Regular User', 'user');

-- +goose Down
DELETE FROM users WHERE email IN ('admin@example.com', 'user@example.com');
```

## Advanced Usage

### Manual Migration

If you need more control, you can use goose directly:

```bash
# Apply specific number of migrations
goose -dir migrations postgres "postgresql://user:pass@localhost:5432/db" up-to 3

# Apply one migration at a time
goose -dir migrations postgres "postgresql://user:pass@localhost:5432/db" up-by-one

# Get current version
goose -dir migrations postgres "postgresql://user:pass@localhost:5432/db" version
```

### Environment-Specific Migrations

Use environment variables for different databases:

```bash
# Development
export DB_DSN="host=localhost port=5432 user=postgres password=postgres dbname=tododb_dev sslmode=disable"
make migrate-up

# Testing
export DB_DSN="host=localhost port=5432 user=postgres password=postgres dbname=tododb_test sslmode=disable"
make migrate-up

# Production
export DB_DSN="host=prod-db port=5432 user=prod_user password=secret dbname=tododb sslmode=require"
make migrate-up
```

## Migration Best Practices

### 1. Always Test Rollbacks

Before applying migrations in production, test both up and down:

```bash
make migrate-up
make migrate-down
make migrate-up
```

### 2. Make Migrations Idempotent

Use `IF EXISTS` and `IF NOT EXISTS`:

```sql
-- +goose Up
CREATE TABLE IF NOT EXISTS users (...);
CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);

-- +goose Down
DROP INDEX IF EXISTS idx_users_email;
DROP TABLE IF EXISTS users;
```

### 3. Keep Migrations Small

One logical change per migration:

- ✅ Good: `00001_create_users_table.sql`
- ✅ Good: `00002_add_users_email_index.sql`
- ❌ Bad: `00001_create_all_tables.sql`

### 4. Never Modify Applied Migrations

Once a migration is applied in any environment:

- ❌ Don't modify it
- ✅ Create a new migration to make changes

### 5. Use Descriptive Names

```bash
# Good names
make migrate-create NAME=add_email_verification_to_users
make migrate-create NAME=create_sessions_table
make migrate-create NAME=add_index_on_user_email

# Bad names
make migrate-create NAME=update_db
make migrate-create NAME=fix_bug
```

### 6. Handle Data Carefully

For data migrations:

- Backup data before migration
- Use transactions where possible
- Test on a copy of production data

### 7. Document Complex Migrations

Add comments explaining why:

```sql
-- +goose Up
-- Migration to support multi-tenancy
-- Related to issue #123 and design doc: docs/multi-tenant-design.md
ALTER TABLE todos ADD COLUMN tenant_id INTEGER REFERENCES tenants(id);
CREATE INDEX idx_todos_tenant ON todos(tenant_id);

-- +goose Down
DROP INDEX IF EXISTS idx_todos_tenant;
ALTER TABLE todos DROP COLUMN IF EXISTS tenant_id;
```

## Troubleshooting

### Migration Failed Partially

If a migration fails partway through:

1. Check the goose_db_version table to see current version
2. Manually fix any partially applied changes
3. Mark migration as failed in goose_db_version if needed
4. Fix the migration file
5. Re-run migration

### Rollback Failed

If rollback fails:

1. Check error message
2. Manually apply the down migration SQL
3. Update goose_db_version table

### Out of Sync

If migrations are out of sync between environments:

```bash
# Check current version
make migrate-status

# Force to specific version (dangerous!)
goose -dir migrations postgres $DB_DSN fix
```

## CI/CD Integration

In your CI/CD pipeline:

```yaml
# Example GitHub Actions
- name: Run Migrations
  run: |
    make migrate-up
  env:
    DB_HOST: ${{ secrets.DB_HOST }}
    DB_USER: ${{ secrets.DB_USER }}
    DB_PASSWORD: ${{ secrets.DB_PASSWORD }}
```

## Migration Checklist

Before deploying migrations to production:

- [ ] Test migration up in development
- [ ] Test migration down (rollback)
- [ ] Test migration up again
- [ ] Review changes with team
- [ ] Test on staging with production-like data
- [ ] Backup production database
- [ ] Have rollback plan ready
- [ ] Monitor application after migration

## Resources

- [Goose Documentation](https://github.com/pressly/goose)
- [SQL Style Guide](https://www.sqlstyle.guide/)
- [PostgreSQL Documentation](https://www.postgresql.org/docs/)

## Summary

Goose provides:

✅ Version control for database schema  
✅ Repeatable deployments  
✅ Easy rollbacks  
✅ Migration history tracking  
✅ Team collaboration on schema changes  

Use it wisely and your database schema will evolve smoothly! 🚀
