import {
  EDITOR_QUERY_SCOPE_OPTIONS,
  queryPreviewSummary,
  queryScopeHasEntity,
  savedQueriesAsOptions,
  selectedSavedQueryText,
} from './editor';
import { contextQueryResultAsListItems } from '../query/query-summary.component';

describe('editor query filters', () => {
  it('keeps saved queries separate from the main query scope filter', () => {
    expect(EDITOR_QUERY_SCOPE_OPTIONS.map((option) => option.value)).not.toContain('saved-query');
    expect(queryScopeHasEntity('project')).toBe(true);
    expect(queryScopeHasEntity('daily')).toBe(false);
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
