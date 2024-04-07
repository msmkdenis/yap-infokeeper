package specification

import (
	"errors"
	"fmt"
	"time"

	pb "github.com/msmkdenis/yap-infokeeper/internal/credit_card/api/grpchandlers/proto"
)

type CreditCardSpecification struct {
	OwnerID       string
	Number        string
	OwnerName     string
	CVVCode       string
	PinCode       string
	Metadata      string
	ExpiresAfter  time.Time
	ExpiresBefore time.Time
	CreatedAfter  time.Time
	CreatedBefore time.Time
}

func NewCreditCardSpecification(ownerID string, in *pb.GetCreditCardRequest) (*CreditCardSpecification, error) {
	spec := &CreditCardSpecification{
		OwnerID:   ownerID,
		Number:    in.Number,
		OwnerName: in.Owner,
		CVVCode:   in.CvvCode,
		PinCode:   in.PinCode,
		Metadata:  in.Metadata,
	}

	if in.ExpiresAfter != "" {
		expiresAfter, err := time.Parse("2006-01-02", in.ExpiresAfter)
		if err != nil {
			return nil, errors.New("expires after must be in format '2006-01-02'")
		}
		spec.ExpiresAfter = expiresAfter
	}

	if in.ExpiresBefore != "" {
		expiresBefore, err := time.Parse("2006-01-02", in.ExpiresBefore)
		if err != nil {
			return nil, errors.New("expires before must be in format '2006-01-02'")
		}
		spec.ExpiresBefore = expiresBefore
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

func (t *CreditCardSpecification) GetQueryArgs(query string) (string, []interface{}) {
	var args []interface{}
	whereCondition := make([]string, 0)
	query += " where "
	whereCondition = append(whereCondition, "owner_id = $")
	args = append(args, t.OwnerID)
	if t.Number != "" {
		whereCondition = append(whereCondition, "number ilike $")
		args = append(args, fmt.Sprintf("%%%s%%", t.Number))
	}
	if t.OwnerName != "" {
		whereCondition = append(whereCondition, "owner_name ilike $")
		args = append(args, fmt.Sprintf("%%%s%%", t.OwnerName))
	}
	if !t.ExpiresAfter.IsZero() {
		whereCondition = append(whereCondition, "expires_at >= $")
		args = append(args, t.ExpiresAfter)
	}
	if !t.ExpiresBefore.IsZero() {
		whereCondition = append(whereCondition, "expires_at <= $")
		args = append(args, t.ExpiresBefore)
	}
	if t.CVVCode != "" {
		whereCondition = append(whereCondition, "cvv_code ilike $")
		args = append(args, fmt.Sprintf("%%%s%%", t.Metadata))
	}
	if t.PinCode != "" {
		whereCondition = append(whereCondition, "pin_code ilike $")
		args = append(args, fmt.Sprintf("%%%s%%", t.PinCode))
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
