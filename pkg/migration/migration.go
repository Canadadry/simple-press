package migration

import (
	"context"
	"database/sql"
	"fmt"
	"io"
	"io/fs"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

const (
	modeMySQL         = "mysql"
	mySqlMigrationExt = ".sql"
	querySeparator    = "-- separator"
)

type Sql interface {
	BeginTx(context.Context, *sql.TxOptions) (*sql.Tx, error)
	ExecContext(context.Context, string, ...any) (sql.Result, error)
	QueryContext(context.Context, string, ...any) (*sql.Rows, error)
}

func Migrate(db Sql, mode string, f fs.ReadDirFS) error {
	ctx, cancel := context.WithTimeout(context.TODO(), 5*60*time.Second)
	defer cancel()

	err := createMigrationTableIfNecessary(db, ctx)
	if err != nil {
		return fmt.Errorf("while looking for migration table : %w", err)
	}
	known, err := getKnownMigration(f)
	if err != nil {
		return fmt.Errorf("while querying list of known migration files : %w", err)
	}
	if len(known) == 0 {
		return fmt.Errorf("no migration file found")
	}
	applied, err := getAppliedMigration(db, ctx)
	if err != nil {
		return fmt.Errorf("while querying list of applied migration files : %w", err)
	}
	if !isCompatible(applied, known) {
		return fmt.Errorf("schema cannot be migrated : know migration not compatible with those already applied \n applied : %v \n known : %v", applied, known)
	}
	for i := len(applied); i < len(known); i++ {
		err = applyMigration(db, mySqlMigrationExt, ctx, f, known[i])
		if err != nil {
			return fmt.Errorf("while applying migration %d : %w", i, err)
		}
	}
	return nil
}

func getKnownMigration(f fs.ReadDirFS) ([]string, error) {
	files, err := f.ReadDir(".")
	if err != nil {
		return nil, fmt.Errorf("reading embed fs : %w", err)
	}
	known := []string{}
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		base := filepath.Base(file.Name())
		if filepath.Ext(base) != mySqlMigrationExt {
			continue
		}
		base = base[:len(base)-len(mySqlMigrationExt)]
		known = append(known, base)
	}

	sort.Strings(known)
	return known, nil
}

func isCompatible(applied, known []string) bool {
	if len(applied) > len(known) {
		return false
	}
	for i := 0; i < len(applied); i++ {
		if applied[i] != known[i] {
			return false
		}
	}
	return true
}

const (
	createMigrationTableQueryMysql = `CREATE TABLE IF NOT EXISTS migration_version (migration VARCHAR(255) );`
)

func createMigrationTableIfNecessary(db Sql, ctx context.Context) error {
	_, err := db.ExecContext(ctx, createMigrationTableQueryMysql)
	if err != nil {
		return fmt.Errorf("executing %s:%w", createMigrationTableQueryMysql, err)
	}
	return nil
}

const getAppliedMigrationQuery = `select migration from migration_version;`

func getAppliedMigration(db Sql, ctx context.Context) ([]string, error) {

	rows, err := db.QueryContext(ctx, getAppliedMigrationQuery)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []string
	for rows.Next() {
		var i string
		err = rows.Scan(&i)
		if err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	sort.Strings(items)
	return items, nil
}

const addMigrationQuery = "Insert into migration_version(migration) VALUES(?);"

func applyMigration(db Sql, ext string, ctx context.Context, f fs.FS, migrationName string) error {

	migrationFile, err := f.Open(migrationName + ext)
	if err != nil {
		return fmt.Errorf("can't open %s%s : %w", migrationName, ext, err)
	}
	migration, err := io.ReadAll(migrationFile)
	if err != nil {
		return fmt.Errorf("can't read %s%s : %w", migrationName, ext, err)
	}
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("cannot start a transaction : %w", err)
	}

	steps := strings.Split(string(migration), querySeparator)
	for i, s := range steps {
		_, err = tx.ExecContext(ctx, s)
		if err != nil {
			_ = tx.Rollback()
			help := fmt.Sprintf("if you think this is a mistake maybe you have forgotten to split your query with '%s'", querySeparator)
			return fmt.Errorf("while applying step %d migration '%s' : %w\n%s", i, s, err, help)
		}
	}
	_, err = tx.ExecContext(ctx, addMigrationQuery, migrationName)
	if err != nil {
		_ = tx.Rollback()
		return fmt.Errorf("while tagging migration '%s' : %w", migrationName, err)
	}
	fmt.Println("Migrating to ", migrationName)
	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("cannot commit transaction %w", err)
	}
	return nil
}
