import { zodResolver } from "@hookform/resolvers/zod";
import { useMutation, useQueryClient } from "@tanstack/react-query";
import { ArrowLeft, CircleNotch } from "@phosphor-icons/react";
import { cloneElement, isValidElement, useEffect } from "react";
import { useForm } from "react-hook-form";
import { Link, useNavigate } from "react-router-dom";
import { api } from "../api/client";
import { orderSchema, type OrderForm } from "../schemas/order";
export function CreateOrderPage() {
  const navigate = useNavigate();
  const client = useQueryClient();
  const form = useForm<OrderForm>({
    resolver: zodResolver(orderSchema),
    defaultValues: {
      customerName: "",
      customerEmail: "",
      merchantName: "",
      amount: 12000,
      currency: "USD",
      installmentCount: 4,
    },
  });
  const dirty = form.formState.isDirty;
  useEffect(() => {
    const warn = (event: BeforeUnloadEvent) => {
      if (dirty) {
        event.preventDefault();
      }
    };
    const guard = (event: MouseEvent) => {
      const target = event.target;
      if (!(target instanceof Element)) return;
      const link = target.closest("a");
      if (dirty && link && !window.confirm("Discard the unfinished order?")) {
        event.preventDefault();
      }
    };
    window.addEventListener("beforeunload", warn);
    document.addEventListener("click", guard);
    return () => {
      window.removeEventListener("beforeunload", warn);
      document.removeEventListener("click", guard);
    };
  }, [dirty]);
  const mutation = useMutation({
    mutationFn: async (raw: OrderForm) => {
      const data = orderSchema.parse(raw);
      const customer = await api.createCustomer({
        name: data.customerName.trim(),
        email: data.customerEmail.trim(),
      });
      const merchant = await api.createMerchant({
        name: data.merchantName.trim(),
      });
      const order = await api.createOrder({
        customer_id: customer.id,
        merchant_id: merchant.id,
        amount: data.amount,
        currency: data.currency,
      });
      await api.createPlan({
        order_id: order.id,
        installment_count: data.installmentCount,
      });
      return order;
    },
    onSuccess: (order) => {
      client.invalidateQueries({ queryKey: ["orders"] });
      form.reset();
      navigate(`/orders/${order.id}`);
    },
  });
  const submit = form.handleSubmit(
    (values) => mutation.mutate(values),
    (errors) => {
      const first = Object.keys(errors)[0] as keyof OrderForm | undefined;
      if (first) form.setFocus(first);
    },
  );
  return (
    <>
      <header className="page-header compact">
        <div>
          <Link className="back-link" to="/orders">
            <ArrowLeft aria-hidden="true" />
            Orders
          </Link>
          <h1>Create an Order</h1>
          <p>Amounts use minor units. For USD, 12000 means $120.00.</p>
        </div>
      </header>
      <section className="form-panel">
        <form onSubmit={submit} noValidate>
          <fieldset disabled={mutation.isPending}>
            <legend>Customer</legend>
            <div className="form-grid">
              <Field
                label="Name"
                error={form.formState.errors.customerName?.message}
              >
                <input
                  {...form.register("customerName")}
                  name="customerName"
                  autoComplete="name"
                  required
                  placeholder="Maya Rivera…"
                  aria-invalid={!!form.formState.errors.customerName}
                />
              </Field>
              <Field
                label="Email"
                error={form.formState.errors.customerEmail?.message}
              >
                <input
                  {...form.register("customerEmail")}
                  name="customerEmail"
                  type="email"
                  autoComplete="email"
                  spellCheck={false}
                  required
                  placeholder="maya@example.test…"
                  aria-invalid={!!form.formState.errors.customerEmail}
                />
              </Field>
            </div>
          </fieldset>
          <fieldset disabled={mutation.isPending}>
            <legend>Order</legend>
            <div className="form-grid">
              <Field
                label="Merchant"
                error={form.formState.errors.merchantName?.message}
              >
                <input
                  {...form.register("merchantName")}
                  name="merchantName"
                  autoComplete="organization"
                  required
                  placeholder="North Harbor Goods…"
                  aria-invalid={!!form.formState.errors.merchantName}
                />
              </Field>
              <Field
                label="Amount in minor units"
                error={form.formState.errors.amount?.message}
              >
                <input
                  {...form.register("amount")}
                  name="amount"
                  type="number"
                  inputMode="numeric"
                  autoComplete="off"
                  min="1"
                  required
                  aria-invalid={!!form.formState.errors.amount}
                />
              </Field>
              <Field label="Currency">
                <select
                  {...form.register("currency")}
                  name="currency"
                  autoComplete="off"
                  required
                >
                  <option value="USD">USD</option>
                  <option value="CLP">CLP</option>
                </select>
              </Field>
              <Field label="Installments">
                <select
                  {...form.register("installmentCount")}
                  name="installmentCount"
                  autoComplete="off"
                  required
                >
                  <option value="2">2 installments</option>
                  <option value="3">3 installments</option>
                  <option value="4">4 installments</option>
                </select>
              </Field>
            </div>
          </fieldset>
          {mutation.isError && (
            <div className="inline-error" role="alert">
              {mutation.error.message} Check the fields and try again.
            </div>
          )}
          <div className="form-actions">
            <Link className="button secondary" to="/orders">
              Cancel
            </Link>
            <button
              className="button primary"
              type="submit"
              disabled={mutation.isPending}
            >
              {mutation.isPending && (
                <CircleNotch className="spin" aria-hidden="true" />
              )}
              Create Order
            </button>
          </div>
          <span className="sr-only" aria-live="polite">
            {mutation.isPending ? "Creating order…" : ""}
          </span>
        </form>
      </section>
    </>
  );
}
function Field({
  label,
  error,
  children,
}: {
  label: string;
  error?: string;
  children: React.ReactElement<{
    name?: string;
    id?: string;
    "aria-describedby"?: string;
  }>;
}) {
  if (!isValidElement(children)) return null;
  const id = `field-${children.props.name ?? label.toLowerCase().replaceAll(" ", "-")}`;
  const errorID = `${id}-error`;
  return (
    <label className="field" htmlFor={id}>
      <span>{label}</span>
      {cloneElement(children, {
        id,
        "aria-describedby": error ? errorID : undefined,
      })}
      {error && (
        <small id={errorID} className="field-error" aria-live="polite">
          {error}
        </small>
      )}
    </label>
  );
}
