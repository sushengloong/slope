package conversation

import (
	"context"
	"regexp"
	"time"

	"encore.dev/storage/sqldb"
	"github.com/go-playground/validator/v10"
	"go.jetify.com/typeid"
	"gorm.io/datatypes"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
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
	ID         string                                `json:"id"`
	CustomerID string                                `json:"customer_id"`
	Channel    Channel                               `json:"channel"`
	Metadata   datatypes.JSONType[map[string]string] `json:"metadata"`
	Status     Status                                `json:"status"`
	Created    time.Time                             `json:"created"`
	Updated    time.Time                             `json:"updated"`
}

type Message struct {
	ID              string                                `json:"id"`
	ConversationID  string                                `json:"conversation_id"`
	Body            string                                `json:"body"`
	ParticipantId   string                                `json:"participant_id"`
	ParticipantType string                                `json:"participant_type"`
	Metadata        datatypes.JSONType[map[string]string] `json:"metadata"`
	Created         time.Time                             `json:"created"`
	Updated         time.Time                             `json:"updated"`
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

type AddMessageParams struct {
	Body            string            `json:"body"`
	ParticipantId   string            `json:"participant_id"`
	ParticipantType string            `json:"participant_type"`
	Metadata        map[string]string `json:"metadata" validate:"min=0,max=20"`
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

//encore:api public method=GET path=/conversations
func (s *Service) List(ctx context.Context) (*ListResponse, error) {
	var convs []*Conversation
	if err := s.db.Find(&convs).Error; err != nil {
		return nil, err
	}
	return &ListResponse{Data: convs}, nil
}

//encore:api public method=GET path=/conversations/:id
func (s *Service) Get(ctx context.Context, id string) (*Conversation, error) {
	var conv Conversation
	if err := s.db.Where("id = $1", id).First(&conv).Error; err != nil {
		return nil, err
	}
	return &conv, nil
}

//encore:api public method=POST path=/conversations
func (s *Service) Start(ctx context.Context, p *StartParams) (*Conversation, error) {
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
	conv := &Conversation{
		ID:         cid.String(),
		CustomerID: p.CustomerID,
		Channel:    Channel(p.Channel),
		Metadata:   datatypes.NewJSONType(metadata),
		Created:    time.Now(),
		Updated:    time.Now(),
		Status:     StatusActive,
	}
	if err := s.db.Create(conv).Error; err != nil {
		return nil, err
	}
	return conv, nil
}

//encore:api public method=POST path=/conversations/:conversationID/messages
func (s *Service) AddMessage(ctx context.Context, conversationID string, p *AddMessageParams) (*Message, error) {
	var conv Conversation
	if err := s.db.Where("id = $1", conversationID).First(&conv).Error; err != nil {
		return nil, err
	}
	mid, err := typeid.WithPrefix("message")
	if err != nil {
		return nil, err
	}
	var metadata map[string]string
	if p.Metadata == nil {
		metadata = make(map[string]string)
	} else {
		metadata = p.Metadata
	}
	msg := &Message{
		ID:              mid.String(),
		ConversationID:  conversationID,
		Body:            p.Body,
		ParticipantId:   p.ParticipantId,
		ParticipantType: p.ParticipantType,
		Metadata:        datatypes.NewJSONType(metadata),
		Created:         time.Now(),
		Updated:         time.Now(),
	}
	if err := s.db.Create(msg).Error; err != nil {
		return nil, err
	}
	return msg, nil
}

//encore:service
type Service struct {
	db *gorm.DB
}

func initService() (*Service, error) {
	db, err := gorm.Open(postgres.New(postgres.Config{
		Conn: db.Stdlib(),
	}))
	if err != nil {
		return nil, err
	}
	return &Service{db: db}, nil
}

var db = sqldb.NewDatabase("conversation", sqldb.DatabaseConfig{
	Migrations: "./migrations",
})
