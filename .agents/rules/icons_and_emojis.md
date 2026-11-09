---
trigger: always_on
---

# Icon and Emoji Standards

## 1. Strict Ban on Emojis in UI
- **Never use emojis or emoji icons anywhere** in the user interface, labels, buttons, badges, selectors, headings, or table cells (e.g., no 📁, 📅, 🎯, 💳, 🛍️, ⚠️, ⏰, 🇮🇳, 🇺🇸, 🇬🇧, 🇦🇪, 📥, ✓, ★, etc.).
- Do not use emoji characters as quick decorations, status indicators, or action shortcuts.

## 2. Use Proper Vector Icons
- Always use proper vector/SVG icon components from the project's established icon library (e.g., `@tabler/icons-react` or `lucide-react`).
- Choose semantic, consistent icons matching the purpose (e.g., `<IconFolder />` for folders/projects, `<IconCalendar />` for sessions/dates, `<IconTarget />` for leads, `<IconCreditCard />` for billings, `<IconShoppingBag />` for orders, `<IconAlertTriangle />` for warnings, `<IconStarFilled />` for ratings, `<IconCheck />` for checkmarks, `<IconInbox />` for empty state).
- When representing country codes in phone or address inputs, use clean text abbreviations (e.g. `IN (+91)`, `US (+1)`, `UK (+44)`) rather than flag emojis.

