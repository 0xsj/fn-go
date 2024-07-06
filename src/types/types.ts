// Enum Type
export enum Status {
  ACTIVE = "ACTIVE",
  INACTIVE = "INACTIVE",
  PENDING = "PENDING",
  APPROVED = "APPROVED",
  REJECTED = "REJECTED",
  UNDER_REVIEW = "UNDER_REVIEW",
  FLAGGED = "FLAGGED",
  NEEDS_ATTENTION = "NEEDS_ATTENTION",
  COMPLETED = "COMPLETED",
  ON_HOLD = "ON_HOLD",
}

// Entity Table
export interface Entity {
  eid: string;
  name: string;
  url: string;
  contact: string;
  employees: string;
  brands: number;
  remarks: string;
  createdAt: Date;
  modifiedAt: Date;
}

// Brand Table
export interface Brand {
  bid: string;
  eid: string;
  name: string;
  url: string;
  contact: string;
  employees: string;
  status: Status;
  remarks: string;
  createdAt: Date;
  modifiedAt: Date;
}

// Employee Table
export interface Employee {
  uid: string;
  bid: string;
  firstName: string;
  lastName: string;
  contact: string;
  position: string;
  gender: string;
  marital: string;
  metrics: string;
  status: Status;
  documents: string;
  accounts: string;
  createdAt: Date;
  updatedAt: Date;
}

// Documents Table
export interface Documents {
  did: string;
  pii: string;
  dlid: string;
  w4id: string;
  esid: string;
  status: Status;
  createdAt: Date;
  updatedAt: Date;
}

// PII Table
export interface PII {
  piid: string;
  ssn: string;
  dob: string;
  status: Status;
  createdAt: Date;
  modifiedAt: Date;
}

// Contact Table
export interface Contact {
  cid: string;
  street: string;
  street2: string;
  city: string;
  state: string;
  zip: string;
  phoneNumber: string;
  email: string;
  status: Status;
  createdAt: Date;
  modifiedAt: Date;
}

// DriversLicense Table
export interface DriversLicense {
  dlid: string;
  image: string; // Assuming `svg` is stored as a string
  url: string;
  status: Status;
  createdAt: Date;
  modifiedAt: Date;
}

// ESig Table
export interface ESig {
  esid: string;
  primary: string;
  secondary: string;
  url: string;
  status: Status;
  createdAt: Date;
  modifiedAt: Date;
}

// W4 Table
export interface W4 {
  w4id: string;
  json: object; // Assuming `json` is an object
  url: string;
  status: Status;
  createdAt: Date;
  modifiedAt: Date;
}

// Metrics Table
export interface Metrics {
  mid: string;
  values: string;
}
