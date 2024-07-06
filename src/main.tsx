import React from "react";
import ReactDOM from "react-dom/client";
import { RouterProvider } from "react-router-dom";
import { router } from "@/router";
import { ThemeProvider } from "@/components/theme-provider";
import "@/index.css";
import i18n from "@/i18n";
import { I18nextProvider } from "react-i18next";
import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { MockProvider } from "./lib/team-context-provider";

const queryClient = new QueryClient();
ReactDOM.createRoot(document.getElementById("root")!).render(
  <React.StrictMode>
    <MockProvider>
      <I18nextProvider i18n={i18n}>
        <QueryClientProvider client={queryClient}>
          <ThemeProvider defaultTheme='dark' storageKey='vite-ui-theme'>
            <RouterProvider router={router} />
          </ThemeProvider>
        </QueryClientProvider>
      </I18nextProvider>
    </MockProvider>
  </React.StrictMode>
);
