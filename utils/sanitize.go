package utils

import (
	"strings"

	"github.com/Nehasirohi07/devSync/models"
)

func SanitizeUser(user models.User) models.User {

	user.Name = strings.TrimSpace(user.Name)

	user.Email = strings.TrimSpace(user.Email)

	user.Email = strings.ToLower(user.Email)

	return user
}

func SanitizeProject(project models.Project) models.Project {

	project.Name = strings.TrimSpace(project.Name)

	return project
}

func SanitizeTask(task models.Task) models.Task {

	task.Title = strings.TrimSpace(task.Title)

	task.Description = strings.TrimSpace(task.Description)

	task.Status = strings.TrimSpace(task.Status)

	return task
}

func SanitizeComment(comment models.Comment) models.Comment {

	comment.Content = strings.TrimSpace(comment.Content)

	return comment
}

func SanitizeActivity(activity models.Activity) models.Activity {

	activity.Action = strings.TrimSpace(activity.Action)

	activity.Details = strings.TrimSpace(activity.Details)

	return activity
}
