import { z } from "zod";

export const orderItemSchema = z.object({
  name: z.string().min(1),
  quantity: z.number().int().min(1),
  image: z.string().min(1),
  price: z.number().min(0),
  productId: z.string(), // ObjectId as string
});

export const shippingAddressSchema = z.object({
  address: z.string().min(1),
  city: z.string().min(1),
  postalCode: z.string().min(1),
  country: z.string().min(1),
});

export const paymentResultSchema = z.object({
  id: z.string().optional(),
  status: z.string().optional(),
  update_time: z.string().optional(),
  email_address: z.string().email().optional(),
});

export const orderSchema = z.object({
  _id: z.string().optional(),
  userId: z.string(),
  orderItems: z.array(orderItemSchema),
  shippingAddress: shippingAddressSchema,
  paymentMethod: z.string().min(1),
  paymentResult: paymentResultSchema.optional(),
  itemsPrice: z.number().min(0).default(0),
  taxPrice: z.number().min(0).default(0),
  shippingPrice: z.number().min(0).default(0),
  totalPrice: z.number().min(0).default(0),
  isPaid: z.boolean().default(false),
  paidAt: z.string().datetime().optional(),
  isDelivered: z.boolean().default(false),
  deliveredAt: z.string().datetime().optional(),
  createdAt: z.iso.datetime().optional(),
  updatedAt: z.iso.datetime().optional(),
});

export type OrderItem = z.infer<typeof orderItemSchema>;
export type ShippingAddress = z.infer<typeof shippingAddressSchema>;
export type PaymentResult = z.infer<typeof paymentResultSchema>;
export type Order = z.infer<typeof orderSchema>;
