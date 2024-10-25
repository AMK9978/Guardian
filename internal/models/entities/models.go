package entities

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Group represents a group of users.
type Group struct {
	ID     primitive.ObjectID `bson:"_id,omitempty"`
	Name   string             `json:"name"`
	Status int                `json:"status"`
	Tasks  *[]Task            `json:"tasks,omitempty"`
}

type GroupMembers struct {
	ID      primitive.ObjectID `bson:"_id,omitempty"`
	UserID  primitive.ObjectID `bson:"_id"`
	GroupID primitive.ObjectID `bson:"_id"`
}

// User represents a user of the system.
type User struct {
	ID       primitive.ObjectID `bson:"_id,omitempty"`
	Name     string             `json:"name"`
	Password string             `json:"-"`
	Status   int                `json:"status"`
	Groups   []Group            `json:"groups"`
	Tasks    []primitive.ObjectID           `json:"tasks,omitempty"`
}

// RefereeModel represents a model used by referees.
type RefereeModel struct {
	ID      primitive.ObjectID `bson:"_id,omitempty"`
	Name    string             `json:"name"`
	Address string             `json:"address"`
	Status  int                `json:"status"`
	ModelID primitive.ObjectID `json:"model_id"`
	Token   string             `json:"token,omitempty"`
}

// TargetModel represents the target model for processing.
type TargetModel struct {
	ID       primitive.ObjectID `bson:"_id,omitempty"`
	Provider string             `json:"provider"`
	Name     string             `json:"name"`
	Address  string             `json:"address"`
	Status   int                `json:"status"`
	Token    string             `json:"token"`
}

// Usage records token consumption for users.
type Usage struct {
	ID                     primitive.ObjectID `json:"_id"`
	UserID                 primitive.ObjectID `json:"user_id"`
	TargetModelID          primitive.ObjectID `json:"target_model_id"`
	InputTokenConsumption  int                `json:"input_token_consumption"`
	OutputTokenConsumption int                `json:"output_token_consumption"`
}

// Task represents a task that can be used in the pipeline.
type Task struct {
	ID     primitive.ObjectID `json:"_id"`
	Type   string             `json:"type"`
	Status int                `json:"status"`
}

// TaskResult represents the result of task in the task pipeline
type TaskResult struct {
	TaskType string
	Success  bool
	Err      error
}
