import type {
  Customer,
  Installment,
  Order,
  Page,
  Payment,
  Plan,
  Refund,
  Transaction,
} from "../types";

const base = import.meta.env.VITE_API_URL ?? "";
export class APIError extends Error {
  constructor(
    public code: string,
    message: string,
    public status: number,
  ) {
    super(message);
  }
}
async function request<T>(path: string, init?: RequestInit): Promise<T> {
  const response = await fetch(`${base}${path}`, {
    ...init,
    headers: { "Content-Type": "application/json", ...init?.headers },
  });
  const body = await response.json().catch(() => null);
  if (!response.ok)
    throw new APIError(
      body?.error?.code ?? "REQUEST_FAILED",
      body?.error?.message ?? "The request failed. Try again.",
      response.status,
    );
  return body as T;
}
export const api = {
  orders: () =>
    request<{ data: Order[] }>("/api/v1/orders").then((v) => v.data),
  order: (id: string) => request<Order>(`/api/v1/orders/${id}`),
  createCustomer: (input: { name: string; email: string }) =>
    request<Customer>("/api/v1/customers", {
      method: "POST",
      body: JSON.stringify(input),
    }),
  createMerchant: (input: { name: string }) =>
    request<{ id: string; name: string }>("/api/v1/merchants", {
      method: "POST",
      body: JSON.stringify(input),
    }),
  createOrder: (input: {
    customer_id: string;
    merchant_id: string;
    amount: number;
    currency: string;
  }) =>
    request<Order>("/api/v1/orders", {
      method: "POST",
      body: JSON.stringify(input),
    }),
  createPlan: (input: { order_id: string; installment_count: number }) =>
    request<Plan>("/api/v1/payment-plans", {
      method: "POST",
      body: JSON.stringify(input),
    }),
  installments: (orderId: string) =>
    request<{ data: Installment[] }>(
      `/api/v1/orders/${orderId}/installments`,
    ).then((v) => v.data),
  pay: (
    input: {
      installment_id: string;
      amount: number;
      simulate_failure: boolean;
    },
    key: string,
  ) =>
    request<Payment>("/api/v1/payments", {
      method: "POST",
      headers: { "Idempotency-Key": key },
      body: JSON.stringify(input),
    }),
  payment: (id: string) => request<Payment>(`/api/v1/payments/${id}`),
  refund: (paymentId: string, input: { amount: number; reason: string }) =>
    request<Refund>(`/api/v1/payments/${paymentId}/refunds`, {
      method: "POST",
      body: JSON.stringify(input),
    }),
  transactions: (customerId: string, page = 1) =>
    request<Page<Transaction>>(
      `/api/v1/customers/${customerId}/transactions?page=${page}&page_size=20`,
    ),
};
