package navigation

import (
	"github.com/a-h/templ"
)

type NavSubItem struct {
	Label  string
	URL    string
	Active bool
}

type NavItem struct {
	Label    string
	URL      string
	Icon     templ.Component
	Active   bool
	Badge    string
	Children []NavSubItem
}

type NavGroup struct {
	Label string
	Items []NavItem
}

type WorkspaceSwitcherProps struct {
	Name     string
	Subtitle string
	LogoIcon templ.Component
	LogoBg   string
}

type UserProfileProps struct {
	Name      string
	Email     string
	AvatarURL string
	Initials  string
}

type SidebarProps struct {
	Workspace         WorkspaceSwitcherProps
	Groups            []NavGroup
	User              UserProfileProps
	SearchURL         string
	SearchPlaceholder string
	ActiveNav         string
}

func hasActiveChild(subs []NavSubItem) bool {
	for _, s := range subs {
		if s.Active {
			return true
		}
	}
	return false
}
