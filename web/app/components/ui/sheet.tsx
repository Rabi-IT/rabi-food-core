import { Drawer } from "@base-ui/react"

type SheetProps = {
  open: boolean
  onClose: () => void
  children: React.ReactNode
}

export function Sheet({ open, onClose, children }: SheetProps) {
  return (
    <Drawer.Root open={open} onOpenChange={(o) => !o && onClose()} swipeDirection="down">
      <Drawer.Portal>
        <Drawer.Backdrop className="fixed inset-0 z-40 bg-black/40 transition-opacity data-[starting-style]:opacity-0 data-[ending-style]:opacity-0" />
        <Drawer.Popup className="fixed inset-x-0 bottom-0 z-50 flex max-h-[90vh] flex-col rounded-t-2xl bg-background shadow-2xl outline-none transition-transform data-[starting-style]:translate-y-full data-[ending-style]:translate-y-full">
          <div className="mx-auto mt-3 h-1 w-12 shrink-0 rounded-full bg-muted-foreground/25" />
          <div className="flex-1 overflow-y-auto">
            {children}
          </div>
        </Drawer.Popup>
      </Drawer.Portal>
    </Drawer.Root>
  )
}
