package service

import (
	m "github.com/fedorkolmykow/avitojob/pkg/models"
	"reflect"
	"testing"
)

func TestPaginationSuccess(t *testing.T){
	perPage := 3
	page := 1
	l := 5
	trs := make([]m.Transaction, l)
	for i := 0; i < l; i++{
		trs[i].UserId = i
	}
	tr := pagination(trs, page, perPage)
	for i := 0; i < perPage; i++{
		if !reflect.DeepEqual(tr[i], trs[(page - 1) * perPage+ i]){
			t.Errorf("wrong pagination:\n%+v\n%+v\n", tr, trs)
		}
	}
}

func TestPaginationShortPage(t *testing.T){
	perPage := 3
	page := 2
	l := 5
	trs := make([]m.Transaction, l)
	for i := 0; i < l; i++{
		trs[i].UserId = i
	}
	tr := pagination(trs, page, perPage)
	for i := 0; i < len(tr); i++{
		if !reflect.DeepEqual(tr[i], trs[(page - 1) * perPage+ i]){
			t.Errorf("wrong pagination:\n%+v\n%+v\n", tr, trs)
		}
	}
}

func TestPaginationLastPage(t *testing.T){
	perPage := 3
	page := 999
	l := 6
	trs := make([]m.Transaction, l)
	for i := 0; i < l; i++{
		trs[i].UserId = i
	}
	tr := pagination(trs, page, perPage)
	for i := 0; i < len(tr); i++{
		if !reflect.DeepEqual(tr[i], trs[3 + i]){
			t.Errorf("wrong pagination:\n%+v\n%+v\n", tr, trs)
		}
	}
}