const TOKEN_KEY = "alfa_start_token";
const USER_KEY = "alfa_start_user_id";
const LEGACY_TOKEN_KEY = "copilot_token";
const LEGACY_USER_KEY = "copilot_user_id";
const EXAMPLE_TOKEN =
  process.env.NEXT_PUBLIC_DEMO_TOKEN ?? "demo-masha-token";

export function getToken(): string | null {
  if (typeof window === "undefined") return null;
  return (
    localStorage.getItem(TOKEN_KEY) || localStorage.getItem(LEGACY_TOKEN_KEY)
  );
}

export function getUserId(): string | null {
  if (typeof window === "undefined") return null;
  return (
    localStorage.getItem(USER_KEY) || localStorage.getItem(LEGACY_USER_KEY)
  );
}

export function hasSession(): boolean {
  return !!getToken();
}

export function setSession(token: string, userId?: string) {
  localStorage.setItem(TOKEN_KEY, token);
  localStorage.removeItem(LEGACY_TOKEN_KEY);
  if (userId) {
    localStorage.setItem(USER_KEY, userId);
    localStorage.removeItem(LEGACY_USER_KEY);
  }
}

export function ensureUserId(userId: string) {
  if (userId) {
    localStorage.setItem(USER_KEY, userId);
    localStorage.removeItem(LEGACY_USER_KEY);
  }
}

export function clearSession() {
  localStorage.removeItem(TOKEN_KEY);
  localStorage.removeItem(USER_KEY);
  localStorage.removeItem(LEGACY_TOKEN_KEY);
  localStorage.removeItem(LEGACY_USER_KEY);
  localStorage.removeItem("alfa_start_session");
  localStorage.removeItem("copilot_session");
}

/** Seeded example cabinet for «Посмотреть пример». */
export function enterExampleSession() {
  setSession(EXAMPLE_TOKEN, "usr_masha_kazan");
}

export function exampleToken() {
  return EXAMPLE_TOKEN;
}
