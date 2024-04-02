package specification

import (
	"errors"
	"fmt"
	"time"

	pb "github.com/msmkdenis/yap-infokeeper/internal/text_data/api/grpchandlers/proto"
)

type TextDataSpecification struct {
	OwnerID       string
	Data          string
	Metadata      string
	CreatedAfter  time.Time
	CreatedBefore time.Time
}

func NewTextDataSpecification(ownerID string, in *pb.GetTextDataRequest) (*TextDataSpecification, error) {
	spec := &TextDataSpecification{
		OwnerID:  ownerID,
		Data:     in.Data,
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

func (t *TextDataSpecification) GetQueryArgs(query string) (string, []interface{}) {
	var args []interface{}
	whereCondition := make([]string, 0)
	query += " where "
	whereCondition = append(whereCondition, "owner_id = $")
	args = append(args, t.OwnerID)
	if t.Data != "" {
		whereCondition = append(whereCondition, "data ilike $")
		args = append(args, fmt.Sprintf("%%%s%%", t.Data))
	}
	if t.Metadata != "" {
		whereCondition = append(whereCondition, "metadata ilike $")
		args = append(args, fmt.Sprintf("%%%s%%", t.Metadata))
	}
	if !t.CreatedAfter.IsZero() {
		whereCondition = append(whereCondition, "created_at >= $")
		args = append(args, t.CreatedAfter)
	}
	if !t.CreatedBefore.IsZero() {
		whereCondition = append(whereCondition, "created_at <= $")
		args = append(args, t.CreatedBefore)
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
