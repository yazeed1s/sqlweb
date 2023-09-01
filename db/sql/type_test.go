package sql

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDbTypeString(t *testing.T) {
	t.Run("MySQL", func(t *testing.T) {
		dbType := MySQL
		expected := "MySQL"
		assert.Equal(t, expected, dbType.String())
	})

	t.Run("Postgres", func(t *testing.T) {
		dbType := PostgreSQL
		expected := "PostgreSQL"
		assert.Equal(t, expected, dbType.String())
	})

	t.Run("SQLite", func(t *testing.T) {
		dbType := SQLite
		expected := "SQLite"
		assert.Equal(t, expected, dbType.String())
	})

	t.Run("Unsupported", func(t *testing.T) {
		dbType := Unsupported
		expected := "Unsupported"
		assert.Equal(t, expected, dbType.String())
	})
}

func TestDbTypeEnumIndex(t *testing.T) {
	t.Run("MySQL", func(t *testing.T) {
		dbType := MySQL
		expected := 1
		assert.Equal(t, expected, dbType.EnumIndex())
	})

	t.Run("Postgres", func(t *testing.T) {
		dbType := PostgreSQL
		expected := 2
		assert.Equal(t, expected, dbType.EnumIndex())
	})

	t.Run("SQLite", func(t *testing.T) {
		dbType := SQLite
		expected := 3
		assert.Equal(t, expected, dbType.EnumIndex())
	})

	t.Run("Unsupported", func(t *testing.T) {
		dbType := Unsupported
		expected := 4
		assert.Equal(t, expected, dbType.EnumIndex())
	})
}
