package testsuite

import "github.com/google/uuid"

type TestSkills []TestSkill

type TestSkill struct {
	ID   uuid.UUID
	Name string
	Type string
}
