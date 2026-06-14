export class ApiError extends Error {
  status: number;
  code: string;
  fields?: { field: string; message: string }[];
  constructor(
    status: number,
    code: string,
    message: string,
    fields?: { field: string; message: string }[],
  ) {
    super(message);
    this.status = status;
    this.code = code;
    this.fields = fields;
  }
}

export async function request<T>(
  method: string,
  path: string,
  body?: unknown,
): Promise<T> {
  const res = await fetch(`/api${path}`, {
    method,
    headers:
      body !== undefined ? { "content-type": "application/json" } : undefined,
    body: body !== undefined ? JSON.stringify(body) : undefined,
  });
  const text = await res.text();
  const data = text ? JSON.parse(text) : undefined;
  if (!res.ok) {
    const e = data?.error ?? {};
    throw new ApiError(
      res.status,
      e.code ?? "error",
      e.message ?? res.statusText,
      e.fields,
    );
  }
  return data as T;
}
