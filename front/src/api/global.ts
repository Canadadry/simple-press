import { apiRequest } from "./api";
import type { Dict } from "./api";

export interface GDefinition {
  definition: Dict;
}

export interface GData {
  data: Dict;
}

export interface ValidationErrors {
  [key: string]: string[];
}

export async function getGlobalDefinition(section: string): Promise<Dict> {
  section = section.replace(/\./g, "_");
  const all = await apiRequest<GDefinition>("/admin/global/definition", "GET");
  return all.definition[section] as Dict;
}

export async function patchGlobalDefinition(section: string, definition: Dict) {
  section = section.replace(/\./g, "_");
  const previous = await apiRequest<GDefinition>(
    "/admin/global/definition",
    "GET",
  );
  return apiRequest<GDefinition>("/admin/global/definition", "PATCH", {
    definition: {
      ...previous.definition,
      [section]: definition,
    },
  });
}

export async function getGlobalData(): Promise<Dict> {
  const all = await apiRequest<GData>("/admin/global/data", "GET");
  return all.data as Dict;
}
export async function patchGlobalData(data: Dict) {
  return apiRequest<GData>("/admin/global/data", "PATCH", {
    data: data,
  });
}
