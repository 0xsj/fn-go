import { IconMoon, IconSun } from "@tabler/icons-react";
import { useTheme } from "@/hooks/use-theme";
import { Button } from "./button";
import { useEffect } from "react";

export const ThemeSwitch = () => {
  const { theme, setTheme } = useTheme();
  useEffect(() => {
    const themeColor = theme === "dark" ? "#020817" : "#fff";
    const metaThemeColor = document.querySelector("meta[name='theme-color']");
    metaThemeColor && metaThemeColor.setAttribute("content", themeColor);
  }, [theme]);

  return (
    <Button
      size={"icon"}
      variant={"ghost"}
      className='rounded-full'
      onClick={() => setTheme(theme === "light" ? "dark" : "light")}
    >
      <IconMoon />
    </Button>
  );
};
