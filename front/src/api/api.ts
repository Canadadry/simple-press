export interface List<T> {
  items: T[];
}

export type Data = string | number | boolean | Dict;
export type Dict = { [key: string]: Data };

export async function apiRequest<T>(
  url: string,
  method: string,
  body: object | null = null,
): Promise<T> {
  const headers = {
    "Content-Type": "application/json",
  };

  const apiUrl = import.meta.env.VITE_API_URL;
  const response = await fetch(apiUrl + url, {
    method,
    headers,
    body: body ? JSON.stringify(body) : null,
  });

  if (!response.ok) {
    const errorData = await response.json();
    throw new Error(
      `Request failed with status ${response.status}: ${JSON.stringify(errorData)}`,
    );
  }

  return response.json();
}
