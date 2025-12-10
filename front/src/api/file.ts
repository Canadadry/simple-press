import { apiRequest } from "./api.ts";
import type { List } from "./api.ts";

const ARTICLE_BASE_URL = "/admin/files";

export interface File {
  name: string;
}

export interface ValidationErrors {
  [key: string]: string[];
}

export async function getFileList() {
  return apiRequest<List<File>>(`${ARTICLE_BASE_URL}`, "GET");
}
