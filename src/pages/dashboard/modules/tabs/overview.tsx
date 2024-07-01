import { TabsContent } from "@/components/ui/tabs"

interface Props {}

export const OverviewTab: React.FC<Props> = () => {
  return (
    <TabsContent value='overview'>
      <div>
        <h1>overview</h1>
      </div>
    </TabsContent>
  )
}
