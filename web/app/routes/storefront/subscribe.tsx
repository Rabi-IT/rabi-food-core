import { useLoaderData, useActionData, useNavigation, useSubmit, Link } from "react-router"
import type { Route } from "./+types/subscribe"
import { resolveTenant } from "~/lib/tenant.server"
import { apiClient } from "~/lib/api.server"
import { useWizard, type WizardData } from "~/components/subscription/use-wizard"
import { StepBasket } from "~/components/subscription/step-basket"
import { StepDelivery } from "~/components/subscription/step-delivery"
import { StepPlan } from "~/components/subscription/step-plan"
import { StepUser } from "~/components/subscription/step-user"
import { StepSummary } from "~/components/subscription/step-summary"
import { useEffect } from "react"
import { cn } from "~/lib/utils"
import type { Product, SubscriptionConfig } from "~/components/subscription/types"


type ActionResult = { ok: true } | { error: string }

export async function loader({ request }: Route.LoaderArgs) {
  const tenant = await resolveTenant(request)
  if (!tenant) throw new Response("Not found", { status: 404 })

  const api = apiClient(tenant.id)

  const [productsRes, config] = await Promise.all([
    api.get<{ data: readonly Product[] }>("/product/", { isActive: "1" }),
    api.get<SubscriptionConfig>("/subscription/config").catch(() => null),
  ])

  if (!config) throw new Response("Not found", { status: 404 })

  return { products: productsRes.data, config }
}

export async function action({ request }: Route.ActionArgs): Promise<ActionResult> {
  const tenant = await resolveTenant(request)
  if (!tenant) return { error: "Tenant não encontrado" }

  const api = apiClient(tenant.id)
  const data = (await request.json()) as WizardData

  try {
    await api.post("/auth/signup", {
      email: data.user.email,
      password: data.user.password,
      name: data.user.name,
      phone: data.user.phone,
      taxId: data.user.taxId,
      street: data.address.street,
      complement: data.address.complement,
      neighborhood: data.address.neighborhood,
      city: data.address.city,
      state: data.address.state,
      zip: data.address.zip,
    })
  } catch {
    return { error: "Não foi possível criar sua conta. Verifique seus dados e tente novamente." }
  }

  let token: string
  try {
    const auth = await api.post<{ accessToken: string }>("/auth/signin", {
      email: data.user.email,
      password: data.user.password,
    })
    token = auth.accessToken
  } catch {
    return { error: "Conta criada, mas não foi possível autenticar. Tente fazer login." }
  }

  try {
    await api.post(
      "/subscription/",
      {
        items: data.items,
        deliveryDays: data.deliveryDays,
        totalCycles: data.totalCycles,
        autoRenew: data.autoRenew,
      },
      token,
    )
  } catch {
    return { error: "Conta criada! Mas houve um problema ao salvar sua assinatura. Faça login para concluir." }
  }

  return { ok: true }
}

const STEP_LABELS = {
  1: "Monte sua cesta",
  2: "Dias de entrega",
  3: "Seu plano",
  4: "Seus dados",
  summary: "Resumo",
} as const

const NUMERIC_STEPS = [1, 2, 3, 4] as const

export default function Subscribe() {
  const { products, config } = useLoaderData<typeof loader>()
  const actionData = useActionData<typeof action>()
  const navigation = useNavigation()
  const submit = useSubmit()
  const { step, setStep, data, setData, reset } = useWizard()

  const isSubmitting = navigation.state !== "idle"
  const actionResult = actionData as ActionResult | undefined

  useEffect(() => {
    if (actionResult && "ok" in actionResult) reset()
  }, [actionResult]) // eslint-disable-line react-hooks/exhaustive-deps

  function handleConfirm() {
    submit(data as unknown as Record<string, string>, {
      method: "post",
      encType: "application/json",
    })
  }

  if (actionResult && "ok" in actionResult) {
    return (
      <div className="min-h-screen bg-background flex flex-col items-center justify-center px-4 text-center">
        <div className="text-5xl mb-4">🎉</div>
        <h1 className="text-2xl font-bold mb-2">Assinatura confirmada!</h1>
        <p className="text-muted-foreground mb-8">
          Sua cesta recorrente está ativa. Em breve você receberá um e-mail de confirmação.
        </p>
        <Link
          to="/"
          className="text-sm text-primary font-medium hover:underline"
        >
          Voltar à loja
        </Link>
      </div>
    )
  }

  const currentStepNumber = step === "summary" ? null : (step as number)
  const stepLabel = STEP_LABELS[step]
      

  return (
    <div className="min-h-screen bg-background text-foreground">
      {/* Header */}
      <header className="sticky top-0 z-10 bg-background border-b px-4 py-3">
        <div className="flex items-center gap-3">
          {step !== 1 ? (
            <button
              type="button"
              onClick={() => {
                if (step === "summary") setStep(4)
                else if (step === 2) setStep(1)
                else if (step === 3) setStep(2)
                else if (step === 4) setStep(3)
              }}
              className="text-sm text-muted-foreground hover:text-foreground transition-colors"
            >
              ←
            </button>
          ) : (
            <Link to="/" className="text-sm text-muted-foreground hover:text-foreground transition-colors">
              ←
            </Link>
          )}

          <div className="flex-1">
            <p className="text-xs text-muted-foreground">
              {currentStepNumber ? `Etapa ${currentStepNumber} de 4` : "Resumo"}
            </p>
            <p className="text-sm font-medium">{stepLabel}</p>
          </div>
        </div>

        {/* Progress bar */}
        {currentStepNumber && (
          <div className="mt-2 flex gap-1">
            {NUMERIC_STEPS.map((s) => (
              <div
                key={s}
                className={cn(
                  "h-1 flex-1 rounded-full transition-colors",
                  s <= currentStepNumber ? "bg-primary" : "bg-muted"
                )}
              />
            ))}
          </div>
        )}
      </header>

      <main className="mx-auto max-w-lg px-4 py-6">
        {!config?.isOpen && (
          <div className="rounded-xl border border-destructive/30 bg-destructive/5 p-4 mb-6 text-center">
            <p className="text-sm text-destructive">
              As assinaturas estão temporariamente indisponíveis. Tente novamente mais tarde.
            </p>
          </div>
        )}

        {step === 1 && (
          <StepBasket
            products={products}
            items={data.items}
            totalCycles={data.totalCycles}
            cycleDiscountRules={config.discountRules}
            onChange={(items) => setData({ items })}
            onNext={() => setStep(2)}
          />
        )}

        {step === 2 && (
          <StepDelivery
            deliveryDays={data.deliveryDays}
            onChange={(deliveryDays) => setData({ deliveryDays })}
            onNext={() => setStep(3)}
          />
        )}

        {step === 3 && (
          <StepPlan
            products={products}
            items={data.items}
            totalCycles={data.totalCycles}
            cycleDiscountRules={config?.discountRules ?? []}
            autoRenew={data.autoRenew}
            onChange={(totalCycles, autoRenew) => setData({ totalCycles, autoRenew })}
            onNext={() => setStep(4)}
          />
        )}

        {step === 4 && (
          <StepUser
            user={data.user}
            address={data.address}
            onChange={(user, address) => setData({ user, address })}
            onNext={() => setStep("summary")}
          />
        )}

        {step === "summary" && (
          <StepSummary
            data={data}
            products={products}
            cycleDiscountRules={config.discountRules}
            cutoffOffsetMinutes={config?.cutoffOffsetMinutes ?? 0}
            error={actionResult && "error" in actionResult ? actionResult.error : undefined}
            isSubmitting={isSubmitting}
            onConfirm={handleConfirm}
            onBack={() => setStep(4)}
          />
        )}
      </main>
    </div>
  )
}
