import { Card } from "@/components/ui/card";
import { TabsContent } from "@/components/ui/tabs";

interface Props {}

export const OverviewTab: React.FC<Props> = () => {
  return (
    <TabsContent value='overview'>
      <div>
        <h1>overview</h1>
        <Card></Card>
      </div>
    </TabsContent>
  );
};
