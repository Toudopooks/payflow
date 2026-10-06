import { useQuery } from "@tanstack/react-query";
import {
  ArrowRight,
  CheckCircle,
  Clock,
  WarningCircle,
} from "@phosphor-icons/react";
import { Link } from "react-router-dom";
import { api } from "../api/client";
import { EmptyState, ErrorState, LoadingRows } from "../components/AsyncState";
import { OrderTable } from "../components/OrderTable";
export function DashboardPage() {
  const query = useQuery({
    queryKey: ["orders"],
    queryFn: api.orders,
    retry: 1,
  });
  const orders = query.data ?? [];
  const completed = orders.filter((o) => o.status === "completed").length;
  const active = orders.filter((o) => o.status === "active").length;
  return (
    <>
      <header className="page-header">
        <div>
          <p className="kicker">Overview</p>
          <h1>Payment Operations</h1>
          <p>
            Track installment orders and payment activity in the local
            simulation.
          </p>
        </div>
        <Link className="button primary" to="/orders/new">
          Create Order
          <ArrowRight aria-hidden="true" />
        </Link>
      </header>
      <section className="metrics" aria-label="Order summary">
        <article>
          <span>Total orders</span>
          <strong>{orders.length}</strong>
          <Clock aria-hidden="true" />
        </article>
        <article>
          <span>Active plans</span>
          <strong>{active}</strong>
          <WarningCircle aria-hidden="true" />
        </article>
        <article>
          <span>Completed</span>
          <strong>{completed}</strong>
          <CheckCircle aria-hidden="true" />
        </article>
      </section>
      <section className="content-section">
        <div className="section-heading">
          <div>
            <h2>Recent Orders</h2>
            <p>Up to 100 orders, newest first.</p>
          </div>
          <Link to="/orders">View all orders</Link>
        </div>
        {query.isLoading ? (
          <LoadingRows />
        ) : query.isError ? (
          <ErrorState
            message={query.error.message}
            retry={() => query.refetch()}
          />
        ) : orders.length ? (
          <OrderTable orders={orders.slice(0, 6)} />
        ) : (
          <EmptyState
            body="Create a customer, merchant, and order to begin."
            action={
              <Link className="button secondary" to="/orders/new">
                Create Order
              </Link>
            }
          />
        )}
      </section>
    </>
  );
}
