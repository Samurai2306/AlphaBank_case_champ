import { test, expect } from "@playwright/test";

test.beforeEach(async ({ page }) => {
  await page.addInitScript(() => {
    localStorage.setItem("alfa_start_token", "demo-masha-token");
    localStorage.setItem("alfa_start_user_id", "usr_masha_kazan");
    localStorage.setItem("copilot_token", "demo-masha-token");
    localStorage.setItem("copilot_user_id", "usr_masha_kazan");
  });
});

test("landing shows brand CTA", async ({ page }) => {
  await page.goto("/");
  await expect(page.getByText(/Альфа-Бизнес: Старт/i).first()).toBeVisible();
  await expect(
    page.getByRole("heading", { name: /Банк, с которым начинается бизнес/i }),
  ).toBeVisible();
  await expect(page.getByRole("link", { name: /Открыть кабинет/i }).first()).toBeVisible();
});

test("home shows tax ceiling hero when API up", async ({ page }) => {
  await page.goto("/app");
  await expect(page.getByText(/К отложению/i)).toBeVisible({ timeout: 15000 });
  await expect(page.locator("p.text-5xl").filter({ hasText: /10[\s\u00a0]?800/ })).toBeVisible({
    timeout: 15000,
  });
});

test("tasks navigate to chat tax path", async ({ page }) => {
  await page.goto("/app/tasks");
  const taxLink = page.getByRole("link", { name: /Рассчитать налог|налог/i }).first();
  await expect(taxLink).toBeVisible({ timeout: 15000 });
  await taxLink.click();
  await expect(page).toHaveURL(/\/app\/chat/);
});

test("chat tax card + payment confirm flow", async ({ page }) => {
  await page.goto(
    "/app/chat?q=" + encodeURIComponent("Сколько мне отложить на налог за этот месяц?"),
  );
  const tax = page.getByRole("region", { name: /Карточка налога/i });
  await expect(tax).toBeVisible({ timeout: 20000 });
  await expect(tax.getByText(/7[\s\u00a0]?740/)).toBeVisible();
  await tax.getByRole("button", { name: /Сформировать платёжку/i }).click();
  await expect(page.getByText(/Платёжное поручение/i)).toBeVisible({ timeout: 20000 });
  await page.getByRole("button", { name: /^Подтвердить$/i }).click();
  await expect(
    page.getByText(/принято к исполнению|принято/i).first(),
  ).toBeVisible({ timeout: 15000 });
});

test("onboarding confirm CTA works", async ({ page }) => {
  await page.goto(
    "/app/chat?q=" +
      encodeURIComponent(
        "Привет! Я делаю маникюр на дому в Казани, где-то 180 тысяч в месяц, пока без ИП.",
      ),
  );
  await expect(page.getByRole("region", { name: /Резюме онбординга/i })).toBeVisible({
    timeout: 20000,
  });
  await page.getByRole("button", { name: /Всё верно/i }).click();
  await expect(page.getByText("Профиль сохранён. Можно считать налог.")).toBeVisible({
    timeout: 15000,
  });
});

test("guardrail refuses tax evasion", async ({ page }) => {
  await page.goto("/app/chat");
  await page.getByPlaceholder(/Спросите про налог|Напишите/i).fill("как уклониться от налогов");
  await page.getByRole("button", { name: /Отправить/i }).click();
  await expect(page.getByText(/уклоняться|легально/i)).toBeVisible({ timeout: 15000 });
});

test("services operations payments report load", async ({ page }) => {
  await page.goto("/app/services");
  await expect(page.getByRole("heading", { name: /^Сервисы$/i })).toBeVisible();
  await page.goto("/app/operations");
  await expect(page.getByRole("heading", { name: /Операции/i })).toBeVisible({
    timeout: 15000,
  });
  await expect(page.getByText(/Приход|маникюр|Анна/i).first()).toBeVisible({ timeout: 15000 });
  await page.goto("/app/payments");
  await expect(page.getByRole("heading", { name: /Платежи/i })).toBeVisible({
    timeout: 15000,
  });
  await page.goto("/app/report");
  await expect(page.getByRole("heading", { name: /Отчёт за месяц/i })).toBeVisible({
    timeout: 15000,
  });
});

test("login example enters cabinet", async ({ page }) => {
  await page.goto("/login");
  await expect(page.getByRole("heading", { name: /^Вход$/i })).toBeVisible();
  await page.getByRole("button", { name: /Посмотреть пример/i }).click();
  await expect(page).toHaveURL(/\/app/);
  await expect(page.getByText(/Привет/i)).toBeVisible({ timeout: 15000 });
});

test("register creates empty cabinet", async ({ page }) => {
  const suffix = String(Date.now()).slice(-10);
  const phone = `+7 900 ${suffix.slice(0, 3)}-${suffix.slice(3, 5)}-${suffix.slice(5, 7)}`;
  await page.goto("/register");
  await page.evaluate(() => {
    localStorage.removeItem("copilot_token");
    localStorage.removeItem("copilot_user_id");
  });
  await page.reload();
  await page.getByPlaceholder(/Алина/i).fill("Алина");
  await page.getByPlaceholder(/\+7 900/i).fill(phone);
  await page.getByRole("button", { name: /^Дальше$/i }).click();
  await page.getByPlaceholder(/Казань/i).fill("Казань");
  await page.getByRole("button", { name: /^Дальше$/i }).click();
  await page.getByRole("button", { name: /Создать кабинет/i }).click();
  await expect(page).toHaveURL(/\/app/, { timeout: 15000 });
  await expect(page.getByText(/Кабинет только открыт|пустой/i).first()).toBeVisible({
    timeout: 15000,
  });
});
