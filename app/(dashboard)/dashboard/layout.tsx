import { HeaderBar } from "@/components/layout/header";
import { SideBar } from "@/components/layout/sidebar";
import type { Metadata } from "next";
import React from "react";

export const metadata: Metadata = {
  title: "",
  description: "",
};

interface DashboardLayoutProps {
  children: React.ReactNode;
}
export default function DashboardLayout({ children }: DashboardLayoutProps) {
  return (
    <>
      <HeaderBar />
      <div className="flex h-screen overflow-hidden">
        <SideBar />
        <main className="flex-1 overflow-hidden pt-16">{children}</main>
      </div>
    </>
  );
}
