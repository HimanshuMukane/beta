package pages

import (
	"context"
	"net/http"
)

type dummyRenderer struct{}

func (dummyRenderer) Render(_ context.Context, w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	_, err := w.Write([]byte(`{"status":"ok"}`))
	return err
}

type DashboardData struct {
	TotalLeadsCount      int
	QualifiedLeadsCount  int
	ActiveDealsCount     int
	PipelineValueMinor   int64
	ActiveProjectsCount  int
	ProjectsDelayedCount int
	UpcomingMilestones   any
	RecentDeals          any
}

func Dashboard(...interface{}) dummyRenderer        { return dummyRenderer{} }
func LeadsList(...interface{}) dummyRenderer        { return dummyRenderer{} }
func LeadCreate(...interface{}) dummyRenderer       { return dummyRenderer{} }
func LeadDetail(...interface{}) dummyRenderer       { return dummyRenderer{} }
func DealsPipeline(...interface{}) dummyRenderer     { return dummyRenderer{} }
func DealCreate(...interface{}) dummyRenderer       { return dummyRenderer{} }
func DealDetail(...interface{}) dummyRenderer       { return dummyRenderer{} }
func ProjectsList(...interface{}) dummyRenderer     { return dummyRenderer{} }
func ProjectCreate(...interface{}) dummyRenderer    { return dummyRenderer{} }
func ProjectDetail(...interface{}) dummyRenderer    { return dummyRenderer{} }
func ScheduleCalendar(...interface{}) dummyRenderer { return dummyRenderer{} }
func QuotesList(...interface{}) dummyRenderer       { return dummyRenderer{} }
func QuoteCreate(...interface{}) dummyRenderer      { return dummyRenderer{} }
func ContactsList(...interface{}) dummyRenderer     { return dummyRenderer{} }
func ContactCreate(...interface{}) dummyRenderer    { return dummyRenderer{} }
func ContactDetail(...interface{}) dummyRenderer    { return dummyRenderer{} }
