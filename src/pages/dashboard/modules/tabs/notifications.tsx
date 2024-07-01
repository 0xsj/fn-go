import { TabsContent } from "@/components/ui/tabs"

interface Props {}

export const NotificationTab: React.FC<Props> = () => {
  return (
    <TabsContent value='notifications'>
      <div>
        <h1>notifications</h1>
      </div>
    </TabsContent>
  )
}
