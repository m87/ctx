import { Card, CardContent } from "@/components/ui/card"
import { Dialog, DialogContent, DialogHeader, DialogTitle } from "@/components/ui/dialog"
import { Input } from "@/components/ui/input"
import { Label } from "@/components/ui/label"
import { cn } from "@/lib/utils"
import { useState, useRef } from "react"

interface TimelineBlock {
  label: string
  start: number // czas w godzinach z minutami i sekundami (np. 13.25 = 13:15)
  end: number
  color?: string // opcjonalnie kolor bloku
}

interface TimelineProps {
  date: string
  blocks: TimelineBlock[]
}

export default function Timeline({ date, blocks }: TimelineProps) {
  const [editDialogOpen, setEditDialogOpen] = useState(false)
  const [editingBlockIndex, setEditingBlockIndex] = useState<number | null>(null)
  const [editedBlock, setEditedBlock] = useState<TimelineBlock | null>(null)
  const [previewBlock, setPreviewBlock] = useState<TimelineBlock | null>(null)
  const containerRef = useRef<HTMLDivElement>(null)

  const toHours = (clientX: number, container: DOMRect) => {
    const offsetX = clientX - container.left
    const percentage = offsetX / container.width
    return Math.max(0, Math.min(24, percentage * 24))
  }

  const handleBlockClick = (index: number) => {
    setEditingBlockIndex(index)
    setEditedBlock({ ...blocks[index] })
    setEditDialogOpen(true)
  }

  const handleBlockSave = () => {
    if (editingBlockIndex === null || !editedBlock) return
    blocks[editingBlockIndex] = editedBlock
    setEditDialogOpen(false)
    setPreviewBlock(null)
  }

  const handleDrag = (e: React.MouseEvent, index: number) => {
    if (!containerRef.current) return
    const rect = containerRef.current.getBoundingClientRect()
    const newStart = toHours(e.clientX, rect)
    const block = blocks[index]
    const duration = block.end - block.start
    const newEnd = Math.min(24, newStart + duration)
    setPreviewBlock({ ...block, start: newStart, end: newEnd })
  }

  const handleResize = (e: React.MouseEvent, index: number, side: "left" | "right") => {
    if (!containerRef.current) return
    const rect = containerRef.current.getBoundingClientRect()
    const time = toHours(e.clientX, rect)
    const block = blocks[index]
    const updatedBlock = { ...block }
    if (side === "left") {
      updatedBlock.start = Math.min(updatedBlock.end - 1 / 3600, time)
    } else {
      updatedBlock.end = Math.max(updatedBlock.start + 1 / 3600, time)
    }
    setPreviewBlock(updatedBlock)
  }

  const applyPreview = (index: number) => {
    if (!previewBlock) return
    blocks[index] = previewBlock
    setPreviewBlock(null)
  }

  return (
      <div className="">
            {/* <div className="mb-2 font-semibold text-sm">{date}</div> */}
            <div ref={containerRef} className="relative h-24 border rounded bg-muted select-none">
              <div className="absolute inset-0 flex">
                {Array.from({ length: 25 }, (_, h) => (
                    <div
                        key={h}
                          className="flex-1 border-r text-xs text-center text-muted-foreground"
                    >
                    </div>
                    
                ))}
              </div>
              {blocks.map((block, index) => {
                const isEditing = editingBlockIndex === index
                const effectiveBlock = isEditing && previewBlock ? previewBlock : block
                const left = `${(effectiveBlock.start / 24) * 100}%`
                const width = `${((effectiveBlock.end - effectiveBlock.start) / 24) * 100}%`

                return (
                    <div
                        key={index}
                        className={cn(
                            "absolute rounded text-white text-sm flex items-center justify-center px-2 cursor-move",
                            block.color ? "" : "bg-primary",
                            isEditing && previewBlock ? "opacity-50 pointer-events-none" : ""
                        )}
                        style={{
                          left,
                          width,
                          backgroundColor: effectiveBlock.color || undefined,
                        }}
                        onMouseDown={(e) => {
                          e.preventDefault()
                          setEditingBlockIndex(index)
                          const move = (ev: MouseEvent) => handleDrag(ev as unknown as React.MouseEvent, index)
                          const up = () => {
                            applyPreview(index)
                            window.removeEventListener("mousemove", move)
                            window.removeEventListener("mouseup", up)
                          }
                          window.addEventListener("mousemove", move)
                          window.addEventListener("mouseup", up)
                        }}
                        onClick={(e) => {
                          if (!previewBlock) handleBlockClick(index)
                        }}
                    >
                      <div
                          className="absolute left-0 top-0 h-full w-2 cursor-ew-resize z-10"
                          onMouseDown={(e) => {
                            e.stopPropagation()
                            setEditingBlockIndex(index)
                            const move = (ev: MouseEvent) => handleResize(ev as unknown as React.MouseEvent, index, "left")
                            const up = () => {
                              applyPreview(index)
                              window.removeEventListener("mousemove", move)
                              window.removeEventListener("mouseup", up)
                            }
                            window.addEventListener("mousemove", move)
                            window.addEventListener("mouseup", up)
                          }}
                      />
                      <div
                          className="absolute right-0 top-0 h-full w-2 cursor-ew-resize z-10"
                          onMouseDown={(e) => {
                            e.stopPropagation()
                            setEditingBlockIndex(index)
                            const move = (ev: MouseEvent) => handleResize(ev as unknown as React.MouseEvent, index, "right")
                            const up = () => {
                              applyPreview(index)
                              window.removeEventListener("mousemove", move)
                              window.removeEventListener("mouseup", up)
                            }
                            window.addEventListener("mousemove", move)
                            window.addEventListener("mouseup", up)
                          }}
                      />
                      {effectiveBlock.label}
                    </div>
                )
              })}
            </div>

        <Dialog open={editDialogOpen} onOpenChange={setEditDialogOpen}>
          <DialogContent>
            <DialogHeader>
              <DialogTitle>Edytuj blok</DialogTitle>
            </DialogHeader>
            {editedBlock && (
                <div className="space-y-4">
                  <div>
                    <Label>Nazwa</Label>
                    <Input
                        value={editedBlock.label}
                        onChange={(e) => setEditedBlock({ ...editedBlock, label: e.target.value })}
                    />
                  </div>
                  <div className="flex gap-4">
                    <div className="flex-1">
                      <Label>Start (hh:mm:ss)</Label>
                      <Input
                          type="time"
                          step="1"
                          value={secondsToTime(editedBlock.start * 3600)}
                          onChange={(e) => setEditedBlock({ ...editedBlock, start: timeToDecimal(e.target.value) })}
                      />
                    </div>
                    <div className="flex-1">
                      <Label>Koniec (hh:mm:ss)</Label>
                      <Input
                          type="time"
                          step="1"
                          value={secondsToTime(editedBlock.end * 3600)}
                          onChange={(e) => setEditedBlock({ ...editedBlock, end: timeToDecimal(e.target.value) })}
                      />
                    </div>
                  </div>
                </div>
            )}
          </DialogContent>
        </Dialog>
      </div>
  )
}

function timeToDecimal(time: string): number {
  const [h, m, s] = time.split(":" ).map(Number)
  return h + m / 60 + s / 3600
}

function secondsToTime(seconds: number): string {
  const h = Math.floor(seconds / 3600)
  const m = Math.floor((seconds % 3600) / 60)
  const s = Math.floor(seconds % 60)
  return [h, m, s].map((n) => String(n).padStart(2, "0")).join(":" )
}