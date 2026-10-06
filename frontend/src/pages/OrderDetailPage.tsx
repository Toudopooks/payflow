import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import {
  ArrowLeft,
  CircleNotch,
  CreditCard,
  Receipt,
} from "@phosphor-icons/react";
import { useState } from "react";
import { Link, useParams } from "react-router-dom";
import { api } from "../api/client";
import { EmptyState, ErrorState, LoadingRows } from "../components/AsyncState";
import { StatusBadge } from "../components/StatusBadge";
import type { Installment, Payment } from "../types";
import { date, money } from "../utils/format";

export function OrderDetailPage() {
  const { id = "" } = useParams();
  const client = useQueryClient();
  const query = useQuery({
    queryKey: ["order", id],
    queryFn: () => api.order(id),
    enabled: !!id,
  });
  const [lastPayment, setLastPayment] = useState<Payment | null>(null);
  const [fail, setFail] = useState(false);
  const payment = useMutation({
    mutationFn: (item: Installment) =>
      api.pay(
        {
          installment_id: item.id,
          amount: item.amount,
          simulate_failure: fail,
        },
        crypto.randomUUID(),
      ),
    onSuccess: (value) => {
      setLastPayment(value);
      client.invalidateQueries({ queryKey: ["order", id] });
    },
  });
  const refund = useMutation({
    mutationFn: (value: Payment) =>
      api.refund(value.id, {
        amount: value.amount,
        reason: "Demo order refund",
      }),
    onSuccess: () => {
      setLastPayment(null);
      client.invalidateQueries({ queryKey: ["order", id] });
    },
  });
  if (query.isLoading) return <LoadingRows />;
  if (query.isError)
    return (
      <ErrorState message={query.error.message} retry={() => query.refetch()} />
    );
  if (!query.data) return null;
  const order = query.data;
  const installments = order.payment_plan?.installments ?? [];
  return (
    <>
      <header className="page-header compact">
        <div>
          <Link className="back-link" to="/orders">
            <ArrowLeft aria-hidden="true" />
            Orders
          </Link>
          <div className="title-row">
            <h1>{order.reference}</h1>
            <StatusBadge status={order.status} />
          </div>
          <p>
            Created {date(order.created_at)}. Reference IDs stay visible for
            debugging.
          </p>
        </div>
        <div className="order-total">
          <span>Total</span>
          <strong>{money(order.amount, order.currency)}</strong>
        </div>
      </header>
      <section className="order-meta" aria-label="Order references">
        <div>
          <span>Customer ID</span>
          <code>{order.customer_id}</code>
        </div>
        <div>
          <span>Merchant ID</span>
          <code>{order.merchant_id}</code>
        </div>
        <Link to={`/transactions?customer=${order.customer_id}`}>
          View Customer Transactions
        </Link>
      </section>
      <section className="detail-grid">
        <article className="content-section">
          <div className="section-heading">
            <div>
              <h2>Installment Plan</h2>
              <p>
                {order.payment_plan?.installment_count ?? 0} scheduled payments
              </p>
            </div>
          </div>
          {installments.length === 0 ? (
            <EmptyState body="This order does not have an installment plan yet." />
          ) : (
            <div className="installments">
              {installments.map((item) => (
                <div className="installment" key={item.id}>
                  <div className="installment-number">{item.sequence}</div>
                  <div className="installment-main">
                    <strong>{money(item.amount, order.currency)}</strong>
                    <span>Due {date(item.due_at)}</span>
                  </div>
                  <StatusBadge status={item.status} />
                  <button
                    className="button secondary pay-button"
                    disabled={
                      payment.isPending ||
                      item.status === "paid" ||
                      item.status === "refunded"
                    }
                    onClick={() => payment.mutate(item)}
                  >
                    {payment.isPending && payment.variables?.id === item.id ? (
                      <CircleNotch className="spin" aria-hidden="true" />
                    ) : (
                      <CreditCard aria-hidden="true" />
                    )}
                    {item.status === "paid" ? "Paid" : "Pay"}
                  </button>
                </div>
              ))}
            </div>
          )}
        </article>
        <aside className="content-section simulation">
          <h2>Payment Simulation</h2>
          <p>
            Turn on the decline switch to test the failure path. It works only
            when backend demo mode is enabled.
          </p>
          <label className="switch-row">
            <input
              type="checkbox"
              checked={fail}
              onChange={(event) => setFail(event.target.checked)}
            />
            <span>Simulate a declined payment</span>
          </label>
          {payment.isError && (
            <div className="inline-error" role="alert">
              {payment.error.message} You can retry the installment.
            </div>
          )}
          {lastPayment && (
            <div className="payment-result" aria-live="polite">
              <Receipt aria-hidden="true" />
              <div>
                <strong>Payment {lastPayment.status}</strong>
                <span className="break-word">{lastPayment.id}</span>
              </div>
              {lastPayment.status === "succeeded" && (
                <button
                  className="button danger"
                  disabled={refund.isPending}
                  onClick={() => {
                    if (
                      window.confirm(
                        `Refund ${money(lastPayment.amount, lastPayment.currency)}?`,
                      )
                    )
                      refund.mutate(lastPayment);
                  }}
                >
                  {refund.isPending && (
                    <CircleNotch className="spin" aria-hidden="true" />
                  )}
                  Refund Payment
                </button>
              )}
            </div>
          )}
          {refund.isError && (
            <div className="inline-error" role="alert">
              {refund.error.message}
            </div>
          )}
        </aside>
      </section>
    </>
  );
}
