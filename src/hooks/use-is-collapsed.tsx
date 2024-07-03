import { useEffect } from "react";
import { useLocalStorage } from "./use-local-storage";

export const useIsCollapsed = () => {
  const [isCollapsed, setIsCollapsed] = useLocalStorage({
    key: "collapsed-sidebar",
    defaultValue: false,
  });

  useEffect(() => {
    const handleResize = () => {
      setIsCollapsed(window.innerWidth < 768 ? false : isCollapsed);
    };
    handleResize();

    window.addEventListener("resize", handleResize);
    return () => {
      window.removeEventListener("resize", handleResize);
    };
  }, [isCollapsed, setIsCollapsed]);
  return [isCollapsed, setIsCollapsed] as const;
};
