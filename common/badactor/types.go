package badactor


import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type BadactorService interface {
	CreateInfraction(ctx context.Context, infraction *Infraction) (*Infraction, error)
	CountInfraction(ctx context.Context, actorName, ruleName string, expireTime time.Time) (int64, error)
	CreateJail(ctx context.Context, jail *Jail) error
	FindJail(ctx context.Context, actorName, ruleName string, releaseTime time.Time) (*Jail, error)
	UpdateJail(ctx context.Context, jail *Jail) error
}

type BadStudio interface {
	IsJailedFor(c context.Context, actorName, ruleName string) bool
	Infraction(c context.Context, actorName, ruleName string) error
	Pardon(c context.Context, actorName, ruleName string) error
}

type Actor string

type Rule struct {
	Name        string        `bson:"name"`
	Message     string        `bson:"message"`
	StrikeLimit int           `bson:"strike_limit"`
	ExpireBase  time.Duration `bson:"expire_base"`
	Sentence    time.Duration `bson:"sentence"`
}

type Infraction struct {
	ID       primitive.ObjectID `bson:"_id"`
	Actor    Actor              `bson:"actor"`
	Rule     *Rule              `bson:"rule"`
	ExpireBy time.Time          `bson:"expire_by"`
}

type Jail struct {
	ID        primitive.ObjectID `bson:"_id"`
	Actor     Actor              `bson:"actor"`
	Rule      *Rule              `bson:"rule"`
	ReleaseBy time.Time          `bson:"release_by"`
	Start     time.Time          `bson:"start"`
}

func newJail(a Actor, r *Rule) *Jail {
	return &Jail{
		Actor:     a,
		Rule:      r,
		ReleaseBy: time.Now().Add(r.Sentence),
		Start:     time.Now(),
	}
}