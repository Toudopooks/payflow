import type { Currency } from "../types";
export function money(amount: number, currency: Currency) {
  const divisor = currency === "CLP" ? 1 : 100;
  return new Intl.NumberFormat(undefined, {
    style: "currency",
    currency,
    maximumFractionDigits: currency === "CLP" ? 0 : 2,
  }).format(amount / divisor);
}
export function date(value: string) {
  return new Intl.DateTimeFormat(undefined, { dateStyle: "medium" }).format(
    new Date(value),
  );
}
