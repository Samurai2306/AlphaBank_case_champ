"use client";

import { motion, useReducedMotion } from "framer-motion";
import type { SduiEnvelope } from "@/lib/api";
import { TaxCard } from "./TaxCard";
import { PaymentDraftCard } from "./PaymentDraftCard";
import { UnitEconomicsChart } from "./UnitEconomicsChart";
import { ComplianceTrafficLight } from "./ComplianceTrafficLight";
import { OnboardingSummary } from "./OnboardingSummary";
import { EnpPiggyBank } from "./EnpPiggyBank";
import { LegalFlagsList } from "./LegalFlagsList";
import { KnowledgeSources } from "./KnowledgeSources";
import { ProductOffers } from "./ProductOffers";

export function GenUIRenderer({
  envelope,
  onAction,
}: {
  envelope: SduiEnvelope;
  onAction?: (action: string, payload?: Record<string, unknown>) => void;
}) {
  const reduce = useReducedMotion();
  const body = (() => {
    switch (envelope.component) {
      case "TaxCard":
        return <TaxCard props={envelope.props} onAction={onAction} />;
      case "PaymentDraftCard":
        return <PaymentDraftCard props={envelope.props} onAction={onAction} />;
      case "UnitEconomicsChart":
        return <UnitEconomicsChart props={envelope.props} />;
      case "ComplianceTrafficLight":
        return <ComplianceTrafficLight props={envelope.props} onAction={onAction} />;
      case "OnboardingSummary":
        return <OnboardingSummary props={envelope.props} onAction={onAction} />;
      case "EnpPiggyBank":
        return <EnpPiggyBank props={envelope.props} />;
      case "LegalFlagsList":
        return <LegalFlagsList props={envelope.props} />;
      case "KnowledgeSources":
        return <KnowledgeSources props={envelope.props} />;
      case "ProductOffers":
        return <ProductOffers props={envelope.props} onAction={onAction} />;
      default:
        return (
          <div className="rounded-card bg-white p-4 text-sm text-ink/70 shadow-soft">
            Карточка временно недоступна. Продолжите в чате или выберите частый
            запрос.
          </div>
        );
    }
  })();

  return (
    <motion.div
      initial={reduce ? false : { opacity: 0, y: 10 }}
      animate={{ opacity: 1, y: 0 }}
      transition={{ duration: 0.28, ease: "easeOut" }}
    >
      {body}
    </motion.div>
  );
}
