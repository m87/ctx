import { computed, effect, inject, signal, untracked } from '@angular/core';
import { injectQuery } from '@tanstack/angular-query-experimental';
import { DashboardQueries } from '../../api/dashboard/dashboard.queries';
import {
  Dashboard,
  DashboardScope,
  DashboardWidgetViewState,
} from '../../api/dashboard/dashboard.service';
import { SettingsQueries } from '../../api/settings/settings.queries';
import {
  DashboardDefinition,
  normalizeDashboardDefinition,
} from '../editor/dashboard/dashboard-definition';
import { TimeZoneService } from '../shared/time-zone.service';

export function injectDashboardMetrics(
  dashboard: () => Dashboard | null,
  scope: () => DashboardScope = () => ({}),
  enabled: () => boolean = () => true,
) {
  const queries = inject(DashboardQueries);
  const timeZone = inject(TimeZoneService);
  const settings = inject(SettingsQueries);
  const settingsQuery = injectQuery(() => settings.settings());
  const viewState = signal<Record<string, DashboardWidgetViewState>>({});
  const requestConfiguration = computed(() => {
    const content = definition();
    if (!content || !dashboard()) return null;
    return {
      workspaceId: dashboard()?.workspaceId ?? '',
      definition: content,
      scope: scope(),
      timeZone: timeZone.effectiveTimeZone(),
      firstDayOfWeek: settingsQuery.data()?.['client.general.firstDay'] === 'Sunday' ? 0 : 1,
      viewState: viewState(),
    };
  });
  const requestSignature = computed(() => {
    const configuration = requestConfiguration();
    if (!configuration) return '';
    return JSON.stringify({
      workspaceId: configuration.workspaceId,
      scope: configuration.scope,
      timeZone: configuration.timeZone,
      firstDayOfWeek: configuration.firstDayOfWeek,
      data: metricDataKey(configuration.definition),
      viewState: configuration.viewState,
    });
  });
  const debouncedConfiguration = signal<ReturnType<typeof requestConfiguration> | null>(null);
  const debouncedSignature = signal('');
  effect((onCleanup) => {
    const signature = requestSignature();
    if (!signature) {
      debouncedConfiguration.set(null);
      debouncedSignature.set('');
      return;
    }
    const timeout = setTimeout(() => {
      debouncedConfiguration.set(untracked(requestConfiguration));
      debouncedSignature.set(signature);
    }, 300);
    onCleanup(() => clearTimeout(timeout));
  });
  const definition = computed(() => {
    try {
      return normalizeDashboardDefinition(dashboard()?.definition ?? {});
    } catch {
      return null;
    }
  });
  let previousDataSignature = '';
  effect(() => {
    const content = definition();
    if (!content) return;
    const signature = JSON.stringify(metricDataKey(content));
    if (previousDataSignature && signature !== previousDataSignature) {
      viewState.update((current) =>
        Object.fromEntries(
          Object.entries(current).map(([id, state]) => [id, { ...state, page: 1 }]),
        ),
      );
    }
    previousDataSignature = signature;
  });
  const query = injectQuery(() => {
    const configuration = debouncedConfiguration();
    const content = configuration?.definition;
    const widgets = content?.widgets.filter((widget) => widget.type !== 'text') ?? [];
    const dataKey = content
      ? { ...metricDataKey(content), viewState: configuration?.viewState ?? {} }
      : {};
    const liveScope = !configuration?.scope.date || configuration.scope.date === timeZone.today();
    const refresh =
      liveScope &&
      widgets.some((widget) => {
        const period =
          widget.properties.period === 'inherit'
            ? (content?.period ?? 'this_week')
            : widget.properties.period;
        const currentRange = [
          'today',
          'this_week',
          'this_month',
          'last_7_days',
          'last_30_days',
          'last_3_months',
          'last_6_months',
          'last_12_months',
        ].includes(period);
        return currentRange && (widget.properties.includeRunning || widget.type === 'heatmap');
      })
        ? 30_000
        : false;
    return queries.evaluate(
      {
        workspaceId: configuration?.workspaceId ?? '',
        definition: content ?? {},
        scope: configuration?.scope ?? {},
        timeZone: configuration?.timeZone ?? timeZone.effectiveTimeZone(),
        firstDayOfWeek: configuration?.firstDayOfWeek ?? 1,
        viewState: configuration?.viewState,
      },
      dataKey,
      !!content && widgets.length > 0 && enabled() && !!configuration,
      refresh,
    );
  });
  const currentConfigurationMatches = computed(() => requestSignature() === debouncedSignature());
  return {
    query,
    firstDayOfWeek: computed(() =>
      settingsQuery.data()?.['client.general.firstDay'] === 'Sunday' ? 0 : 1,
    ),
    results: computed(() =>
      enabled() && currentConfigurationMatches() ? (query.data()?.widgets ?? {}) : {},
    ),
    error: computed(() =>
      !enabled() || !currentConfigurationMatches()
        ? ''
        : query.isError()
          ? 'Unable to load metric data.'
          : query.isPaused()
            ? 'Waiting for a connection.'
            : '',
    ),
    setWidgetViewState(widgetId: string, state: DashboardWidgetViewState) {
      viewState.update((current) => ({ ...current, [widgetId]: state }));
    },
  };
}

function metricDataKey(content: DashboardDefinition) {
  return {
    query: content.query,
    period: content.period,
    widgets: content.widgets
      .filter((widget) => widget.type !== 'text')
      .map((widget) => ({
        id: widget.id,
        query: widget.query,
        metric: 'metric' in widget.properties ? widget.properties.metric : undefined,
        period: widget.properties.period,
        startDate: widget.properties.startDate,
        endDate: widget.properties.endDate,
        projectId: widget.properties.projectId,
        contextId: widget.properties.contextId,
        tag: widget.properties.tag,
        archived: widget.properties.archived,
        includeRunning: widget.properties.includeRunning,
        ...(widget.type === 'radar'
          ? { axesBy: widget.properties.axesBy, axisKeys: [...widget.properties.axisKeys].sort() }
          : {}),
        ...(widget.type === 'bar'
          ? {
              groupBy: widget.properties.groupBy,
              bucket: widget.properties.bucket,
              sort: widget.properties.sort,
              limit: widget.properties.limit,
              showOther: widget.properties.showOther,
            }
          : {}),
        ...(widget.type === 'line'
          ? {
              bucket: widget.properties.bucket,
              seriesBy: widget.properties.seriesBy,
              seriesLimit: widget.properties.seriesLimit,
              valueMode: widget.properties.valueMode,
            }
          : {}),
        ...(widget.type === 'table'
          ? {
              groupBy: widget.properties.groupBy,
              columns: widget.properties.columns
                .map(({ id, metric }) => ({ id, metric }))
                .sort((a, b) => a.id.localeCompare(b.id)),
              sort: widget.properties.sort,
              pageSize: widget.properties.pageSize,
            }
          : {}),
      })),
  };
}
