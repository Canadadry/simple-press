package paginator

import (
// "fmt"
// "strings"
)

import (
	"net/http"
	"strconv"
)

const (
	two = 2
)

func ComputeSliceBound(limit int, offset int, length int) (int, int) {
	if length <= 0 {
		return 0, 0
	}
	if offset < 0 {
		offset = 0
	}
	if limit < 0 {
		limit = 1
	}
	if limit > 100 {
		limit = 100
	}
	start := offset
	end := start + limit
	if start > length {
		start = length - 1
	}
	if end > length {
		end = length
	}
	return start, end
}

type Page struct {
	Limit  int
	Offset int
	Filter map[string]string
	Order  map[string]string
}

//
// func PaginateFind(db *gorm.DB, p Page, scopes []func(*gorm.DB) *gorm.DB, table interface{}) ([]int, int, error) {
// 	idFieldName := "id"
// 	tabler, ok := table.(schema.Tabler)
// 	if ok {
// 		idFieldName = tabler.TableName() + ".id"
// 	}
// 	var ids []int
// 	err := db.Scopes(scopes...).Model(table).Distinct().Pluck(idFieldName, &ids).Error
// 	if err != nil {
// 		return nil, 0, fmt.Errorf("can't load ids of %T: %w", table, err)
// 	}
// 	start, end := ComputeSliceBound(p.Limit, p.Offset, len(ids))
// 	return ids[start:end], len(ids), nil
// }
//
// type Join struct {
// 	Query      string
// 	Dependency string
// }
//
// type Joins map[string]Join
//
// func loadJoin(jName string, joins Joins, joinAlreadyLoaded map[string]bool) []func(*gorm.DB) *gorm.DB {
// 	if joinAlreadyLoaded[jName] {
// 		return nil
// 	}
// 	j, ok := joins[jName]
// 	if !ok {
// 		return nil
// 	}
// 	joinAlreadyLoaded[jName] = true
//
// 	scopes := []func(*gorm.DB) *gorm.DB{}
// 	scopes = append(scopes, loadJoin(j.Dependency, joins, joinAlreadyLoaded)...)
// 	scopes = append(scopes, func(db *gorm.DB) *gorm.DB {
// 		return db.Joins(j.Query)
// 	})
// 	return scopes
// }
//
// func IsLike(field string, value string) func(db *gorm.DB) *gorm.DB {
// 	return func(db *gorm.DB) *gorm.DB {
// 		return db.Where(field+" LIKE ?", "%"+value+"%")
// 	}
// }
//
// func GtThan(field string, value string) func(db *gorm.DB) *gorm.DB {
// 	return func(db *gorm.DB) *gorm.DB {
// 		return db.Where(field+" > ?", value)
// 	}
// }
//
// func LoThan(field string, value string) func(db *gorm.DB) *gorm.DB {
// 	return func(db *gorm.DB) *gorm.DB {
// 		return db.Where(field+" < ?", value)
// 	}
// }
//
// func OrderBy(field string, value string) func(db *gorm.DB) *gorm.DB {
// 	if value != "ASC" {
// 		value = "DESC"
// 	}
// 	return func(db *gorm.DB) *gorm.DB {
// 		return db.Order(field + " " + value)
// 	}
// }
//
// type Filter struct {
// 	Query func(field string, value string) func(db *gorm.DB) *gorm.DB
// 	Field string
// }
//
// type Filters map[string]Filter
//
// func BuildFilterScope(joins Joins, filters Filters, queryFilter map[string]string, orderBy map[string]string) []func(*gorm.DB) *gorm.DB {
//
// 	joinAlreadyLoaded := map[string]bool{}
//
// 	scopes := []func(*gorm.DB) *gorm.DB{}
//
// 	for fName, f := range filters {
// 		fValue, ok := queryFilter[fName]
// 		if !ok {
// 			continue
// 		}
// 		if f.Query == nil {
// 			continue
// 		}
// 		d := getDependencyFromFieldName(f.Field)
// 		scopes = append(scopes, loadJoin(d, joins, joinAlreadyLoaded)...)
// 		scopes = append(scopes, f.Query(f.Field, fValue))
// 	}
//
// 	for fName, f := range filters {
// 		orderValue, ok := orderBy[fName]
// 		if !ok {
// 			continue
// 		}
// 		d := getDependencyFromFieldName(f.Field)
// 		scopes = append(scopes, loadJoin(d, joins, joinAlreadyLoaded)...)
// 		scopes = append(scopes, OrderBy(f.Field, orderValue))
// 	}
// 	return scopes
// }
//
// func getDependencyFromFieldName(field string) string {
// 	parts := strings.Split(field, ".")
// 	if len(parts) != 2 {
// 		return ""
// 	}
// 	return parts[0]
// }

func PageFromRequest(r *http.Request, name string, min int) int {
	pageInt, _ := strconv.Atoi(r.FormValue(name))
	if pageInt < min {
		return min
	}
	return pageInt
}

type Pages struct {
	Current  int
	Previous int
	Next     int
	Total    int
	All      []int
	Link     string
}

func New(current, total, max int, link string) Pages {
	if max > total {
		max = total
	}
	if max <= 0 {
		max = 3
	}
	p := Pages{
		Current:  current,
		Total:    total,
		Previous: current - 1,
		Next:     current + 1,
		Link:     link,
	}
	if p.Current < 0 {
		p.Current = 0
	}
	if p.Next > p.Total {
		p.Next = 0
	}
	start := current - max/two
	if start <= 0 {
		start = 1
	}
	if start+max > total {
		start = start - ((start + max - 1) - total)
	}
	for i := 0; i < max; i++ {
		p.All = append(p.All, i+start)
	}
	return p
}
