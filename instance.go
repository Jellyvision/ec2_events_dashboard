package ec2_events_dashboard

import "github.com/aws/aws-sdk-go/service/ec2"

// Instance represents an EC2 instance
type Instance struct {
	Client      *ec2.EC2
	Status      *ec2.InstanceStatus
	Reservation *ec2.Reservation // TODO: change this to *ec2.Instance
}

// InstancesWithEvents returns a []*instance which are instances in regions of the specified []*ec2.EC2 with any scheduled events
func InstancesWithEvents(clients ...*ec2.EC2) ([]*Instance, error) {
	instancesChan := make(chan []*Instance)

	for _, client := range clients {
		go func(c *ec2.EC2) {
			// TODO: handle errors here, maybe retry? or return an error for the whole function?
			someInstances, _ := instancesWithEvents(c)
			instancesChan <- someInstances
		}(client)
	}

	allInstances := []*Instance{}
	for range clients {
		someInstances := <-instancesChan
		allInstances = append(allInstances, someInstances...)
	}

	return allInstances, nil
}

func instancesWithEvents(client *ec2.EC2) ([]*Instance, error) {
	describeStatusesResponse, err := client.DescribeInstanceStatus(nil)
	if err != nil {
		return nil, err
	}

	instanceIdsWithEvents := []*string{}
	instanceStatusesWithEvents := []*ec2.InstanceStatus{}
	statuses := describeStatusesResponse.InstanceStatuses
	for _, status := range statuses {
		if len(status.Events) != 0 {
			instanceIdsWithEvents = append(instanceIdsWithEvents, status.InstanceId)
			instanceStatusesWithEvents = append(instanceStatusesWithEvents, status)
		}
	}

	describeInstancesOutput, err := client.DescribeInstances(
		&ec2.DescribeInstancesInput{
			InstanceIds: instanceIdsWithEvents,
		},
	)
	if err != nil {
		return nil, err
	}

	return newInstances(client, instanceStatusesWithEvents, describeInstancesOutput.Reservations), nil
}

func newInstances(client *ec2.EC2, statuses []*ec2.InstanceStatus, reservations []*ec2.Reservation) []*Instance {
	instances := []*Instance{}
	for i := 0; i < len(statuses); i++ {
		instances = append(instances, &Instance{Client: client, Status: statuses[i], Reservation: reservations[i]})
	}
	return instances
}

func newInstance(client *ec2.EC2, status *ec2.InstanceStatus, reservation *ec2.Reservation) *Instance {
	return &Instance{Client: client, Status: status, Reservation: reservation}
}
