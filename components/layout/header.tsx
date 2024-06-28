import { cn } from "@/lib/utils";
// import Link from "next/link";
import { ThemeToggle } from "./theme/theme-toggle";

interface Props {}

export const HeaderBar: React.FC<Props> = () => {
  return (
    <div className="supports-backdrop-blur:bg-background/60 fixed left-0 right-0 top-0 z-20 border-b bg-background/95 backdrop-blur">
      <nav className="flex h-14 items-center justify-between px-4">
        <div className="hidden lg:block"></div>
        <div className={cn("block lg:hidden")}></div>
        <div className="flex items-center gap-2">
          <ThemeToggle />
        </div>
      </nav>
    </div>
  );
};
