import { Card, CardContent } from "@/components/ui/card"
import { Button } from "@/components/ui/button"
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

interface DayTimeline {
  date: string
  blocks: TimelineBlock[]
}

interface TimelineProps {
  data: DayTimeline[]
}

export default function Timeline({ data }: TimelineProps) {
  const [selectedDate, setSelectedDate] = useState(data[0]?.date ?? "")
  const [editDialogOpen, setEditDialogOpen] = useState(false)
  const [editingBlockIndex, setEditingBlockIndex] = useState<number | null>(null)
  const [editedBlock, setEditedBlock] = useState<TimelineBlock | null>(null)
  const [dragging, setDragging] = useState(false)
  const [previewBlock, setPreviewBlock] = useState<TimelineBlock | null>(null)
  const containerRef = useRef<HTMLDivElement>(null)

  const selectedDayIndex = data.findIndex((d) => d.date === selectedDate)
  const selectedDay = data[selectedDayIndex]
  const hours = Array.from({ length: 25 }, (_, i) => i)

  const handleBlockClick = (index: number) => {
    setEditingBlockIndex(index)
    setEditedBlock({ ...selectedDay.blocks[index] })
    setEditDialogOpen(true)
  }

  const handleBlockSave = () => {
    if (editingBlockIndex === null || !editedBlock) return
    const updatedData = [...data]
    updatedData[selectedDayIndex].blocks[editingBlockIndex] = editedBlock
    setEditDialogOpen(false)
    setPreviewBlock(null)
  }

  const toHours = (clientX: number, container: DOMRect) => {
    const offsetX = clientX - container.left
    const percentage = offsetX / container.width
    return Math.max(0, Math.min(24, (percentage * 24)))
  }

  const handleDrag = (e: React.MouseEvent, index: number) => {
    if (!containerRef.current) return
    const rect = containerRef.current.getBoundingClientRect()
    const newStart = toHours(e.clientX, rect)
    const block = selectedDay.blocks[index]
    const duration = block.end - block.start
    const newEnd = Math.min(24, newStart + duration)
    const updatedBlock = { ...block, start: newStart, end: newEnd }
    setPreviewBlock(updatedBlock)
  }

  const handleResize = (e: React.MouseEvent, index: number, side: 'left' | 'right') => {
    if (!containerRef.current) return
    const rect = containerRef.current.getBoundingClientRect()
    const time = toHours(e.clientX, rect)
    const block = selectedDay.blocks[index]
    const updatedBlock = { ...block }
    if (side === 'left') {
      updatedBlock.start = Math.min(updatedBlock.end - 1 / 3600, time) // minimalna zmiana: 1 sekunda
    } else {
      updatedBlock.end = Math.max(updatedBlock.start + 1 / 3600, time)
    }
    setPreviewBlock(updatedBlock)
  }

  const applyPreview = (index: number) => {
    if (!previewBlock) return
    const updatedData = [...data]
    updatedData[selectedDayIndex].blocks[index] = previewBlock
    setPreviewBlock(null)
  }

  return (
    <div className="p-4 space-y-4">
      {/* Pasek dni */}
      <div className="flex gap-2 overflow-x-auto">
        {data.map((day) => (
          <Button
            key={day.date}
            variant={day.date === selectedDate ? "default" : "outline"}
            onClick={() => setSelectedDate(day.date)}
          >
            {day.date}
          </Button>
        ))}
      </div>

      {/* Timeline */}
      <Card className="overflow-x-auto">
        <CardContent>
          <div ref={containerRef} className="relative h-24 border rounded bg-muted select-none">
            {/* Osie czasu */}
            <div className="absolute inset-0 flex">
              {hours.map((h) => (
                <div
                  key={h}
                  className="flex-1 border-r text-xs text-center text-muted-foreground"
                >
                  {String(h).padStart(2, '0')}:00
                </div>
              ))}
            </div>

            {/* Bloki czasu */}
            {selectedDay?.blocks.map((block, i) => {
              const isPreview = previewBlock && editingBlockIndex === i
              const effectiveBlock = isPreview ? previewBlock : block
              const effectiveLeft = `${(effectiveBlock.start / 24) * 100}%`
              const effectiveWidth = `${((effectiveBlock.end - effectiveBlock.start) / 24) * 100}%`

              return (
                <div
                  key={i}
                  className={cn(
                    "absolute top-8 h-8 rounded text-white text-sm flex items-center justify-center px-2 cursor-move",
                    block.color ? '' : 'bg-primary',
                    isPreview ? 'opacity-50 pointer-events-none' : ''
                  )}
                  style={{
                    left: effectiveLeft,
                    width: effectiveWidth,
                    backgroundColor: effectiveBlock.color || undefined,
                  }}
                  onMouseDown={(e) => {
                    e.preventDefault()
                    setEditingBlockIndex(i)
                    const move = (ev: MouseEvent) => handleDrag(ev as unknown as React.MouseEvent, i)
                    const up = () => {
                      applyPreview(i)
                      setDragging(false)
                      window.removeEventListener('mousemove', move)
                      window.removeEventListener('mouseup', up)
                    }
                    window.addEventListener('mousemove', move)
                    window.addEventListener('mouseup', up)
                    setDragging(true)
                  }}
                  onClick={(e) => {
                    if (!dragging && !isPreview) handleBlockClick(i)
                  }}
                >
                  <div
                    className="absolute left-0 top-0 h-full w-2 cursor-ew-resize z-10"
                    onMouseDown={(e) => {
                      e.stopPropagation()
                      setEditingBlockIndex(i)
                      const move = (ev: MouseEvent) => handleResize(ev as unknown as React.MouseEvent, i, 'left')
                      const up = () => {
                        applyPreview(i)
                        window.removeEventListener('mousemove', move)
                        window.removeEventListener('mouseup', up)
                      }
                      window.addEventListener('mousemove', move)
                      window.addEventListener('mouseup', up)
                    }}
                  />
                  <div
                    className="absolute right-0 top-0 h-full w-2 cursor-ew-resize z-10"
                    onMouseDown={(e) => {
                      e.stopPropagation()
                      setEditingBlockIndex(i)
                      const move = (ev: MouseEvent) => handleResize(ev as unknown as React.MouseEvent, i, 'right')
                      const up = () => {
                        applyPreview(i)
                        window.removeEventListener('mousemove', move)
                        window.removeEventListener('mouseup', up)
                      }
                      window.addEventListener('mousemove', move)
                      window.addEventListener('mouseup', up)
                    }}
                  />
                  {effectiveBlock.label}
                </div>
              )
            })}
          </div>
        </CardContent>
      </Card>

      {/* Dialog edycji */}
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
              <Button onClick={handleBlockSave}>Zapisz</Button>
            </div>
          )}
        </DialogContent>
      </Dialog>
    </div>
  )
}

function timeToDecimal(time: string): number {
  const [h, m, s] = time.split(":").map(Number)
  return h + m / 60 + s / 3600
}

function secondsToTime(seconds: number): string {
  const h = Math.floor(seconds / 3600)
  const m = Math.floor((seconds % 3600) / 60)
  const s = Math.floor(seconds % 60)
  return [h, m, s].map((n) => String(n).padStart(2, "0")).join(":")
}