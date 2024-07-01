import { TabsContent } from "@/components/ui/tabs"

interface Props {}

export const ReportsTab: React.FC<Props> = () => {
  return (
    <TabsContent value='reports'>
      <div>
        <h1>reports</h1>
      </div>
    </TabsContent>
  )
}
