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

export type LinkRuleValues = object;

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

  linkify(text: string, values: LinkRuleValues = {}): readonly LinkifiedTextPart[] {
    let parts: LinkifiedTextPart[] = [{ text }];
    const templateValues = { name: text, ...values };

    for (const rule of this.linkRules()) {
      const expression = this.createExpression(rule);
      if (expression === null) {
        continue;
      }

      parts = parts.flatMap((part) =>
        part.href === undefined
          ? this.applyRule(part.text, expression, rule.link, templateValues)
          : [part],
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

  private applyRule(
    text: string,
    expression: RegExp,
    linkTemplate: string,
    values: LinkRuleValues,
  ): LinkifiedTextPart[] {
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
        href: resolveLinkTemplate(linkTemplate, match, values),
      });
      cursor = match.index + match[0].length;
    }

    if (cursor < text.length) {
      parts.push({ text: text.slice(cursor) });
    }

    return parts.length > 0 ? parts : [{ text }];
  }
}

export function resolveLinkTemplate(
  template: string,
  match: RegExpExecArray,
  values: LinkRuleValues,
): string {
  const contextName = Object.prototype.hasOwnProperty.call(values, 'name')
    ? (values as Record<string, unknown>)['name']
    : undefined;
  const templateValues = {
    ...values,
    name: typeof contextName === 'string' ? contextName.replace(match[0], '').trim() : undefined,
  };
  const withValues = template.replace(
    /\$\{([A-Za-z_][A-Za-z0-9_]*(?:\.[A-Za-z_][A-Za-z0-9_]*)*)\}/g,
    (placeholder, placeholderName: string) => {
      const value = placeholderName.split('.').reduce<unknown>((current, segment) => {
        if (
          typeof current !== 'object' ||
          current === null ||
          !Object.prototype.hasOwnProperty.call(current, segment)
        ) {
          return undefined;
        }

        return (current as Record<string, unknown>)[segment];
      }, templateValues);

      if (typeof value !== 'string' && typeof value !== 'number' && typeof value !== 'boolean') {
        return placeholder;
      }

      return encodeURIComponent(String(value));
    },
  );

  return withValues.replace(/\$(\$|&|\d{1,2})/g, (_token, reference: string) => {
    if (reference === '$') {
      return '$';
    }
    if (reference === '&' || reference === '0') {
      return match[0];
    }

    return match[Number(reference)] ?? '';
  });
}
