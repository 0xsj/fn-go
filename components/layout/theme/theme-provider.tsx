"use client";

import { ThemeProvider as NextThemeProvider } from "next-themes";
import { type ThemeProviderProps } from "next-themes/dist/types";

interface Props extends ThemeProviderProps {
  children: React.ReactNode;
}

export const ThemeProvider: React.FC<Props> = ({ children, ...props }) => {
  return (
    <>
      <NextThemeProvider {...props}>{children}</NextThemeProvider>
    </>
  );
};
