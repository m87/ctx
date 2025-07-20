package tui

import (
	"testing"

	"github.com/m87/ctx/core"
	"github.com/stretchr/testify/assert"
)

func TestList(t *testing.T) {
	session := core.CreateTestSession()
	output := List(*session)
	expected := "- Test Context\n- Test2 Context\n"
	assert.NotEmpty(t, output)
	assert.Equal(t, expected, output)
}

func TestListFull(t *testing.T) {
	session := core.CreateTestSession()
	output := ListFull(*session)

	expected := "- [test-context] Test Context\n\t[test-interval] 2025-02-02 12:12:12 - 2025-02-02 13:12:12\n\t[test-interval-2] 2025-02-02 13:12:12 - 2025-02-02 15:12:12\n- [test-context-2] Test2 Context\n\t[test-interval-2-1] 2025-02-02 12:12:12 - 2025-02-02 13:12:12\n\t[test-interval-2-2] 2025-02-02 13:12:12 - 2025-02-02 15:12:12\n"
	assert.NotEmpty(t, output)
	assert.Equal(t, expected, output)
}

func TestListJson(t *testing.T) {
	session := core.CreateTestSession()
	output := ListJson(*session)
	expected := "[{\"id\":\"test-context\",\"description\":\"Test Context\",\"comments\":null,\"state\":0,\"duration\":7200000000000,\"intervals\":{\"test-interval-1\":{\"id\":\"test-interval\",\"start\":{\"time\":\"2025-02-02T12:12:12Z\",\"timezone\":\"UTC\"},\"end\":{\"time\":\"2025-02-02T13:12:12Z\",\"timezone\":\"UTC\"},\"duration\":3600000000000,\"labels\":null},\"test-interval-2\":{\"id\":\"test-interval-2\",\"start\":{\"time\":\"2025-02-02T13:12:12Z\",\"timezone\":\"UTC\"},\"end\":{\"time\":\"2025-02-02T15:12:12Z\",\"timezone\":\"UTC\"},\"duration\":3600000000000,\"labels\":null}},\"labels\":[\"test1-2\",\"test1-1\"]},{\"id\":\"test-context-2\",\"description\":\"Test2 Context\",\"comments\":null,\"state\":0,\"duration\":7200000000000,\"intervals\":{\"test-interval-2-1\":{\"id\":\"test-interval-2-1\",\"start\":{\"time\":\"2025-02-02T12:12:12Z\",\"timezone\":\"UTC\"},\"end\":{\"time\":\"2025-02-02T13:12:12Z\",\"timezone\":\"UTC\"},\"duration\":3600000000000,\"labels\":null},\"test-interval-2-2\":{\"id\":\"test-interval-2-2\",\"start\":{\"time\":\"2025-02-02T13:12:12Z\",\"timezone\":\"UTC\"},\"end\":{\"time\":\"2025-02-02T15:12:12Z\",\"timezone\":\"UTC\"},\"duration\":3600000000000,\"labels\":null}},\"labels\":[\"test2-2\",\"test2-1\"]}]"
	assert.NotEmpty(t, output)
	assert.Equal(t, expected, output)
}
