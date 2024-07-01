import { Layout } from "@/components/custom/layout"
import { ThemeSwitch } from "@/components/custom/theme-switch"
import { Sidebar } from "@/components/sidebar"
import { UserNav } from "@/components/user-nav"
import { useState } from "react"

export default function Dashboard() {
  return (
    <Layout>
      <Layout.Header sticky>
        <div className='ml-auto flex items-center space-x-4'>
          <ThemeSwitch />
          <UserNav />
        </div>
      </Layout.Header>
      <Layout.Body>
        <div></div>
      </Layout.Body>
    </Layout>
  )
}
