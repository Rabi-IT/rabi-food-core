import i18next from "i18next"
import HttpBackend from "i18next-http-backend"
import { StrictMode, startTransition } from "react"
import { hydrateRoot } from "react-dom/client"
import { I18nextProvider, initReactI18next } from "react-i18next"
import { HydratedRouter } from "react-router/dom"
import { getInitialNamespaces } from "remix-i18next/client"
import i18nConfig from "~/i18n"

async function main() {
  await i18next
    .use(initReactI18next)
    .use(HttpBackend)
    .init({
      ...i18nConfig,
      ns: getInitialNamespaces(),
      backend: { loadPath: "/locales/{{lng}}/{{ns}}.json" },
      detection: { order: ["htmlTag"] },
    })

  startTransition(() => {
    hydrateRoot(
      document,
      <StrictMode>
        <I18nextProvider i18n={i18next}>
          <HydratedRouter />
        </I18nextProvider>
      </StrictMode>
    )
  })
}

main()
