import type { Lead } from "../types/lead"
import type { Order } from "../types/order"
import type { BillingInvoice } from "../types/billing"
import type { ProjectSession } from "../types/session"

export function sanitizePhoneNumber(countryCode?: string, mobileNo?: string): string {
  var rawNumber: string = (mobileNo || "").trim()
  var digits: string = rawNumber.replace(/\D/g, "")
  if (digits === "") {
    return ""
  }
  var codeDigits: string = (countryCode || "").replace(/\D/g, "")
  if (digits.length === 10 && codeDigits !== "") {
    return codeDigits + digits
  }
  if (digits.length === 10 && codeDigits === "") {
    return "91" + digits
  }
  return digits
}

export function buildWhatsAppUrl(phoneNumber: string, message: string): string {
  var cleanPhone: string = phoneNumber.replace(/\D/g, "")
  var encodedText: string = encodeURIComponent(message)
  if (cleanPhone !== "") {
    return "https://wa.me/" + cleanPhone + "?text=" + encodedText
  }
  return "https://wa.me/?text=" + encodedText
}

export function buildMailtoUrl(emailAddress: string, subject: string, body: string): string {
  var trimmedEmail: string = (emailAddress || "").trim()
  var encodedSubject: string = encodeURIComponent(subject)
  var encodedBody: string = encodeURIComponent(body)
  return "mailto:" + trimmedEmail + "?subject=" + encodedSubject + "&body=" + encodedBody
}

export function formatLeadQuotationSummary(lead: Lead): string {
  var messageLines: string[] = []
  messageLines.push("Hello " + (lead.name || "Client") + ",")
  messageLines.push("")
  messageLines.push("Thank you for contacting us regarding your upcoming " + (lead.projectName || lead.projectType || "photography") + " project.")
  messageLines.push("")
  messageLines.push("Here is a summary of your quotation:")
  messageLines.push("• Package/Type: " + (lead.packageName || lead.quoteType || "Standard Package"))
  if (lead.projectStartDate) {
    messageLines.push("• Start Date: " + lead.projectStartDate + (lead.projectEndDate ? " to " + lead.projectEndDate : ""))
  }
  if (lead.projectAddress) {
    messageLines.push("• Venue/Location: " + lead.projectAddress)
  }
  messageLines.push("• Estimated Total: ₹" + (lead.quoteAmount || 0).toLocaleString("en-IN"))
  if (lead.selectedDeliverables && lead.selectedDeliverables.length > 0) {
    messageLines.push("• Deliverables Included: " + lead.selectedDeliverables.join(", "))
  }
  messageLines.push("")
  messageLines.push("Please let us know if you would like to proceed or require any customizations.")
  messageLines.push("")
  messageLines.push("Warm regards,")
  messageLines.push(lead.assignee?.name || "Studio Team")
  return messageLines.join("\n")
}

export function formatOrderQuotationSummary(order: Order): string {
  var messageLines: string[] = []
  messageLines.push("Hello " + (order.clientName || "Client") + ",")
  messageLines.push("")
  messageLines.push("Here is the booking & order summary for " + order.orderNumber + ":")
  messageLines.push("• Order Date: " + order.orderDate)
  messageLines.push("• Order Type: " + order.orderType)
  messageLines.push("• Stage: " + order.stage)
  messageLines.push("• Total Amount: ₹" + order.totalAmount.toLocaleString("en-IN"))
  messageLines.push("• Paid Amount: ₹" + order.paidAmount.toLocaleString("en-IN"))
  messageLines.push("• Balance Due: ₹" + order.balanceAmount.toLocaleString("en-IN"))
  if (order.items && order.items.length > 0) {
    messageLines.push("• Items:")
    var itemIndex: number = 0
    for (itemIndex = 0; itemIndex < order.items.length; itemIndex++) {
      var lineItem = order.items[itemIndex]
      messageLines.push("   - " + lineItem.name + " (Qty: " + lineItem.quantity + ") - ₹" + lineItem.subTotal.toLocaleString("en-IN"))
    }
  }
  messageLines.push("")
  messageLines.push("Thank you for choosing our studio!")
  return messageLines.join("\n")
}

export function formatInvoiceSummary(invoice: BillingInvoice): string {
  var messageLines: string[] = []
  messageLines.push("Hello " + (invoice.clientName || "Client") + ",")
  messageLines.push("")
  messageLines.push("Here is your invoice summary from the studio:")
  messageLines.push("• Invoice Number: " + invoice.id)
  messageLines.push("• Description: " + invoice.invoiceName)
  if (invoice.projectsName) {
    messageLines.push("• Project: " + invoice.projectsName)
  }
  messageLines.push("• Invoice Date: " + invoice.invoiceDate)
  messageLines.push("• Total Amount: ₹" + invoice.totalAmount.toLocaleString("en-IN"))
  messageLines.push("• Amount Paid: ₹" + invoice.paidAmount.toLocaleString("en-IN"))
  messageLines.push("• Amount Pending: ₹" + invoice.pendingAmount.toLocaleString("en-IN"))
  messageLines.push("• Payment Status: " + invoice.paymentStatus)
  messageLines.push("")
  messageLines.push("Please let us know once the payment is made. Thank you!")
  return messageLines.join("\n")
}

export function formatSessionConfirmation(session: ProjectSession): string {
  var messageLines: string[] = []
  messageLines.push("Photoshoot Confirmation: " + (session.title || session.projectName || "Photography Session"))
  messageLines.push("")
  messageLines.push("Dear Client,")
  messageLines.push("Here are the confirmed details for your upcoming photoshoot:")
  messageLines.push("• Date: " + session.date)
  messageLines.push("• Call Time / Timing: " + (session.timeRange || session.startTime || "TBD") + (session.shootTiming ? " (" + session.shootTiming + ")" : ""))
  if (session.duration) {
    messageLines.push("• Duration: " + session.duration)
  }
  messageLines.push("• Venue: " + (session.venue || "To be determined") + (session.city ? ", " + session.city : ""))
  if (session.locationLink) {
    messageLines.push("• Venue Map Link: " + session.locationLink)
  }
  if (session.members && session.members.length > 0) {
    messageLines.push("• Crew / Photographers: " + session.members.join(", "))
  }
  if (session.sessionManager) {
    messageLines.push("• Session Manager / Contact: " + session.sessionManager)
  }
  if (session.remarks) {
    messageLines.push("• Notes: " + session.remarks)
  }
  messageLines.push("")
  messageLines.push("Please ensure you arrive 15 minutes before the scheduled call time.")
  messageLines.push("We look forward to capturing your moments!")
  return messageLines.join("\n")
}

export function formatLeadSchedule(lead: Lead): string {
  var messageLines: string[] = []
  messageLines.push("Photoshoot Schedule Details: " + (lead.projectName || lead.name))
  messageLines.push("")
  messageLines.push("Hello " + (lead.name || "Client") + ",")
  messageLines.push("Here are the schedule and location details for your shoot:")
  messageLines.push("• Shoot Type: " + lead.projectType)
  if (lead.projectStartDate) {
    messageLines.push("• Call Date: " + lead.projectStartDate + (lead.projectEndDate ? " to " + lead.projectEndDate : ""))
  }
  if (lead.projectAddress) {
    messageLines.push("• Venue Location: " + lead.projectAddress)
  }
  if (lead.assignee) {
    messageLines.push("• Assigned Photographer/POC: " + lead.assignee.name)
  }
  if (lead.projectSegments && lead.projectSegments.length > 0) {
    messageLines.push("• Event Segments:")
    var segmentIndex: number = 0
    for (segmentIndex = 0; segmentIndex < lead.projectSegments.length; segmentIndex++) {
      var segment = lead.projectSegments[segmentIndex]
      messageLines.push("   - " + segment.name + " (" + segment.type + "): " + segment.startDate + " to " + segment.endDate)
    }
  }
  messageLines.push("")
  messageLines.push("For any updates or adjustments, please feel free to reach out.")
  return messageLines.join("\n")
}
