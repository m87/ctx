import { Action, Selector, State, StateContext } from '@ngxs/store';

export type WorkspaceStateModel = {
  selectedWorkspaceId: string | null;
};

export class SelectWorkspace {
  static readonly type = '[Workspace] Select';

  constructor(public workspaceId: string | null) {}
}

@State<WorkspaceStateModel>({
  name: 'workspace',
  defaults: {
    selectedWorkspaceId: null,
  },
})
export class WorkspaceState {
  @Selector()
  static selectedWorkspaceId(state: WorkspaceStateModel): string | null {
    return state.selectedWorkspaceId;
  }

  @Action(SelectWorkspace)
  selectWorkspace(ctx: StateContext<WorkspaceStateModel>, action: SelectWorkspace): void {
    const state = ctx.getState();
    ctx.setState({
      ...state,
      selectedWorkspaceId: action.workspaceId,
    });
  }
}
