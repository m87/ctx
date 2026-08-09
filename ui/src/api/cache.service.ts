import { inject, Injectable } from '@angular/core';
import { QueryClient, QueryKey } from '@tanstack/angular-query-experimental';
import { contextQueryKeys } from './context.queries';
import { intervalQueryKeys } from './interval.queries';
import { settingsQueryKeys } from './settings.queries';
import { workspaceQueryKeys } from './workspace.queries';

@Injectable({ providedIn: 'root' })
export class CacheService {
  private readonly queryClient = inject(QueryClient);

  afterContextCreate() {
    return this.invalidate(contextQueryKeys.lists(), workspaceQueryKeys.stats());
  }

  afterProjectChange(workspaceId: string) {
    return this.invalidate(
      workspaceQueryKeys.statsFor(workspaceId),
      workspaceQueryKeys.detail(workspaceId),
    );
  }

  afterActiveIntervalChange() {
    return this.invalidate(
      contextQueryKeys.lists(),
      contextQueryKeys.active(),
      contextQueryKeys.intervals(),
      contextQueryKeys.stats(),
      intervalQueryKeys.dayStats(),
      intervalQueryKeys.days(),
      workspaceQueryKeys.stats(),
    );
  }

  afterContextMetadataChange(contextId: string, includeActive = false) {
    return this.invalidate(
      contextQueryKeys.lists(),
      contextQueryKeys.detail(contextId),
      intervalQueryKeys.dayStats(),
      intervalQueryKeys.days(),
      workspaceQueryKeys.stats(),
      ...(includeActive ? [contextQueryKeys.active()] : []),
    );
  }

  afterContextDelete(contextId: string) {
    this.remove(
      contextQueryKeys.detail(contextId),
      contextQueryKeys.intervalsFor(contextId),
      contextQueryKeys.statsFor(contextId),
    );

    return this.invalidate(
      contextQueryKeys.lists(),
      contextQueryKeys.active(),
      intervalQueryKeys.dayStats(),
      intervalQueryKeys.days(),
      workspaceQueryKeys.stats(),
    );
  }

  afterIntervalChange(...contextIds: string[]) {
    const queryKeys: QueryKey[] = [
      intervalQueryKeys.days(),
      intervalQueryKeys.dayStats(),
      workspaceQueryKeys.stats(),
    ];

    for (const contextId of new Set(contextIds.filter(Boolean))) {
      queryKeys.push(
        contextQueryKeys.intervalsFor(contextId),
        contextQueryKeys.statsFor(contextId),
        contextQueryKeys.detail(contextId),
      );
    }

    return this.invalidate(...queryKeys);
  }

  afterIntervalDelete(intervalId: string, contextId: string) {
    this.remove(intervalQueryKeys.detail(intervalId));
    return this.afterIntervalChange(contextId);
  }

  afterSettingsSave(settingKeys: string[]) {
    return this.invalidate(
      settingsQueryKeys.settings(),
      ...settingKeys.map((key) => settingsQueryKeys.setting(key)),
    );
  }

  afterWorkspaceListChange() {
    return this.invalidate(workspaceQueryKeys.list());
  }

  afterWorkspaceDelete(workspaceId: string) {
    this.remove(workspaceQueryKeys.detail(workspaceId), workspaceQueryKeys.statsFor(workspaceId));

    return this.afterWorkspaceListChange();
  }

  afterWorkspaceUpdate(workspaceId: string) {
    return this.invalidate(
      workspaceQueryKeys.list(),
      workspaceQueryKeys.statsFor(workspaceId),
      workspaceQueryKeys.detail(workspaceId),
    );
  }

  private async invalidate(...queryKeys: QueryKey[]) {
    await Promise.all(
      queryKeys.map((queryKey) => this.queryClient.invalidateQueries({ queryKey })),
    );
  }

  private remove(...queryKeys: QueryKey[]) {
    for (const queryKey of queryKeys) {
      this.queryClient.removeQueries({ queryKey });
    }
  }
}
