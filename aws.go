package main

import (
	"errors"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/ec2instanceconnect"
	log "github.com/sirupsen/logrus"
)

// AwsInstance An AWS instance
type AwsInstance struct {
	session *session.Session
	id      string
	data    *ec2.DescribeInstancesOutput
}

// NewAwsInstance returns a new AWS instance
func NewAwsInstance(region, instanceID string) (*AwsInstance, error) {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(region),
	})
	if err != nil {
		return nil, err
	}

	svc := ec2.New(sess)
	input := &ec2.DescribeInstancesInput{
		InstanceIds: []*string{
			aws.String(instanceID),
		},
	}

	instanceData, err := svc.DescribeInstances(input)
	if err != nil {
		return nil, err
	}

	return &AwsInstance{
		session: sess,
		id:      instanceID,
		data:    instanceData,
	}, nil
}

func (a *AwsInstance) IP(public bool) (string, error) {
	if public {
		ip, err := a.publicIP()
		if err != nil {
			return "", err
		}
		return ip, nil
	}

	ip, err := a.privateIP()
	if err != nil {
		return "", err
	}
	return ip, nil
}

func (a *AwsInstance) privateIP() (string, error) {
	ip := a.data.Reservations[0].Instances[0].PrivateIpAddress
	if ip == nil {
		return "", errors.New("could not find public ip")
	}
	return *ip, nil
}

func (a *AwsInstance) publicIP() (string, error) {
	ip := a.data.Reservations[0].Instances[0].PublicIpAddress
	if ip == nil {
		return "", errors.New("could not find public ip")
	}
	return *ip, nil
}

func (a *AwsInstance) sendPublicKey(user, publicKey string) error {
	svc := ec2instanceconnect.New(a.session)
	input := &ec2instanceconnect.SendSSHPublicKeyInput{
		AvailabilityZone: aws.String(*a.data.Reservations[0].Instances[0].Placement.AvailabilityZone),
		InstanceId:       aws.String(a.id),
		InstanceOSUser:   aws.String(user),
		SSHPublicKey:     aws.String(publicKey),
	}

	_, err := svc.SendSSHPublicKey(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case ec2instanceconnect.ErrCodeAuthException:
				log.Errorln(ec2instanceconnect.ErrCodeAuthException, aerr.Error())
			case ec2instanceconnect.ErrCodeInvalidArgsException:
				log.Errorln(ec2instanceconnect.ErrCodeInvalidArgsException, aerr.Error())
			case ec2instanceconnect.ErrCodeServiceException:
				log.Errorln(ec2instanceconnect.ErrCodeServiceException, aerr.Error())
			case ec2instanceconnect.ErrCodeThrottlingException:
				log.Errorln(ec2instanceconnect.ErrCodeThrottlingException, aerr.Error())
			case ec2instanceconnect.ErrCodeEC2InstanceNotFoundException:
				log.Errorln(ec2instanceconnect.ErrCodeEC2InstanceNotFoundException, aerr.Error())
			default:
				log.Errorln(aerr.Error())
			}
		} else {
			log.Errorln(err.Error())
		}
		return err
	}

	// log.Debug(result)
	return nil
}
