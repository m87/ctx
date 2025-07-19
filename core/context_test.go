package core

import (
	"testing"

	"github.com/stretchr/testify/assert"
)




func TestRenameContext(t *testing.T) {
	session := CreateTestSession()
	assert.Len(t, session.State.Contexts, 1)

	err := session.RenameContext(TEST_ID, "newId", "new")
	assert.NoError(t, err)
	assert.Len(t, session.State.Contexts, 1)
	assert.Equal(t, "new", session.State.Contexts["newId"].Description)

}

