import { useState } from "react"
import { useTranslation } from "react-i18next"
import { isValidPhoneNumber } from "react-phone-number-input"
import { Button } from "~/components/ui/button"
import { PhoneInput } from "~/components/ui/phone-input"
import { cn } from "~/lib/utils"
import type { WizardData } from "./use-wizard"

type AddressData = WizardData["address"]

type ExistingProfile = {
  readonly street: string
  readonly complement: string
  readonly neighborhood: string
  readonly city: string
  readonly state: string
  readonly zip: string
  readonly phone: string
}

type Props = {
  readonly phone: string
  readonly address: AddressData
  readonly existingProfile: ExistingProfile | null
  readonly isLoadingProfile: boolean
  readonly onChangePhone: (phone: string) => void
  readonly onChangeAddress: (address: AddressData) => void
  readonly onNext: () => void
  readonly onChangeIdentity?: () => void
  readonly isSaving: boolean
  readonly saveError?: string
}

type AddressMode = "confirm" | "edit"

const ZIP_LENGTH = 8

const inputClass =
  "w-full rounded-lg border border-border bg-background px-3 py-2 text-sm outline-none focus:border-primary focus:ring-2 focus:ring-primary/20 transition-colors"

function Field({
  label,
  error,
  children,
}: {
  readonly label: string
  readonly error?: string
  readonly children: React.ReactNode
}) {
  return (
    <div className="space-y-1">
      <label className="text-sm font-medium">{label}</label>
      {children}
      {error && <p className="text-xs text-destructive">{error}</p>}
    </div>
  )
}

function validate(
  phone: string,
  address: AddressData,
  t: (key: string) => string,
): Readonly<Record<string, string>> {
  const errors: Record<string, string> = {}
  if (!phone || !isValidPhoneNumber(phone)) errors.phone = t("subscription.address.phoneInvalid")
  if (!address.zip.trim()) errors.zip = t("subscription.address.zipRequired")
  if (!address.street.trim()) errors.street = t("subscription.address.streetRequired")
  if (!address.neighborhood.trim()) errors.neighborhood = t("subscription.address.neighborhoodRequired")
  if (!address.city.trim()) errors.city = t("subscription.address.cityRequired")
  if (!address.state.trim()) errors.state = t("subscription.address.stateRequired")
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
  onChangeIdentity,
  isSaving,
  saveError,
}: Props) {
  const { t } = useTranslation()
  const [touched, setTouched] = useState(false)
  const [mode, setMode] = useState<AddressMode>("confirm")
  const [cepLoading, setCepLoading] = useState(false)
  const [cepError, setCepError] = useState<string | null>(null)
  const [addressReady, setAddressReady] = useState(
    () => address.street.trim() !== "" || address.city.trim() !== "",
  )

  async function handleZipChange(raw: string) {
    onChangeAddress({ ...address, zip: raw })
    setCepError(null)

    const digits = raw.replace(/\D/g, "")
    if (digits.length !== ZIP_LENGTH) {
      setAddressReady(false)
      return
    }

    setCepLoading(true)
    try {
      const res = await fetch(`https://viacep.com.br/ws/${digits}/json/`)
      const json = await res.json()
      if (json.erro) {
        setCepError(t("subscription.address.zipNotFound"))
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
      setCepError(t("subscription.address.zipError"))
      setAddressReady(true)
    } finally {
      setCepLoading(false)
    }
  }

  function handleNext() {
    setTouched(true)
    if (Object.keys(validate(phone, address, t)).length === 0) onNext()
  }

  function handleUseExisting(p: ExistingProfile) {
    onChangePhone(p.phone)
    onChangeAddress({
      street: p.street,
      complement: p.complement,
      neighborhood: p.neighborhood,
      city: p.city,
      state: p.state,
      zip: p.zip,
    })
    onNext()
  }

  const errors = touched ? validate(phone, address, t) : {}
  const hasExisting = !!existingProfile && existingProfile.street.trim() !== ""

  if (isLoadingProfile) {
    return (
      <div className="flex items-center justify-center py-16 text-sm text-muted-foreground">
        Carregando…
      </div>
    )
  }

  if (hasExisting && mode === "confirm") {
    const p = existingProfile!
    return (
      <div className="flex flex-col space-y-4">
        <p className="text-sm text-muted-foreground">{t("subscription.address.existingHint")}</p>
        <div className="rounded-xl border bg-card p-4 space-y-1 text-sm">
          {p.phone && (
            <p>
              <span className="text-muted-foreground">{t("subscription.address.existingPhone")}:</span>{" "}
              {p.phone}
            </p>
          )}
          <p>{p.street}</p>
          {p.complement && <p>{p.complement}</p>}
          <p>{p.neighborhood}</p>
          <p>{p.city} — {p.state}</p>
          <p>{p.zip}</p>
        </div>
        <div className="sticky bottom-0 border-t bg-background p-4 space-y-2">
          <Button
            className="w-full"
            onClick={() => handleUseExisting(p)}
            disabled={isSaving}
          >
            {isSaving ? t("subscription.address.saving") : t("subscription.address.useExisting")}
          </Button>
          <button
            type="button"
            onClick={() => setMode("edit")}
            className="w-full text-sm text-primary hover:underline py-1"
          >
            {t("subscription.address.updateAddress")}
          </button>
          {onChangeIdentity && (
            <button
              type="button"
              onClick={onChangeIdentity}
              className="w-full text-xs text-muted-foreground hover:text-foreground transition-colors py-1"
            >
              {t("subscription.address.notYou")}
            </button>
          )}
        </div>
      </div>
    )
  }

  return (
    <div className="flex flex-col space-y-4">
      <Field label={t("subscription.address.phone")} error={errors.phone}>
        <PhoneInput
          value={phone}
          onChange={onChangePhone}
          className={cn(inputClass, errors.phone && "border-destructive")}
          autoComplete="tel"
        />
      </Field>

      <Field label={t("subscription.address.zip")} error={cepError ?? errors.zip}>
        <div className="relative">
          <input
            type="text"
            value={address.zip}
            onChange={(e) => handleZipChange(e.target.value)}
            className={cn(inputClass, (cepError ?? errors.zip) && "border-destructive")}
            placeholder={t("subscription.address.zipPlaceholder")}
            autoComplete="postal-code"
            maxLength={9}
          />
          {cepLoading && (
            <span className="absolute right-3 top-1/2 -translate-y-1/2 text-xs text-muted-foreground animate-pulse">
              {t("subscription.address.zipSearching")}
            </span>
          )}
        </div>
      </Field>

      {addressReady && (
        <>
          <Field label={t("subscription.address.street")} error={errors.street}>
            <input
              type="text"
              value={address.street}
              onChange={(e) => onChangeAddress({ ...address, street: e.target.value })}
              className={cn(inputClass, errors.street && "border-destructive")}
              placeholder={t("subscription.address.streetPlaceholder")}
              autoComplete="street-address"
            />
          </Field>

          <Field label={t("subscription.address.complement")}>
            <input
              type="text"
              value={address.complement}
              onChange={(e) => onChangeAddress({ ...address, complement: e.target.value })}
              className={inputClass}
              placeholder={t("subscription.address.complementPlaceholder")}
            />
          </Field>

          <Field label={t("subscription.address.neighborhood")} error={errors.neighborhood}>
            <input
              type="text"
              value={address.neighborhood}
              onChange={(e) => onChangeAddress({ ...address, neighborhood: e.target.value })}
              className={cn(inputClass, errors.neighborhood && "border-destructive")}
              placeholder={t("subscription.address.neighborhoodPlaceholder")}
            />
          </Field>

          <div className="grid grid-cols-3 gap-3">
            <div className="col-span-2">
              <Field label={t("subscription.address.city")} error={errors.city}>
                <input
                  type="text"
                  value={address.city}
                  onChange={(e) => onChangeAddress({ ...address, city: e.target.value })}
                  className={cn(inputClass, errors.city && "border-destructive")}
                  placeholder={t("subscription.address.cityPlaceholder")}
                />
              </Field>
            </div>
            <Field label={t("subscription.address.state")} error={errors.state}>
              <input
                type="text"
                value={address.state}
                onChange={(e) => onChangeAddress({ ...address, state: e.target.value })}
                className={cn(inputClass, errors.state && "border-destructive")}
                placeholder={t("subscription.address.statePlaceholder")}
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

      <div className="sticky bottom-0 border-t bg-background p-4 space-y-2">
        <Button className="w-full" onClick={handleNext} disabled={isSaving || cepLoading}>
          {isSaving ? t("subscription.address.saving") : t("subscription.address.continue")}
        </Button>
        {onChangeIdentity && (
          <button
            type="button"
            onClick={onChangeIdentity}
            className="w-full text-xs text-muted-foreground hover:text-foreground transition-colors py-1"
          >
            {t("subscription.address.notYou")}
          </button>
        )}
      </div>
    </div>
  )
}
