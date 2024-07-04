/* eslint-disable react-refresh/only-export-components */
import { generateMock } from "@anatine/zod-mock";
import {
  BaseModel,
  BaseUser,
  Address,
  SignatureObject,
  ApplicationForm,
  Schedule,
  JobItem,
  LegalItems,
  ScheduleType,
  W4Form,
  Entity,
  Employee,
  Handbook,
  HandbookSection,
  MaritalStatus,
  JobPosition,
} from "@/types/types";
import { string, z } from "zod";

// BaseModel Schema
const BaseModelSchema = z.object({
  id: z.string().uuid(),
  createdAt: z.string().datetime(),
  updatedAt: z.string().datetime(),
});

// Address Schema
const AddressSchema = z.object({
  street: z.string().min(1).max(50),
  city: z.string().min(1).max(50),
  state: z.string().min(2).max(50),
  zip: z.string().min(5).max(10).regex(/^\d+$/, "Invalid ZIP code"),
  country: z.string().min(2).max(50).optional(),
});

// SignatureObject Schema
const SignatureObjectSchema = z.object({
  svg: z.string().min(1),
  type: z.string().min(1).optional(),
});

// BaseUser Schema
const BaseUserSchema = BaseModelSchema.extend({
  name: z.string().min(1).max(100),
  address: AddressSchema,
  email: z.string().email().max(100),
  phone: z
    .string()
    .min(10)
    .max(15)
    .regex(/^\+?\d{10,15}$/, "Invalid phone number"),
});

// ScheduleType Schema
const ScheduleTypeSchema = z.enum(["lunch", "dinner"]);

// Schedule Schema
const ScheduleSchema = z.object({
  days: z.record(
    z.object({
      type: ScheduleTypeSchema,
      remarks: z.string().max(200).optional(),
    })
  ),
});

// JobItem Schema
const JobItemSchema = z.object({
  employer: z.string().min(1).max(50),
  position: z.string().min(1).max(50),
  address: AddressSchema,
  remarks: z.string().max(50),
  startDate: z.string().max(50),
  endDate: z.string().max(50),
});

// LegalItems Schema
const LegalItemsSchema = z.object({
  field1: z.boolean(),
  explanation1: z.string().max(100).optional(),
  field2: z.boolean(),
  explanation2: z.string().max(100).optional(),
  field3: z.boolean(),
  explanation3: z.string().max(100).optional(),
});

// MaritalStatus Schema
const MaritalStatusSchema = z.enum(["single", "married"]);

// Entity Schema
const EntitySchema = z.object({
  storeName: z.string().min(1).max(100),
  address: AddressSchema,
  employees: z.array(
    BaseUserSchema.extend({
      position: z.string().min(1).max(50),
      social: z
        .string()
        .min(9)
        .max(11)
        .regex(/^\d{3}-\d{2}-\d{4}$/, "Invalid social security number"), // Updated to reflect valid SSN format
      gender: z.string().max(10),
      maritalStatus: MaritalStatusSchema,
    })
  ),
});

// HandbookSection Schema
const HandbookSectionSchema: z.ZodSchema<HandbookSection> = z.lazy(() =>
  z.object({
    title: z.string().max(50),
    content: z.string(),
    subsections: z.array(HandbookSectionSchema).optional(),
  })
);

// Handbook Schema
const HandbookSchema = BaseModelSchema.extend({
  title: z.string().min(1).max(100),
  sections: z.array(HandbookSectionSchema),
});

// W4Form Schema
const W4FormSchema = BaseUserSchema.extend({
  marital: MaritalStatusSchema,
  step3a: z.number().min(0),
  step3b: z.number().min(0),
  step3c: z.number().min(0),
  step4a: z.number().min(0),
  step4b: z.number().min(0),
  step4c: z.number().min(0),
  signature: SignatureObjectSchema,
  entity: EntitySchema,
});

export {
  BaseModelSchema,
  BaseUserSchema,
  AddressSchema,
  SignatureObjectSchema,
  ScheduleSchema,
  JobItemSchema,
  LegalItemsSchema,
  MaritalStatusSchema,
  EntitySchema,
  W4FormSchema,
  HandbookSchema,
  HandbookSectionSchema,
};
