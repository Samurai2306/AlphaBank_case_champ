import { chromium } from "playwright";
import { mkdirSync } from "fs";
import { join } from "path";

const BASE = process.env.PLAYWRIGHT_BASE_URL ?? "http://localhost:3000";
const out = join(process.cwd(), "audit-shots");
mkdirSync(out, { recursive: true });

const pages = [
  ["01-landing", "/"],
  ["02-login", "/login"],
  ["03-tasks", "/app/tasks"],
  ["04-home", "/app"],
  ["05-operations", "/app/operations"],
  ["06-services", "/app/services"],
  ["07-piggy", "/app/piggy"],
  ["08-profile", "/app/profile"],
  ["09-alerts", "/app/alerts"],
  ["10-journey", "/app/journey"],
  ["11-chat", "/app/chat"],
  ["12-payments", "/app/payments"],
  ["13-report", "/app/report"],
];

const issues = [];

async function shot(page, name) {
  await page.screenshot({ path: join(out, `${name}.png`), fullPage: true });
}

async function checkOverflow(page, label) {
  const overflow = await page.evaluate(() => {
    const doc = document.documentElement;
    return {
      scrollWidth: doc.scrollWidth,
      clientWidth: doc.clientWidth,
      bodyScroll: document.body.scrollWidth,
    };
  });
  if (overflow.scrollWidth > overflow.clientWidth + 2) {
    issues.push(`${label}: horizontal overflow ${overflow.scrollWidth}>${overflow.clientWidth}`);
  }
}

async function runViewport(width, height, prefix) {
  const browser = await chromium.launch();
  const context = await browser.newContext({
    viewport: { width, height },
    locale: "ru-RU",
  });
  const page = await context.newPage();
  page.on("pageerror", (e) => issues.push(`${prefix} pageerror: ${e.message}`));
  page.on("console", (msg) => {
    if (msg.type() === "error") issues.push(`${prefix} console: ${msg.text()}`);
  });

  for (const [name, path] of pages) {
    await page.goto(BASE + path, { waitUntil: "networkidle", timeout: 30000 });
    await page.waitForTimeout(600);
    await shot(page, `${prefix}-${name}`);
    await checkOverflow(page, `${prefix} ${path}`);
  }

  // Chat tax flow
  await page.goto(
    BASE +
      "/app/chat?q=" +
      encodeURIComponent("Сколько мне отложить на налог за этот месяц?"),
    { waitUntil: "networkidle" },
  );
  await page.waitForTimeout(2500);
  await shot(page, `${prefix}-12-chat-tax`);
  const taxVisible = await page.getByRole("region", { name: /Карточка налога/i }).isVisible().catch(() => false);
  if (!taxVisible) issues.push(`${prefix}: TaxCard not visible`);

  // Onboarding
  await page.goto(
    BASE +
      "/app/chat?q=" +
      encodeURIComponent(
        "Привет! Я делаю маникюр на дому в Казани, где-то 180 тысяч в месяц, пока без ИП.",
      ),
    { waitUntil: "networkidle" },
  );
  await page.waitForTimeout(2500);
  await shot(page, `${prefix}-13-chat-onboarding`);

  // Home segments
  await page.goto(BASE + "/app", { waitUntil: "networkidle" });
  await page.waitForTimeout(800);
  for (const label of ["Сегодня", "Путь"]) {
    const btn = page.getByRole("button", { name: label });
    if (await btn.count()) {
      await btn.click();
      await page.waitForTimeout(300);
      await shot(page, `${prefix}-14-home-${label}`);
    }
  }

  await browser.close();
}

await runViewport(1280, 800, "desk");
await runViewport(390, 844, "mobi");

console.log("SHOTS_DIR", out);
console.log("ISSUES_COUNT", issues.length);
for (const i of issues) console.log("ISSUE:", i);
