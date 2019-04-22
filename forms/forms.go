package forms

import (
	cmsf "backend/forms/cmsforms"
	uf "backend/forms/userforms"
)

type (
	AssignmentAggQuery cmsf.AssignmentAgg

	CourseAggQuery        cmsf.CourseAgg
	CourseAddUserForm     cmsf.CourseAddUser
	CourseBulkAddUserForm cmsf.CourseBulkAddUser

	CreateAssignmentPreForm  cmsf.CreateAssignmentPreParse
	CreateAssignmentPostForm cmsf.CreateAssignmentPostParse
	CreateCourseForm         cmsf.CreateCourse

	GradeAggQuery cmsf.GradeAgg

	UserLoginForm    uf.LoginForm
	UserRegisterForm uf.RegisterForm
)
