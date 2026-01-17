export function removeLastPathSegment(path: string): string {
  if (!path) {
    return path;
  }

  const normalizedPath = path.replace(/\/+$/, "");
  const lastSlashIndex = normalizedPath.lastIndexOf("/");

  if (lastSlashIndex === -1) {
    return "";
  }

  return normalizedPath.substring(0, lastSlashIndex);
}
