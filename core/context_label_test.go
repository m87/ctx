package core

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLabelContext(t *testing.T) {
	session := CreateTestSession()

	assert.Len(t, session.State.Contexts[TEST_ID].Labels, 2)

	err := session.LabelContext(TEST_ID, "label")

	assert.NoError(t, err)
	assert.Len(t, session.State.Contexts[TEST_ID].Labels, 3)
	assert.Contains(t, session.State.Contexts[TEST_ID].Labels, "label")

}

func TestDeleteLabelContext(t *testing.T) {
	session := CreateTestSession()

	err := session.LabelContext(TEST_ID, "label")
	assert.NoError(t, err)
	assert.Len(t, session.State.Contexts[TEST_ID].Labels, 3)

	err = session.DeleteLabelContext(TEST_ID, "label")
	assert.NoError(t, err)
	assert.Len(t, session.State.Contexts[TEST_ID].Labels, 2)

}

func TestLabelContextUnique(t *testing.T) {
	session := CreateTestSession()

	assert.Len(t, session.State.Contexts[TEST_ID].Labels, 2)

	err := session.LabelContext(TEST_ID, "labelNew")
	err = session.LabelContext(TEST_ID, "labelNew")
	err = session.LabelContext(TEST_ID, "labelNew2")

	assert.NoError(t, err)
	assert.Len(t, session.State.Contexts[TEST_ID].Labels, 4)
	assert.Contains(t, session.State.Contexts[TEST_ID].Labels, "labelNew")
	assert.Contains(t, session.State.Contexts[TEST_ID].Labels, "labelNew2")

}
