import { ColumnDef } from "@tanstack/react-table";
import { Checkbox } from "@/components/ui/checkbox";
import { DataTableColumnHeader } from "./data-table-column-header";
import { DataTableRowActions } from "./data-table-row-actions";
import { Employee } from "@/lib/team-context-provider";

export const columns: ColumnDef<Employee>[] = [
  {
    id: "select",
    header: ({ table }) => (
      <Checkbox
        checked={
          table.getIsAllPageRowsSelected() ||
          (table.getIsSomePageRowsSelected() && "indeterminate")
        }
        onCheckedChange={(value) => table.toggleAllPageRowsSelected(!!value)}
        aria-label='Select all'
        className='translate-y-[2px]'
      />
    ),
    cell: ({ row }) => (
      <Checkbox
        checked={row.getIsSelected()}
        onCheckedChange={(value) => row.toggleSelected(!!value)}
        aria-label='Select row'
        className='translate-y-[2px]'
      />
    ),
    enableSorting: false,
    enableHiding: false,
  },
  {
    accessorKey: "uid",
    header: ({ column }) => (
      <DataTableColumnHeader column={column} title='uid' />
    ),
    cell: ({ row }) => <div className='w-[80px]'>{row.getValue("uid")}</div>,
    enableSorting: false,
    enableHiding: false,
  },
  {
    accessorKey: "firstName",
    header: ({ column }) => (
      <DataTableColumnHeader column={column} title='firstName' />
    ),
    cell: ({ row }) => (
      <div className='w-[80px]'>{row.getValue("firstName")}</div>
    ),
    enableSorting: false,
    enableHiding: false,
  },
  {
    accessorKey: "lastName",
    header: ({ column }) => (
      <DataTableColumnHeader column={column} title='lastName' />
    ),
    cell: ({ row }) => (
      <div className='w-[80px]'>{row.getValue("lastName")}</div>
    ),
    enableSorting: false,
    enableHiding: false,
  },
  {
    accessorKey: "gender",
    header: ({ column }) => (
      <DataTableColumnHeader column={column} title='gender' />
    ),
    cell: ({ row }) => <div className='w-[80px]'>{row.getValue("gender")}</div>,
    enableSorting: false,
    enableHiding: false,
  },

  {
    accessorKey: "position",
    header: ({ column }) => (
      <DataTableColumnHeader column={column} title='position' />
    ),
    cell: ({ row }) => (
      <div className='w-[80px]'>{row.getValue("position")}</div>
    ),
    enableSorting: false,
    enableHiding: false,
  },
  {
    accessorKey: "documents",
    header: ({ column }) => (
      <DataTableColumnHeader column={column} title='documents' />
    ),
    cell: ({ row }) => (
      <div className='w-[80px]'>{row.getValue("documents")}</div>
    ),
    enableSorting: false,
    enableHiding: false,
  },
  {
    accessorKey: "status",
    header: ({ column }) => (
      <DataTableColumnHeader column={column} title='status' />
    ),
    cell: ({ row }) => <div className='w-[80px]'>{row.getValue("status")}</div>,
    enableSorting: false,
    enableHiding: false,
  },
  {
    accessorKey: "createdAt",
    header: ({ column }) => (
      <DataTableColumnHeader column={column} title='createdAt' />
    ),
    cell: ({ row }) => (
      <div className='w-[80px]'>{row.getValue("createdAt")}</div>
    ),
    enableSorting: false,
    enableHiding: false,
  },
  {
    accessorKey: "modifiedAt",
    header: ({ column }) => (
      <DataTableColumnHeader column={column} title='modifiedAt' />
    ),
    cell: ({ row }) => (
      <div className='w-[80px]'>{row.getValue("modifiedAt")}</div>
    ),
    enableSorting: false,
    enableHiding: false,
  },

  // {
  //   accessorKey: "title",
  //   header: ({ column }) => (
  //     <DataTableColumnHeader column={column} title='Title' />
  //   ),
  //   cell: ({ row }) => {
  //     const label = labels.find((label) => label.value === row.original.label);

  //     return (
  //       <div className='flex space-x-2'>
  //         {label && <Badge variant='outline'>{label.label}</Badge>}
  //         <span className='max-w-32 truncate font-medium sm:max-w-72 md:max-w-[31rem]'>
  //           {row.getValue("title")}
  //         </span>
  //       </div>
  //     );
  //   },
  // },
  // {
  //   accessorKey: "status",
  //   header: ({ column }) => (
  //     <DataTableColumnHeader column={column} title='Status' />
  //   ),
  //   cell: ({ row }) => {
  //     const status = statuses.find(
  //       (status) => status.value === row.getValue("status")
  //     );

  //     if (!status) {
  //       return null;
  //     }

  //     return (
  //       <div className='flex w-[100px] items-center'>
  //         {status.icon && (
  //           <status.icon className='mr-2 h-4 w-4 text-muted-foreground' />
  //         )}
  //         <span>{status.label}</span>
  //       </div>
  //     );
  //   },
  //   filterFn: (row, id, value) => {
  //     return value.includes(row.getValue(id));
  //   },
  // },
  // {
  //   accessorKey: "priority",
  //   header: ({ column }) => (
  //     <DataTableColumnHeader column={column} title='Priority' />
  //   ),
  //   cell: ({ row }) => {
  //     const priority = priorities.find(
  //       (priority) => priority.value === row.getValue("priority")
  //     );

  //     if (!priority) {
  //       return null;
  //     }

  //     return (
  //       <div className='flex items-center'>
  //         {priority.icon && (
  //           <priority.icon className='mr-2 h-4 w-4 text-muted-foreground' />
  //         )}
  //         <span>{priority.label}</span>
  //       </div>
  //     );
  //   },
  //   filterFn: (row, id, value) => {
  //     return value.includes(row.getValue(id));
  //   },
  // },
  {
    id: "actions",
    cell: ({ row }) => <DataTableRowActions row={row} />,
  },
];
