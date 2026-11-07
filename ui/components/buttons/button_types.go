package buttons

import (
	"github.com/a-h/templ"
)

type ButtonVariant string
type ButtonSize string
type ButtonType string

const (
	ButtonVariantDefault     ButtonVariant = "default"
	ButtonVariantDestructive ButtonVariant = "destructive"
	ButtonVariantOutline     ButtonVariant = "outline"
	ButtonVariantSecondary   ButtonVariant = "secondary"
	ButtonVariantGhost       ButtonVariant = "ghost"
	ButtonVariantLink        ButtonVariant = "link"
)

const (
	ButtonTypeButton ButtonType = "button"
	ButtonTypeReset  ButtonType = "reset"
	ButtonTypeSubmit ButtonType = "submit"
)

const (
	ButtonSizeDefault ButtonSize = "default"
	ButtonSizeXs      ButtonSize = "xs"
	ButtonSizeSm      ButtonSize = "sm"
	ButtonSizeLg      ButtonSize = "lg"
	ButtonSizeIcon    ButtonSize = "icon"
	ButtonSizeIconXs  ButtonSize = "icon-xs"
	ButtonSizeIconSm  ButtonSize = "icon-sm"
	ButtonSizeIconLg  ButtonSize = "icon-lg"
)

type ButtonProps struct {
	ID         string
	Class      string
	Attributes templ.Attributes
	Variant    ButtonVariant
	Size       ButtonSize
	Href       string
	Target     string
	Disabled   bool
	Type       ButtonType
	Form       string
}
