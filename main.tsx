import React from 'react'
import ReactDOM from 'react-dom/client'
import {
  RouterProvider,
} from 'react-router-dom'
import {router} from '@/router'
import { ThemeProvider } from '@/components/theme-provider'
import './index.css'
import { TooltipProvider } from '@/components/ui/tooltip'
import { ToastProvider } from '@/components/ui/toast'



ReactDOM.createRoot(document.getElementById('root')!).render(
  <React.StrictMode>
    <ThemeProvider defaultTheme='dark' storageKey='vite-ui-theme'>
      <TooltipProvider>
        <ToastProvider/>
      <RouterProvider router={router} />
      </TooltipProvider>
    </ThemeProvider>
  </React.StrictMode>
)
