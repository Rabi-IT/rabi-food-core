import { useState } from "react"
import { Button } from "~/components/ui/button"
import { cn } from "~/lib/utils"
import type { WizardData } from "./use-wizard"

type AddressData = WizardData["address"]

type ExistingProfile = {
  street: string; complement: string; neighborhood: string
  city: string; state: string; zip: string; phone: string
}

type Props = {
  phone: string
  address: AddressData
  existingProfile: ExistingProfile | null
  isLoadingProfile: boolean
  onChangePhone: (phone: string) => void
  onChangeAddress: (address: AddressData) => void
  onNext: () => void
  isSaving: boolean
  saveError?: string
}

function Field({ label, error, children }: { label: string; error?: string; children: React.ReactNode }) {
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

function validate(phone: string, address: AddressData): Record<string, string> {
  const errors: Record<string, string> = {}
  if (!phone.trim()) errors.phone = "Telefone obrigatório"
  if (!address.zip.trim()) errors.zip = "CEP obrigatório"
  if (!address.street.trim()) errors.street = "Rua obrigatória"
  if (!address.neighborhood.trim()) errors.neighborhood = "Bairro obrigatório"
  if (!address.city.trim()) errors.city = "Cidade obrigatória"
  if (!address.state.trim()) errors.state = "Estado obrigatório"
  return errors
}

export function StepAddress({
  phone,
  address,
  existingProfile,
  isLoadingProfile,
  onChangePhone,
  onChangeAddress,
  onNext,
  isSaving,
  saveError,
}: Props) {
  const [touched, setTouched] = useState(false)
  const [mode, setMode] = useState<"confirm" | "edit">("confirm")
  const [cepLoading, setCepLoading] = useState(false)
  const [cepError, setCepError] = useState<string | null>(null)
  const [addressReady, setAddressReady] = useState(
    () => address.street.trim() !== "" || address.city.trim() !== ""
  )

  async function handleZipChange(raw: string) {
    onChangeAddress({ ...address, zip: raw })
    setCepError(null)

    const digits = raw.replace(/\D/g, "")
    if (digits.length !== 8) {
      setAddressReady(false)
      return
    }

    setCepLoading(true)
    try {
      const res = await fetch(`https://viacep.com.br/ws/${digits}/json/`)
      const json = await res.json()
      if (json.erro) {
        setCepError("CEP não encontrado")
        setAddressReady(false)
        return
      }
      onChangeAddress({
        ...address,
        zip: raw,
        street: json.logradouro ?? "",
        neighborhood: json.bairro ?? "",
        city: json.localidade ?? "",
        state: json.uf ?? "",
      })
      setAddressReady(true)
    } catch {
      setCepError("Erro ao buscar CEP. Preencha manualmente.")
      setAddressReady(true)
    } finally {
      setCepLoading(false)
    }
  }

  function handleNext() {
    setTouched(true)
    if (Object.keys(validate(phone, address)).length === 0) onNext()
  }

  const errors = touched ? validate(phone, address) : {}
  const hasExisting = existingProfile && existingProfile.street.trim() !== ""

  if (isLoadingProfile) {
    return (
      <div className="flex items-center justify-center py-16 text-sm text-muted-foreground">
        Carregando…
      </div>
    )
  }

  if (hasExisting && mode === "confirm") {
    const p = existingProfile
    return (
      <div className="flex flex-col pb-32 space-y-4">
        <p className="text-sm text-muted-foreground">
          Encontramos um endereço cadastrado. Deseja usá-lo para esta entrega?
        </p>
        <div className="rounded-xl border bg-card p-4 space-y-1 text-sm">
          {p.phone && <p><span className="text-muted-foreground">Telefone:</span> {p.phone}</p>}
          <p>{p.street}</p>
          {p.complement && <p>{p.complement}</p>}
          <p>{p.neighborhood}</p>
          <p>{p.city} — {p.state}</p>
          <p>{p.zip}</p>
        </div>
        <div className="fixed inset-x-0 bottom-0 border-t bg-background p-4 space-y-2">
          <Button
            className="w-full"
            onClick={() => {
              onChangePhone(p.phone)
              onChangeAddress({ street: p.street, complement: p.complement, neighborhood: p.neighborhood, city: p.city, state: p.state, zip: p.zip })
              onNext()
            }}
            disabled={isSaving}
          >
            {isSaving ? "Salvando…" : "Usar este endereço"}
          </Button>
          <button
            type="button"
            onClick={() => setMode("edit")}
            className="w-full text-sm text-primary hover:underline py-1"
          >
            Atualizar endereço
          </button>
        </div>
      </div>
    )
  }

  return (
    <div className="flex flex-col pb-32 space-y-4">
      <Field label="Telefone" error={errors.phone}>
        <input
          type="tel"
          value={phone}
          onChange={(e) => onChangePhone(e.target.value)}
          className={cn(inputClass, errors.phone && "border-destructive")}
          placeholder="(11) 9 0000-0000"
          autoComplete="tel"
        />
      </Field>

      <Field label="CEP" error={cepError ?? errors.zip}>
        <div className="relative">
          <input
            type="text"
            value={address.zip}
            onChange={(e) => handleZipChange(e.target.value)}
            className={cn(inputClass, (cepError || errors.zip) && "border-destructive")}
            placeholder="00000-000"
            autoComplete="postal-code"
            maxLength={9}
          />
          {cepLoading && (
            <span className="absolute right-3 top-1/2 -translate-y-1/2 text-xs text-muted-foreground animate-pulse">
              Buscando…
            </span>
          )}
        </div>
      </Field>

      {addressReady && (
        <>
          <Field label="Rua e número" error={errors.street}>
            <input
              type="text"
              value={address.street}
              onChange={(e) => onChangeAddress({ ...address, street: e.target.value })}
              className={cn(inputClass, errors.street && "border-destructive")}
              placeholder="Rua Exemplo, 123"
              autoComplete="street-address"
            />
          </Field>

          <Field label="Complemento">
            <input
              type="text"
              value={address.complement}
              onChange={(e) => onChangeAddress({ ...address, complement: e.target.value })}
              className={inputClass}
              placeholder="Apto, bloco… (opcional)"
            />
          </Field>

          <Field label="Bairro" error={errors.neighborhood}>
            <input
              type="text"
              value={address.neighborhood}
              onChange={(e) => onChangeAddress({ ...address, neighborhood: e.target.value })}
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
                  onChange={(e) => onChangeAddress({ ...address, city: e.target.value })}
                  className={cn(inputClass, errors.city && "border-destructive")}
                  placeholder="Cidade"
                />
              </Field>
            </div>
            <Field label="Estado" error={errors.state}>
              <input
                type="text"
                value={address.state}
                onChange={(e) => onChangeAddress({ ...address, state: e.target.value })}
                className={cn(inputClass, errors.state && "border-destructive")}
                placeholder="SP"
                maxLength={2}
              />
            </Field>
          </div>
        </>
      )}

      {saveError && (
        <p className="text-sm text-destructive rounded-lg border border-destructive/30 bg-destructive/5 p-3">
          {saveError}
        </p>
      )}

      <div className="fixed inset-x-0 bottom-0 border-t bg-background p-4">
        <Button className="w-full" onClick={handleNext} disabled={isSaving || cepLoading}>
          {isSaving ? "Salvando…" : "Continuar"}
        </Button>
      </div>
    </div>
  )
}
