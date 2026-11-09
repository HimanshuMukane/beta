import * as React from 'react';
import {
  KanbanBoard,
  KanbanBoardProvider,
  KanbanBoardColumn,
  KanbanBoardColumnHeader,
  KanbanBoardColumnTitle,
  KanbanBoardColumnList,
  KanbanBoardColumnListItem,
  KanbanBoardCard,
  KanbanColorCircle,
  KanbanBoardExtraMargin,
  type KanbanBoardCircleColor,
  type KanbanBoardDropDirection,
} from './kanban-primitive';
import { Button } from '../ui/button';
import { Input } from '../ui/input';
import { Badge } from '../ui/badge';
import { Avatar, AvatarFallback, AvatarImage } from '../ui/avatar';
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from '../ui/dropdown-menu';
import {
  Sheet,
  SheetContent,
  SheetHeader,
  SheetTitle,
  SheetDescription,
  SheetFooter,
} from '../ui/sheet';
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogDescription,
  DialogFooter,
} from '../ui/dialog';
import {
  IconSearch,
  IconPlus,
  IconDotsVertical,
  IconX,
  IconFilter,
  IconSortAscending,
  IconSortDescending,
  IconLayoutKanban,
  IconTable,
  IconCalendar,
  IconTrash,
  IconEdit,
  IconCopy,
  IconPaperclip,
  IconMessageCircle,
  IconUser,
} from '@tabler/icons-react';
import { cn } from '@repo/ui/lib/utils';
import { toast } from 'sonner';

export interface KanbanTask {
  id: string;
  title: string;
  status: string;
  customerName?: string;
  dealValue?: number;
  priority?: 'low' | 'medium' | 'high' | 'urgent';
  assignee?: {
    name: string;
    avatar?: string;
    initials?: string;
  };
  dueDate?: string;
  tags?: string[];
  commentsCount?: number;
  attachmentsCount?: number;
  checklistProgress?: {
    completed: number;
    total: number;
  };
  description?: string;
  [key: string]: unknown;
}

export interface KanbanColumnDef {
  id: string;
  title: string;
  color?: KanbanBoardCircleColor;
}

export interface KanbanBoardProps {
  columns: KanbanColumnDef[];
  tasks: KanbanTask[];
  onTasksChange?: (tasks: KanbanTask[]) => void;
  onTaskMove?: (taskId: string, targetStatus: string, newIndex: number) => Promise<boolean> | boolean | void;
  onTaskClick?: (task: KanbanTask) => void;
  onAddTask?: (status: string) => void;
  onDeleteTask?: (taskId: string) => void;
  title?: string;
  description?: string;
  viewMode?: 'kanban' | 'table';
  onViewModeChange?: (mode: 'kanban' | 'table') => void;
  searchPlaceholder?: string;
  className?: string;
}

export function Kanban({
  columns: initialColumns,
  tasks: initialTasks,
  onTasksChange,
  onTaskMove,
  onTaskClick,
  onAddTask,
  onDeleteTask,
  title,
  description,
  viewMode = 'kanban',
  onViewModeChange,
  searchPlaceholder = 'Search records...',
  className,
}: KanbanBoardProps) {
  const [columns, setColumns] = React.useState<KanbanColumnDef[]>(initialColumns);
  const [tasks, setTasks] = React.useState<KanbanTask[]>(initialTasks);
  const [searchQuery, setSearchQuery] = React.useState('');
  const [priorityFilter, setPriorityFilter] = React.useState<string | null>(null);
  const [sortBy, setSortBy] = React.useState<'title' | 'dealValue' | 'dueDate' | 'none'>('none');
  const [sortDirection, setSortDirection] = React.useState<'asc' | 'desc'>('asc');

  // Detail Drawer state
  const [activeTask, setActiveTask] = React.useState<KanbanTask | null>(null);

  // New task modal state
  const [isNewTaskOpen, setIsNewTaskOpen] = React.useState(false);
  const [newTaskStatus, setNewTaskStatus] = React.useState<string>(columns[0]?.id || '');
  const [newTaskTitle, setNewTaskTitle] = React.useState('');
  const [newTaskCustomer, setNewTaskCustomer] = React.useState('');
  const [newTaskValue, setNewTaskValue] = React.useState<string>('');
  const [newTaskPriority, setNewTaskPriority] = React.useState<KanbanTask['priority']>('medium');

  // Sync external tasks changes
  React.useEffect(() => {
    setTasks(initialTasks);
  }, [initialTasks]);

  // Sync columns
  React.useEffect(() => {
    setColumns(initialColumns);
  }, [initialColumns]);

  // Filter & Search Logic
  const filteredTasks = React.useMemo(() => {
    return tasks.filter((t) => {
      if (priorityFilter && t.priority !== priorityFilter) return false;
      if (!searchQuery) return true;
      const q = searchQuery.toLowerCase();
      return (
        t.title.toLowerCase().includes(q) ||
        (t.customerName && t.customerName.toLowerCase().includes(q)) ||
        (t.assignee?.name && t.assignee.name.toLowerCase().includes(q)) ||
        (t.tags && t.tags.some((tag) => tag.toLowerCase().includes(q))) ||
        t.id.toLowerCase().includes(q)
      );
    });
  }, [tasks, searchQuery, priorityFilter]);

  // Group into column tasks with sorting
  const columnTasksMap = React.useMemo(() => {
    const map = new Map<string, KanbanTask[]>();
    columns.forEach((col) => map.set(col.id, []));

    filteredTasks.forEach((task) => {
      const list = map.get(task.status) || [];
      list.push(task);
      map.set(task.status, list);
    });

    if (sortBy !== 'none') {
      columns.forEach((col) => {
        const list = map.get(col.id) || [];
        list.sort((a, b) => {
          let cmp = 0;
          if (sortBy === 'title') cmp = a.title.localeCompare(b.title);
          else if (sortBy === 'dealValue') cmp = (a.dealValue || 0) - (b.dealValue || 0);
          else if (sortBy === 'dueDate') {
            cmp = (a.dueDate || '').localeCompare(b.dueDate || '');
          }
          return sortDirection === 'asc' ? cmp : -cmp;
        });
        map.set(col.id, list);
      });
    }

    return map;
  }, [columns, filteredTasks, sortBy, sortDirection]);

  // Column aggregate totals (e.g. deals value)
  const columnAggregates = React.useMemo(() => {
    const aggregates = new Map<string, { totalValue: number; count: number }>();
    columns.forEach((col) => {
      const colItems = columnTasksMap.get(col.id) || [];
      const total = colItems.reduce((acc, t) => acc + (t.dealValue || 0), 0);
      aggregates.set(col.id, { totalValue: total, count: colItems.length });
    });
    return aggregates;
  }, [columns, columnTasksMap]);

  // Handle Drag and Drop End
  const handleItemMove = async (
    activeId: string,
    targetColumnId: string,
    targetIndex: number
  ) => {
    const previousTasks = [...tasks];
    const task = tasks.find((t) => t.id === activeId);
    if (!task) return;

    if (task.status === targetColumnId) {
      // Reordering within same column
      const colList = [...(columnTasksMap.get(targetColumnId) || [])];
      const oldIdx = colList.findIndex((t) => t.id === activeId);
      if (oldIdx === -1 || oldIdx === targetIndex) return;

      colList.splice(oldIdx, 1);
      colList.splice(targetIndex, 0, task);

      const otherTasks = tasks.filter((t) => t.status !== targetColumnId);
      const nextTasks = [...otherTasks, ...colList];
      setTasks(nextTasks);
      onTasksChange?.(nextTasks);
      return;
    }

    // Moving across columns: Optimistic update
    const updated = tasks.map((t) =>
      t.id === activeId ? { ...t, status: targetColumnId } : t
    );
    setTasks(updated);
    onTasksChange?.(updated);

    if (onTaskMove) {
      try {
        const ok = await onTaskMove(activeId, targetColumnId, targetIndex);
        if (ok === false) {
          setTasks(previousTasks);
          onTasksChange?.(previousTasks);
          toast.error('Failed to move record, rolled back.');
        } else {
          toast.success(`Moved to ${columns.find((c) => c.id === targetColumnId)?.title || targetColumnId}`);
        }
      } catch {
        setTasks(previousTasks);
        onTasksChange?.(previousTasks);
        toast.error('Network error moving record, rolled back.');
      }
    }
  };

  const handleCreateTask = () => {
    if (!newTaskTitle.trim()) {
      toast.error('Please enter a record title');
      return;
    }

    const newTask: KanbanTask = {
      id: `task-${Date.now().toString(36)}`,
      title: newTaskTitle.trim(),
      status: newTaskStatus,
      customerName: newTaskCustomer.trim() || undefined,
      dealValue: newTaskValue ? parseFloat(newTaskValue) : undefined,
      priority: newTaskPriority,
      dueDate: new Date(Date.now() + 7 * 86400000).toISOString().split('T')[0],
      tags: ['New'],
    };

    const next = [newTask, ...tasks];
    setTasks(next);
    onTasksChange?.(next);
    toast.success('Record created');

    // Reset modal
    setNewTaskTitle('');
    setNewTaskCustomer('');
    setNewTaskValue('');
    setNewTaskPriority('medium');
    setIsNewTaskOpen(false);
  };

  const handleDeleteTask = (taskId: string) => {
    const next = tasks.filter((t) => t.id !== taskId);
    setTasks(next);
    onTasksChange?.(next);
    onDeleteTask?.(taskId);
    if (activeTask?.id === taskId) {
      setActiveTask(null);
    }
    toast.success('Record deleted');
  };

  return (
    <div className={cn('flex flex-col gap-2.5 w-full', className)}>
      {/* 1. Board Toolbar */}
      <div className="flex flex-col gap-2 sm:flex-row sm:items-center sm:justify-between border-b pb-2 px-1">
        <div className="flex items-center gap-3">
          {title && (
            <div>
              <h2 className="text-sm font-semibold tracking-tight">{title}</h2>
              {description && (
                <p className="text-[11px] text-muted-foreground">{description}</p>
              )}
            </div>
          )}
          <div className="relative w-48 sm:w-64">
            <IconSearch className="absolute left-2.5 top-1/2 -translate-y-1/2 h-3.5 w-3.5 text-muted-foreground" />
            <Input
              value={searchQuery}
              onChange={(e) => setSearchQuery(e.target.value)}
              placeholder={searchPlaceholder}
              className="h-8 pl-8 pr-7 text-xs bg-card"
            />
            {searchQuery && (
              <button
                type="button"
                onClick={() => setSearchQuery('')}
                className="absolute right-2 top-1/2 -translate-y-1/2 text-muted-foreground hover:text-foreground"
              >
                <IconX className="h-3 w-3" />
              </button>
            )}
          </div>
        </div>

        <div className="flex items-center gap-1.5 self-end sm:self-auto">
          {/* Priority filter */}
          <DropdownMenu>
            <DropdownMenuTrigger asChild>
              <Button variant="outline" size="sm" className="h-8 gap-1 text-xs">
                <IconFilter className="h-3.5 w-3.5 text-muted-foreground" />
                <span>Priority</span>
                {priorityFilter && (
                  <Badge variant="secondary" className="ml-1 h-4 px-1 text-[10px] capitalize">
                    {priorityFilter}
                  </Badge>
                )}
              </Button>
            </DropdownMenuTrigger>
            <DropdownMenuContent align="end" className="w-36 text-xs">
              <DropdownMenuItem onClick={() => setPriorityFilter(null)}>
                All Priorities
              </DropdownMenuItem>
              <DropdownMenuSeparator />
              {(['urgent', 'high', 'medium', 'low'] as const).map((p) => (
                <DropdownMenuItem key={p} onClick={() => setPriorityFilter(p)} className="capitalize">
                  {p}
                </DropdownMenuItem>
              ))}
            </DropdownMenuContent>
          </DropdownMenu>

          {/* Sort Menu */}
          <DropdownMenu>
            <DropdownMenuTrigger asChild>
              <Button variant="outline" size="sm" className="h-8 gap-1 text-xs">
                {sortDirection === 'asc' ? (
                  <IconSortAscending className="h-3.5 w-3.5 text-muted-foreground" />
                ) : (
                  <IconSortDescending className="h-3.5 w-3.5 text-muted-foreground" />
                )}
                <span>Sort</span>
              </Button>
            </DropdownMenuTrigger>
            <DropdownMenuContent align="end" className="w-40 text-xs">
              <DropdownMenuItem onClick={() => setSortBy('none')}>Default</DropdownMenuItem>
              <DropdownMenuItem onClick={() => { setSortBy('dealValue'); setSortDirection((d) => (d === 'asc' ? 'desc' : 'asc')) }}>
                Deal Value
              </DropdownMenuItem>
              <DropdownMenuItem onClick={() => { setSortBy('dueDate'); setSortDirection((d) => (d === 'asc' ? 'desc' : 'asc')) }}>
                Due Date
              </DropdownMenuItem>
              <DropdownMenuItem onClick={() => { setSortBy('title'); setSortDirection((d) => (d === 'asc' ? 'desc' : 'asc')) }}>
                Title
              </DropdownMenuItem>
            </DropdownMenuContent>
          </DropdownMenu>

          {/* View Mode Toggle */}
          {onViewModeChange && (
            <div className="flex items-center rounded-md border bg-muted/30 p-0.5">
              <Button
                variant={viewMode === 'kanban' ? 'secondary' : 'ghost'}
                size="sm"
                onClick={() => onViewModeChange('kanban')}
                className="h-7 w-7 p-0"
                title="Kanban view"
              >
                <IconLayoutKanban className="h-3.5 w-3.5" />
              </Button>
              <Button
                variant={viewMode === 'table' ? 'secondary' : 'ghost'}
                size="sm"
                onClick={() => onViewModeChange('table')}
                className="h-7 w-7 p-0"
                title="Table view"
              >
                <IconTable className="h-3.5 w-3.5" />
              </Button>
            </div>
          )}

          {/* Add Record Button */}
          <Button
            size="sm"
            onClick={() => {
              if (onAddTask) {
                onAddTask(columns[0]?.id || '');
              } else {
                setNewTaskStatus(columns[0]?.id || '');
                setIsNewTaskOpen(true);
              }
            }}
            className="h-8 gap-1 text-xs"
          >
            <IconPlus className="h-3.5 w-3.5" />
            <span>Add Record</span>
          </Button>
        </div>
      </div>

      {/* Filter chips */}
      {(searchQuery || priorityFilter) && (
        <div className="flex items-center gap-1.5 px-1">
          {searchQuery && (
            <Badge variant="secondary" className="gap-1 text-[11px] h-5 px-1.5">
              <span>Search: {searchQuery}</span>
              <IconX className="h-3 w-3 cursor-pointer" onClick={() => setSearchQuery('')} />
            </Badge>
          )}
          {priorityFilter && (
            <Badge variant="secondary" className="gap-1 text-[11px] h-5 px-1.5 capitalize">
              <span>Priority: {priorityFilter}</span>
              <IconX className="h-3 w-3 cursor-pointer" onClick={() => setPriorityFilter(null)} />
            </Badge>
          )}
          <Button
            variant="ghost"
            size="sm"
            onClick={() => {
              setSearchQuery('');
              setPriorityFilter(null);
            }}
            className="h-5 px-1.5 text-[10px] text-muted-foreground hover:text-foreground"
          >
            Clear all
          </Button>
        </div>
      )}

      {/* 2. Horizontal Scrollable Kanban Columns */}
      <div className="overflow-x-auto pb-4 pt-1">
        <KanbanBoardProvider>
          <KanbanBoard>
            {columns.map((column) => {
              const colItems = columnTasksMap.get(column.id) || [];
              const aggregate = columnAggregates.get(column.id);

              return (
                <KanbanBoardColumn
                  key={column.id}
                  columnId={column.id}
                  onDropOverColumn={(data) => {
                    try {
                      const parsed = JSON.parse(data);
                      if (parsed && parsed.id) {
                        handleItemMove(parsed.id, column.id, colItems.length);
                      }
                    } catch {}
                  }}
                  className="w-72 shrink-0 rounded-md border border-border bg-card/60 p-2 shadow-xs"
                >
                  <KanbanBoardColumnHeader className="flex items-center justify-between pb-2 border-b border-border/60">
                  <div className="flex items-center gap-1.5 min-w-0">
                    <KanbanColorCircle color={column.color || 'blue'} />
                    <KanbanBoardColumnTitle columnId={column.id} className="text-xs font-semibold truncate text-foreground">
                      {column.title}
                    </KanbanBoardColumnTitle>
                    <Badge variant="secondary" className="h-4 px-1 text-[10px] font-normal">
                      {colItems.length}
                    </Badge>
                  </div>

                  <div className="flex items-center gap-0.5">
                    {aggregate && aggregate.totalValue > 0 && (
                      <span className="text-[10px] font-medium text-muted-foreground tabular-nums mr-1">
                        ${aggregate.totalValue.toLocaleString()}
                      </span>
                    )}
                    <Button
                      variant="ghost"
                      size="sm"
                      onClick={() => {
                        setNewTaskStatus(column.id);
                        setIsNewTaskOpen(true);
                      }}
                      className="h-6 w-6 p-0 hover:bg-muted"
                      title="Add card"
                    >
                      <IconPlus className="h-3 w-3 text-muted-foreground" />
                    </Button>
                  </div>
                </KanbanBoardColumnHeader>

                {/* Column Items */}
                <KanbanBoardColumnList className="flex flex-col gap-1.5 pt-2 min-h-[300px]">
                  {colItems.map((task, idx) => (
                    <KanbanBoardColumnListItem
                      key={task.id}
                      cardId={task.id}
                      onDropOverListItem={(data, direction: KanbanBoardDropDirection) => {
                        try {
                          const parsed = JSON.parse(data);
                          if (parsed && parsed.id) {
                            const newIdx = direction === 'bottom' ? idx + 1 : idx;
                            handleItemMove(parsed.id, column.id, newIdx);
                          }
                        } catch {}
                      }}
                    >
                      <KanbanBoardCard
                        data={task}
                        onClick={() => {
                          setActiveTask(task);
                          onTaskClick?.(task);
                        }}
                        className="group relative flex flex-col gap-1.5 rounded-md border border-border bg-card p-2.5 shadow-xs transition-all hover:border-primary/40 hover:shadow-sm cursor-grab select-none text-left w-full"
                      >
                        {/* Top: Customer or Category & Priority */}
                        <div className="flex items-center justify-between gap-1 w-full">
                          <span className="text-[10px] font-medium text-muted-foreground truncate max-w-[150px]">
                            {task.customerName || task.id}
                          </span>
                          <div className="flex items-center gap-1">
                            {task.priority && (
                              <span
                                className={cn(
                                  'h-1.5 w-1.5 rounded-full',
                                  task.priority === 'urgent' && 'bg-red-500',
                                  task.priority === 'high' && 'bg-amber-500',
                                  task.priority === 'medium' && 'bg-blue-500',
                                  task.priority === 'low' && 'bg-emerald-500'
                                )}
                                title={`Priority: ${task.priority}`}
                              />
                            )}
                            <DropdownMenu>
                              <DropdownMenuTrigger asChild onClick={(e) => e.stopPropagation()}>
                                <span
                                  role="button"
                                  tabIndex={0}
                                  onKeyDown={(e) => {
                                    if (e.key === 'Enter' || e.key === ' ') {
                                      e.preventDefault();
                                      e.stopPropagation();
                                      e.currentTarget.click();
                                    }
                                  }}
                                  className="opacity-0 group-hover:opacity-100 p-0.5 text-muted-foreground hover:text-foreground transition-opacity cursor-pointer inline-flex items-center justify-center rounded"
                                >
                                  <IconDotsVertical className="h-3 w-3" />
                                </span>
                              </DropdownMenuTrigger>
                              <DropdownMenuContent align="end" className="w-36 text-xs">
                                <DropdownMenuItem onClick={(e) => { e.stopPropagation(); setActiveTask(task); }}>
                                  <IconEdit className="mr-1.5 h-3.5 w-3.5" />
                                  Edit Record
                                </DropdownMenuItem>
                                <DropdownMenuItem onClick={(e) => { e.stopPropagation(); navigator.clipboard.writeText(task.id); toast.success('ID copied'); }}>
                                  <IconCopy className="mr-1.5 h-3.5 w-3.5" />
                                  Copy ID
                                </DropdownMenuItem>
                                <DropdownMenuSeparator />
                                <DropdownMenuItem
                                  onClick={(e) => { e.stopPropagation(); handleDeleteTask(task.id); }}
                                  className="text-destructive"
                                >
                                  <IconTrash className="mr-1.5 h-3.5 w-3.5" />
                                  Delete
                                </DropdownMenuItem>
                              </DropdownMenuContent>
                            </DropdownMenu>
                          </div>
                        </div>

                        {/* Title */}
                        <div className="text-xs font-medium text-foreground leading-snug line-clamp-2 w-full">
                          {task.title}
                        </div>

                        {/* Middle: Deal value / Metadata */}
                        {task.dealValue !== undefined && (
                          <div className="text-xs font-semibold text-primary tabular-nums">
                            ${task.dealValue.toLocaleString()}
                          </div>
                        )}

                        {/* Bottom: Assignee & Due Date / Indicators */}
                        <div className="flex items-center justify-between pt-1 border-t border-border/40 text-[10px] text-muted-foreground w-full">
                          <div className="flex items-center gap-1.5">
                            {task.assignee ? (
                              <Avatar className="h-5 w-5">
                                <AvatarImage src={task.assignee.avatar} />
                                <AvatarFallback className="text-[9px]">
                                  {task.assignee.initials || task.assignee.name.slice(0, 2).toUpperCase()}
                                </AvatarFallback>
                              </Avatar>
                            ) : (
                              <IconUser className="h-3.5 w-3.5 opacity-50" />
                            )}
                            {task.dueDate && (
                              <span className="flex items-center gap-0.5">
                                <IconCalendar className="h-3 w-3 opacity-60" />
                                {task.dueDate}
                              </span>
                            )}
                          </div>

                          <div className="flex items-center gap-1.5">
                            {task.commentsCount !== undefined && task.commentsCount > 0 && (
                              <span className="flex items-center gap-0.5">
                                <IconMessageCircle className="h-3 w-3 opacity-70" />
                                {task.commentsCount}
                              </span>
                            )}
                            {task.attachmentsCount !== undefined && task.attachmentsCount > 0 && (
                              <span className="flex items-center gap-0.5">
                                <IconPaperclip className="h-3 w-3 opacity-70" />
                                {task.attachmentsCount}
                              </span>
                            )}
                          </div>
                        </div>
                      </KanbanBoardCard>
                    </KanbanBoardColumnListItem>
                  ))}

                  {colItems.length === 0 && (
                    <div className="flex flex-col items-center justify-center p-6 text-center border border-dashed rounded-md border-border/60 text-muted-foreground">
                      <p className="text-[11px]">No records</p>
                      <button
                        type="button"
                        onClick={() => {
                          setNewTaskStatus(column.id);
                          setIsNewTaskOpen(true);
                        }}
                        className="text-[10px] text-primary hover:underline mt-1"
                      >
                        + Add record
                      </button>
                    </div>
                  )}
                </KanbanBoardColumnList>
              </KanbanBoardColumn>
            );
          })}
          <KanbanBoardExtraMargin />
        </KanbanBoard>
      </KanbanBoardProvider>
    </div>

      {/* 3. Detail Drawer (Sheet) */}
      <Sheet open={!!activeTask} onOpenChange={(open) => !open && setActiveTask(null)}>
        <SheetContent className="w-full sm:max-w-md p-4 flex flex-col justify-between">
          {activeTask && (
            <>
              <SheetHeader className="text-left space-y-1">
                <div className="flex items-center justify-between">
                  <Badge variant="outline" className="text-xs">
                    {columns.find((c) => c.id === activeTask.status)?.title || activeTask.status}
                  </Badge>
                  <span className="text-[10px] text-muted-foreground font-mono">
                    {activeTask.id}
                  </span>
                </div>
                <SheetTitle className="text-base font-bold text-foreground">
                  {activeTask.title}
                </SheetTitle>
                <SheetDescription className="text-xs text-muted-foreground">
                  {activeTask.customerName ? `Client: ${activeTask.customerName}` : 'Record details'}
                </SheetDescription>
              </SheetHeader>

              <div className="space-y-4 py-4 text-xs">
                {activeTask.dealValue !== undefined && (
                  <div className="rounded-md border p-2.5 bg-muted/20">
                    <span className="text-muted-foreground block text-[10px]">Deal Value</span>
                    <span className="text-lg font-bold text-primary tabular-nums">
                      ${activeTask.dealValue.toLocaleString()}
                    </span>
                  </div>
                )}

                <div className="grid grid-cols-2 gap-3">
                  <div>
                    <label className="text-[10px] text-muted-foreground block mb-1">Priority</label>
                    <Badge variant="secondary" className="capitalize">
                      {activeTask.priority || 'medium'}
                    </Badge>
                  </div>
                  <div>
                    <label className="text-[10px] text-muted-foreground block mb-1">Due Date</label>
                    <span className="font-medium text-foreground">
                      {activeTask.dueDate || 'No date set'}
                    </span>
                  </div>
                </div>

                {activeTask.assignee && (
                  <div>
                    <label className="text-[10px] text-muted-foreground block mb-1">Assigned To</label>
                    <div className="flex items-center gap-2">
                      <Avatar className="h-6 w-6">
                        <AvatarImage src={activeTask.assignee.avatar} />
                        <AvatarFallback className="text-[10px]">
                          {activeTask.assignee.initials || 'U'}
                        </AvatarFallback>
                      </Avatar>
                      <span className="font-medium text-foreground">{activeTask.assignee.name}</span>
                    </div>
                  </div>
                )}

                {activeTask.tags && activeTask.tags.length > 0 && (
                  <div>
                    <label className="text-[10px] text-muted-foreground block mb-1">Tags</label>
                    <div className="flex flex-wrap gap-1">
                      {activeTask.tags.map((tag) => (
                        <Badge key={tag} variant="outline" className="text-[10px]">
                          {tag}
                        </Badge>
                      ))}
                    </div>
                  </div>
                )}

                {activeTask.description && (
                  <div>
                    <label className="text-[10px] text-muted-foreground block mb-1">Notes / Description</label>
                    <p className="p-2 rounded-md bg-muted/30 text-foreground">
                      {activeTask.description}
                    </p>
                  </div>
                )}
              </div>

              <SheetFooter className="flex-row items-center justify-between sm:justify-between border-t pt-3">
                <Button
                  variant="destructive"
                  size="sm"
                  onClick={() => handleDeleteTask(activeTask.id)}
                  className="h-8 gap-1 text-xs"
                >
                  <IconTrash className="h-3.5 w-3.5" />
                  <span>Delete</span>
                </Button>
                <Button
                  variant="outline"
                  size="sm"
                  onClick={() => setActiveTask(null)}
                  className="h-8 text-xs"
                >
                  Close
                </Button>
              </SheetFooter>
            </>
          )}
        </SheetContent>
      </Sheet>

      {/* 4. Add Record Modal (Dialog) */}
      <Dialog open={isNewTaskOpen} onOpenChange={setIsNewTaskOpen}>
        <DialogContent className="sm:max-w-md p-4">
          <DialogHeader className="text-left space-y-1">
            <DialogTitle className="text-sm font-semibold">Create New Record</DialogTitle>
            <DialogDescription className="text-xs text-muted-foreground">
              Add a new task or deal into the {columns.find((c) => c.id === newTaskStatus)?.title || 'pipeline'}.
            </DialogDescription>
          </DialogHeader>

          <div className="space-y-3 py-2 text-xs">
            <div>
              <label className="text-xs font-medium text-muted-foreground block mb-1">
                Title / Deal Name *
              </label>
              <Input
                value={newTaskTitle}
                onChange={(e) => setNewTaskTitle(e.target.value)}
                placeholder="e.g. Wedding Photography Package"
                className="h-8 text-xs"
                autoFocus
              />
            </div>

            <div className="grid grid-cols-2 gap-2">
              <div>
                <label className="text-xs font-medium text-muted-foreground block mb-1">
                  Customer / Client
                </label>
                <Input
                  value={newTaskCustomer}
                  onChange={(e) => setNewTaskCustomer(e.target.value)}
                  placeholder="e.g. Sarah Jenkins"
                  className="h-8 text-xs"
                />
              </div>
              <div>
                <label className="text-xs font-medium text-muted-foreground block mb-1">
                  Deal Value ($)
                </label>
                <Input
                  type="number"
                  value={newTaskValue}
                  onChange={(e) => setNewTaskValue(e.target.value)}
                  placeholder="1500"
                  className="h-8 text-xs"
                />
              </div>
            </div>

            <div>
              <label className="text-xs font-medium text-muted-foreground block mb-1">
                Priority
              </label>
              <div className="flex gap-1.5">
                {(['low', 'medium', 'high', 'urgent'] as const).map((p) => (
                  <Button
                    key={p}
                    type="button"
                    variant={newTaskPriority === p ? 'default' : 'outline'}
                    size="sm"
                    onClick={() => setNewTaskPriority(p)}
                    className="h-7 text-xs flex-1 capitalize"
                  >
                    {p}
                  </Button>
                ))}
              </div>
            </div>
          </div>

          <DialogFooter className="flex-row items-center justify-end gap-2 border-t pt-3">
            <Button
              variant="outline"
              size="sm"
              onClick={() => setIsNewTaskOpen(false)}
              className="h-8 text-xs"
            >
              Cancel
            </Button>
            <Button
              size="sm"
              onClick={handleCreateTask}
              className="h-8 text-xs"
            >
              Create Record
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>
    </div>
  );
}
