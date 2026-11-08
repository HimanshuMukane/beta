"use client"

import * as React from "react"
import {
  useTable,
  type ColumnDef,
  type ColumnFiltersState,
  type ColumnVisibilityState,
  type RowData,
  type SortingState,
  type Row,
} from "@tanstack/react-table"
import {
  IconChevronLeft,
  IconChevronRight,
  IconChevronsLeft,
  IconChevronsRight,
  IconAdjustmentsHorizontal,
  IconSearch,
  IconX,
  IconFilter,
  IconDownload,
  IconDotsVertical,
  IconArrowsSort,
  IconSortAscending,
  IconSortDescending,
  IconRefresh,
  IconAlertCircle,
  IconInbox,
} from "@tabler/icons-react"

import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "../ui/table"
import { Button } from "../ui/button"
import { Input } from "../ui/input"
import { Badge } from "../ui/badge"
import {
  DropdownMenu,
  DropdownMenuCheckboxItem,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuRadioGroup,
  DropdownMenuRadioItem,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from "../ui/dropdown-menu"
import {
  Popover,
  PopoverContent,
  PopoverTrigger,
} from "../ui/popover"
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "../ui/select"
import {
  Sheet,
  SheetContent,
  SheetDescription,
  SheetHeader,
  SheetTitle,
} from "../ui/sheet"
import {
  AlertDialog,
  AlertDialogAction,
  AlertDialogCancel,
  AlertDialogContent,
  AlertDialogDescription,
  AlertDialogFooter,
  AlertDialogHeader,
  AlertDialogTitle,
} from "../ui/alert-dialog"
import { Skeleton } from "../ui/skeleton"
import { Tooltip, TooltipContent, TooltipProvider, TooltipTrigger } from "../ui/tooltip"
import { toast } from "sonner"
import { features, type DataTableFeatures } from "./data-table-features"
import { cn } from "@repo/ui/lib/utils"

export type TableDensity = "compact" | "comfortable"

export interface DataTableFilterOption {
  label: string
  value: string
  icon?: React.ComponentType<{ className?: string }>
}

export interface DataTableFilterConfig {
  columnId: string
  title: string
  type: "select" | "multi-select" | "text" | "number-range"
  options?: DataTableFilterOption[]
}

export interface DataTableBulkAction<TData> {
  label: string
  icon?: React.ComponentType<{ className?: string }>
  variant?: "default" | "destructive" | "outline" | "secondary"
  requiresConfirmation?: boolean
  confirmationTitle?: string
  confirmationDescription?: (selectedCount: number) => string
  action: (selectedRows: TData[], clearSelection: () => void) => void | Promise<void>
}

export interface DataTableRowAction<TData> {
  label: string
  icon?: React.ComponentType<{ className?: string }>
  variant?: "default" | "destructive"
  separatorBefore?: boolean
  tooltip?: string
  action: (row: TData) => void
}

export interface DataTableProps<TData extends RowData> {
  columns: ColumnDef<DataTableFeatures, TData>[]
  data: TData[]
  searchKey?: string
  searchPlaceholder?: string
  searchableColumns?: string[]
  filterConfigs?: DataTableFilterConfig[]
  bulkActions?: DataTableBulkAction<TData>[]
  rowActions?: DataTableRowAction<TData>[]
  onRowClick?: (row: TData) => void
  renderDetailDrawer?: (row: TData, close: () => void) => React.ReactNode
  detailDrawerTitle?: string | ((row: TData) => string)
  detailDrawerDescription?: string | ((row: TData) => string)
  isLoading?: boolean
  error?: string | null
  onRetry?: () => void
  emptyTitle?: string
  emptyDescription?: string
  emptyAction?: {
    label: string
    onClick: () => void
  }
  defaultDensity?: TableDensity
  enableExport?: boolean
  exportFilename?: string
}

export function DataTable<TData extends RowData>({
  columns,
  data,
  searchKey,
  searchPlaceholder = "Search records... (Press ⌘K)",
  searchableColumns,
  filterConfigs = [],
  bulkActions = [],
  rowActions = [],
  onRowClick,
  renderDetailDrawer,
  detailDrawerTitle = "Record Details",
  detailDrawerDescription,
  isLoading = false,
  error = null,
  onRetry,
  emptyTitle = "No records found",
  emptyDescription = "There are no records in this workspace yet.",
  emptyAction,
  defaultDensity = "comfortable",
  enableExport = true,
  exportFilename = "export.csv",
}: DataTableProps<TData>) {
  const [sorting, setSorting] = React.useState<SortingState>([])
  const [columnFilters, setColumnFilters] = React.useState<ColumnFiltersState>([])
  const [globalFilter, setGlobalFilter] = React.useState<string>("")
  const [columnVisibility, setColumnVisibility] = React.useState<ColumnVisibilityState>({})
  const [rowSelection, setRowSelection] = React.useState<Record<string, true>>({})
  const [density, setDensity] = React.useState<TableDensity>(() => {
    if (typeof window !== "undefined") {
      return (localStorage.getItem("data-table-density") as TableDensity) || defaultDensity
    }
    return defaultDensity
  })

  const [activeDrawerRow, setActiveDrawerRow] = React.useState<TData | null>(null)

  const [pendingBulkAction, setPendingBulkAction] = React.useState<{
    action: DataTableBulkAction<TData>
    count: number
  } | null>(null)
  const [isBulkProcessing, setIsBulkProcessing] = React.useState(false)

  const searchInputRef = React.useRef<HTMLInputElement>(null)

  const handleDensityChange = (val: TableDensity) => {
    setDensity(val)
    if (typeof window !== "undefined") {
      localStorage.setItem("data-table-density", val)
    }
  }

  React.useEffect(() => {
    const handleKeyDown = (e: KeyboardEvent) => {
      if ((e.metaKey || e.ctrlKey) && e.key.toLowerCase() === "k") {
        const target = e.target as HTMLElement
        if (target && (target.tagName === "INPUT" || target.tagName === "TEXTAREA" || target.isContentEditable)) {
          return
        }
        e.preventDefault()
        searchInputRef.current?.focus()
      }
    }
    window.addEventListener("keydown", handleKeyDown)
    return () => window.removeEventListener("keydown", handleKeyDown)
  }, [])

  const tableColumns = React.useMemo(() => {
    if (!rowActions || rowActions.length === 0) {
      return columns
    }

    const actionColumn: ColumnDef<DataTableFeatures, TData> = {
      id: "actions",
      enableHiding: false,
      enableSorting: false,
      header: () => <span className="sr-only">Actions</span>,
      cell: ({ row }) => {
        const rowData = row.original
        return (
          <div
            className="flex items-center justify-end"
            onClick={(e) => e.stopPropagation()}
          >
            <DropdownMenu>
              <TooltipProvider>
                <Tooltip>
                  <TooltipTrigger asChild>
                    <DropdownMenuTrigger asChild>
                      <Button
                        variant="ghost"
                        size="icon"
                        className="h-8 w-8 p-0 text-muted-foreground hover:text-foreground"
                        aria-label="Open row actions"
                      >
                        <IconDotsVertical className="h-4 w-4" />
                      </Button>
                    </DropdownMenuTrigger>
                  </TooltipTrigger>
                  <TooltipContent>Row actions</TooltipContent>
                </Tooltip>
              </TooltipProvider>
              <DropdownMenuContent align="end" className="w-44">
                <DropdownMenuLabel className="text-xs font-normal text-muted-foreground">
                  Actions
                </DropdownMenuLabel>
                <DropdownMenuSeparator />
                {rowActions.map((action, idx) => {
                  const Icon = action.icon
                  return (
                    <React.Fragment key={idx}>
                      {action.separatorBefore && <DropdownMenuSeparator />}
                      <DropdownMenuItem
                        className={cn(
                          action.variant === "destructive" &&
                            "text-destructive focus:bg-destructive/10 focus:text-destructive"
                        )}
                        onClick={() => action.action(rowData)}
                      >
                        {Icon && <Icon className="mr-2 h-4 w-4" />}
                        <span>{action.label}</span>
                      </DropdownMenuItem>
                    </React.Fragment>
                  )
                })}
              </DropdownMenuContent>
            </DropdownMenu>
          </div>
        )
      },
    }

    return [...columns, actionColumn]
  }, [columns, rowActions])

  const globalFilterFn = React.useCallback(
    (row: Row<DataTableFeatures, TData>, _columnId: string, filterValue: string) => {
      if (!filterValue) return true
      const searchStr = filterValue.toLowerCase().trim()
      if (!searchStr) return true

      if (searchableColumns && searchableColumns.length > 0) {
        return searchableColumns.some((colId) => {
          const val = row.getValue(colId)
          if (val === null || val === undefined) return false
          return String(val).toLowerCase().includes(searchStr)
        })
      }

      if (searchKey) {
        const val = row.getValue(searchKey)
        if (val === null || val === undefined) return false
        return String(val).toLowerCase().includes(searchStr)
      }

      const original = row.original as Record<string, unknown>
      return Object.values(original).some((val) => {
        if (val === null || val === undefined) return false
        if (typeof val === "object") return false
        return String(val).toLowerCase().includes(searchStr)
      })
    },
    [searchKey, searchableColumns]
  )

  const table = useTable({
    features,
    data,
    columns: tableColumns,
    onSortingChange: setSorting,
    onColumnFiltersChange: setColumnFilters,
    onGlobalFilterChange: setGlobalFilter,
    globalFilterFn,
    onColumnVisibilityChange: setColumnVisibility,
    onRowSelectionChange: setRowSelection,
    state: {
      sorting,
      columnFilters,
      globalFilter,
      columnVisibility,
      rowSelection,
    },
  })

  const handleExportCSV = () => {
    try {
      const rows = table.getFilteredRowModel().rows
      if (!rows.length) {
        toast.info("No records to export")
        return
      }

      const visibleColumns = table
        .getAllLeafColumns()
        .filter((col) => col.id !== "select" && col.id !== "actions" && col.getIsVisible())

      const headers = visibleColumns.map((c) => `"${c.id.replace(/"/g, '""')}"`).join(",")
      const csvRows = rows.map((row) => {
        return visibleColumns
          .map((col) => {
            const val = row.getValue(col.id)
            if (val === null || val === undefined) return '""'
            const strVal = String(val).replace(/"/g, '""')
            return `"${strVal}"`
          })
          .join(",")
      })

      const csvContent = [headers, ...csvRows].join("\n")
      const blob = new Blob([csvContent], { type: "text/csv;charset=utf-8;" })
      const url = URL.createObjectURL(blob)
      const link = document.createElement("a")
      link.setAttribute("href", url)
      link.setAttribute("download", exportFilename)
      document.body.appendChild(link)
      link.click()
      document.body.removeChild(link)
      URL.revokeObjectURL(url)

      toast.success(`Exported ${rows.length} record(s) to CSV`)
    } catch {
      toast.error("Failed to export data to CSV")
    }
  }

  const activeFiltersCount = columnFilters.length + (globalFilter ? 1 : 0)

  const clearAllFilters = () => {
    setColumnFilters([])
    setGlobalFilter("")
  }

  const removeFilter = (id: string) => {
    if (id === "search") {
      setGlobalFilter("")
    } else {
      setColumnFilters((prev) => prev.filter((f) => f.id !== id))
    }
  }

  const selectedRowsModel = table.getFilteredSelectedRowModel().rows
  const selectedCount = selectedRowsModel.length
  const totalFilteredCount = table.getFilteredRowModel().rows.length
  const totalRowCount = data.length

  const handleBulkActionClick = (action: DataTableBulkAction<TData>) => {
    if (action.requiresConfirmation) {
      setPendingBulkAction({
        action,
        count: selectedCount,
      })
    } else {
      executeBulkAction(action)
    }
  }

  const executeBulkAction = async (action: DataTableBulkAction<TData>) => {
    try {
      setIsBulkProcessing(true)
      const selectedData = selectedRowsModel.map((r) => r.original)
      await action.action(selectedData, () => setRowSelection({}))
    } catch {
      toast.error("Bulk action failed")
    } finally {
      setIsBulkProcessing(false)
      setPendingBulkAction(null)
    }
  }

  return (
    <div className="space-y-4">
      {selectedCount > 0 && (
        <div className="flex flex-wrap items-center justify-between gap-3 rounded-lg border border-primary/20 bg-primary/5 px-4 py-2.5 text-sm transition-all animate-in fade-in-50">
          <div className="flex items-center gap-2">
            <Badge variant="default" className="font-medium">
              {selectedCount} selected
            </Badge>
            <span className="text-muted-foreground text-xs md:text-sm">
              of {totalFilteredCount} matching record(s)
            </span>
          </div>

          <div className="flex items-center gap-2">
            {bulkActions.map((action, idx) => {
              const Icon = action.icon
              return (
                <Button
                  key={idx}
                  size="sm"
                  variant={action.variant || "outline"}
                  onClick={() => handleBulkActionClick(action)}
                  disabled={isBulkProcessing}
                  className="h-8 gap-1.5 text-xs font-medium"
                >
                  {Icon && <Icon className="h-3.5 w-3.5" />}
                  {action.label}
                </Button>
              )
            })}
            <Button
              size="sm"
              variant="ghost"
              onClick={() => setRowSelection({})}
              disabled={isBulkProcessing}
              className="h-8 text-xs text-muted-foreground hover:text-foreground"
            >
              Clear selection
            </Button>
          </div>
        </div>
      )}

      <div className="flex flex-col gap-3 md:flex-row md:items-center md:justify-between">
        <div className="flex flex-1 flex-wrap items-center gap-2">
          <div className="relative w-full sm:w-[260px] lg:w-[320px]">
            <IconSearch className="absolute left-3 top-1/2 h-4 w-4 -translate-y-1/2 text-muted-foreground pointer-events-none" />
            <Input
              ref={searchInputRef}
              placeholder={searchPlaceholder}
              value={globalFilter}
              onChange={(e) => setGlobalFilter(e.target.value)}
              className="h-9 pl-9 pr-8 text-sm"
            />
            {globalFilter && (
              <button
                type="button"
                onClick={() => setGlobalFilter("")}
                className="absolute right-2.5 top-1/2 -translate-y-1/2 text-muted-foreground hover:text-foreground p-0.5 rounded-sm"
                aria-label="Clear search"
              >
                <IconX className="h-3.5 w-3.5" />
              </button>
            )}
          </div>

          {filterConfigs.length > 0 && (
            <Popover>
              <PopoverTrigger asChild>
                <Button variant="outline" size="sm" className="h-9 gap-1.5">
                  <IconFilter className="h-4 w-4 text-muted-foreground" />
                  <span>Filters</span>
                  {columnFilters.length > 0 && (
                    <Badge variant="secondary" className="ml-1 h-5 px-1.5 text-[11px]">
                      {columnFilters.length}
                    </Badge>
                  )}
                </Button>
              </PopoverTrigger>
              <PopoverContent className="w-80 p-4" align="start">
                <div className="space-y-4">
                  <div className="flex items-center justify-between pb-2 border-b">
                    <span className="font-semibold text-sm">Filter Records</span>
                    {columnFilters.length > 0 && (
                      <Button
                        variant="ghost"
                        size="sm"
                        onClick={() => setColumnFilters([])}
                        className="h-7 px-2 text-xs text-muted-foreground hover:text-foreground"
                      >
                        Reset
                      </Button>
                    )}
                  </div>

                  <div className="space-y-3 max-h-[300px] overflow-y-auto pr-1">
                    {filterConfigs.map((config) => {
                      const currentVal = (table.getColumn(config.columnId)?.getFilterValue() as string) || ""

                      if (config.type === "select" && config.options) {
                        return (
                          <div key={config.columnId} className="space-y-1.5">
                            <label className="text-xs font-medium text-muted-foreground">
                              {config.title}
                            </label>
                            <Select
                              value={currentVal || "ALL"}
                              onValueChange={(val) => {
                                const targetCol = table.getColumn(config.columnId)
                                if (val === "ALL") {
                                  targetCol?.setFilterValue(undefined)
                                } else {
                                  targetCol?.setFilterValue(val)
                                }
                              }}
                            >
                              <SelectTrigger className="h-8 text-xs">
                                <SelectValue placeholder={`All ${config.title}`} />
                              </SelectTrigger>
                              <SelectContent>
                                <SelectItem value="ALL" className="text-xs">
                                  All {config.title}
                                </SelectItem>
                                {config.options.map((opt) => (
                                  <SelectItem key={opt.value} value={opt.value} className="text-xs">
                                    {opt.label}
                                  </SelectItem>
                                ))}
                              </SelectContent>
                            </Select>
                          </div>
                        )
                      }

                      return (
                        <div key={config.columnId} className="space-y-1.5">
                          <label className="text-xs font-medium text-muted-foreground">
                            {config.title}
                          </label>
                          <Input
                            placeholder={`Filter by ${config.title.toLowerCase()}...`}
                            value={currentVal}
                            onChange={(e) => {
                              const targetCol = table.getColumn(config.columnId)
                              targetCol?.setFilterValue(e.target.value || undefined)
                            }}
                            className="h-8 text-xs"
                          />
                        </div>
                      )
                    })}
                  </div>
                </div>
              </PopoverContent>
            </Popover>
          )}

          {activeFiltersCount > 0 && (
            <Button
              variant="ghost"
              size="sm"
              onClick={clearAllFilters}
              className="h-9 px-2.5 text-xs text-muted-foreground hover:text-foreground"
            >
              Clear all ({activeFiltersCount})
            </Button>
          )}
        </div>

        <div className="flex items-center gap-2 self-end md:self-auto">
          {enableExport && (
            <TooltipProvider>
              <Tooltip>
                <TooltipTrigger asChild>
                  <Button
                    variant="outline"
                    size="sm"
                    onClick={handleExportCSV}
                    className="h-9 gap-1.5 text-xs"
                  >
                    <IconDownload className="h-4 w-4 text-muted-foreground" />
                    <span className="hidden sm:inline">Export</span>
                  </Button>
                </TooltipTrigger>
                <TooltipContent>Export current table view to CSV</TooltipContent>
              </Tooltip>
            </TooltipProvider>
          )}

          <DropdownMenu>
            <DropdownMenuTrigger asChild>
              <Button variant="outline" size="sm" className="h-9 gap-1.5 text-xs">
                <IconAdjustmentsHorizontal className="h-4 w-4 text-muted-foreground" />
                <span>View</span>
              </Button>
            </DropdownMenuTrigger>
            <DropdownMenuContent align="end" className="w-56">
              <DropdownMenuLabel className="text-xs">Density</DropdownMenuLabel>
              <DropdownMenuRadioGroup
                value={density}
                onValueChange={(val) => handleDensityChange(val as TableDensity)}
              >
                <DropdownMenuRadioItem value="compact" className="text-xs">
                  Compact (44px)
                </DropdownMenuRadioItem>
                <DropdownMenuRadioItem value="comfortable" className="text-xs">
                  Comfortable (56px)
                </DropdownMenuRadioItem>
              </DropdownMenuRadioGroup>

              <DropdownMenuSeparator />
              <DropdownMenuLabel className="text-xs">Toggle Columns</DropdownMenuLabel>
              <div className="max-h-48 overflow-y-auto">
                {table
                  .getAllColumns()
                  .filter((col) => col.id !== "select" && col.id !== "actions" && col.getCanHide())
                  .map((column) => {
                    return (
                      <DropdownMenuCheckboxItem
                        key={column.id}
                        className="capitalize text-xs"
                        checked={column.getIsVisible()}
                        onCheckedChange={(val) => column.toggleVisibility(!!val)}
                      >
                        {column.id}
                      </DropdownMenuCheckboxItem>
                    )
                  })}
              </div>
            </DropdownMenuContent>
          </DropdownMenu>
        </div>
      </div>

      {activeFiltersCount > 0 && (
        <div className="flex flex-wrap items-center gap-1.5 pt-1">
          <span className="text-xs text-muted-foreground mr-1">Active filters:</span>
          {globalFilter && (
            <Badge
              variant="secondary"
              className="gap-1 px-2 py-0.5 text-xs font-normal hover:bg-secondary/80"
            >
              <span>Search: &ldquo;{globalFilter}&rdquo;</span>
              <button
                type="button"
                onClick={() => removeFilter("search")}
                className="hover:text-foreground text-muted-foreground ml-0.5"
                aria-label="Remove search filter"
              >
                <IconX className="h-3 w-3" />
              </button>
            </Badge>
          )}

          {columnFilters.map((filter) => {
            const config = filterConfigs.find((c) => c.columnId === filter.id)
            const label = config ? config.title : filter.id
            const displayVal = String(filter.value)
            return (
              <Badge
                key={filter.id}
                variant="secondary"
                className="gap-1 px-2 py-0.5 text-xs font-normal hover:bg-secondary/80"
              >
                <span>
                  {label}: <strong className="font-medium">{displayVal}</strong>
                </span>
                <button
                  type="button"
                  onClick={() => removeFilter(filter.id)}
                  className="hover:text-foreground text-muted-foreground ml-0.5"
                  aria-label={`Remove ${label} filter`}
                >
                  <IconX className="h-3 w-3" />
                </button>
              </Badge>
            )
          })}
        </div>
      )}

      <div className="overflow-hidden rounded-lg border border-border bg-card shadow-sm">
        <div className="relative w-full overflow-x-auto">
          <Table>
            <TableHeader className="bg-muted/40 border-b">
              {table.getHeaderGroups().map((headerGroup) => (
                <TableRow key={headerGroup.id} className="hover:bg-transparent">
                  {headerGroup.headers.map((header) => {
                    const canSort = header.column.getCanSort()
                    const isSorted = header.column.getIsSorted()

                    return (
                      <TableHead
                        key={header.id}
                        className={cn(
                          "h-12 px-4 text-xs font-semibold tracking-wider text-muted-foreground select-none uppercase",
                          header.id === "select" && "w-10 pr-0",
                          header.id === "actions" && "w-12 text-right"
                        )}
                      >
                        {header.isPlaceholder ? null : canSort ? (
                          <button
                            type="button"
                            className="flex items-center gap-1.5 hover:text-foreground transition-colors group focus:outline-none"
                            onClick={header.column.getToggleSortingHandler()}
                          >
                            <span>
                              <table.FlexRender header={header} />
                            </span>
                            <span className="flex items-center text-muted-foreground group-hover:text-foreground">
                              {isSorted === "asc" ? (
                                <IconSortAscending className="h-3.5 w-3.5 text-primary" />
                              ) : isSorted === "desc" ? (
                                <IconSortDescending className="h-3.5 w-3.5 text-primary" />
                              ) : (
                                <IconArrowsSort className="h-3.5 w-3.5 opacity-40 group-hover:opacity-100 transition-opacity" />
                              )}
                            </span>
                          </button>
                        ) : (
                          <table.FlexRender header={header} />
                        )}
                      </TableHead>
                    )
                  })}
                </TableRow>
              ))}
            </TableHeader>

            <TableBody>
              {isLoading ? (
                Array.from({ length: 5 }).map((_, rIdx) => (
                  <TableRow
                    key={`skeleton-${rIdx}`}
                    className={density === "compact" ? "h-11" : "h-14"}
                  >
                    {table.getVisibleLeafColumns().map((_, cIdx) => (
                      <TableCell key={`skeleton-cell-${cIdx}`} className="px-4 py-2">
                        <Skeleton className="h-4 w-full max-w-[80%]" />
                      </TableCell>
                    ))}
                  </TableRow>
                ))
              ) : error ? (
                <TableRow>
                  <TableCell
                    colSpan={tableColumns.length}
                    className="h-56 text-center"
                  >
                    <div className="flex flex-col items-center justify-center gap-2 max-w-sm mx-auto text-center">
                      <div className="rounded-full bg-destructive/10 p-3 text-destructive">
                        <IconAlertCircle className="h-6 w-6" />
                      </div>
                      <h3 className="font-medium text-sm text-foreground">Failed to load data</h3>
                      <p className="text-xs text-muted-foreground">{error}</p>
                      {onRetry && (
                        <Button
                          variant="outline"
                          size="sm"
                          onClick={onRetry}
                          className="mt-2 h-8 gap-1.5 text-xs"
                        >
                          <IconRefresh className="h-3.5 w-3.5" />
                          Retry
                        </Button>
                      )}
                    </div>
                  </TableCell>
                </TableRow>
              ) : table.getRowModel().rows?.length ? (
                table.getRowModel().rows.map((row) => {
                  const isSelected = row.getIsSelected()
                  const hasClick = !!onRowClick || !!renderDetailDrawer

                  return (
                    <TableRow
                      key={row.id}
                      data-state={isSelected && "selected"}
                      onClick={() => {
                        if (renderDetailDrawer) {
                          setActiveDrawerRow(row.original)
                        }
                        if (onRowClick) {
                          onRowClick(row.original)
                        }
                      }}
                      className={cn(
                        "transition-colors",
                        density === "compact" ? "h-11 py-1" : "h-14 py-3",
                        hasClick && "cursor-pointer hover:bg-muted/60",
                        isSelected && "bg-muted/80 font-medium"
                      )}
                    >
                      {row.getVisibleCells().map((cell) => (
                        <TableCell
                          key={cell.id}
                          className={cn(
                            "px-4 align-middle text-sm",
                            density === "compact" ? "py-1.5" : "py-3",
                            cell.column.id === "select" && "w-10 pr-0",
                            cell.column.id === "actions" && "w-12 text-right"
                          )}
                        >
                          <table.FlexRender cell={cell} />
                        </TableCell>
                      ))}
                    </TableRow>
                  )
                })
              ) : (
                <TableRow>
                  <TableCell
                    colSpan={tableColumns.length}
                    className="h-64 text-center"
                  >
                    <div className="flex flex-col items-center justify-center gap-2 max-w-sm mx-auto text-center py-6">
                      <div className="rounded-full bg-muted p-3 text-muted-foreground">
                        <IconInbox className="h-6 w-6" />
                      </div>
                      <h3 className="font-semibold text-sm text-foreground">
                        {activeFiltersCount > 0 ? "No matching records found" : emptyTitle}
                      </h3>
                      <p className="text-xs text-muted-foreground">
                        {activeFiltersCount > 0
                          ? "Try adjusting your search criteria or resetting active filters to find records."
                          : emptyDescription}
                      </p>
                      {activeFiltersCount > 0 ? (
                        <Button
                          variant="outline"
                          size="sm"
                          onClick={clearAllFilters}
                          className="mt-2 h-8 text-xs"
                        >
                          Reset Filters
                        </Button>
                      ) : emptyAction ? (
                        <Button
                          variant="default"
                          size="sm"
                          onClick={emptyAction.onClick}
                          className="mt-2 h-8 text-xs"
                        >
                          {emptyAction.label}
                        </Button>
                      ) : null}
                    </div>
                  </TableCell>
                </TableRow>
              )}
            </TableBody>
          </Table>
        </div>
      </div>

      <div className="flex flex-col gap-3 sm:flex-row sm:items-center sm:justify-between px-1">
        <div className="text-xs text-muted-foreground flex items-center gap-2">
          {totalFilteredCount > 0 ? (
            <span>
              Showing{" "}
              <strong className="font-semibold text-foreground">
                {table.state.pagination.pageIndex * table.state.pagination.pageSize + 1}
              </strong>
              –
              <strong className="font-semibold text-foreground">
                {Math.min(
                  (table.state.pagination.pageIndex + 1) * table.state.pagination.pageSize,
                  totalFilteredCount
                )}
              </strong>{" "}
              of <strong className="font-semibold text-foreground">{totalFilteredCount}</strong>{" "}
              records
              {totalFilteredCount !== totalRowCount && ` (filtered from ${totalRowCount})`}
            </span>
          ) : (
            <span>0 records</span>
          )}
        </div>

        <div className="flex items-center gap-4 self-end sm:self-auto">
          <div className="flex items-center gap-2">
            <span className="text-xs text-muted-foreground hidden sm:inline">Rows per page:</span>
            <Select
              value={`${table.state.pagination.pageSize}`}
              onValueChange={(val) => table.setPageSize(Number(val))}
            >
              <SelectTrigger className="h-8 w-[72px] text-xs">
                <SelectValue placeholder={table.state.pagination.pageSize} />
              </SelectTrigger>
              <SelectContent side="top">
                {[10, 25, 50, 100].map((size) => (
                  <SelectItem key={size} value={`${size}`} className="text-xs">
                    {size}
                  </SelectItem>
                ))}
              </SelectContent>
            </Select>
          </div>

          <div className="text-xs font-medium text-muted-foreground">
            Page {table.state.pagination.pageIndex + 1} of {Math.max(1, table.getPageCount())}
          </div>

          <div className="flex items-center space-x-1">
            <Button
              variant="outline"
              size="icon"
              className="h-8 w-8"
              onClick={() => table.setPageIndex(0)}
              disabled={!table.getCanPreviousPage()}
              aria-label="Go to first page"
            >
              <IconChevronsLeft className="h-4 w-4" />
            </Button>
            <Button
              variant="outline"
              size="icon"
              className="h-8 w-8"
              onClick={() => table.previousPage()}
              disabled={!table.getCanPreviousPage()}
              aria-label="Go to previous page"
            >
              <IconChevronLeft className="h-4 w-4" />
            </Button>
            <Button
              variant="outline"
              size="icon"
              className="h-8 w-8"
              onClick={() => table.nextPage()}
              disabled={!table.getCanNextPage()}
              aria-label="Go to next page"
            >
              <IconChevronRight className="h-4 w-4" />
            </Button>
            <Button
              variant="outline"
              size="icon"
              className="h-8 w-8"
              onClick={() => table.setPageIndex(table.getPageCount() - 1)}
              disabled={!table.getCanNextPage()}
              aria-label="Go to last page"
            >
              <IconChevronsRight className="h-4 w-4" />
            </Button>
          </div>
        </div>
      </div>

      {renderDetailDrawer && (
        <Sheet
          open={!!activeDrawerRow}
          onOpenChange={(open) => {
            if (!open) setActiveDrawerRow(null)
          }}
        >
          <SheetContent side="right" className="w-full sm:max-w-md overflow-y-auto">
            <SheetHeader>
              <SheetTitle>
                {activeDrawerRow
                  ? typeof detailDrawerTitle === "function"
                    ? detailDrawerTitle(activeDrawerRow)
                    : detailDrawerTitle
                  : ""}
              </SheetTitle>
              {detailDrawerDescription && activeDrawerRow && (
                <SheetDescription>
                  {typeof detailDrawerDescription === "function"
                    ? detailDrawerDescription(activeDrawerRow)
                    : detailDrawerDescription}
                </SheetDescription>
              )}
            </SheetHeader>
            <div className="py-6">
              {activeDrawerRow && renderDetailDrawer(activeDrawerRow, () => setActiveDrawerRow(null))}
            </div>
          </SheetContent>
        </Sheet>
      )}

      <AlertDialog
        open={!!pendingBulkAction}
        onOpenChange={(open) => {
          if (!open) setPendingBulkAction(null)
        }}
      >
        <AlertDialogContent>
          <AlertDialogHeader>
            <AlertDialogTitle>
              {pendingBulkAction?.action.confirmationTitle || "Confirm Bulk Action"}
            </AlertDialogTitle>
            <AlertDialogDescription>
              {pendingBulkAction?.action.confirmationDescription
                ? pendingBulkAction.action.confirmationDescription(pendingBulkAction.count)
                : `Are you sure you want to perform this action on ${pendingBulkAction?.count} selected record(s)? This action cannot be undone.`}
            </AlertDialogDescription>
          </AlertDialogHeader>
          <AlertDialogFooter>
            <AlertDialogCancel disabled={isBulkProcessing}>Cancel</AlertDialogCancel>
            <AlertDialogAction
              disabled={isBulkProcessing}
              className={cn(
                pendingBulkAction?.action.variant === "destructive" &&
                  "bg-destructive text-destructive-foreground hover:bg-destructive/90"
              )}
              onClick={(e) => {
                e.preventDefault()
                if (pendingBulkAction) {
                  executeBulkAction(pendingBulkAction.action)
                }
              }}
            >
              {isBulkProcessing ? "Processing..." : "Confirm"}
            </AlertDialogAction>
          </AlertDialogFooter>
        </AlertDialogContent>
      </AlertDialog>
    </div>
  )
}
