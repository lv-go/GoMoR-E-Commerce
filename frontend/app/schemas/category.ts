import { z } from "zod";

export const categorySchema = z.object({
  _id: z.string().optional(),
  name: z.string().trim().min(1).max(32),
});

export type Category = z.infer<typeof categorySchema>;
