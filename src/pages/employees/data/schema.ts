import { z } from "zod";

export const taskSchema = z.object({
  id: z.string(),
  title: z.string(),
  status: z.string(),
  label: z.string(),
  priority: z.string(),
});

export const userSchema = z.object({
  uuidv4: z.string().optional(),
  firstName: z.string().optional(),
  lastName: z.string().optional(),
  maritalStatus: z.string().optional(),
  phone: z.string().optional(),
  email: z.string().optional(),
  ssn: z.string().optional(),
  dob: z.string().optional(),
  street: z.string().optional(),
  street2: z.string().optional(),
  city: z.string().optional(),
  state: z.string().optional(),
  zip: z.string().optional(),
  company: z.string().optional(),
  title: z.string().optional(),
  dateCreated: z.string().optional(),
  dateModified: z.string().optional(),
});

export type Task = z.infer<typeof taskSchema>;
export type User = z.infer<typeof userSchema>;
