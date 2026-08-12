import { computed, inject, Injectable } from '@angular/core';
import { Store } from '@ngxs/store';
import { injectQuery } from '@tanstack/angular-query-experimental';
import { WorkspaceQueries } from '../../api/workspace/workspace.queries';
import { LinkRule } from '../../api/workspace/workspace.service';
import { WorkspaceState } from '../sidebar/workspace.state';

export type LinkifiedTextPart = {
  text: string;
  href?: string;
};

@Injectable({
  providedIn: 'root',
})
export class LinkRulesService {
  private store = inject(Store);
  private workspaceQueries = inject(WorkspaceQueries);
  private activeWorkspaceId = this.store.selectSignal(WorkspaceState.selectedWorkspaceId);
  private workspaceQuery = injectQuery(() => {
    const workspaceId = this.activeWorkspaceId();

    return {
      ...this.workspaceQueries.get(workspaceId ?? ''),
      enabled: workspaceId !== null,
    };
  });

  private readonly linkRules = computed<readonly LinkRule[]>(
    () => this.workspaceQuery.data()?.properties?.linkRules ?? [],
  );

  linkify(text: string): readonly LinkifiedTextPart[] {
    let parts: LinkifiedTextPart[] = [{ text }];

    for (const rule of this.linkRules()) {
      const expression = this.createExpression(rule);
      if (expression === null) {
        continue;
      }

      parts = parts.flatMap((part) =>
        part.href === undefined ? this.applyRule(part.text, expression, rule.link) : [part],
      );
    }

    return parts;
  }

  private createExpression(rule: LinkRule): RegExp | null {
    if (rule.regexp.length === 0 || rule.link.length === 0) {
      return null;
    }

    try {
      return new RegExp(rule.regexp, 'g');
    } catch {
      return null;
    }
  }

  private applyRule(text: string, expression: RegExp, linkTemplate: string): LinkifiedTextPart[] {
    const parts: LinkifiedTextPart[] = [];
    let cursor = 0;
    let match: RegExpExecArray | null;

    expression.lastIndex = 0;
    while ((match = expression.exec(text)) !== null) {
      if (match[0].length === 0) {
        expression.lastIndex = match.index + 1;
        continue;
      }

      if (match.index > cursor) {
        parts.push({ text: text.slice(cursor, match.index) });
      }

      parts.push({
        text: match[0],
        href: this.resolveLink(linkTemplate, match),
      });
      cursor = match.index + match[0].length;
    }

    if (cursor < text.length) {
      parts.push({ text: text.slice(cursor) });
    }

    return parts.length > 0 ? parts : [{ text }];
  }

  private resolveLink(template: string, match: RegExpExecArray): string {
    return template.replace(/\$(\$|&|\d{1,2})/g, (_token, reference: string) => {
      if (reference === '$') {
        return '$';
      }
      if (reference === '&' || reference === '0') {
        return match[0];
      }

      return match[Number(reference)] ?? '';
    });
  }
}
