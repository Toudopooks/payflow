import { Route, Routes } from "react-router-dom";
import { useEffect } from "react";
import { useLocation } from "react-router-dom";
import { AppShell } from "./layout/AppShell";
import { CreateOrderPage } from "./pages/CreateOrderPage";
import { DashboardPage } from "./pages/DashboardPage";
import { NotFoundPage } from "./pages/NotFoundPage";
import { OrderDetailPage } from "./pages/OrderDetailPage";
import { OrdersPage } from "./pages/OrdersPage";
import { TransactionsPage } from "./pages/TransactionsPage";
export default function App() {
  const location = useLocation();
  useEffect(() => {
    const path = location.pathname;
    const page =
      path === "/"
        ? "Overview"
        : path === "/orders"
          ? "Orders"
          : path === "/orders/new"
            ? "Create Order"
            : path.startsWith("/orders/")
              ? "Order Details"
              : path === "/transactions"
                ? "Transactions"
                : "Not Found";
    document.title = `${page} | PayFlow`;
  }, [location.pathname]);
  return (
    <Routes>
      <Route element={<AppShell />}>
        <Route index element={<DashboardPage />} />
        <Route path="orders" element={<OrdersPage />} />
        <Route path="orders/new" element={<CreateOrderPage />} />
        <Route path="orders/:id" element={<OrderDetailPage />} />
        <Route path="transactions" element={<TransactionsPage />} />
        <Route path="*" element={<NotFoundPage />} />
      </Route>
    </Routes>
  );
}
