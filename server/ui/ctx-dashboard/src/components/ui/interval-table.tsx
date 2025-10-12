import * as React from "react"
import {
    type UniqueIdentifier,
} from "@dnd-kit/core"
import {
    SortableContext,
    useSortable,
    verticalListSortingStrategy,
} from "@dnd-kit/sortable"
import { CSS } from "@dnd-kit/utilities"
import {
    IconChevronLeft,
    IconChevronRight,
    IconChevronsLeft,
    IconChevronsRight,
} from "@tabler/icons-react"
import {
    ColumnDef,
    ColumnFiltersState,
    flexRender,
    getCoreRowModel,
    getFacetedRowModel,
    getFacetedUniqueValues,
    getFilteredRowModel,
    getPaginationRowModel,
    getSortedRowModel,
    Row,
    SortingState,
    useReactTable,
    VisibilityState,
} from "@tanstack/react-table"
import { z } from "zod"
import { Button } from "@/components/ui/button"
import { Label } from "@/components/ui/label"
import {
    Select,
    SelectContent,
    SelectItem,
    SelectTrigger,
    SelectValue,
} from "@/components/ui/select"
import {
    Table,
    TableBody,
    TableCell,
    TableHead,
    TableHeader,
    TableRow,
} from "@/components/ui/table"
import { Edit, Trash } from "lucide-react"
import { api, ZonedDateTime } from "@/api/api"
import { useMutation, useQueryClient } from "@tanstack/react-query"
import { Sheet, SheetClose, SheetContent, SheetDescription, SheetFooter, SheetHeader, SheetTitle, SheetTrigger } from "./sheet"
import { Input } from "./input"
import { DateTimeInput } from "./datetime"
import { Dialog, DialogClose, DialogContent, DialogDescription, DialogFooter, DialogHeader, DialogTitle, DialogTrigger } from "./dialog"
import { Item, ItemActions, ItemContent, ItemDescription, ItemTitle } from "./item"
import { Separator } from "./separator"
export const schema = z.object({
    id: z.string(),
    ctxId: z.string(),
    start: z.instanceof(ZonedDateTime),
    end: z.instanceof(ZonedDateTime),
    summary: z.string(),
})
const columns: ColumnDef<z.infer<typeof schema>>[] = [
    {
        accessorKey: "startTime",
        header: "Start time",
        cell: ({ row }) => {
            return <div className="flex justify-center">{row.original.start.toString()}</div>
        },
        enableHiding: false,
    },
    {
        accessorKey: "endTime",
        header: "End time",
        cell: ({ row }) => (
            <div className="flex justify-center">{row.original.end.toString()}</div>
        ),
    },
    {
        accessorKey: "summary",
        header: "Summary",
        cell: ({ row }) => (
            <div className="flex justify-center">{row.original.summary}</div>
        ),
    },
    {
        id: "actions",
        cell: ({ row }) => (
            <ActionCell interval={row.original}></ActionCell>
        ),
    },
]
function DraggableRow({ row }: { row: Row<z.infer<typeof schema>> }) {
    const { transform, transition, setNodeRef, isDragging } = useSortable({
        id: row.original.id,
    })
    return (
        <TableRow
            data-state={row.getIsSelected() && "selected"}
            data-dragging={isDragging}
            ref={setNodeRef}
            className="relative z-0 data-[dragging=true]:z-10 data-[dragging=true]:opacity-80"
            style={{
                transform: CSS.Transform.toString(transform),
                transition: transition,
            }}
        >
            {row.getVisibleCells().map((cell) => (
                <TableCell key={cell.id}>
                    {flexRender(cell.column.columnDef.cell, cell.getContext())}
                </TableCell>
            ))}
        </TableRow>
    )
}
export function DataTable({
    data: initialData,
}: {
    data: z.infer<typeof schema>[]
}) {
    const [data, setData] = React.useState(() => initialData)
    const [rowSelection, setRowSelection] = React.useState({})
    const [columnVisibility, setColumnVisibility] =
        React.useState<VisibilityState>({})
    const [columnFilters, setColumnFilters] = React.useState<ColumnFiltersState>(
        []
    )
    const [sorting, setSorting] = React.useState<SortingState>([])
    const [pagination, setPagination] = React.useState({
        pageIndex: 0,
        pageSize: 10,
    })
    const dataIds = React.useMemo<UniqueIdentifier[]>(
        () => data?.map(({ id }) => id) || [],
        [data]
    )
    const table = useReactTable({
        data,
        columns,
        state: {
            sorting,
            columnVisibility,
            rowSelection,
            columnFilters,
            pagination,
        },
        getRowId: (row) => row.id.toString(),
        enableRowSelection: true,
        onRowSelectionChange: setRowSelection,
        onSortingChange: setSorting,
        onColumnFiltersChange: setColumnFilters,
        onColumnVisibilityChange: setColumnVisibility,
        onPaginationChange: setPagination,
        getCoreRowModel: getCoreRowModel(),
        getFilteredRowModel: getFilteredRowModel(),
        getPaginationRowModel: getPaginationRowModel(),
        getSortedRowModel: getSortedRowModel(),
        getFacetedRowModel: getFacetedRowModel(),
        getFacetedUniqueValues: getFacetedUniqueValues(),
    })
    return (
        <>
            <div className="overflow-hidden rounded-lg border">
                <Table>
                    <TableHeader className="bg-muted sticky top-0 z-10">
                        {table.getHeaderGroups().map((headerGroup) => (
                            <TableRow key={headerGroup.id}>
                                {headerGroup.headers.map((header) => {
                                    return (
                                        <TableHead key={header.id} colSpan={header.colSpan} className="text-center">
                                            {header.isPlaceholder
                                                ? null
                                                : flexRender(
                                                    header.column.columnDef.header,
                                                    header.getContext()
                                                )}
                                        </TableHead>
                                    )
                                })}
                            </TableRow>
                        ))}
                    </TableHeader>
                    <TableBody className="**:data-[slot=table-cell]:first:w-8">
                        {table.getRowModel().rows?.length ? (
                            <SortableContext
                                items={dataIds}
                                strategy={verticalListSortingStrategy}
                            >
                                {table.getRowModel().rows.map((row) => (
                                    <DraggableRow key={row.id} row={row} />
                                ))}
                            </SortableContext>
                        ) : (
                            <TableRow>
                                <TableCell
                                    colSpan={columns.length}
                                    className="h-24 text-center"
                                >
                                    No results.
                                </TableCell>
                            </TableRow>
                        )}
                    </TableBody>
                </Table>
            </div>
            {data.length > table.getState().pagination.pageSize &&
                <div className="flex items-center justify-between px-4">
                    <div className="text-muted-foreground hidden flex-1 text-sm lg:flex">
                    </div>
                    <div className="flex w-full items-center gap-8 lg:w-fit">
                        <div className="hidden items-center gap-2 lg:flex">
                            <Label htmlFor="rows-per-page" className="text-sm font-medium">
                                Rows per page
                            </Label>
                            <Select
                                value={`${table.getState().pagination.pageSize}`}
                                onValueChange={(value) => {
                                    table.setPageSize(Number(value))
                                }}
                            >
                                <SelectTrigger size="sm" className="w-20" id="rows-per-page">
                                    <SelectValue
                                        placeholder={table.getState().pagination.pageSize}
                                    />
                                </SelectTrigger>
                                <SelectContent side="top">
                                    {[10, 20, 30, 40, 50].map((pageSize) => (
                                        <SelectItem key={pageSize} value={`${pageSize}`}>
                                            {pageSize}
                                        </SelectItem>
                                    ))}
                                </SelectContent>
                            </Select>
                        </div>
                        <div className="flex w-fit items-center justify-center text-sm font-medium">
                            Page {table.getState().pagination.pageIndex + 1} of{" "}
                            {table.getPageCount()}
                        </div>
                        <div className="ml-auto flex items-center gap-2 lg:ml-0">
                            <Button
                                variant="outline"
                                className="hidden h-8 w-8 p-0 lg:flex"
                                onClick={() => table.setPageIndex(0)}
                                disabled={!table.getCanPreviousPage()}
                            >
                                <span className="sr-only">Go to first page</span>
                                <IconChevronsLeft />
                            </Button>
                            <Button
                                variant="outline"
                                className="size-8"
                                size="icon"
                                onClick={() => table.previousPage()}
                                disabled={!table.getCanPreviousPage()}
                            >
                                <span className="sr-only">Go to previous page</span>
                                <IconChevronLeft />
                            </Button>
                            <Button
                                variant="outline"
                                className="size-8"
                                size="icon"
                                onClick={() => table.nextPage()}
                                disabled={!table.getCanNextPage()}
                            >
                                <span className="sr-only">Go to next page</span>
                                <IconChevronRight />
                            </Button>
                            <Button
                                variant="outline"
                                className="hidden size-8 lg:flex"
                                size="icon"
                                onClick={() => table.setPageIndex(table.getPageCount() - 1)}
                                disabled={!table.getCanNextPage()}
                            >
                                <span className="sr-only">Go to last page</span>
                                <IconChevronsRight />
                            </Button>
                        </div>
                    </div>
                </div>
            }
        </>
    )
}

function ActionCell({ interval }: {
    interval: {
        id: string,
        ctxId: string,
        start: ZonedDateTime,
        end: ZonedDateTime,
        summary: string,
    }
}) {
    const qc = useQueryClient();
    const deleteMutation = useMutation(api.intervals.deleteMutation(qc))

    const handleSubmit = (e: React.FormEvent<HTMLFormElement>) => {
        e.preventDefault();
    }

    return (
        <div className="flex gap-2 justify-end">
            <Dialog>
                <DialogTrigger asChild>
                    <Button variant="ghost"><Edit /></Button>
                </DialogTrigger>

                <DialogContent className="grid grid-rows-[auto_1fr_auto] sm:max-w-[520px] p-0">
                    <form className="contents" onSubmit={handleSubmit}>
                        <div className="p-6 pb-4">
                            <DialogHeader>
                                <DialogTitle>Edit interval</DialogTitle>
                                <DialogDescription>Adjust start and end times.</DialogDescription>
                            </DialogHeader>
                        </div>

                        <div className="p-6 pt-0 overflow-auto">
                            <div className="grid gap-4 sm:grid-cols-2">
                                <div className="grid gap-2">
                                    <Label htmlFor="start">Start</Label>
                                    <DateTimeInput id="start" datetime={interval.start.toDateTime()} editable />
                                </div>
                                <div className="grid gap-2">
                                    <Label htmlFor="end">End</Label>
                                    <DateTimeInput id="end" datetime={interval.end.toDateTime()} editable />
                                </div>
                            </div>

                            <div className="mt-6">
                                <Item variant="outline" className="border-destructive/30 bg-destructive/5">
                                    <ItemContent>
                                        <ItemTitle className="text-destructive">Delete interval</ItemTitle>
                                        <ItemDescription className="text-muted-foreground"> This action is irreversible </ItemDescription>
                                    </ItemContent>
                                    <ItemActions>
                                        <Button type="button" variant="destructive" size="sm" onClick={() => deleteMutation.mutate({ ctxId: interval.ctxId, id: interval.id })}>
                                            <Trash></Trash>
                                        </Button>
                                    </ItemActions>
                                </Item>
                            </div>
                        </div>

                        <DialogFooter className="sm:space-x-0 gap-2 flex-row justify-end border-t bg-background p-4">
                            <DialogClose asChild>
                                <Button type="button" variant="outline">Cancel</Button>
                            </DialogClose>
                            <Button type="submit">Save changes</Button>
                        </DialogFooter>
                    </form>
                </DialogContent>
            </Dialog>
        </div>
    );
}