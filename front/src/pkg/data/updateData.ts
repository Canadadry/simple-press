import { Dict } from "../../api/api";

export function updateData(data: Dict, fullPath: string, value: string): Dict {
  const parts = fullPath.split(".");
  const result: Dict = JSON.parse(JSON.stringify(data));
  let current: Dict = result;

  for (let i = 0; i < parts.length - 1; i++) {
    const key = parts[i];
    const next = current[key];
    if (typeof next != "object") {
      return result;
    }
    current = next;
  }
  setValue(current, parts[parts.length - 1], value);
  return result;
}

function setValue(obj: Dict, key: string, value: string) {
  if (!(key in obj)) {
    return;
  }
  switch (typeof obj[key]) {
    case "string":
      obj[key] = value;
      break;
    case "number":
      obj[key] = Number(value);
      break;
    case "boolean":
      obj[key] = value === "true";
      break;
    default:
      obj[key] = value;
      break;
  }
}
