package handlers

import (
	"net/http"
	"chi-mysql-server/ui/pages"
)

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	rows := []pages.TaskRow{
		{
			Index: 1,
			Name:  "Dashboard",
			Assignees: []pages.Assignee{
				{Initials: "D", BgColor: "#e17055", Name: "Dishita Vaswani"},
			},
			Status:     "TO DO",
			StatusType: "todo",
		},
		{
			Index: 2,
			Name:  "Calendar",
			Assignees: []pages.Assignee{
				{Initials: "D", BgColor: "#e17055", Name: "Dishita Vaswani"},
			},
			Status:     "TO DO",
			StatusType: "todo",
		},
		{
			Index: 3,
			Name:  "Charts",
			Assignees: []pages.Assignee{
				{Initials: "A", BgColor: "#0984e3", Name: "Aniket"},
				{Initials: "V", BgColor: "#6c5ce7", Name: "Vanshita"},
			},
			Status:     "TO DO",
			StatusType: "todo",
		},
		{
			Index: 4,
			Name:  "User Profile",
			Assignees: []pages.Assignee{
				{Initials: "A", BgColor: "#0984e3", Name: "Aniket"},
				{Initials: "V", BgColor: "#6c5ce7", Name: "Vanshita"},
			},
			Status:     "TO DO",
			StatusType: "todo",
		},
		{
			Index: 5,
			Name:  "Form Elements",
			Assignees: []pages.Assignee{
				{Initials: "D", BgColor: "#e17055", Name: "Dishita"},
				{Initials: "A", BgColor: "#0984e3", Name: "Aniket"},
			},
			Status:     "IN PROGRESS",
			StatusType: "inprogress",
		},
		{
			Index: 6,
			Name:  "Form Layout",
			Assignees: []pages.Assignee{
				{Initials: "M", BgColor: "#a29bfe", Name: "Mahek"},
				{Initials: "J", BgColor: "#0984e3", Name: "Janvi"},
			},
			Status:     "IN PROGRESS",
			StatusType: "inprogress",
		},
		{
			Index: 7,
			Name:  "Pages Half Until Error Pa...",
			Assignees: []pages.Assignee{
				{Initials: "J", BgColor: "#0984e3", Name: "Janvi Talreja"},
			},
			Status:     "TO DO",
			StatusType: "todo",
		},
		{
			Index: 8,
			Name:  "Invoice",
			Assignees: []pages.Assignee{
				{Initials: "M", BgColor: "#a29bfe", Name: "Mahek Narwani"},
			},
			Status:     "TO DO",
			StatusType: "todo",
		},
		{
			Index: 9,
			Name:  "UI elements till modal",
			Assignees: []pages.Assignee{
				{Initials: "M", BgColor: "#a29bfe", Name: "Mahek Narwani"},
			},
			Status:     "TO DO",
			StatusType: "todo",
		},
	}

	component := pages.Dashboard("Frontend Project | ClickUp", rows)
	component.Render(r.Context(), w)
}
