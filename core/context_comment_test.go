package core

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateComment(t *testing.T) {
	session := CreateTestSession()
	comment := Comment{
		Id:      "test-comment",
		Content: "This is a test comment",
	}
	err := session.SaveContextComment("test-context", comment)
	if err != nil {
		t.Fatal(err)
	}
	assert.Contains(t, session.State.Contexts["test-context"].Comments, "test-comment")
}

func TestDeleteComment(t *testing.T) {
	session := CreateTestSession()
	comment := Comment{
		Id:      "test-comment",
		Content: "This is a test comment",
	}
	err := session.SaveContextComment("test-context", comment)
	if err != nil {
		t.Fatal(err)
	}
	assert.Contains(t, session.State.Contexts["test-context"].Comments, "test-comment")

	err = session.DeleteContextComment("test-context", "test-comment")
	if err != nil {
		t.Fatal(err)
	}
	assert.NotContains(t, session.State.Contexts["test-context"].Comments, "test-comment")
}
