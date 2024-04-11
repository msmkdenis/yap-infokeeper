package specification

import (
	"errors"
	"fmt"
	"time"

	pb "github.com/msmkdenis/yap-infokeeper/internal/proto/credential"
)

type CredentialSpecification struct {
	OwnerID       string
	Login         string
	Password      string
	Metadata      string
	CreatedAfter  time.Time
	CreatedBefore time.Time
}

func NewCredentialSpecification(ownerID string, in *pb.GetCredentialRequest) (*CredentialSpecification, error) {
	spec := &CredentialSpecification{
		OwnerID:  ownerID,
		Login:    in.Login,
		Password: in.Password,
		Metadata: in.Metadata,
	}

	if in.CreatedAfter != "" {
		after, err := time.Parse("2006-01-02", in.CreatedAfter)
		if err != nil {
			return nil, errors.New("created after must be in format '2006-01-02'")
		}
		spec.CreatedAfter = after
	}

	if in.CreatedBefore != "" {
		before, err := time.Parse("2006-01-02", in.CreatedBefore)
		if err != nil {
			return nil, errors.New("created before must be in format '2006-01-02'")
		}
		spec.CreatedBefore = before
	}

	return spec, nil
}

func (c *CredentialSpecification) GetQueryArgs(query string) (string, []interface{}) {
	var args []interface{}
	whereCondition := make([]string, 0)
	query += " where "
	whereCondition = append(whereCondition, "owner_id = $")
	args = append(args, c.OwnerID)
	if c.Login != "" {
		whereCondition = append(whereCondition, "login ilike $")
		args = append(args, fmt.Sprintf("%%%s%%", c.Login))
	}
	if c.Password != "" {
		whereCondition = append(whereCondition, "password ilike $")
		args = append(args, fmt.Sprintf("%%%s%%", c.Password))
	}
	if c.Metadata != "" {
		whereCondition = append(whereCondition, "metadata ilike $")
		args = append(args, fmt.Sprintf("%%%s%%", c.Metadata))
	}
	if !c.CreatedAfter.IsZero() {
		whereCondition = append(whereCondition, "created_at >= $")
		args = append(args, c.CreatedAfter)
	}
	if !c.CreatedBefore.IsZero() {
		whereCondition = append(whereCondition, "created_at <= $")
		args = append(args, c.CreatedBefore)
	}

	var counter int
	for i, clause := range whereCondition {
		if i == 0 {
			counter++
			query = query + clause + fmt.Sprint(counter)
		} else {
			counter++
			query = query + " and " + clause + fmt.Sprint(counter)
		}
	}

	return query, args
}
