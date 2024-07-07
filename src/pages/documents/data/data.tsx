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
import { IconBook, IconFile, IconPdf } from "@tabler/icons-react";
// import { IconExclamationCircle } from "@tabler/icons-react";

export const dockLabels = [
  {
    value: "w4",
    label: "w4",
    icon: <IconPdf size={20} />,
  },
  {
    value: "application",
    label: "application",
    icon: <IconFile size={20} />,
  },
  {
    value: "handbook",
    label: "handbook",
    icon: <IconBook size={20} />,
  },
];

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
    label: "canceled",
    icon: CrossCircledIcon,
  },
  {
    value: "flagged",
    label: "flagged",
    icon: ExclamationTriangleIcon,
  },
];
