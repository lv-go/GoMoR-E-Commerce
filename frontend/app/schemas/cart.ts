import z from "zod";

export const cartItemSchema = z.object({
  _id: z.string().min(1),
  name: z.string().min(1),
  price: z.number().min(0),
  image: z.string().min(1),
  countInStock: z.number().min(0),
  quantity: z.number().min(1),
});

export type CartItem = z.infer<typeof cartItemSchema>;

export const cartSchema = z.object({
  cartItems: z.array(
    cartItemSchema
  ),
  shippingAddress: z.object({
    address: z.string().min(1),
    city: z.string().min(1),
    postalCode: z.string().min(1),
    country: z.string().min(1),
  }),
  paymentMethod: z.string().min(1),
  itemsPrice: z.number().min(1),
  shippingPrice: z.number().min(1),
  taxPrice: z.number().min(1),
  totalPrice: z.number().min(1),
});

export type Cart = z.infer<typeof cartSchema>;

export const newCart = (): Cart => ({
  cartItems: [],
  shippingAddress: {
    address: "",
    city: "",
    postalCode: "",
    country: "",
  },
  paymentMethod: "PayPal",
  itemsPrice: 0,
  shippingPrice: 0,
  taxPrice: 0,
  totalPrice: 0,
});
