import { render, screen } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import { ErrorState, LoadingRows } from "./AsyncState";
test("announces loading state", () => {
  render(<LoadingRows />);
  expect(screen.getByLabelText("Loading…")).toHaveAttribute(
    "aria-busy",
    "true",
  );
});
test("offers recovery from errors", async () => {
  const retry = vi.fn();
  render(<ErrorState message="Network unavailable." retry={retry} />);
  await userEvent.click(screen.getByRole("button", { name: "Retry" }));
  expect(retry).toHaveBeenCalledOnce();
});
