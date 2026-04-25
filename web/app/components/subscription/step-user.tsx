import { useState } from "react"
import { Button } from "~/components/ui/button"
import { cn } from "~/lib/utils"
import type { WizardData } from "./use-wizard"

type UserData = WizardData["user"]
type AddressData = WizardData["address"]

type Props = {
  user: UserData
  address: AddressData
  onChange: (user: UserData, address: AddressData) => void
  onNext: () => void
}

function Field({
  label,
  error,
  children,
}: {
  label: string
  error?: string
  children: React.ReactNode
}) {
  return (
    <div className="space-y-1">
      <label className="text-sm font-medium">{label}</label>
      {children}
      {error && <p className="text-xs text-destructive">{error}</p>}
    </div>
  )
}

const inputClass =
  "w-full rounded-lg border border-border bg-background px-3 py-2 text-sm outline-none focus:border-primary focus:ring-2 focus:ring-primary/20 transition-colors"

export function StepUser({ user, address, onChange, onNext }: Props) {
  const [touched, setTouched] = useState(false)

  function updateUser(field: keyof UserData, value: string) {
    onChange({ ...user, [field]: value }, address)
  }

  function updateAddress(field: keyof AddressData, value: string) {
    onChange(user, { ...address, [field]: value })
  }

  function validate(): Record<string, string> {
    const errors: Record<string, string> = {}
    if (!user.name.trim()) errors.name = "Nome obrigatório"
    if (!user.phone.trim()) errors.phone = "Telefone obrigatório"
    if (!user.email.trim()) errors.email = "E-mail obrigatório"
    if (user.password.length < 8) errors.password = "Mínimo 8 caracteres"
    if (!user.taxId.trim()) errors.taxId = "CPF obrigatório"
    if (!address.street.trim()) errors.street = "Rua obrigatória"
    if (!address.neighborhood.trim()) errors.neighborhood = "Bairro obrigatório"
    if (!address.city.trim()) errors.city = "Cidade obrigatória"
    if (!address.state.trim()) errors.state = "Estado obrigatório"
    if (!address.zip.trim()) errors.zip = "CEP obrigatório"
    return errors
  }

  function handleNext() {
    setTouched(true)
    const errors = validate()
    if (Object.keys(errors).length === 0) onNext()
  }

  const errors = touched ? validate() : {}

  return (
    <div className="flex flex-col pb-32">
      <section className="space-y-4">
        <h3 className="text-sm font-semibold uppercase tracking-wide text-muted-foreground">
          Dados pessoais
        </h3>

        <Field label="Nome completo" error={errors.name}>
          <input
            type="text"
            value={user.name}
            onChange={(e) => updateUser("name", e.target.value)}
            className={cn(inputClass, errors.name && "border-destructive")}
            placeholder="Seu nome"
          />
        </Field>

        <Field label="Telefone" error={errors.phone}>
          <input
            type="tel"
            value={user.phone}
            onChange={(e) => updateUser("phone", e.target.value)}
            className={cn(inputClass, errors.phone && "border-destructive")}
            placeholder="(11) 9 0000-0000"
          />
        </Field>

        <Field label="E-mail" error={errors.email}>
          <input
            type="email"
            value={user.email}
            onChange={(e) => updateUser("email", e.target.value)}
            className={cn(inputClass, errors.email && "border-destructive")}
            placeholder="seu@email.com"
          />
        </Field>

        <Field label="Senha" error={errors.password}>
          <input
            type="password"
            value={user.password}
            onChange={(e) => updateUser("password", e.target.value)}
            className={cn(inputClass, errors.password && "border-destructive")}
            placeholder="Mínimo 8 caracteres"
          />
        </Field>

        <Field label="CPF" error={errors.taxId}>
          <input
            type="text"
            value={user.taxId}
            onChange={(e) => updateUser("taxId", e.target.value)}
            className={cn(inputClass, errors.taxId && "border-destructive")}
            placeholder="000.000.000-00"
          />
        </Field>
      </section>

      <section className="mt-8 space-y-4">
        <h3 className="text-sm font-semibold uppercase tracking-wide text-muted-foreground">
          Endereço de entrega
        </h3>

        <Field label="CEP" error={errors.zip}>
          <input
            type="text"
            value={address.zip}
            onChange={(e) => updateAddress("zip", e.target.value)}
            className={cn(inputClass, errors.zip && "border-destructive")}
            placeholder="00000-000"
          />
        </Field>

        <Field label="Rua e número" error={errors.street}>
          <input
            type="text"
            value={address.street}
            onChange={(e) => updateAddress("street", e.target.value)}
            className={cn(inputClass, errors.street && "border-destructive")}
            placeholder="Rua Exemplo, 123"
          />
        </Field>

        <Field label="Complemento">
          <input
            type="text"
            value={address.complement}
            onChange={(e) => updateAddress("complement", e.target.value)}
            className={inputClass}
            placeholder="Apto, bloco… (opcional)"
          />
        </Field>

        <Field label="Bairro" error={errors.neighborhood}>
          <input
            type="text"
            value={address.neighborhood}
            onChange={(e) => updateAddress("neighborhood", e.target.value)}
            className={cn(inputClass, errors.neighborhood && "border-destructive")}
            placeholder="Seu bairro"
          />
        </Field>

        <div className="grid grid-cols-3 gap-3">
          <div className="col-span-2">
            <Field label="Cidade" error={errors.city}>
              <input
                type="text"
                value={address.city}
                onChange={(e) => updateAddress("city", e.target.value)}
                className={cn(inputClass, errors.city && "border-destructive")}
                placeholder="Cidade"
              />
            </Field>
          </div>
          <Field label="Estado" error={errors.state}>
            <input
              type="text"
              value={address.state}
              onChange={(e) => updateAddress("state", e.target.value)}
              className={cn(inputClass, errors.state && "border-destructive")}
              placeholder="SP"
              maxLength={2}
            />
          </Field>
        </div>
      </section>

      <div className="fixed inset-x-0 bottom-0 border-t bg-background p-4">
        <Button className="w-full" onClick={handleNext}>
          Revisar assinatura
        </Button>
      </div>
    </div>
  )
}
