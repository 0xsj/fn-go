import { TabsContent } from "@/components/ui/tabs"

interface Props {}

export const TasksTab: React.FC<Props> = () => {
  return (
    <TabsContent value='tasks'>
      <div>
        <h1>tasks</h1>
      </div>
    </TabsContent>
  )
}
