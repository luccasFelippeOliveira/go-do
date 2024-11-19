package cmd

import (
	"reflect"
	"testing"
	"time"
)

func TestInsert(t *testing.T) {
	type args struct {
		todo *Todo
	}

	tests := []struct {
		args    args
		want    *TodoEntity
		name    string
		fields  TodoRepository
		wantErr bool
	}{
		{
			name: "Insert a Todo with status",
			args: args{
				todo: &Todo{
					Description: "Todo Description",
					Status:      StatusDone,
				},
			},
			fields: TodoRepository{
				TodoList: make([]TodoEntity, 0),
				GenerateId: func() string {
					return "123"
				},
				Clock: func() time.Time {
					return time.Date(2024, time.November, 10, 0, 0, 0, 0, time.UTC)
				},
			},
			want: &TodoEntity{
				Entity{
					Id:        "123",
					CreatedAt: time.Date(2024, time.November, 10, 0, 0, 0, 0, time.UTC),
					UpdatedAt: time.Date(2024, time.November, 10, 0, 0, 0, 0, time.UTC),
				},
				Todo{
					Description: "Todo Description",
					Status:      StatusDone,
				},
			},
			wantErr: false,
		},
		{
			name: "Insert a Todo without status",
			args: args{
				todo: &Todo{
					Description: "No Status",
				},
			},
			fields: TodoRepository{
				TodoList: make([]TodoEntity, 0),
				GenerateId: func() string {
					return "123"
				},
				Clock: func() time.Time {
					return time.Date(2024, time.November, 10, 0, 0, 0, 0, time.UTC)
				},
			},
			want: &TodoEntity{
				Entity{
					Id:        "123",
					CreatedAt: time.Date(2024, time.November, 10, 0, 0, 0, 0, time.UTC),
					UpdatedAt: time.Date(2024, time.November, 10, 0, 0, 0, 0, time.UTC),
				},
				Todo{
					Description: "No Status",
					Status:      StatusNotDone,
				},
			},
			wantErr: false,
		},
		{
			name: "Insert a Todo without Description",
			args: args{
				todo: &Todo{},
			},
			fields: TodoRepository{
				TodoList: make([]TodoEntity, 0),
				GenerateId: func() string {
					return "123"
				},
				Clock: func() time.Time {
					return time.Date(2024, time.November, 10, 0, 0, 0, 0, time.UTC)
				},
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			repository := tc.fields

			got, err := repository.Insert(tc.args.todo)
			if (err != nil) != tc.wantErr {
				t.Errorf("TodoRepository.Insert error %v, wantsErr %v", err, tc.wantErr)
				return
			}
			if err != nil && tc.wantErr {
				return
			}

			if !reflect.DeepEqual(got, tc.want) {
				t.Errorf("TodoRepository.Insert() = %v, want %v", got, tc.want)
			}
		})
	}
}

func TestFetchAll(t *testing.T) {
	tests := []struct {
		name    string
		fields  TodoRepository
		want    []TodoEntity
		wantErr bool
	}{
		{
			name: "Fetches all",
			fields: TodoRepository{
				GenerateId: func() string {
					return ""
				},
				TodoList: []TodoEntity{
					{
						Entity{
							Id: "1",
						},
						Todo{
							Description: "Description 1",
							Status:      StatusDone,
						},
					},
					{
						Entity{
							Id: "2",
						},
						Todo{
							Description: "Description 2",
							Status:      StatusDone,
						},
					},
				},
			},
			want: []TodoEntity{
				{
					Entity{
						Id: "1",
					},
					Todo{
						Description: "Description 1",
						Status:      StatusDone,
					},
				},
				{
					Entity{
						Id: "2",
					},
					Todo{
						Description: "Description 2",
						Status:      StatusDone,
					},
				},
			},
			wantErr: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			repository := &tc.fields

			got, err := repository.FetchAll()
			if (err != nil) != tc.wantErr {
				t.Errorf("TodoRepository.FetchAll() error %v, wantsErr %v", err, tc.wantErr)
				return
			}
			if err != nil && tc.wantErr {
				return
			}

			if !reflect.DeepEqual(got, tc.want) {
				t.Errorf("TodoRepository.FetchAll() = %v, want %v", got, tc.want)
			}
		})
	}
}

func TestFetchQuery(t *testing.T) {
	type args map[string]string
	tests := []struct {
		name    string
		fields  TodoRepository
		args    args
		want    []TodoEntity
		wantErr bool
	}{
		{
			name: "Fetch by query id, should return",
			fields: TodoRepository{
				TodoList: []TodoEntity{
					{
						Entity{
							Id: "1234",
						},
						Todo{
							Description: "Description 1234",
							Status:      StatusDone,
						},
					},
					{
						Entity{
							Id: "1235",
						},
						Todo{
							Description: "Description 1235",
							Status:      StatusNotDone,
						},
					},
				},
			},
			args: args{
				"Id": "1234",
			},
			want: []TodoEntity{
				{
					Entity{
						Id: "1234",
					},
					Todo{
						Description: "Description 1234",
						Status:      StatusDone,
					},
				},
			},
			wantErr: false,
		},
		{
			name: "Fetch by query status done, should return all todos that are done",
			fields: TodoRepository{
				TodoList: []TodoEntity{
					{
						Entity{
							Id: "1234",
						},
						Todo{
							Description: "Description 1234",
							Status:      StatusDone,
						},
					},
					{
						Entity{
							Id: "1235",
						},
						Todo{
							Description: "Description 1235",
							Status:      StatusNotDone,
						},
					},
				},
			},
			args: args{
				"Status": "Done",
			},
			want: []TodoEntity{
				{
					Entity{
						Id: "1234",
					},
					Todo{
						Description: "Description 1234",
						Status:      StatusDone,
					},
				},
			},
			wantErr: false,
		},
		{
			name: "Fetch by query createAt less than than 2024-11-10, should return all todos that where created before 2024-11-10",
			fields: TodoRepository{
				TodoList: []TodoEntity{
					{
						Entity{
							Id:        "1234",
							CreatedAt: time.Date(2024, 11, 11, 0, 0, 0, 0, time.UTC),
						},
						Todo{
							Description: "Description 1234",
							Status:      StatusDone,
						},
					},
					{
						Entity{
							Id:        "1235",
							CreatedAt: time.Date(2024, 11, 9, 0, 0, 0, 0, time.UTC),
						},
						Todo{
							Description: "Description 1235",
							Status:      StatusNotDone,
						},
					},
				},
			},
			args: args{
				"CreatedAt_lt": "2024-11-10",
			},
			want: []TodoEntity{
				{
					Entity{
						Id:        "1235",
						CreatedAt: time.Date(2024, 11, 9, 0, 0, 0, 0, time.UTC),
					},
					Todo{
						Description: "Description 1235",
						Status:      StatusNotDone,
					},
				},
			},
			wantErr: false,
		},
		{
			name: "Fetch by query createAt greater than than 2024-11-10, should return all todos that where created after 2024-11-10",
			fields: TodoRepository{
				TodoList: []TodoEntity{
					{
						Entity{
							Id:        "1234",
							CreatedAt: time.Date(2024, 11, 11, 0, 0, 0, 0, time.UTC),
						},
						Todo{
							Description: "Description 1234",
							Status:      StatusDone,
						},
					},
					{
						Entity{
							Id:        "1235",
							CreatedAt: time.Date(2024, 11, 9, 0, 0, 0, 0, time.UTC),
						},
						Todo{
							Description: "Description 1235",
							Status:      StatusNotDone,
						},
					},
				},
			},
			args: args{
				"CreatedAt_gt": "2024-11-10",
			},
			want: []TodoEntity{
				{
					Entity{
						Id:        "1234",
						CreatedAt: time.Date(2024, 11, 11, 0, 0, 0, 0, time.UTC),
					},
					Todo{
						Description: "Description 1234",
						Status:      StatusDone,
					},
				},
			},
			wantErr: false,
		},
		{
			name: "Fetch by query sort by CreateAt in ascending order",
			fields: TodoRepository{
				TodoList: []TodoEntity{
					{
						Entity{
							Id:        "1234",
							CreatedAt: time.Date(2024, 11, 11, 0, 0, 0, 0, time.UTC),
						},
						Todo{
							Description: "Description 1234",
							Status:      StatusDone,
						},
					},
					{
						Entity{
							Id:        "1235",
							CreatedAt: time.Date(2024, 11, 9, 0, 0, 0, 0, time.UTC),
						},
						Todo{
							Description: "Description 1235",
							Status:      StatusNotDone,
						},
					},
				},
			},
			args: args{
				"SortBy": "CreatedAt",
				"Sort":   "asc",
			},
			want: []TodoEntity{
				{
					Entity{
						Id:        "1235",
						CreatedAt: time.Date(2024, 11, 9, 0, 0, 0, 0, time.UTC),
					},
					Todo{
						Description: "Description 1235",
						Status:      StatusNotDone,
					},
				},
				{
					Entity{
						Id:        "1234",
						CreatedAt: time.Date(2024, 11, 11, 0, 0, 0, 0, time.UTC),
					},
					Todo{
						Description: "Description 1234",
						Status:      StatusDone,
					},
				},
			},
			wantErr: false,
		},
		{
			name: "Fetch by query sort by CreateAt in descending order",
			fields: TodoRepository{
				TodoList: []TodoEntity{
					{
						Entity{
							Id:        "1234",
							CreatedAt: time.Date(2024, 11, 11, 0, 0, 0, 0, time.UTC),
						},
						Todo{
							Description: "Description 1234",
							Status:      StatusDone,
						},
					},
					{
						Entity{
							Id:        "1235",
							CreatedAt: time.Date(2024, 11, 9, 0, 0, 0, 0, time.UTC),
						},
						Todo{
							Description: "Description 1235",
							Status:      StatusNotDone,
						},
					},
				},
			},
			args: args{
				"SortBy": "CreatedAt",
				"Sort":   "desc",
			},
			want: []TodoEntity{
				{
					Entity{
						Id:        "1234",
						CreatedAt: time.Date(2024, 11, 11, 0, 0, 0, 0, time.UTC),
					},
					Todo{
						Description: "Description 1234",
						Status:      StatusDone,
					},
				},
				{
					Entity{
						Id:        "1235",
						CreatedAt: time.Date(2024, 11, 9, 0, 0, 0, 0, time.UTC),
					},
					Todo{
						Description: "Description 1235",
						Status:      StatusNotDone,
					},
				},
			},
			wantErr: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			repository := tc.fields

			got, err := repository.FetchByQuery(tc.args)

			if (err != nil) != tc.wantErr {
				t.Errorf("TodoRepository.FetchByQuery() error %v, wantsErr %v", err, tc.wantErr)
			}

			if err != nil && tc.wantErr {
				return
			}

			if !reflect.DeepEqual(got, tc.want) {
				t.Errorf("TodoRepository.FetchByQuery() = %v, want %v", got, tc.want)
			}
		})
	}
}

func TestUpdate(t *testing.T) {
	type args struct {
		id    string
		model Todo
	}
	tests := []struct {
		want    *TodoEntity
		args    args
		name    string
		fields  TodoRepository
		wantErr bool
	}{
		{
			name: "Update todo description",
			args: args{
				id: "1234",
				model: Todo{
					Description: "New Description",
				},
			},
			fields: TodoRepository{
				TodoList: []TodoEntity{
					{
						Entity{
							Id: "1234",
						},
						Todo{
							Description: "Description 1234",
							Status:      StatusDone,
						},
					},
					{
						Entity{
							Id: "1235",
						},
						Todo{
							Description: "Description 1235",
							Status:      StatusNotDone,
						},
					},
				},
				Clock: func() time.Time {
					return time.Date(2024, time.November, 10, 0, 0, 0, 0, time.UTC)
				},
			},
			want: &TodoEntity{
				Entity{
					Id:        "1234",
					UpdatedAt: time.Date(2024, time.November, 10, 0, 0, 0, 0, time.UTC),
				},
				Todo{
					Description: "New Description",
					Status:      StatusDone,
				},
			},
		},
		{
			name: "Update todo status",
			args: args{
				id: "1234",
				model: Todo{
					Status: StatusNotDone,
				},
			},
			fields: TodoRepository{
				TodoList: []TodoEntity{
					{
						Entity{
							Id: "1234",
						},
						Todo{
							Description: "Description 1234",
							Status:      StatusDone,
						},
					},
					{
						Entity{
							Id: "1235",
						},
						Todo{
							Description: "Description 1235",
							Status:      StatusNotDone,
						},
					},
				},
				Clock: func() time.Time {
					return time.Date(2024, time.November, 10, 0, 0, 0, 0, time.UTC)
				},
			},
			want: &TodoEntity{
				Entity{
					Id:        "1234",
					UpdatedAt: time.Date(2024, time.November, 10, 0, 0, 0, 0, time.UTC),
				},
				Todo{
					Description: "Description 1234",
					Status:      StatusNotDone,
				},
			},
		},
		{
			name: "Update a entity with invalid id results in error",
			args: args{
				id:    "12345",
				model: Todo{},
			},
			fields: TodoRepository{
				TodoList: []TodoEntity{
					{
						Entity{
							Id: "1234",
						},
						Todo{
							Description: "Description 1234",
							Status:      StatusDone,
						},
					},
					{
						Entity{
							Id: "1235",
						},
						Todo{
							Description: "Description 1235",
							Status:      StatusNotDone,
						},
					},
				},
				Clock: func() time.Time {
					return time.Date(2024, time.November, 10, 0, 0, 0, 0, time.UTC)
				},
			},
			wantErr: true,
		},
		{
			name: "Update a entity with invalid model results in error",
			args: args{
				id:    "1234",
				model: Todo{},
			},
			fields: TodoRepository{
				TodoList: []TodoEntity{
					{
						Entity{
							Id: "1234",
						},
						Todo{
							Description: "Description 1234",
							Status:      StatusDone,
						},
					},
					{
						Entity{
							Id: "1235",
						},
						Todo{
							Description: "Description 1235",
							Status:      StatusNotDone,
						},
					},
				},
				Clock: func() time.Time {
					return time.Date(2024, time.November, 10, 0, 0, 0, 0, time.UTC)
				},
			},
			wantErr: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			repository := tc.fields

			got, err := repository.Update(tc.args.id, tc.args.model)

			if (err != nil) != tc.wantErr {
				t.Errorf("TodoRepository.Update() error %v, wantsErr %v", err, tc.wantErr)
			}

			if err != nil && tc.wantErr {
				return
			}

			if !reflect.DeepEqual(got, tc.want) {
				t.Errorf("TodoRepository.Update() = %v, want %v", got, tc.want)
			}
		})
	}
}

func TestDelete(t *testing.T) {
	type args struct {
		id string
	}
	tests := []struct {
		want    *TodoEntity
		args    args
		name    string
		fields  TodoRepository
		wantErr bool
	}{
		{
			name: "Delete a todo",
			args: args{
				id: "1234",
			},
			fields: TodoRepository{
				TodoList: []TodoEntity{
					{
						Entity{
							Id: "1234",
						},
						Todo{
							Description: "Description 1234",
							Status:      StatusDone,
						},
					},
					{
						Entity{
							Id: "1235",
						},
						Todo{
							Description: "Description 1235",
							Status:      StatusNotDone,
						},
					},
				},
			},
			want: &TodoEntity{
				Entity{
					Id: "1234",
				},
				Todo{
					Description: "Description 1234",
					Status:      StatusDone,
				},
			},
		},
		{
			name: "Delete a todo with a invalid id",
			args: args{
				id: "12345",
			},
			fields: TodoRepository{
				TodoList: []TodoEntity{
					{
						Entity{
							Id: "1234",
						},
						Todo{
							Description: "Description 1234",
							Status:      StatusDone,
						},
					},
					{
						Entity{
							Id: "1235",
						},
						Todo{
							Description: "Description 1235",
							Status:      StatusNotDone,
						},
					},
				},
			},
			wantErr: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			repository := tc.fields

			got, err := repository.Delete(tc.args.id)

			if (err != nil) != tc.wantErr {
				t.Errorf("TodoRepository.Delete() error %v, wantsErr %v", err, tc.wantErr)
			}

			if err != nil && tc.wantErr {
				return
			}

			if !reflect.DeepEqual(got, tc.want) {
				t.Errorf("TodoRepository.Delete() = %v, want %v", got, tc.want)
			}
		})
	}
}
