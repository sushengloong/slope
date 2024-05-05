package conversation

import (
	"context"
	"fmt"
	"regexp"
	"time"

	"github.com/go-playground/validator/v10"
	"go.jetify.com/typeid"
)

type Channel string

const (
	ChannelWeb   Channel = "web"
	ChannelEmail Channel = "email"
)

type Status string

const (
	StatusObserving Status = "observing"
	StatusActive    Status = "active"
	StatusCancelled Status = "cancelled"
	StatusFinished  Status = "finished"
	StatusFailed    Status = "failed"
)

type Conversation struct {
	ID         string            `json:"id"`
	CustomerID string            `json:"customer_id"`
	Channel    Channel           `json:"channel"`
	Metadata   map[string]string `json:"metadata"`
	Created    time.Time         `json:"created"`
	Updated    time.Time         `json:"updated"`
	Status     Status            `json:"status"`
}

type ListResponse struct {
	// Sites is the list of monitored sites.
	Data []*Conversation `json:"data"`
}

type StartParams struct {
	CustomerID string            `json:"customer_id" validate:"required,customerId"`
	Channel    string            `json:"channel" validate:"required,oneof=web email"`
	Metadata   map[string]string `json:"metadata" validate:"min=0,max=20"`
}

func (p *StartParams) Validate() error {
	validate := validator.New(validator.WithRequiredStructEnabled())
	validate.RegisterValidation("customerId", validateCustomerId)
	// TODO: more user friendly error message
	return validate.Struct(p)
}

func validateCustomerId(fl validator.FieldLevel) bool {
	re, err := regexp.Compile("^[a-zA-Z0-9\\-_+=]+$")
	if err != nil {
		panic(err)
	}
	fval := fl.Field().String()
	return re.MatchString(fval)
}

var conversations = map[string]*Conversation{}

//encore:api public method=GET path=/conversations
func List(ctx context.Context) (*ListResponse, error) {
	convs := getConversationValues(conversations)
	return &ListResponse{Data: convs}, nil
}

//encore:api public method=GET path=/conversations/:id
func Get(ctx context.Context, id string) (*Conversation, error) {
	conv, ok := conversations[id]
	if !ok {
		return nil, fmt.Errorf("No conversation found with id: %s", id)
	}
	return conv, nil
}

//encore:api public method=POST path=/conversations
func Start(ctx context.Context, p *StartParams) (*Conversation, error) {
	cid, err := typeid.WithPrefix("conversation")
	if err != nil {
		return nil, err
	}
	var metadata map[string]string
	if p.Metadata == nil {
		metadata = make(map[string]string)
	} else {
		metadata = p.Metadata
	}
	conv := Conversation{
		ID:         cid.String(),
		CustomerID: p.CustomerID,
		Channel:    Channel(p.Channel),
		Metadata:   metadata,
		Created:    time.Now(),
		Updated:    time.Now(),
		Status:     StatusActive,
	}
	conversations[conv.ID] = &conv
	return &conv, nil
}

func getConversationValues(m map[string]*Conversation) []*Conversation {
	values := make([]*Conversation, 0, len(m))
	for _, value := range m {
		values = append(values, value)
	}
	return values
}
