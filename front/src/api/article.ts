import { apiRequest } from "./api.ts";
import type { List } from "./api.ts";

const ARTICLE_BASE_URL = "/admin/articles";

export interface Article {
  title: string;
  author: string;
  image: string;
  slug: string;
  content: string;
  draft: boolean;
  layout?: number;
  layout_id?: number;
  blocks?: Array<{ name: string; value: number }>;
}

export interface ValidationErrors {
  [key: string]: string[];
}

export async function getArticleList() {
  return apiRequest<List<Article>>(`${ARTICLE_BASE_URL}`, "GET");
}

export async function postArticleAdd(data: Article) {
  return apiRequest<Article>(`${ARTICLE_BASE_URL}/add`, "POST", data);
}

export async function getArticleEdit(slug: string) {
  return apiRequest<Article>(`${ARTICLE_BASE_URL}/${slug}/edit`, "GET");
}

export async function postArticleEditMetadata(slug: string, metadata: Article) {
  return apiRequest<Article>(
    `${ARTICLE_BASE_URL}/${slug}/edit/metadata`,
    "POST",
    { ...metadata, layout: metadata.layout_id },
  );
}

export async function postArticleEditContent(slug: string, content: string) {
  return apiRequest<Article>(
    `${ARTICLE_BASE_URL}/${slug}/edit/content`,
    "POST",
    {
      content,
    },
  );
}

export async function postArticleEditBlockEdit(
  slug: string,
  blockData: object,
) {
  return apiRequest<Article>(
    `${ARTICLE_BASE_URL}/${slug}/edit/block_edit`,
    "POST",
    blockData,
  );
}

export async function postArticleEditBlockAdd(slug: string, blockData: object) {
  return apiRequest<Article>(
    `${ARTICLE_BASE_URL}/${slug}/edit/block_add`,
    "POST",
    blockData,
  );
}

export async function getArticlePreview(slug: string) {
  return apiRequest<Article>(`${ARTICLE_BASE_URL}/${slug}/preview`, "GET");
}
