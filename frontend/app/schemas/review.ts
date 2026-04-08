import z from "zod";

export const reviewSchema = z.object({
  _id: z.string().optional(),
  name: z.string().min(1),
  rating: z.number().min(0).max(5),
  comment: z.string().min(1),
  user: z.string(), // ObjectId as string
  createdAt: z.iso.datetime().optional(),
  updatedAt: z.iso.datetime().optional(),
});

export type Review = z.infer<typeof reviewSchema>;
