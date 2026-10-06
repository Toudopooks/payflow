import {
  ChartLineUp,
  ListBullets,
  Plus,
  Receipt,
  X,
} from "@phosphor-icons/react";
import { useEffect, useRef, useState } from "react";
import { NavLink, Outlet, useLocation } from "react-router-dom";

const items = [
  { to: "/", label: "Overview", icon: ChartLineUp },
  { to: "/orders", label: "Orders", icon: ListBullets },
  { to: "/orders/new", label: "Create Order", icon: Plus },
  { to: "/transactions", label: "Transactions", icon: Receipt },
];
export function AppShell() {
  const [open, setOpen] = useState(false);
  const [mobile, setMobile] = useState(
    () => window.matchMedia("(max-width: 820px)").matches,
  );
  const location = useLocation();
  const mainRef = useRef<HTMLElement>(null);
  const menuButtonRef = useRef<HTMLButtonElement>(null);
  const sidebarRef = useRef<HTMLElement>(null);
  useEffect(() => {
    mainRef.current?.focus();
  }, [location.pathname]);
  useEffect(() => {
    if (!open) return;
    sidebarRef.current?.querySelector<HTMLElement>("a")?.focus();
    const close = (event: KeyboardEvent) => {
      if (event.key === "Escape") {
        setOpen(false);
        menuButtonRef.current?.focus();
      }
    };
    document.addEventListener("keydown", close);
    return () => document.removeEventListener("keydown", close);
  }, [open]);
  useEffect(() => {
    const query = window.matchMedia("(max-width: 820px)");
    const update = () => setMobile(query.matches);
    query.addEventListener("change", update);
    return () => query.removeEventListener("change", update);
  }, []);
  return (
    <div className="app-shell">
      <a className="skip-link" href="#main-content">
        Skip to content
      </a>
      <header className="mobile-bar">
        <NavLink className="brand" to="/" translate="no">
          PayFlow
        </NavLink>
        <button
          ref={menuButtonRef}
          className="icon-button"
          aria-label={open ? "Close navigation" : "Open navigation"}
          aria-expanded={open}
          aria-controls="primary-sidebar"
          onClick={() => setOpen((v) => !v)}
        >
          {open ? <X aria-hidden="true" /> : <ListBullets aria-hidden="true" />}
        </button>
      </header>
      <aside
        ref={sidebarRef}
        id="primary-sidebar"
        aria-label="Application navigation"
        aria-hidden={mobile && !open}
        inert={mobile && !open ? true : undefined}
        className={`sidebar ${open ? "is-open" : ""}`}
      >
        <NavLink className="brand" to="/" translate="no">
          Pay<span>Flow</span>
        </NavLink>
        <p className="sidebar-label">Payment operations</p>
        <nav aria-label="Primary navigation">
          {items.map(({ to, label, icon: Icon }) => (
            <NavLink
              key={to}
              end={to === "/"}
              to={to}
              onClick={() => setOpen(false)}
            >
              <Icon size={18} aria-hidden="true" />
              {label}
            </NavLink>
          ))}
        </nav>
        <div className="demo-note">
          <strong>Simulation only</strong>
          <span>No card data is collected.</span>
        </div>
      </aside>
      {open && (
        <button
          className="sidebar-scrim"
          aria-label="Close navigation"
          onClick={() => setOpen(false)}
        />
      )}
      <main
        ref={mainRef}
        id="main-content"
        tabIndex={-1}
        inert={open ? true : undefined}
      >
        <Outlet />
      </main>
    </div>
  );
}
