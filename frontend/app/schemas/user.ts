import { z } from "zod";

export const userSchema = z.object({
  _id: z.string().optional(),
  username: z.string().min(1),
  email: z.string().email(),
  password: z.string().min(6).optional(), // Optional for some UI contexts but required in model
  isAdmin: z.boolean().default(false),
  createdAt: z.string().datetime().optional(),
  updatedAt: z.string().datetime().optional(),
});

export type User = z.infer<typeof userSchema>;
