package aws

import (
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/davecgh/go-spew/spew"
	log "github.com/sirupsen/logrus"
)

func _getInstances(sess *Session, instInput *ec2.DescribeInstancesInput) (*ec2.DescribeInstancesOutput, error) {
	svc := ec2.New(sess.session)
	instanceData, err := svc.DescribeInstances(instInput)
	return instanceData, err
}

func getInstanceFromID(sess *Session, id string) ([]*ec2.Reservation, error) {
	input := &ec2.DescribeInstancesInput{
		InstanceIds: []*string{
			aws.String(id),
		},
	}

	a, err := _getInstances(sess, input)
	return a.Reservations, err
}

// Lookup the instance ID by using the instance's Name tag
func getInstanceFromNameTag(sess *Session, name string) ([]*ec2.Reservation, error) {
	input := &ec2.DescribeInstancesInput{
		Filters: []*ec2.Filter{
			{
				Name: aws.String("tag:Name"),
				Values: []*string{
					aws.String(name + "*"),
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
		return nil, fmt.Errorf("failed to search for tag: %s, err:%s", name, err)
	}

	if instanceData.Reservations == nil {
		*(input.Filters[0].Values[0]) = *(input.Filters[0].Values[0]) + "*"

		log.Debug(spew.Sdump(input))

		instanceData, err = _getInstances(sess, input)

		if err != nil {
			return nil, fmt.Errorf("failed to search for tag: %s, err:%s", name, err)
		}
	}

	if len(instanceData.Reservations) > 0 {
		return instanceData.Reservations, nil
	}

	return nil, errors.New("could not find instance")

}

func getTagValue(inst *ec2.Instance) string {
	for _, t := range inst.Tags {
		if *t.Key == "Name" {
			return *t.Value
		}
	}
	return ""
}
