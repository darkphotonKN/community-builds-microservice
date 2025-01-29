package member

import (
	"fmt"

	"github.com/darkphotonKN/community-builds/internal/models"
	"github.com/darkphotonKN/community-builds/internal/utils/errorutils"
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

	fmt.Println("Error when creating member:", err)

	if err != nil {
		return errorutils.AnalyzeDBErr(err)
	}

	return nil
}

func (r *MemberRepository) UpdatePassword(params MemberUpdatePasswordParams) error {
	query := `UPDATE members SET password = :password WHERE id = :id`

	result, err := r.DB.NamedExec(query, params)
	if err != nil {
		return errorutils.AnalyzeDBErr(err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check affected rows: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("no member found with id: %v", params.ID)
	}

	return nil
}

func (r *MemberRepository) UpdateInfo(params MemberUpdateInfoParams, userId uuid.UUID) error {
	query := `UPDATE members SET name = :name, status = :status WHERE id = :id`

	result, err := r.DB.NamedExec(query, params)

	fmt.Println("result", result)
	if err != nil {
		return errorutils.AnalyzeDBErr(err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check affected rows: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("no member found with id: %v", params.ID)
	}

	return nil
}

func (r *MemberRepository) GetByIdWithPassword(id uuid.UUID) (*models.Member, error) {
	query := `SELECT * FROM members WHERE members.id = $1`

	var member models.Member

	err := r.DB.Get(&member, query, id)

	if err != nil {
		return nil, err
	}

	return &member, nil
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

	err := r.DB.Get(&member, query, email)
	fmt.Println("Error:", err)

	if err != nil {
		return nil, err
	}

	return &member, nil
}

func (r *MemberRepository) CreateDefaultMembers(members []CreateDefaultMember) error {
	query := `
	INSERT INTO members(id, email, name, password, status)
	VALUES(:id, :email, :name, :password, :status)
	ON CONFLICT (id) DO NOTHING
	`
	_, err := r.DB.NamedExec(query, members)

	fmt.Printf("DEBUG: Error when creating default member: %s\n", err)

	if err != nil {
		return errorutils.AnalyzeDBErr(err)
	}

	return nil
}
