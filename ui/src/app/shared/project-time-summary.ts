import type { Context, ProjectMetadata } from '../../api/context/context.service';

export const UNASSIGNED_PROJECT_ID = '__unassigned_project__';

export interface ProjectTimeSummary {
  id: string;
  name: string;
  duration: number;
  percentage: number;
  contextCount: number;
  project?: ProjectMetadata;
}

interface ContextTimeStats {
  contextId: string;
  duration: number;
}

export function summarizeContextsByProject(source: {
  contexts: readonly Context[];
  contextStats: readonly ContextTimeStats[];
}): ProjectTimeSummary[] {
  const contextsById = new Map(source.contexts.map((context) => [context.id, context]));
  const durationsByProject = new Map<
    string,
    Omit<ProjectTimeSummary, 'duration' | 'percentage'> & { duration: number }
  >();

  for (const contextStats of source.contextStats) {
    if (!Number.isFinite(contextStats.duration) || contextStats.duration <= 0) {
      continue;
    }

    const project = contextsById.get(contextStats.contextId)?.project;
    const projectId = project?.id ?? UNASSIGNED_PROJECT_ID;
    const summary = durationsByProject.get(projectId) ?? {
      id: projectId,
      name: project?.name ?? 'No project',
      duration: 0,
      contextCount: 0,
      project,
    };

    summary.duration += contextStats.duration;
    summary.contextCount += 1;
    durationsByProject.set(projectId, summary);
  }

  const totalDuration = [...durationsByProject.values()].reduce(
    (total, summary) => total + summary.duration,
    0,
  );

  return [...durationsByProject.values()]
    .map((summary) => ({
      ...summary,
      percentage: totalDuration > 0 ? (summary.duration / totalDuration) * 100 : 0,
    }))
    .sort((left, right) => right.duration - left.duration || left.name.localeCompare(right.name));
}
