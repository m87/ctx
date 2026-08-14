import { Action, NgxsAfterBootstrap, Selector, State, StateContext } from '@ngxs/store';

export type WorkspaceStateModel = {
  selectedWorkspaceId: string | null;
  initialized: boolean;
  selectedProjectId?: string | null;
};

export class SelectWorkspace {
  static readonly type = '[Workspace] Select';

  constructor(public workspaceId: string | null) {}
}

export class SelectProject {
  static readonly type = '[Workspace] Select Project';

  constructor(public projectId: string | null) {}
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

  @Selector()
  static selectedProjectId(state: WorkspaceStateModel): string | null | undefined {
    return state.selectedProjectId;
  }

  @Action(SelectProject)
  selectProject(ctx: StateContext<WorkspaceStateModel>, action: SelectProject): void {
    const state = ctx.getState();
    ctx.setState({
      ...state,
      selectedProjectId: action.projectId,
    });
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
      selectedProjectId: null,
    });
  }
}
