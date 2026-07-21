export const INTENTS = [
  "TAX_CALC",
  "LEGAL_REVIEW",
  "ONBOARDING",
  "GENERAL_QA",
  "TRANSACTION",
  "UNIT_ECON",
  "PIGGY",
  "COMPLIANCE",
  "GUARD",
] as const;

export type Intent = (typeof INTENTS)[number];

export const SDUI_COMPONENTS = [
  "TaxCard",
  "PaymentDraftCard",
  "UnitEconomicsChart",
  "ComplianceTrafficLight",
  "OnboardingSummary",
  "EnpPiggyBank",
  "LegalFlagsList",
  "KnowledgeSources",
] as const;

export type SduiComponent = (typeof SDUI_COMPONENTS)[number];

export type SduiEnvelope = {
  schema_version: number;
  component: SduiComponent | string;
  props: Record<string, unknown>;
};
