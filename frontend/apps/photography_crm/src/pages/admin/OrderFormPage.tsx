import * as React from 'react'
import {
  IconArrowLeft,
  IconPlus,
  IconTrash,
  IconDeviceFloppy,
} from '@tabler/icons-react'
import { Button } from '@repo/ui/components/ui/button'
import { Input } from '@repo/ui/components/ui/input'
import { Label } from '@repo/ui/components/ui/label'
import {
  Card,
  CardContent,
  CardHeader,
  CardTitle,
  CardDescription,
} from '@repo/ui/components/ui/card'
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@repo/ui/components/ui/select'
import { Badge } from '@repo/ui/components/ui/badge'
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogFooter,
} from '@repo/ui/components/ui/dialog'
import { toast } from '@repo/ui/components/ui/sonner'
import { cn } from '@repo/ui/lib/utils'
import type {
  Order,
  OrderLineItem,
  OrderExpense,
  OrderStage,
  OrderPaymentStatus,
  OrderPaymentRecord,
} from '../../types/order'
import { LeadSourceCombobox } from '../../components/LeadSourceCombobox'
import { PhoneInputField } from '../../components/PhoneInputField'
import type { LeadSource } from '../../types/lead'
import { useOrderPresets } from '../../hooks/useOrderPresets'

export const ORDER_TYPES: string[] = [
  'Walk-in',
  'Studio Session',
  'Print & Album',
  'Framing',
  'Commercial Shoot',
  'Passport/Visa Photo',
  'Custom Merch',
  'Online',
]

export const POPULAR_LINE_ITEMS = [
  { name: '12x36 Photobook (Matte Finish - 40 pgs)', cost: 12500 },
  { name: 'Acrylic Wall Frame (24x36)', cost: 6500 },
  { name: 'Canvas Gallery Wrap (20x30)', cost: 4200 },
  { name: 'Premium Wooden Photo Box (50 prints)', cost: 3500 },
  { name: 'Mini Pocket Books (Set of 3)', cost: 2800 },
  { name: 'Digital Soft Copies High-Res (USB Drive)', cost: 4500 },
  { name: 'Passport/Visa Print Sheet (8 copies)', cost: 450 },
  { name: 'Studio Portrait Session Print (12x18)', cost: 1800 },
  { name: 'Custom Calendar with 12 Shoots', cost: 2200 },
]

interface OrderFormPageProps {
  initialOrder?: Order | null
  onBack: () => void
  onSave: (order: Order) => void
  existingProjectOptions?: { id: string; name: string }[]
}

export function OrderFormPage({
  initialOrder,
  onBack,
  onSave,
  existingProjectOptions = [],
}: OrderFormPageProps) {
  // Client Details
  const [orderNumber] = React.useState(
    initialOrder?.orderNumber || `ESPE/WALKIN/${Math.floor(100 + Math.random() * 900)}`
  )
  const [clientName, setClientName] = React.useState(initialOrder?.clientName || '')
  const [countryCode, setCountryCode] = React.useState(initialOrder?.countryCode || '+91')
  const [mobileNo, setMobileNo] = React.useState(initialOrder?.mobileNo || '')
  const [email, setEmail] = React.useState(initialOrder?.email || '')
  const [address, setAddress] = React.useState(initialOrder?.address || '')
  const [leadSource, setLeadSource] = React.useState<string>(initialOrder?.leadSource || 'Instagram Ads')

  // Orders Details
  const [referenceNo, setReferenceNo] = React.useState(initialOrder?.referenceNo || '')

  const { presets } = useOrderPresets()
  const [orderDate, setOrderDate] = React.useState(
    initialOrder?.orderDate || new Date().toISOString().replace('T', ' ').substring(0, 16)
  )
  const [orderType, setOrderType] = React.useState<string>(initialOrder?.orderType || 'Print & Album')
  const [stage, setStage] = React.useState<OrderStage>(initialOrder?.stage || 'Order Confirmation')
  const [assignee, setAssignee] = React.useState<string>(initialOrder?.assignee || 'Amit Verma')
  const [relatedProjectId, setRelatedProjectId] = React.useState<string>(
    initialOrder?.relatedProjectId || ''
  )

  // Orders Items
  const [items, setItems] = React.useState<OrderLineItem[]>(() => {
    if (initialOrder?.items && initialOrder.items.length > 0) {
      return initialOrder.items
    }
    return [
      {
        id: `item-${Date.now()}`,
        name: '',
        quantity: 1,
        itemCost: 0,
        subTotal: 0,
      },
    ]
  })

  // Orders Expenses
  const [expenses, setExpenses] = React.useState<OrderExpense[]>(() => {
    if (initialOrder?.expenses && initialOrder.expenses.length > 0) {
      return initialOrder.expenses
    }
    return [
      {
        id: `exp-${Date.now()}`,
        date: new Date().toISOString().replace('T', ' ').substring(0, 16),
        name: '',
        amount: 0,
      },
    ]
  })

  // Notes & Discount / Tax
  const [notes, setNotes] = React.useState(initialOrder?.notes || '')
  const [discount, setDiscount] = React.useState<number>(initialOrder?.discount || 0)
  const [showDiscountInput, setShowDiscountInput] = React.useState((initialOrder?.discount || 0) > 0)
  const [taxPercent, setTaxPercent] = React.useState<number>(initialOrder?.taxPercent ?? 18)
  const [showTaxInput, setShowTaxInput] = React.useState((initialOrder?.taxPercent || 0) > 0)

  // Payments
  const [payments, setPayments] = React.useState<OrderPaymentRecord[]>(
    initialOrder?.payments || []
  )
  const [isPaymentModalOpen, setIsPaymentModalOpen] = React.useState(false)
  const [paymentAmountInput, setPaymentAmountInput] = React.useState<number>(0)
  const [paymentMethodInput, setPaymentMethodInput] = React.useState<
    'Cash' | 'Card' | 'UPI' | 'Bank Transfer' | 'Other'
  >('UPI')
  const [paymentRefInput, setPaymentRefInput] = React.useState('')

  // Line Item Calculations
  const handleItemChange = (
    id: string,
    field: keyof OrderLineItem,
    value: string | number
  ) => {
    setItems((prev) =>
      prev.map((item) => {
        if (item.id !== id) return item
        const updated = { ...item, [field]: value }
        if (field === 'quantity' || field === 'itemCost') {
          const qty = field === 'quantity' ? Number(value) || 0 : item.quantity
          const cost = field === 'itemCost' ? Number(value) || 0 : item.itemCost
          updated.subTotal = qty * cost
        }
        return updated
      })
    )
  }

  const handleSelectPredefinedItem = (id: string, selectedName: string) => {
    const presetItem = presets.find((p) => p.name === selectedName)
    const cost = presetItem ? presetItem.amount : (POPULAR_LINE_ITEMS.find((p) => p.name === selectedName)?.cost || 0)
    setItems((prev) =>
      prev.map((item) => {
        if (item.id !== id) return item
        return {
          ...item,
          name: selectedName,
          itemCost: cost,
          subTotal: (item.quantity || 1) * cost,
        }
      })
    )
  }

  const handleAddItem = () => {
    setItems((prev) => [
      ...prev,
      {
        id: `item-${Date.now()}`,
        name: '',
        quantity: 1,
        itemCost: 0,
        subTotal: 0,
      },
    ])
  }

  const handleRemoveItem = (id: string) => {
    if (items.length <= 1) {
      toast.info('Order must have at least one line item')
      return
    }
    setItems((prev) => prev.filter((it) => it.id !== id))
  }

  // Expense Handlers
  const handleExpenseChange = (
    id: string,
    field: keyof OrderExpense,
    value: string | number
  ) => {
    setExpenses((prev) =>
      prev.map((exp) => {
        if (exp.id !== id) return exp
        return { ...exp, [field]: value }
      })
    )
  }

  const handleAddExpense = () => {
    setExpenses((prev) => [
      ...prev,
      {
        id: `exp-${Date.now()}`,
        date: new Date().toISOString().replace('T', ' ').substring(0, 16),
        name: '',
        amount: 0,
      },
    ])
  }

  const handleRemoveExpense = (id: string) => {
    setExpenses((prev) => prev.filter((exp) => exp.id !== id))
  }

  // Financial Computations
  const subTotal = React.useMemo(() => {
    return items.reduce((sum, item) => sum + (Number(item.subTotal) || 0), 0)
  }, [items])

  const effectiveTaxAmount = React.useMemo(() => {
    if (!showTaxInput || taxPercent <= 0) return 0
    const taxableAmount = Math.max(0, subTotal - discount)
    return Math.round((taxableAmount * taxPercent) / 100)
  }, [subTotal, discount, showTaxInput, taxPercent])

  const totalAmount = React.useMemo(() => {
    const afterDiscount = Math.max(0, subTotal - (showDiscountInput ? discount : 0))
    return afterDiscount + effectiveTaxAmount
  }, [subTotal, discount, showDiscountInput, effectiveTaxAmount])

  const paidAmount = React.useMemo(() => {
    return payments.reduce((sum, p) => sum + (Number(p.amount) || 0), 0)
  }, [payments])

  const balanceAmount = React.useMemo(() => {
    return Math.max(0, totalAmount - paidAmount)
  }, [totalAmount, paidAmount])

  const paymentStatus: OrderPaymentStatus = React.useMemo(() => {
    if (totalAmount === 0) return 'Paid'
    if (paidAmount >= totalAmount) return 'Paid'
    if (paidAmount > 0) return 'Partially Paid'
    return 'Pending'
  }, [totalAmount, paidAmount])

  // Payment dialog submit
  const handleMakePayment = () => {
    if (paymentAmountInput <= 0) {
      toast.error('Payment amount must be greater than 0')
      return
    }
    const newPayment: OrderPaymentRecord = {
      id: `pay-${Date.now()}`,
      amount: paymentAmountInput,
      date: new Date().toISOString().replace('T', ' ').substring(0, 16),
      method: paymentMethodInput,
      reference: paymentRefInput || undefined,
    }
    setPayments((prev) => [...prev, newPayment])
    setIsPaymentModalOpen(false)
    setPaymentAmountInput(0)
    setPaymentRefInput('')
    toast.success(`Payment of ₹${newPayment.amount.toLocaleString('en-IN')} recorded`)
  }

  const handleSaveOrder = (e: React.FormEvent) => {
    e.preventDefault()

    if (!clientName.trim()) {
      toast.error('Client name is required')
      return
    }
    if (!mobileNo.trim()) {
      toast.error('Mobile number is required')
      return
    }

    const orderToSave: Order = {
      id: initialOrder?.id || `ord-${Date.now()}`,
      orderNumber,
      clientName: clientName.trim(),
      countryCode,
      mobileNo: mobileNo.trim(),
      email: email.trim() || undefined,
      address: address.trim() || undefined,
      leadSource,
      referenceNo: referenceNo.trim() || undefined,
      orderDate,
      orderType,
      stage,
      assignee,
      items: items.filter((it) => it.name.trim() !== '' || it.itemCost > 0),
      expenses: expenses.filter((ex) => ex.name.trim() !== '' || ex.amount > 0),
      payments,
      notes: notes.trim() || undefined,
      discount: showDiscountInput ? discount : 0,
      taxPercent: showTaxInput ? taxPercent : 0,
      taxAmount: effectiveTaxAmount,
      subTotal,
      totalAmount,
      paidAmount,
      balanceAmount,
      paymentStatus,
      createdDate: initialOrder?.createdDate || new Date().toISOString().split('T')[0],
      relatedProjectId: relatedProjectId || undefined,
    }

    onSave(orderToSave)
    toast.success(initialOrder ? 'Order updated successfully' : 'New order created successfully')
  }

  return (
    <form onSubmit={handleSaveOrder} className="flex flex-col gap-6 max-w-6xl mx-auto pb-16">
      {/* Top action bar */}
      <div className="flex items-center justify-between border-b pb-4">
        <div className="flex items-center gap-3">
          <Button
            type="button"
            variant="ghost"
            size="sm"
            onClick={onBack}
            className="flex items-center gap-1.5 text-muted-foreground hover:text-foreground"
          >
            <IconArrowLeft className="h-4 w-4" />
            Back
          </Button>
          <div className="flex items-center gap-2">
            <span className="text-xl font-bold tracking-tight text-foreground">
              # {orderNumber}
            </span>
            <Badge variant="outline" className="text-xs bg-muted/40">
              {stage}
            </Badge>
          </div>
        </div>

        <div className="flex items-center gap-2">
          <Button type="button" variant="outline" size="sm" onClick={onBack}>
            Cancel
          </Button>
          <Button
            type="submit"
            size="sm"
            className="bg-emerald-600 hover:bg-emerald-700 text-white font-medium gap-1.5 shadow-xs"
          >
            <IconDeviceFloppy className="h-4 w-4" />
            Save Order
          </Button>
        </div>
      </div>

      {/* Grid: Client Details & Orders Details (Screenshot 2 Match) */}
      <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
        {/* Client Details */}
        <Card className="shadow-2xs border-border/80">
          <CardHeader className="pb-3 border-b border-border/40">
            <CardTitle className="text-sm font-bold text-foreground">Client Details</CardTitle>
          </CardHeader>
          <CardContent className="space-y-3.5 pt-4">
            <div className="space-y-1">
              <Label className="text-xs font-semibold text-muted-foreground">
                NAME <span className="text-destructive">*</span>
              </Label>
              <Input
                value={clientName}
                onChange={(e) => setClientName(e.target.value)}
                placeholder="Name"
                className="h-9 text-xs"
                required
              />
            </div>

            <div className="space-y-1">
              <Label className="text-xs font-semibold text-muted-foreground">
                MOBILE NO <span className="text-destructive">*</span>
              </Label>
              <PhoneInputField
                countryCode={countryCode}
                mobileNo={mobileNo}
                onCountryCodeChange={setCountryCode}
                onMobileNoChange={setMobileNo}
                placeholder="Mobile No"
                required
              />
            </div>

            <div className="space-y-1">
              <Label className="text-xs font-semibold text-muted-foreground">EMAIL ADDRESS</Label>
              <Input
                type="email"
                value={email}
                onChange={(e) => setEmail(e.target.value)}
                placeholder="Email Address"
                className="h-9 text-xs"
              />
            </div>

            <div className="space-y-1">
              <Label className="text-xs font-semibold text-muted-foreground">ADDRESS</Label>
              <textarea
                value={address}
                onChange={(e: React.ChangeEvent<HTMLTextAreaElement>) => setAddress(e.target.value)}
                placeholder="Address"
                rows={3}
                className="w-full rounded-md border border-input bg-background px-3 py-2 text-xs shadow-xs focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring resize-none"
              />
            </div>

            <div className="space-y-1">
              <Label className="text-xs font-semibold text-muted-foreground">
                LEAD SOURCE <span className="text-destructive">*</span>
              </Label>
              <LeadSourceCombobox
                value={leadSource}
                onChange={(source: LeadSource) => setLeadSource(source)}
                placeholder="Select Lead Source"
                className="h-9 text-xs"
              />
            </div>
          </CardContent>
        </Card>

        {/* Orders Details */}
        <Card className="shadow-2xs border-border/80">
          <CardHeader className="pb-3 border-b border-border/40">
            <CardTitle className="text-sm font-bold text-foreground">Orders Details</CardTitle>
          </CardHeader>
          <CardContent className="space-y-3.5 pt-4">
            <div className="space-y-1">
              <Label className="text-xs font-semibold text-muted-foreground">REFERENCE NO</Label>
              <Input
                value={referenceNo}
                onChange={(e) => setReferenceNo(e.target.value)}
                placeholder="Reference No"
                className="h-9 text-xs"
              />
            </div>

            <div className="space-y-1">
              <Label className="text-xs font-semibold text-muted-foreground">
                ORDERS DATE <span className="text-destructive">*</span>
              </Label>
              <Input
                type="text"
                value={orderDate}
                onChange={(e) => setOrderDate(e.target.value)}
                placeholder="YYYY-MM-DD HH:mm"
                className="h-9 text-xs"
                required
              />
            </div>

            <div className="space-y-1">
              <Label className="text-xs font-semibold text-muted-foreground">
                ORDERS TYPE <span className="text-destructive">*</span>
              </Label>
              <Select value={orderType} onValueChange={setOrderType}>
                <SelectTrigger className="h-9 text-xs">
                  <SelectValue placeholder="Orders Type" />
                </SelectTrigger>
                <SelectContent>
                  {ORDER_TYPES.map((t) => (
                    <SelectItem key={t} value={t} className="text-xs">
                      {t}
                    </SelectItem>
                  ))}
                </SelectContent>
              </Select>
            </div>

            <div className="space-y-1">
              <Label className="text-xs font-semibold text-muted-foreground">STAGE</Label>
              <Select value={stage} onValueChange={(val) => setStage(val as OrderStage)}>
                <SelectTrigger className="h-9 text-xs">
                  <SelectValue placeholder="Stage" />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="Order Confirmation" className="text-xs">
                    Order Confirmation
                  </SelectItem>
                  <SelectItem value="Product Ready" className="text-xs">
                    Product Ready
                  </SelectItem>
                  <SelectItem value="Ready For Delivery" className="text-xs">
                    Ready For Delivery
                  </SelectItem>
                  <SelectItem value="Product Delivered" className="text-xs">
                    Product Delivered
                  </SelectItem>
                </SelectContent>
              </Select>
            </div>

            <div className="space-y-1">
              <Label className="text-xs font-semibold text-muted-foreground">ASSIGNEE</Label>
              <Select value={assignee} onValueChange={setAssignee}>
                <SelectTrigger className="h-9 text-xs">
                  <SelectValue placeholder="Assignee" />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="Amit Verma" className="text-xs">
                    Amit Verma (Lead Manager)
                  </SelectItem>
                  <SelectItem value="Pooja Hegde" className="text-xs">
                    Pooja Hegde (Studio Lead)
                  </SelectItem>
                  <SelectItem value="Vikram Joshi" className="text-xs">
                    Vikram Joshi (Lab Specialist)
                  </SelectItem>
                  <SelectItem value="Rohan Sen" className="text-xs">
                    Rohan Sen (Operations)
                  </SelectItem>
                </SelectContent>
              </Select>
            </div>

            {existingProjectOptions.length > 0 && (
              <div className="space-y-1 pt-1">
                <Label className="text-xs font-semibold text-muted-foreground">
                  LINK TO PROJECT (OPTIONAL)
                </Label>
                <Select value={relatedProjectId} onValueChange={setRelatedProjectId}>
                  <SelectTrigger className="h-9 text-xs">
                    <SelectValue placeholder="Select Related Project (if any)" />
                  </SelectTrigger>
                  <SelectContent>
                    <SelectItem value="" className="text-xs">
                      None (Standalone Order)
                    </SelectItem>
                    {existingProjectOptions.map((proj) => (
                      <SelectItem key={proj.id} value={proj.id} className="text-xs">
                        {proj.name}
                      </SelectItem>
                    ))}
                  </SelectContent>
                </Select>
              </div>
            )}
          </CardContent>
        </Card>
      </div>

      {/* Orders Items Section (Screenshot 2 Match) */}
      <Card className="shadow-2xs border-border/80">
        <CardHeader className="pb-3 border-b border-border/40 flex flex-row items-center justify-between">
          <div>
            <CardTitle className="text-sm font-bold text-foreground">Orders Items</CardTitle>
            <CardDescription className="text-xs text-muted-foreground">
              Add individual products, prints, frames, or custom deliverables
            </CardDescription>
          </div>
          <Button
            type="button"
            variant="outline"
            size="sm"
            onClick={handleAddItem}
            className="text-xs gap-1 h-8 text-primary border-primary/30 hover:bg-primary/5"
          >
            <IconPlus className="h-3.5 w-3.5" />
            Add Items
          </Button>
        </CardHeader>
        <CardContent className="pt-4">
          <div className="overflow-x-auto">
            <table className="w-full text-xs">
              <thead>
                <tr className="border-b bg-muted/30 text-muted-foreground text-left">
                  <th className="py-2.5 px-3 font-semibold rounded-l-md">Orders Line Item</th>
                  <th className="py-2.5 px-3 font-semibold w-24">Quantity</th>
                  <th className="py-2.5 px-3 font-semibold w-32">Item Cost (₹)</th>
                  <th className="py-2.5 px-3 font-semibold w-32">Sub Total (₹)</th>
                  <th className="py-2.5 px-3 font-semibold w-12 text-center rounded-r-md"></th>
                </tr>
              </thead>
              <tbody className="divide-y divide-border/40">
                {items.map((item) => (
                  <tr key={item.id} className="hover:bg-muted/10">
                    <td className="py-2.5 px-3">
                      <div className="flex flex-col gap-1.5">
                        <Input
                          value={item.name}
                          onChange={(e) => handleItemChange(item.id, 'name', e.target.value)}
                          placeholder="Line item description (e.g. 12x36 Photobook)"
                          className="h-8 text-xs"
                        />
                        {/* Quick preset selector */}
                        <div className="flex items-center gap-1.5 flex-wrap">
                          <span className="text-[10px] text-muted-foreground">Presets:</span>
                          {(presets.length > 0
                            ? presets.map((p) => ({ name: p.name, cost: p.amount }))
                            : POPULAR_LINE_ITEMS
                          )
                            .slice(0, 5)
                            .map((preset) => (
                              <button
                                key={preset.name}
                                type="button"
                                onClick={() => handleSelectPredefinedItem(item.id, preset.name)}
                                className="text-[10px] text-primary hover:underline bg-primary/5 px-1.5 py-0.5 rounded border border-primary/20"
                              >
                                {preset.name.split('(')[0]}
                              </button>
                            ))}
                        </div>
                      </div>
                    </td>
                    <td className="py-2.5 px-3 align-top">
                      <Input
                        type="number"
                        min="1"
                        value={item.quantity}
                        onChange={(e) =>
                          handleItemChange(item.id, 'quantity', parseInt(e.target.value) || 1)
                        }
                        className="h-8 text-xs text-center"
                      />
                    </td>
                    <td className="py-2.5 px-3 align-top">
                      <Input
                        type="number"
                        min="0"
                        value={item.itemCost}
                        onChange={(e) =>
                          handleItemChange(item.id, 'itemCost', parseFloat(e.target.value) || 0)
                        }
                        className="h-8 text-xs text-right"
                      />
                    </td>
                    <td className="py-2.5 px-3 align-top">
                      <div className="h-8 flex items-center justify-end font-semibold text-foreground px-2 bg-muted/20 rounded border border-border/40">
                        ₹{item.subTotal.toLocaleString('en-IN')}
                      </div>
                    </td>
                    <td className="py-2.5 px-3 text-center align-top">
                      <Button
                        type="button"
                        variant="ghost"
                        size="sm"
                        onClick={() => handleRemoveItem(item.id)}
                        className="h-8 w-8 p-0 text-muted-foreground hover:text-destructive"
                      >
                        <IconTrash className="h-4 w-4" />
                      </Button>
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        </CardContent>
      </Card>

      {/* Note & Financial Summary Section (Screenshot 2 Match) */}
      <div className="grid grid-cols-1 md:grid-cols-2 gap-6 items-start">
        {/* Note */}
        <Card className="shadow-2xs border-border/80">
          <CardHeader className="pb-2 border-b border-border/40">
            <CardTitle className="text-xs font-bold text-muted-foreground uppercase">
              Note
            </CardTitle>
          </CardHeader>
          <CardContent className="pt-3">
            <textarea
              value={notes}
              onChange={(e: React.ChangeEvent<HTMLTextAreaElement>) => setNotes(e.target.value)}
              placeholder="Add client instructions, album binding preferences, delivery notes..."
              rows={4}
              className="w-full rounded-md border border-input bg-background px-3 py-2 text-xs shadow-xs focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring resize-none"
            />
          </CardContent>
        </Card>

        {/* Financial Summary */}
        <div className="bg-card rounded-xl border border-border/80 p-5 shadow-2xs space-y-4">
          <div className="flex items-center justify-between text-sm">
            <span className="font-medium text-muted-foreground">Sub Total</span>
            <span className="font-bold text-foreground">₹{subTotal.toLocaleString('en-IN')}</span>
          </div>

          {/* Discount toggle / input */}
          {showDiscountInput ? (
            <div className="flex items-center justify-between gap-3 text-xs">
              <span className="text-muted-foreground">Discount (₹)</span>
              <div className="flex items-center gap-2">
                <Input
                  type="number"
                  min="0"
                  value={discount}
                  onChange={(e) => setDiscount(Math.max(0, parseFloat(e.target.value) || 0))}
                  className="h-7 w-24 text-right text-xs"
                />
                <Button
                  type="button"
                  variant="ghost"
                  size="sm"
                  onClick={() => {
                    setDiscount(0)
                    setShowDiscountInput(false)
                  }}
                  className="h-7 px-1.5 text-xs text-muted-foreground"
                >
                  Remove
                </Button>
              </div>
            </div>
          ) : (
            <div className="flex justify-end">
              <button
                type="button"
                onClick={() => setShowDiscountInput(true)}
                className="text-xs font-semibold text-emerald-600 hover:text-emerald-700 dark:text-emerald-400 hover:underline"
              >
                + Add Discount
              </button>
            </div>
          )}

          {/* Tax toggle / input */}
          {showTaxInput ? (
            <div className="flex items-center justify-between gap-3 text-xs">
              <span className="text-muted-foreground">Tax ({taxPercent}% GST)</span>
              <div className="flex items-center gap-2">
                <Input
                  type="number"
                  min="0"
                  max="100"
                  value={taxPercent}
                  onChange={(e) => setTaxPercent(parseFloat(e.target.value) || 0)}
                  className="h-7 w-16 text-right text-xs"
                />
                <span className="font-semibold text-foreground">
                  ₹{effectiveTaxAmount.toLocaleString('en-IN')}
                </span>
                <Button
                  type="button"
                  variant="ghost"
                  size="sm"
                  onClick={() => {
                    setTaxPercent(0)
                    setShowTaxInput(false)
                  }}
                  className="h-7 px-1.5 text-xs text-muted-foreground"
                >
                  Remove
                </Button>
              </div>
            </div>
          ) : (
            <div className="flex justify-end">
              <button
                type="button"
                onClick={() => {
                  setTaxPercent(18)
                  setShowTaxInput(true)
                }}
                className="text-xs font-semibold text-emerald-600 hover:text-emerald-700 dark:text-emerald-400 hover:underline"
              >
                + Add Tax
              </button>
            </div>
          )}

          <div className="pt-3 border-t border-border/80 flex items-center justify-between">
            <span className="text-base font-bold text-foreground">Total</span>
            <span className="text-xl font-extrabold text-foreground tracking-tight">
              ₹{totalAmount.toLocaleString('en-IN')}
            </span>
          </div>

          <div className="flex items-center justify-between text-xs text-muted-foreground pt-1">
            <span>Paid: ₹{paidAmount.toLocaleString('en-IN')}</span>
            <span
              className={cn(
                'font-bold',
                balanceAmount > 0
                  ? 'text-amber-600 dark:text-amber-400'
                  : 'text-emerald-600 dark:text-emerald-400'
              )}
            >
              Balance: ₹{balanceAmount.toLocaleString('en-IN')}
            </span>
          </div>

          <div className="pt-2">
            <Button
              type="button"
              onClick={() => {
                setPaymentAmountInput(balanceAmount > 0 ? balanceAmount : totalAmount)
                setIsPaymentModalOpen(true)
              }}
              className="w-full bg-emerald-600 hover:bg-emerald-700 text-white font-semibold text-xs shadow-xs"
            >
              Make Payment
            </Button>
          </div>
        </div>
      </div>

      {/* Orders Expense Section (Screenshot 2 Match) */}
      <Card className="shadow-2xs border-border/80">
        <CardHeader className="pb-3 border-b border-border/40 flex flex-row items-center justify-between">
          <div>
            <CardTitle className="text-sm font-bold text-foreground">Orders Expense</CardTitle>
            <CardDescription className="text-xs text-muted-foreground">
              Track album printing costs, lab processing, framing materials, or logistics
            </CardDescription>
          </div>
          <Button
            type="button"
            variant="outline"
            size="sm"
            onClick={handleAddExpense}
            className="text-xs gap-1 h-8 text-emerald-600 border-emerald-600/30 hover:bg-emerald-50 dark:hover:bg-emerald-950/20"
          >
            <IconPlus className="h-3.5 w-3.5" />
            Add Expense
          </Button>
        </CardHeader>
        <CardContent className="pt-4 space-y-3">
          {expenses.length === 0 ? (
            <p className="text-xs text-muted-foreground italic">No expenses recorded for this order.</p>
          ) : (
            expenses.map((expense) => (
              <div key={expense.id} className="flex items-center gap-2">
                <Input
                  value={expense.date}
                  onChange={(e) => handleExpenseChange(expense.id, 'date', e.target.value)}
                  placeholder="YYYY-MM-DD HH:mm"
                  className="w-44 h-8 text-xs"
                />
                <Input
                  value={expense.name}
                  onChange={(e) => handleExpenseChange(expense.id, 'name', e.target.value)}
                  placeholder="Expense Name (e.g. Lab Printing Cost)"
                  className="flex-1 h-8 text-xs"
                />
                <Input
                  type="number"
                  min="0"
                  value={expense.amount}
                  onChange={(e) =>
                    handleExpenseChange(expense.id, 'amount', parseFloat(e.target.value) || 0)
                  }
                  placeholder="Amount"
                  className="w-28 h-8 text-xs text-right"
                />
                <Button
                  type="button"
                  variant="ghost"
                  size="sm"
                  onClick={() => handleRemoveExpense(expense.id)}
                  className="h-8 w-8 p-0 text-muted-foreground hover:text-destructive"
                >
                  <IconTrash className="h-4 w-4" />
                </Button>
              </div>
            ))
          )}
        </CardContent>
      </Card>

      {/* Payment Dialog Modal */}
      <Dialog open={isPaymentModalOpen} onOpenChange={setIsPaymentModalOpen}>
        <DialogContent className="max-w-md">
          <DialogHeader>
            <DialogTitle className="text-base font-bold">Record Payment</DialogTitle>
          </DialogHeader>

          <div className="space-y-4 py-2 text-xs">
            <div className="rounded-lg bg-muted/40 p-3 flex justify-between items-center">
              <span className="text-muted-foreground">Order Balance Due</span>
              <span className="font-bold text-sm text-foreground">
                ₹{balanceAmount.toLocaleString('en-IN')}
              </span>
            </div>

            <div className="space-y-1">
              <Label className="text-xs font-semibold">Payment Amount (₹)</Label>
              <Input
                type="number"
                min="1"
                value={paymentAmountInput}
                onChange={(e) => setPaymentAmountInput(parseFloat(e.target.value) || 0)}
                className="h-9 text-xs"
              />
            </div>

            <div className="space-y-1">
              <Label className="text-xs font-semibold">Payment Method</Label>
              <Select
                value={paymentMethodInput}
                onValueChange={(val) => setPaymentMethodInput(val as any)}
              >
                <SelectTrigger className="h-9 text-xs">
                  <SelectValue />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="UPI">UPI / QR Code</SelectItem>
                  <SelectItem value="Card">Credit / Debit Card</SelectItem>
                  <SelectItem value="Cash">Cash</SelectItem>
                  <SelectItem value="Bank Transfer">Bank Transfer (NEFT/IMPS)</SelectItem>
                  <SelectItem value="Other">Other</SelectItem>
                </SelectContent>
              </Select>
            </div>

            <div className="space-y-1">
              <Label className="text-xs font-semibold">Transaction / Reference Note</Label>
              <Input
                value={paymentRefInput}
                onChange={(e) => setPaymentRefInput(e.target.value)}
                placeholder="e.g. UPI/283918239 or Cash receipt"
                className="h-9 text-xs"
              />
            </div>
          </div>

          <DialogFooter>
            <Button
              type="button"
              variant="outline"
              size="sm"
              onClick={() => setIsPaymentModalOpen(false)}
            >
              Cancel
            </Button>
            <Button
              type="button"
              size="sm"
              onClick={handleMakePayment}
              className="bg-emerald-600 hover:bg-emerald-700 text-white font-medium"
            >
              Confirm Payment
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>
    </form>
  )
}
