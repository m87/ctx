package core

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestRenameContext(t *testing.T) {
	session := CreateTestSession()
	assert.Len(t, session.State.Contexts, 2)

	err := session.RenameContext(TEST_ID, "newId", "new")
	assert.NoError(t, err)
	assert.Len(t, session.State.Contexts, 2)
	assert.Equal(t, "new", session.State.Contexts["newId"].Description)

}

func TestDeleteContext(t *testing.T) {
	session := CreateTestSession()
	assert.Len(t, session.State.Contexts, 2)
	err := session.Delete(TEST_ID)
	assert.NoError(t, err)
	assert.Len(t, session.State.Contexts, 1)
}

func TestDeleteContextNotFound(t *testing.T) {
	session := CreateTestSession()
	assert.Len(t, session.State.Contexts, 2)

	err := session.Delete("not-found")
	assert.Error(t, err)
	assert.Len(t, session.State.Contexts, 2)
}

func TestDeleteActiveContext(t *testing.T) {
	session := CreateTestSession()
	assert.Len(t, session.State.Contexts, 2)
	session.State.CurrentId = TEST_ID

	err := session.Delete(TEST_ID)
	assert.Error(t, err)
	assert.Len(t, session.State.Contexts, 2)
	assert.Equal(t, TEST_ID, session.State.CurrentId)
}

func TestMergeContext(t *testing.T) {
	session := CreateTestSession()
	assert.Len(t, session.State.Contexts, 2)

	err := session.MergeContext(TEST_ID_2, TEST_ID)
	assert.NoError(t, err)
	assert.Len(t, session.State.Contexts, 1)
	assert.Len(t, session.State.Contexts[TEST_ID].Intervals, 4)
	assert.Len(t, session.State.Contexts[TEST_ID].Labels, 4)
	assert.Equal(t, session.State.Contexts[TEST_ID].Duration, 4*time.Hour)
}

func TestMergContextActiveContext(t *testing.T) {
	session := CreateTestSession()
	session.State.CurrentId = TEST_ID_2
	err := session.MergeContext(TEST_ID_2, TEST_ID)
	assert.Error(t, err)
}

func TestMergContextNotExistContext(t *testing.T) {
	session := CreateTestSession()
	err := session.MergeContext(TEST_ID_2, "not-found")
	assert.Error(t, err)

	err = session.MergeContext("not-found", TEST_ID)
	assert.Error(t, err)
}

func TestMergeSameContext(t *testing.T) {
	session := CreateTestSession()
	err := session.MergeContext(TEST_ID, TEST_ID)
	assert.Error(t, err)
}
