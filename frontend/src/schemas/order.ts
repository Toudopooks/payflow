import { z } from "zod";
export const orderSchema = z.object({
  customerName: z.string().trim().min(2, "Enter at least 2 characters."),
  customerEmail: z.email("Enter a valid email address."),
  merchantName: z.string().trim().min(2, "Enter at least 2 characters."),
  amount: z.coerce.number().int().positive("Amount must be greater than zero."),
  currency: z.enum(["USD", "CLP"]),
  installmentCount: z.coerce.number().int().min(2).max(4),
});
export type OrderForm = z.input<typeof orderSchema>;
