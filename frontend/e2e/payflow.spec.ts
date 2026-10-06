import { expect, test } from "@playwright/test";
test("creates an order, pays an installment, and requests a refund", async ({
  page,
}) => {
  await page.goto("/orders/new");
  await page.getByLabel("Name").fill("Maya Rivera");
  await page.getByLabel("Email").fill(`maya-${Date.now()}@example.test`);
  await page.getByLabel("Merchant").fill("North Harbor Goods");
  await page.getByLabel("Amount in minor units").fill("12000");
  await page.getByLabel("Installments").selectOption("4");
  await page.getByRole("button", { name: "Create Order" }).click();
  await expect(page.getByRole("heading", { name: /PF-/ })).toBeVisible();
  await page.getByRole("button", { name: "Pay" }).first().click();
  await expect(page.getByText("Payment succeeded")).toBeVisible();
  page.once("dialog", (dialog) => dialog.accept());
  await page.getByRole("button", { name: "Refund Payment" }).click();
  await expect(page.getByText("Refund Payment")).not.toBeVisible();
});
