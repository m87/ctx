package core

import (
	"testing"

	"github.com/stretchr/testify/assert"
)



func TestLabelContext(t *testing.T) {
	session := CreateTestSession()

	assert.Len(t, session.State.Contexts[TEST_ID].Labels, 0)

	err := session.LabelContext(TEST_ID, "label")

	assert.NoError(t, err)
	assert.Len(t, session.State.Contexts[TEST_ID].Labels, 1)
	assert.Equal(t, session.State.Contexts[TEST_ID].Labels[0], "label")

}


func TestDeleteLabelContext(t *testing.T) {
	session := CreateTestSession()

	err := session.LabelContext(TEST_ID, "label")
	assert.NoError(t, err)
	assert.Len(t, session.State.Contexts[TEST_ID].Labels, 1)

	err = session.DeleteLabelContext(TEST_ID, "label")
	assert.NoError(t, err)
	assert.Len(t, session.State.Contexts[TEST_ID].Labels, 0)

}


func TestLabelContextUnique(t *testing.T) {
	session := CreateTestSession()

	assert.Len(t, session.State.Contexts[TEST_ID].Labels, 0)

	err := session.LabelContext(TEST_ID, "label")
	err = session.LabelContext(TEST_ID, "label")
	err = session.LabelContext(TEST_ID, "label2")

	assert.NoError(t, err)
	assert.Len(t, session.State.Contexts[TEST_ID].Labels, 2)
	assert.Equal(t, session.State.Contexts[TEST_ID].Labels[0], "label")
	assert.Equal(t, session.State.Contexts[TEST_ID].Labels[1], "label2")


}


