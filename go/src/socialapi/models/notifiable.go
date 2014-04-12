package models

import (
	"github.com/koding/bongo"
	"time"
)

var (
	NOTIFIER_LIMIT = 3
)

type Notifiable interface {
	// users that will be notified are fetched while creating notification
	GetNotifiedUsers() ([]int64, error)
	GetType() string
	GetTargetId() int64
	FetchActors() (*ActorContainer, error)
	SetTargetId(int64)
	SetListerId(int64)
}

type InteractionNotification struct {
	TargetId   int64
	Type       string
	ListerId   int64
	NotifierId int64
}

func (n *InteractionNotification) GetNotifiedUsers() ([]int64, error) {
	i := NewInteraction()
	i.MessageId = n.TargetId
	return i.FetchInteractorIds(&bongo.Pagination{})
}

func (n *InteractionNotification) GetType() string {
	return n.Type
}

func (n *InteractionNotification) GetTargetId() int64 {
	return n.TargetId
}

func (n *InteractionNotification) SetTargetId(targetId int64) {
	n.TargetId = targetId
}

func (n *InteractionNotification) FetchActors() (*ActorContainer, error) {
	var count int
	i := NewInteraction()
	p := &bongo.Pagination{
		Limit: NOTIFIER_LIMIT,
	}
	i.MessageId = n.TargetId

	actors, err := i.FetchInteractorIdsWithCount(p, &count)
	if err != nil {
		return nil, err
	}

	ac := NewActorContainer()
	ac.LatestActors = actors
	ac.Count = count

	return ac, nil
}

func (n *InteractionNotification) SetListerId(listerId int64) {
	n.ListerId = listerId
}

func NewInteractionNotification(notificationType string) *InteractionNotification {
	return &InteractionNotification{Type: notificationType}
}

type ReplyNotification struct {
	TargetId   int64
	ListerId   int64
	NotifierId int64
}

func (n *ReplyNotification) GetNotifiedUsers() ([]int64, error) {
	// fetch all repliers
	cm := NewChannelMessage()
	cm.Id = n.TargetId

	p := &bongo.Pagination{}
	replierIds, err := cm.FetchReplierIds(p, true, time.Time{})

	if err != nil {
		return nil, err
	}

	// regress notifier from notified users
	filteredRepliers := make([]int64, 0)
	for _, replierId := range replierIds {
		if replierId != n.NotifierId {
			filteredRepliers = append(filteredRepliers, replierId)
		}
	}

	return filteredRepliers, nil
}

func (n *ReplyNotification) GetType() string {
	return NotificationContent_TYPE_COMMENT
}

func (n *ReplyNotification) GetTargetId() int64 {
	return n.TargetId
}

func (n *ReplyNotification) SetTargetId(targetId int64) {
	n.TargetId = targetId
}

func (n *ReplyNotification) FetchActors() (*ActorContainer, error) {
	mr := NewMessageReply()
	mr.MessageId = n.TargetId

	// we are gonna fetch actors after notified users first reply
	if err := mr.FetchFirstAccountReply(n.ListerId); err != nil {
		return nil, err
	}

	cm := NewChannelMessage()
	cm.Id = n.TargetId
	cm.AccountId = n.ListerId

	p := &bongo.Pagination{
		Limit: NOTIFIER_LIMIT,
	}

	// for preparing Actor Container we need latest actors and total replier count
	var count int
	actors, err := cm.FetchReplierIdsWithCount(p, &count, mr.CreatedAt)
	if err != nil {
		return nil, err
	}

	ac := NewActorContainer()
	ac.LatestActors = actors
	ac.Count = count

	return ac, nil
}

func (n *ReplyNotification) SetListerId(listerId int64) {
	n.ListerId = listerId
}

func NewReplyNotification() *ReplyNotification {
	return &ReplyNotification{}
}
