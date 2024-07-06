CREATE TYPE "Status" AS ENUM (
  'ACTIVE',
  'INACTIVE',
  'PENDING',
  'APPROVED',
  'REJECTED',
  'UNDER_REVIEW',
  'FLAGGED',
  'NEEDS_ATTENTION',
  'COMPLETED',
  'ON_HOLD'
);

CREATE TABLE "Entity" (
  "eid" varchar PRIMARY KEY,
  "name" varchar,
  "url" varchar,
  "contact" varchar,
  "employees" varchar,
  "brands" number,
  "remarks" varchar,
  "createdAt" timestamp,
  "modifiedAt" timestamp DEFAULT (now())
);

CREATE TABLE "Brand" (
  "bid" varchar PRIMARY KEY,
  "eid" varchar,
  "name" varchar,
  "url" varchar,
  "contact" varchar,
  "status" Status,
  "remarks" varchar,
  "createdAt" timestamp,
  "modifiedAt" timestamp DEFAULT (now())
);

CREATE TABLE "Location" (
  "lid" varchar PRIMARY KEY,
  "bid" varchar,
  "name" varchar,
  "address" varchar,
  "city" varchar,
  "state" varchar,
  "zip" varchar,
  "contact" varchar,
  "employees" varchar,
  "createdAt" timestamp,
  "modifiedAt" timestamp DEFAULT (now())
);

CREATE TABLE "Employee" (
  "uid" varchar PRIMARY KEY,
  "bid" varchar,
  "firstName" varchar,
  "lastName" varchar,
  "contact" varchar,
  "position" varchar,
  "gender" varchar,
  "marital" string,
  "metrics" varchar,
  "status" Status,
  "documents" varchar,
  "accounts" varchar,
  "createdAt" timestamp,
  "updatedAt" timestamp DEFAULT (now())
);

CREATE TABLE "Documents" (
  "did" varchar PRIMARY KEY,
  "pii" varchar,
  "dlid" varchar,
  "w4id" varchar,
  "esid" varchar,
  "status" Status,
  "createdAt" timestamp,
  "updatedAt" timestamp DEFAULT (now())
);

CREATE TABLE "PII" (
  "piid" varchar PRIMARY KEY,
  "ssn" varchar,
  "dob" varchar,
  "status" Status,
  "createdAt" timestamp,
  "modifiedAt" timestamp DEFAULT (now())
);

CREATE TABLE "Contact" (
  "cid" varchar PRIMARY KEY,
  "street" string,
  "street2" string,
  "city" string,
  "state" string,
  "zip" string,
  "phoneNumber" string,
  "email" string,
  "status" Status,
  "createdAt" timestamp,
  "modifiedAt" timestamp DEFAULT (now())
);

CREATE TABLE "DriversLicense" (
  "dlid" varchar PRIMARY KEY,
  "image" svg,
  "url" varchar,
  "status" Status,
  "createdAt" timestamp,
  "modifiedAt" timestamp DEFAULT (now())
);

CREATE TABLE "ESig" (
  "esid" string PRIMARY KEY,
  "primary" varchar,
  "secondary" varchar,
  "url" varchar,
  "status" Status,
  "createdAt" timestamp,
  "modifiedAt" timestamp DEFAULT (now())
);

CREATE TABLE "W4" (
  "w4id" string PRIMARY KEY,
  "json" json,
  "url" varchar,
  "status" Status,
  "createdAt" timestamp,
  "modifiedAt" timestamp DEFAULT (now())
);

CREATE TABLE "Metrics" (
  "mid" varchar PRIMARY KEY,
  "values" varchar
);

ALTER TABLE "Entity" ADD FOREIGN KEY ("contact") REFERENCES "Contact" ("cid");

ALTER TABLE "Entity" ADD FOREIGN KEY ("employees") REFERENCES "Employee" ("uid");

ALTER TABLE "Entity" ADD FOREIGN KEY ("brands") REFERENCES "Brand" ("bid");

ALTER TABLE "Entity" ADD FOREIGN KEY ("eid") REFERENCES "Brand" ("eid");

ALTER TABLE "Brand" ADD FOREIGN KEY ("contact") REFERENCES "Contact" ("cid");

ALTER TABLE "Brand" ADD FOREIGN KEY ("bid") REFERENCES "Location" ("bid");

ALTER TABLE "Location" ADD FOREIGN KEY ("contact") REFERENCES "Contact" ("cid");

ALTER TABLE "Location" ADD FOREIGN KEY ("employees") REFERENCES "Employee" ("uid");

ALTER TABLE "Brand" ADD FOREIGN KEY ("bid") REFERENCES "Employee" ("bid");

ALTER TABLE "Employee" ADD FOREIGN KEY ("contact") REFERENCES "Contact" ("cid");

ALTER TABLE "Metrics" ADD FOREIGN KEY ("mid") REFERENCES "Employee" ("metrics");

ALTER TABLE "Employee" ADD FOREIGN KEY ("documents") REFERENCES "Documents" ("did");

ALTER TABLE "Documents" ADD FOREIGN KEY ("pii") REFERENCES "PII" ("piid");

ALTER TABLE "Documents" ADD FOREIGN KEY ("dlid") REFERENCES "DriversLicense" ("dlid");

ALTER TABLE "Documents" ADD FOREIGN KEY ("w4id") REFERENCES "W4" ("w4id");

ALTER TABLE "Documents" ADD FOREIGN KEY ("esid") REFERENCES "ESig" ("esid");
