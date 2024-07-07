import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { TabsContent } from "@/components/ui/tabs";
import { IconPdf, IconFile, IconBook } from "@tabler/icons-react";
import { DataTable } from "../table/data-table";
import { employees } from "@/data/data";
import { columns } from "../table/columns";
import { documentTasks, tasks } from "../../data/tasks";

export const AllTab: React.FC = () => {
  return (
    <TabsContent value='all' className='space-y-4'>
      <div className='grid gap-3 sm:grid-cols-1 lg:grid-cols-3'>
        <Card className=''>
          <CardHeader>
            <IconPdf size={18} />
          </CardHeader>
          <CardTitle></CardTitle>
          <CardContent>
            <div>W4</div>
            <div>year / version</div>
            <div>date stamp</div>
          </CardContent>
        </Card>
        <Card className=''>
          <CardHeader>
            <IconFile size={18} />
          </CardHeader>
          <CardTitle></CardTitle>
          <CardContent>
            <div>Application</div>
            <div>year / version</div>
            <div>date stamp</div>
          </CardContent>
        </Card>
        <Card className=''>
          <CardHeader>
            <IconBook size={18} />
          </CardHeader>
          <CardTitle></CardTitle>
          <CardContent>
            <div>Employee handbook</div>
            <div>year / version</div>
            <div>date stamp</div>
          </CardContent>
        </Card>
      </div>
      <div className='-mx-4 flex-1 overflow-auto px-4 py-1 lg:flex-row lg:space-x-12 lg:space-y-0'>
        <DataTable data={documentTasks} columns={columns} />
      </div>
    </TabsContent>
  );
};
