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

func TestCreateContext(t *testing.T) {
	session := CreateTestSession()
	assert.Len(t, session.State.Contexts, 2)

	err := session.CreateContext("newId", "new")
	assert.NoError(t, err)
	assert.Len(t, session.State.Contexts, 3)
	assert.Equal(t, "new", session.State.Contexts["newId"].Description)
}

func TestCreateExistingContext(t *testing.T) {
	session := CreateTestSession()
	assert.Len(t, session.State.Contexts, 2)

	err := session.CreateContext(TEST_ID, "new")
	assert.Error(t, err)
	assert.Len(t, session.State.Contexts, 2)
}

func TestCreateContextWithEmptyDescription(t *testing.T) {
	session := CreateTestSession()
	assert.Len(t, session.State.Contexts, 2)

	err := session.CreateContext("newId", "")
	assert.Error(t, err)
	assert.Len(t, session.State.Contexts, 2)
}

func TestSwitchContext(t *testing.T) {
	session := CreateTestSession()
	assert.Len(t, session.State.Contexts, 2)

	err := session.Switch(TEST_ID_2)
	assert.NoError(t, err)
	assert.Equal(t, TEST_ID_2, session.State.CurrentId)

	err = session.Switch(TEST_ID)
	assert.NoError(t, err)
	assert.Equal(t, TEST_ID, session.State.CurrentId)

	for _, v := range session.State.Contexts[TEST_ID_2].Intervals {
		assert.False(t, v.End.Time.IsZero())
	}
}

func TestSwitchContextNotExists(t *testing.T) {
	session := CreateTestSession()
	assert.Len(t, session.State.Contexts, 2)

	err := session.Switch("not-found")
	assert.Error(t, err)
	assert.Equal(t, "", session.State.CurrentId)
}

func TestSwitchToActiveContext(t *testing.T) {
	session := CreateTestSession()
	assert.Len(t, session.State.Contexts, 2)

	session.State.CurrentId = TEST_ID
	err := session.Switch(TEST_ID)
	assert.Error(t, err)
	assert.Equal(t, TEST_ID, session.State.CurrentId)
}

func TestSwitchToEmptyContext(t *testing.T) {
	session := CreateTestSession()
	assert.Len(t, session.State.Contexts, 2)

	session.State.CurrentId = ""
	err := session.Switch(TEST_ID)
	assert.NoError(t, err)
	assert.Equal(t, TEST_ID, session.State.CurrentId)
}

func TestCreateIfNotExistsAndSwitch(t *testing.T) {
	session := CreateTestSession()
	assert.Len(t, session.State.Contexts, 2)

	err := session.CreateIfNotExistsAndSwitch("newId", "new")
	assert.NoError(t, err)
	assert.Len(t, session.State.Contexts, 3)
	assert.Equal(t, "new", session.State.Contexts["newId"].Description)
	assert.Equal(t, "newId", session.State.CurrentId)

	err = session.CreateIfNotExistsAndSwitch(TEST_ID, "new")
	assert.NoError(t, err)
	assert.Len(t, session.State.Contexts, 3)
	assert.Equal(t, TEST_ID, session.State.CurrentId)
}

func TestSearchContext(t *testing.T) {
	session := CreateTestSession()
	assert.Len(t, session.State.Contexts, 2)

	ctxs, err := session.Search("Test")
	assert.NoError(t, err)
	assert.Len(t, ctxs, 2)

	ctxs, err = session.Search("Test2")
	assert.NoError(t, err)
	assert.Len(t, ctxs, 1)
	assert.Equal(t, "Test2 Context", ctxs[0].Description)

	ctxs, err = session.Search("NonExistent")
	assert.NoError(t, err)
	assert.Len(t, ctxs, 0)
}
func TestSearchContextWithRegex(t *testing.T) {
	session := CreateTestSession()
	assert.Len(t, session.State.Contexts, 2)

	ctxs, err := session.Search("Test.*Context")
	assert.NoError(t, err)
	assert.Len(t, ctxs, 2)

	ctxs, err = session.Search("Test2 Context")
	assert.NoError(t, err)
	assert.Len(t, ctxs, 1)
	assert.Equal(t, "Test2 Context", ctxs[0].Description)

	ctxs, err = session.Search("NonExistent")
	assert.NoError(t, err)
	assert.Len(t, ctxs, 0)
}
