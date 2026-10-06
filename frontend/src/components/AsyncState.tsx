import { ArrowClockwise, Tray } from "@phosphor-icons/react";
export function LoadingRows() {
  return (
    <div className="skeleton-list" aria-label="Loading…" aria-busy="true">
      {[1, 2, 3].map((item) => (
        <div className="skeleton-row" key={item} />
      ))}
    </div>
  );
}
export function ErrorState({
  message,
  retry,
}: {
  message: string;
  retry?: () => void;
}) {
  return (
    <section className="state-panel" role="alert">
      <h2>Couldn’t Load This View</h2>
      <p>{message}</p>
      {retry && (
        <button className="button secondary" onClick={retry}>
          <ArrowClockwise aria-hidden="true" />
          Retry
        </button>
      )}
    </section>
  );
}
export function EmptyState({
  title = "Nothing Here Yet",
  body,
  action,
}: {
  title?: string;
  body: string;
  action?: React.ReactNode;
}) {
  return (
    <section className="state-panel">
      <Tray size={28} aria-hidden="true" />
      <h2>{title}</h2>
      <p>{body}</p>
      {action}
    </section>
  );
}
