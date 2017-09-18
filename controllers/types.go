package controllers

type HarborProjects struct {
	Project_id    float64 `json:project_id`
	Name          string  `json:name`
	Repo_count    float64 `json:repo_count`
	Creation_time string  `json:creation_time`
	Update_time   string  `json:update_time`
}

// owner_id is float64 1
// creation_time_str is string
// deleted is float64 0
// owner_name is string
// update_time is string 2017-08-14T02:51:14Z
// project_id is float64 9
// creation_time is string 2017-08-14T02:51:14Z
// public is float64 1
// Togglable is of a type I don't know how to handle
// current_user_role_id is float64 0
// repo_count is float64 0
// name is string uat
