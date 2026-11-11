import * as React from "react"
import { useEditor, EditorContent } from "@tiptap/react"
import StarterKit from "@tiptap/starter-kit"
import {
  IconBold,
  IconItalic,
  IconStrikethrough,
  IconCode,
  IconList,
  IconListNumbers,
  IconH1,
  IconH2,
  IconBlockquote,
  IconArrowBackUp,
  IconArrowForwardUp,
} from "../../icons/tabler-icons"
import { Button } from "../ui/button"
import { cn } from "@repo/ui/lib/utils"

export interface RichTextEditorProps {
  content?: string
  onChange?: (html: string) => void
  placeholder?: string
  editable?: boolean
  className?: string
}

export function RichTextEditor({
  content = "",
  onChange,
  placeholder: _placeholder,
  editable = true,
  className,
}: RichTextEditorProps) {
  const editor = useEditor({
    extensions: [StarterKit],
    content,
    editable,
    onUpdate: ({ editor }) => {
      onChange?.(editor.getHTML())
    },
  })

  React.useEffect(() => {
    if (editor && content !== editor.getHTML()) {
      editor.commands.setContent(content)
    }
  }, [content, editor])

  if (!editor) {
    return (
      <div className={cn("h-32 rounded-md border border-input bg-background animate-pulse", className)} />
    )
  }

  return (
    <div className={cn("rounded-md border border-input bg-background shadow-xs overflow-hidden", className)}>
      {editable && (
        <div className="flex flex-wrap items-center gap-0.5 border-b bg-muted/40 p-1">
          <Button
            type="button"
            variant="ghost"
            size="sm"
            onClick={() => editor.chain().focus().toggleBold().run()}
            className={cn("h-7 w-7 p-0", editor.isActive("bold") && "bg-muted text-primary")}
            title="Bold"
          >
            <IconBold className="h-3.5 w-3.5" />
          </Button>
          <Button
            type="button"
            variant="ghost"
            size="sm"
            onClick={() => editor.chain().focus().toggleItalic().run()}
            className={cn("h-7 w-7 p-0", editor.isActive("italic") && "bg-muted text-primary")}
            title="Italic"
          >
            <IconItalic className="h-3.5 w-3.5" />
          </Button>
          <Button
            type="button"
            variant="ghost"
            size="sm"
            onClick={() => editor.chain().focus().toggleStrike().run()}
            className={cn("h-7 w-7 p-0", editor.isActive("strike") && "bg-muted text-primary")}
            title="Strikethrough"
          >
            <IconStrikethrough className="h-3.5 w-3.5" />
          </Button>
          <Button
            type="button"
            variant="ghost"
            size="sm"
            onClick={() => editor.chain().focus().toggleCode().run()}
            className={cn("h-7 w-7 p-0", editor.isActive("code") && "bg-muted text-primary")}
            title="Code"
          >
            <IconCode className="h-3.5 w-3.5" />
          </Button>

          <div className="h-4 w-px bg-border mx-1" />

          <Button
            type="button"
            variant="ghost"
            size="sm"
            onClick={() => editor.chain().focus().toggleHeading({ level: 1 }).run()}
            className={cn("h-7 w-7 p-0", editor.isActive("heading", { level: 1 }) && "bg-muted text-primary")}
            title="Heading 1"
          >
            <IconH1 className="h-3.5 w-3.5" />
          </Button>
          <Button
            type="button"
            variant="ghost"
            size="sm"
            onClick={() => editor.chain().focus().toggleHeading({ level: 2 }).run()}
            className={cn("h-7 w-7 p-0", editor.isActive("heading", { level: 2 }) && "bg-muted text-primary")}
            title="Heading 2"
          >
            <IconH2 className="h-3.5 w-3.5" />
          </Button>
          <Button
            type="button"
            variant="ghost"
            size="sm"
            onClick={() => editor.chain().focus().toggleBulletList().run()}
            className={cn("h-7 w-7 p-0", editor.isActive("bulletList") && "bg-muted text-primary")}
            title="Bullet list"
          >
            <IconList className="h-3.5 w-3.5" />
          </Button>
          <Button
            type="button"
            variant="ghost"
            size="sm"
            onClick={() => editor.chain().focus().toggleOrderedList().run()}
            className={cn("h-7 w-7 p-0", editor.isActive("orderedList") && "bg-muted text-primary")}
            title="Numbered list"
          >
            <IconListNumbers className="h-3.5 w-3.5" />
          </Button>
          <Button
            type="button"
            variant="ghost"
            size="sm"
            onClick={() => editor.chain().focus().toggleBlockquote().run()}
            className={cn("h-7 w-7 p-0", editor.isActive("blockquote") && "bg-muted text-primary")}
            title="Blockquote"
          >
            <IconBlockquote className="h-3.5 w-3.5" />
          </Button>

          <div className="h-4 w-px bg-border mx-1" />

          <Button
            type="button"
            variant="ghost"
            size="sm"
            onClick={() => editor.chain().focus().undo().run()}
            disabled={!editor.can().undo()}
            className="h-7 w-7 p-0"
            title="Undo"
          >
            <IconArrowBackUp className="h-3.5 w-3.5" />
          </Button>
          <Button
            type="button"
            variant="ghost"
            size="sm"
            onClick={() => editor.chain().focus().redo().run()}
            disabled={!editor.can().redo()}
            className="h-7 w-7 p-0"
            title="Redo"
          >
            <IconArrowForwardUp className="h-3.5 w-3.5" />
          </Button>
        </div>
      )}
      <div className="p-3 text-xs min-h-[120px] focus-within:outline-none">
        <EditorContent editor={editor} className="prose prose-sm dark:prose-invert max-w-none focus:outline-none" />
      </div>
    </div>
  )
}
