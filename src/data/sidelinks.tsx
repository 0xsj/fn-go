import {
  IconBoxSeam,
  IconChartHistogram,
  IconLayoutDashboard,
  IconMessages,
  IconRouteAltLeft,
  IconTruck,
  IconUsers,
} from "@tabler/icons-react"

export interface NavLink {
  title: string
  label?: string
  href: string
  icon: JSX.Element
  isActive: boolean
}

export interface SideLink extends NavLink {
  sub?: NavLink[]
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
    icon: <IconMessages size={18} />,
    isActive: false,
  },
  {
    title: "Documents",
    label: "",
    href: "/documents",
    icon: <IconUsers size={18} />,
    isActive: false,
  },
  {
    title: "Tasks",
    label: "",
    href: "/tasks",
    icon: <IconChartHistogram size={18} />,
    isActive: false,
  },
  {
    title: "Kanban",
    label: "",
    href: "/kanban",
    icon: <IconChartHistogram size={18} />,
    isActive: false,
  },
]
