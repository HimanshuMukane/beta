import * as React from "react"
import {
  IconPlus,
  IconTrash,
  IconChevronDown,
  IconCheck,
  IconTag,
  IconCornerDownRight,
} from "../../icons/tabler-icons"
import { Button } from "../ui/button"
import { Input } from "../ui/input"
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from "../ui/dropdown-menu"
import {
  Popover,
  PopoverContent,
  PopoverTrigger,
} from "../ui/popover"
import { cn } from "@repo/ui/lib/utils"

export type FilterFieldType = "select" | "multi-select" | "text" | "number" | "date"

export interface FilterOption {
  label: string
  value: string
  icon?: React.ComponentType<{ className?: string }>
}

export interface FilterColumnConfig {
  columnId: string
  title: string
  type: FilterFieldType
  options?: FilterOption[]
}

export type FilterOperator =
  | "is"
  | "is_not"
  | "is_set"
  | "is_not_set"
  | "contains"
  | "does_not_contain"
  | "starts_with"
  | "ends_with"
  | "gt"
  | "gte"
  | "lt"
  | "lte"
  | "before"
  | "after"

export interface FilterRuleItem {
  id: string
  type: "rule"
  field: string
  operator: FilterOperator
  value?: string | string[]
}

export interface FilterGroupItem {
  id: string
  type: "group"
  conjunction: "and" | "or"
  items: (FilterRuleItem | FilterGroupItem)[]
}

export interface FilterBuilderProps {
  configs: FilterColumnConfig[]
  rootGroup: FilterGroupItem
  onChange: (group: FilterGroupItem) => void
  onClear: () => void
}

const OPERATOR_LABELS: Record<FilterOperator, string> = {
  is: "Is",
  is_not: "Is not",
  is_set: "Is set",
  is_not_set: "Is not set",
  contains: "Contains",
  does_not_contain: "Does not contain",
  starts_with: "Starts with",
  ends_with: "Ends with",
  gt: "Greater than",
  gte: "Greater than or equal",
  lt: "Less than",
  lte: "Less than or equal",
  before: "Is before",
  after: "Is after",
}

export function getOperatorsForType(type: FilterFieldType): FilterOperator[] {
  switch (type) {
    case "select":
      return ["is", "is_not", "is_set", "is_not_set"]
    case "multi-select":
      return ["contains", "does_not_contain", "is_set", "is_not_set"]
    case "date":
      return ["is", "before", "after", "is_set", "is_not_set"]
    case "number":
      return ["is", "is_not", "gt", "gte", "lt", "lte", "is_set", "is_not_set"]
    case "text":
    default:
      return ["contains", "does_not_contain", "is", "starts_with", "ends_with", "is_set", "is_not_set"]
  }
}

export function evaluateFilterTree(rowValueGetter: (colId: string) => any, group: FilterGroupItem): boolean {
  if (!group.items || group.items.length === 0) return true

  const results = group.items.map((item) => {
    if (item.type === "group") {
      return evaluateFilterTree(rowValueGetter, item)
    }

    const cellVal = rowValueGetter(item.field)
    const { operator, value } = item

    if (operator === "is_set") {
      return cellVal !== null && cellVal !== undefined && cellVal !== ""
    }
    if (operator === "is_not_set") {
      return cellVal === null || cellVal === undefined || cellVal === ""
    }

    if (cellVal === null || cellVal === undefined) {
      return false
    }

    const strCell = String(cellVal).toLowerCase()
    const strVal = typeof value === "string" ? value.toLowerCase() : ""

    switch (operator) {
      case "is":
        if (Array.isArray(value)) {
          return value.some((v) => String(v).toLowerCase() === strCell)
        }
        return strCell === strVal
      case "is_not":
        if (Array.isArray(value)) {
          return !value.some((v) => String(v).toLowerCase() === strCell)
        }
        return strCell !== strVal
      case "contains":
        if (Array.isArray(cellVal)) {
          if (Array.isArray(value)) {
            return value.some((v) => cellVal.map(String).includes(String(v)))
          }
          return cellVal.map(String).includes(String(value))
        }
        return strCell.includes(strVal)
      case "does_not_contain":
        if (Array.isArray(cellVal)) {
          if (Array.isArray(value)) {
            return !value.some((v) => cellVal.map(String).includes(String(v)))
          }
          return !cellVal.map(String).includes(String(value))
        }
        return !strCell.includes(strVal)
      case "starts_with":
        return strCell.startsWith(strVal)
      case "ends_with":
        return strCell.endsWith(strVal)
      case "gt":
        return Number(cellVal) > Number(value)
      case "gte":
        return Number(cellVal) >= Number(value)
      case "lt":
        return Number(cellVal) < Number(value)
      case "lte":
        return Number(cellVal) <= Number(value)
      case "before":
        return new Date(cellVal).getTime() < new Date(String(value)).getTime()
      case "after":
        return new Date(cellVal).getTime() > new Date(String(value)).getTime()
      default:
        return true
    }
  })

  if (group.conjunction === "or") {
    return results.some(Boolean)
  }
  return results.every(Boolean)
}

export function FilterBuilder({ configs, rootGroup, onChange, onClear }: FilterBuilderProps) {
  const handleAddRule = (targetGroup: FilterGroupItem) => {
    const defaultField = configs[0]?.columnId || ""
    const defaultType = configs[0]?.type || "text"
    const defaultOp = getOperatorsForType(defaultType)[0] || "is"

    const newRule: FilterRuleItem = {
      id: "rule_" + Math.random().toString(36).slice(2, 9),
      type: "rule",
      field: defaultField,
      operator: defaultOp,
      value: "",
    }

    const updated = addRuleToGroup(rootGroup, targetGroup.id, newRule)
    onChange(updated)
  }

  const handleAddNestedGroup = (parentGroupId: string) => {
    const defaultField = configs[0]?.columnId || ""
    const defaultType = configs[0]?.type || "text"
    const defaultOp = getOperatorsForType(defaultType)[0] || "is"

    const newGroup: FilterGroupItem = {
      id: "grp_" + Math.random().toString(36).slice(2, 9),
      type: "group",
      conjunction: "and",
      items: [
        {
          id: "rule_" + Math.random().toString(36).slice(2, 9),
          type: "rule",
          field: defaultField,
          operator: defaultOp,
          value: "",
        },
      ],
    }

    const updated = addGroupToGroup(rootGroup, parentGroupId, newGroup)
    onChange(updated)
  }

  const handleRemoveItem = (itemId: string) => {
    const updated = removeItemFromGroup(rootGroup, itemId)
    onChange(updated)
  }

  const handleUpdateRule = (ruleId: string, patch: Partial<FilterRuleItem>) => {
    const updated = updateRuleInGroup(rootGroup, ruleId, patch, configs)
    onChange(updated)
  }

  const handleUpdateGroupConjunction = (groupId: string, conjunction: "and" | "or") => {
    const updated = updateGroupInGroup(rootGroup, groupId, { conjunction })
    onChange(updated)
  }

  const totalRules = countRules(rootGroup)

  return (
    <div className="flex flex-col gap-3 p-3.5 bg-neutral-900 border border-neutral-800 text-neutral-100 rounded-xl shadow-2xl min-w-[560px] max-w-[700px]">
      <FilterGroupView
        group={rootGroup}
        isRoot={true}
        configs={configs}
        onAddRule={handleAddRule}
        onAddNestedGroup={handleAddNestedGroup}
        onRemoveItem={handleRemoveItem}
        onUpdateRule={handleUpdateRule}
        onUpdateGroupConjunction={handleUpdateGroupConjunction}
      />

      <div className="flex items-center justify-between pt-2 border-t border-neutral-800/80">
        <Button
          type="button"
          variant="outline"
          size="sm"
          onClick={() => handleAddRule(rootGroup)}
          className="h-8 gap-1.5 bg-neutral-800/80 border-neutral-700/80 hover:bg-neutral-800 text-neutral-200 text-xs font-medium rounded-lg"
        >
          <IconPlus className="h-3.5 w-3.5" />
          <span>Add filter</span>
        </Button>

        {totalRules > 0 && (
          <Button
            type="button"
            variant="ghost"
            size="sm"
            onClick={onClear}
            className="h-8 px-2.5 text-xs text-red-400 hover:text-red-300 hover:bg-red-950/40 rounded-lg"
          >
            Clear all
          </Button>
        )}
      </div>
    </div>
  )
}

interface FilterGroupViewProps {
  group: FilterGroupItem
  isRoot: boolean
  configs: FilterColumnConfig[]
  onAddRule: (group: FilterGroupItem) => void
  onAddNestedGroup: (parentGroupId: string) => void
  onRemoveItem: (itemId: string) => void
  onUpdateRule: (ruleId: string, patch: Partial<FilterRuleItem>) => void
  onUpdateGroupConjunction: (groupId: string, conjunction: "and" | "or") => void
}

function FilterGroupView({
  group,
  isRoot,
  configs,
  onAddRule,
  onAddNestedGroup,
  onRemoveItem,
  onUpdateRule,
  onUpdateGroupConjunction,
}: FilterGroupViewProps) {
  return (
    <div
      className={cn(
        "flex flex-col gap-2 rounded-lg transition-colors",
        !isRoot && "p-2.5 bg-neutral-950/50 border border-neutral-800/90 relative ml-4"
      )}
    >
      {group.items.map((item, index) => {
        const isFirst = index === 0

        if (item.type === "group") {
          return (
            <div key={item.id} className="relative flex flex-col gap-1">
              <div className="flex items-center gap-2 mb-1">
                <ConjunctionSelector
                  conjunction={group.conjunction}
                  onChange={(c) => onUpdateGroupConjunction(group.id, c)}
                  isFirst={isFirst}
                />
                <span className="text-[11px] text-neutral-400 font-mono tracking-wide">
                  Nested condition
                </span>
                <button
                  type="button"
                  onClick={() => onRemoveItem(item.id)}
                  className="ml-auto p-1 text-neutral-400 hover:text-red-400 transition-colors"
                  title="Delete group"
                >
                  <IconTrash className="h-3.5 w-3.5" />
                </button>
              </div>

              <FilterGroupView
                group={item}
                isRoot={false}
                configs={configs}
                onAddRule={onAddRule}
                onAddNestedGroup={onAddNestedGroup}
                onRemoveItem={onRemoveItem}
                onUpdateRule={onUpdateRule}
                onUpdateGroupConjunction={onUpdateGroupConjunction}
              />
            </div>
          )
        }

        return (
          <FilterRowView
            key={item.id}
            rule={item}
            isFirst={isFirst}
            conjunction={group.conjunction}
            configs={configs}
            onConjunctionChange={(c) => onUpdateGroupConjunction(group.id, c)}
            onUpdateRule={(patch) => onUpdateRule(item.id, patch)}
            onRemove={() => onRemoveItem(item.id)}
          />
        )
      })}

      <div className="pt-0.5">
        <button
          type="button"
          onClick={() => onAddNestedGroup(group.id)}
          className="text-[12px] text-neutral-400 hover:text-neutral-200 transition-colors inline-flex items-center gap-1.5 py-1 px-1"
        >
          <IconCornerDownRight className="h-3 w-3 text-neutral-500" />
          <span>Add nested filter</span>
        </button>
      </div>
    </div>
  )
}

interface FilterRowViewProps {
  rule: FilterRuleItem
  isFirst: boolean
  conjunction: "and" | "or"
  configs: FilterColumnConfig[]
  onConjunctionChange: (conjunction: "and" | "or") => void
  onUpdateRule: (patch: Partial<FilterRuleItem>) => void
  onRemove: () => void
}

function FilterRowView({
  rule,
  isFirst,
  conjunction,
  configs,
  onConjunctionChange,
  onUpdateRule,
  onRemove,
}: FilterRowViewProps) {
  const currentConfig = configs.find((c) => c.columnId === rule.field) || configs[0]
  const fieldType = currentConfig?.type || "text"
  const availableOperators = getOperatorsForType(fieldType)

  const isUnaryOperator = rule.operator === "is_set" || rule.operator === "is_not_set"

  return (
    <div className="flex items-center gap-2 text-xs">
      {/* 1. Conjunction / Where Label */}
      <div className="w-[72px] shrink-0 flex items-center justify-end">
        {isFirst ? (
          <span className="text-neutral-400 font-medium text-xs pr-1">Where</span>
        ) : (
          <DropdownMenu>
            <DropdownMenuTrigger asChild>
              <button
                type="button"
                className="flex items-center justify-between w-full px-2 py-1.5 rounded-lg bg-neutral-800/80 hover:bg-neutral-800 border border-neutral-700/80 text-neutral-200 font-medium text-xs transition-colors"
              >
                <span>{conjunction.toUpperCase()}</span>
                <IconChevronDown className="h-3 w-3 text-neutral-400" />
              </button>
            </DropdownMenuTrigger>
            <DropdownMenuContent align="start" className="w-24 bg-neutral-900 border-neutral-800 text-neutral-200">
              <DropdownMenuItem
                onClick={() => onConjunctionChange("and")}
                className="cursor-pointer text-xs focus:bg-neutral-800 focus:text-white"
              >
                AND
              </DropdownMenuItem>
              <DropdownMenuItem
                onClick={() => onConjunctionChange("or")}
                className="cursor-pointer text-xs focus:bg-neutral-800 focus:text-white"
              >
                OR
              </DropdownMenuItem>
            </DropdownMenuContent>
          </DropdownMenu>
        )}
      </div>

      {/* 2. Column / Field Dropdown */}
      <DropdownMenu>
        <DropdownMenuTrigger asChild>
          <button
            type="button"
            className="flex items-center justify-between min-w-[125px] px-3 py-1.5 rounded-lg bg-neutral-950/80 hover:bg-neutral-800 border border-neutral-800 text-neutral-100 font-medium text-xs transition-colors"
          >
            <span className="truncate">{currentConfig?.title || rule.field}</span>
            <IconChevronDown className="h-3 w-3 ml-2 text-neutral-400 shrink-0" />
          </button>
        </DropdownMenuTrigger>
        <DropdownMenuContent align="start" className="w-44 bg-neutral-900 border-neutral-800 text-neutral-200">
          {configs.map((config) => (
            <DropdownMenuItem
              key={config.columnId}
              onClick={() => {
                const newType = config.type
                const newOps = getOperatorsForType(newType)
                onUpdateRule({
                  field: config.columnId,
                  operator: newOps[0] || "is",
                  value: "",
                })
              }}
              className="flex items-center justify-between cursor-pointer text-xs focus:bg-neutral-800 focus:text-white"
            >
              <span>{config.title}</span>
              {rule.field === config.columnId && <IconCheck className="h-3.5 w-3.5 text-neutral-300" />}
            </DropdownMenuItem>
          ))}
        </DropdownMenuContent>
      </DropdownMenu>

      {/* 3. Operator Dropdown */}
      <DropdownMenu>
        <DropdownMenuTrigger asChild>
          <button
            type="button"
            className="flex items-center justify-between min-w-[95px] px-3 py-1.5 rounded-lg bg-neutral-950/80 hover:bg-neutral-800 border border-neutral-800 text-neutral-100 text-xs transition-colors"
          >
            <span>{OPERATOR_LABELS[rule.operator] || rule.operator}</span>
            <IconChevronDown className="h-3 w-3 ml-2 text-neutral-400 shrink-0" />
          </button>
        </DropdownMenuTrigger>
        <DropdownMenuContent align="start" className="w-36 bg-neutral-900 border-neutral-800 text-neutral-200">
          {availableOperators.map((op) => (
            <DropdownMenuItem
              key={op}
              onClick={() => onUpdateRule({ operator: op })}
              className="flex items-center justify-between cursor-pointer text-xs focus:bg-neutral-800 focus:text-white"
            >
              <span>{OPERATOR_LABELS[op]}</span>
              {rule.operator === op && <IconCheck className="h-3.5 w-3.5 text-neutral-300" />}
            </DropdownMenuItem>
          ))}
        </DropdownMenuContent>
      </DropdownMenu>

      {/* 4. Value Input / Select */}
      <div className="flex-1 min-w-[170px]">
        {isUnaryOperator ? (
          <div className="h-8 px-2.5 flex items-center text-[11px] text-neutral-500 italic rounded-lg bg-neutral-950/40 border border-transparent">
            No value needed
          </div>
        ) : fieldType === "select" && currentConfig?.options ? (
          <DropdownMenu>
            <DropdownMenuTrigger asChild>
              <button
                type="button"
                className="flex items-center justify-between w-full h-8 px-3 py-1.5 rounded-lg bg-neutral-950/80 hover:bg-neutral-800 border border-neutral-800 text-neutral-200 text-xs transition-colors"
              >
                <span className={cn(!rule.value && "text-neutral-500")}>
                  {currentConfig.options.find((o) => o.value === rule.value)?.label || "Select option"}
                </span>
                <IconChevronDown className="h-3 w-3 text-neutral-400 shrink-0" />
              </button>
            </DropdownMenuTrigger>
            <DropdownMenuContent align="start" className="w-48 bg-neutral-900 border-neutral-800 text-neutral-200">
              {currentConfig.options.map((opt) => (
                <DropdownMenuItem
                  key={opt.value}
                  onClick={() => onUpdateRule({ value: opt.value })}
                  className="flex items-center justify-between cursor-pointer text-xs focus:bg-neutral-800 focus:text-white"
                >
                  <div className="flex items-center gap-1.5">
                    {opt.icon && <opt.icon className="h-3.5 w-3.5" />}
                    <span>{opt.label}</span>
                  </div>
                  {rule.value === opt.value && <IconCheck className="h-3.5 w-3.5 text-neutral-300" />}
                </DropdownMenuItem>
              ))}
            </DropdownMenuContent>
          </DropdownMenu>
        ) : fieldType === "multi-select" && currentConfig?.options ? (
          <MultiSelectTagsPicker
            options={currentConfig.options}
            selectedValues={Array.isArray(rule.value) ? rule.value : rule.value ? [rule.value] : []}
            onChange={(vals) => onUpdateRule({ value: vals })}
          />
        ) : fieldType === "date" ? (
          <Input
            type="date"
            value={typeof rule.value === "string" ? rule.value : ""}
            onChange={(e) => onUpdateRule({ value: e.target.value })}
            className="h-8 text-xs bg-neutral-950/80 border-neutral-800 text-neutral-200 [color-scheme:dark]"
          />
        ) : fieldType === "number" ? (
          <Input
            type="number"
            placeholder="Enter number..."
            value={typeof rule.value === "string" ? rule.value : ""}
            onChange={(e) => onUpdateRule({ value: e.target.value })}
            className="h-8 text-xs bg-neutral-950/80 border-neutral-800 text-neutral-200 placeholder:text-neutral-500"
          />
        ) : (
          <Input
            type="text"
            placeholder="Enter value..."
            value={typeof rule.value === "string" ? rule.value : ""}
            onChange={(e) => onUpdateRule({ value: e.target.value })}
            className="h-8 text-xs bg-neutral-950/80 border-neutral-800 text-neutral-200 placeholder:text-neutral-500"
          />
        )}
      </div>

      {/* 5. Delete Action */}
      <button
        type="button"
        onClick={onRemove}
        className="p-1.5 rounded-md hover:bg-neutral-800 text-neutral-400 hover:text-red-400 transition-colors shrink-0"
        title="Remove filter"
      >
        <IconTrash className="h-4 w-4" />
      </button>
    </div>
  )
}

function ConjunctionSelector({
  conjunction,
  onChange,
  isFirst,
}: {
  conjunction: "and" | "or"
  onChange: (val: "and" | "or") => void
  isFirst: boolean
}) {
  if (isFirst) {
    return <span className="text-neutral-400 font-medium text-xs pr-1">Where</span>
  }

  return (
    <DropdownMenu>
      <DropdownMenuTrigger asChild>
        <button
          type="button"
          className="flex items-center gap-1 px-2 py-1 rounded bg-neutral-800 hover:bg-neutral-700 border border-neutral-700 text-neutral-200 text-xs font-medium"
        >
          <span>{conjunction.toUpperCase()}</span>
          <IconChevronDown className="h-3 w-3 text-neutral-400" />
        </button>
      </DropdownMenuTrigger>
      <DropdownMenuContent align="start" className="w-20 bg-neutral-900 border-neutral-800 text-neutral-200">
        <DropdownMenuItem
          onClick={() => onChange("and")}
          className="cursor-pointer text-xs focus:bg-neutral-800 focus:text-white"
        >
          AND
        </DropdownMenuItem>
        <DropdownMenuItem
          onClick={() => onChange("or")}
          className="cursor-pointer text-xs focus:bg-neutral-800 focus:text-white"
        >
          OR
        </DropdownMenuItem>
      </DropdownMenuContent>
    </DropdownMenu>
  )
}

function MultiSelectTagsPicker({
  options,
  selectedValues,
  onChange,
}: {
  options: FilterOption[]
  selectedValues: string[]
  onChange: (vals: string[]) => void
}) {
  const toggleValue = (val: string) => {
    if (selectedValues.includes(val)) {
      onChange(selectedValues.filter((v) => v !== val))
    } else {
      onChange([...selectedValues, val])
    }
  }

  return (
    <Popover>
      <PopoverTrigger asChild>
        <button
          type="button"
          className="flex items-center justify-between w-full h-8 px-3 py-1.5 rounded-lg bg-neutral-950/80 hover:bg-neutral-800 border border-neutral-800 text-neutral-200 text-xs transition-colors"
        >
          <div className="flex items-center gap-1.5 overflow-hidden">
            <IconTag className="h-3.5 w-3.5 text-neutral-400 shrink-0" />
            <span className={cn("truncate", selectedValues.length === 0 && "text-neutral-500")}>
              {selectedValues.length === 0
                ? "Select tags"
                : `${selectedValues.length} tag${selectedValues.length > 1 ? "s" : ""} selected`}
            </span>
          </div>
          <IconChevronDown className="h-3 w-3 text-neutral-400 shrink-0" />
        </button>
      </PopoverTrigger>
      <PopoverContent align="start" className="w-52 p-2 bg-neutral-900 border-neutral-800 text-neutral-200">
        <div className="space-y-1">
          {options.map((opt) => {
            const isChecked = selectedValues.includes(opt.value)
            return (
              <button
                key={opt.value}
                type="button"
                onClick={() => toggleValue(opt.value)}
                className="w-full flex items-center justify-between px-2 py-1.5 rounded text-xs hover:bg-neutral-800 text-left transition-colors"
              >
                <span>{opt.label}</span>
                {isChecked && <IconCheck className="h-3.5 w-3.5 text-neutral-300" />}
              </button>
            )
          })}
        </div>
      </PopoverContent>
    </Popover>
  )
}

function addRuleToGroup(
  root: FilterGroupItem,
  targetGroupId: string,
  newRule: FilterRuleItem
): FilterGroupItem {
  if (root.id === targetGroupId) {
    return { ...root, items: [...root.items, newRule] }
  }
  return {
    ...root,
    items: root.items.map((item) => {
      if (item.type === "group") {
        return addRuleToGroup(item, targetGroupId, newRule)
      }
      return item
    }),
  }
}

function addGroupToGroup(
  root: FilterGroupItem,
  parentGroupId: string,
  newGroup: FilterGroupItem
): FilterGroupItem {
  if (root.id === parentGroupId) {
    return { ...root, items: [...root.items, newGroup] }
  }
  return {
    ...root,
    items: root.items.map((item) => {
      if (item.type === "group") {
        return addGroupToGroup(item, parentGroupId, newGroup)
      }
      return item
    }),
  }
}

function removeItemFromGroup(root: FilterGroupItem, itemId: string): FilterGroupItem {
  return {
    ...root,
    items: root.items
      .filter((item) => item.id !== itemId)
      .map((item) => {
        if (item.type === "group") {
          return removeItemFromGroup(item, itemId)
        }
        return item
      }),
  }
}

function updateRuleInGroup(
  root: FilterGroupItem,
  ruleId: string,
  patch: Partial<FilterRuleItem>,
  configs: FilterColumnConfig[]
): FilterGroupItem {
  return {
    ...root,
    items: root.items.map((item) => {
      if (item.type === "rule" && item.id === ruleId) {
        return { ...item, ...patch }
      }
      if (item.type === "group") {
        return updateRuleInGroup(item, ruleId, patch, configs)
      }
      return item
    }),
  }
}

function updateGroupInGroup(
  root: FilterGroupItem,
  groupId: string,
  patch: Partial<FilterGroupItem>
): FilterGroupItem {
  if (root.id === groupId) {
    return { ...root, ...patch }
  }
  return {
    ...root,
    items: root.items.map((item) => {
      if (item.type === "group") {
        return updateGroupInGroup(item, groupId, patch)
      }
      return item
    }),
  }
}

function countRules(group: FilterGroupItem): number {
  return group.items.reduce((acc, item) => {
    if (item.type === "rule") return acc + 1
    if (item.type === "group") return acc + countRules(item)
    return acc
  }, 0)
}
