package cmd

import (
	"cmp"
	"errors"
	"fmt"
	"regexp"
	"slices"
	"strings"
	"time"
)

type TodoStatus string

const (
	StatusDone    TodoStatus = "Done"
	StatusNotDone TodoStatus = "NotDone"
)

type Entity struct {
	CreatedAt time.Time
	UpdatedAt time.Time
	Id        string
}

type Todo struct {
	Description string
	Status      TodoStatus
}

type TodoEntity struct {
	Entity
	Todo
}

type GenerateId func() string

type TodoRepository struct {
	GenerateId GenerateId
	TodoList   []TodoEntity
}

func (r *TodoRepository) Insert(todo *Todo) (*TodoEntity, error) {
	if len(todo.Description) == 0 {
		return nil, errors.New("description is not valid, it must be a valid string")
	}

	var status TodoStatus
	if len(todo.Status) != 0 {
		status = todo.Status
	} else {
		status = StatusNotDone
	}

	todoEntity := &TodoEntity{
		Entity{
			Id:        r.GenerateId(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		Todo{
			Description: todo.Description,
			Status:      status,
		},
	}

	// Insert the entity into the TodoList
	r.TodoList = append(r.TodoList, *todoEntity)

	return todoEntity, nil
}

func (r *TodoRepository) FetchAll() ([]TodoEntity, error) {
	if r.TodoList == nil {
		return nil, errors.New("repository not initialized")
	}
	return r.TodoList, nil
}

func (r *TodoRepository) filterById(id string) *TodoEntity {
	index := slices.IndexFunc(r.TodoList, func(todoEntity TodoEntity) bool {
		return todoEntity.Id == id
	})

	if index > 0 {
		return &r.TodoList[index]
	} else {
		return nil
	}
}

func filterString(field string, value string) bool {
	match, err := regexp.MatchString(value, field)
	if err != nil {
		return false
	}

	return match
}

func filterDate(field time.Time, value time.Time) bool {
	truncField := field.Truncate(time.Hour * 24)
	truncValue := value.Truncate(time.Hour * 24)

	return truncField.Equal(truncValue)
}

func validateQuery(query map[string]string) error {
	for qf, qv := range query {
		switch qf {
		case "Id":
			continue
		case "CreatedAt":
		case "UpdatedAt":
		case "CreatedAt_lt":
		case "UpdatedAt_lt":
		case "CreatedAt_gt":
		case "UpdatedAt_gt":
			// Check if is valid format. YYYY-MM-dd
			_, err := time.Parse("YYY-MM-dd", qv)
			if err != nil {
				return fmt.Errorf("Invalid time format %e", err)
			}
		case "Description":
			continue
		case "Status":
			if qv != string(StatusDone) && qv != string(StatusNotDone) {
				return fmt.Errorf("Invalid Status query value")
			}
			continue
		case "Sort":
			if qv != "asc" && qv != "desc" {
				return fmt.Errorf("Invalid Status query value")
			}
			continue
		case "SortBy":
			if qv != "Id" && qv != "CreatedAt" && qv != "UpdatedAt" && qv != "Description" {
				return fmt.Errorf("Invalid sort, sort by only accepts Id, CreatedAt, UpdatedAt, Description")
			}
			continue
		default:
			return fmt.Errorf("Invalid query field. Got %v", qf)
		}
	}

	_, hasSort := query["Sort"]
	_, hasSortBy := query["SortBy"]

	if hasSort && !hasSortBy {
		return fmt.Errorf("If the query has a sort, then the field must be defined by SortBy")
	}

	if hasSortBy && !hasSort {
		return fmt.Errorf("If the query has a sort by a field, then the direction must be defined by Sort")
	}

	return nil
}

// matchDate takes two arguments and compares then.
// will return 0 if both are equal
// will return 1 if d1 > d2
// will return -1 if d1 < d2
func matchDate(d1 time.Time, d2 time.Time) int {
	d1Trunc := d1.Truncate(24 * time.Hour)
	d2Trunc := d2.Truncate(24 * time.Hour)

	return d1Trunc.Compare(d2Trunc)
}

var (
	createdAtRegex = regexp.MustCompile(`CreatedAt`)
	updatedAtRegex = regexp.MustCompile(`UpdatedAt`)
)

func matchQuery(entity *TodoEntity, query map[string]string) bool {
	if entity == nil {
		return false
	}

	isMatch := true

	for qf, qv := range query {
		switch field := qf; {
		case field == "Id":
			isMatch = isMatch && entity.Id == qv
		case createdAtRegex.MatchString(field):
			qvDate, _ := time.Parse("2006-01-02", qv)
			res := matchDate(entity.CreatedAt, qvDate)
			if strings.HasSuffix(qf, "_lt") {
				isMatch = isMatch && res < 0
			} else if strings.HasSuffix(qf, "_gt") {
				isMatch = isMatch && res > 0
			} else {
				isMatch = isMatch && res == 0
			}
		case updatedAtRegex.MatchString(field):
			qvDate, _ := time.Parse("2006-01-02", qv)
			res := matchDate(entity.UpdatedAt, qvDate)
			if strings.HasSuffix(qf, "_lt") {
				isMatch = isMatch && res < 0
			} else if strings.HasSuffix(qf, "_gt") {
				isMatch = isMatch && res > 0
			} else {
				isMatch = isMatch && res == 0
			}
		case field == "Description":
			isMatch = isMatch && entity.Description == qv
		case field == "Status":
			isMatch = isMatch && string(entity.Status) == qv
		}
	}

	return isMatch
}

func sortQuery(entity1 *TodoEntity, entity2 *TodoEntity, sortBy string, order string) int {
	switch sortBy {
	case "Id":
		c := cmp.Compare(entity1.Id, entity2.Id)
		if order == "asc" {
			return c
		} else {
			return -c
		}
	case "CreatedAt":
		c := matchDate(entity1.CreatedAt, entity2.CreatedAt)
		if order == "asc" {
			return c
		} else {
			return -c
		}
	case "UpdatedAt":
		c := matchDate(entity1.UpdatedAt, entity2.UpdatedAt)
		if order == "asc" {
			return c
		} else {
			return -c
		}
	case "Description":
		c := cmp.Compare(entity1.Description, entity2.Description)
		if order == "asc" {
			return c
		} else {
			return -c
		}
	}
	return 0
}

func (r *TodoRepository) FetchByQuery(query map[string]string) ([]TodoEntity, error) {
	if r.TodoList == nil {
		return nil, errors.New("repostitory not initialized")
	}

	// Validate the query.
	queryErr := validateQuery(query)
	if queryErr != nil {
		return nil, queryErr
	}

	sortDirection, hasSort := query["Sort"]
	sortField := query["SortBy"]

	result := make([]TodoEntity, 0)
	for _, t := range r.TodoList {
		if matchQuery(&t, query) {
			result = append(result, t)
		}
	}

	if hasSort && len(result) > 0 {
		slices.SortFunc(result, func(e1 TodoEntity, e2 TodoEntity) int {
			return sortQuery(&e1, &e2, sortField, sortDirection)
		})
	}

	return result, nil
}
