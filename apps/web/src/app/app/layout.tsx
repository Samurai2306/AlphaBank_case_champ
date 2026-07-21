import { Suspense } from "react";
import { AppNav } from "@/components/AppNav";
import { AuthGate } from "@/components/AuthGate";

export default function AppLayout({ children }: { children: React.ReactNode }) {
  return (
    <div className="min-h-screen pb-24">
      <Suspense
        fallback={
          <div className="flex min-h-[40vh] items-center justify-center text-sm text-ink/50">
            Открываем кабинет…
          </div>
        }
      >
        <AuthGate>{children}</AuthGate>
      </Suspense>
      <AppNav />
    </div>
  );
}
