import { apiRequest } from "./api.ts";
import type { List } from "./api.ts";

const TEMPLATE_BASE_URL = "/admin/templates";

export interface Template {
  id: number;
  name: string;
  content: string;
}

export interface ValidationErrors {
  [key: string]: string[];
}

export async function getTemplateList() {
  return apiRequest<List<Template>>(`${TEMPLATE_BASE_URL}`, "GET");
}

export async function postTemplateAdd(name: string) {
  return apiRequest<Template>(`${TEMPLATE_BASE_URL}/add`, "POST", {
    name: name,
  });
}

export async function getTemplateEdit(name: string) {
  return apiRequest<Template>(`${TEMPLATE_BASE_URL}/${name}/edit`, "GET");
}

export async function postTemplateEdit(previous_name: string, l: Template) {
  return apiRequest<Template>(
    `${TEMPLATE_BASE_URL}/${previous_name}/edit`,
    "POST",
    {
      name: l.name,
      content: l.content,
    },
  );
}
