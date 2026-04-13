import { z } from "zod";

export const userSchema = z.object({
  _id: z.string().optional(),
  displayName: z.string().min(1),
  email: z.email(),
  password: z.string().min(6).optional(), // Optional for some UI contexts but required in model
  isAdmin: z.boolean().default(false),
  createdAt: z.iso.datetime().optional(),
  updatedAt: z.iso.datetime().optional(),
});

export type User = z.infer<typeof userSchema>;
