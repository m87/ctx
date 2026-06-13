import { Action, NgxsAfterBootstrap, Selector, State, StateContext } from '@ngxs/store';

export type WorkspaceStateModel = {
  selectedWorkspaceId: string | null;
  initialized: boolean;
};

export class SelectWorkspace {
  static readonly type = '[Workspace] Select';

  constructor(public workspaceId: string | null) {}
}

@State<WorkspaceStateModel>({
  name: 'workspace',
  defaults: {
    selectedWorkspaceId: null,
    initialized: false,
  },
})
export class WorkspaceState implements NgxsAfterBootstrap {
  @Selector()
  static selectedWorkspaceId(state: WorkspaceStateModel): string | null {
    return state.selectedWorkspaceId;
  }

  @Selector()
  static initialized(state: WorkspaceStateModel): boolean {
    return state.initialized;
  }

  ngxsAfterBootstrap(ctx: StateContext<WorkspaceStateModel>): void {
    ctx.patchState({ initialized: true });
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
