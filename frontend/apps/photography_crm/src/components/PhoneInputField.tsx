import {
  Input,
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
  cn,
} from '@repo/ui'
import { COUNTRY_CODES } from '../lib/formatters'

interface PhoneInputFieldProps {
  countryCode: string
  mobileNo: string
  onCountryCodeChange: (code: string) => void
  onMobileNoChange: (mobile: string) => void
  placeholder?: string
  required?: boolean
  className?: string
  disabled?: boolean
}

export function PhoneInputField({
  countryCode,
  mobileNo,
  onCountryCodeChange,
  onMobileNoChange,
  placeholder = '98765 43210',
  required = false,
  className,
  disabled = false,
}: PhoneInputFieldProps) {
  return (
    <div className={cn('flex items-center gap-2', className)}>
      <Select value={countryCode} onValueChange={onCountryCodeChange} disabled={disabled}>
        <SelectTrigger className="w-28 shrink-0 h-10 text-xs font-medium">
          <SelectValue />
        </SelectTrigger>
        <SelectContent>
          {COUNTRY_CODES.map((c) => (
            <SelectItem key={c.code} value={c.code} className="text-xs">
              {c.label}
            </SelectItem>
          ))}
        </SelectContent>
      </Select>
      <Input
        type="tel"
        value={mobileNo}
        onChange={(e) => onMobileNoChange(e.target.value)}
        placeholder={placeholder}
        className="h-10 text-sm flex-1"
        required={required}
        disabled={disabled}
      />
    </div>
  )
}

