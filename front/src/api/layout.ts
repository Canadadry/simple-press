import { apiRequest } from "./api";
import type { List } from "./api";

const LAYOUT_BASE_URL = "/admin/layouts";

export interface Layout {
  id: number;
  name: string;
  content: string;
}

export interface ValidationErrors {
  [key: string]: string[];
}

export async function getLayoutList() {
  return apiRequest<List<Layout>>(`${LAYOUT_BASE_URL}`, "GET");
}

export async function postLayoutAdd(name: string) {
  return apiRequest<Layout>(`${LAYOUT_BASE_URL}/add`, "POST", { name: name });
}

export async function getLayoutEdit(name: string) {
  return apiRequest<Layout>(`${LAYOUT_BASE_URL}/${name}/edit`, "GET");
}

export async function postLayoutEdit(name: string, l: Layout) {
  return apiRequest<Layout>(`${LAYOUT_BASE_URL}/${name}/edit`, "POST", {
    name: l.name,
    content: l.content,
  });
}
