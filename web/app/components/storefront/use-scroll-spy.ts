import { useEffect, useRef, useState } from "react"

type Item = { readonly id: string; readonly name: string }

const DEAD_ZONE_RATIO = 0.3

export function useScrollSpy(categories: readonly Item[], titleGap: number) {
  const [activeCategoryId, setActiveCategoryId] = useState<string | null>(categories[0]?.id ?? null)
  const [offset, setOffset] = useState(0)
  const barRef = useRef<HTMLDivElement>(null)
  const offsetRef = useRef(0)

  useEffect(() => {
    const bar = barRef.current
    if (!bar) return

    const observer = new ResizeObserver(([entry]) => {
      const next = entry.contentRect.height + titleGap
      offsetRef.current = next
      setOffset(next)
    })

    observer.observe(bar)
    return () => observer.disconnect()
  }, [titleGap])

  useEffect(() => {
    if (categories.length === 0) return

    function handleScroll() {
      const off = offsetRef.current
      let active = categories[0]?.id ?? null

      for (let i = 0; i < categories.length; i++) {
        const el = document.getElementById(`category-${categories[i].id}`)
        if (!el) continue
        if (el.getBoundingClientRect().top > off) continue

        active = categories[i].id

        const next = categories[i + 1]
        if (!next) continue
        const nextEl = document.getElementById(`category-${next.id}`)
        if (!nextEl) continue

        const visibleHeight = nextEl.getBoundingClientRect().top - off
        const sectionHeight = el.getBoundingClientRect().height
        if (visibleHeight < sectionHeight * DEAD_ZONE_RATIO) {
          active = next.id
        }
      }

      setActiveCategoryId(active)
    }

    window.addEventListener("scroll", handleScroll, { passive: true })
    return () => window.removeEventListener("scroll", handleScroll)
  }, [categories])

  useEffect(() => {
    if (!activeCategoryId || !barRef.current) return
    const btn = barRef.current.querySelector(`[data-category="${activeCategoryId}"]`)
    btn?.scrollIntoView({ behavior: "smooth", block: "nearest", inline: "center" })
  }, [activeCategoryId])

  function scrollToCategory(id: string) {
    const el = document.getElementById(`category-${id}`)
    if (!el) return
    const top = el.getBoundingClientRect().top + window.scrollY - offsetRef.current
    window.scrollTo({ top, behavior: "smooth" })
  }

  return { activeCategoryId, offset, barRef, scrollToCategory }
}
