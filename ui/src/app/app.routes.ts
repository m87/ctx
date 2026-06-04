import { Routes } from '@angular/router';

export const routes: Routes = [
  {
    path: '',
    pathMatch: 'full',
    redirectTo: 'day',
  },
  {
    path: 'day',
    loadComponent: () => import('./day/day.component').then((m) => m.DayComponent),
  },
  {
    path: 'day/:date',
    loadComponent: () => import('./day/day.component').then((m) => m.DayComponent),
  },
  {
    path: 'workspace',
    loadComponent: () => import('./workspace/workspace.component').then((m) => m.WorkspaceComponent),
  },
  {
    path: 'workspace/:id',
    loadComponent: () => import('./workspace/workspace.component').then((m) => m.WorkspaceComponent),
  },
  {
    path: 'context/:id',
    loadComponent: () => import('./context/context.component').then((m) => m.ContextComponent),
  },
];
