export interface List<T> {
  items: T[];
}

export type Data = string | number | boolean | Dict;
export type Dict = { [key: string]: Data };
export async function apiRequest<T>(
  url: string,
  method: string,
  body: BodyInit | null = null,
): Promise<T> {
  const apiUrl = import.meta.env.VITE_API_URL;

  const headers: HeadersInit = {};

  let requestBody: BodyInit | null = null;

  if (body instanceof FormData) {
    requestBody = body;
  } else if (body !== null) {
    headers["Content-Type"] = "application/json";
    requestBody = JSON.stringify(body);
  }

  const response = await fetch(apiUrl + url, {
    method,
    headers,
    body: requestBody,
  });

  if (!response.ok) {
    const errorData = await response.json();
    throw new Error(
      `Request failed with status ${response.status}: ${JSON.stringify(errorData)}`,
    );
  }

  return response.json();
}
