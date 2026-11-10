import * as React from 'react'
import {
  IconDotsVertical,
  IconEye,
  IconEdit,
  IconTrash,
  IconUserCheck,
} from '@tabler/icons-react'
import { Badge } from '@repo/ui/components/ui/badge'
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from '@repo/ui/components/ui/dropdown-menu'
import { Button } from '@repo/ui/components/ui/button'
import { cn } from '@repo/ui/lib/utils'
import type { Order, OrderStage } from '../types/order'
import { useWorkflows } from '../hooks/useWorkflows'

interface OrderKanbanViewProps {
  orders: Order[]
  onOrderClick: (order: Order) => void
  onEditOrder: (order: Order) => void
  onDeleteOrder: (orderId: string) => void
  onStageChange: (orderId: string, newStage: OrderStage) => void
}

interface StageColumn {
  id: OrderStage
  title: string
  colorBadge: string
  accentBorder: string
}

const STAGE_COLUMNS: StageColumn[] = [
  {
    id: 'Order Confirmation',
    title: 'Order Confirmation',
    colorBadge: 'bg-blue-500/10 text-blue-700 dark:text-blue-300 border-blue-500/20',
    accentBorder: 'border-t-blue-500',
  },
  {
    id: 'Product Ready',
    title: 'Product Ready',
    colorBadge: 'bg-amber-500/10 text-amber-700 dark:text-amber-300 border-amber-500/20',
    accentBorder: 'border-t-amber-500',
  },
  {
    id: 'Ready For Delivery',
    title: 'Ready For Delivery',
    colorBadge: 'bg-purple-500/10 text-purple-700 dark:text-purple-300 border-purple-500/20',
    accentBorder: 'border-t-purple-500',
  },
  {
    id: 'Product Delivered',
    title: 'Product Delivered',
    colorBadge: 'bg-emerald-500/10 text-emerald-700 dark:text-emerald-300 border-emerald-500/20',
    accentBorder: 'border-t-emerald-500',
  },
]

export function OrderKanbanView({
  orders,
  onOrderClick,
  onEditOrder,
  onDeleteOrder,
  onStageChange,
}: OrderKanbanViewProps) {
  const { workflows } = useWorkflows('order')
  const [draggedOrderId, setDraggedOrderId] = React.useState<string | null>(null)
  const [dragOverColumn, setDragOverColumn] = React.useState<string | null>(null)

  const kanbanColumns = React.useMemo(() => {
    if (workflows && workflows.length > 0) {
      return workflows.map((w) => ({
        id: w.stage as OrderStage,
        title: w.stage,
        color: w.color_code,
      }))
    }
    return STAGE_COLUMNS.map((s) => ({
      id: s.id,
      title: s.title,
      color: '#3b82f6',
    }))
  }, [workflows])

  const handleDragStart = (e: React.DragEvent, orderId: string) => {
    e.dataTransfer.setData('text/plain', orderId)
    e.dataTransfer.effectAllowed = 'move'
    setDraggedOrderId(orderId)
  }

  const handleDragEnd = () => {
    setDraggedOrderId(null)
    setDragOverColumn(null)
  }

  const handleDragOver = (e: React.DragEvent, stage: string) => {
    e.preventDefault()
    e.dataTransfer.dropEffect = 'move'
    if (dragOverColumn !== stage) {
      setDragOverColumn(stage)
    }
  }

  const handleDrop = (e: React.DragEvent, stage: string) => {
    e.preventDefault()
    const orderId = e.dataTransfer.getData('text/plain') || draggedOrderId
    if (orderId) {
      onStageChange(orderId, stage as OrderStage)
    }
    setDraggedOrderId(null)
    setDragOverColumn(null)
  }

  return (
    <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4 pb-6">
      {kanbanColumns.map((col) => {
        const columnOrders = orders.filter((o) => o.stage === col.id)
        const totalStageAmount = columnOrders.reduce((sum, o) => sum + o.totalAmount, 0)
        const isTarget = dragOverColumn === col.id

        return (
          <div
            key={col.id}
            onDragOver={(e) => handleDragOver(e, col.id)}
            onDragLeave={() => setDragOverColumn(null)}
            onDrop={(e) => handleDrop(e, col.id)}
            style={{ borderTopColor: col.color }}
            className={cn(
              'flex flex-col rounded-xl border border-border/80 bg-muted/20 p-3 min-h-[520px] transition-colors border-t-4',
              isTarget && 'bg-primary/5 border-primary/40 ring-2 ring-primary/20'
            )}
          >
            {/* Column Header */}
            <div className="flex items-center justify-between pb-2 border-b border-border/40 mb-3">
              <div className="flex items-center gap-2">
                <span className="font-semibold text-xs text-foreground tracking-tight">
                  {col.title}
                </span>
                <span
                  style={{
                    backgroundColor: `${col.color}20`,
                    borderColor: `${col.color}40`,
                    color: col.color,
                  }}
                  className="inline-flex items-center justify-center rounded-full px-2 py-0.5 text-[11px] font-bold border"
                >
                  {columnOrders.length}
                </span>
              </div>
              <span className="text-[11px] font-medium text-muted-foreground">
                ₹{totalStageAmount.toLocaleString('en-IN')}
              </span>
            </div>

            {/* Column Cards */}
            <div className="flex flex-col gap-2.5 flex-1 overflow-y-auto">
              {columnOrders.length === 0 ? (
                <div className="flex flex-1 items-center justify-center rounded-lg border border-dashed border-border/60 p-6 text-center text-xs text-muted-foreground">
                  No orders in this stage
                </div>
              ) : (
                columnOrders.map((order) => {
                  const isDragging = draggedOrderId === order.id

                  return (
                    <div
                      key={order.id}
                      draggable
                      onDragStart={(e) => handleDragStart(e, order.id)}
                      onDragEnd={handleDragEnd}
                      onClick={() => onOrderClick(order)}
                      className={cn(
                        'group relative flex flex-col gap-2 rounded-lg border border-border bg-card p-3 shadow-2xs hover:shadow-md hover:border-primary/40 transition cursor-grab active:cursor-grabbing',
                        isDragging && 'opacity-40 scale-95'
                      )}
                    >
                      {/* Card Header */}
                      <div className="flex items-start justify-between gap-2">
                        <div className="flex flex-col min-w-0">
                          <span className="text-[11px] font-bold text-primary tracking-tight">
                            {order.orderNumber}
                          </span>
                          <span className="font-semibold text-xs text-foreground truncate">
                            {order.clientName}
                          </span>
                        </div>

                        <div className="flex items-center" onClick={(e) => e.stopPropagation()}>
                          <DropdownMenu>
                            <DropdownMenuTrigger asChild>
                              <Button
                                variant="ghost"
                                size="sm"
                                className="h-6 w-6 p-0 opacity-40 group-hover:opacity-100"
                              >
                                <IconDotsVertical className="h-3.5 w-3.5" />
                              </Button>
                            </DropdownMenuTrigger>
                            <DropdownMenuContent align="end" className="w-36 text-xs">
                              <DropdownMenuItem onClick={() => onOrderClick(order)}>
                                <IconEye className="h-3.5 w-3.5 mr-2 text-muted-foreground" />
                                View Details
                              </DropdownMenuItem>
                              <DropdownMenuItem onClick={() => onEditOrder(order)}>
                                <IconEdit className="h-3.5 w-3.5 mr-2 text-muted-foreground" />
                                Edit Order
                              </DropdownMenuItem>
                              <DropdownMenuItem
                                onClick={() => onDeleteOrder(order.id)}
                                className="text-destructive focus:text-destructive"
                              >
                                <IconTrash className="h-3.5 w-3.5 mr-2 text-destructive" />
                                Delete
                              </DropdownMenuItem>
                            </DropdownMenuContent>
                          </DropdownMenu>
                        </div>
                      </div>

                      {/* Type and Items Badge */}
                      <div className="flex flex-wrap items-center gap-1.5 text-[11px]">
                        <span className="inline-flex items-center px-1.5 py-0.5 rounded bg-muted text-muted-foreground font-medium text-[10px]">
                          {order.orderType}
                        </span>
                        <span className="text-muted-foreground text-[10px]">
                          • {order.items.length} item{order.items.length !== 1 ? 's' : ''}
                        </span>
                      </div>

                      {/* Financial info */}
                      <div className="mt-1 pt-2 border-t border-border/50 flex items-center justify-between text-xs">
                        <div className="flex flex-col">
                          <span className="text-[10px] text-muted-foreground">Total</span>
                          <span className="font-semibold text-foreground">
                            ₹{order.totalAmount.toLocaleString('en-IN')}
                          </span>
                        </div>

                        <div className="flex flex-col items-end">
                          <span className="text-[10px] text-muted-foreground">Balance</span>
                          <span
                            className={cn(
                              'font-semibold text-xs',
                              order.balanceAmount > 0
                                ? 'text-amber-600 dark:text-amber-400'
                                : 'text-emerald-600 dark:text-emerald-400'
                            )}
                          >
                            ₹{order.balanceAmount.toLocaleString('en-IN')}
                          </span>
                        </div>
                      </div>

                      {/* Card Footer: Assignee & Payment Status */}
                      <div className="flex items-center justify-between pt-1 text-[10px] text-muted-foreground">
                        <span className="truncate max-w-[100px] flex items-center gap-1">
                          <IconUserCheck className="h-3 w-3 shrink-0" />
                          {order.assignee}
                        </span>
                        <Badge
                          variant={
                            order.paymentStatus === 'Paid'
                              ? 'success'
                              : order.paymentStatus === 'Partially Paid'
                              ? 'warning'
                              : 'secondary'
                          }
                          className="text-[9px] px-1.5 py-0"
                        >
                          {order.paymentStatus}
                        </Badge>
                      </div>
                    </div>
                  )
                })
              )}
            </div>
          </div>
        )
      })}
    </div>
  )
}

