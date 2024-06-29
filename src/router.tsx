import { Suspense, lazy } from "react"
import { createBrowserRouter, Route, Routes } from "react-router-dom"
import Loader from "@/components/loader"

export const router = createBrowserRouter([
  {
    path: "/",
    lazy: async () => {
      const RootLayout = await import("./components/root-layout")
      return { Component: RootLayout.default }
    },
    children: [
      {
        index: true,
        lazy: async () => ({
          Component: (await import("@/pages/dashboard")).default,
        }),
      },
      {
        path: "tasks",
        lazy: async () => ({
          Component: (await import("@/pages/tasks")).default,
        }),
      },
      {
        path: "documents",
        lazy: async () => ({
          Component: (await import("@/pages/documents")).default,
        }),
      },
      {
        path: "employees",
        lazy: async () => ({
          Component: (await import("@/pages/employees")).default,
        }),
      },
      {
        path: "kanban",
        lazy: async () => ({
          Component: (await import("@/pages/kanban")).default,
        }),
      },
    ],
  },
])
