import { apiRequest } from "./api";
import type { List } from "./api";

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

export async function postFile(file: Blob, filename: string, archive: boolean) {
  const formData = new FormData();
  formData.append("content", file);
  formData.append("name", filename);
  if (archive) {
    formData.append("archive", "true");
  }
  return apiRequest<void>(`${ARTICLE_BASE_URL}/add`, "POST", formData);
}

export async function deleteFile(path: string) {
  return apiRequest<void>(`${ARTICLE_BASE_URL}/${path}/delete`, "DELETE");
}
