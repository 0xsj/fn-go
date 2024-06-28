"use client";

import { cn } from "@/lib/utils";
import { ChevronLeftIcon } from "@radix-ui/react-icons";
import { SidebarNav } from "../sidebar-nav";
import { navItems } from "@/constants/data";

type Props = {
  className?: string;
};
export const SideBar: React.FC<Props> = ({ className, ...props }) => {
  return (
    <nav
      {...props}
      className={cn(
        `relative hidden h-screen flex-none border-r z-10 pt-20 md:block`,
        className,
      )}
    >
      <ChevronLeftIcon
        className={cn(
          `absolute-right-3 top-20 cursor-pointer rounded-full border bg-background text-3xl text-foreground`,
        )}
      />
      <div className="space-y-4 py-4">
        <div className="px-3 py-2">
          <div className="mt-3 space-y-1">
            <SidebarNav items={navItems} />
          </div>
        </div>
      </div>
    </nav>
  );
};
