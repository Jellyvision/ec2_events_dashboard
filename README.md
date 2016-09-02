EC2 Events Dashboard
====================

A tool to show all of the scheduled events for your EC2 instances, in all AWS regions, from all of your AWS accounts, in one place

Usage
-----

This example shows using 3 different AWS accounts, but you can use as many or as few as you want

    $ git clone git@github.com:Jellyvision/ec2_events_dashboard.git
    $ cd ec2_events_dashboard/cmd/ec2_events_dashboard
    $ go run main.go \
        --creds $ACCOUNT_1_ACCESS_KEY_ID:$ACCOUNT_1_SECRET_ACCESS_KEY \
        --creds $ACCOUNT_2_ACCESS_KEY_ID:$ACCOUNT_2_SECRET_ACCESS_KEY \
        --creds $ACCOUNT_3_ACCESS_KEY_ID:$ACCOUNT_3_SECRET_ACCESS_KEY

Then go to `http://localhost:3000` in your browser and you'll see something like this:

![dashboard screenshot](https://github.com/Jellyvision/ec2_events_dashboard/blob/master/screenshot.png)
