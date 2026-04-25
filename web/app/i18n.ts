import type { InitOptions } from "i18next"

export const DEFAULT_LOCALE = "pt-BR" as const

export const supportedLngs = ["pt-BR", "en"] as const
export type SupportedLng = (typeof supportedLngs)[number]

export default {
  supportedLngs,
  fallbackLng: DEFAULT_LOCALE,
  defaultNS: "common",
  ns: ["common"],
} satisfies InitOptions
