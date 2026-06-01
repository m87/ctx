import { Action, Selector, State, StateContext } from '@ngxs/store';

export type WorkspaceItem = {
  id: string;
  name: string;
  contextsCount: number;
  updatedLabel: string;
};

export type WorkspaceStateModel = {
  selectedWorkspaceId: string;
  workspaces: WorkspaceItem[];
};

const DEFAULT_WORKSPACES: WorkspaceItem[] = [
  { id: 'ws-1', name: 'Personal', contextsCount: 12, updatedLabel: 'active today' },
  { id: 'ws-2', name: 'Product Team', contextsCount: 7, updatedLabel: 'updated yesterday' },
  { id: 'ws-3', name: 'Research Vault', contextsCount: 19, updatedLabel: 'updated 2 days ago' },
];


export class SelectWorkspace {
  static readonly type = '[Workspace] Select';

  constructor(public workspaceId: string) {}
}

@State<WorkspaceStateModel>({
  name: 'workspace',
  defaults: {
    selectedWorkspaceId: DEFAULT_WORKSPACES[0].id,
    workspaces: DEFAULT_WORKSPACES,
  },
})
export class WorkspaceState {
  @Selector()
  static selectedWorkspaceId(state: WorkspaceStateModel): string {
    return state.selectedWorkspaceId;
  }

  @Selector()
  static workspaces(state: WorkspaceStateModel): WorkspaceItem[] {
    return state.workspaces;
  }

  @Action(SelectWorkspace)
  selectWorkspace(ctx: StateContext<WorkspaceStateModel>, action: SelectWorkspace): void {
    const state = ctx.getState();
    if (!state.workspaces.some((workspace) => workspace.id === action.workspaceId)) {
      return;
    }

    const nextState = {
      ...state,
      selectedWorkspaceId: action.workspaceId,
    };

    ctx.setState(nextState);
  }
}
