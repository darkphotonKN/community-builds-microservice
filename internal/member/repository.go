package member

import (
	"fmt"

	"github.com/darkphotonKN/community-builds/internal/models"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type MemberRepository struct {
	DB *sqlx.DB
}

func NewMemberRepository(db *sqlx.DB) *MemberRepository {
	return &MemberRepository{
		DB: db,
	}
}

func (r *MemberRepository) Create(member models.Member) error {
	query := `INSERT INTO members (name, email, password) VALUES (:name, :email, :password)`

	_, err := r.DB.NamedExec(query, member)

	if err != nil {
		return err
	}

	return nil
}

func (r *MemberRepository) GetById(id uuid.UUID) (*models.Member, error) {
	query := `SELECT * FROM members WHERE members.id = $1`

	var member models.Member

	err := r.DB.Get(&member, query, id)

	if err != nil {
		return nil, err
	}

	// Remove password from the struct
	member.Password = ""

	return &member, nil
}

func (r *MemberRepository) GetMemberByEmail(email string) (*models.Member, error) {
	var member models.Member
	query := `SELECT * FROM members WHERE members.email = $1`

	fmt.Println("Querying member with email:", email)

	err := r.DB.Get(&member, query, email)
	fmt.Println("Error:", err)

	if err != nil {
		return nil, err
	}

	return &member, nil
}
