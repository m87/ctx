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

func TestDeleteContext(t *testing.T) {
	session := CreateTestSession()
	assert.Len(t, session.State.Contexts, 1)

	err := session.Delete(TEST_ID)
	assert.NoError(t, err)
	assert.Len(t, session.State.Contexts, 0)
}

func TestDeleteContextNotFound(t *testing.T) {
	session := CreateTestSession()
	assert.Len(t, session.State.Contexts, 1)

	err := session.Delete("not-found")
	assert.Error(t, err)
	assert.Len(t, session.State.Contexts, 1)
}

func TestDeleteActiveContext(t *testing.T) {
	session := CreateTestSession()
	assert.Len(t, session.State.Contexts, 1)
	session.State.CurrentId = TEST_ID

	err := session.Delete(TEST_ID)
	assert.Error(t, err)
	assert.Len(t, session.State.Contexts, 1)
	assert.Equal(t, TEST_ID, session.State.CurrentId)
}
