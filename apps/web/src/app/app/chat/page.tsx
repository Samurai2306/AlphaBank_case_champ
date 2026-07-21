"use client";

import { Suspense, useEffect, useMemo, useRef, useState } from "react";
import Link from "next/link";
import { useRouter, useSearchParams } from "next/navigation";
import {
  ArrowLeft,
  ChevronDown,
  FileText,
  LayoutGrid,
  ListChecks,
  Mic,
  SendHorizontal,
  Sparkles,
} from "lucide-react";
import { GenUIRenderer } from "@/components/genui/GenUIRenderer";
import {
  cancelDraft,
  confirmDraft,
  patchOnboarding,
  streamChat,
  uploadLegal,
  type ChatPart,
  type SduiEnvelope,
} from "@/lib/api";
import { RichText } from "@/lib/richText";
import { QUICK_CHIPS, SCENARIOS } from "@/lib/scenarios";

type Suggestion = { label: string; text: string };

type Msg = {
  id: string;
  role: "user" | "assistant";
  parts: ChatPart[];
  suggestions?: Suggestion[];
};

function ChatInner() {
  const router = useRouter();
  const params = useSearchParams();
  const initialQ = params.get("q") ?? "";
  const [input, setInput] = useState("");
  const [busy, setBusy] = useState(false);
  const [listening, setListening] = useState(false);
  const [msgs, setMsgs] = useState<Msg[]>([]);
  const [menuOpen, setMenuOpen] = useState(false);
  const [liveSuggestions, setLiveSuggestions] = useState<Suggestion[]>([]);
  const boot = useRef(false);
  const endRef = useRef<HTMLDivElement>(null);
  const menuRef = useRef<HTMLDivElement>(null);
  const recognitionRef = useRef<{ stop: () => void } | null>(null);
  const fileRef = useRef<HTMLInputElement>(null);

  const disclaimer =
    "Альфа-Бизнес: Старт. Рекомендации справочные — финальное решение за вами.";

  useEffect(() => {
    endRef.current?.scrollIntoView({ behavior: "smooth" });
  }, [msgs, busy, liveSuggestions]);

  useEffect(() => {
    if (boot.current || !initialQ) return;
    boot.current = true;
    void send(initialQ);
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [initialQ]);

  useEffect(() => {
    if (!menuOpen) return;
    function onDoc(e: MouseEvent) {
      if (menuRef.current && !menuRef.current.contains(e.target as Node)) {
        setMenuOpen(false);
      }
    }
    function onKey(e: KeyboardEvent) {
      if (e.key === "Escape") setMenuOpen(false);
    }
    document.addEventListener("mousedown", onDoc);
    document.addEventListener("keydown", onKey);
    return () => {
      document.removeEventListener("mousedown", onDoc);
      document.removeEventListener("keydown", onKey);
    };
  }, [menuOpen]);

  function runSuggestion(s: Suggestion) {
    if (s.text.startsWith("__href__:")) {
      router.push(s.text.slice("__href__:".length));
      return;
    }
    void send(s.text);
  }

  function historyPayload(extraUser?: string) {
    const turns: { role: string; content: string }[] = [];
    for (const m of msgs) {
      const text = m.parts
        .filter((p): p is { type: "text"; text: string } => p.type === "text")
        .map((p) => p.text)
        .join("\n")
        .trim();
      if (!text) continue;
      turns.push({ role: m.role, content: text.slice(0, 800) });
    }
    if (extraUser) {
      turns.push({ role: "user", content: extraUser.slice(0, 800) });
    }
    return turns.slice(-8);
  }

  async function send(text: string) {
    const message = text.trim();
    if (!message || busy) return;
    setInput("");
    setMenuOpen(false);
    setLiveSuggestions([]);
    const history = historyPayload();
    const userMsg: Msg = {
      id: crypto.randomUUID(),
      role: "user",
      parts: [{ type: "text", text: message }],
    };
    const asstId = crypto.randomUUID();
    setMsgs((m) => [
      ...m,
      userMsg,
      { id: asstId, role: "assistant", parts: [] },
    ]);
    setBusy(true);
    let textAcc = "";
    const envelopes: SduiEnvelope[] = [];
    let suggestions: Suggestion[] = [];
    try {
      await streamChat(
        message,
        (event, data) => {
        if (event === "token") {
          const t = (data as { text?: string }).text ?? "";
          textAcc += t;
          setMsgs((all) =>
            all.map((msg) =>
              msg.id === asstId
                ? {
                    ...msg,
                    parts: [
                      { type: "text", text: textAcc },
                      ...envelopes.map((e) => ({
                        type: "sdui" as const,
                        envelope: e,
                      })),
                    ],
                    suggestions,
                  }
                : msg,
            ),
          );
        }
        if (event === "sdui") {
          const env = data as SduiEnvelope;
          envelopes.push(env);
          setMsgs((all) =>
            all.map((msg) =>
              msg.id === asstId
                ? {
                    ...msg,
                    parts: [
                      ...(textAcc
                        ? [{ type: "text" as const, text: textAcc }]
                        : []),
                      ...envelopes.map((e) => ({
                        type: "sdui" as const,
                        envelope: e,
                      })),
                    ],
                    suggestions,
                  }
                : msg,
            ),
          );
        }
        if (event === "suggestions") {
          const items = (data as { items?: Suggestion[] }).items ?? [];
          suggestions = items;
          setLiveSuggestions(items);
          setMsgs((all) =>
            all.map((msg) =>
              msg.id === asstId ? { ...msg, suggestions: items } : msg,
            ),
          );
        }
        if (event === "error") {
          const err = data as { message?: string };
          textAcc = err.message ?? "Не удалось обработать запрос. Попробуйте ещё раз.";
          setMsgs((all) =>
            all.map((msg) =>
              msg.id === asstId
                ? {
                    ...msg,
                    parts: [{ type: "text", text: textAcc }],
                    suggestions,
                  }
                : msg,
            ),
          );
        }
        },
        history,
      );
      if (!textAcc && envelopes.length === 0) {
        setMsgs((all) =>
          all.map((msg) =>
            msg.id === asstId
              ? {
                  ...msg,
                  parts: [
                    {
                      type: "text",
                      text: "Пустой ответ. Нажмите частый запрос ниже или сформулируйте иначе.",
                    },
                  ],
                }
              : msg,
          ),
        );
      }
    } catch (e) {
      setMsgs((all) =>
        all.map((msg) =>
          msg.id === asstId
            ? {
                ...msg,
                parts: [
                  {
                    type: "text",
                    text: `Не удалось связаться с API. ${(e as Error).message}`,
                  },
                ],
              }
            : msg,
        ),
      );
    } finally {
      setBusy(false);
    }
  }

  function patchDraftStatus(draftId: string, status: string) {
    setMsgs((all) =>
      all.map((msg) => ({
        ...msg,
        parts: msg.parts.map((p) => {
          if (p.type !== "sdui" || p.envelope.component !== "PaymentDraftCard") {
            return p;
          }
          if (String(p.envelope.props.draft_id ?? "") !== draftId) return p;
          return {
            type: "sdui" as const,
            envelope: {
              ...p.envelope,
              props: { ...p.envelope.props, status },
            },
          };
        }),
      })),
    );
  }

  function patchOnboardingConfirmed() {
    setMsgs((all) =>
      all.map((msg) => ({
        ...msg,
        parts: msg.parts.map((p) => {
          if (p.type !== "sdui" || p.envelope.component !== "OnboardingSummary") {
            return p;
          }
          return {
            type: "sdui" as const,
            envelope: {
              ...p.envelope,
              props: { ...p.envelope.props, confirmed: true },
            },
          };
        }),
      })),
    );
  }

  async function onAction(action: string, payload?: Record<string, unknown>) {
    if (action === "create_payment_draft") {
      await send("Сформируй платёжку.");
      return;
    }
    if (action === "confirm_onboarding") {
      try {
        await patchOnboarding({
          business_sphere: String(payload?.business_sphere ?? ""),
          city: String(payload?.city ?? ""),
          monthly_revenue_estimate: Number(payload?.monthly_revenue_estimate ?? 0),
        });
        patchOnboardingConfirmed();
        const followUps: Suggestion[] = [
          {
            label: "Посчитать налог",
            text: "Сколько мне отложить на налог за этот месяц?",
          },
          {
            label: "Безубыточность",
            text: "Посчитай, сколько клиентов в день нужно, если аренда кабинета 40к, расходники 200, цена услуги 1500.",
          },
        ];
        setLiveSuggestions(followUps);
        setMsgs((all) => [
          ...all,
          {
            id: crypto.randomUUID(),
            role: "assistant",
            parts: [
              {
                type: "text",
                text: "Профиль сохранён. Дальше логично посчитать налог НПД или точку безубыточности — выберите шаг ниже.",
              },
            ],
            suggestions: followUps,
          },
        ]);
      } catch (e) {
        setMsgs((all) => [
          ...all,
          {
            id: crypto.randomUUID(),
            role: "assistant",
            parts: [{ type: "text", text: (e as Error).message }],
          },
        ]);
      }
      return;
    }
    if (action === "confirm_draft") {
      const id = String(payload?.draft_id ?? "");
      if (!id) return;
      try {
        await confirmDraft(id);
        patchDraftStatus(id, "confirmed");
        const followUps: Suggestion[] = [
          { label: "К платежам", text: "__href__:/app/payments" },
          { label: "Включить копилку", text: "Включи копилку 6%." },
        ];
        setLiveSuggestions(followUps);
        setMsgs((all) => [
          ...all,
          {
            id: crypto.randomUUID(),
            role: "assistant",
            parts: [
              {
                type: "text",
                text: "Черновик подтверждён — реального списания со счёта пока нет. Историю смотрите в «Платежи»; параллельно можно включить копилку ЕНП.",
              },
            ],
            suggestions: followUps,
          },
        ]);
      } catch (e) {
        setMsgs((all) => [
          ...all,
          {
            id: crypto.randomUUID(),
            role: "assistant",
            parts: [{ type: "text", text: (e as Error).message }],
          },
        ]);
      }
      return;
    }
    if (action === "cancel_draft") {
      const id = String(payload?.draft_id ?? "");
      if (!id) return;
      try {
        await cancelDraft(id);
        patchDraftStatus(id, "cancelled");
      } catch (e) {
        setMsgs((all) => [
          ...all,
          {
            id: crypto.randomUUID(),
            role: "assistant",
            parts: [{ type: "text", text: (e as Error).message }],
          },
        ]);
      }
      return;
    }
    if (action === "send_message") {
      const text = String(payload?.text ?? "").trim();
      if (text) await send(text);
      return;
    }
    if (action === "navigate") {
      const href = String(payload?.href ?? "").trim();
      if (href.startsWith("/")) router.push(href);
    }
  }

  async function onLegalFile(file: File | undefined) {
    if (!file || busy) return;
    if (file.size > 10 * 1024 * 1024) {
      setMsgs((all) => [
        ...all,
        {
          id: crypto.randomUUID(),
          role: "assistant",
          parts: [{ type: "text", text: "Файл больше 10 МБ — загрузите меньший PDF." }],
        },
      ]);
      return;
    }
    setBusy(true);
    setLiveSuggestions([]);
    try {
      const res = await uploadLegal(file);
      const followUps: Suggestion[] = [
        { label: "Что с штрафом?", text: "А штраф за выход из договора нормальный?" },
        { label: "Проверить ИНН", text: "Проверь ИНН 1650987654" },
        {
          label: "Безубыточность",
          text: "Посчитай, сколько клиентов в день нужно, если аренда кабинета 40к, расходники 200, цена услуги 1500.",
        },
      ];
      setLiveSuggestions(followUps);
      const narrative =
        res.narrative?.trim() ||
        "Разобрала загруженный договор. Ниже — на что обратить внимание до подписания. Это предварительный разбор, не юридическое заключение.";
      const charsNote =
        typeof res.text_chars === "number" && res.text_chars > 0
          ? `\n\nИз PDF прочитано ≈${res.text_chars} символов текста.`
          : "\n\nТекст из PDF почти не извлечён — показала типовые риски аренды. Загрузите текстовый PDF для точечного разбора.";
      setMsgs((all) => [
        ...all,
        {
          id: crypto.randomUUID(),
          role: "user",
          parts: [{ type: "text", text: `Загружен договор: ${file.name}` }],
        },
        {
          id: crypto.randomUUID(),
          role: "assistant",
          parts: [
            {
              type: "text",
              text: narrative + charsNote,
            },
            {
              type: "sdui",
              envelope: {
                component: "LegalFlagsList",
                props: {
                  document_name: res.document_name,
                  flags: res.flags,
                  mode: res.mode,
                },
              },
            },
          ],
          suggestions: followUps,
        },
      ]);
    } catch (e) {
      setMsgs((all) => [
        ...all,
        {
          id: crypto.randomUUID(),
          role: "assistant",
          parts: [
            {
              type: "text",
              text: `Не удалось загрузить файл. ${(e as Error).message}`,
            },
          ],
        },
      ]);
    } finally {
      setBusy(false);
    }
  }

  const thinking = useMemo(() => busy, [busy]);

  function toggleVoice() {
    type Rec = {
      lang: string;
      continuous: boolean;
      interimResults: boolean;
      onresult:
        | ((ev: {
            resultIndex: number;
            results: ArrayLike<{
              isFinal?: boolean;
              0: { transcript: string };
            }>;
          }) => void)
        | null;
      onend: (() => void) | null;
      onerror: ((ev?: { error?: string }) => void) | null;
      start: () => void;
      stop: () => void;
    };
    const w = window as Window & {
      SpeechRecognition?: new () => Rec;
      webkitSpeechRecognition?: new () => Rec;
    };
    const Ctor = w.SpeechRecognition || w.webkitSpeechRecognition;
    if (!Ctor) {
      setMsgs((all) => [
        ...all,
        {
          id: crypto.randomUUID(),
          role: "assistant",
          parts: [
            {
              type: "text",
              text: "Голосовой ввод недоступен в этом браузере. Используйте Chrome или Edge и разрешите микрофон.",
            },
          ],
        },
      ]);
      return;
    }
    if (listening && recognitionRef.current) {
      recognitionRef.current.stop();
      setListening(false);
      return;
    }
    const rec = new Ctor();
    rec.lang = "ru-RU";
    rec.continuous = false;
    rec.interimResults = true;
    let finalText = "";
    rec.onresult = (ev) => {
      let interim = "";
      for (let i = ev.resultIndex; i < ev.results.length; i++) {
        const piece = ev.results[i]?.[0]?.transcript ?? "";
        if (ev.results[i]?.isFinal) {
          finalText = `${finalText} ${piece}`.trim();
        } else {
          interim += piece;
        }
      }
      const shown = `${finalText} ${interim}`.trim();
      if (shown) setInput(shown);
    };
    rec.onend = () => {
      setListening(false);
      recognitionRef.current = null;
      const text = finalText.trim();
      if (text.length >= 2 && !busy) {
        void send(text);
      }
    };
    rec.onerror = (ev) => {
      setListening(false);
      recognitionRef.current = null;
      if (ev?.error === "not-allowed") {
        setMsgs((all) => [
          ...all,
          {
            id: crypto.randomUUID(),
            role: "assistant",
            parts: [
              {
                type: "text",
                text: "Нет доступа к микрофону. Разрешите микрофон в браузере и нажмите на иконку снова.",
              },
            ],
          },
        ]);
      }
    };
    recognitionRef.current = rec;
    setListening(true);
    try {
      rec.start();
    } catch {
      setListening(false);
    }
  }

  const groups = useMemo(() => {
    const map = new Map<string, typeof SCENARIOS>();
    for (const s of SCENARIOS) {
      const list = map.get(s.group) ?? [];
      list.push(s);
      map.set(s.group, list);
    }
    return [...map.entries()];
  }, []);

  const chipsToShow =
    liveSuggestions.length > 0
      ? liveSuggestions
      : msgs.length === 0
        ? QUICK_CHIPS.map((c) => ({ label: c.label, text: c.text }))
        : [];

  return (
    <main className="mx-auto flex min-h-screen max-w-3xl flex-col px-4 pb-36 pt-4">
      <header className="mb-3 flex items-start justify-between gap-3">
        <div>
          <Link
            href="/app"
            className="inline-flex items-center gap-1.5 text-sm text-ink/50 hover:text-ink"
          >
            <ArrowLeft className="h-4 w-4" aria-hidden />
            Главная
          </Link>
          <h1 className="mt-1 flex items-center gap-2 text-xl font-bold">
            <Sparkles className="h-5 w-5 text-brand" aria-hidden />
            Альфа-Бизнес: Старт
          </h1>
          <p className="mt-0.5 text-xs text-ink/50">
            Налоги, платежи, договоры и подсказки по делу
          </p>
        </div>
        <div className="flex gap-2">
          <Link
            href="/app/services"
            className="inline-flex items-center gap-1 rounded-pill bg-white px-3 py-2 text-xs font-semibold text-ink/70 shadow-soft"
            title="Каталог сервисов"
          >
            <LayoutGrid className="h-3.5 w-3.5" aria-hidden />
            Сервисы
          </Link>
          <Link
            href="/app/tasks"
            className="inline-flex items-center gap-1 rounded-pill bg-brand px-3 py-2 text-xs font-semibold text-white"
            title="Рекомендуемые задачи"
          >
            <ListChecks className="h-3.5 w-3.5" aria-hidden />
            Задачи
          </Link>
        </div>
      </header>

      <div className="relative mb-3" ref={menuRef}>
        <button
          type="button"
          onClick={() => setMenuOpen((v) => !v)}
          className="btn-touch flex w-full items-center justify-between rounded-2xl bg-white px-4 py-3 text-left text-sm shadow-soft"
          aria-expanded={menuOpen}
        >
          <span className="flex items-center gap-2 font-medium">
            <Sparkles className="h-4 w-4 text-brand" aria-hidden />
            Частые запросы
          </span>
          <ChevronDown
            className={`h-4 w-4 text-ink/50 transition ${menuOpen ? "rotate-180" : ""}`}
            aria-hidden
          />
        </button>
        {menuOpen ? (
          <div className="absolute left-0 right-0 z-30 mt-2 max-h-72 overflow-auto rounded-2xl bg-white p-2 shadow-[0_16px_40px_rgba(11,17,23,0.14)] ring-1 ring-platinum/60">
            {groups.map(([group, items]) => (
              <div key={group} className="mb-2">
                <p className="px-2 py-1 text-[11px] font-semibold uppercase tracking-wide text-ink/40">
                  {group}
                </p>
                {items.map((s) => (
                  <button
                    key={s.id}
                    type="button"
                    disabled={busy}
                    onClick={() => void send(s.text)}
                    className="flex w-full flex-col rounded-xl px-3 py-2.5 text-left hover:bg-canvas"
                  >
                    <span className="text-sm font-semibold">{s.label}</span>
                    <span className="text-xs text-ink/50">{s.hint}</span>
                  </button>
                ))}
              </div>
            ))}
          </div>
        ) : null}
      </div>

      <div className="flex-1 space-y-4">
        {msgs.length === 0 ? (
          <div className="rounded-[28px] bg-white p-5 shadow-soft">
            <p className="font-semibold">С чего начать</p>
            <p className="mt-1 text-sm text-ink/60">
              Выберите частый запрос или подсказку — ответ придёт карточкой с
              кнопками. Можно писать своими словами.
            </p>
            <ul className="mt-4 space-y-2 text-sm text-ink/70">
              <li>1. Онбординг → налог → платёжка</li>
              <li>2. Безубыточность → проверка ИНН</li>
              <li>3. Копилка ЕНП → платёжка к сроку</li>
            </ul>
          </div>
        ) : null}

        {msgs.map((m) => (
          <div
            key={m.id}
            className={`flex ${m.role === "user" ? "justify-end" : "justify-start"}`}
          >
            <div
              className={`max-w-[92%] space-y-3 ${
                m.role === "user"
                  ? "rounded-card bg-ink px-4 py-3 text-white"
                  : ""
              }`}
            >
              {m.parts.map((p, idx) =>
                p.type === "text" ? (
                  <p
                    key={idx}
                    className={`whitespace-pre-wrap text-sm leading-relaxed ${
                      m.role === "assistant"
                        ? "rounded-card bg-white px-4 py-3 shadow-soft"
                        : ""
                    }`}
                  >
                    <RichText text={p.text} />
                  </p>
                ) : (
                  <GenUIRenderer
                    key={idx}
                    envelope={p.envelope}
                    onAction={onAction}
                  />
                ),
              )}
            </div>
          </div>
        ))}

        {thinking ? (
          <div className="inline-flex items-center gap-2 rounded-card bg-white px-4 py-3 text-sm text-ink/60 shadow-soft">
            <span className="flex gap-1">
              <i className="h-2 w-2 animate-bounce rounded-full bg-brand [animation-delay:-0.2s]" />
              <i className="h-2 w-2 animate-bounce rounded-full bg-brand [animation-delay:-0.1s]" />
              <i className="h-2 w-2 animate-bounce rounded-full bg-brand" />
            </span>
            Считаю…
          </div>
        ) : null}

        {chipsToShow.length > 0 && !busy ? (
          <div className="flex flex-wrap gap-2" aria-label="Подсказки">
            {chipsToShow.map((c) => (
              <button
                key={c.label + c.text}
                type="button"
                onClick={() => runSuggestion(c)}
                className="btn-touch inline-flex items-center rounded-pill bg-white px-3.5 py-2 text-xs font-semibold text-ink shadow-soft ring-1 ring-platinum/70 hover:ring-brand/40"
              >
                {c.label}
              </button>
            ))}
          </div>
        ) : null}

        <div ref={endRef} />
      </div>

      <div className="fixed bottom-0 left-0 right-0 border-t border-platinum/60 bg-canvas/95 backdrop-blur">
        <p className="mx-auto max-w-3xl px-4 pt-2 text-[11px] leading-snug text-ink/45">
          {disclaimer}
        </p>
        <form
          className="mx-auto flex max-w-3xl items-center gap-2 px-4 py-3"
          onSubmit={(e) => {
            e.preventDefault();
            void send(input);
          }}
        >
          <input
            ref={fileRef}
            type="file"
            accept="application/pdf,.pdf"
            className="hidden"
            onChange={(e) => {
              void onLegalFile(e.target.files?.[0]);
              e.target.value = "";
            }}
          />
          <button
            type="button"
            title="Загрузить договор PDF"
            onClick={() => fileRef.current?.click()}
            className="btn-touch inline-flex h-11 w-11 shrink-0 items-center justify-center rounded-full border border-platinum bg-white text-ink/70 hover:border-brand hover:text-brand"
          >
            <FileText className="h-4 w-4" aria-hidden />
            <span className="sr-only">Документ PDF</span>
          </button>
          <button
            type="button"
            onClick={toggleVoice}
            aria-pressed={listening}
            title="Голосовой ввод"
            className={`btn-touch inline-flex h-11 w-11 shrink-0 items-center justify-center rounded-full border ${
              listening
                ? "border-brand bg-rose text-brand"
                : "border-platinum bg-white text-ink/70 hover:border-brand"
            }`}
          >
            <Mic className="h-4 w-4" aria-hidden />
            <span className="sr-only">Микрофон</span>
          </button>
          <input
            value={input}
            onChange={(e) => setInput(e.target.value)}
            placeholder="Спросите про налог, платёжку, ИНН…"
            className="min-w-0 flex-1 rounded-pill border border-platinum bg-white px-4 py-3 text-sm outline-none ring-brand focus:ring-2"
          />
          <button
            type="submit"
            disabled={busy || !input.trim()}
            title="Отправить"
            className="btn-touch inline-flex h-11 w-11 shrink-0 items-center justify-center rounded-full bg-brand text-white disabled:opacity-50"
          >
            <SendHorizontal className="h-4 w-4" aria-hidden />
            <span className="sr-only">Отправить</span>
          </button>
        </form>
      </div>
    </main>
  );
}

export default function ChatPage() {
  return (
    <Suspense fallback={<main className="p-6">Загрузка чата…</main>}>
      <ChatInner />
    </Suspense>
  );
}
