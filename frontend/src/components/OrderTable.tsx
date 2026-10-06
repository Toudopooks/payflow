import { Link } from "react-router-dom";
import type { Order } from "../types";
import { date, money } from "../utils/format";
import { StatusBadge } from "./StatusBadge";
export function OrderTable({ orders }: { orders: Order[] }) {
  return (
    <div className="table-wrap">
      <table>
        <thead>
          <tr>
            <th>Order</th>
            <th>Created</th>
            <th>Status</th>
            <th className="numeric">Total</th>
          </tr>
        </thead>
        <tbody>
          {orders.map((order) => (
            <tr key={order.id}>
              <td>
                <Link className="primary-link" to={`/orders/${order.id}`}>
                  {order.reference}
                </Link>
                <span className="cell-subtitle">{order.id}</span>
              </td>
              <td>{date(order.created_at)}</td>
              <td>
                <StatusBadge status={order.status} />
              </td>
              <td className="numeric">{money(order.amount, order.currency)}</td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  );
}
