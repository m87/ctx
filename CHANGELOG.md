## 0.5.3
- enhancement: include archived contexts in the search results

## 0.5.2
- fix: invalidate cache on context switch
- fix: change intervals sorting to descending order
- feature: add context archiving feature
  - archived contexts are hidden from the search and context list, but can be accessed via daily summary or workspace view
  - archived contexts can be restored or permanently deleted
  - archived contexts are read-only and cannot be modified
  - archived contexts are not included in the summary and time distribution calculations

## 0.5.1
- dev: nod library update

## 0.5.0

- feature: workspaces — create, rename, select, and delete workspaces to organize contexts
- enhancement: error handling — display toast notifications for failed queries and mutations
- enhancement: API errors — return consistent error codes and descriptions from server endpoints
- enhancement: data integrity tool — check for broken contexts, intervals, and workspaces, and attempt auto-repair
- dev: python script to generate test database with random contexts, intervals, and workspaces

On the first run after upgrading, a default workspace is created and all existing contexts are assigned to it. You can change workspace name later.
You can verify this by going to the settings view and checking data integrity. Most of the issues may be resolved with auto-repair, but if you encounter any problems, please report them in the GitHub repository.

## 0.4.0
- feature: settings view — currently supports first day of week and light/dark theme
- enhancement: UI — display application version number
- fix: wrong 'top context' widget layout on mobile
- style: favicon

## 0.3.2

- enhancement: search UI — create tile, grouped results (Today/date), per-context time badges
- enhancement: daily summary — show first start and last end times

## 0.3.1

- fix: long context names overlap calendar component
- refactor: merge ui repo into main repo and remove submodule
- fix: time distribution calculation for a given day
- fix: hidden timeline on mobile web browser
