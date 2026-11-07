package cards

import (
	"github.com/a-h/templ"
)

type CardSize string

const (
	CardSizeDefault CardSize = "default"
	CardSizeSm      CardSize = "sm"
)

type CardProps struct {
	ID         string
	Class      string
	Attributes templ.Attributes
	Size       CardSize
}

type CardHeaderProps struct {
	ID         string
	Class      string
	Attributes templ.Attributes
}

type CardTitleProps struct {
	ID         string
	Class      string
	Attributes templ.Attributes
}

type CardDescriptionProps struct {
	ID         string
	Class      string
	Attributes templ.Attributes
}

type CardActionProps struct {
	ID         string
	Class      string
	Attributes templ.Attributes
}

type CardContentProps struct {
	ID         string
	Class      string
	Attributes templ.Attributes
}

type CardFooterProps struct {
	ID         string
	Class      string
	Attributes templ.Attributes
}
