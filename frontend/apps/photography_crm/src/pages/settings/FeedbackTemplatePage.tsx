import * as React from 'react'
import {
  IconPlus,
  IconTrash,
  IconEdit,
  IconDeviceFloppy,
  IconStarFilled,
  IconCheck,
  IconMessageHeart,
  IconAlignLeft,
} from '@tabler/icons-react'
import {
  Button,
  Input,
  Label,
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogDescription,
  DialogFooter,
  Tooltip,
  TooltipTrigger,
  TooltipContent,
  toast,
} from '@repo/ui'

export type QuestionType = 'rating' | 'textarea'

export interface FeedbackQuestion {
  id: string
  title: string
  type: QuestionType
  placeholder?: string
}

const INITIAL_QUESTIONS: FeedbackQuestion[] = [
  {
    id: 'q-1',
    title: 'How satisfied were you with the birthday photography session overall?',
    type: 'rating',
  },
  {
    id: 'q-2',
    title: "How was the photographer's interaction with kids and guests?",
    type: 'rating',
  },
  {
    id: 'q-3',
    title: 'How likely are you to recommend our birthday photography services to others?',
    type: 'rating',
  },
  {
    id: 'q-4',
    title: 'Was there anything we could improve for future events?',
    type: 'textarea',
    placeholder: 'Share your suggestions or specific areas for improvement...',
  },
  {
    id: 'q-5',
    title: 'Any suggestions or additional services you would like us to offer?',
    type: 'textarea',
    placeholder: 'e.g. photobooth props, instant polaroid prints...',
  },
  {
    id: 'q-6',
    title: 'Would you like to leave a testimonial or review for us?',
    type: 'textarea',
    placeholder: 'Write a brief quote that we can feature on our client wall...',
  },
]

export function FeedbackTemplatePage() {
  const [templateName, setTemplateName] = React.useState('Birthday Shoot Feedback')
  const [questions, setQuestions] = React.useState<FeedbackQuestion[]>(INITIAL_QUESTIONS)

  // Question editing modal
  const [editingQuestion, setEditingQuestion] = React.useState<FeedbackQuestion | null>(null)
  const [modalTitle, setModalTitle] = React.useState('')
  const [modalType, setModalType] = React.useState<QuestionType>('rating')
  const [isModalOpen, setIsModalOpen] = React.useState(false)

  const handleOpenAdd = () => {
    setEditingQuestion(null)
    setModalTitle('')
    setModalType('rating')
    setIsModalOpen(true)
  }

  const handleOpenEdit = (q: FeedbackQuestion) => {
    setEditingQuestion(q)
    setModalTitle(q.title)
    setModalType(q.type)
    setIsModalOpen(true)
  }

  const handleSaveQuestion = () => {
    if (!modalTitle.trim()) {
      toast.error('Question title cannot be blank')
      return
    }

    if (editingQuestion) {
      setQuestions((prev) =>
        prev.map((q) =>
          q.id === editingQuestion.id ? { ...q, title: modalTitle.trim(), type: modalType } : q
        )
      )
      toast.success('Question updated')
    } else {
      const newQ: FeedbackQuestion = {
        id: `q-${Date.now()}`,
        title: modalTitle.trim(),
        type: modalType,
        placeholder: modalType === 'textarea' ? 'Type feedback here...' : undefined,
      }
      setQuestions((prev) => [...prev, newQ])
      toast.success('Question added')
    }

    setIsModalOpen(false)
  }

  const handleDeleteQuestion = (id: string) => {
    setQuestions((prev) => prev.filter((q) => q.id !== id))
    toast.success('Question removed')
  }

  const handleSaveAll = () => {
    toast.success(`Feedback template "${templateName}" saved`)
  }

  return (
    <div className="max-w-3xl mx-auto space-y-6 pb-16">
      {/* Top Header Bar (Matches Screenshot 4) */}
      <div className="flex items-center justify-between border-b pb-4">
        <div className="flex items-center gap-2.5">
          <span className="flex h-8 w-8 items-center justify-center rounded-lg bg-primary/10 text-primary border border-primary/20">
            <IconMessageHeart className="h-4 w-4" />
          </span>
          <div>
            <h1 className="text-xl font-bold tracking-tight text-foreground">
              Edit Feedback Template
            </h1>
            <p className="text-xs text-muted-foreground">
              Configure post-event client survey questions and rating metrics
            </p>
          </div>
        </div>

        <Button
          type="button"
          onClick={handleSaveAll}
          className="h-8 text-xs font-semibold gap-1.5 bg-primary text-primary-foreground shadow-xs"
        >
          <IconDeviceFloppy className="h-3.5 w-3.5" />
          Save
        </Button>
      </div>

      <div className="bg-card border border-border/80 shadow-2xs rounded-xl p-6 sm:p-8 space-y-6">
        {/* Template Name Input */}
        <div className="space-y-1.5 max-w-sm">
          <Label className="text-xs font-semibold text-foreground">Feedback Template Name</Label>
          <Input
            value={templateName}
            onChange={(e) => setTemplateName(e.target.value)}
            className="h-9 text-xs bg-background font-medium"
            placeholder="e.g. Wedding Client Survey"
          />
        </div>

        <div className="space-y-4 pt-2">
          <div className="flex items-center justify-between border-b pb-2">
            <h2 className="text-xs font-bold text-foreground uppercase tracking-wider">
              Feedback Questions
            </h2>
            <span className="text-[11px] text-muted-foreground">
              {questions.length} questions configured
            </span>
          </div>

          <div className="space-y-6">
            {questions.map((q, idx) => (
              <div
                key={q.id}
                className="group relative p-3 rounded-lg border border-transparent hover:border-border/60 hover:bg-muted/15 transition-all"
              >
                <div className="flex items-start justify-between gap-4">
                  <div className="space-y-2 flex-1">
                    <p className="text-xs font-semibold text-foreground leading-snug">
                      <span className="text-muted-foreground font-normal mr-1">{idx + 1}.</span>
                      {q.title}
                    </p>

                    {/* Question Interactive Visual Representation */}
                    {q.type === 'rating' ? (
                      <div className="flex items-center gap-1 text-muted/80 pt-0.5">
                        {[1, 2, 3, 4, 5].map((star) => (
                          <IconStarFilled
                            key={star}
                            className="h-4 w-4 text-muted-foreground/30 hover:text-amber-400 cursor-pointer transition-colors"
                          />
                        ))}
                      </div>
                    ) : (
                      <div className="h-16 w-full rounded border border-border/70 bg-background/50 p-2 text-[11px] text-muted-foreground/60 italic">
                        {q.placeholder || 'Client long-form text response field...'}
                      </div>
                    )}
                  </div>

                  {/* Actions (Matches Screenshot 4 pencil & minus/trash) */}
                  <div className="flex items-center gap-1 shrink-0 pt-0.5">
                    <Tooltip>
                      <TooltipTrigger asChild>
                        <Button
                          type="button"
                          variant="ghost"
                          size="icon"
                          onClick={() => handleOpenEdit(q)}
                          className="h-7 w-7 text-muted-foreground hover:text-foreground"
                        >
                          <IconEdit className="h-3.5 w-3.5" />
                        </Button>
                      </TooltipTrigger>
                      <TooltipContent side="top">
                        Edit question
                      </TooltipContent>
                    </Tooltip>

                    <Tooltip>
                      <TooltipTrigger asChild>
                        <Button
                          type="button"
                          variant="ghost"
                          size="icon"
                          onClick={() => handleDeleteQuestion(q.id)}
                          className="h-7 w-7 text-muted-foreground hover:text-destructive hover:bg-destructive/10"
                        >
                          <IconTrash className="h-3.5 w-3.5" />
                        </Button>
                      </TooltipTrigger>
                      <TooltipContent side="top">
                        Delete question
                      </TooltipContent>
                    </Tooltip>
                  </div>
                </div>
              </div>
            ))}
          </div>

          <div className="pt-2">
            <Button
              type="button"
              variant="ghost"
              onClick={handleOpenAdd}
              className="text-xs text-primary font-semibold hover:bg-primary/10 gap-1.5 h-8 px-2"
            >
              <IconPlus className="h-3.5 w-3.5" />
              Add Question
            </Button>
          </div>
        </div>
      </div>

      {/* Add / Edit Question Modal */}
      <Dialog open={isModalOpen} onOpenChange={setIsModalOpen}>
        <DialogContent className="sm:max-w-md">
          <DialogHeader>
            <DialogTitle className="text-sm font-bold">
              {editingQuestion ? 'Edit Question' : 'Add Feedback Question'}
            </DialogTitle>
            <DialogDescription className="text-xs text-muted-foreground">
              Define the question prompt and choose the response format.
            </DialogDescription>
          </DialogHeader>

          <div className="space-y-4 py-2">
            <div className="space-y-1.5">
              <Label className="text-xs font-semibold">Question Prompt</Label>
              <Input
                value={modalTitle}
                onChange={(e) => setModalTitle(e.target.value)}
                placeholder="e.g. How satisfied were you with the shoot?"
                className="h-8 text-xs"
                autoFocus
              />
            </div>

            <div className="space-y-1.5">
              <Label className="text-xs font-semibold">Response Type</Label>
              <Select
                value={modalType}
                onValueChange={(val) => setModalType(val as QuestionType)}
              >
                <SelectTrigger className="h-8 text-xs bg-background">
                  <SelectValue />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="rating">
                    <span className="flex items-center gap-1.5">
                      <IconStarFilled className="h-3.5 w-3.5 text-amber-500" />
                      5-Star Rating
                    </span>
                  </SelectItem>
                  <SelectItem value="textarea">
                    <span className="flex items-center gap-1.5">
                      <IconAlignLeft className="h-3.5 w-3.5 text-blue-500" />
                      Long Text Feedback
                    </span>
                  </SelectItem>
                </SelectContent>
              </Select>
            </div>
          </div>

          <DialogFooter>
            <Button
              type="button"
              variant="outline"
              onClick={() => setIsModalOpen(false)}
              className="h-8 text-xs"
            >
              Cancel
            </Button>
            <Button
              type="button"
              onClick={handleSaveQuestion}
              className="h-8 text-xs font-semibold gap-1 bg-primary text-primary-foreground"
            >
              <IconCheck className="h-3.5 w-3.5" />
              {editingQuestion ? 'Save Changes' : 'Add Question'}
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>
    </div>
  )
}

