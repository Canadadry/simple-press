type Dict = Record<string, unknown>;

export async function parseFormData(req: Request, def: Dict): Promise<Dict> {
  const text = await req.text();
  const params = new URLSearchParams(text);
  return extractData(params, def, "");
}

function extractData(values: URLSearchParams, def: Dict, prefix: string): Dict {
  const result: Dict = {};

  for (const [key, val] of Object.entries(def)) {
    const fullKey = prefix ? `${prefix}.${key}` : key;

    if (isObject(val)) {
      result[key] = extractData(values, val as Dict, fullKey);
      continue;
    }

    if (typeof val === "string") {
      if (values.has(fullKey)) {
        result[key] = values.get(fullKey) ?? val;
      } else {
        result[key] = val;
      }
      continue;
    }

    if (typeof val === "boolean") {
      if (values.has(fullKey)) {
        result[key] = (values.get(fullKey) ?? "") === "true";
      } else {
        result[key] = val;
      }
      continue;
    }

    if (typeof val === "number") {
      if (values.has(fullKey)) {
        const raw = values.get(fullKey) ?? "";
        if (raw === "") {
          result[key] = 0;
          continue;
        }
        const num = Number(raw);
        if (Number.isNaN(num)) {
          throw new Error(`invalid number for key ${fullKey}`);
        }
        result[key] = num;
      } else {
        result[key] = val;
      }
      continue;
    }

    result[key] = val;
  }

  return result;
}

function isObject(value: unknown): value is Dict {
  return typeof value === "object" && value !== null && !Array.isArray(value);
}
