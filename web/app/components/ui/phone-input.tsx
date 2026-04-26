import BasePhoneInput from "react-phone-number-input/input"
import { cn } from "~/lib/utils"

type Props = {
  readonly value: string
  readonly onChange: (value: string) => void
  readonly className?: string
  readonly id?: string
  readonly autoComplete?: string
}

const PHONE_PLACEHOLDER = "(11) 99999-9999"
const PHONE_MAX_LENGTH = 15 // (11) 99999-9999

export function PhoneInput({ value, onChange, className, id, autoComplete }: Props) {
  return (
    <BasePhoneInput
      id={id}
      country="BR"
      value={value}
      onChange={(v) => onChange(v ?? "")}
      placeholder={PHONE_PLACEHOLDER}
      maxLength={PHONE_MAX_LENGTH}
      className={cn(className)}
      autoComplete={autoComplete}
    />
  )
}
