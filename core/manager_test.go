package core

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

type mockContextRepository struct {
	contexts             []*Context
	listByWorkspaceErr   error
	listByWorkspaceCalls int
	listedWorkspaceID    string
}

func (r *mockContextRepository) GetById(string) (*Context, error) {
	return nil, nil
}

func (r *mockContextRepository) Save(*Context) (string, error) {
	return "", nil
}

func (r *mockContextRepository) Delete(string) error {
	return nil
}

func (r *mockContextRepository) List() ([]*Context, error) {
	return nil, nil
}

func (r *mockContextRepository) ListByWorkspace(workspaceID string) ([]*Context, error) {
	r.listByWorkspaceCalls++
	r.listedWorkspaceID = workspaceID
	return r.contexts, r.listByWorkspaceErr
}

func (r *mockContextRepository) GetActive() (*Context, error) {
	return nil, nil
}

type mockWorkspaceRepository struct {
	deleteErr          error
	deleteCalls        int
	deletedWorkspaceID string
}

func (r *mockWorkspaceRepository) GetById(string) (*Workspace, error) {
	return nil, nil
}

func (r *mockWorkspaceRepository) Save(*Workspace) (string, error) {
	return "", nil
}

func (r *mockWorkspaceRepository) Delete(workspaceID string) error {
	r.deleteCalls++
	r.deletedWorkspaceID = workspaceID
	return r.deleteErr
}

func (r *mockWorkspaceRepository) List() ([]*Workspace, error) {
	return nil, nil
}

func TestContextManagerDeleteWorkspaceDeletesUnusedWorkspace(t *testing.T) {
	contextRepo := &mockContextRepository{}
	workspaceRepo := &mockWorkspaceRepository{}
	manager := NewContextManager(nil, contextRepo, nil, workspaceRepo)

	err := manager.DeleteWorkspace("workspace-1")

	require.NoError(t, err)
	require.Equal(t, 1, contextRepo.listByWorkspaceCalls)
	require.Equal(t, "workspace-1", contextRepo.listedWorkspaceID)
	require.Equal(t, 1, workspaceRepo.deleteCalls)
	require.Equal(t, "workspace-1", workspaceRepo.deletedWorkspaceID)
}

func TestContextManagerDeleteWorkspaceReturnsWorkspaceInUseError(t *testing.T) {
	contextRepo := &mockContextRepository{
		contexts: []*Context{{Id: "context-1", WorkspaceId: "workspace-1"}},
	}
	workspaceRepo := &mockWorkspaceRepository{}
	manager := NewContextManager(nil, contextRepo, nil, workspaceRepo)

	err := manager.DeleteWorkspace("workspace-1")

	var workspaceInUseErr *WorkspaceInUseError
	require.ErrorAs(t, err, &workspaceInUseErr)
	require.Equal(t, "workspace-1", workspaceInUseErr.WorkspaceId)
	require.Equal(t, 0, workspaceRepo.deleteCalls)
}

func TestContextManagerDeleteWorkspaceReturnsContextRepositoryError(t *testing.T) {
	wantErr := errors.New("list contexts failed")
	contextRepo := &mockContextRepository{listByWorkspaceErr: wantErr}
	workspaceRepo := &mockWorkspaceRepository{}
	manager := NewContextManager(nil, contextRepo, nil, workspaceRepo)

	err := manager.DeleteWorkspace("workspace-1")

	require.ErrorIs(t, err, wantErr)
	require.Equal(t, 0, workspaceRepo.deleteCalls)
}

func TestContextManagerDeleteWorkspaceReturnsWorkspaceRepositoryError(t *testing.T) {
	wantErr := errors.New("delete workspace failed")
	contextRepo := &mockContextRepository{}
	workspaceRepo := &mockWorkspaceRepository{deleteErr: wantErr}
	manager := NewContextManager(nil, contextRepo, nil, workspaceRepo)

	err := manager.DeleteWorkspace("workspace-1")

	require.ErrorIs(t, err, wantErr)
	require.Equal(t, 1, workspaceRepo.deleteCalls)
	require.Equal(t, "workspace-1", workspaceRepo.deletedWorkspaceID)
}
