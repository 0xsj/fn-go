// Base Types
export type BaseModel = {
  id: string;
  createdAt: string;
  updatedAt: string;
};

export type BaseUser = BaseModel & {
  name: string;
  address: Address;
  email: string;
  phone: string;
};

export type Address = {
  street: string;
  city: string;
  state: string;
  zip: string;
  country?: string; // Optional country field
};

export type SignatureObject = {
  svg: string; // Assuming SVG string for signature
  type?: string; // Optional type field for future formats
};

// Application Form Model
export type ApplicationForm = BaseModel & {
  fullName: string; // first, last - fillable
  position: string; // could be an enum - fillable
  address: Address; // fillable
  startDate: string; // fillable
  salaryDesired: string; // fillable
  contactInfo: string; // fillable
  legalStatus: boolean; // fillable - Are you a green card holder? No visa? Citizen?
  dob: string; // fillable - Are you 21 years or older?
  schedule: Schedule; // fillable - Availability schedule
  workExperience: JobItem[]; // fillable - Previous work experiences
  other: LegalItems; // fillable - Legal information
  signature: SignatureObject; // fillable - User's signature
};

export type Schedule = {
  days: {
    [key: string]: {
      type: ScheduleType;
      remarks: string;
    };
  };
};

export type JobItem = {
  employer: string; // name of the employer - fillable
  position: string; // previous position - fillable
  address: Address; // fillable
  remarks: string; // reason for leaving - fillable
  startDate: string; // fillable
  endDate: string; // fillable
};

export type LegalItems = {
  field1: boolean; // if yes, explain - fillable
  explanation1?: string; // explanation for field1 - fillable if field1 is true
  field2: boolean; // if yes, explain - fillable
  explanation2?: string; // explanation for field2 - fillable if field2 is true
  field3: boolean; // if yes, explain - fillable
  explanation3?: string; // explanation for field3 - fillable if field3 is true
};

export enum ScheduleType {
  LUNCH = "lunch",
  DINNER = "dinner",
}

// W4 Form Model
export type W4Form = BaseUser & {
  marital: MaritalStatus; // fillable enum
  step3a: number; // algorithmically calculated
  step3b: number; // algorithmically calculated
  step3c: number; // step3a + step3b
  step4a: number; // fillable
  step4b: number; // fillable
  step4c: number; // fillable
  signature: SignatureObject; // user signature
  entity: Entity; // store entity information
};

export type Entity = {
  storeName: string;
  address: Address;
  employees: Employee[];
};

export type Employee = BaseUser & {
  position: string;
  social: number;
  gender: string;
  maritalStatus: MaritalStatus;
};

export type Handbook = BaseModel & {
  title: string;
  sections: HandbookSection[];
};

export type HandbookSection = {
  title: string;
  content: string;
  subsections?: HandbookSection[];
};

export enum MaritalStatus {
  SINGLE = "single",
  MARRIED = "married",
}

export enum JobPosition {
  MANAGER = "manager",
  SERVER = "server",
  CHEF = "chef",
  BACKLINE = "backline",
  FOH = "foh",
}
