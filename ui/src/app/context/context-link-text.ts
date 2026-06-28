export type ContextLinkTextPart = {
  text: string;
  href?: string;
};

export type ContextLinkRule = {
  regex?: string | null;
  linkTemplate?: string | null;
};

type ContextLinkMatch = {
  start: number;
  end: number;
  order: number;
  text: string;
  href: string;
};

export function linkifyContextText(
  text: string | null | undefined,
  rules: readonly ContextLinkRule[] | null | undefined,
): ContextLinkTextPart[] {
  const source = text ?? '';

  if (source.length === 0 || !rules || rules.length === 0) {
    return [{ text: source }];
  }

  const matches = findContextLinkMatches(source, rules);
  if (matches.length === 0) {
    return [{ text: source }];
  }

  const parts: ContextLinkTextPart[] = [];
  let lastIndex = 0;
  for (const match of matches) {
    if (match.start < lastIndex) {
      continue;
    }

    if (match.start > lastIndex) {
      parts.push({ text: source.slice(lastIndex, match.start) });
    }

    parts.push({ text: match.text, href: match.href });
    lastIndex = match.end;
  }

  if (lastIndex < source.length) {
    parts.push({ text: source.slice(lastIndex) });
  }

  return parts.length > 0 ? parts : [{ text: source }];
}

function findContextLinkMatches(
  source: string,
  rules: readonly ContextLinkRule[],
): ContextLinkMatch[] {
  const matches: ContextLinkMatch[] = [];

  rules.forEach((rule, order) => {
    const pattern = rule.regex?.trim() ?? '';
    const template = rule.linkTemplate?.trim() ?? '';
    if (pattern.length === 0 || template.length === 0) {
      return;
    }

    let regex: RegExp;
    try {
      regex = new RegExp(pattern, 'g');
    } catch {
      return;
    }

    let match: RegExpExecArray | null;
    while ((match = regex.exec(source)) !== null) {
      const matchedText = match[0];
      if (matchedText.length === 0) {
        regex.lastIndex = match.index + 1;
        continue;
      }

      const href = normalizeContextLinkHref(applyContextLinkTemplate(match, template));
      if (!href) {
        continue;
      }

      matches.push({
        start: match.index,
        end: match.index + matchedText.length,
        order,
        text: matchedText,
        href,
      });
    }
  });

  return matches.sort((left, right) => {
    if (left.start !== right.start) {
      return left.start - right.start;
    }
    if (left.order !== right.order) {
      return left.order - right.order;
    }
    return right.end - left.end;
  });
}

function applyContextLinkTemplate(match: RegExpExecArray, template: string): string {
  return template.replace(
    /\$(\d+)|\$\{([A-Za-z][A-Za-z0-9_]*)\}|\$<([A-Za-z][A-Za-z0-9_]*)>/g,
    (placeholder: string, index?: string, braceName?: string, angleName?: string) => {
      if (index !== undefined) {
        return match[Number(index)] ?? '';
      }

      const groupName = braceName ?? angleName;
      if (groupName) {
        return match.groups?.[groupName] ?? '';
      }

      return placeholder;
    },
  );
}

function normalizeContextLinkHref(href: string): string | undefined {
  const trimmed = href.trim();
  if (/^https?:\/\//i.test(trimmed)) {
    return trimmed;
  }

  if (trimmed.startsWith('/')) {
    return trimmed;
  }

  return undefined;
}
