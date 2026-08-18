import { summarizeContextsByProject, UNASSIGNED_PROJECT_ID } from './project-time-summary';

describe('project time summary', () => {
  it('aggregates tracked time by project and sorts projects by duration', () => {
    const summaries = summarizeContextsByProject({
      contexts: [
        {
          id: 'context-1',
          name: 'Planning',
          workspaceId: 'workspace-1',
          project: { id: 'project-1', name: 'Website' },
        },
        {
          id: 'context-2',
          name: 'Implementation',
          workspaceId: 'workspace-1',
          project: { id: 'project-1', name: 'Website' },
        },
        {
          id: 'context-3',
          name: 'Research',
          workspaceId: 'workspace-1',
          project: { id: 'project-2', name: 'Mobile app' },
        },
      ],
      contextStats: [
        { contextId: 'context-1', duration: 30 },
        { contextId: 'context-2', duration: 60 },
        { contextId: 'context-3', duration: 45 },
      ],
    });

    expect(summaries).toEqual([
      {
        id: 'project-1',
        name: 'Website',
        duration: 90,
        percentage: (90 / 135) * 100,
        contextCount: 2,
        project: { id: 'project-1', name: 'Website' },
      },
      {
        id: 'project-2',
        name: 'Mobile app',
        duration: 45,
        percentage: (45 / 135) * 100,
        contextCount: 1,
        project: { id: 'project-2', name: 'Mobile app' },
      },
    ]);
  });

  it('includes tracked contexts without a project in a separate group', () => {
    const summaries = summarizeContextsByProject({
      contexts: [{ id: 'context-1', name: 'Admin', workspaceId: 'workspace-1' }],
      contextStats: [
        { contextId: 'context-1', duration: 30 },
        { contextId: 'missing-context', duration: 15 },
      ],
    });

    expect(summaries).toEqual([
      {
        id: UNASSIGNED_PROJECT_ID,
        name: 'No project',
        duration: 45,
        percentage: 100,
        contextCount: 2,
        project: undefined,
      },
    ]);
  });

  it('omits contexts without tracked time', () => {
    expect(
      summarizeContextsByProject({
        contexts: [{ id: 'context-1', name: 'Planning', workspaceId: 'workspace-1' }],
        contextStats: [{ contextId: 'context-1', duration: 0 }],
      }),
    ).toEqual([]);
  });
});
