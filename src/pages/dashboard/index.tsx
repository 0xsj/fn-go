import { Button } from "@/components/custom/button"
import { Layout } from "@/components/custom/layout"
import { ThemeSwitch } from "@/components/custom/theme-switch"
import { Tabs, TabsList, TabsTrigger } from "@/components/ui/tabs"
import { UserNav } from "@/components/user-nav"
import { OverviewTab } from "./modules/tabs/overview"
import { ReportsTab } from "./modules/tabs/reports"
import { NotificationTab } from "./modules/tabs/notifications"
import { TasksTab } from "./modules/tabs/tasks"
import { SearchBox } from "@/components/custom/search-box"
import { HeaderNav } from "@/components/header-nav"

export default function Dashboard() {
  return (
    <Layout>
      <Layout.Header sticky>
        <HeaderNav links={navItems} />
        <div className='ml-auto flex items-center space-x-4'>
          <SearchBox />
          <ThemeSwitch />
          <UserNav />
        </div>
      </Layout.Header>
      <Layout.Body>
        <div className='mb-2 flex items-center justify-between space-y-2'>
          <h1 className='text-2xl font-bold tracking-tight'>Overview</h1>
          <div className='flex items-center space-x-2'>
            <Button>Download</Button>
          </div>
        </div>
        <Tabs
          orientation='vertical'
          defaultValue='overview'
          className='space-y-4'
        >
          <div>
            <TabsList>
              <TabsTrigger value='overview'>Overview</TabsTrigger>
              <TabsTrigger value='reports'>Reports</TabsTrigger>
              <TabsTrigger value='tasks'>Tasks</TabsTrigger>
              <TabsTrigger value='notifications'>Notifications</TabsTrigger>
            </TabsList>
          </div>
          <OverviewTab />
          <ReportsTab />
          <NotificationTab />
          <TasksTab />
        </Tabs>
      </Layout.Body>
    </Layout>
  )
}

const navItems = [
  {
    title: "Overview",
    href: "/",
    isActive: true,
  },
  {
    title: "Github",
    href: "https://github.com/",
    isActive: true,
  },
]
