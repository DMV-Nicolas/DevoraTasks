package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	mockdb "github.com/DMV-Nicolas/DevoraTasks/db/mock"
	db "github.com/DMV-Nicolas/DevoraTasks/db/sqlc"
	"github.com/DMV-Nicolas/DevoraTasks/util"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestCreateTaskAPI(t *testing.T) {
	user, _ := randomUser(t)
	task := randomTask(user.ID)

	testCases := []struct {
		name          string
		body          map[string]any
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			body: map[string]any{
				"user_id":     task.UserID,
				"title":       task.Title,
				"description": task.Description,
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.CreateTaskParams{
					UserID:      task.UserID,
					Title:       task.Title,
					Description: task.Description,
				}
				store.EXPECT().
					CreateTask(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(task, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusCreated, recorder.Code)
				requireBodyMatchTask(t, recorder.Body, task)
			},
		},
		{
			name: "InternalServerError",
			body: map[string]any{
				"user_id":     task.UserID,
				"title":       task.Title,
				"description": task.Description,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					CreateTask(gomock.Any(), gomock.Any()).
					Times(1).
					Return(db.Task{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name: "InvalidUserID",
			body: map[string]any{
				"user_id":     task.UserID,
				"title":       task.Title,
				"description": task.Description,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					CreateTask(gomock.Any(), gomock.Any()).
					Times(1).
					Return(db.Task{}, db.ErrForeignKeyViolation)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "InvalidTitle",
			body: map[string]any{
				"user_id":     task.UserID,
				"title":       "",
				"description": task.Description,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					CreateTask(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)
			tc.buildStubs(store)

			// marshal data body to json
			data, err := json.Marshal(tc.body)
			require.NoError(t, err)

			// start test server and send request
			server := NewServer(store)
			recorder := httptest.NewRecorder()

			url := "/tasks"
			request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
			require.NoError(t, err)

			server.router.ServeHTTP(recorder, request)

			tc.checkResponse(t, recorder)
		})
	}
}

func TestGetTaskAPI(t *testing.T) {
	user, _ := randomUser(t)
	task := randomTask(user.ID)

	testCases := []struct {
		name          string
		taskID        int64
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:   "OK",
			taskID: task.ID,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetTask(gomock.Any(), gomock.Eq(task.ID)).
					Times(1).
					Return(task, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchTask(t, recorder.Body, task)
			},
		},
		{
			name:   "NotFound",
			taskID: task.ID,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetTask(gomock.Any(), gomock.Any()).
					Times(1).
					Return(db.Task{}, sql.ErrNoRows)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name:   "InternalServerError",
			taskID: task.ID,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetTask(gomock.Any(), gomock.Any()).
					Times(1).
					Return(db.Task{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name:   "InvalidID",
			taskID: 0,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetTask(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)
			tc.buildStubs(store)

			// start test server and send request
			server := NewServer(store)
			recorder := httptest.NewRecorder()

			url := fmt.Sprintf("/tasks/%d", tc.taskID)
			request, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			server.router.ServeHTTP(recorder, request)

			tc.checkResponse(t, recorder)
		})
	}
}

func TestListTasksAPI(t *testing.T) {
	user, _ := randomUser(t)

	n := 5
	tasks := make([]db.Task, n)
	for i := 0; i < n; i++ {
		tasks[i] = randomTask(user.ID)
	}

	type Query struct {
		offset int32
		limit  int32
	}

	testCases := []struct {
		name          string
		query         Query
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			query: Query{
				offset: 0,
				limit:  int32(n),
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.ListTasksParams{
					Offset: 0,
					Limit:  5,
				}

				store.EXPECT().
					ListTasks(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(tasks, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchTasks(t, recorder.Body, tasks)
			},
		},
		{
			name: "InternalServerError",
			query: Query{
				offset: 0,
				limit:  int32(n),
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					ListTasks(gomock.Any(), gomock.Any()).
					Times(1).
					Return([]db.Task{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name: "InvalidOffset",
			query: Query{
				offset: -1,
				limit:  int32(n),
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					ListTasks(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "InvalidLimit",
			query: Query{
				offset: 0,
				limit:  0,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					ListTasks(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)
			tc.buildStubs(store)

			// start test server and send request
			server := NewServer(store)
			recorder := httptest.NewRecorder()

			url := "/tasks"
			request, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			// Add query parameters to request URL
			q := request.URL.Query()
			q.Add("offset", fmt.Sprintf("%d", tc.query.offset))
			q.Add("limit", fmt.Sprintf("%d", tc.query.limit))
			request.URL.RawQuery = q.Encode()

			server.router.ServeHTTP(recorder, request)

			tc.checkResponse(t, recorder)
		})
	}
}

func TestUpdateTaskAPI(t *testing.T) {
	user, _ := randomUser(t)
	task1 := randomTask(user.ID)
	task2 := randomTask(user.ID)
	task2.ID = task1.ID
	task2.Done = true

	testCases := []struct {
		name          string
		body          map[string]any
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			body: map[string]any{
				"id":          task2.ID,
				"title":       task2.Title,
				"description": task2.Description,
				"done":        task2.Done,
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.UpdateTaskParams{
					ID:          task2.ID,
					Title:       task2.Title,
					Description: task2.Description,
					Done:        task2.Done,
				}

				store.EXPECT().
					UpdateTask(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(task2, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchTask(t, recorder.Body, task2)
			},
		},
		{
			name: "NotFound",
			body: map[string]any{
				"id":          task2.ID,
				"title":       task2.Title,
				"description": task2.Description,
				"done":        task2.Done,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					UpdateTask(gomock.Any(), gomock.Any()).
					Times(1).
					Return(db.Task{}, sql.ErrNoRows)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name: "InternalServerError",
			body: map[string]any{
				"id":          task2.ID,
				"title":       task2.Title,
				"description": task2.Description,
				"done":        task2.Done,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					UpdateTask(gomock.Any(), gomock.Any()).
					Times(1).
					Return(db.Task{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name: "InvalidID",
			body: map[string]any{
				"id":          0,
				"title":       task2.Title,
				"description": task2.Description,
				"done":        task2.Done,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					UpdateTask(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "InvalidTitle",
			body: map[string]any{
				"id":          task2.ID,
				"title":       "",
				"description": task2.Description,
				"done":        task2.Done,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					UpdateTask(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)
			tc.buildStubs(store)

			// marshal data body to json
			data, err := json.Marshal(tc.body)
			require.NoError(t, err)

			// start test server and send request
			server := NewServer(store)
			recorder := httptest.NewRecorder()

			url := "/tasks"
			request, err := http.NewRequest(http.MethodPut, url, bytes.NewReader(data))
			require.NoError(t, err)

			server.router.ServeHTTP(recorder, request)

			tc.checkResponse(t, recorder)
		})
	}
}

func TestDeleteTaskAPI(t *testing.T) {
	user, _ := randomUser(t)
	task := randomTask(user.ID)

	testCases := []struct {
		name          string
		body          map[string]any
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			body: map[string]any{"id": task.ID},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					DeleteTask(gomock.Any(), gomock.Eq(task.ID)).
					Times(1).
					Return(nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNoContent, recorder.Code)
			},
		},
		{
			name: "NotFound",
			body: map[string]any{"id": task.ID},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					DeleteTask(gomock.Any(), gomock.Any()).
					Times(1).
					Return(sql.ErrNoRows)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name: "InternalServerError",
			body: map[string]any{"id": task.ID},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					DeleteTask(gomock.Any(), gomock.Any()).
					Times(1).
					Return(sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name: "InvalidID",
			body: map[string]any{"id": 0},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					DeleteTask(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)
			tc.buildStubs(store)

			// marshal data body to json
			data, err := json.Marshal(tc.body)
			require.NoError(t, err)

			// start test server and send request
			server := NewServer(store)
			recorder := httptest.NewRecorder()

			url := "/tasks"
			request, err := http.NewRequest(http.MethodDelete, url, bytes.NewReader(data))
			require.NoError(t, err)

			server.router.ServeHTTP(recorder, request)

			tc.checkResponse(t, recorder)
		})
	}
}

func randomTask(userID int64) db.Task {
	return db.Task{
		ID:          util.RandomID(),
		UserID:      userID,
		Title:       util.RandomTitle(),
		Description: util.RandomPassword(100),
	}
}

func requireBodyMatchTask(t *testing.T, body *bytes.Buffer, task db.Task) {
	data, err := io.ReadAll(body)
	require.NoError(t, err)

	var gotTask db.Task
	err = json.Unmarshal(data, &gotTask)
	require.NoError(t, err)
	require.Equal(t, task, gotTask)
}

func requireBodyMatchTasks(t *testing.T, body *bytes.Buffer, tasks []db.Task) {
	data, err := io.ReadAll(body)
	require.NoError(t, err)

	var gotTasks []db.Task
	err = json.Unmarshal(data, &gotTasks)
	require.NoError(t, err)
	require.Equal(t, tasks, gotTasks)
}
