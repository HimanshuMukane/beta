package messages

type BubbleVariant string

const (
	BubbleVariantDefault     BubbleVariant = "default"
	BubbleVariantSecondary   BubbleVariant = "secondary"
	BubbleVariantMuted       BubbleVariant = "muted"
	BubbleVariantTinted      BubbleVariant = "tinted"
	BubbleVariantOutline     BubbleVariant = "outline"
	BubbleVariantGhost       BubbleVariant = "ghost"
	BubbleVariantDestructive BubbleVariant = "destructive"
)

type BubbleProps struct {
	ID      string
	Class   string
	Variant BubbleVariant
}

type BubbleGroupProps struct {
	ID    string
	Class string
}

type BubbleContentProps struct {
	ID    string
	Class string
}

type BubbleReactionsProps struct {
	ID    string
	Class string
}
