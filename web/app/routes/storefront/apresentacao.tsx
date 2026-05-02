import { Link } from "react-router"
import { useEffect, useRef } from "react"

const BENEFITS = [
  {
    icon: "📦",
    title: "Organização",
    description:
      "Você nunca mais esquece de comprar o que consome toda semana. Sua cesta chega no dia e horário que você escolheu.",
  },
  {
    icon: "💰",
    title: "Preço justo",
    description:
      "Pagar antecipado gera desconto real. Quanto mais ciclos, maior a economia — sem pegadinhas.",
  },
  {
    icon: "🌿",
    title: "Frescor garantido",
    description:
      "Alimentos chegam com regularidade, colhidos ou separados próximos à entrega. Sem acúmulo, sem desperdício.",
  },
  {
    icon: "🏘️",
    title: "Impacto local",
    description:
      "Negócios locais dependem menos de empréstimos e juros abusivos. Aqui não tem comissão por venda — o que você paga chega inteiro para quem produz.",
  },
] as const

const HOW_IT_WORKS = [
  { step: "1", label: "Escolha os produtos que quer receber" },
  { step: "2", label: "Defina os dias e o horário de entrega" },
  { step: "3", label: "Escolha quantas entregas quer pagar antecipado" },
  { step: "4", label: "Confirme — e pronto, sua cesta está ativa" },
] as const

function useReveal() {
  const ref = useRef<HTMLDivElement>(null)

  useEffect(() => {
    const el = ref.current
    if (!el) return

    const observer = new IntersectionObserver(
      (entries) => {
        entries.forEach((entry) => {
          if (!entry.isIntersecting) return
          entry.target.classList.add("is-visible")
          observer.unobserve(entry.target)
        })
      },
      { threshold: 0.15 },
    )

    const targets = el.querySelectorAll(".reveal")
    targets.forEach((t) => observer.observe(t))

    return () => observer.disconnect()
  }, [])

  return ref
}

export default function Apresentacao() {
  const ref = useReveal()

  return (
    <div ref={ref} className="min-h-screen bg-background text-foreground">
      <style>{`
        .reveal {
          opacity: 0;
          transform: translateY(24px);
          transition: opacity 0.6s ease, transform 0.6s ease;
        }
        .reveal.is-visible {
          opacity: 1;
          transform: translateY(0);
        }
        .reveal-delay-1 { transition-delay: 0.1s; }
        .reveal-delay-2 { transition-delay: 0.2s; }
        .reveal-delay-3 { transition-delay: 0.3s; }
        .reveal-delay-4 { transition-delay: 0.4s; }
      `}</style>

      {/* Hero */}
      <section className="relative overflow-hidden bg-foreground text-background px-6 py-24 text-center">
        <div className="relative z-10 mx-auto max-w-xl">
          <p className="text-sm font-medium uppercase tracking-widest opacity-60 mb-4">
            Assinatura de alimentos
          </p>
          <h1 className="text-4xl sm:text-5xl font-bold leading-tight mb-6">
            Receba o que você precisa,{" "}
            <span className="opacity-60">sempre.</span>
          </h1>
          <p className="text-lg opacity-70 mb-10 leading-relaxed">
            Monte sua cesta, escolha os dias de entrega e esqueça de ter que lembrar.
            Com desconto por comprometimento — sem fidelidade forçada.
          </p>
          <Link
            to="/"
            className="inline-block rounded-full bg-background text-foreground px-8 py-3 text-sm font-semibold hover:opacity-90 transition-opacity"
          >
            Ver produtos
          </Link>
        </div>

        {/* Decorative blobs */}
        <div
          aria-hidden
          className="absolute -top-20 -right-20 w-72 h-72 rounded-full opacity-10"
          style={{ background: "radial-gradient(circle, white 0%, transparent 70%)" }}
        />
        <div
          aria-hidden
          className="absolute -bottom-16 -left-16 w-56 h-56 rounded-full opacity-10"
          style={{ background: "radial-gradient(circle, white 0%, transparent 70%)" }}
        />
      </section>

      {/* Benefits */}
      <section className="px-6 py-20 mx-auto max-w-2xl">
        <h2 className="reveal text-2xl font-bold text-center mb-12">
          Por que assinar?
        </h2>
        <div className="space-y-8">
          {BENEFITS.map((b, i) => (
            <div
              key={b.title}
              className={`reveal reveal-delay-${i + 1} flex gap-5 items-start`}
            >
              <div className="text-4xl shrink-0 w-12 text-center">{b.icon}</div>
              <div>
                <h3 className="font-semibold text-lg mb-1">{b.title}</h3>
                <p className="text-muted-foreground leading-relaxed">{b.description}</p>
              </div>
            </div>
          ))}
        </div>
      </section>

      {/* How it works */}
      <section className="px-6 py-20 bg-muted/40">
        <div className="mx-auto max-w-lg">
          <h2 className="reveal text-2xl font-bold text-center mb-12">
            Como funciona
          </h2>
          <ol className="space-y-6">
            {HOW_IT_WORKS.map((item, i) => (
              <li
                key={item.step}
                className={`reveal reveal-delay-${i + 1} flex items-start gap-4`}
              >
                <span className="shrink-0 w-8 h-8 rounded-full bg-foreground text-background flex items-center justify-center text-sm font-bold">
                  {item.step}
                </span>
                <p className="pt-1 text-sm leading-relaxed">{item.label}</p>
              </li>
            ))}
          </ol>
        </div>
      </section>

      {/* FAQ inline — cancellation */}
      <section className="px-6 py-20 mx-auto max-w-lg">
        <h2 className="reveal text-2xl font-bold text-center mb-10">
          Perguntas frequentes
        </h2>
        <div className="space-y-8">
          <div className="reveal reveal-delay-1">
            <h3 className="font-semibold mb-1">Posso cancelar a qualquer momento?</h3>
            <p className="text-muted-foreground text-sm leading-relaxed">
              Sim. Se quiser cancelar dentro de 7 dias após a primeira entrega, reembolsamos
              as entregas não realizadas automaticamente. Depois desse prazo, a assinatura
              segue entregando o que já foi pago — sem multa.
            </p>
          </div>
          <div className="reveal reveal-delay-2">
            <h3 className="font-semibold mb-1">E se eu não quiser renovar?</h3>
            <p className="text-muted-foreground text-sm leading-relaxed">
              Nenhum problema. Cancele a renovação quando quiser. As entregas pagas
              continuam chegando até o fim do ciclo.
            </p>
          </div>
          <div className="reveal reveal-delay-3">
            <h3 className="font-semibold mb-1">Posso trocar um produto da cesta?</h3>
            <p className="text-muted-foreground text-sm leading-relaxed">
              Sim — substituição é diferente de cancelamento. Você pode trocar qualquer
              produto da cesta a qualquer momento.
            </p>
          </div>
        </div>
      </section>

      {/* CTA */}
      <section className="px-6 py-20 bg-foreground text-background text-center">
        <div className="mx-auto max-w-sm">
          <h2 className="text-2xl font-bold mb-4">Pronto para começar?</h2>
          <p className="opacity-70 text-sm mb-8 leading-relaxed">
            Escolha um produto na vitrine e monte sua cesta em menos de 2 minutos.
          </p>
          <Link
            to="/"
            className="inline-block rounded-full bg-background text-foreground px-8 py-3 text-sm font-semibold hover:opacity-90 transition-opacity"
          >
            Ver produtos
          </Link>
        </div>
      </section>
    </div>
  )
}
