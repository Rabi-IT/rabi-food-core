import { useState, useRef } from "react"
import { useTranslation } from "react-i18next"
import { isValidPhoneNumber } from "react-phone-number-input"
import { Button } from "~/components/ui/button"
import { PhoneInput } from "~/components/ui/phone-input"
import { cn } from "~/lib/utils"

type UserData = {
  readonly name: string
  readonly email: string
  readonly phone: string
  readonly token: string
}

type AccountMode = "new" | "existing"

type Props = {
  readonly user: UserData
  readonly subStep: "identity" | "code"
  readonly onChange: (user: UserData) => void
  readonly onBack: () => void
  readonly onSendOtp: (name: string, email: string, phone: string) => void
  readonly onVerifyOtp: (email: string, code: string) => void
  readonly onResetOtp: () => void
  readonly isSending: boolean
  readonly isVerifying: boolean
  readonly otpError?: string
}

const OTP_CODE_LENGTH = 6
const OTP_TOTAL_STEPS = 2

const inputClass =
  "w-full rounded-lg border border-border bg-background px-3 py-2 text-sm outline-none focus:border-primary focus:ring-2 focus:ring-primary/20 transition-colors"

function validateIdentity(
  user: UserData,
  mode: AccountMode,
): Readonly<Record<string, string>> {
  const errors: Record<string, string> = {}
  if (!user.email.trim()) errors.email = "emailRequired"
  if (mode !== "new") return errors
  if (!user.name.trim()) errors.name = "nameRequired"
  if (!user.phone || !isValidPhoneNumber(user.phone)) errors.phone = "phoneInvalid"
  return errors
}

export function StepUser({
  user,
  subStep,
  onChange,
  onBack,
  onSendOtp,
  onVerifyOtp,
  onResetOtp,
  isSending,
  isVerifying,
  otpError,
}: Props) {
  const { t } = useTranslation()
  const [touched, setTouched] = useState(false)
  const [code, setCode] = useState("")
  const [mode, setMode] = useState<AccountMode | undefined>(undefined)
  const codeInputRef = useRef<HTMLInputElement>(null)

  const errors = touched && mode ? validateIdentity(user, mode) : {}

  function handleBack() {
    if (subStep === "code") {
      onResetOtp()
      setCode("")
      return
    }
    onBack()
  }

  function handleSend() {
    setTouched(true)
    if (!mode) return
    const errs = validateIdentity(user, mode)
    if (Object.keys(errs).length > 0) return
    onSendOtp(
      mode === "new" ? user.name : "",
      user.email,
      mode === "new" ? user.phone : "",
    )
  }

  function handleVerify() {
    if (code.length < OTP_CODE_LENGTH) return
    onVerifyOtp(user.email, code)
  }

  function handleCodeChange(e: React.ChangeEvent<HTMLInputElement>) {
    setCode(e.target.value.replace(/\D/g, "").slice(0, OTP_CODE_LENGTH))
  }

  return (
    <div className="flex flex-col pb-32">
      <div className="flex items-center gap-3 mb-6">
        <button
          type="button"
          onClick={handleBack}
          className="text-sm text-muted-foreground hover:text-foreground transition-colors"
        >
          ←
        </button>
        <div className="flex gap-1 flex-1">
          <div className="h-1 flex-1 rounded-full bg-primary" />
          <div
            className={cn(
              "h-1 flex-1 rounded-full transition-colors",
              subStep === "code" ? "bg-primary" : "bg-muted",
            )}
          />
        </div>
      </div>

      {subStep === "code" && (
        <VerifySection
          email={user.email}
          code={code}
          codeInputRef={codeInputRef}
          otpError={otpError}
          isVerifying={isVerifying}
          onChange={handleCodeChange}
          onVerify={handleVerify}
          onReset={() => { onResetOtp(); setCode("") }}
          t={t}
        />
      )}

      {subStep === "identity" && (
        <section className="space-y-4">
          <div className="mb-2">
            <p className="text-xs text-muted-foreground">
              {t("subscription.user.step", { current: 1, total: OTP_TOTAL_STEPS })}
            </p>
            <p className="text-sm font-medium">{t("subscription.user.title")}</p>
          </div>

          {!mode && (
            <div className="space-y-3">
              <p className="text-sm text-muted-foreground">{t("subscription.user.hasAccount")}</p>
              <div className="grid grid-cols-2 gap-3">
                <button
                  type="button"
                  onClick={() => setMode("existing")}
                  className="rounded-lg border border-border py-3 text-sm font-medium hover:border-primary hover:text-primary transition-colors"
                >
                  {t("subscription.user.hasAccountYes")}
                </button>
                <button
                  type="button"
                  onClick={() => setMode("new")}
                  className="rounded-lg border border-border py-3 text-sm font-medium hover:border-primary hover:text-primary transition-colors"
                >
                  {t("subscription.user.hasAccountNo")}
                </button>
              </div>
            </div>
          )}

          {mode && (
            <IdentityFields
              user={user}
              mode={mode}
              errors={errors}
              onChange={onChange}
              onResetMode={() => { setMode(undefined); setTouched(false) }}
              t={t}
            />
          )}
        </section>
      )}

      <div className="fixed inset-x-0 bottom-0 border-t bg-background p-4">
        {subStep === "identity" && (
          <Button
            className="w-full"
            onClick={handleSend}
            disabled={isSending || !mode}
          >
            {isSending ? t("subscription.user.sending") : t("subscription.user.send")}
          </Button>
        )}
        {subStep === "code" && (
          <Button
            className="w-full"
            onClick={handleVerify}
            disabled={code.length < OTP_CODE_LENGTH || isVerifying}
          >
            {isVerifying
              ? t("subscription.user.verify.confirming")
              : t("subscription.user.verify.confirm")}
          </Button>
        )}
      </div>
    </div>
  )
}

type IdentityFieldsProps = {
  readonly user: UserData
  readonly mode: AccountMode
  readonly errors: Readonly<Record<string, string>>
  readonly onChange: (user: UserData) => void
  readonly onResetMode: () => void
  readonly t: (key: string) => string
}

function IdentityFields({ user, mode, errors, onChange, onResetMode, t }: IdentityFieldsProps) {
  return (
    <>
      {mode === "new" && (
        <div className="space-y-1">
          <label className="text-sm font-medium">{t("subscription.user.name")}</label>
          <input
            type="text"
            value={user.name}
            onChange={(e) => onChange({ ...user, name: e.target.value })}
            className={cn(inputClass, errors.name && "border-destructive")}
            placeholder={t("subscription.user.namePlaceholder")}
            autoComplete="name"
          />
          {errors.name && (
            <p className="text-xs text-destructive">{t(`subscription.user.${errors.name}`)}</p>
          )}
        </div>
      )}

      <div className="space-y-1">
        <label className="text-sm font-medium">{t("subscription.user.email")}</label>
        <input
          type="email"
          value={user.email}
          onChange={(e) => onChange({ ...user, email: e.target.value })}
          className={cn(inputClass, errors.email && "border-destructive")}
          placeholder={t("subscription.user.emailPlaceholder")}
          autoComplete="email"
        />
        {errors.email && (
          <p className="text-xs text-destructive">{t(`subscription.user.${errors.email}`)}</p>
        )}
      </div>

      {mode === "new" && (
        <div className="space-y-1">
          <label className="text-sm font-medium">{t("subscription.user.phone")}</label>
          <PhoneInput
            value={user.phone}
            onChange={(v) => onChange({ ...user, phone: v })}
            className={cn(inputClass, errors.phone && "border-destructive")}
            autoComplete="tel"
          />
          {errors.phone && (
            <p className="text-xs text-destructive">{t(`subscription.user.${errors.phone}`)}</p>
          )}
        </div>
      )}

      <p className="text-xs text-muted-foreground">{t("subscription.user.otpHint")}</p>

      <button
        type="button"
        onClick={onResetMode}
        className="text-xs text-primary hover:underline"
      >
        {t("subscription.user.backToSelection")}
      </button>
    </>
  )
}

type VerifySectionProps = {
  readonly email: string
  readonly code: string
  readonly codeInputRef: React.RefObject<HTMLInputElement | null>
  readonly otpError?: string
  readonly isVerifying: boolean
  readonly onChange: (e: React.ChangeEvent<HTMLInputElement>) => void
  readonly onVerify: () => void
  readonly onReset: () => void
  readonly t: (key: string, opts?: Record<string, unknown>) => string
}

function VerifySection({
  email,
  code,
  codeInputRef,
  otpError,
  onChange,
  onReset,
  t,
}: VerifySectionProps) {
  return (
    <section className="space-y-6">
      <div className="mb-2">
        <p className="text-xs text-muted-foreground">
          {t("subscription.user.step", { current: 2, total: OTP_TOTAL_STEPS })}
        </p>
        <p className="text-sm font-medium">{t("subscription.user.verify.title")}</p>
      </div>

      <p className="text-sm text-muted-foreground">
        {t("subscription.user.verify.hint", { email })
          .split("<1>")
          .flatMap((part, i) => {
            if (i === 0) return [part]
            const [bold, rest] = part.split("</1>")
            return [<span key={i} className="font-medium text-foreground">{bold}</span>, rest]
          })}
      </p>

      <div className="space-y-1">
        <label className="text-sm font-medium">{t("subscription.user.verify.codeLabel")}</label>
        <input
          ref={codeInputRef}
          type="text"
          inputMode="numeric"
          pattern="[0-9]*"
          maxLength={OTP_CODE_LENGTH}
          value={code}
          onChange={onChange}
          className={cn(
            inputClass,
            "text-center text-2xl font-bold tracking-widest",
            otpError && "border-destructive",
          )}
          placeholder={t("subscription.user.verify.codePlaceholder")}
          autoComplete="one-time-code"
          autoFocus
        />
        {otpError && <p className="text-xs text-destructive">{otpError}</p>}
      </div>

      <button
        type="button"
        onClick={onReset}
        className="text-xs text-primary hover:underline"
      >
        {t("subscription.user.verify.resend")}
      </button>
    </section>
  )
}
