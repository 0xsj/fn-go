import { z } from "zod";

// Enum type for Status
export const StatusSchema = z.enum([
  "ACTIVE",
  "INACTIVE",
  "PENDING",
  "APPROVED",
  "REJECTED",
  "UNDER_REVIEW",
  "FLAGGED",
  "NEEDS_ATTENTION",
  "COMPLETED",
  "ON_HOLD",
]);

// Schema for Contact
export const ContactSchema = z.object({
  cid: z.string(),
  street: z.string().nullable(),
  street2: z.string().nullable(),
  city: z.string().nullable(),
  state: z.string().nullable(),
  zip: z.string().nullable(),
  phoneNumber: z.string().nullable(),
  email: z.string().nullable(),
  status: StatusSchema,
  createdAt: z.date(),
  modifiedAt: z.date().default(new Date()),
});

// Schema for Entity
export const EntitySchema = z.object({
  eid: z.string(),
  name: z.string().nullable(),
  url: z.string().nullable(),
  contact: ContactSchema.optional().nullable(), // Relation to Contact
  employees: z.array(z.string()).nullable(), // List of Employee IDs
  brands: z.array(z.string()).nullable(), // List of Brand IDs
  remarks: z.string().nullable(),
  createdAt: z.date(),
  modifiedAt: z.date().default(new Date()),
});

// Schema for Brand
export const Brand = z.object({
  bid: z.string(),
  eid: z.string().nullable(), // Foreign Key to Entity
  name: z.string().nullable(),
  url: z.string().nullable(),
  contact: ContactSchema.optional().nullable(), // Relation to Contact
  employees: z.array(z.string()).nullable(), // List of Employee IDs
  status: StatusSchema,
  remarks: z.string().nullable(),
  createdAt: z.date(),
  modifiedAt: z.date().default(new Date()),
});

// Schema for Employee
export const Employee = z.object({
  uid: z.string(),
  bid: z.string().nullable(), // Foreign Key to Brand
  firstName: z.string().nullable(),
  lastName: z.string().nullable(),
  contact: ContactSchema.optional().nullable(), // Relation to Contact
  position: z.string().nullable(),
  gender: z.string().nullable(),
  marital: z.string().nullable(),
  metrics: z.string().nullable(),
  status: StatusSchema,
  documents: z.array(z.string()).nullable(), // List of Document IDs
  accounts: z.string().nullable(),
  createdAt: z.date(),
  updatedAt: z.date().default(new Date()),
});

// Schema for Documents
export const Documents = z.object({
  did: z.string(),
  pii: z.string().nullable(), // Foreign Key to PII
  dlid: z.string().nullable(), // Foreign Key to DriversLicense
  w4id: z.string().nullable(), // Foreign Key to W4
  esid: z.string().nullable(), // Foreign Key to ESig
  status: StatusSchema,
  createdAt: z.date(),
  updatedAt: z.date().default(new Date()),
});

// Schema for PII
export const PII = z.object({
  piid: z.string(),
  ssn: z.string().nullable(),
  dob: z.string().nullable(),
  status: StatusSchema,
  createdAt: z.date(),
  modifiedAt: z.date().default(new Date()),
});

// Schema for DriversLicense
export const DriversLicense = z.object({
  dlid: z.string(),
  image: z.string().nullable(), // Assuming svg as string for simplicity
  url: z.string().nullable(),
  status: StatusSchema,
  createdAt: z.date(),
  modifiedAt: z.date().default(new Date()),
});

// Schema for ESig
export const ESig = z.object({
  esid: z.string(),
  primary: z.string().nullable(),
  secondary: z.string().nullable(),
  url: z.string().nullable(),
  status: StatusSchema,
  createdAt: z.date(),
  modifiedAt: z.date().default(new Date()),
});

// Schema for W4
export const W4 = z.object({
  w4id: z.string(),
  json: z.any(), // Assuming JSON as any type
  url: z.string().nullable(),
  status: StatusSchema,
  createdAt: z.date(),
  modifiedAt: z.date().default(new Date()),
});

// Schema for Metrics
export const Metrics = z.object({
  mid: z.string(),
  values: z.string().nullable(),
});
