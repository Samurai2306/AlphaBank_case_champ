"use client";

import {
  CartesianGrid,
  Line,
  LineChart,
  ReferenceDot,
  ReferenceLine,
  ResponsiveContainer,
  Tooltip,
  XAxis,
  YAxis,
} from "recharts";
import { formatRub } from "@/lib/api";

function clientsWord(n: number): string {
  const abs = Math.abs(Math.round(n)) % 100;
  const d = abs % 10;
  if (abs > 10 && abs < 20) return "клиентов";
  if (d === 1) return "клиент";
  if (d >= 2 && d <= 4) return "клиента";
  return "клиентов";
}

export function UnitEconomicsChart({
  props,
}: {
  props: Record<string, unknown>;
}) {
  const series = (props.series as { units: number; profit: number }[]) ?? [];
  const be = Number(props.breakeven_units_per_day ?? 0);
  const recommended = Number(props.recommended_units_per_day ?? be);
  const taxPerUnit = Number(props.tax_per_unit ?? 0);
  const days = Number(props.working_days_per_month ?? 22);
  const bePoint = series.find((p) => p.units >= be) ?? series[0];
  const zeroPoint = series.find((p) => p.units === 0);

  return (
    <div
      className="rounded-card bg-white p-5 shadow-soft"
      role="region"
      aria-label="Точка безубыточности"
    >
      <p className="font-semibold">Точка безубыточности</p>
      <p className="mt-2 text-3xl font-bold tabular">
        {be}{" "}
        <span className="text-lg font-medium text-ink/60">
          {clientsWord(be)}/день
        </span>
      </p>
      {recommended > be ? (
        <p className="mt-1 text-sm text-ink/70">
          С запасом 15%: ≈ {recommended} {clientsWord(recommended)}/день
        </p>
      ) : null}
      <p className="mt-2 text-sm text-ink/60">
        Чистая маржа {formatRub(Number(props.margin_per_unit ?? 0))} ₽ · цена{" "}
        {formatRub(Number(props.price_per_unit ?? 0))} ₽ · фикс{" "}
        {formatRub(Number(props.fixed_costs_monthly ?? 0))} ₽
        {taxPerUnit > 0
          ? ` · налог ~${formatRub(taxPerUnit)} ₽/клиент`
          : ""}
        {` · ${days} раб. дней`}
      </p>
      {series.length === 0 ? (
        <p className="mt-4 rounded-2xl bg-canvas px-3 py-3 text-sm text-ink/60">
          График недоступен — ориентир по числу клиентов выше.
        </p>
      ) : (
        <div className="mt-4 h-48">
          <ResponsiveContainer width="100%" height="100%">
            <LineChart
              data={series}
              margin={{ top: 8, right: 8, left: 0, bottom: 0 }}
            >
              <CartesianGrid stroke="#eee" strokeDasharray="3 3" />
              <XAxis
                dataKey="units"
                tick={{ fontSize: 11 }}
                label={{
                  value: "Клиентов в день",
                  position: "insideBottom",
                  offset: -2,
                  fontSize: 11,
                }}
              />
              <YAxis
                tick={{ fontSize: 11 }}
                tickFormatter={(v) => formatRub(Number(v))}
                width={56}
                label={{
                  value: "Прибыль, ₽",
                  angle: -90,
                  position: "insideLeft",
                  fontSize: 11,
                }}
              />
              <Tooltip
                formatter={(v) => [`${formatRub(Number(v))} ₽`, "Прибыль"]}
                labelFormatter={(u) =>
                  `${u} ${clientsWord(Number(u))}/день`
                }
              />
              <ReferenceLine y={0} stroke="#94a3b8" strokeDasharray="4 4" />
              <Line
                type="linear"
                dataKey="profit"
                stroke="#EF3124"
                strokeWidth={2}
                dot={{ r: 3, fill: "#EF3124" }}
                activeDot={{ r: 5 }}
              />
              {zeroPoint ? (
                <ReferenceDot
                  x={0}
                  y={zeroPoint.profit}
                  r={5}
                  fill="#64748b"
                  stroke="#fff"
                  strokeWidth={2}
                />
              ) : null}
              {bePoint ? (
                <ReferenceDot
                  x={bePoint.units}
                  y={bePoint.profit}
                  r={6}
                  fill="#EF3124"
                  stroke="#fff"
                  strokeWidth={2}
                />
              ) : null}
            </LineChart>
          </ResponsiveContainer>
        </div>
      )}
    </div>
  );
}
