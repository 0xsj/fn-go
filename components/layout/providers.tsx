"use client";
import { ThemeProvider } from "next-themes";
import React from "react";

interface Props {
  children: React.ReactNode;
}
export const Providers: React.FC<Props> = ({ children }) => {
  return (
    <>
      <ThemeProvider>{children}</ThemeProvider>
    </>
  );
};
