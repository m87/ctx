import { ApplicationConfig, provideBrowserGlobalErrorListeners } from '@angular/core';
import { provideRouter } from '@angular/router';
import { withNgxsStoragePlugin } from '@ngxs/storage-plugin';
import { provideStore } from '@ngxs/store';

import { routes } from './app.routes';
import { provideHttpClient } from '@angular/common/http';
import { provideTanStackQuery, QueryClient } from '@tanstack/angular-query-experimental';
import { WorkspaceState } from './sidebar/workspace.state';

export const appConfig: ApplicationConfig = {
  providers: [
    provideBrowserGlobalErrorListeners(),
    provideRouter(routes),
    provideHttpClient(),
    provideStore(
      [WorkspaceState],
      withNgxsStoragePlugin({
        keys: ['workspace.selectedWorkspaceId'],
      }),
    ),
    provideTanStackQuery(new QueryClient()),
  ],
};
