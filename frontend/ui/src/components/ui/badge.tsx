import * as React from "react"
import { cva, type VariantProps } from "class-variance-authority"

import { cn } from "@repo/ui/lib/utils"

const badgeVariants = cva(
  "inline-flex items-center rounded-full border px-2.5 py-0.5 text-xs font-semibold transition-colors focus:outline-none focus:ring-2 focus:ring-ring focus:ring-offset-2",
  {
    variants: {
      variant: {
        default:
          "border-transparent bg-primary text-primary-foreground hover:bg-primary/90",
        secondary:
          "border-transparent bg-secondary text-secondary-foreground hover:bg-secondary/80",
        destructive:
          "border-transparent bg-destructive/10 text-destructive border border-destructive/20 hover:bg-destructive/15",
        success:
          "border-transparent bg-emerald-500/10 text-emerald-700 dark:text-emerald-400 border border-emerald-500/20 hover:bg-emerald-500/15",
        warning:
          "border-transparent bg-amber-500/10 text-amber-700 dark:text-amber-400 border border-amber-500/20 hover:bg-amber-500/15",
        info:
          "border-transparent bg-blue-500/10 text-blue-700 dark:text-blue-400 border border-blue-500/20 hover:bg-blue-500/15",
        neutral:
          "border-transparent bg-slate-500/10 text-slate-700 dark:text-slate-400 border border-slate-500/20 hover:bg-slate-500/15",
        outline: "text-foreground border-border",
      },
    },
    defaultVariants: {
      variant: "default",
    },
  }
)

export interface BadgeProps
  extends React.HTMLAttributes<HTMLDivElement>,
    VariantProps<typeof badgeVariants> {}

function Badge({ className, variant, ...props }: BadgeProps) {
  return (
    <div className={cn(badgeVariants({ variant }), className)} {...props} />
  )
}

export { Badge, badgeVariants }
