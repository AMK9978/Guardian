package entities

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Group represents a group of users.
type Group struct {
	ID     primitive.ObjectID `bson:"_id,omitempty"`
	Name   string             `json:"name"`
	Status int                `json:"status"`
	Users  []User             `json:"users"`
}

// User represents a user of the system.
type User struct {
	ID       primitive.ObjectID `bson:"_id,omitempty"`
	Name     string             `json:"name"`
	Password string             `json:"-"`
	Status   int                `json:"status"`
	Groups   []Group            `json:"groups"`
}

// AIModel represents an AI model's metadata.
type AIModel struct {
	Name    string `json:"name"`
	Address string `json:"address"`
	Status  int    `json:"status"`
}

// RefereeModel represents a model used by referees.
type RefereeModel struct {
	ModelID primitive.ObjectID `json:"model_id"`
	Token   string             `json:"token,omitempty"`
}

// TargetModel represents the target model for processing.
type TargetModel struct {
	ModelID primitive.ObjectID `json:"model_id"`
	Token   string             `json:"token"`
}

// Usage records token consumption for users.
type Usage struct {
	UserID                 primitive.ObjectID `json:"user_id"`
	TargetModelID          primitive.ObjectID `json:"target_model_id"`
	InputTokenConsumption  int                `json:"input_token_consumption"`
	OutputTokenConsumption int                `json:"output_token_consumption"`
}

// Task represents a task that can be used in the pipeline.
type Task struct {
	Type string `json:"type"`
}

// Pipeline contains the related processing tasks for each
type Pipeline struct {
	UserTasks  *[]UserTask  `json:"user_tasks,omitempty"`
	GroupTasks *[]GroupTask `json:"group_tasks,omitempty"`
}

// UserTask links a user to a task.
type UserTask struct {
	UserID primitive.ObjectID `json:"user_id"`
	Task   Task               `json:"task"`
}

// GroupTask links a group to a task.
type GroupTask struct {
	GroupID primitive.ObjectID `json:"group_id"`
	Task    Task               `json:"task"`
}

// TaskResult represents the result of task in the task pipeline
type TaskResult struct {
	TaskType string
	Success  bool
	Err      error
}
