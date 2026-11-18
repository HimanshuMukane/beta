import * as React from "react"
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogFooter,
} from "@repo/ui/components/ui/dialog"
import { Button } from "@repo/ui/components/ui/button"
import { Input } from "@repo/ui/components/ui/input"
import { Label } from "@repo/ui/components/ui/label"
import { Textarea } from "@repo/ui/components/ui/textarea"
import { toast } from "@repo/ui/components/ui/sonner"
import {
  IconBrandWhatsapp,
  IconMail,
  IconCopy,
  IconExternalLink,
} from "@repo/ui/icons/tabler-icons"
import {
  buildWhatsAppUrl,
  buildMailtoUrl,
  sanitizePhoneNumber,
} from "../lib/communicationUtils"
import { apiClient } from "@repo/ui/lib/api-client"

export interface ShareCommunicationModalProps {
  open: boolean
  onOpenChange: (open: boolean) => void
  title: string
  clientName: string
  clientPhone?: string
  clientEmail?: string
  defaultSubject: string
  defaultMessage: string
}

export function ShareCommunicationModal({
  open,
  onOpenChange,
  title,
  clientName,
  clientPhone = "",
  clientEmail = "",
  defaultSubject,
  defaultMessage,
}: ShareCommunicationModalProps) {
  const [phoneNumber, setPhoneNumber] = React.useState<string>(clientPhone)
  const [emailAddress, setEmailAddress] = React.useState<string>(clientEmail)
  const [subject, setSubject] = React.useState<string>(defaultSubject)
  const [message, setMessage] = React.useState<string>(defaultMessage)
  const [isSendingBackendEmail, setIsSendingBackendEmail] = React.useState<boolean>(false)

  React.useEffect(() => {
    setPhoneNumber(clientPhone)
    setEmailAddress(clientEmail)
    setSubject(defaultSubject)
    setMessage(defaultMessage)
  }, [clientPhone, clientEmail, defaultSubject, defaultMessage, open])

  const handleCopyMessage = () => {
    if (navigator && navigator.clipboard) {
      navigator.clipboard.writeText(message)
      toast.success("Message content copied to clipboard")
    }
  }

  const handleTriggerWhatsApp = () => {
    var sanitized: string = sanitizePhoneNumber("", phoneNumber)
    var url: string = buildWhatsAppUrl(sanitized, message)
    window.open(url, "_blank", "noopener,noreferrer")
    toast.success("Opening WhatsApp with pre-composed message")
    onOpenChange(false)
  }

  const handleTriggerMailto = () => {
    if (!emailAddress) {
      toast.error("Recipient email address is missing")
      return
    }
    var mailtoUrl: string = buildMailtoUrl(emailAddress, subject, message)
    window.location.href = mailtoUrl
    toast.success("Triggered default email client")
    onOpenChange(false)
  }

  const handleBackendSMTPSend = async () => {
    if (!emailAddress) {
      toast.error("Recipient email address is required")
      return
    }
    setIsSendingBackendEmail(true)
    try {
      await apiClient<{ status: string; message: string }>("/api/communication/send-email", {
        method: "POST",
        body: JSON.stringify({
          recipientEmail: emailAddress,
          recipientName: clientName,
          subject: subject,
          bodyText: message,
        }),
      })
      toast.success("Email dispatched successfully via backend SMTP!")
      onOpenChange(false)
    } catch {
      var fallbackUrl: string = buildMailtoUrl(emailAddress, subject, message)
      window.location.href = fallbackUrl
      toast.info("Direct SMTP server offline - opened system email client (mailto:) instead")
      onOpenChange(false)
    } finally {
      setIsSendingBackendEmail(false)
    }
  }

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="sm:max-w-xl">
        <DialogHeader>
          <DialogTitle className="flex items-center gap-2 text-base font-bold text-foreground">
            <IconMail className="h-4 w-4 text-primary" />
            <span>{title}</span>
          </DialogTitle>
          <p className="text-xs text-muted-foreground">
            Share quotes, confirmations, and billing summaries directly with {clientName || "the client"}.
          </p>
        </DialogHeader>

        <div className="space-y-4 py-2 text-xs">
          <div className="grid grid-cols-1 sm:grid-cols-2 gap-3">
            <div className="space-y-1.5">
              <Label className="text-xs font-semibold text-foreground flex items-center gap-1.5">
                <IconBrandWhatsapp className="h-3.5 w-3.5 text-emerald-600" />
                WhatsApp Phone
              </Label>
              <Input
                value={phoneNumber}
                onChange={(e) => setPhoneNumber(e.target.value)}
                placeholder="+91 9876543210"
                className="h-8 text-xs font-mono"
              />
            </div>
            <div className="space-y-1.5">
              <Label className="text-xs font-semibold text-foreground flex items-center gap-1.5">
                <IconMail className="h-3.5 w-3.5 text-blue-600" />
                Recipient Email
              </Label>
              <Input
                value={emailAddress}
                onChange={(e) => setEmailAddress(e.target.value)}
                placeholder="client@example.com"
                className="h-8 text-xs font-mono"
              />
            </div>
          </div>

          <div className="space-y-1.5">
            <Label className="text-xs font-semibold text-foreground">Subject Line</Label>
            <Input
              value={subject}
              onChange={(e) => setSubject(e.target.value)}
              placeholder="Quotation / Booking Details"
              className="h-8 text-xs"
            />
          </div>

          <div className="space-y-1.5">
            <div className="flex items-center justify-between">
              <Label className="text-xs font-semibold text-foreground">Message Summary</Label>
              <button
                type="button"
                onClick={handleCopyMessage}
                className="text-[11px] text-muted-foreground hover:text-foreground flex items-center gap-1 transition"
              >
                <IconCopy className="h-3 w-3" />
                Copy
              </button>
            </div>
            <Textarea
              rows={8}
              value={message}
              onChange={(e) => setMessage(e.target.value)}
              className="text-xs font-mono leading-relaxed resize-y"
            />
          </div>
        </div>

        <DialogFooter className="gap-2 sm:gap-2 flex-wrap items-center">
          <Button
            type="button"
            variant="outline"
            size="sm"
            onClick={() => onOpenChange(false)}
            className="text-xs h-8"
          >
            Cancel
          </Button>

          <Button
            type="button"
            size="sm"
            onClick={handleTriggerWhatsApp}
            className="bg-emerald-600 hover:bg-emerald-700 text-white text-xs h-8 gap-1.5"
          >
            <IconBrandWhatsapp className="h-3.5 w-3.5" />
            <span>Open WhatsApp</span>
            <IconExternalLink className="h-3 w-3 opacity-70" />
          </Button>

          <Button
            type="button"
            variant="secondary"
            size="sm"
            onClick={handleTriggerMailto}
            className="text-xs h-8 gap-1.5"
            title="Open in native email client"
          >
            <IconMail className="h-3.5 w-3.5 text-blue-600" />
            <span>Open mailto:</span>
          </Button>

          <Button
            type="button"
            size="sm"
            disabled={isSendingBackendEmail}
            onClick={handleBackendSMTPSend}
            className="bg-primary hover:bg-primary/90 text-primary-foreground text-xs h-8 gap-1.5"
            title="Send directly using backend SMTP service"
          >
            <IconMail className="h-3.5 w-3.5" />
            <span>{isSendingBackendEmail ? "Dispatching..." : "Dispatch Email (SMTP)"}</span>
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  )
}
