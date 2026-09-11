import { getToken, exampleToken } from "@/lib/session";

const API_URL = (process.env.NEXT_PUBLIC_API_URL || "/api/v1").replace(/\/$/, "");

export type SduiEnvelope = {
  schema_version?: number;
  component: string;
  props: Record<string, unknown>;
};

export type ChatPart =
  | { type: "text"; text: string }
  | { type: "sdui"; envelope: SduiEnvelope };

function authToken(): string {
  return getToken() ?? exampleToken();
}

function friendlyError(status: number, body: string): string {
  try {
    const j = JSON.parse(body) as { error?: { message?: string; code?: string } };
    if (j.error?.message) return j.error.message;
  } catch {
    /* raw */
  }
  if (status === 401) return "Нужна авторизация.";
  if (status >= 500) return "Сервер временно недоступен.";
  return body.slice(0, 180) || "Ошибка запроса";
}

async function apiFetch<T>(path: string, init?: RequestInit): Promise<T> {
  const res = await fetch(`${API_URL}${path}`, {
    ...init,
    headers: {
      Authorization: `Bearer ${authToken()}`,
      "Content-Type": "application/json",
      ...(init?.headers ?? {}),
    },
    cache: "no-store",
  });
  if (!res.ok) {
    throw new Error(friendlyError(res.status, await res.text()));
  }
  return res.json() as Promise<T>;
}

export type Profile = {
  user_id: string;
  display_name: string;
  full_name?: string;
  phone?: string;
  email?: string;
  inn?: string;
  business_sphere: string;
  city: string;
  address?: string;
  cjm_level: number;
  monthly_revenue_estimate: number;
  tax_regime?: string;
  account_masked?: string;
  bank_name?: string;
  registered_at?: string;
  npd_status?: string;
  okved_codes?: string[];
};

export type AuthResponse = { token: string; profile: Profile };

export async function loginUser(phone: string): Promise<AuthResponse> {
  const res = await fetch(`${API_URL}/login`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ phone }),
  });
  if (!res.ok) throw new Error(friendlyError(res.status, await res.text()));
  return res.json() as Promise<AuthResponse>;
}

export async function registerAccount(body: {
  display_name: string;
  phone: string;
  email?: string;
  business_sphere: string;
  city: string;
  monthly_revenue_estimate: number;
  tax_regime: string;
}): Promise<AuthResponse> {
  const res = await fetch(`${API_URL}/register`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(body),
  });
  if (!res.ok) throw new Error(friendlyError(res.status, await res.text()));
  return res.json() as Promise<AuthResponse>;
}

export function getHome() {
  return apiFetch<{
    profile: Profile;
    tax_due: number;
    tax_estimate?: number;
    tax_rate_label?: string;
    period: string;
    piggy_gap?: number;
    tax_due_date?: string;
    tax_days_left?: number;
    empty_cabinet?: boolean;
    piggy: {
      enabled: boolean;
      rate_percent: number;
      balance: number;
      last_contribution: number;
      target_amount?: number;
    };
    insights: {
      id: string;
      text: string;
      why?: string;
      cta: string;
      href: string;
    }[];
  }>("/home");
}

export function getMe() {
  return apiFetch<Profile>("/me");
}

export function patchOnboarding(body: Partial<Profile> & { monthly_revenue_estimate?: number }) {
  return apiFetch<Profile>("/me/onboarding", {
    method: "PATCH",
    body: JSON.stringify(body),
  });
}

export function patchMe(body: Partial<Profile> & { monthly_revenue_estimate?: number }) {
  return apiFetch<Profile>("/me", {
    method: "PATCH",
    body: JSON.stringify(body),
  });
}

export function getTransactions() {
  return apiFetch<{
    items: {
      id: string;
      amount: number;
      direction: string;
      counterparty_name: string;
      counterparty_inn?: string;
      booked_at: string;
      category: string;
      comment?: string;
    }[];
    summary: { income: number; expense: number; net: number; currency: string };
  }>("/transactions");
}

export function getPiggy() {
  return apiFetch<{
    enabled: boolean;
    rate_percent: number;
    balance: number;
    last_contribution: number;
    currency: string;
    target_amount?: number;
  }>("/piggy");
}

export function patchPiggy(body: { enabled?: boolean; rate_percent?: number }) {
  return apiFetch<{
    enabled: boolean;
    rate_percent: number;
    balance: number;
    last_contribution: number;
    target_amount?: number;
  }>("/piggy", { method: "PATCH", body: JSON.stringify(body) });
}

export function getPiggyContributions() {
  return apiFetch<{
    items: { id: string; amount: number; booked_at: string; source_ref: string }[];
  }>("/piggy/contributions");
}

export function getAlerts() {
  return apiFetch<{
    items: {
      id: string;
      severity: string;
      title: string;
      text: string;
      cta: string;
      href: string;
    }[];
  }>("/alerts");
}

export function getTasks() {
  return apiFetch<{
    title: string;
    items: { id: string; label: string; text: string; href: string }[];
  }>("/tasks");
}

export function getPayments() {
  return apiFetch<{
    items: {
      draft_id: string;
      amount: number;
      purpose: string;
      payee_name: string;
      status: string;
      created_at: string;
      updated_at: string;
    }[];
  }>("/payments");
}

export function getTaxCalendar() {
  return apiFetch<{
    period_label: string;
    due_date: string;
    days_left: number;
    amount_due: number;
    amount_saved: number;
    amount_gap: number;
    status: string;
  }>("/calendar/tax");
}

export function getMonthReport() {
  return apiFetch<{
    period: string;
    income: number;
    expense: number;
    net: number;
    service_count: number;
    expense_by_category: Record<string, number>;
    tax_estimate: number;
    currency: string;
  }>("/report/month");
}

export function confirmDraft(id: string) {
  return apiFetch<{ draft_id: string; status: string }>(
    `/payments/draft/${id}/confirm`,
    { method: "POST", body: "{}" },
  );
}

export function cancelDraft(id: string) {
  return apiFetch<{ draft_id: string; status: string }>(
    `/payments/draft/${id}/cancel`,
    { method: "POST", body: "{}" },
  );
}

export async function uploadLegal(file: File): Promise<{
  document_id: string;
  document_name: string;
  component: string;
  mode?: string;
  narrative?: string;
  text_chars?: number;
  flags: {
    severity: string;
    title: string;
    quote: string;
    recommendation: string;
  }[];
}> {
  const form = new FormData();
  form.append("file", file);
  const res = await fetch(`${API_URL}/documents/legal`, {
    method: "POST",
    headers: { Authorization: `Bearer ${authToken()}` },
    body: form,
  });
  if (!res.ok) {
    throw new Error(friendlyError(res.status, await res.text()));
  }
  return res.json();
}

export async function streamChat(
  message: string,
  onEvent: (event: string, data: unknown) => void,
  history?: { role: string; content: string }[],
): Promise<void> {
  const res = await fetch(`${API_URL}/chat`, {
    method: "POST",
    headers: {
      Authorization: `Bearer ${authToken()}`,
      "Content-Type": "application/json",
    },
    body: JSON.stringify({ message, history: history ?? [] }),
  });
  if (!res.ok || !res.body) {
    throw new Error(friendlyError(res.status, await res.text().catch(() => "")));
  }
  const reader = res.body.getReader();
  const decoder = new TextDecoder();
  let buf = "";
  while (true) {
    const { done, value } = await reader.read();
    if (done) break;
    buf += decoder.decode(value, { stream: true });
    const chunks = buf.split("\n\n");
    buf = chunks.pop() ?? "";
    for (const chunk of chunks) {
      const lines = chunk.split("\n");
      let event = "message";
      let data = "";
      for (const line of lines) {
        if (line.startsWith("event:")) event = line.slice(6).trim();
        if (line.startsWith("data:")) data += line.slice(5).trim();
      }
      if (!data) continue;
      try {
        onEvent(event, JSON.parse(data));
      } catch {
        onEvent(event, data);
      }
    }
  }
}

export function formatRub(n: number) {
  return new Intl.NumberFormat("ru-RU").format(Math.round(n));
}
