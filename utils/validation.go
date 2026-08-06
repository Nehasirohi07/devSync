package utils

import (
	"errors"
	"regexp"

	"github.com/Nehasirohi07/devSync/models"
)

var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)

func ValidateUser(user models.User) error {

	if user.Name == "" {
		return errors.New("Name is required")
	}

	if len(user.Name) > 100 {
		return errors.New("Name must be less than 100 characters")
	}

	if user.Email == "" {
		return errors.New("Email is required")
	}

	if len(user.Email) > 150 {
		return errors.New("Email must be less than 150 characters")
	}

	if !emailRegex.MatchString(user.Email) {
		return errors.New("Invalid email format")
	}

	if user.Password == "" {
		return errors.New("Password is required")
	}

	if len(user.Password) < 6 {
		return errors.New("Password must bhi contain 6 characters")
	}

	if len(user.Password) > 100 {
		return errors.New("Password must be 100 characters or less")
	}

	return nil
}

func ValidateProject(project models.Project) error {

	if project.Name == "" {
		return errors.New("Project name is required")
	}

	if len(project.Name) > 100 {
		return errors.New("project name must be 100 characters or less")
	}

	if project.ManagerID <= 0 {
		return errors.New("invalid manager id")
	}

	return nil
}

func ValidateTask(task models.Task) error {

	if task.ProjectID <= 0 {
		return errors.New("Invalid project ID")
	}

	if task.AssignedTo <= 0 {
		return errors.New("Invalid assigned user ID")
	}

	if task.Title == "" {
		return errors.New("Task title is required")
	}

	if len(task.Title) > 100 {
		return errors.New("Task title must be 100 characters or less")
	}

	if task.Status != "pending" && task.Status != "in_progress" && task.Status != "completed" {
		return errors.New("Invalid task status")
	}

	if task.Progress < 0 || task.Progress > 100 {
		return errors.New("Progress must be between 0 and 100 ")
	}

	return nil
}

func ValidateComment(comment models.Comment) error {

	if comment.TaskID <= 0 {
		return errors.New("Invalid task ID")
	}

	if comment.UserID <= 0 {
		return errors.New("Invalid user ID")
	}

	if comment.Content == "" {
		return errors.New("Comment content is required")
	}

	return nil
}

func ValidateActivity(activity models.Activity) error {

	if activity.UserID <= 0 {
		return errors.New("Invalid user ID")
	}

	if activity.TaskID <= 0 {
		return errors.New("Invalid task ID")
	}

	if activity.Action == "" {
		return errors.New("Activity action is required")
	}

	if len(activity.Action) > 100 {
		return errors.New("Activity action must be 100 characters or less")
	}

	return nil
}
