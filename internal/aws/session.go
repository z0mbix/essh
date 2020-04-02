package aws

import (
	"errors"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	log "github.com/sirupsen/logrus"
	"github.com/z0mbix/essh/internal/config"
)

// Session An AWS API session
type Session struct {
	session *session.Session
	region  string
}

// NewSession A new AWS Session
func NewSession(region string) (*Session, error) {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(region),
	})
	if err != nil {
		return nil, err
	}

	return &Session{
		session: sess,
		region:  region,
	}, nil
}

func (sess *Session) GetReservations() ([]*ec2.Reservation, error) {
	switch config.SearchMode {
	case config.SearchModeTag:
		log.Debugf("using Name tag %s to find instance id", config.SearchValue)
		return getInstanceFromNameTag(sess, config.SearchValue)

	case config.SearchModeInst:
		return getInstanceFromID(sess, config.SearchValue)

	case config.SearchModeMenu:
		return getInstanceFromNameTag(sess, "*")
	}

	return nil, errors.New("invalid search mode")
}
