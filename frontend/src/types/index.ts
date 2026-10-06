export type Currency = "USD" | "CLP";
export interface Customer {
  id: string;
  name: string;
  email: string;
  created_at: string;
}
export interface Merchant {
  id: string;
  name: string;
  created_at: string;
}
export interface Installment {
  id: string;
  payment_plan_id: string;
  sequence: number;
  amount: number;
  status: "pending" | "paid" | "failed" | "refunded";
  due_at: string;
  paid_at?: string;
}
export interface Plan {
  id: string;
  order_id: string;
  installment_count: number;
  status: string;
  installments: Installment[];
}
export interface Order {
  id: string;
  reference: string;
  customer_id: string;
  merchant_id: string;
  amount: number;
  currency: Currency;
  status: string;
  payment_plan?: Plan;
  created_at: string;
}
export interface Payment {
  id: string;
  installment_id: string;
  amount: number;
  currency: Currency;
  status: string;
  failure_reason?: string;
  refunded_amount: number;
  created_at: string;
}
export interface Refund {
  id: string;
  payment_id: string;
  amount: number;
  status: string;
  reason: string;
  created_at: string;
}
export interface Transaction {
  id: string;
  order_id: string;
  type: string;
  status: string;
  amount: number;
  currency: Currency;
  reference_id: string;
  created_at: string;
}
export interface Page<T> {
  data: T[];
  page: number;
  page_size: number;
  total: number;
  total_pages: number;
}
