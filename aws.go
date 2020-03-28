package main

import (
	"errors"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/ec2instanceconnect"
	"github.com/davecgh/go-spew/spew"
	log "github.com/sirupsen/logrus"
)

// AwsSession An AWS API session
type AwsSession struct {
	session *session.Session
	region  string
}

// NewAwsSession A new AWS Session
func NewAwsSession(region string) (*AwsSession, error) {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(region),
	})
	if err != nil {
		return nil, err
	}

	return &AwsSession{
		session: sess,
		region:  region,
	}, nil
}

func _getInstances(sess *AwsSession, instInput *ec2.DescribeInstancesInput) (*ec2.DescribeInstancesOutput, error) {
	svc := ec2.New(sess.session)
	instanceData, err := svc.DescribeInstances(instInput)
	if err != nil {
		return nil, err
	}

	return instanceData, nil
}

// Lookup the instance ID by using the instance's Name tag
func getInstanceIDFromNameTag(sess *AwsSession, name string) (string, error) {

	input := &ec2.DescribeInstancesInput{
		Filters: []*ec2.Filter{
			{
				Name: aws.String("tag:Name"),
				Values: []*string{
					aws.String(name),
				},
			},
			{
				Name: aws.String("instance-state-name"),
				Values: []*string{
					aws.String("running"),
				},
			},
		}}

	instanceData, err := _getInstances(sess, input)

	if err != nil {
		log.Fatal("agggggg")
	}

	if len(instanceData.Reservations) > 0 {
		return *instanceData.Reservations[0].Instances[0].InstanceId, nil
	}

	if instanceData.Reservations == nil {
		*(input.Filters[0].Values[0]) = *(input.Filters[0].Values[0]) + "*"
		spew.Dump(input)
		instanceData, err = _getInstances(sess, input)

		if err != nil {
			log.Fatal("ttttttt")
		}

		if len(instanceData.Reservations) > 0 {
			return *instanceData.Reservations[0].Instances[0].InstanceId, nil
		}

		os.Exit(1)

	}

	return "", errors.New("could not find instance")

}

// AwsInstance An AWS instance
type AwsInstance struct {
	session *AwsSession
	id      string
	data    *ec2.DescribeInstancesOutput
}

// NewAwsInstance returns a new AWS instance
func NewAwsInstance(sess *AwsSession, instanceID string) (*AwsInstance, error) {
	svc := ec2.New(sess.session)
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
		return "", errors.New("could not find private ip")
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
	svc := ec2instanceconnect.New(a.session.session)
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

	return nil
}
