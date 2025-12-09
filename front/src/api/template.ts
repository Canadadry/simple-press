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

export async function postTemplateAdd(data: Template) {
  return apiRequest<Template>(`${TEMPLATE_BASE_URL}/add`, "POST", data);
}

export async function getTemplateEdit(name: string) {
  return apiRequest<Template>(`${TEMPLATE_BASE_URL}/${name}/edit`, "GET");
}

export async function postTemplateEditBlockEdit(l: Template) {
  return apiRequest<Template>(`${TEMPLATE_BASE_URL}/${l.name}/edit`, "POST", {
    name: l.name,
    content: l.content,
  });
}
