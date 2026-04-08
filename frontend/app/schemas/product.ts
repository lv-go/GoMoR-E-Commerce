import { z } from "zod";
import { reviewSchema } from "./review";

export const productSchema = z.object({
  _id: z.string().optional(),
  name: z.string().min(1),
  image: z.string().min(1),
  brand: z.string().min(1),
  quantity: z.number().int().min(0),
  categoryId: z.string(), // ObjectId as string
  description: z.string().min(1),
  reviews: z.array(reviewSchema).default([]),
  rating: z.number().min(0).max(5).default(0),
  numReviews: z.number().int().min(0).default(0),
  price: z.number().min(0).default(0),
  countInStock: z.number().int().min(0).default(0),
  createdAt: z.iso.datetime().optional(),
  updatedAt: z.iso.datetime().optional(),
});

export type Product = z.infer<typeof productSchema>;

export const newProduct = (): Product => ({
  name: "",
  image: "",
  brand: "",
  quantity: 0,
  categoryId: "",
  description: "",
  reviews: [],
  rating: 0,
  numReviews: 0,
  price: 0,
  countInStock: 0,
  createdAt: "",
  updatedAt: "",
});

export const getId = (item: Product) => item._id || "";

export type ProductFilters = {
  category?: string;
  price?: string;
  brand?: string;
};
