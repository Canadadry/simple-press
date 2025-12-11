export interface DynamicFormRenderer {
  beginForm(): void;
  endForm(): void;
  beginObject(label: string): void;
  endObject(): void;
  input(label: string, name: string, inputType: string, value: string): void;
  checkbox(label: string, name: string, checked: boolean): void;
}

export function render(data: unknown, r: DynamicFormRenderer): Error | null {
  try {
    r.beginForm();
    renderValue(data, "", "", r);
    r.endForm();
    return null;
  } catch (e) {
    return e instanceof Error ? e : new Error(String(e));
  }
}

function renderValue(
  val: unknown,
  key: string,
  path: string,
  r: DynamicFormRenderer,
): void {
  // ------------------------
  // case: object / map[string]any
  // ------------------------
  if (isPlainObject(val)) {
    if (key !== "") {
      r.beginObject(key);
    }

    const entries = Object.entries(val);
    const sortedKeys = entries.map(([k]) => k).sort();

    for (const k of sortedKeys) {
      const nextVal = (val as Record<string, unknown>)[k];
      const fullPath = joinPath(path, k);
      renderValue(nextVal, k, fullPath, r);
    }

    if (key !== "") {
      r.endObject();
    }
    return;
  }

  // ------------------------
  // case: string
  // ------------------------
  if (typeof val === "string") {
    r.input(key, path, "text", val);
    return;
  }

  // ------------------------
  // case: numeric types (Go: float64, float32, int, etc.)
  // ------------------------
  if (typeof val === "number") {
    r.input(key, path, "number", String(val));
    return;
  }

  // ------------------------
  // case: boolean
  // ------------------------
  if (typeof val === "boolean") {
    r.checkbox(key, path, val);
    return;
  }

  // ------------------------
  // unsupported type
  // ------------------------
  throw new Error(`at ${path} not handle value type ${typeof val}`);
}

function joinPath(...parts: string[]): string {
  const nonEmpty = parts.filter((p) => p !== "");
  return nonEmpty.join(".");
}

function isPlainObject(v: unknown): v is Record<string, unknown> {
  return typeof v === "object" && v !== null && !Array.isArray(v);
}
