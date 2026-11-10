import * as React from 'react'
import {
  Card,
  CardContent,
  Button,
  Badge,
  Checkbox,
  DataTable,
  toast,
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogFooter,
  cn,
} from '@repo/ui'
import { createColumnHelper } from '@tanstack/react-table'
import type {
  DataTableFeatures,
  DataTableFilterConfig,
  DataTableBulkAction,
  DataTableRowAction,
} from '@repo/ui'
import {
  IconTable,
  IconLayoutKanban,
  IconEdit,
  IconTrash,
  IconEye,
} from '@tabler/icons-react'
import type { Order, OrderStage, OrderPaymentStatus } from '../types/order'
import { OrderKanbanView } from '../components/OrderKanbanView'

interface OrdersPageProps {
  orders: Order[]
  onAddNewOrder: () => void
  onEditOrder: (order: Order) => void
  onDeleteOrder: (orderId: string) => void
  onStageChange: (orderId: string, newStage: OrderStage) => void
  onViewOrder?: (order: Order) => void
}

const columnHelper = createColumnHelper<DataTableFeatures, Order>()

const STAGES: (OrderStage | 'All')[] = [
  'All',
  'Order Confirmation',
  'Product Ready',
  'Ready For Delivery',
  'Product Delivered',
]

export function OrdersPage({
  orders,
  onAddNewOrder,
  onEditOrder,
  onDeleteOrder,
  onStageChange,
}: OrdersPageProps) {
  const [viewMode, setViewMode] = React.useState<'table' | 'kanban'>('table')
  const [stageFilter, setStageFilter] = React.useState<OrderStage | 'All'>('All')
  const [paymentFilter] = React.useState<OrderPaymentStatus | 'All'>('All')

  // Selected Order for detail modal
  const [selectedOrder, setSelectedOrder] = React.useState<Order | null>(null)

  // Filtered orders
  const filteredOrders = React.useMemo(() => {
    return orders.filter((o) => {
      if (stageFilter !== 'All' && o.stage !== stageFilter) return false
      if (paymentFilter !== 'All' && o.paymentStatus !== paymentFilter) return false
      return true
    })
  }, [orders, stageFilter, paymentFilter])

  // Top KPI metric widgets as seen in Screenshot 1:
  // Count | Pending (₹) | Paid (₹) | Total Orders (₹)
  const metrics = React.useMemo(() => {
    const count = orders.length
    const totalOrdersAmount = orders.reduce((sum, o) => sum + o.totalAmount, 0)
    const paidAmount = orders.reduce((sum, o) => sum + o.paidAmount, 0)
    const pendingAmount = orders.reduce((sum, o) => sum + o.balanceAmount, 0)
    return { count, pendingAmount, paidAmount, totalOrdersAmount }
  }, [orders])

  // Table Columns (Exactly as requested by user):
  // ID, Client Name, Orders Type, Stage, Amount, Paid Amount, Assignee, Balance Amount, Orders Date, Payment Status, Created Date
  const columns = React.useMemo(
    () =>
      columnHelper.columns([
        columnHelper.display({
          id: 'select',
          header: ({ table }) => (
            <Checkbox
              checked={
                table.getIsAllPageRowsSelected() ||
                (table.getIsSomePageRowsSelected() && 'indeterminate')
              }
              onCheckedChange={(value: boolean | 'indeterminate') =>
                table.toggleAllPageRowsSelected(!!value)
              }
              aria-label="Select all"
            />
          ),
          cell: ({ row }) => (
            <Checkbox
              checked={row.getIsSelected()}
              onCheckedChange={(value: boolean | 'indeterminate') =>
                row.toggleSelected(!!value)
              }
              aria-label="Select row"
            />
          ),
          enableSorting: false,
          enableHiding: false,
        }),
        columnHelper.accessor('orderNumber', {
          header: 'ID',
          cell: (info) => (
            <button
              type="button"
              onClick={() => setSelectedOrder(info.row.original)}
              className="font-mono text-xs font-bold text-primary hover:underline hover:text-primary/80 tracking-tight"
            >
              {info.getValue()}
            </button>
          ),
          meta: { filterKey: 'orderNumber' },
        }),
        columnHelper.accessor('clientName', {
          header: 'Client Name',
          cell: (info) => {
            const order = info.row.original
            return (
              <div className="flex flex-col">
                <span className="font-semibold text-xs text-foreground">{info.getValue()}</span>
                <span className="text-[11px] text-muted-foreground">{order.mobileNo}</span>
              </div>
            )
          },
          meta: { filterKey: 'clientName' },
        }),
        columnHelper.accessor('orderType', {
          header: 'Orders Type',
          cell: (info) => (
            <span className="inline-flex items-center px-2 py-0.5 rounded-md text-[11px] font-medium bg-muted/50 text-foreground/80 border border-border/50">
              {info.getValue()}
            </span>
          ),
          meta: { filterKey: 'orderType' },
        }),
        columnHelper.accessor('stage', {
          header: 'Stage',
          cell: (info) => {
            const stage = info.getValue()
            const stageConfig: Record<string, { bg: string; text: string; border: string; dot: string }> = {
              'Order Confirmation': {
                bg: 'bg-sky-500/10',
                text: 'text-sky-400',
                border: 'border-sky-500/20',
                dot: 'bg-sky-400',
              },
              'Product Ready': {
                bg: 'bg-amber-500/10',
                text: 'text-amber-400',
                border: 'border-amber-500/20',
                dot: 'bg-amber-400',
              },
              'Ready For Delivery': {
                bg: 'bg-violet-500/10',
                text: 'text-violet-400',
                border: 'border-violet-500/20',
                dot: 'bg-violet-400',
              },
              'Product Delivered': {
                bg: 'bg-emerald-500/10',
                text: 'text-emerald-400',
                border: 'border-emerald-500/20',
                dot: 'bg-emerald-400',
              },
            }

            const current = stageConfig[stage] || {
              bg: 'bg-slate-500/10',
              text: 'text-slate-400',
              border: 'border-slate-500/20',
              dot: 'bg-slate-400',
            }

            return (
              <span
                className={cn(
                  'inline-flex items-center gap-1.5 px-2.5 py-0.5 rounded-full text-[11px] font-medium tracking-tight border shadow-2xs backdrop-blur-xs transition-colors',
                  current.bg,
                  current.text,
                  current.border
                )}
              >
                <span className={cn('h-1.5 w-1.5 rounded-full shrink-0', current.dot)} />
                <span>{stage}</span>
              </span>
            )
          },
          meta: { filterKey: 'stage' },
        }),
        columnHelper.accessor('totalAmount', {
          header: 'Amount',
          cell: (info) => (
            <span className="font-semibold text-xs text-foreground">
              ₹{info.getValue().toLocaleString('en-IN')}
            </span>
          ),
        }),
        columnHelper.accessor('paidAmount', {
          header: 'Paid Amount',
          cell: (info) => (
            <span className="text-xs text-emerald-600 dark:text-emerald-400 font-medium">
              ₹{info.getValue().toLocaleString('en-IN')}
            </span>
          ),
        }),
        columnHelper.accessor('assignee', {
          header: 'Assignee',
          cell: (info) => (
            <span className="text-xs text-muted-foreground">{info.getValue() || '—'}</span>
          ),
        }),
        columnHelper.accessor('balanceAmount', {
          header: 'Balance Amount',
          cell: (info) => {
            const val = info.getValue()
            return (
              <span
                className={cn(
                  'font-semibold text-xs',
                  val > 0
                    ? 'text-amber-600 dark:text-amber-400'
                    : 'text-emerald-600 dark:text-emerald-400'
                )}
              >
                ₹{val.toLocaleString('en-IN')}
              </span>
            )
          },
        }),
        columnHelper.accessor('orderDate', {
          header: 'Orders Date',
          cell: (info) => (
            <span className="text-xs text-muted-foreground font-mono">{info.getValue()}</span>
          ),
        }),
        columnHelper.accessor('paymentStatus', {
          header: 'Payment Status',
          cell: (info) => {
            const status = info.getValue()
            return (
              <Badge
                variant={
                  status === 'Paid'
                    ? 'success'
                    : status === 'Partially Paid'
                      ? 'warning'
                      : 'secondary'
                }
                className="text-[10px] px-2 py-0.5 font-medium"
              >
                {status}
              </Badge>
            )
          },
        }),
        columnHelper.accessor('createdDate', {
          header: 'Created Date',
          cell: (info) => (
            <span className="text-xs text-muted-foreground">{info.getValue()}</span>
          ),
        }),
      ]),
    []
  )

  // Filters for tanstack table
  const filterConfigs: DataTableFilterConfig[] = React.useMemo(
    () => [
      {
        columnId: 'orderType',
        title: 'Order Type',
        type: 'select',
        options: [
          { label: 'Walk-in', value: 'Walk-in' },
          { label: 'Studio Session', value: 'Studio Session' },
          { label: 'Print & Album', value: 'Print & Album' },
          { label: 'Framing', value: 'Framing' },
          { label: 'Custom Merch', value: 'Custom Merch' },
        ],
      },
      {
        columnId: 'stage',
        title: 'Stage',
        type: 'select',
        options: STAGES.filter((s) => s !== 'All').map((s) => ({ label: s, value: s })),
      },
    ],
    []
  )

  const rowActions: DataTableRowAction<Order>[] = React.useMemo(
    () => [
      {
        label: 'View Details',
        icon: IconEye,
        action: (order: Order) => setSelectedOrder(order),
      },
      {
        label: 'Edit Order',
        icon: IconEdit,
        action: (order: Order) => onEditOrder(order),
      },
      {
        label: 'Delete Order',
        icon: IconTrash,
        variant: 'destructive',
        action: (order: Order) => onDeleteOrder(order.id),
      },
    ],
    [onEditOrder, onDeleteOrder]
  )

  const bulkActions: DataTableBulkAction<Order>[] = React.useMemo(
    () => [
      {
        label: 'Mark as Product Ready',
        action: (selectedRows: Order[], clearSelection: () => void) => {
          selectedRows.forEach((o) => onStageChange(o.id, 'Product Ready'))
          clearSelection()
          toast.success(`Updated ${selectedRows.length} orders to Product Ready`)
        },
      },
      {
        label: 'Mark as Ready For Delivery',
        action: (selectedRows: Order[], clearSelection: () => void) => {
          selectedRows.forEach((o) => onStageChange(o.id, 'Ready For Delivery'))
          clearSelection()
          toast.success(`Updated ${selectedRows.length} orders to Ready For Delivery`)
        },
      },
      {
        label: 'Mark as Product Delivered',
        action: (selectedRows: Order[], clearSelection: () => void) => {
          selectedRows.forEach((o) => onStageChange(o.id, 'Product Delivered'))
          clearSelection()
          toast.success(`Updated ${selectedRows.length} orders to Product Delivered`)
        },
      },
    ],
    [onStageChange]
  )

  return (
    <div className="flex flex-col gap-5">
      {/* Top Bar Match (Screenshot 1): List/Board Toggle + Filters + KPIs + New Orders Button */}
      <div className="flex flex-wrap items-center justify-between gap-4 bg-card p-3 rounded-xl border shadow-2xs">
        {/* Left: View Mode Toggle & Quick Filter Pills */}
        <div className="flex items-center gap-3">
          <div className="flex items-center rounded-lg border bg-muted/30 p-0.5">
            <button
              type="button"
              onClick={() => setViewMode('table')}
              className={cn(
                'flex items-center gap-1.5 px-3 py-1.5 text-xs font-semibold rounded-md transition',
                viewMode === 'table'
                  ? 'bg-foreground text-background shadow-xs'
                  : 'text-muted-foreground hover:text-foreground'
              )}
            >
              <IconTable className="h-3.5 w-3.5" />
              List
            </button>
            <button
              type="button"
              onClick={() => setViewMode('kanban')}
              className={cn(
                'flex items-center gap-1.5 px-3 py-1.5 text-xs font-semibold rounded-md transition',
                viewMode === 'kanban'
                  ? 'bg-foreground text-background shadow-xs'
                  : 'text-muted-foreground hover:text-foreground'
              )}
            >
              <IconLayoutKanban className="h-3.5 w-3.5" />
              Board
            </button>
          </div>

          {/* Quick Stage Filter Badges */}
          <div className="hidden md:flex items-center gap-1 border-l pl-3">
            {STAGES.map((s) => (
              <button
                key={s}
                type="button"
                onClick={() => setStageFilter(s)}
                className={cn(
                  'px-2.5 py-1 text-[11px] rounded-full transition font-medium',
                  stageFilter === s
                    ? 'bg-primary/15 text-primary border border-primary/30 font-semibold'
                    : 'text-muted-foreground hover:bg-muted/50'
                )}
              >
                {s}
              </button>
            ))}
          </div>
        </div>

        {/* Right: KPIs (Count, Pending, Paid, Total Orders) + New Orders Button */}
        <div className="flex items-center gap-6">
          <div className="flex items-center gap-5 text-right">
            <div className="flex flex-col">
              <span className="text-sm font-bold text-foreground">{metrics.count}</span>
              <span className="text-[11px] font-medium text-muted-foreground">Count</span>
            </div>
            <div className="flex flex-col">
              <span className="text-sm font-bold text-amber-600 dark:text-amber-400">
                ₹{metrics.pendingAmount.toLocaleString('en-IN')}
              </span>
              <span className="text-[11px] font-medium text-muted-foreground">Pending</span>
            </div>
            <div className="flex flex-col">
              <span className="text-sm font-bold text-emerald-600 dark:text-emerald-400">
                ₹{metrics.paidAmount.toLocaleString('en-IN')}
              </span>
              <span className="text-[11px] font-medium text-muted-foreground">Paid</span>
            </div>
            <div className="flex flex-col">
              <span className="text-sm font-bold text-foreground">
                ₹{metrics.totalOrdersAmount.toLocaleString('en-IN')}
              </span>
              <span className="text-[11px] font-medium text-muted-foreground">Total Orders</span>
            </div>
          </div>

          <Button
            type="button"
            onClick={onAddNewOrder}
            className="bg-emerald-600 hover:bg-emerald-700 text-white font-semibold text-xs px-4 h-9 shadow-xs"
          >
            New Orders
          </Button>
        </div>
      </div>

      {/* Main View: List View (DataTable) or Board View (Kanban) */}
      {viewMode === 'table' ? (
        <Card className="shadow-2xs border-border/80">
          <CardContent className="p-0">
            <DataTable
              data={filteredOrders}
              columns={columns}
              rowActions={rowActions}
              bulkActions={bulkActions}
              filterConfigs={filterConfigs}
              searchableColumns={['orderNumber', 'clientName', 'orderType']}
              searchPlaceholder="Search orders by ID, client or type..."
            />
          </CardContent>
        </Card>
      ) : (
        <OrderKanbanView
          orders={filteredOrders}
          onOrderClick={(order) => setSelectedOrder(order)}
          onEditOrder={onEditOrder}
          onDeleteOrder={onDeleteOrder}
          onStageChange={onStageChange}
        />
      )}

      {/* Order Detail Modal */}
      {selectedOrder && (
        <Dialog open={!!selectedOrder} onOpenChange={(open) => !open && setSelectedOrder(null)}>
          <DialogContent className="max-w-2xl">
            <DialogHeader>
              <div className="flex items-center justify-between pr-6">
                <div>
                  <DialogTitle className="text-base font-bold flex items-center gap-2">
                    <span>{selectedOrder.orderNumber}</span>
                    <Badge variant="outline" className="text-xs">
                      {selectedOrder.stage}
                    </Badge>
                  </DialogTitle>
                  <p className="text-xs text-muted-foreground mt-0.5">
                    Order Date: {selectedOrder.orderDate} • Assignee: {selectedOrder.assignee}
                  </p>
                </div>
                <Badge
                  variant={
                    selectedOrder.paymentStatus === 'Paid'
                      ? 'success'
                      : selectedOrder.paymentStatus === 'Partially Paid'
                        ? 'warning'
                        : 'secondary'
                  }
                  className="text-xs"
                >
                  {selectedOrder.paymentStatus}
                </Badge>
              </div>
            </DialogHeader>

            <div className="space-y-4 py-2 text-xs">
              {/* Client Info Banner */}
              <div className="rounded-lg bg-muted/30 p-3 grid grid-cols-2 gap-3 border">
                <div>
                  <span className="text-[10px] text-muted-foreground font-semibold uppercase">
                    Client Name
                  </span>
                  <p className="font-semibold text-foreground">{selectedOrder.clientName}</p>
                </div>
                <div>
                  <span className="text-[10px] text-muted-foreground font-semibold uppercase">
                    Phone & Email
                  </span>
                  <p className="font-medium text-foreground">
                    {selectedOrder.mobileNo} {selectedOrder.email ? `• ${selectedOrder.email}` : ''}
                  </p>
                </div>
                {selectedOrder.address && (
                  <div className="col-span-2">
                    <span className="text-[10px] text-muted-foreground font-semibold uppercase">
                      Delivery Address
                    </span>
                    <p className="text-foreground">{selectedOrder.address}</p>
                  </div>
                )}
              </div>

              {/* Items List */}
              <div className="space-y-2">
                <span className="font-bold text-foreground text-xs">Order Line Items</span>
                <div className="border rounded-lg overflow-hidden">
                  <table className="w-full text-xs">
                    <thead className="bg-muted/40 border-b text-muted-foreground text-left">
                      <tr>
                        <th className="p-2 font-medium">Item Description</th>
                        <th className="p-2 font-medium text-center w-16">Qty</th>
                        <th className="p-2 font-medium text-right w-24">Price</th>
                        <th className="p-2 font-medium text-right w-24">Subtotal</th>
                      </tr>
                    </thead>
                    <tbody className="divide-y">
                      {selectedOrder.items.map((it) => (
                        <tr key={it.id}>
                          <td className="p-2 font-medium text-foreground">{it.name}</td>
                          <td className="p-2 text-center text-muted-foreground">{it.quantity}</td>
                          <td className="p-2 text-right text-muted-foreground">
                            ₹{it.itemCost.toLocaleString('en-IN')}
                          </td>
                          <td className="p-2 text-right font-semibold text-foreground">
                            ₹{it.subTotal.toLocaleString('en-IN')}
                          </td>
                        </tr>
                      ))}
                    </tbody>
                  </table>
                </div>
              </div>

              {/* Order Expenses */}
              {selectedOrder.expenses.length > 0 && (
                <div className="space-y-2">
                  <span className="font-bold text-foreground text-xs">Order Expenses / Costs</span>
                  <div className="border rounded-lg overflow-hidden bg-muted/10">
                    <table className="w-full text-xs">
                      <tbody className="divide-y">
                        {selectedOrder.expenses.map((ex) => (
                          <tr key={ex.id}>
                            <td className="p-2 text-muted-foreground">{ex.date}</td>
                            <td className="p-2 font-medium text-foreground">{ex.name}</td>
                            <td className="p-2 text-right font-semibold text-muted-foreground">
                              ₹{ex.amount.toLocaleString('en-IN')}
                            </td>
                          </tr>
                        ))}
                      </tbody>
                    </table>
                  </div>
                </div>
              )}

              {/* Financial Summary */}
              <div className="flex justify-between items-start pt-2 border-t">
                <div className="space-y-1">
                  {selectedOrder.notes && (
                    <p className="text-muted-foreground italic text-[11px] max-w-xs">
                      Note: {selectedOrder.notes}
                    </p>
                  )}
                </div>

                <div className="w-56 space-y-1.5 text-xs">
                  <div className="flex justify-between">
                    <span className="text-muted-foreground">Subtotal</span>
                    <span className="font-semibold">
                      ₹{selectedOrder.subTotal.toLocaleString('en-IN')}
                    </span>
                  </div>
                  {selectedOrder.discount > 0 && (
                    <div className="flex justify-between text-emerald-600">
                      <span>Discount</span>
                      <span>-₹{selectedOrder.discount.toLocaleString('en-IN')}</span>
                    </div>
                  )}
                  {selectedOrder.taxAmount > 0 && (
                    <div className="flex justify-between text-muted-foreground">
                      <span>Tax ({selectedOrder.taxPercent}%)</span>
                      <span>+₹{selectedOrder.taxAmount.toLocaleString('en-IN')}</span>
                    </div>
                  )}
                  <div className="flex justify-between font-bold text-sm pt-1 border-t">
                    <span>Total Amount</span>
                    <span>₹{selectedOrder.totalAmount.toLocaleString('en-IN')}</span>
                  </div>
                  <div className="flex justify-between text-muted-foreground">
                    <span>Paid Amount</span>
                    <span className="text-emerald-600 font-semibold">
                      ₹{selectedOrder.paidAmount.toLocaleString('en-IN')}
                    </span>
                  </div>
                  <div className="flex justify-between font-bold">
                    <span>Balance Due</span>
                    <span className="text-amber-600">
                      ₹{selectedOrder.balanceAmount.toLocaleString('en-IN')}
                    </span>
                  </div>
                </div>
              </div>
            </div>

            <DialogFooter className="gap-2">
              <Button
                type="button"
                variant="outline"
                size="sm"
                onClick={() => setSelectedOrder(null)}
              >
                Close
              </Button>
              <Button
                type="button"
                size="sm"
                onClick={() => {
                  const ord = selectedOrder
                  setSelectedOrder(null)
                  onEditOrder(ord)
                }}
                className="bg-primary text-primary-foreground"
              >
                Edit Order
              </Button>
            </DialogFooter>
          </DialogContent>
        </Dialog>
      )}
    </div>
  )
}
