import {
  DASHBOARD_COLUMN_COUNT,
  DASHBOARD_MINIMUM_ROW_COUNT,
  DASHBOARD_TYPE_OPTIONS,
  dashboardTargetId,
  dashboardTypeRequiresTarget,
  dashboardsAsOptions,
  dashboardGridRowCount,
  queryPreviewSummary,
  savedQueriesAsOptions,
  selectedSavedQueryText,
} from './editor';
import { contextQueryResultAsListItems } from '../query/query-summary.component';

describe('editor helpers', () => {
  it('uses a 16 by 32 minimum dashboard grid and expands for widgets', () => {
    expect(DASHBOARD_COLUMN_COUNT).toBe(16);
    expect(DASHBOARD_MINIMUM_ROW_COUNT).toBe(32);
    expect(dashboardGridRowCount([])).toBe(32);
    expect(
      dashboardGridRowCount([
        { y: 4, height: 8 },
        { y: 30, height: 7 },
      ]),
    ).toBe(37);
  });

  it('defines custom and insight dashboard types with the required targets', () => {
    expect(DASHBOARD_TYPE_OPTIONS.map((option) => option.value)).toEqual([
      'custom',
      'workspace',
      'project',
      'context',
      'daily',
    ]);
    expect(dashboardTypeRequiresTarget('project')).toBe(true);
    expect(dashboardTypeRequiresTarget('context')).toBe(true);
    expect(dashboardTypeRequiresTarget('daily')).toBe(false);
    expect(dashboardTypeRequiresTarget('custom')).toBe(false);
    expect(dashboardTargetId('workspace', 'workspace-1', '')).toBe('workspace-1');
    expect(dashboardTargetId('daily', 'workspace-1', '')).toBe('workspace-1');
    expect(dashboardTargetId('project', 'workspace-1', 'project-1')).toBe('project-1');
    expect(dashboardTargetId('custom', 'workspace-1', '')).toBeUndefined();
  });

  it('maps saved dashboards to the dashboard selector', () => {
    expect(
      dashboardsAsOptions([
        {
          id: 'dashboard-1',
          workspaceId: 'workspace-1',
          type: 'custom',
          name: 'Focus',
          definition: {},
        },
        {
          id: 'dashboard-2',
          workspaceId: 'workspace-1',
          type: 'project',
          targetId: 'project-1',
          name: 'Project overview',
          definition: {},
        },
      ]),
    ).toEqual([
      {
        value: 'dashboard-1',
        label: 'Focus',
        description: 'Custom',
        keywords: ['custom'],
      },
      {
        value: 'dashboard-2',
        label: 'Project overview',
        description: 'Project insight',
        keywords: ['project'],
      },
    ]);
  });

  it('maps saved queries to side-panel options and resolves the selected query text', () => {
    const queries = [
      {
        id: 'query-1',
        workspaceId: 'workspace-1',
        name: 'Long meetings',
        query: 'project ~ "Meetings" and duration > 1h',
      },
    ];

    expect(savedQueriesAsOptions(queries)).toEqual([
      {
        value: 'query-1',
        label: 'Long meetings',
        description: 'project ~ "Meetings" and duration > 1h',
        keywords: ['project ~ "Meetings" and duration > 1h'],
      },
    ]);
    expect(selectedSavedQueryText(queries, 'query-1')).toBe(
      'project ~ "Meetings" and duration > 1h',
    );
    expect(selectedSavedQueryText(queries, 'missing')).toBe('');
  });

  it('summarizes query preview results', () => {
    expect(
      queryPreviewSummary({
        workspaceId: 'workspace-1',
        query: 'ignored by passthrough interpreter',
        contexts: [
          {
            id: 'context-1',
            workspaceId: 'workspace-1',
            name: 'Focus',
            status: 'stopped',
          },
        ],
        contextStats: [],
        totalDuration: 90 * 60 * 1_000_000_000,
        totalSessions: 3,
      }),
    ).toBe('1 context · 1h 30m · 3 sessions');
    expect(queryPreviewSummary(undefined)).toBe('');
  });

  it('maps query results to the context list used by daily and workspace views', () => {
    const items = contextQueryResultAsListItems({
      workspaceId: 'workspace-1',
      query: 'ignored by passthrough interpreter',
      contexts: [
        {
          id: 'context-1',
          workspaceId: 'workspace-1',
          name: 'Focus',
          project: { id: 'project-1', name: 'Ctx' },
        },
      ],
      contextStats: [
        {
          contextId: 'context-1',
          duration: 90 * 60 * 1_000_000_000,
          percentage: 75,
          intervalCount: 3,
        },
      ],
      totalDuration: 90 * 60 * 1_000_000_000,
      totalSessions: 3,
    });

    expect(items).toEqual([
      expect.objectContaining({
        id: 'context-1',
        name: 'Focus',
        duration: '1h 30m',
        durationValue: 90 * 60 * 1_000_000_000,
        percentage: 75,
        sessions: 3,
        project: { id: 'project-1', name: 'Ctx' },
      }),
    ]);
  });
});
