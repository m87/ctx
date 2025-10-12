package core

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateContextsExists(t *testing.T) {
	session := CreateTestSession()
	session.Free()
	assert.NoError(t, session.ValidateContextsExist(TEST_ID, TEST_ID_2))
	assert.Error(t, session.ValidateContextsExist(TEST_ID, TEST_ID_2, "test333"))

}

func TestValidateActiveInterval(t *testing.T) {
	session := CreateTestSession()
	session.CreateIfNotExistsAndSwitch("new-ctx", "test")

	for k, _ := range session.State.Contexts["new-ctx"].Intervals {
		err := session.ValidateActiveInterval("new-ctx", k)
		assert.Error(t, err)
		assert.ErrorContains(t, err, "interval is active")
	}

	session.Free()

	for k, _ := range session.State.Contexts["new-ctx"].Intervals {
		assert.NoError(t, session.ValidateActiveInterval("new-ctx", k))
	}
}
