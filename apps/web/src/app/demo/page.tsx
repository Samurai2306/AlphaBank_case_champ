"use client";

import { useEffect } from "react";
import { useRouter } from "next/navigation";
import { enterExampleSession } from "@/lib/session";

/** Jury/demo entry: open Маша cabinet on tasks. */
export default function DemoRedirect() {
  const router = useRouter();
  useEffect(() => {
    enterExampleSession();
    router.replace("/app/tasks");
  }, [router]);
  return (
    <main className="flex min-h-screen items-center justify-center text-sm text-ink/50">
      Открываем пример кабинета…
    </main>
  );
}
