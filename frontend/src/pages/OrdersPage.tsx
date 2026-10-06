import { useQuery } from "@tanstack/react-query";
import { Plus } from "@phosphor-icons/react";
import { Link } from "react-router-dom";
import { api } from "../api/client";
import { EmptyState, ErrorState, LoadingRows } from "../components/AsyncState";
import { OrderTable } from "../components/OrderTable";
export function OrdersPage() {
  const query = useQuery({ queryKey: ["orders"], queryFn: api.orders });
  return (
    <>
      <header className="page-header">
        <div>
          <p className="kicker">Orders</p>
          <h1>Installment Orders</h1>
          <p>Review order value, plan status, and payment progress.</p>
        </div>
        <Link className="button primary" to="/orders/new">
          <Plus aria-hidden="true" />
          Create Order
        </Link>
      </header>
      <section className="content-section">
        {query.isLoading ? (
          <LoadingRows />
        ) : query.isError ? (
          <ErrorState
            message={query.error.message}
            retry={() => query.refetch()}
          />
        ) : query.data?.length ? (
          <OrderTable orders={query.data} />
        ) : (
          <EmptyState body="No orders have been created." />
        )}
      </section>
    </>
  );
}
