package migrations

import (
	"context"
	"embed"
	"fmt"
	"log"
	"log/slog"
	"strings"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/jackc/pgx/v5/pgxpool"
)

//go:embed *.sql
var fs embed.FS

func RunDBMigration(ctx context.Context, dbInstance *pgxpool.Pool, logger *slog.Logger) error {
	d, err := iofs.New(fs, ".")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(dbInstance.Config().ConnString())
	migration, err := migrate.NewWithSourceInstance("iofs", d, dbInstance.Config().ConnString())
	if err != nil {
		return fmt.Errorf("cannot create new migrate instance: %v", err)
	}

	if err = migration.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("failed to run migrate up: %v", err)
	}

	// Start a new transaction.
	tx, err := dbInstance.Begin(ctx)
	if err != nil {
		log.Fatal(fmt.Errorf("unable to begin transaction: %w", err))
	}
	defer tx.Rollback(ctx)

	rows, err := tx.Query(ctx, "SELECT table_name FROM information_schema.tables WHERE table_schema = 'public'")
	if err != nil {
		log.Fatal(fmt.Errorf("unable to query db: %w", err))
	}

	defer rows.Close()

	var tableNames []string

	for rows.Next() {
		var tableName string
		if err := rows.Scan(&tableName); err != nil {
			log.Fatal(fmt.Errorf("unable to scan row: %w", err))
		}
		fmt.Println("Table Name:", tableName)

		tableNames = append(tableNames, tableName)
	}

	if err := rows.Err(); err != nil {
		log.Fatal(fmt.Errorf("error while iterating over rows: %w", err))
	}

	for _, t := range tableNames {
		if t == "schema_migrations" || strings.Contains(t, "_logs") {
			continue
		}

		tbl := t
		if !strings.HasSuffix(t, "ies") {
			tbl, _ = strings.CutSuffix(t, "s")
		}

		// log table
		query := `
				CREATE TABLE IF NOT EXISTS ` + t + `_logs (
					id BIGSERIAL PRIMARY KEY NOT NULL,
					` + tbl + `_id UUID REFERENCES ` + t + `(id),
					date TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
					action TEXT,
					data_before JSONB NULL,
					data_after JSONB NULL
				);`

		_, err := tx.Exec(ctx, query)
		if err != nil {
			log.Fatal(fmt.Errorf("table: %s. unable to execute SQL statement: %w", t, err))
		}

		// trigger function
		triggerFunctionQuery := fmt.Sprintf(`
		 CREATE OR REPLACE FUNCTION %s_trigger_function() RETURNS TRIGGER AS $$
		 DECLARE
			 log_table TEXT := '%s_logs';
			 log_action VARCHAR(10);
			 log_data_before JSONB;
			 log_data_after JSONB;
		 BEGIN
			 IF TG_OP = 'INSERT' THEN
				 log_action := 'insert';
				 log_data_before := '{}'::JSONB;
				 log_data_after := to_jsonb(NEW) - 'log_id';
			 ELSIF TG_OP = 'UPDATE' THEN
				 log_action := 'update';
				 log_data_before := to_jsonb(OLD) - 'log_id';
				 log_data_after := to_jsonb(NEW) - 'log_id';
			 END IF;
			 
			 INSERT INTO %s_logs (%s_id, date, action, data_before, data_after)
			 VALUES (NEW.id, current_timestamp, log_action, log_data_before, log_data_after);
 
			 RETURN NEW;
		 END;
		 $$ LANGUAGE plpgsql;`, t, t, t, tbl)

		_, err = tx.Exec(ctx, triggerFunctionQuery)
		if err != nil {
			log.Fatal(fmt.Errorf("table: %s. unable to execute trigger function SQL statement: %w", t, err))
		}

		insertTriggerQuery := fmt.Sprintf(`
	 DO $$ BEGIN
		 IF NOT EXISTS (
			 SELECT 1
			 FROM   pg_trigger
			 WHERE  tgrelid = '%s'::regclass
			 AND    tgname = '%s_insert_trigger'
		 ) THEN
			 CREATE TRIGGER %s_insert_trigger
			 AFTER INSERT ON %s
			 FOR EACH ROW
			 EXECUTE FUNCTION %s_trigger_function();
		 END IF;
	 END $$;`, t, t, t, t, t)

		_, err = tx.Exec(ctx, insertTriggerQuery)
		if err != nil {
			return fmt.Errorf("table: %s. unable to execute insert trigger SQL statement: %w", t, err)
		}

		updateTriggerQuery := fmt.Sprintf(`
	 DO $$ BEGIN
		 IF NOT EXISTS (
			 SELECT 1
			 FROM   pg_trigger
			 WHERE  tgrelid = '%s'::regclass
			 AND    tgname = '%s_update_trigger'
		 ) THEN
			 CREATE TRIGGER %s_update_trigger
			 AFTER UPDATE ON %s
			 FOR EACH ROW
			 EXECUTE FUNCTION %s_trigger_function();
		 END IF;
	 END $$;`, t, t, t, t, t)

		_, err = tx.Exec(ctx, updateTriggerQuery)
		if err != nil {
			return fmt.Errorf("table: %s. unable to execute update trigger SQL statement: %w", t, err)
		}

	}

	// Commit the transaction if everything is successful.
	err = tx.Commit(ctx)
	if err != nil {
		log.Fatal(fmt.Errorf("unable to commit transaction: %w", err))
	}

	logger.Info("db migrated successfully")
	return nil
}
