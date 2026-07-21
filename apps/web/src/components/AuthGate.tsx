"use client";

import { useEffect, useState } from "react";
import { usePathname, useRouter, useSearchParams } from "next/navigation";
import { enterExampleSession, hasSession } from "@/lib/session";

/** Redirects unauthenticated users to /register; ?example=1 opens Маша. */
export function AuthGate({ children }: { children: React.ReactNode }) {
  const router = useRouter();
  const pathname = usePathname();
  const search = useSearchParams();
  const [ready, setReady] = useState(false);

  useEffect(() => {
    if (search.get("example") === "1") {
      enterExampleSession();
      router.replace(pathname || "/app");
      setReady(true);
      return;
    }
    if (!hasSession()) {
      const next = encodeURIComponent(pathname || "/app");
      router.replace(`/login?next=${next}`);
      return;
    }
    setReady(true);
  }, [pathname, router, search]);

  if (!ready) {
    return (
      <div className="flex min-h-[40vh] items-center justify-center text-sm text-ink/50">
        Открываем кабинет…
      </div>
    );
  }
  return <>{children}</>;
}
