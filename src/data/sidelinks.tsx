import {
  IconBoxSeam,
  IconChartHistogram,
  IconLayoutDashboard,
  IconMessages,
  IconRouteAltLeft,
  IconTruck,
  IconPaperclip,
  IconSubtask,
  IconLayoutKanban,
  IconUsers,
  IconApps,
} from "@tabler/icons-react";

export interface NavLink {
  title: string;
  label?: string;
  href: string;
  icon: JSX.Element;
  isActive: boolean;
}

export interface SideLink extends NavLink {
  sub?: NavLink[];
}

export const sidelinks: SideLink[] = [
  {
    title: "Overview",
    label: "",
    href: "/",
    icon: <IconLayoutDashboard size={18} />,
    isActive: true,
  },
  {
    title: "Employees",
    label: "9",
    href: "/employees",
    icon: <IconUsers size={18} />,
    isActive: false,
  },
  {
    title: "Documents",
    label: "",
    href: "/documents",
    icon: <IconPaperclip size={18} />,
    isActive: false,
  },
  {
    title: "Apps",
    label: "",
    href: "/apps",
    icon: <IconApps size={18} />,
    isActive: false,
  },
  // {
  //   title: "Tasks",
  //   label: "",
  //   href: "/tasks",
  //   icon: <IconSubtask size={18} />,
  //   isActive: false,
  // },
  // {
  //   title: "Kanban",
  //   label: "",
  //   href: "/kanban",
  //   icon: <IconLayoutKanban size={18} />,
  //   isActive: false,
  // },
];
