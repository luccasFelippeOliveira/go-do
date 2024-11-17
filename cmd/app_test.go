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
			want: &TodoEntity{
				Entity{
					Id:        "123",
					CreatedAt: time.Time{},
					UpdatedAt: time.Time{},
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
			want: &TodoEntity{
				Entity{
					Id:        "123",
					CreatedAt: time.Time{},
					UpdatedAt: time.Time{},
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
			want:    nil,
			wantErr: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			repository := &TodoRepository{
				TodoList: make([]TodoEntity, 0),
				GenerateId: func() string {
					return "123"
				},
			}

			got, err := repository.Insert(tc.args.todo)
			if (err != nil) != tc.wantErr {
				t.Errorf("TodoRepository.Insert error %v, wantsErr %v", err, tc.wantErr)
				return
			}
			if err != nil && tc.wantErr {
				return
			}

			got.CreatedAt = time.Time{}
			got.UpdatedAt = time.Time{}

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
}

func TestDelete(t *testing.T) {
}
