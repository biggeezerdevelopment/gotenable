package tio

import (
	"context"
	"encoding/json"
	"strconv"
	"time"

	"github.com/tenable/gotenable/pkg/base"
)

// AuditLogAPI handles audit log operations.
type AuditLogAPI struct {
	client *Client
}

// AuditEvent represents an audit log event.
type AuditEvent struct {
	ID          string    `json:"id"`
	Action      string    `json:"action"`
	Crud        string    `json:"crud"`
	Description string    `json:"description"`
	Fields      []AuditField `json:"fields,omitempty"`
	IsAnonymous bool      `json:"is_anonymous"`
	IsFailure   bool      `json:"is_failure"`
	Received    time.Time `json:"received"`
	Target      AuditTarget `json:"target"`
	Actor       AuditActor  `json:"actor"`
}

// AuditField represents a field in an audit event.
type AuditField struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

// AuditTarget represents the target of an audit event.
type AuditTarget struct {
	ID   string `json:"id"`
	Type string `json:"type"`
	Name string `json:"name"`
}

// AuditActor represents the actor of an audit event.
type AuditActor struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// AuditLogOptions represents options for querying the audit log.
type AuditLogOptions struct {
	FromDate *time.Time
	ToDate   *time.Time
	Actor    string
	Target   string
	Action   string
}

// Events retrieves audit log events.
func (a *AuditLogAPI) Events(ctx context.Context, opts *AuditLogOptions) *base.Iterator[AuditEvent] {
	fetcher := func(ctx context.Context, offset, limit int) (json.RawMessage, *base.PaginationInfo, error) {
		params := map[string]string{
			"limit":  strconv.Itoa(limit),
			"offset": strconv.Itoa(offset),
		}

		if opts != nil {
			if opts.FromDate != nil {
				params["f"] = "date.gt:" + opts.FromDate.Format(time.RFC3339)
			}
			if opts.ToDate != nil {
				if _, ok := params["f"]; ok {
					params["f"] += ",date.lt:" + opts.ToDate.Format(time.RFC3339)
				} else {
					params["f"] = "date.lt:" + opts.ToDate.Format(time.RFC3339)
				}
			}
			if opts.Actor != "" {
				if _, ok := params["f"]; ok {
					params["f"] += ",actor.id.match:" + opts.Actor
				} else {
					params["f"] = "actor.id.match:" + opts.Actor
				}
			}
			if opts.Target != "" {
				if _, ok := params["f"]; ok {
					params["f"] += ",target.id.match:" + opts.Target
				} else {
					params["f"] = "target.id.match:" + opts.Target
				}
			}
			if opts.Action != "" {
				if _, ok := params["f"]; ok {
					params["f"] += ",action.match:" + opts.Action
				} else {
					params["f"] = "action.match:" + opts.Action
				}
			}
		}

		var result struct {
			Events     []AuditEvent `json:"events"`
			Pagination struct {
				Total  int `json:"total"`
				Limit  int `json:"limit"`
				Offset int `json:"offset"`
			} `json:"pagination"`
		}

		_, err := a.client.GetWithParams(ctx, "audit-log/v1/events", params, &result)
		if err != nil {
			return nil, nil, err
		}

		data, _ := json.Marshal(result.Events)
		return data, &base.PaginationInfo{
			Total:  result.Pagination.Total,
			Limit:  result.Pagination.Limit,
			Offset: result.Pagination.Offset,
		}, nil
	}

	transformer := func(data json.RawMessage) ([]AuditEvent, error) {
		var items []AuditEvent
		err := json.Unmarshal(data, &items)
		return items, err
	}

	return base.NewIterator(ctx, fetcher, transformer)
}

