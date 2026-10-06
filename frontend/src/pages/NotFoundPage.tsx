import { ArrowLeft } from "@phosphor-icons/react";
import { Link } from "react-router-dom";
export function NotFoundPage() {
  return (
    <section className="state-panel full-page">
      <h1>Page Not Found</h1>
      <p>The route does not exist. Return to the payment overview.</p>
      <Link className="button primary" to="/">
        <ArrowLeft aria-hidden="true" />
        Back to Overview
      </Link>
    </section>
  );
}
