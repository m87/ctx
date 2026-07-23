import { Component, signal } from '@angular/core';

type LinkRuleField = 'regexp' | 'link';

type LinkRuleDraft = {
  id: number;
  regexp: string;
  link: string;
};

@Component({
  selector: 'ctx-sidebar-settings-link-rules-section',
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
        class="h-10 rounded-md border px-4 text-[14px] font-medium text-foreground hover:bg-muted/50"
        (click)="addRule()"
      >
        Add link rule
      </button>

      <div class="space-y-3">
        @for (rule of linkRules(); track rule.id; let index = $index) {
          <div class="rounded-lg border p-4">
            <div class="mb-3 flex items-center justify-between gap-3">
              <div class="font-medium text-foreground">Rule {{ index + 1 }}</div>
              <button
                type="button"
                class="h-8 rounded-md px-3 text-[12px] font-medium text-destructive hover:bg-destructive/10"
                [attr.aria-label]="'Remove rule ' + (index + 1)"
                (click)="removeRule(rule.id)"
              >
                Remove
              </button>
            </div>

            <div class="grid gap-4 sm:grid-cols-2">
              <label class="space-y-1.5 text-[13px]">
                <span class="font-medium text-foreground">RegExp</span>
                <input
                  type="text"
                  class="h-10 w-full rounded-md border border-border bg-background px-3 text-foreground placeholder:text-muted-foreground"
                  placeholder="([A-Z]{3})"
                  [value]="rule.regexp"
                  (input)="updateRule(rule.id, 'regexp', getInputValue($event))"
                />
              </label>

              <label class="space-y-1.5 text-[13px]">
                <span class="font-medium text-foreground">Link</span>
                <input
                  type="text"
                  class="h-10 w-full rounded-md border border-border bg-background px-3 text-foreground placeholder:text-muted-foreground"
                  placeholder="http://test.pl/$1"
                  [value]="rule.link"
                  (input)="updateRule(rule.id, 'link', getInputValue($event))"
                />
              </label>
            </div>

            <div class="mt-3 text-[12px] text-muted-foreground">
              Use capture groups such as <code>$1</code> in the link template.
            </div>
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
  private nextRuleId = 2;

  readonly linkRules = signal<LinkRuleDraft[]>([
    {
      id: 1,
      regexp: '',
      link: '',
    },
  ]);

  addRule(): void {
    this.linkRules.update((rules) => [
      {
        id: this.nextRuleId++,
        regexp: '',
        link: '',
      },
      ...rules,
    ]);
  }

  updateRule(id: number, field: LinkRuleField, value: string): void {
    this.linkRules.update((rules) =>
      rules.map((rule) => (rule.id === id ? { ...rule, [field]: value } : rule)),
    );
  }

  removeRule(id: number): void {
    this.linkRules.update((rules) => rules.filter((rule) => rule.id !== id));
  }

  getInputValue(event: Event): string {
    return (event.target as HTMLInputElement).value;
  }
}
