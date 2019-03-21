package forms

import (
	cmsf "backend/forms/cmsforms"
	uf "backend/forms/userforms"
)

type (
	AssignmentAggQuery cmsf.AssignmentAgg

	CourseAggQuery cmsf.CourseAgg
	CourseAddUserForm cmsf.CourseAddUser
	CourseBulkAddUserForm cmsf.CourseBulkAddUser

	CreateAssignmentForm cmsf.CreateAssignment
	CreateCourseForm cmsf.CreateCourse

	UserLoginForm uf.LoginForm
	UserRegisterForm uf.RegisterForm
)
