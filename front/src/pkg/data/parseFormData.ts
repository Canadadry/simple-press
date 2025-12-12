export type Data = string | number | boolean | Dict;
export type Dict = { [key: string]: Data };
export type FlatDict = { [key: string]: string };

export function extractData(values: FlatDict, def: Dict, prefix: string): Dict {
  const result: Dict = {};

  for (const [key, val] of Object.entries(def)) {
    const fullKey = prefix ? `${prefix}.${key}` : key;
    switch (typeof val) {
      case "object":
        result[key] = extractData(values, val as Dict, fullKey);
        break;
      case "number":
        result[key] = values[fullKey] ? Number(values[fullKey]) : val;
        break;
      case "boolean":
        result[key] = values[fullKey] ? values[fullKey] === "true" : val;
        break;
      case "string":
        result[key] = values[fullKey] ? values[fullKey] : val;
        break;
      default:
        result[key] = values[fullKey] ? JSON.stringify(values[fullKey]) : val;
    }
  }
  return result;
}
