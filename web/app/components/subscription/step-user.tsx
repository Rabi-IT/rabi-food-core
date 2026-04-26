import { useState, useRef } from "react"
import { Button } from "~/components/ui/button"
import { cn } from "~/lib/utils"

type UserData = { name: string; email: string; phone: string; token: string }

type Props = {
  user: UserData
  subStep: "identity" | "code"
  onChange: (user: UserData) => void
  onBack: () => void
  onSendOtp: (name: string, email: string) => void
  onVerifyOtp: (email: string, code: string) => void
  onResetOtp: () => void
  isSending: boolean
  isVerifying: boolean
  otpError?: string
}

const inputClass =
  "w-full rounded-lg border border-border bg-background px-3 py-2 text-sm outline-none focus:border-primary focus:ring-2 focus:ring-primary/20 transition-colors"

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
  const [touched, setTouched] = useState(false)
  const [code, setCode] = useState("")
  const codeInputRef = useRef<HTMLInputElement>(null)

  const nameError = touched && !user.name.trim() ? "Nome obrigatório" : undefined
  const emailError = touched && !user.email.trim() ? "E-mail obrigatório" : undefined

  function handleSend() {
    setTouched(true)
    if (!user.name.trim() || !user.email.trim()) return
    onSendOtp(user.name, user.email)
  }

  function handleVerify() {
    if (code.length < 6) return
    onVerifyOtp(user.email, code)
  }

  return (
    <div className="flex flex-col pb-32">
      <div className="flex items-center gap-3 mb-6">
        <button
          type="button"
          onClick={() => subStep === "code" ? (onResetOtp(), setCode("")) : onBack()}
          className="text-sm text-muted-foreground hover:text-foreground transition-colors"
        >
          ←
        </button>
        <div className="flex gap-1 flex-1">
          <div className="h-1 flex-1 rounded-full bg-primary" />
          <div className={cn("h-1 flex-1 rounded-full transition-colors", subStep === "code" ? "bg-primary" : "bg-muted")} />
        </div>
      </div>

      {subStep === "identity" && (
        <section className="space-y-4">
          <div className="mb-2">
            <p className="text-xs text-muted-foreground">Etapa 1 de 2</p>
            <p className="text-sm font-medium">Quem é você?</p>
          </div>

          <div className="space-y-1">
            <label className="text-sm font-medium">Nome completo</label>
            <input
              type="text"
              value={user.name}
              onChange={(e) => onChange({ ...user, name: e.target.value })}
              className={cn(inputClass, nameError && "border-destructive")}
              placeholder="Seu nome"
              autoComplete="name"
            />
            {nameError && <p className="text-xs text-destructive">{nameError}</p>}
          </div>

          <div className="space-y-1">
            <label className="text-sm font-medium">E-mail</label>
            <input
              type="email"
              value={user.email}
              onChange={(e) => onChange({ ...user, email: e.target.value })}
              className={cn(inputClass, emailError && "border-destructive")}
              placeholder="seu@email.com"
              autoComplete="email"
            />
            {emailError && <p className="text-xs text-destructive">{emailError}</p>}
          </div>

          <p className="text-xs text-muted-foreground">
            Vamos enviar um código para confirmar seu e-mail.
          </p>
        </section>
      )}

      {subStep === "code" && (
        <section className="space-y-6">
          <div className="mb-2">
            <p className="text-xs text-muted-foreground">Etapa 2 de 2</p>
            <p className="text-sm font-medium">Confirme seu e-mail</p>
          </div>

          <p className="text-sm text-muted-foreground">
            Enviamos um código para <span className="font-medium text-foreground">{user.email}</span>.
            Verifique sua caixa de entrada.
          </p>

          <div className="space-y-1">
            <label className="text-sm font-medium">Código de verificação</label>
            <input
              ref={codeInputRef}
              type="text"
              inputMode="numeric"
              pattern="[0-9]*"
              maxLength={6}
              value={code}
              onChange={(e) => setCode(e.target.value.replace(/\D/g, "").slice(0, 6))}
              className={cn(
                inputClass,
                "text-center text-2xl font-bold tracking-widest",
                otpError && "border-destructive"
              )}
              placeholder="000000"
              autoComplete="one-time-code"
              autoFocus
            />
            {otpError && <p className="text-xs text-destructive">{otpError}</p>}
          </div>

          <button
            type="button"
            onClick={() => { onResetOtp(); setCode("") }}
            className="text-xs text-primary hover:underline"
          >
            Reenviar código ou corrigir e-mail
          </button>
        </section>
      )}

      <div className="fixed inset-x-0 bottom-0 border-t bg-background p-4">
        {subStep === "identity" ? (
          <Button className="w-full" onClick={handleSend} disabled={isSending}>
            {isSending ? "Enviando…" : "Enviar código"}
          </Button>
        ) : (
          <Button className="w-full" onClick={handleVerify} disabled={code.length < 6 || isVerifying}>
            {isVerifying ? "Verificando…" : "Confirmar"}
          </Button>
        )}
      </div>
    </div>
  )
}
