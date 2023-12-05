package db

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/DMV-Nicolas/DevoraTasks/util"
	"github.com/stretchr/testify/require"
)

func createRandomTask(t *testing.T, name string) Task {
	user := createRandomUser(t)
	arg := CreateTaskParams{
		UserID:      user.ID,
		Title:       util.RandomTitle(),
		Description: name,
	}

	task, err := testQueries.CreateTask(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, task)

	require.Equal(t, arg.UserID, task.UserID)
	require.Equal(t, arg.Title, task.Title)
	require.Equal(t, arg.Description, task.Description)

	require.NotZero(t, task.CreatedAt)

	return task
}

func TestCreateTask(t *testing.T) {
	createRandomTask(t, "Create Task")
}

func TestGetTask(t *testing.T) {
	task1 := createRandomTask(t, "Get Task")
	task2, err := testQueries.GetTask(context.Background(), task1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, task2)

	require.Equal(t, task1.UserID, task2.UserID)
	require.Equal(t, task1.Title, task2.Title)
	require.Equal(t, task1.Description, task2.Description)

	require.WithinDuration(t, task1.CreatedAt, task2.CreatedAt, time.Second)
}

func TestListTasks(t *testing.T) {
	for i := 0; i < 5; i++ {
		createRandomTask(t, "List task")
	}

	arg := ListTasksParams{
		Limit:  5,
		Offset: 0,
	}

	tasks, err := testQueries.ListTasks(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, tasks, 5)

	for _, task := range tasks {
		require.NotEmpty(t, task)
	}
}

func TestUpdateTask(t *testing.T) {
	task1 := createRandomTask(t, "Update Task")
	arg := UpdateTaskParams{
		ID:          task1.ID,
		Title:       util.RandomTitle() + ":)",
		Description: "Updated task",
		Done:        true,
	}

	task2, err := testQueries.UpdateTask(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, task2)

	require.Equal(t, task1.ID, task2.ID)
	require.Equal(t, arg.Title, task2.Title)
	require.Equal(t, arg.Description, task2.Description)
	require.Equal(t, arg.Done, task2.Done)

	require.WithinDuration(t, task1.CreatedAt, task2.CreatedAt, time.Second)
}

func TestDeleteTask(t *testing.T) {
	task1 := createRandomTask(t, "Delete task")
	err := testQueries.DeleteTask(context.Background(), task1.ID)
	require.NoError(t, err)

	task2, err := testQueries.GetTask(context.Background(), task1.ID)
	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, task2)
}
