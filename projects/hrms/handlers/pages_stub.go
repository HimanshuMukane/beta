package handlers

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

type pagesStub struct{}

var pages pagesStub

func (pagesStub) Dashboard(...interface{}) dummyRenderer            { return dummyRenderer{} }
func (pagesStub) WorkersList(...interface{}) dummyRenderer          { return dummyRenderer{} }
func (pagesStub) WorkerCreate(...interface{}) dummyRenderer         { return dummyRenderer{} }
func (pagesStub) WorkerDetail(...interface{}) dummyRenderer         { return dummyRenderer{} }
func (pagesStub) FinancialsOverview(...interface{}) dummyRenderer   { return dummyRenderer{} }
func (pagesStub) PaymentCreate(...interface{}) dummyRenderer        { return dummyRenderer{} }
func (pagesStub) DeductionsManager(...interface{}) dummyRenderer    { return dummyRenderer{} }
func (pagesStub) EmploymentTypesManager(...interface{}) dummyRenderer { return dummyRenderer{} }

