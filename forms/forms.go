package forms

import (
	cmsf "backend/forms/cmsforms"
	uf "backend/forms/userforms"
)

type (
	AssignmentAggQuery cmsf.AssignmentAgg

	CreateAssignmentForm cmsf.CreateAssignment
	CourseAggQuery cmsf.CourseAgg

	UserLoginForm uf.LoginForm
	UserRegisterForm uf.RegisterForm
)