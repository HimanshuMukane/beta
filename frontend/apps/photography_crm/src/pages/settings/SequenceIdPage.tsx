import * as React from 'react'
import {
  IconDeviceFloppy,
  IconId,
} from "@repo/ui/icons/tabler-icons"
import {
  Card,
  CardContent,
  CardHeader,
  CardTitle,
  CardDescription,
} from '@repo/ui/components/ui/card'
import { Button } from '@repo/ui/components/ui/button'
import { Input } from '@repo/ui/components/ui/input'
import { Label } from '@repo/ui/components/ui/label'
import { toast } from '@repo/ui/components/ui/sonner'
import { apiClient } from '@repo/ui/lib/api-client'
import type { SequenceIdSettings, SequenceConfig } from '../../types/settings'

export function SequenceIdPage() {
  const [activeTab, setActiveTab] = React.useState<'sequence' | 'visibility'>('sequence')
  const [settings, setSettings] = React.useState<SequenceIdSettings>({
    leads: { prefix: 'ESPE/LQuote/', nextNumber: '101' },
    client: { prefix: 'ESPE/CLIENT/', nextNumber: '101' },
    orders: { prefix: 'ESPE/WALKIN/', nextNumber: '100' },
    expense: { prefix: 'ESPE/EX/', nextNumber: '100' },
    invoice: { prefix: 'ESPE/INVOICE/', nextNumber: '101' },
    projects: { prefix: 'ESPE/EVENT/', nextNumber: '101' },
    studioTaxNo: '27AABCE1234F1Z5',
    cancellation: { prefix: 'ESPE/CANCEL/', nextNumber: '100' },
    showClientTaxNo: true,
  })

  React.useEffect(() => {
    // Connect to sequence_id API table
    apiClient<any>('/api/sequence-id')
      .then((seqData) => {
        if (seqData && seqData.leads) {
          setSettings((prev) => ({
            ...prev,
            leads: seqData.leads || prev.leads,
            projects: seqData.projects || prev.projects,
            invoice: seqData.invoice || prev.invoice,
            orders: seqData.orders || prev.orders,
            expense: seqData.expense || prev.expense,
            cancellation: seqData.cancellation || prev.cancellation,
            client: seqData.client || prev.client,
          }))
        }
      })
      .catch((err: any) => {
        toast.error(err?.message || 'Failed to load sequence configuration')
      })
  }, [])

  const handleSequenceChange = (
    key: keyof Omit<SequenceIdSettings, 'showClientTaxNo' | 'studioTaxNo'>,
    field: keyof SequenceConfig,
    val: string
  ) => {
    setSettings((prev) => ({
      ...prev,
      [key]: {
        ...prev[key],
        [field]: val,
      },
    }))
  }

  const handleSave = async () => {
    try {
      await apiClient('/api/sequence-id', {
        method: 'POST',
        body: JSON.stringify({
          leads: settings.leads,
          projects: settings.projects,
          invoice: settings.invoice,
          orders: settings.orders,
          expense: settings.expense,
          cancellation: settings.cancellation,
          client: settings.client,
        }),
      })

      toast.success('Sequence & Visibility settings saved successfully')
    } catch (err: any) {
      toast.error(err.message || 'Failed to save settings')
    }
  }

  return (
    <div className="space-y-6 max-w-4xl mx-auto pb-16">
      <div className="flex flex-col sm:flex-row sm:items-center justify-between gap-4 border-b pb-4">
        <div className="flex items-center gap-3">
          <span className="flex h-10 w-10 shrink-0 items-center justify-center rounded-xl bg-primary/10 text-primary border border-primary/20">
            <IconId className="h-5 w-5" />
          </span>
          <div>
            <h1 className="text-xl font-bold tracking-tight text-foreground">Sequence ID & Tax Settings</h1>
            <p className="text-xs text-muted-foreground">
              Configure auto-increment prefixes and fiscal visibility on quotations and invoices.
            </p>
          </div>
        </div>

        <Button
          type="button"
          onClick={handleSave}
          className="h-8 text-xs font-semibold gap-1.5 bg-primary text-primary-foreground self-start sm:self-auto shadow-xs"
        >
          <IconDeviceFloppy className="h-3.5 w-3.5" />
          Save Changes
        </Button>
      </div>

      {/* Top Tab Pill Switcher (Matches Screenshot 3) */}
      <div className="flex justify-center pt-1">
        <div className="inline-flex rounded-xl bg-muted/60 p-1 border shadow-2xs">
          <button
            type="button"
            onClick={() => setActiveTab('sequence')}
            className={`px-6 py-2 text-xs font-bold rounded-lg transition ${activeTab === 'sequence'
                ? 'bg-background text-foreground shadow-xs'
                : 'text-muted-foreground hover:text-foreground'
              }`}
          >
            Sequence ID
          </button>
          <button
            type="button"
            onClick={() => setActiveTab('visibility')}
            className={`px-6 py-2 text-xs font-bold rounded-lg transition ${activeTab === 'visibility'
                ? 'bg-background text-foreground shadow-xs'
                : 'text-muted-foreground hover:text-foreground'
              }`}
          >
            Visibility
          </button>
        </div>
      </div>

      {/* SEQUENCE ID TAB (Matches Screenshot 3) */}
      {activeTab === 'sequence' && (
        <div className="space-y-3">
          {[
            { key: 'leads', label: 'Leads' },
            { key: 'projects', label: 'Projects' },
            { key: 'invoice', label: 'Invoice' },
            { key: 'orders', label: 'Orders' },
            { key: 'expense', label: 'Expense' },
            { key: 'cancellation', label: 'Cancellation' },
            { key: 'client', label: 'Client' },
          ].map(({ key, label }) => {
            const seqKey = key as keyof Omit<
              SequenceIdSettings,
              'showClientTaxNo' | 'studioTaxNo'
            >
            const config = settings[seqKey]

            return (
              <Card
                key={key}
                className="border border-border/80 shadow-2xs hover:border-primary/40 transition rounded-xl"
              >
                <CardContent className="p-4 flex items-center justify-between gap-4">
                  <span className="font-semibold text-xs text-foreground min-w-[100px]">
                    {label}
                  </span>

                  <div className="flex items-center gap-2 flex-1 max-w-md justify-end">
                    <Input
                      value={config.prefix}
                      onChange={(e) => handleSequenceChange(seqKey, 'prefix', e.target.value)}
                      placeholder="Prefix"
                      className="h-9 text-xs font-mono w-44 bg-background"
                    />
                    <Input
                      value={config.nextNumber}
                      onChange={(e) => handleSequenceChange(seqKey, 'nextNumber', e.target.value)}
                      placeholder="Start No"
                      className="h-9 text-xs font-mono w-28 bg-background text-center"
                    />
                  </div>
                </CardContent>
              </Card>
            )
          })}
        </div>
      )}

      {/* VISIBILITY TAB (User explicitly requested):
          "in sequence id
           visibility tab will have:-
           Show Client Tax No (toggle)
           Studio Tax No (input)"
      */}
      {activeTab === 'visibility' && (
        <div className="space-y-4">
          <Card className="border border-border/80 shadow-2xs rounded-xl">
            <CardHeader className="pb-3 border-b">
              <CardTitle className="text-sm font-bold text-foreground">Tax & Invoicing Visibility</CardTitle>
              <CardDescription className="text-xs text-muted-foreground">
                Control studio tax identifiers and whether customer GST/tax numbers appear on documents.
              </CardDescription>
            </CardHeader>
            <CardContent className="p-5 space-y-5">
              {/* Show Client Tax No Toggle */}
              <div className="flex items-center justify-between">
                <div className="space-y-0.5">
                  <Label className="text-xs font-bold text-foreground">Show Client Tax No</Label>
                  <p className="text-[11px] text-muted-foreground">
                    Display client GSTIN or VAT registration number on invoices and quotes
                  </p>
                </div>
                <button
                  type="button"
                  role="switch"
                  aria-checked={settings.showClientTaxNo}
                  onClick={() =>
                    setSettings((prev) => ({ ...prev, showClientTaxNo: !prev.showClientTaxNo }))
                  }
                  className={`relative inline-flex h-6 w-11 shrink-0 cursor-pointer rounded-full border-2 border-transparent transition-colors duration-200 ease-in-out focus:outline-none focus:ring-2 focus:ring-primary focus:ring-offset-2 ${settings.showClientTaxNo ? 'bg-emerald-600' : 'bg-muted'
                    }`}
                >
                  <span
                    className={`pointer-events-none inline-block h-5 w-5 transform rounded-full bg-white shadow-lg ring-0 transition duration-200 ease-in-out ${settings.showClientTaxNo ? 'translate-x-5' : 'translate-x-0'
                      }`}
                  />
                </button>
              </div>

              {/* Studio Tax No Input */}
              <div className="space-y-1.5 pt-3 border-t">
                <Label className="text-xs font-bold text-foreground">Studio Tax No</Label>
                <Input
                  value={settings.studioTaxNo}
                  onChange={(e) =>
                    setSettings((prev) => ({ ...prev, studioTaxNo: e.target.value }))
                  }
                  placeholder="e.g. 27AABCE1234F1Z5 (GSTIN / VAT No)"
                  className="h-9 text-xs font-mono max-w-md bg-background"
                />
                <p className="text-[10px] text-muted-foreground">
                  This tax registration number will be printed on all generated PDF invoices and receipts.
                </p>
              </div>
            </CardContent>
          </Card>
        </div>
      )}

      {/* Save Button */}
      <div className="flex justify-end pt-2">
        <Button
          type="button"
          onClick={handleSave}
          className="bg-emerald-600 hover:bg-emerald-700 text-white font-semibold text-xs px-6 h-9 shadow-xs gap-1.5"
        >
          <IconDeviceFloppy className="h-4 w-4" />
          Save Settings
        </Button>
      </div>
    </div>
  )
}
