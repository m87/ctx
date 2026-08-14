export type ProjectPicker = {
  contextName: string;
  projectQuery: string;
};

export function parseProjectPicker(searchTerm: string): ProjectPicker | null {
  const hashIndex = searchTerm.lastIndexOf('#');
  if (hashIndex < 0) {
    return null;
  }
  return {
    contextName: searchTerm.slice(0, hashIndex).trim(),
    projectQuery: searchTerm.slice(hashIndex + 1).trim(),
  };
}
