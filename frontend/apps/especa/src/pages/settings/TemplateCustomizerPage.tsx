import * as React from 'react'
import {
  IconDeviceFloppy,
} from "@repo/ui/icons/tabler-icons"
import {
  Card,
  CardContent,
  CardHeader,
  CardTitle,
} from '@repo/ui/components/ui/card'
import { Button } from '@repo/ui/components/ui/button'
import { Input } from '@repo/ui/components/ui/input'
import { Label } from '@repo/ui/components/ui/label'
import { toast } from '@repo/ui/components/ui/sonner'

interface TemplateCustomizerProps {
  title: string
  subtitle: string
  type: 'form' | 'quotation' | 'invoice' | 'contract' | 'feedback'
}

export function TemplateCustomizerPage({
  title,
  subtitle,
  type,
}: TemplateCustomizerProps) {
  const [headerTitle, setHeaderTitle] = React.useState(`${title} - especa Studios`)
  const [termsText, setTermsText] = React.useState(
    '1. 50% advance upon contract signing.\n2. Delivery of color graded soft copies in 14 business days.\n3. Raw video footage retains copyrights by studio.'
  )
  const [accentColor, setAccentColor] = React.useState('#059669') // Emerald
  const [footerNote, setFooterNote] = React.useState('Thank you for choosing our visual artistry.')

  const handleSave = () => {
    toast.success(`${title} design & layout saved`)
  }

  return (
    <div className="space-y-6 max-w-4xl mx-auto pb-16">
      <div className="flex items-center justify-between border-b pb-4">
        <div>
          <h1 className="text-xl font-bold tracking-tight text-foreground">{title}</h1>
          <p className="text-xs text-muted-foreground">{subtitle}</p>
        </div>
        <Button
          type="button"
          onClick={handleSave}
          className="bg-emerald-600 hover:bg-emerald-700 text-white font-semibold text-xs h-8 gap-1.5 shadow-xs"
        >
          <IconDeviceFloppy className="h-3.5 w-3.5" />
          Save Template
        </Button>
      </div>

      <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
        {/* Configuration settings */}
        <div className="space-y-4">
          <Card className="border border-border/80 shadow-2xs rounded-xl">
            <CardHeader className="pb-3 border-b">
              <CardTitle className="text-xs font-bold text-foreground uppercase tracking-wider">
                Branding & Headings
              </CardTitle>
            </CardHeader>
            <CardContent className="p-4 space-y-3.5 text-xs">
              <div className="space-y-1">
                <Label className="text-xs font-semibold text-foreground">Document Header Title</Label>
                <Input
                  value={headerTitle}
                  onChange={(e) => setHeaderTitle(e.target.value)}
                  className="h-9 text-xs"
                />
              </div>

              <div className="space-y-1">
                <Label className="text-xs font-semibold text-foreground">Theme Accent Color</Label>
                <div className="flex items-center gap-3">
                  <input
                    type="color"
                    value={accentColor}
                    onChange={(e) => setAccentColor(e.target.value)}
                    className="h-8 w-14 rounded border cursor-pointer bg-background"
                  />
                  <span className="font-mono text-xs text-muted-foreground">{accentColor}</span>
                </div>
              </div>

              <div className="space-y-1">
                <Label className="text-xs font-semibold text-foreground">Footer / Sign-off Text</Label>
                <Input
                  value={footerNote}
                  onChange={(e) => setFooterNote(e.target.value)}
                  className="h-9 text-xs"
                />
              </div>
            </CardContent>
          </Card>

          <Card className="border border-border/80 shadow-2xs rounded-xl">
            <CardHeader className="pb-3 border-b">
              <CardTitle className="text-xs font-bold text-foreground uppercase tracking-wider">
                Terms & Conditions / Guidelines
              </CardTitle>
            </CardHeader>
            <CardContent className="p-4 space-y-2">
              <textarea
                value={termsText}
                onChange={(e) => setTermsText(e.target.value)}
                rows={5}
                className="w-full rounded-md border border-input bg-background px-3 py-2 text-xs shadow-xs focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring resize-none font-sans"
              />
            </CardContent>
          </Card>
        </div>

        {/* Live Preview Box */}
        <div className="space-y-2">
          <span className="text-xs font-bold uppercase tracking-wider text-muted-foreground">
            Live Preview
          </span>
          <div className="rounded-xl border border-border bg-white text-slate-900 p-6 shadow-sm min-h-[420px] flex flex-col justify-between text-xs">
            <div className="space-y-4">
              <div
                className="pb-3 border-b-2 flex items-center justify-between"
                style={{ borderColor: accentColor }}
              >
                <div>
                  <h3 className="font-bold text-sm tracking-tight text-slate-900">{headerTitle}</h3>
                  <p className="text-[10px] text-slate-500">Premium Photography & Cinema</p>
                </div>
                <div
                  className="px-2 py-0.5 rounded text-[10px] font-bold text-white uppercase"
                  style={{ backgroundColor: accentColor }}
                >
                  {type.toUpperCase()}
                </div>
              </div>

              <div className="bg-slate-50 p-3 rounded border text-[11px] space-y-1 text-slate-600">
                <div className="flex justify-between font-medium">
                  <span>Client: Priya & Rahul</span>
                  <span>Date: 09 Nov 2026</span>
                </div>
                <div className="flex justify-between">
                  <span>Project: Royal Wedding Cinema</span>
                  <span>Doc No: #ESPE-2026-001</span>
                </div>
              </div>

              <div className="space-y-1.5 pt-2">
                <span className="text-[10px] font-bold uppercase text-slate-400">Terms & Policy</span>
                <p className="text-[10px] text-slate-600 whitespace-pre-line leading-relaxed">
                  {termsText}
                </p>
              </div>
            </div>

            <div className="pt-4 border-t text-center text-[10px] text-slate-500">
              {footerNote}
            </div>
          </div>
        </div>
      </div>
    </div>
  )
}
