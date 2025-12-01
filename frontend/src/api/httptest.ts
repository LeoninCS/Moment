import http from "./http";

export interface HttpTestPayload {
  method: string;
  url: string;
  header: Record<string, string>;
  body: string;
}

export interface HttpTestResponse {
  status: string;
  status_code: number;
  header: Record<string, string[]>;
  body: string;
}

export function sendHttpTest(payload: HttpTestPayload) {
  return http.post<HttpTestResponse>("/httptest/do", payload);
}
