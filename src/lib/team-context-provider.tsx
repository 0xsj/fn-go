import React, {
  createContext,
  useState,
  useContext,
  HTMLAttributes,
  Provider,
  ReactNode,
} from "react";
import {
  bcd,
  yardHouse,
  entities,
  employees,
  kopan,
  brands,
} from "@/data/data";

interface MockContextType {
  mockData: MockData;
  updatedMockData: (upDatedData: MockData) => void;
}

const MockContext = createContext<MockContextType | undefined>(undefined);

interface ProviderProps {
  children: ReactNode;
}

export const MockProvider: React.FC<ProviderProps> = ({ children }) => {
  const initialMockData: MockData = {
    entities: entities,
    brands: brands,
    locations: {
      kopan: kopan,
      bcd: bcd,
      yardHouse: yardHouse,
    },
    employees: employees,
  };

  const [mockData, setMockData] = useState<MockData>(initialMockData);

  const updatedMockData = (updatedData: MockData) => {
    setMockData(updatedData);
  };

  return (
    <MockContext.Provider value={{ mockData, updatedMockData }}>
      {children}
    </MockContext.Provider>
  );
};

// eslint-disable-next-line react-refresh/only-export-components
export const useMockContext = (): MockContextType => {
  const context = useContext(MockContext);
  if (!context) {
    throw new Error("useMockContext must be used within a MockProvider");
  }
  return context;
};

// Define Status ENUM
export type Status =
  | "ACTIVE"
  | "INACTIVE"
  | "PENDING"
  | "APPROVED"
  | "REJECTED"
  | "UNDER_REVIEW"
  | "FLAGGED"
  | "NEEDS_ATTENTION"
  | "COMPLETED"
  | "ON_HOLD";

// Interface for Entity table
export interface Entity {
  eid: string;
  name: string;
  url?: string;
  contact: string; // Reference to Contact.cid
  employees: Employee[]; // Reference to Employee.uid
  brands: Brand[]; // Reference to Brand.bid
  remarks?: string;
  createdAt: Date | string;
  modifiedAt: Date | string;
}

// Interface for Brand table
export interface Brand {
  bid: string;
  eid: string; // Reference to Entity.eid
  name: string;
  url?: string;
  contact: string;
  status: Status;
  remarks?: string;
  createdAt: Date | string;
  modifiedAt: Date | string;
}

// Interface for Location table
export interface Location {
  lid: string;
  bid: string; // Reference to Brand.bid
  name: string;
  address?: string;
  city?: string;
  state?: string;
  zip?: string;
  contact: string; // Reference to Contact.cid
  employees: Employee[]; // Reference to Employee.uid
  createdAt: Date | string;
  modifiedAt: Date | string;
}

// Interface for Employee table
export interface Employee {
  uid: string;
  bid: string; // Reference to Brand.bid
  firstName: string;
  lastName: string;
  contact: string; // Reference to Contact.cid
  position?: string;
  gender?: string;
  marital?: string;
  metrics?: string; // Reference to Metrics.mid
  status: Status;
  documents: string[]; // Reference to Documents.did
  accounts?: string[];
  createdAt: Date | string;
  modifiedAt: Date | string;
}

// Interface for Documents table
export interface Documents {
  did: string;
  pii?: string; // Reference to PII.piid
  dlid?: string; // Reference to DriversLicense.dlid
  w4id?: string; // Reference to W4.w4id
  esid?: string; // Reference to ESig.esid
  status: Status;
  createdAt: Date | string;
  modifiedAt: Date | string;
}

// Interface for PII table
export interface PII {
  piid: string;
  ssn?: string;
  dob?: string;
  status: Status;
  createdAt: Date | string;
  modifiedAt: Date | string;
}

// Interface for Contact table
export interface Contact {
  cid: string;
  street?: string;
  street2?: string;
  city?: string;
  state?: string;
  zip?: string;
  phoneNumber?: string;
  email?: string;
  status: Status;
  createdAt: Date | string;
  modifiedAt: Date | string;
}

// Interface for DriversLicense table
export interface DriversLicense {
  dlid: string;
  image?: string; // SVG image data
  url?: string;
  status: Status;
  createdAt: Date | string;
  modifiedAt: Date | string;
}

// Interface for ESig table
export interface ESig {
  esid: string;
  primary?: string;
  secondary?: string;
  url?: string;
  status: Status;
  createdAt: Date | string;
  modifiedAt: Date | string;
}

// Interface for W4 table
export interface W4 {
  w4id: string;
  json?: object; // JSON data
  url?: string;
  status: Status;
  createdAt: Date | string;
  modifiedAt: Date | string;
}

// Interface for Metrics table
export interface Metrics {
  mid: string;
  values?: string;
}

export interface MockData {
  entities: Entity[];
  brands: Brand[];
  locations: {
    kopan: Location[];
    bcd: Location[];
    yardHouse: Location[];
  };
  employees: Employee[];
}
