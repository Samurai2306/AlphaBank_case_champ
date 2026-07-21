import Link from "next/link";

export default function NotFound() {
  return (
    <main className="mx-auto flex min-h-screen max-w-lg flex-col justify-center px-4">
      <p className="text-sm text-ink/50">404</p>
      <h1 className="mt-2 text-3xl font-bold">Страница не найдена</h1>
      <p className="mt-3 text-ink/60">
        Вернитесь в кабинет или на главную страницу продукта.
      </p>
      <div className="mt-6 flex flex-wrap gap-3">
        <Link href="/login" className="rounded-pill bg-brand px-5 py-2.5 text-sm font-semibold text-white">
          Войти
        </Link>
        <Link href="/app?example=1" className="rounded-pill bg-ink px-5 py-2.5 text-sm font-semibold text-white">
          Пример Маши
        </Link>
        <Link href="/" className="rounded-pill bg-white px-5 py-2.5 text-sm font-medium shadow-soft">
          На сайт
        </Link>
      </div>
    </main>
  );
}
