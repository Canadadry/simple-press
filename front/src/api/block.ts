import { apiRequest } from "./api";
import type { List, Dict } from "./api";

const BASE_URL = "/admin/blocks";

export interface Block {
  name: string;
  content: string;
  definition: Dict;
}

export interface ValidationErrors {
  [key: string]: string[];
}

export async function getBlockList() {
  return apiRequest<List<Block>>(`${BASE_URL}`, "GET");
}

export async function postBlockAdd(data: Block) {
  return apiRequest<Block>(`${BASE_URL}/add`, "POST", data);
}

export async function getBlockEdit(name: string) {
  return apiRequest<Block>(`${BASE_URL}/${name}/edit`, "GET");
}

export async function postBlockEdit(previous_name: string, l: Block) {
  return apiRequest<Block>(`${BASE_URL}/${previous_name}/edit`, "POST", {
    name: l.name,
    content: l.content,
    definition: l.definition,
  });
}
