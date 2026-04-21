import type { InitOptions } from "i18next"

export const supportedLngs = ["pt-BR", "en"] as const
export type SupportedLng = (typeof supportedLngs)[number]

export default {
  supportedLngs,
  fallbackLng: "pt-BR",
  defaultNS: "common",
  ns: ["common"],
} satisfies InitOptions
