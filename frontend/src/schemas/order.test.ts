import { orderSchema } from "./order";
test("rejects invalid order input", () => {
  const result = orderSchema.safeParse({
    customerName: "A",
    customerEmail: "bad",
    merchantName: "M",
    amount: 0,
    currency: "EUR",
    installmentCount: 6,
  });
  expect(result.success).toBe(false);
});
test("accepts supported plan", () => {
  const result = orderSchema.safeParse({
    customerName: "Maya Rivera",
    customerEmail: "maya@example.test",
    merchantName: "North Harbor Goods",
    amount: 10000,
    currency: "USD",
    installmentCount: 3,
  });
  expect(result.success).toBe(true);
});
