package badactor

import (
	"context"
	"sync"
	"time"

	"github.com/sleekservices/ServiceRenderer/common/errors"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const pardonBuffer = -2 * time.Second

type Studio struct {
	guard    *sync.RWMutex
	rules    map[string]*Rule
	badactor BadactorService
}

func NewStudio(ctx context.Context, enabled bool, badactor BadactorService) *Studio {
	studio := &Studio{
		guard:    &sync.RWMutex{},
		rules:    make(map[string]*Rule),
		badactor: badactor,
	}

	if enabled {
		studio.AddRule(&Rule{
			Name:        "Login",
			Message:     "You have failed to login too many times",
			StrikeLimit: 5,
			ExpireBase:  time.Second * 30,
			Sentence:    time.Minute * 5,
		})
	}
	return studio
}

func (st *Studio) AddRule(r *Rule) {
	st.guard.Lock()
	defer st.guard.Unlock()

	st.rules[r.Name] = r
}

// Infraction accepts an ActorName and RuleName creates an infraction record and determines if the actor should be jailed
func (st *Studio) Infraction(c context.Context, actorName string, ruleName string) error {
	st.guard.RLock()
	defer st.guard.RUnlock()

	rule, ok := st.rules[ruleName]
	if !ok {
		return errors.ErrorLog(errors.ErrBadActor, "No rule found")
	}

	actor := Actor(actorName)
	err := st.createInfraction(c, actor, rule)
	if err != nil {
		return err
	}

	count, err := st.Strikes(c, actorName, ruleName)
	if err != nil {
		return err
	}

	if int(count) < rule.StrikeLimit {
		return nil
	}

	err = st.jail(c, actor, rule)
	return err
}

// createInfraction takes an Actor and Rule and creates an Infraction
func (st *Studio) createInfraction(ctx context.Context, actor Actor, rule *Rule) error {
	releaseTime := time.Now().Add(rule.ExpireBase)
	infraction := &Infraction{
		Actor:    actor,
		Rule:     rule,
		ExpireBy: releaseTime,
	}
	infraction.ID = primitive.NewObjectID()
	_, err := st.badactor.CreateInfraction(ctx, infraction)
	return err
}

// Strikes accepts an ActorName and a RuleName and returns the total infractions an Actor holds for a particular Rule that hasn't expired
func (st *Studio) Strikes(ctx context.Context, actorName string, ruleName string) (int64, error) {
	expired := time.Now()
	return st.badactor.CountInfraction(ctx, actorName, ruleName, expired)
}

// IsJailedFor accepts an ActorName and a RuleName and returns a bool if the Actor is Jailed for that particular Rule
func (st *Studio) IsJailedFor(ctx context.Context, actorName string, ruleName string) bool {
	term := time.Now()
	_, err := st.badactor.FindJail(ctx, actorName, ruleName, term)
	if err != nil {
		return err.Error() != errors.ErrorLog(errors.ErrNoDocumentFound).Error()
	}
	return true
}

// Pardon accepts an ActorName and a RuleName and sets the release date of the sentence to now
func (st *Studio) Pardon(ctx context.Context, actorName string, ruleName string) error {
	term := time.Now().Add(pardonBuffer)
	jail, err := st.badactor.FindJail(ctx, actorName, ruleName, term)
	if err != nil {
		return err
	}
	// set release date to now
	jail.ReleaseBy = term
	err = st.badactor.UpdateJail(ctx, jail)
	return err
}

// jail the actor if the Limit has been reached
func (st *Studio) jail(ctx context.Context, actor Actor, rule *Rule) error {
	jail := newJail(actor, rule)
	err := st.badactor.CreateJail(ctx, jail)

	return err
}
