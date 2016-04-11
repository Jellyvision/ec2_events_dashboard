package ec2_events_dashboard

import "github.com/aws/aws-sdk-go/service/ec2"

// Instance represents an EC2 instance
type Instance struct {
	Client   *ec2.EC2
	Status   *ec2.InstanceStatus
	Instance *ec2.Instance
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

	instanceChan := make(chan *Instance)
	instancesCount := 0

	// TODO: This makes one http request for getting instances statuses, then one request per instace with events
	// it would be possible to just make one request with a slice of instance IDs, but would need to sort a []*ec2.Instance
	statuses := describeStatusesResponse.InstanceStatuses
	for _, status := range statuses {
		if len(status.Events) != 0 {
			instancesCount++

			go func(s *ec2.InstanceStatus) {
				describeInstancesOutput, _ := client.DescribeInstances(
					&ec2.DescribeInstancesInput{
						InstanceIds: []*string{s.InstanceId},
					},
				)

				instanceChan <- &Instance{Client: client, Status: s, Instance: describeInstancesOutput.Reservations[0].Instances[0]}
			}(status)
		}
	}

	instances := []*Instance{}
	for i := 0; i < instancesCount; i++ {
		instances = append(instances, <-instanceChan)
	}

	return instances, nil
}
