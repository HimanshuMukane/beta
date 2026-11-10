import * as React from 'react'
import {
  IconPlus,
  IconTrash,
  IconEdit,
  IconCheck,
  IconX,
  IconDatabase,
} from '@tabler/icons-react'
import { Card, CardContent } from '@repo/ui/components/ui/card'
import { Button } from '@repo/ui/components/ui/button'
import { Input } from '@repo/ui/components/ui/input'
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogDescription,
  DialogFooter,
} from '@repo/ui/components/ui/dialog'
import { toast } from '@repo/ui/components/ui/sonner'
import type { MasterDataCategory, MasterDataItem } from '../../types/settings'
import { settingsApi } from '../../api/settingsApi'

export function MasterDataPage() {
  const [categories, setCategories] = React.useState<Record<string, MasterDataCategory>>({})
  const [activeCategoryKey, setActiveCategoryKey] = React.useState<string | null>(null)

  React.useEffect(() => {
    settingsApi.getSettings().then((s) => {
      if (s.masterData && Object.keys(s.masterData).length > 0) {
        setCategories(s.masterData)
      }
    }).catch(console.error)
  }, [])

  // Drawer / Modal item editing state
  const [editingItemId, setEditingItemId] = React.useState<string | null>(null)
  const [editItemName, setEditItemName] = React.useState('')
  const [newItemName, setNewItemName] = React.useState('')

  const activeCategory = activeCategoryKey ? categories[activeCategoryKey] : null

  const handleOpenCategory = (key: string) => {
    setActiveCategoryKey(key)
    setEditingItemId(null)
    setEditItemName('')
    setNewItemName('')
  }

  const handleCloseCategory = () => {
    setActiveCategoryKey(null)
    setEditingItemId(null)
    setEditItemName('')
    setNewItemName('')
  }

  const handleAddItem = () => {
    if (!newItemName.trim() || !activeCategoryKey) return
    const newItem: MasterDataItem = {
      id: `item-${Date.now()}`,
      name: newItemName.trim(),
    }

    const updated = {
      ...categories,
      [activeCategoryKey]: {
        ...categories[activeCategoryKey],
        items: [...(categories[activeCategoryKey]?.items || []), newItem],
      },
    }
    setCategories(updated)
    settingsApi.updateSettings({ masterData: updated }).catch(console.error)
    setNewItemName('')
    toast.success(`Added "${newItem.name}" to ${activeCategory?.title}`)
  }

  const handleStartEdit = (item: MasterDataItem) => {
    setEditingItemId(item.id)
    setEditItemName(item.name)
  }

  const handleSaveEdit = (itemId: string) => {
    if (!editItemName.trim() || !activeCategoryKey) return
    const updated = {
      ...categories,
      [activeCategoryKey]: {
        ...categories[activeCategoryKey],
        items: (categories[activeCategoryKey]?.items || []).map((it) =>
          it.id === itemId ? { ...it, name: editItemName.trim() } : it
        ),
      },
    }
    setCategories(updated)
    settingsApi.updateSettings({ masterData: updated }).catch(console.error)
    setEditingItemId(null)
    setEditItemName('')
    toast.success('Updated successfully')
  }

  const handleDeleteItem = (itemId: string) => {
    if (!activeCategoryKey) return
    const updated = {
      ...categories,
      [activeCategoryKey]: {
        ...categories[activeCategoryKey],
        items: (categories[activeCategoryKey]?.items || []).filter((it) => it.id !== itemId),
      },
    }
    setCategories(updated)
    settingsApi.updateSettings({ masterData: updated }).catch(console.error)
    toast.success('Item deleted')
  }

  return (
    <div className="space-y-6 max-w-4xl mx-auto pb-12">
      <div className="flex flex-col sm:flex-row sm:items-center justify-between gap-4 border-b pb-4">
        <div className="flex items-center gap-3">
          <span className="flex h-10 w-10 shrink-0 items-center justify-center rounded-xl bg-primary/10 text-primary border border-primary/20">
            <IconDatabase className="h-5 w-5" />
          </span>
          <div>
            <h1 className="text-xl font-bold tracking-tight text-foreground">Master Data</h1>
            <p className="text-xs text-muted-foreground">
              Configure core business parameters—projects, expenses, tax percentage, schedules, and statuses.
            </p>
          </div>
        </div>
      </div>

      {/* Types Section */}
      <div className="space-y-3">
        <span className="text-xs font-bold uppercase tracking-wider text-muted-foreground">
          Types
        </span>
        <div className="space-y-3">
          {[
            categories['projects'],
            categories['expense'],
            categories['tax'],
            categories['orders'],
            categories['session-timings'],
          ].map((cat) => (
            <Card
              key={cat.key}
              className="border border-border/80 shadow-2xs hover:border-primary/40 transition rounded-xl"
            >
              <CardContent className="p-4 flex items-center justify-between">
                <div>
                  <h3 className="text-sm font-semibold text-foreground">{cat.title}</h3>
                  <p className="text-xs text-muted-foreground mt-0.5">{cat.items.length} items configured</p>
                </div>
                <Button
                  variant="outline"
                  size="sm"
                  onClick={() => handleOpenCategory(cat.key)}
                  className="rounded-full px-5 text-xs font-semibold h-8"
                >
                  Manage
                </Button>
              </CardContent>
            </Card>
          ))}
        </div>
      </div>

      {/* Status Section */}
      <div className="space-y-3 pt-2">
        <span className="text-xs font-bold uppercase tracking-wider text-muted-foreground">
          Status
        </span>
        <Card className="border border-border/80 shadow-2xs hover:border-primary/40 transition rounded-xl">
          <CardContent className="p-4 flex items-center justify-between">
            <div>
              <h3 className="text-sm font-semibold text-foreground">
                {categories['projects-status'].title}
              </h3>
              <p className="text-xs text-muted-foreground mt-0.5">
                {categories['projects-status'].items.length} workflow statuses
              </p>
            </div>
            <Button
              variant="outline"
              size="sm"
              onClick={() => handleOpenCategory('projects-status')}
              className="rounded-full px-5 text-xs font-semibold h-8"
            >
              Manage
            </Button>
          </CardContent>
        </Card>
      </div>

      {/* Detail Drawer / Modal for each category */}
      {/* User prompt requirement:
          "for each will open a drawer which will have table with 3 columns:
           1. Type
           2. edit
           3. delete
           and at bottom a input field, and instead of edit/delete there will be a add button which will be col span 2 to cover both edit/delete rows to add new data"
      */}
      {activeCategory && (
        <Dialog open={!!activeCategoryKey} onOpenChange={(open) => !open && handleCloseCategory()}>
          <DialogContent className="max-w-xl">
            <DialogHeader>
              <DialogTitle className="text-base font-bold flex items-center gap-2">
                <IconDatabase className="h-4 w-4 text-primary" />
                Manage {activeCategory.title}
              </DialogTitle>
              <DialogDescription className="text-xs text-muted-foreground">
                {activeCategory.subtitle}
              </DialogDescription>
            </DialogHeader>

            <div className="py-2">
              <div className="rounded-xl border overflow-hidden">
                <table className="w-full text-xs">
                  <thead className="bg-muted/40 border-b">
                    <tr className="text-muted-foreground">
                      <th className="py-2.5 px-3 text-left font-semibold">Type</th>
                      <th className="py-2.5 px-3 text-center font-semibold w-16">Edit</th>
                      <th className="py-2.5 px-3 text-center font-semibold w-16">Delete</th>
                    </tr>
                  </thead>
                  <tbody className="divide-y divide-border/60">
                    {activeCategory.items.map((item) => (
                      <tr key={item.id} className="hover:bg-muted/15 transition">
                        <td className="py-2 px-3">
                          {editingItemId === item.id ? (
                            <div className="flex items-center gap-1.5">
                              <Input
                                value={editItemName}
                                onChange={(e) => setEditItemName(e.target.value)}
                                className="h-7 text-xs"
                                autoFocus
                              />
                              <Button
                                size="sm"
                                variant="ghost"
                                onClick={() => handleSaveEdit(item.id)}
                                className="h-7 w-7 p-0 text-emerald-600 hover:text-emerald-700"
                              >
                                <IconCheck className="h-3.5 w-3.5" />
                              </Button>
                              <Button
                                size="sm"
                                variant="ghost"
                                onClick={() => setEditingItemId(null)}
                                className="h-7 w-7 p-0 text-muted-foreground"
                              >
                                <IconX className="h-3.5 w-3.5" />
                              </Button>
                            </div>
                          ) : (
                            <span className="font-medium text-foreground">{item.name}</span>
                          )}
                        </td>
                        <td className="py-2 px-3 text-center">
                          <Button
                            variant="ghost"
                            size="sm"
                            onClick={() => handleStartEdit(item)}
                            className="h-7 w-7 p-0 text-muted-foreground hover:text-primary"
                          >
                            <IconEdit className="h-3.5 w-3.5" />
                          </Button>
                        </td>
                        <td className="py-2 px-3 text-center">
                          <Button
                            variant="ghost"
                            size="sm"
                            onClick={() => handleDeleteItem(item.id)}
                            className="h-7 w-7 p-0 text-muted-foreground hover:text-destructive"
                          >
                            <IconTrash className="h-3.5 w-3.5" />
                          </Button>
                        </td>
                      </tr>
                    ))}
                  </tbody>

                  {/* Bottom input row with Col Span 2 add button as explicitly requested */}
                  <tfoot className="border-t bg-muted/20">
                    <tr>
                      <td className="p-2.5">
                        <Input
                          value={newItemName}
                          onChange={(e) => setNewItemName(e.target.value)}
                          placeholder={`Enter new ${activeCategory.title.toLowerCase()} type...`}
                          className="h-8 text-xs bg-background"
                          onKeyDown={(e) => {
                            if (e.key === 'Enter') {
                              e.preventDefault()
                              handleAddItem()
                            }
                          }}
                        />
                      </td>
                      <td colSpan={2} className="p-2.5 text-center">
                        <Button
                          type="button"
                          onClick={handleAddItem}
                          className="w-full h-8 text-xs font-semibold bg-emerald-600 hover:bg-emerald-700 text-white gap-1 shadow-xs"
                        >
                          <IconPlus className="h-3.5 w-3.5" />
                          Add
                        </Button>
                      </td>
                    </tr>
                  </tfoot>
                </table>
              </div>
            </div>

            <DialogFooter>
              <Button type="button" variant="outline" size="sm" onClick={handleCloseCategory}>
                Done
              </Button>
            </DialogFooter>
          </DialogContent>
        </Dialog>
      )}
    </div>
  )
}
