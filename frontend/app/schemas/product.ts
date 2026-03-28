import { z } from "zod";

export const reviewSchema = z.object({
  _id: z.string().optional(),
  name: z.string().min(1),
  rating: z.number().min(0).max(5),
  comment: z.string().min(1),
  user: z.string(), // ObjectId as string
  createdAt: z.string().datetime().optional(),
  updatedAt: z.string().datetime().optional(),
});

export const productSchema = z.object({
  _id: z.string().optional(),
  name: z.string().min(1),
  image: z.string().min(1),
  brand: z.string().min(1),
  quantity: z.number().int().min(0),
  category: z.string(), // ObjectId as string
  description: z.string().min(1),
  reviews: z.array(reviewSchema).default([]),
  rating: z.number().min(0).max(5).default(0),
  numReviews: z.number().int().min(0).default(0),
  price: z.number().min(0).default(0),
  countInStock: z.number().int().min(0).default(0),
  createdAt: z.string().datetime().optional(),
  updatedAt: z.string().datetime().optional(),
});

export type Review = z.infer<typeof reviewSchema>;
export type Product = z.infer<typeof productSchema>;
