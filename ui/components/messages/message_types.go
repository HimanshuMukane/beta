package messages

type Align string

const (
	AlignStart Align = "start"
	AlignEnd   Align = "end"
)

type MessageProps struct {
	ID    string
	Class string
	Align Align
}

type GroupProps struct {
	ID    string
	Class string
}

type MessageAvatarProps struct {
	ID    string
	Class string
}

type ContentProps struct {
	ID    string
	Class string
}

type HeaderProps struct {
	ID    string
	Class string
}

type FooterProps struct {
	ID    string
	Class string
}
