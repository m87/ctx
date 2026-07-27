import { Component, computed, inject, linkedSignal, signal } from '@angular/core';
import { NgIcon, provideIcons } from '@ng-icons/core';
import { lucideCheck, lucidePencil, lucideTrash2, lucideX } from '@ng-icons/lucide';
import { Store } from '@ngxs/store';
import { injectMutation, injectQuery } from '@tanstack/angular-query-experimental';
import { WorkspaceQueries } from '../../api/workspace.quries';
import { WorkspaceState } from './workspace.state';
import { WorkspaceMutations } from '../../api/workspace.mutations';
import { LinkRule } from '../../api/workspace.service';

type LinkRuleField = keyof LinkRule;

type LinkRuleEdit = {
  originalRegexp: string;
  isNew: boolean;
  draft: LinkRule;
};

@Component({
  selector: 'ctx-sidebar-settings-link-rules-section',
  imports: [NgIcon],
  providers: [provideIcons({ lucideCheck, lucidePencil, lucideTrash2, lucideX })],
  template: `
    <div class="space-y-6">
      <div class="space-y-1.5">
        <div class="text-foreground font-medium text-[15px]">Link rules</div>
        <div class="text-[13px] sm:text-[14px]">
          Turn matching text into links by pairing a regular expression with a link template.
        </div>
      </div>

      <button
        type="button"
        class="h-10 rounded-md border px-4 text-[14px] font-medium text-foreground hover:bg-muted/50 disabled:cursor-not-allowed disabled:opacity-50"
        [disabled]="!canAddRule()"
        (click)="addRule()"
      >
        Add link rule
      </button>

      <div class="space-y-3">
        @for (rule of linkRules(); track $index; let index = $index) {
          <div class="rounded-lg border p-4">
            <div class="mb-3 flex items-center justify-between gap-3">
              <div class="font-medium text-foreground">Rule {{ index + 1 }}</div>
              <div class="flex items-center gap-1">
                @if (isEditingRule(rule.regexp)) {
                  <button
                    type="button"
                    class="h-8 w-8 rounded-md text-foreground hover:bg-muted/60 disabled:cursor-not-allowed disabled:opacity-50 flex items-center justify-center"
                    [attr.aria-label]="'Save rule ' + (index + 1)"
                    title="Save"
                    [disabled]="!canSaveRule()"
                    (click)="saveRule()"
                  >
                    <ng-icon name="lucideCheck" class="text-[14px]"></ng-icon>
                  </button>
                  <button
                    type="button"
                    class="h-8 w-8 rounded-md text-muted-foreground hover:bg-muted/60 flex items-center justify-center"
                    [attr.aria-label]="'Cancel editing rule ' + (index + 1)"
                    title="Cancel"
                    (click)="cancelEdit()"
                  >
                    <ng-icon name="lucideX" class="text-[14px]"></ng-icon>
                  </button>
                } @else {
                  <button
                    type="button"
                    class="h-8 w-8 rounded-md text-foreground hover:bg-muted/60 flex items-center justify-center"
                    [attr.aria-label]="'Edit rule ' + (index + 1)"
                    title="Edit"
                    [disabled]="editingRule() !== null || isSaving()"
                    (click)="editRule(rule)"
                  >
                    <ng-icon name="lucidePencil" class="text-[14px]"></ng-icon>
                  </button>
                  <button
                    type="button"
                    class="h-8 w-8 rounded-md text-destructive hover:bg-destructive/10 flex items-center justify-center"
                    [attr.aria-label]="'Remove rule ' + (index + 1)"
                    title="Remove"
                    [disabled]="editingRule() !== null || isSaving()"
                    (click)="removeRule(rule.regexp)"
                  >
                    <ng-icon name="lucideTrash2" class="text-[14px]"></ng-icon>
                  </button>
                }
              </div>
            </div>

            <div class="grid gap-4 sm:grid-cols-2">
              <label class="space-y-1.5 text-[13px]">
                <span class="font-medium text-foreground">RegExp</span>
                <input
                  type="text"
                  class="h-10 w-full rounded-md border border-border bg-background px-3 text-foreground placeholder:text-muted-foreground"
                  placeholder="([A-Z]{3})"
                  [value]="getRuleFieldValue(rule, 'regexp')"
                  [disabled]="!isEditingRule(rule.regexp)"
                  (input)="updateDraft('regexp', getInputValue($event))"
                />
              </label>

              <label class="space-y-1.5 text-[13px]">
                <span class="font-medium text-foreground">Link</span>
                <input
                  type="text"
                  class="h-10 w-full rounded-md border border-border bg-background px-3 text-foreground placeholder:text-muted-foreground"
                  placeholder="http://test.pl/$1"
                  [value]="getRuleFieldValue(rule, 'link')"
                  [disabled]="!isEditingRule(rule.regexp)"
                  (input)="updateDraft('link', getInputValue($event))"
                />
              </label>
            </div>

            <div class="mt-3 text-[12px] text-muted-foreground">
              Use capture groups such as <code>$1</code> in the link template.
            </div>
            @if (isEditingRule(rule.regexp) && !canSaveRule()) {
              <div class="mt-2 text-[12px] text-destructive">
                A rule with this RegExp already exists.
              </div>
            }
          </div>
        } @empty {
          <div class="rounded-lg border border-dashed p-6 text-center text-[13px]">
            No link rules yet. Add a rule to get started.
          </div>
        }
      </div>
    </div>
  `,
})
export class SidebarSettingsLinkRulesSectionComponent {
  private store = inject(Store);
  private workspaceQueries = inject(WorkspaceQueries);
  private workspaceMutations = inject(WorkspaceMutations);
  private updateWorkspaceMutation = injectMutation(() => this.workspaceMutations.update());
  private activeWorkspaceId = this.store.selectSignal(WorkspaceState.selectedWorkspaceId);
  private getWorkspaceQuery = injectQuery(() => {
    const workspaceId = this.activeWorkspaceId();

    return {
      ...this.workspaceQueries.get(workspaceId ?? ''),
      enabled: workspaceId !== null,
    };
  });

  readonly linkRules = linkedSignal<LinkRule[]>(() => {
    const workspace = this.getWorkspaceQuery.data();
    if (!workspace) {
      return [];
    }

    return [
      ...new Map(
        (workspace.properties.linkRules ?? []).map((rule) => [rule.regexp, rule]),
      ).values(),
    ];
  });

  readonly editingRule = signal<LinkRuleEdit | null>(null);
  readonly isSaving = computed(() => this.updateWorkspaceMutation.isPending());

  readonly canAddRule = computed(
    () =>
      this.editingRule() === null &&
      !this.isSaving() &&
      !this.linkRules().some((rule) => rule.regexp === ''),
  );

  readonly canSaveRule = computed(() => {
    const edit = this.editingRule();
    if (edit === null) {
      return false;
    }

    return !this.linkRules().some(
      (rule) => rule.regexp === edit.draft.regexp && rule.regexp !== edit.originalRegexp,
    );
  });

  addRule(): void {
    if (!this.canAddRule()) {
      return;
    }

    const rule: LinkRule = {
      regexp: '',
      link: '',
    };

    this.linkRules.update((rules) => [rule, ...rules]);
    this.editingRule.set({
      originalRegexp: rule.regexp,
      isNew: true,
      draft: { ...rule },
    });
  }

  editRule(rule: LinkRule): void {
    if (this.editingRule() !== null || this.isSaving()) {
      return;
    }

    this.editingRule.set({
      originalRegexp: rule.regexp,
      isNew: false,
      draft: { ...rule },
    });
  }

  updateDraft(field: LinkRuleField, value: string): void {
    this.editingRule.update((edit) =>
      edit === null
        ? null
        : {
            ...edit,
            draft: {
              ...edit.draft,
              [field]: value,
            },
          },
    );
  }

  saveRule(): void {
    const edit = this.editingRule();
    if (edit === null || !this.canSaveRule() || this.isSaving()) {
      return;
    }

    this.linkRules.update((rules) =>
      rules.map((rule) => (rule.regexp === edit.originalRegexp ? { ...edit.draft } : rule)),
    );
    this.editingRule.set(null);
    this.patchWorkspaceProperties();
  }

  cancelEdit(): void {
    const edit = this.editingRule();
    if (edit?.isNew) {
      this.linkRules.update((rules) => rules.filter((rule) => rule.regexp !== edit.originalRegexp));
    }

    this.editingRule.set(null);
  }

  removeRule(regexp: string): void {
    if (this.editingRule() !== null || this.isSaving()) {
      return;
    }

    this.linkRules.update((rules) => rules.filter((rule) => rule.regexp !== regexp));
    this.patchWorkspaceProperties();
  }

  getInputValue(event: Event): string {
    return (event.target as HTMLInputElement).value;
  }

  isEditingRule(regexp: string): boolean {
    return this.editingRule()?.originalRegexp === regexp;
  }

  getRuleFieldValue(rule: LinkRule, field: LinkRuleField): string {
    const edit = this.editingRule();
    return edit?.originalRegexp === rule.regexp ? edit.draft[field] : rule[field];
  }

  private patchWorkspaceProperties(): void {
    const workspaceId = this.activeWorkspaceId();
    if (!workspaceId) {
      return;
    }

    const workspace = this.getWorkspaceQuery.data();
    if (!workspace) {
      return;
    }

    const patchedWorkspace = {
      ...workspace,
      properties: {
        ...workspace.properties,
        linkRules: this.linkRules(),
      },
    };

    this.updateWorkspaceMutation.mutate(patchedWorkspace);
  }
}
