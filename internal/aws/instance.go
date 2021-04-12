package aws

import (
	"errors"

	"github.com/apex/log"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/ec2instanceconnect"
)

// AwsInstance An AWS instance
type Instance struct {
	session *Session
	data    *ec2.Instance

	// extracted here for convenience
	ID        string
	Public    bool
	ConnectIP string
	NameTag   string
}

// NewInstance returns a new AWS instance
func NewInstance(sess *Session, inst *ec2.Instance, publicIP bool) (*Instance, error) {
	ai := Instance{
		session: sess,
		ID:      *inst.InstanceId,
		data:    inst,
		Public:  publicIP,
	}

	var err error
	ai.ConnectIP, err = ai.IP(publicIP)
	if err != nil {
		return nil, err
	}
	ai.NameTag = getTagValue(inst)
	return &ai, nil
}

// IP is the IP address of the instance
func (i *Instance) IP(public bool) (string, error) {
	if public {
		ip, err := i.publicIP()
		if err != nil {
			return "", err
		}
		return ip, nil
	}

	ip, err := i.privateIP()
	if err != nil {
		return "", err
	}
	return ip, nil
}

func (i *Instance) privateIP() (string, error) {
	ip := i.data.PrivateIpAddress
	if ip == nil {
		return "", errors.New("could not find private ip")
	}
	return *ip, nil
}

func (i *Instance) publicIP() (string, error) {
	ip := i.data.PublicIpAddress
	if ip == nil {
		return "", errors.New("could not find public ip")
	}
	return *ip, nil
}

func (i *Instance) SendPublicKey(user, publicKey string) error {
	log.Debug("pushing public key to instance")

	svc := ec2instanceconnect.New(i.session.session)
	input := &ec2instanceconnect.SendSSHPublicKeyInput{
		AvailabilityZone: aws.String(*i.data.Placement.AvailabilityZone),
		InstanceId:       aws.String(i.ID),
		InstanceOSUser:   aws.String(user),
		SSHPublicKey:     aws.String(publicKey),
	}

	_, err := svc.SendSSHPublicKey(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case ec2instanceconnect.ErrCodeAuthException:
				log.Errorf(ec2instanceconnect.ErrCodeAuthException, aerr.Error())
			case ec2instanceconnect.ErrCodeInvalidArgsException:
				log.Errorf(ec2instanceconnect.ErrCodeInvalidArgsException, aerr.Error())
			case ec2instanceconnect.ErrCodeServiceException:
				log.Errorf(ec2instanceconnect.ErrCodeServiceException, aerr.Error())
			case ec2instanceconnect.ErrCodeThrottlingException:
				log.Errorf(ec2instanceconnect.ErrCodeThrottlingException, aerr.Error())
			case ec2instanceconnect.ErrCodeEC2InstanceNotFoundException:
				log.Errorf(ec2instanceconnect.ErrCodeEC2InstanceNotFoundException, aerr.Error())
			default:
				log.Errorf(aerr.Error())
			}
		} else {
			log.Errorf(err.Error())
		}
		return err
	}

	return nil
}
