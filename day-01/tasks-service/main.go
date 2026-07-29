package main

import (
	"errors"
	"fmt"
	"time"
)

type Task struct {
	ID        int64
	UserID    int64
	Title     string
	Completed bool
	CreatedAt time.Time
}

type TaskRepository interface{
	CreateTask(task Task) (Task, error)
	GetTaskByID(taskID int64) (Task, error)
	UpdateTask(task Task) (Task, error)
}



type TaskService struct {
	repo TaskRepository
}

func NewTaskService(taskRepo TaskRepository) *TaskService{
	return &TaskService{
		repo: taskRepo,
	}
}

type MemoryTasksRepository struct{
	tasks map[int64]Task
}


var (
	ErrInvalidUserID = errors.New("invalid user ID")
	ErrEmptyTitle = errors.New("title is empty")
	ErrCannotFindTask = errors.New("cannot find task")
	ErrInvalidTaskID = errors.New("invalid task id")
	ErrAccessDenied    = errors.New("access denied")
)

func (r *MemoryTasksRepository) CreateTask(task Task) (Task, error) {
	newID := int64(len(r.tasks)+1)
	task.ID = newID
	task.CreatedAt = time.Now()
	r.tasks[newID] =  task

	return task, nil
}

func (r *MemoryTasksRepository) GetTaskByID(taskID int64) (Task, error){
	task, ok := r.tasks[taskID]
	if !ok {
		return Task{}, ErrCannotFindTask
	}
	return task, nil
}

func (r *MemoryTasksRepository) UpdateTask(task Task) (Task, error) {
	r.tasks[task.ID] = task
	
	return task, nil
}

func (s *TaskService) CreateTask(userID int64, title string) (Task, error) {
	if userID <= 0 {
		return Task{}, ErrInvalidUserID
	}
	if title == "" {
		return Task{}, ErrEmptyTitle
	}

	task:= Task{
		UserID: userID,
		Title: title,
		Completed: false,
	}

	createdTask, err := s.repo.CreateTask(task)
	if err != nil {
		return Task{}, err
	}

	return createdTask, nil
}

func (s *TaskService) GetTaskByID(taskID int64, userID int64) (Task, error) {
	if taskID <= 0 {
		return Task{}, ErrInvalidTaskID
	}
	if userID <= 0 {
		return Task{}, ErrInvalidUserID
	}

	task, err := s.repo.GetTaskByID(taskID) 
	if err != nil {
		return Task{}, err
	}

	if task.UserID != userID {
		return Task{}, ErrAccessDenied
	}

	return task, nil
}


func( s *TaskService) CompleteTask(taskID int64, userID int64) (Task, error) {
	if taskID <= 0 {
		return Task{}, ErrInvalidTaskID
	}
	if userID <= 0 {
		return Task{}, ErrInvalidUserID
	}
	
	task, err := s.repo.GetTaskByID(taskID)
	if err != nil {
		return Task{}, err
	}
	if task.UserID != userID {
		return Task{}, ErrAccessDenied
	}

	task.Completed = true

	updatedTask, err := s.repo.UpdateTask(task)
	if err != nil {
		return Task{}, err
	}

	return updatedTask, nil
}



func main() {
	repo := &MemoryTasksRepository{
		tasks: make(map[int64]Task),
	}
	
	service := NewTaskService(repo)

	newTask, err := service.CreateTask(1, "make breakfast")
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	fmt.Println("task created:", newTask)

	foundTask, err := service.GetTaskByID(1, 1)
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	fmt.Println("found task:", foundTask)

	updateTask, err := service.CompleteTask(newTask.ID, newTask.UserID)
	if err != nil {
		fmt.Println("error", err)
		return
	}
	fmt.Println("updated task:", updateTask)
	
}