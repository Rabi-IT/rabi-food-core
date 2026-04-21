import { resolve } from "node:path"
import { createCookie } from "react-router"
import { RemixI18Next } from "remix-i18next/server"
import FsBackend from "i18next-fs-backend"
import i18nConfig from "~/i18n"

export const localeCookie = createCookie("i18n", {
  path: "/",
  sameSite: "lax",
  maxAge: 60 * 60 * 24 * 365,
})

export default new RemixI18Next({
  detection: {
    supportedLanguages: i18nConfig.supportedLngs as unknown as string[],
    fallbackLanguage: i18nConfig.fallbackLng,
    order: ["cookie", "header"],
    cookie: localeCookie,
  },
  i18next: {
    ...i18nConfig,
    backend: {
      loadPath: resolve("./public/locales/{{lng}}/{{ns}}.json"),
    },
  },
  plugins: [FsBackend],
})
