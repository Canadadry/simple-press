import { apiRequest } from "./api";
import type { List } from "./api";

const BASE_URL = "/admin/files";

export interface File {
  name: string;
}

export interface FileTree {
  path: string;
  files: string[];
  folders: string[];
}

export interface ValidationErrors {
  [key: string]: string[];
}

export async function getFileList() {
  return apiRequest<List<File>>(`${BASE_URL}`, "GET");
}

export async function getFileTree(path: string) {
  if (path !== "" && path[0] != "/") {
    path = "/" + path;
  }
  return apiRequest<FileTree>(`${BASE_URL}/tree${path}`, "GET");
}

function cleanPath(path: string): string {
  if (!path) {
    return "";
  }

  return path
    .replace(/\/{2,}/g, "/") // remplace // par /
    .replace(/^\/+/, "") // supprime / au début
    .replace(/\/+$/, ""); // supprime / à la fin
}

export async function postFile(file: Blob, filename: string, archive: boolean) {
  const formData = new FormData();
  formData.append("content", file);
  formData.append("name", cleanPath(filename));
  if (archive) {
    formData.append("archive", "true");
  }
  return apiRequest<void>(`${BASE_URL}/add`, "POST", formData);
}

export async function deleteFile(path: string) {
  return apiRequest<void>(`${BASE_URL}/${path}/delete`, "DELETE");
}
