"use client";

import Link from "next/link";
import { useEffect, useState, type ReactNode } from "react";

const API_URL = (process.env.NEXT_PUBLIC_API_URL || "/api/v1").replace(/\/$/, "");
const TOKEN_KEY = "alfa_admin_token";

type Profile = {
  user_id: string;
  display_name: string;
  full_name?: string;
  phone?: string;
  email?: string;
  city?: string;
  business_sphere?: string;
  tax_regime?: string;
  monthly_revenue_estimate?: number;
  cjm_level?: number;
  registered_at?: string;
};

type Draft = {
  draft_id?: string;
  id?: string;
  amount: number;
  purpose: string;
  status: string;
  payee_name?: string;
};

type ChatLine = { role: string; content: string };

type UserCard = {
  profile: Profile;
  txn_count: number;
  draft_count: number;
  chat_count: number;
  piggy_balance: number;
  piggy_on: boolean;
};

type UserDetail = UserCard & { drafts: Draft[]; chat: ChatLine[] };

type ProfileForm = {
  display_name: string;
  full_name: string;
  phone: string;
  email: string;
  city: string;
  business_sphere: string;
  tax_regime: string;
  monthly_revenue_estimate: string;
  cjm_level: string;
};

const emptyCreate = {
  display_name: "",
  phone: "",
  city: "",
  business_sphere: "",
  tax_regime: "NPD",
  monthly_revenue_estimate: "",
};

function rub(n: number) {
  return new Intl.NumberFormat("ru-RU").format(Math.round(n)) + " ₽";
}

function draftKey(d: Draft) {
  return d.draft_id || d.id || d.purpose;
}

function apiMessage(text: string, fallback: string) {
  try {
    const data = JSON.parse(text) as { error?: { message?: string } };
    if (data.error?.message) return data.error.message;
  } catch {
    /* ответ не JSON */
  }
  return fallback;
}

async function adminFetch<T>(path: string, token: string, init?: RequestInit): Promise<T> {
  const res = await fetch(`${API_URL}${path}`, {
    ...init,
    headers: {
      Authorization: `Bearer ${token}`,
      "Content-Type": "application/json",
      ...(init?.headers ?? {}),
    },
    cache: "no-store",
  });
  const text = await res.text();
  if (!res.ok) {
    const err = new Error(apiMessage(text, "Ошибка запроса")) as Error & { status?: number };
    err.status = res.status;
    throw err;
  }
  if (!text) return {} as T;
  return JSON.parse(text) as T;
}

function profileForm(profile: Profile): ProfileForm {
  return {
    display_name: profile.display_name || "",
    full_name: profile.full_name || "",
    phone: profile.phone || "",
    email: profile.email || "",
    city: profile.city || "",
    business_sphere: profile.business_sphere || "",
    tax_regime: profile.tax_regime || "NPD",
    monthly_revenue_estimate: String(profile.monthly_revenue_estimate ?? ""),
    cjm_level: String(profile.cjm_level ?? 1),
  };
}

export default function AdminPage() {
  const [token, setToken] = useState<string | null>(null);
  const [ready, setReady] = useState(false);
  const [password, setPassword] = useState("");
  const [err, setErr] = useState("");
  const [busy, setBusy] = useState(false);
  const [items, setItems] = useState<UserCard[]>([]);
  const [selected, setSelected] = useState<UserDetail | null>(null);
  const [creating, setCreating] = useState(false);
  const [createForm, setCreateForm] = useState(emptyCreate);
  const [profile, setProfile] = useState<ProfileForm | null>(null);
  const [draftAmount, setDraftAmount] = useState("");
  const [draftPurpose, setDraftPurpose] = useState("");
  const [draftPayee, setDraftPayee] = useState("");

  useEffect(() => {
    const saved = localStorage.getItem(TOKEN_KEY);
    setToken(saved);
    setReady(true);
  }, []);

  useEffect(() => {
    if (!token) return;
    let stop = false;
    async function load() {
      try {
        const data = await adminFetch<{ items: UserCard[] }>("/admin/users", token!);
        if (!stop) setItems(data.items ?? []);
      } catch (e) {
        if (stop) return;
        const status = (e as { status?: number }).status;
        if (status === 401) {
          localStorage.removeItem(TOKEN_KEY);
          setToken(null);
          setErr("Сессия админки сброшена. Войдите снова.");
          return;
        }
        setErr("Не удалось обновить список");
      }
    }
    void load();
    const timer = window.setInterval(() => void load(), 4000);
    return () => {
      stop = true;
      window.clearInterval(timer);
    };
  }, [token]);

  useEffect(() => {
    if (!selected) {
      setProfile(null);
      return;
    }
    setProfile(profileForm(selected.profile));
  }, [selected]);

  async function enter() {
    setBusy(true);
    setErr("");
    try {
      const res = await fetch(`${API_URL}/admin/login`, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ password }),
      });
      if (!res.ok) {
        setErr("Неверный пароль");
        setBusy(false);
        return;
      }
      const data = (await res.json()) as { token: string };
      localStorage.setItem(TOKEN_KEY, data.token);
      setToken(data.token);
      setPassword("");
    } catch {
      setErr("Сервер недоступен");
    } finally {
      setBusy(false);
    }
  }

  async function openUser(id: string) {
    if (!token) return;
    setErr("");
    try {
      const detail = await adminFetch<UserDetail>(`/admin/users/${id}`, token);
      setSelected(detail);
    } catch (e) {
      setErr(e instanceof Error ? e.message : "Не удалось открыть кабинет");
    }
  }

  async function createUser() {
    if (!token) return;
    setBusy(true);
    setErr("");
    try {
      const revenue = Number(createForm.monthly_revenue_estimate.replace(",", "."));
      const detail = await adminFetch<UserDetail>("/admin/users", token, {
        method: "POST",
        body: JSON.stringify({
          display_name: createForm.display_name.trim(),
          phone: createForm.phone.trim(),
          city: createForm.city.trim(),
          business_sphere: createForm.business_sphere.trim(),
          tax_regime: createForm.tax_regime,
          monthly_revenue_estimate: Number.isFinite(revenue) ? revenue : 0,
        }),
      });
      setCreating(false);
      setCreateForm(emptyCreate);
      setSelected(detail);
      setItems((prev) => {
        const rest = prev.filter((row) => row.profile.user_id !== detail.profile.user_id);
        return [detail, ...rest];
      });
    } catch (e) {
      setErr(e instanceof Error ? e.message : "Не удалось создать кабинет");
    } finally {
      setBusy(false);
    }
  }

  async function saveProfile() {
    if (!token || !selected || !profile) return;
    setBusy(true);
    setErr("");
    try {
      const revenue = Number(profile.monthly_revenue_estimate.replace(",", "."));
      const detail = await adminFetch<UserDetail>(`/admin/users/${selected.profile.user_id}`, token, {
        method: "PATCH",
        body: JSON.stringify({
          display_name: profile.display_name.trim(),
          full_name: profile.full_name.trim(),
          phone: profile.phone.trim(),
          email: profile.email.trim(),
          city: profile.city.trim(),
          business_sphere: profile.business_sphere.trim(),
          tax_regime: profile.tax_regime,
          monthly_revenue_estimate: Number.isFinite(revenue) ? revenue : 0,
          cjm_level: Number(profile.cjm_level),
        }),
      });
      setSelected(detail);
    } catch (e) {
      setErr(e instanceof Error ? e.message : "Не удалось сохранить профиль");
    } finally {
      setBusy(false);
    }
  }

  async function removeUser() {
    if (!token || !selected) return;
    const name = selected.profile.display_name;
    if (!window.confirm(`Удалить кабинет «${name}»? Вход этого клиента перестанет работать.`)) return;
    setBusy(true);
    setErr("");
    try {
      await adminFetch(`/admin/users/${selected.profile.user_id}`, token, { method: "DELETE" });
      setItems((prev) => prev.filter((row) => row.profile.user_id !== selected.profile.user_id));
      setSelected(null);
    } catch (e) {
      setErr(e instanceof Error ? e.message : "Не удалось удалить кабинет");
    } finally {
      setBusy(false);
    }
  }

  async function togglePiggy() {
    if (!token || !selected) return;
    setBusy(true);
    setErr("");
    try {
      const detail = await adminFetch<UserDetail>(`/admin/users/${selected.profile.user_id}/piggy`, token, {
        method: "PATCH",
        body: JSON.stringify({ enabled: !selected.piggy_on }),
      });
      setSelected(detail);
    } catch (e) {
      setErr(e instanceof Error ? e.message : "Не удалось изменить копилку");
    } finally {
      setBusy(false);
    }
  }

  async function addDraft() {
    if (!token || !selected) return;
    const amount = Number(draftAmount.replace(",", "."));
    setBusy(true);
    setErr("");
    try {
      const detail = await adminFetch<UserDetail>(`/admin/users/${selected.profile.user_id}/drafts`, token, {
        method: "POST",
        body: JSON.stringify({
          amount,
          purpose: draftPurpose.trim(),
          payee_name: draftPayee.trim(),
        }),
      });
      setSelected(detail);
      setDraftAmount("");
      setDraftPurpose("");
      setDraftPayee("");
    } catch (e) {
      setErr(e instanceof Error ? e.message : "Не удалось создать черновик");
    } finally {
      setBusy(false);
    }
  }

  async function setDraftStatus(draft: Draft, action: "confirm" | "cancel") {
    if (!token || !selected) return;
    const id = draft.draft_id || draft.id;
    if (!id) return;
    if (action === "confirm") {
      const ok = window.confirm(
        `Подтвердить черновик на ${rub(draft.amount)}? На стенде меняется только статус, деньги не списываются.`,
      );
      if (!ok) return;
    }
    setBusy(true);
    setErr("");
    try {
      const detail = await adminFetch<UserDetail>(
        `/admin/users/${selected.profile.user_id}/drafts/${id}/${action}`,
        token,
        { method: "POST", body: "{}" },
      );
      setSelected(detail);
    } catch (e) {
      setErr(e instanceof Error ? e.message : "Не удалось обновить черновик");
    } finally {
      setBusy(false);
    }
  }

  function leave() {
    localStorage.removeItem(TOKEN_KEY);
    setToken(null);
    setItems([]);
    setSelected(null);
    setCreating(false);
  }

  if (!ready) {
    return (
      <main className="flex min-h-screen items-center justify-center text-sm text-ink/50">
        Загрузка…
      </main>
    );
  }

  if (!token) {
    return (
      <main className="mx-auto flex min-h-screen max-w-md flex-col justify-center px-4 py-10">
        <Link href="/login" className="text-sm text-ink/50 hover:text-ink">
          ← Ко входу
        </Link>
        <p className="mt-8 text-sm font-semibold text-brand">Альфа-Бизнес: Старт</p>
        <h1 className="mt-2 text-3xl font-bold">Админка стенда</h1>
        <p className="mt-2 text-sm text-ink/60">
          Кабинеты, открытые на стенде. Демо-пароль: start-admin.
        </p>
        <label className="mt-6 block text-sm">
          <span className="text-ink/60">Пароль</span>
          <input
            type="password"
            value={password}
            onChange={(e) => setPassword(e.target.value)}
            onKeyDown={(e) => {
              if (e.key === "Enter") void enter();
            }}
            className="mt-1 w-full rounded-2xl border border-platinum bg-white px-4 py-3 text-ink outline-none focus:ring-2 focus:ring-brand/30"
          />
        </label>
        {err ? <p className="mt-3 text-sm text-brand">{err}</p> : null}
        <button
          type="button"
          disabled={busy || !password.trim()}
          onClick={() => void enter()}
          className="mt-6 w-full rounded-pill bg-brand py-3.5 text-sm font-semibold text-white disabled:opacity-50"
        >
          {busy ? "Входим…" : "Войти"}
        </button>
      </main>
    );
  }

  return (
    <main className="mx-auto min-h-screen max-w-6xl px-4 py-8">
      <div className="flex flex-wrap items-end justify-between gap-3">
        <div>
          <p className="text-sm font-semibold text-brand">Альфа-Бизнес: Старт</p>
          <h1 className="mt-1 text-3xl font-bold">Кабинеты</h1>
        </div>
        <div className="flex gap-2">
          <button
            type="button"
            onClick={() => setCreating((v) => !v)}
            className="rounded-pill bg-brand px-4 py-2 text-sm font-semibold text-white"
          >
            {creating ? "Закрыть форму" : "Новый кабинет"}
          </button>
          <button
            type="button"
            onClick={leave}
            className="rounded-pill bg-white px-4 py-2 text-sm font-semibold text-ink shadow-soft ring-1 ring-platinum"
          >
            Выйти
          </button>
        </div>
      </div>

      {err ? <p className="mt-4 text-sm text-brand">{err}</p> : null}

      {creating ? (
        <form
          className="mt-4 grid gap-3 rounded-[28px] bg-white p-5 shadow-soft ring-1 ring-platinum sm:grid-cols-2"
          onSubmit={(e) => {
            e.preventDefault();
            void createUser();
          }}
        >
          <Field label="Имя">
            <TextInput
              value={createForm.display_name}
              onChange={(display_name) => setCreateForm({ ...createForm, display_name })}
              required
            />
          </Field>
          <Field label="Телефон">
            <TextInput
              value={createForm.phone}
              onChange={(phone) => setCreateForm({ ...createForm, phone })}
              required
            />
          </Field>
          <Field label="Город">
            <TextInput value={createForm.city} onChange={(city) => setCreateForm({ ...createForm, city })} />
          </Field>
          <Field label="Сфера">
            <TextInput
              value={createForm.business_sphere}
              onChange={(business_sphere) => setCreateForm({ ...createForm, business_sphere })}
            />
          </Field>
          <Field label="Режим">
            <RegimeSelect
              value={createForm.tax_regime}
              onChange={(tax_regime) => setCreateForm({ ...createForm, tax_regime })}
            />
          </Field>
          <Field label="Оборот в месяц, ₽">
            <TextInput
              value={createForm.monthly_revenue_estimate}
              onChange={(monthly_revenue_estimate) => setCreateForm({ ...createForm, monthly_revenue_estimate })}
              inputMode="decimal"
            />
          </Field>
          <div className="sm:col-span-2">
            <button
              type="submit"
              disabled={busy || !createForm.display_name.trim() || !createForm.phone.trim()}
              className="rounded-pill bg-brand px-5 py-2.5 text-sm font-semibold text-white disabled:opacity-50"
            >
              {busy ? "Сохраняем…" : "Создать кабинет"}
            </button>
          </div>
        </form>
      ) : null}

      <div className="mt-6 grid gap-4 lg:grid-cols-[1fr_1fr]">
        <section className="overflow-x-auto rounded-[28px] bg-white shadow-soft ring-1 ring-platinum">
          <table className="w-full text-left text-sm">
            <thead className="bg-[#F2F3F5] text-xs uppercase tracking-wide text-ink/50">
              <tr>
                <th className="px-4 py-3 font-semibold">Клиент</th>
                <th className="px-4 py-3 font-semibold">Город</th>
                <th className="px-4 py-3 font-semibold">Оборот</th>
                <th className="px-4 py-3 font-semibold">Копилка</th>
              </tr>
            </thead>
            <tbody>
              {items.length === 0 ? (
                <tr>
                  <td colSpan={4} className="px-4 py-6 text-sm text-ink/50">
                    Кабинетов пока нет
                  </td>
                </tr>
              ) : null}
              {items.map((row) => {
                const active = selected?.profile.user_id === row.profile.user_id;
                return (
                  <tr key={row.profile.user_id} className={active ? "bg-rose/40" : "hover:bg-[#FAFAFB]"}>
                    <td className="px-4 py-3">
                      <button
                        type="button"
                        onClick={() => void openUser(row.profile.user_id)}
                        className="text-left"
                      >
                        <span className="font-semibold">{row.profile.display_name}</span>
                        <span className="mt-0.5 block text-xs text-ink/50">
                          {row.profile.phone || "без телефона"} · уровень {row.profile.cjm_level ?? 0}
                        </span>
                      </button>
                    </td>
                    <td className="px-4 py-3 text-ink/70">{row.profile.city || "—"}</td>
                    <td className="px-4 py-3 tabular-nums">{rub(row.profile.monthly_revenue_estimate ?? 0)}</td>
                    <td className="px-4 py-3 tabular-nums">{rub(row.piggy_balance)}</td>
                  </tr>
                );
              })}
            </tbody>
          </table>
        </section>

        <aside className="rounded-[28px] bg-white p-5 shadow-soft ring-1 ring-platinum">
          {selected && profile ? (
            <div>
              <p className="text-xs uppercase tracking-wide text-ink/45">Кабинет</p>
              <h2 className="mt-1 text-2xl font-bold">{selected.profile.display_name}</h2>
              <p className="mt-1 text-sm text-ink/55">
                Регистрация {selected.profile.registered_at || "—"} · операций {selected.txn_count}
              </p>

              <form
                className="mt-4 grid gap-3 sm:grid-cols-2"
                onSubmit={(e) => {
                  e.preventDefault();
                  void saveProfile();
                }}
              >
                <Field label="Имя">
                  <TextInput
                    value={profile.display_name}
                    onChange={(display_name) => setProfile({ ...profile, display_name })}
                    required
                  />
                </Field>
                <Field label="ФИО">
                  <TextInput
                    value={profile.full_name}
                    onChange={(full_name) => setProfile({ ...profile, full_name })}
                  />
                </Field>
                <Field label="Телефон">
                  <TextInput value={profile.phone} onChange={(phone) => setProfile({ ...profile, phone })} required />
                </Field>
                <Field label="Почта">
                  <TextInput value={profile.email} onChange={(email) => setProfile({ ...profile, email })} />
                </Field>
                <Field label="Город">
                  <TextInput value={profile.city} onChange={(city) => setProfile({ ...profile, city })} />
                </Field>
                <Field label="Сфера">
                  <TextInput
                    value={profile.business_sphere}
                    onChange={(business_sphere) => setProfile({ ...profile, business_sphere })}
                  />
                </Field>
                <Field label="Режим">
                  <RegimeSelect
                    value={profile.tax_regime}
                    onChange={(tax_regime) => setProfile({ ...profile, tax_regime })}
                  />
                </Field>
                <Field label="Уровень CJM">
                  <select
                    value={profile.cjm_level}
                    onChange={(e) => setProfile({ ...profile, cjm_level: e.target.value })}
                    className="mt-1 w-full rounded-2xl border border-platinum bg-white px-4 py-3 text-ink outline-none focus:ring-2 focus:ring-brand/30"
                  >
                    {[1, 2, 3, 4, 5].map((level) => (
                      <option key={level} value={level}>
                        {level}
                      </option>
                    ))}
                  </select>
                </Field>
                <Field label="Оборот в месяц, ₽">
                  <TextInput
                    value={profile.monthly_revenue_estimate}
                    onChange={(monthly_revenue_estimate) => setProfile({ ...profile, monthly_revenue_estimate })}
                    inputMode="decimal"
                  />
                </Field>
                <div className="flex items-end">
                  <button
                    type="submit"
                    disabled={busy}
                    className="w-full rounded-pill bg-brand py-3 text-sm font-semibold text-white disabled:opacity-50"
                  >
                    Сохранить профиль
                  </button>
                </div>
              </form>

              <div className="mt-4 flex flex-wrap items-center justify-between gap-2 rounded-2xl bg-[#F2F3F5] px-3 py-3 text-sm">
                <span>
                  Копилка {selected.piggy_on ? "включена" : "выключена"} · {rub(selected.piggy_balance)}
                </span>
                <button
                  type="button"
                  disabled={busy}
                  onClick={() => void togglePiggy()}
                  className="rounded-pill bg-white px-3 py-1.5 text-sm font-semibold text-ink ring-1 ring-platinum disabled:opacity-50"
                >
                  {selected.piggy_on ? "Выключить" : "Включить"}
                </button>
              </div>

              <h3 className="mt-5 text-sm font-semibold">Платёжки</h3>
              <p className="mt-1 text-xs text-ink/50">
                Подтверждение меняет статус черновика. Деньги со счёта не списываются.
              </p>
              {selected.drafts?.length ? (
                <ul className="mt-2 space-y-2">
                  {selected.drafts.map((d) => (
                    <li key={draftKey(d)} className="rounded-2xl bg-[#F2F3F5] px-3 py-2 text-sm">
                      <div className="flex flex-wrap items-start justify-between gap-2">
                        <div>
                          <span className="font-semibold">{rub(d.amount)}</span>
                          <span className="text-ink/50"> · {d.status}</span>
                          <p className="text-ink/70">{d.purpose}</p>
                          {d.payee_name ? <p className="text-xs text-ink/45">{d.payee_name}</p> : null}
                        </div>
                        {d.status === "draft" ? (
                          <div className="flex gap-2">
                            <button
                              type="button"
                              disabled={busy}
                              onClick={() => void setDraftStatus(d, "confirm")}
                              className="rounded-pill bg-brand px-3 py-1.5 text-xs font-semibold text-white disabled:opacity-50"
                            >
                              Подтвердить
                            </button>
                            <button
                              type="button"
                              disabled={busy}
                              onClick={() => void setDraftStatus(d, "cancel")}
                              className="rounded-pill bg-white px-3 py-1.5 text-xs font-semibold text-ink ring-1 ring-platinum disabled:opacity-50"
                            >
                              Отменить
                            </button>
                          </div>
                        ) : null}
                      </div>
                    </li>
                  ))}
                </ul>
              ) : (
                <p className="mt-2 text-sm text-ink/45">Черновиков нет</p>
              )}

              <form
                className="mt-3 grid gap-2"
                onSubmit={(e) => {
                  e.preventDefault();
                  void addDraft();
                }}
              >
                <TextInput
                  value={draftAmount}
                  onChange={setDraftAmount}
                  inputMode="decimal"
                  placeholder="Сумма, ₽"
                  aria-label="Сумма черновика"
                />
                <TextInput
                  value={draftPurpose}
                  onChange={setDraftPurpose}
                  placeholder="Назначение"
                  aria-label="Назначение платежа"
                />
                <TextInput
                  value={draftPayee}
                  onChange={setDraftPayee}
                  placeholder="Получатель"
                  aria-label="Получатель"
                />
                <button
                  type="submit"
                  disabled={busy || !draftAmount.trim()}
                  className="rounded-pill bg-white py-2.5 text-sm font-semibold text-ink ring-1 ring-platinum disabled:opacity-50"
                >
                  Создать черновик
                </button>
              </form>

              <h3 className="mt-5 text-sm font-semibold">Чат</h3>
              {selected.chat?.length ? (
                <ul className="mt-2 space-y-2">
                  {selected.chat.slice(-6).map((line, i) => (
                    <li key={i} className="text-sm text-ink/70">
                      <span className="font-semibold text-ink">{line.role}: </span>
                      {line.content}
                    </li>
                  ))}
                </ul>
              ) : (
                <p className="mt-2 text-sm text-ink/45">Сообщений ещё нет</p>
              )}

              <button
                type="button"
                disabled={busy}
                onClick={() => void removeUser()}
                className="mt-5 text-sm font-semibold text-brand disabled:opacity-50"
              >
                Удалить кабинет
              </button>
            </div>
          ) : (
            <p className="text-sm text-ink/55">Выберите кабинет в списке слева.</p>
          )}
        </aside>
      </div>
    </main>
  );
}

function Field({ label, children }: { label: string; children: ReactNode }) {
  return (
    <label className="block text-sm">
      <span className="text-ink/60">{label}</span>
      {children}
    </label>
  );
}

function TextInput({
  value,
  onChange,
  required,
  inputMode,
  placeholder,
  "aria-label": ariaLabel,
}: {
  value: string;
  onChange: (value: string) => void;
  required?: boolean;
  inputMode?: "decimal" | "text";
  placeholder?: string;
  "aria-label"?: string;
}) {
  return (
    <input
      value={value}
      required={required}
      inputMode={inputMode}
      placeholder={placeholder}
      aria-label={ariaLabel}
      onChange={(e) => onChange(e.target.value)}
      className="mt-1 w-full rounded-2xl border border-platinum bg-white px-4 py-3 text-ink outline-none focus:ring-2 focus:ring-brand/30"
    />
  );
}

function RegimeSelect({ value, onChange }: { value: string; onChange: (value: string) => void }) {
  return (
    <select
      value={value}
      onChange={(e) => onChange(e.target.value)}
      className="mt-1 w-full rounded-2xl border border-platinum bg-white px-4 py-3 text-ink outline-none focus:ring-2 focus:ring-brand/30"
    >
      <option value="NPD">НПД</option>
      <option value="IP">ИП</option>
      <option value="UNDECIDED">Ещё не выбран</option>
    </select>
  );
}
