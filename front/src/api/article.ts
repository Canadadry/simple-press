import { apiRequest } from "./api";
import type { List, Dict } from "./api";

const ARTICLE_BASE_URL = "/admin/articles";

export interface BlockData {
  id: number;
  name: string;
  data: Dict;
}

export interface Article {
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

export async function postArticleEditBlockEdit(slug: string, data: BlockData) {
  return apiRequest<Article>(
    `${ARTICLE_BASE_URL}/${slug}/edit/block_edit`,
    "POST",
    {
      block_id: data.id,
      block_data: data.data,
    },
  );
}

export async function postArticleEditBlockAdd(slug: string, block: number) {
  return apiRequest<Article>(
    `${ARTICLE_BASE_URL}/${slug}/edit/block_add`,
    "POST",
    { new_block: block },
  );
}

export async function getArticlePreview(slug: string) {
  return apiRequest<Article>(`${ARTICLE_BASE_URL}/${slug}/preview`, "GET");
}
