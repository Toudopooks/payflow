import { useQuery } from "@tanstack/react-query";
import { MagnifyingGlass } from "@phosphor-icons/react";
import { FormEvent, useState } from "react";
import { useSearchParams } from "react-router-dom";
import { api } from "../api/client";
import { EmptyState, ErrorState, LoadingRows } from "../components/AsyncState";
import { StatusBadge } from "../components/StatusBadge";
import { date, money } from "../utils/format";
export function TransactionsPage() {
  const [params, setParams] = useSearchParams();
  const customer = params.get("customer") ?? "";
  const [input, setInput] = useState(customer);
  const query = useQuery({
    queryKey: ["transactions", customer],
    queryFn: () => api.transactions(customer),
    enabled: !!customer,
    retry: false,
  });
  const submit = (event: FormEvent) => {
    event.preventDefault();
    const value = input.trim();
    setParams(value ? { customer: value } : {});
  };
  return (
    <>
      <header className="page-header">
        <div>
          <p className="kicker">Ledger</p>
          <h1>Transaction History</h1>
          <p>Look up the immutable event stream for a customer ID.</p>
        </div>
      </header>
      <section className="content-section">
        <form className="search-form" onSubmit={submit}>
          <label htmlFor="customer-id">Customer ID</label>
          <div>
            <input
              id="customer-id"
              name="customer"
              value={input}
              onChange={(event) => setInput(event.target.value)}
              autoComplete="off"
              spellCheck={false}
              placeholder="Paste a customer UUID…"
            />
            <button className="button primary" type="submit">
              <MagnifyingGlass aria-hidden="true" />
              Find Transactions
            </button>
          </div>
        </form>
        {!customer ? (
          <EmptyState
            title="Enter a Customer ID"
            body="Customer IDs are returned when a customer is created."
          />
        ) : query.isLoading ? (
          <LoadingRows />
        ) : query.isError ? (
          <ErrorState message={query.error.message} />
        ) : query.data?.data.length ? (
          <div className="timeline">
            {query.data.data.map((item) => (
              <article key={item.id}>
                <div>
                  <strong>{item.type.replaceAll("_", " ")}</strong>
                  <span>{date(item.created_at)}</span>
                </div>
                <StatusBadge status={item.status} />
                <strong className="numeric">
                  {money(item.amount, item.currency)}
                </strong>
              </article>
            ))}
          </div>
        ) : (
          <EmptyState body="This customer has no transaction events." />
        )}
      </section>
    </>
  );
}
