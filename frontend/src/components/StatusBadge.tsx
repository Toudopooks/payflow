const labels: Record<string, string> = {
  pending: "Pending",
  active: "Active",
  completed: "Completed",
  paid: "Paid",
  succeeded: "Succeeded",
  failed: "Failed",
  refunded: "Refunded",
  cancelled: "Cancelled",
};
export function StatusBadge({ status }: { status: string }) {
  return (
    <span className={`status status-${status}`}>
      <span aria-hidden="true" className="status-mark" />
      {labels[status] ?? status}
    </span>
  );
}
