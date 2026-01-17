import { apiRequest } from "./api";
import type { List, Dict } from "./api";

const BASE_URL = "/admin/articles";

export interface BlockData {
  id: number;
  position: number;
  name: string;
  data: Dict;
}

export interface Article {
  id: number;
  title: string;
  author: string;
  image: string;
  slug: string;
  content: string;
  draft: boolean;
  layout_id: number;
  layouts: Array<{ name: string; value: number }>;
  blocks: Array<{ name: string; value: number }>;
  block_datas: Array<BlockData>;
}

export interface ArticleTree {
  path: string;
  articles: Article[];
  folders: string[];
}

export interface ValidationErrors {
  [key: string]: string[];
}

export async function getArticleList() {
  return apiRequest<List<Article>>(`${BASE_URL}`, "GET");
}

export async function getArticleTree(path: string) {
  if (path !== "" && path[0] != "/") {
    path = "/" + path;
  }
  return apiRequest<ArticleTree>(`${BASE_URL}/tree${path}`, "GET");
}

export async function postArticleAdd(article: {
  title: string;
  author: string;
}) {
  return apiRequest<Article>(`${BASE_URL}/add`, "POST", article);
}

export async function getArticle(id: number) {
  return apiRequest<Article>(`${BASE_URL}/${id}/edit`, "GET");
}

export async function postArticleEditMetadata(id: number, metadata: Article) {
  return apiRequest<Article>(`${BASE_URL}/${id}/edit/metadata`, "POST", {
    ...metadata,
    layout: metadata.layout_id,
  });
}

export async function postArticleEditContent(id: number, content: string) {
  return apiRequest<Article>(`${BASE_URL}/${id}/edit/content`, "POST", {
    content,
  });
}

export async function postArticleEditBlockEdiPosition(data: BlockData) {
  return apiRequest<Article>(`${BASE_URL}/block/${data.id}/edit`, "PATCH", {
    block_position: data.position,
  });
}
export async function postArticleEditBlockEdit(data: BlockData) {
  return apiRequest<Article>(`${BASE_URL}/block/${data.id}/edit`, "PATCH", {
    block_position: data.position,
    block_data: data.data,
  });
}

export async function deleteArticleEditBlockEdit(data: BlockData) {
  return apiRequest<void>(`${BASE_URL}/block/${data.id}/delete`, "DELETE");
}

export async function postArticleEditBlockAdd(
  id: number,
  block: number,
  position: number,
) {
  return apiRequest<Article>(`${BASE_URL}/${id}/edit/block_add`, "POST", {
    new_block: block,
    position: position,
  });
}

export async function getArticlePreview(id: number) {
  return apiRequest<Article>(`${BASE_URL}/${id}/preview`, "GET");
}
