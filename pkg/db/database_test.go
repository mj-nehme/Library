package db

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestInitDB tests the database connection initialization
func TestDB(t *testing.T) {
	db := SetupTest(t)
	err := db.Teardown()
	assert.NoError(t, err)
}
