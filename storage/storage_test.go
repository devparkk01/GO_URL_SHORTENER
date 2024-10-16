package storage

import (
	"os"
	"testing"

	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/require"
)

func TestNewURLStore(t *testing.T) {

	t.Run("DB_PATH is not set", func(t *testing.T) {
		_ = os.Unsetenv("DB_PATH")
		urlStore, err := NewURLStore()
		require.NotNil(t, err)
		require.Nil(t, urlStore)
	})

	t.Run("DB_PATH set to in memory", func(t *testing.T) {
		dbPath := ":memory:"
		_ = os.Setenv("DB_PATH", dbPath)
		urlStore, err := NewURLStore()
		require.Nil(t, err)
		require.NotNil(t, urlStore)
		require.NotNil(t, urlStore.db)

	})

	t.Run("Custom DB_PATH", func(t *testing.T) {
		dbPath := "database.sqlite3"
		_ = os.Setenv("DB_PATH", dbPath)
		urlStore, err := NewURLStore()
		require.Nil(t, err)
		require.NotNil(t, urlStore)
		require.NotNil(t, urlStore.db)

		defer func() {
			_ = os.Remove(dbPath)
		}()
	})

}
