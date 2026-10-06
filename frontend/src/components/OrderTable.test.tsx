import { render, screen } from "@testing-library/react";
import { MemoryRouter } from "react-router-dom";
import { OrderTable } from "./OrderTable";

test("renders order status and locale-aware money", () => {
  render(
    <MemoryRouter>
      <OrderTable
        orders={[
          {
            id: "order-1",
            reference: "PF-A12C",
            customer_id: "customer-1",
            merchant_id: "merchant-1",
            amount: 12000,
            currency: "USD",
            status: "active",
            created_at: "2026-01-10T12:00:00Z",
          },
        ]}
      />
    </MemoryRouter>,
  );
  expect(screen.getByRole("link", { name: "PF-A12C" })).toBeInTheDocument();
  expect(screen.getByText("Active")).toBeInTheDocument();
  expect(screen.getByText(/120\.00/)).toBeInTheDocument();
});
