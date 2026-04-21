import { createReadableStreamFromReadable } from "@react-router/node"
import i18next from "i18next"
import FsBackend from "i18next-fs-backend"
import { resolve } from "node:path"
import { PassThrough } from "node:stream"
import { renderToPipeableStream } from "react-dom/server"
import { I18nextProvider, initReactI18next } from "react-i18next"
import type { EntryContext } from "react-router"
import { ServerRouter } from "react-router"
import i18nServer from "~/i18n.server"
import i18nConfig from "~/i18n"

const ABORT_DELAY = 5000

export default async function handleRequest(
  request: Request,
  responseStatusCode: number,
  responseHeaders: Headers,
  routerContext: EntryContext
) {
  const locale = await i18nServer.getLocale(request)
  const ns = i18nServer.getRouteNamespaces(routerContext)

  const instance = i18next.createInstance()
  await instance
    .use(initReactI18next)
    .use(FsBackend)
    .init({
      ...i18nConfig,
      lng: locale,
      ns,
      backend: { loadPath: resolve("./public/locales/{{lng}}/{{ns}}.json") },
    })

  return new Promise<Response>((resolve, reject) => {
    let shellRendered = false

    const { pipe, abort } = renderToPipeableStream(
      <I18nextProvider i18n={instance}>
        <ServerRouter context={routerContext} url={request.url} />
      </I18nextProvider>,
      {
        onShellReady() {
          shellRendered = true
          const body = new PassThrough()
          const stream = createReadableStreamFromReadable(body)
          responseHeaders.set("Content-Type", "text/html")
          resolve(
            new Response(stream, {
              headers: responseHeaders,
              status: responseStatusCode,
            })
          )
          pipe(body)
        },
        onShellError(error: unknown) {
          reject(error)
        },
        onError(error: unknown) {
          responseStatusCode = 500
          if (shellRendered) console.error(error)
        },
      }
    )

    setTimeout(abort, ABORT_DELAY)
  })
}
