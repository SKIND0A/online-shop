package postgres

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

// EnsureSchema добавляет недостающие колонки под текущий код (без полноценного migrate).
// Нужно, если применили только старую 000001 без 000002_user_display_name.
func EnsureSchema(ctx context.Context, pool *pgxpool.Pool) error {
	_, err := pool.Exec(ctx, `
		ALTER TABLE users ADD COLUMN IF NOT EXISTS display_name VARCHAR(120) NOT NULL DEFAULT '';
	`)
	return err
}
