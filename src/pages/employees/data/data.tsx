import { Status } from "@/lib/team-context-provider";
import {
  ArrowDownIcon,
  ArrowRightIcon,
  ArrowUpIcon,
  CheckCircledIcon,
  CircleIcon,
  CrossCircledIcon,
  QuestionMarkCircledIcon,
  StopwatchIcon,
  ExclamationTriangleIcon,
} from "@radix-ui/react-icons";
// import { IconExclamationCircle } from "@tabler/icons-react";

export const labels = [
  {
    value: "bug",
    label: "Bug",
  },
  {
    value: "feature",
    label: "Feature",
  },
  {
    value: "documentation",
    label: "Documentation",
  },
];

export const statuses = [
  {
    value: "ACTIVE",
    label: "ACTIVE",
    icon: QuestionMarkCircledIcon,
  },
  {
    value: "PENDING",
    label: "PENDING",
    icon: CircleIcon,
  },
  {
    value: "APPROVED",
    label: "APPROVED",
    icon: StopwatchIcon,
  },
  {
    value: "REJECTED",
    label: "REJECTED",
    icon: CrossCircledIcon,
  },
  {
    value: "FLAGGED",
    label: "FLAGGED",
    icon: CrossCircledIcon,
  },
];

export const priorities = [
  {
    label: "Low",
    value: "low",
    icon: ArrowDownIcon,
  },
  {
    label: "Medium",
    value: "medium",
    icon: ArrowRightIcon,
  },
  {
    label: "High",
    value: "high",
    icon: ArrowUpIcon,
  },
];

export const applicationStatus = [
  {
    label: "done",
    value: "done",
    icon: CheckCircledIcon,
  },
  {
    label: "in progress",
    value: "in progress",
    icon: StopwatchIcon,
  },
  {
    value: "canceled",
    label: "Canceled",
    icon: CrossCircledIcon,
  },
  {
    value: "needs attention",
    label: "needs attention",
    icon: ExclamationTriangleIcon,
  },
];
