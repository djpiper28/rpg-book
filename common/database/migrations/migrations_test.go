package migrations_test

import (
	"errors"
	"os"
	"testing"

	"github.com/djpiper28/rpg-book/common/database/migrations"
	"github.com/djpiper28/rpg-book/common/database/sqlite3"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/require"
)

func TestMigrationsEmptyList(t *testing.T) {
	name := uuid.New().String() + sqlite3.DbExtension
	defer os.Remove(name)

	db, err := sqlite3.New(name)
	defer db.Close()
	require.NoError(t, err)

	m := migrations.New([]migrations.Migration{})
	err = m.Migrate(db)
	require.NoError(t, err)
}

func TestMigrationsOneTable(t *testing.T) {
	name := uuid.New().String() + sqlite3.DbExtension
	defer os.Remove(name)

	db, err := sqlite3.New(name)
	defer db.Close()
	require.NoError(t, err)

	var preProcessDone, postProcessDone, testDone bool
	m := migrations.New([]migrations.Migration{
		{
			PreProcess: func(tx *sqlx.Tx) error {
				preProcessDone = true
				require.NotNil(t, tx)
				return nil
			},
			PostProcess: func(tx *sqlx.Tx) error {
				postProcessDone = true
				require.NotNil(t, tx)
				return nil
			},
			Sql: `
  CREATE TABLE test_table (
    id UUID PRIMARY KEY
  );
      `,
			Test: func(tx *sqlx.Tx) error {
				testDone = true
				require.NotNil(t, tx)
				_, err := tx.Exec("INSERT INTO test_table (id) values('9e48d824-13ff-4e6f-aac7-fbe42498556e');")
				require.NoError(t, err)
				return nil
			},
		},
	})
	err = m.Migrate(db)
	require.NoError(t, err)

	require.True(t, preProcessDone)
	require.True(t, postProcessDone)
	require.True(t, testDone)

	_, err = db.Db.Exec("INSERT INTO test_table (id) values('578d07ce-5ad8-49ac-92b5-5b52c127fdcc');")
	require.NoError(t, err)
}

func TestMigrationsNilFunctions(t *testing.T) {
	name := uuid.New().String() + sqlite3.DbExtension
	defer os.Remove(name)

	db, err := sqlite3.New(name)
	defer db.Close()
	require.NoError(t, err)

	m := migrations.New([]migrations.Migration{
		{
			Sql: `
    CREATE TABLE test_table (
      id UUID PRIMARY KEY
    );
      `,
		},
	})
	err = m.Migrate(db)
	require.NoError(t, err)

	_, err = db.Db.Exec(`
    INSERT INTO test_table (id) 
    VALUES ('578d07ce-5ad8-49ac-92b5-5b52c127fdcc');`)
	require.NoError(t, err)
}

func TestMigrationsTwo(t *testing.T) {
	name := uuid.New().String() + sqlite3.DbExtension
	defer os.Remove(name)

	db, err := sqlite3.New(name)
	defer db.Close()
	require.NoError(t, err)

	m := migrations.New([]migrations.Migration{
		{
			Sql: `
  CREATE TABLE test_table (
    id UUID PRIMARY KEY
  );
      `,
		},
		{
			Sql: `
  ALTER TABLE test_table ADD COLUMN name TEXT;
      `,
		},
	})
	err = m.Migrate(db)
	require.NoError(t, err)

	_, err = db.Db.Exec(`
    INSERT INTO test_table (id, name) 
    VALUES ('578d07ce-5ad8-49ac-92b5-5b52c127fdcc', 'Uwu 123');`)
	require.NoError(t, err)
}

func TestExecutedMigrationsAreNotExecutedAgain(t *testing.T) {
	name := uuid.New().String() + sqlite3.DbExtension
	defer os.Remove(name)

	db, err := sqlite3.New(name)
	require.NoError(t, err)

	m := migrations.New([]migrations.Migration{
		{
			Sql: `
    CREATE TABLE test_table (
      id UUID PRIMARY KEY
    );
      `,
		},
	})
	err = m.Migrate(db)
	require.NoError(t, err)

	_, err = db.Db.Exec(`
    INSERT INTO test_table (id) 
    VALUES ('578d07ce-5ad8-49ac-92b5-5b52c127fdcc');`)
	require.NoError(t, err)

	db.Close()

	db, err = sqlite3.New(name)
	defer db.Close()
	require.NoError(t, err)

	m = migrations.New([]migrations.Migration{
		{
			Sql: `
      INVALID SQL
      `,
			PreProcess: func(tx *sqlx.Tx) error {
				t.Log("This migration should never have been executed")
				t.Fail()
				return nil
			},
		},
	})
	err = m.Migrate(db)
	require.NoError(t, err)

	_, err = db.Db.Exec(`
    INSERT INTO test_table (id) 
    VALUES ('263e06d8-7e07-4d12-9479-d7faf0156743');`)
	require.NoError(t, err)
}

func TestMigrationsDifferentTransactions(t *testing.T) {
	name := uuid.New().String() + sqlite3.DbExtension
	defer os.Remove(name)

	db, err := sqlite3.New(name)
	require.NoError(t, err)

	m := migrations.New([]migrations.Migration{
		{
			Sql: `
    CREATE TABLE test_table (
      id UUID PRIMARY KEY
    );
      `,
		},
	})
	err = m.Migrate(db)
	require.NoError(t, err)

	_, err = db.Db.Exec(`
    INSERT INTO test_table (id) 
    VALUES ('578d07ce-5ad8-49ac-92b5-5b52c127fdcc');`)
	require.NoError(t, err)

	db.Close()

	db, err = sqlite3.New(name)
	defer db.Close()
	require.NoError(t, err)

	m = migrations.New([]migrations.Migration{
		{
			Sql: `
    CREATE TABLE test_table (
      id UUID PRIMARY KEY
    );
      `,
		},
		{
			Sql: `
    ALTER TABLE test_table ADD COLUMN foo TEXT DEFAULT('baz');
      `,
		},
	})
	err = m.Migrate(db)
	require.NoError(t, err)

	_, err = db.Db.Exec(`
    INSERT INTO test_table (id) 
    VALUES ('263e06d8-7e07-4d12-9479-d7faf0156743');`)
	require.NoError(t, err)
}

func TestPreProcessError(t *testing.T) {
	name := uuid.New().String() + sqlite3.DbExtension
	defer os.Remove(name)

	db, err := sqlite3.New(name)
	defer db.Close()
	require.NoError(t, err)

	m := migrations.New([]migrations.Migration{
		{
			PreProcess: func(tx *sqlx.Tx) error {
				return errors.New("Mocked error")
			},
			Sql: `
  SHOULD NOT BE EXECUTED;
      `,
		},
	})
	err = m.Migrate(db)
	require.Error(t, err)
}

func TestSqlError(t *testing.T) {
	name := uuid.New().String() + sqlite3.DbExtension
	defer os.Remove(name)

	db, err := sqlite3.New(name)
	defer db.Close()
	require.NoError(t, err)

	m := migrations.New([]migrations.Migration{
		{
			Sql: `
  SHOULD CAUSE AN ERROR;
      `,
		},
	})
	err = m.Migrate(db)
	require.Error(t, err)
}

func TestPostProcessError(t *testing.T) {
	name := uuid.New().String() + sqlite3.DbExtension
	defer os.Remove(name)

	db, err := sqlite3.New(name)
	defer db.Close()
	require.NoError(t, err)

	m := migrations.New([]migrations.Migration{
		{
			Sql: `
  CREATE TABLE test_table (
    id UUID PRIMARY KEY
  );
      `,
			PostProcess: func(tx *sqlx.Tx) error {
				return errors.New("Mocked error")
			},
		},
	})
	err = m.Migrate(db)
	require.Error(t, err)
}

func TestSqlTestError(t *testing.T) {
	name := uuid.New().String() + sqlite3.DbExtension
	defer os.Remove(name)

	db, err := sqlite3.New(name)
	defer db.Close()
	require.NoError(t, err)

	m := migrations.New([]migrations.Migration{
		{
			Sql: `
  CREATE TABLE test_table (
    id UUID PRIMARY KEY
  );
      `,
			Test: func(tx *sqlx.Tx) error {
				return errors.New("Mocked error")
			},
		},
	})
	err = m.Migrate(db)
	require.Error(t, err)
}
