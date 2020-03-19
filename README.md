# essh

SSH to EC2 instances using in memory, ephemeral ssh keys and EC2 instance connect


## Description

`essh` does the following:

- Generates a one time RSA ssh keypair in memory
- Adds the private key to you ssh agent (for a few seconds)
- Pushes the public key to the instance using [ec2-instance connect](https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/Connect-using-EC2-Instance-Connect.html)
- `ssh` to the instance using the private IP address (public IP can be used with `-p`), using user `ec2-user` by default


## Requirements

As `essh` uses AWS APIs, you will need you have valid credentials configured. If you're using this tool, then I'm presuming that you know how to do this, if not [see here](https://docs.aws.amazon.com/cli/latest/userguide/cli-chap-configure.html).

You should set the region with the `-r`/`--region` flag, or by setting the environment variable `AWS_DEFAULT_REGION`


## Usage

```
$ essh --help
Usage of essh:
  -d, --debug             Enable debug logging
  -r, --region string     AWS Region
  -p, --use-public-ip     Use the public ip instead of the private ip address
  -u, --username string   UNIX user name (default "ec2-user")
pflag: help requested
```

Connect to an instance's private IP with ssh as user `ec2-user` (the default):

```
$ essh i-02fab0d7dd3ab737b
```

Connect to an instance's public IP with ssh as user `ec2-user`:

```
$ essh -p i-02fab0d7dd3ab737b
```

Connect as user `ubuntu` passing the flags `-A`, `-4` and the command `uptime` to the ssh command:

```
$ essh -u ubuntu i-02fab0d7dd3ab737b -- -A -4 uptime
```

Connect to an instance by its name tag:

```
$ essh prod-web1
```


## Examples

Run with debug logging enabled:

```
$ essh -d i-03faf0d7dd3ab737a
DEBU[2020-03-16T22:49:31.384] aws region not set, trying AWS_DEFAULT_REGION environment variable
DEBU[2020-03-16T22:49:31.384] aws region found in AWS_DEFAULT_REGION environment variable: eu-west-1
DEBU[2020-03-16T22:49:31.680] looking up ip of instance id: i-03faf0d7dd3ab737a
DEBU[2020-03-16T22:49:31.863] adding key to agent
DEBU[2020-03-16T22:49:31.864] host: 10.200.3.25
DEBU[2020-03-16T22:49:31.864] pushing public key to instance
INFO[2020-03-16T22:49:32.728] running command: ssh -l ec2-user 10.200.3.25
Last login: Mon Mar 16 22:49:14 2020 from ip-10-200-42-219.eu-west-1.compute.internal

       __|  __|_  )
       _|  (     /   Amazon Linux 2 AMI
      ___|\___|___|

https://aws.amazon.com/amazon-linux-2/
No packages needed for security; 6 packages available
Run "sudo yum update" to apply all updates.
[ec2-user@ip-10-200-3-25 ~]$
```

Run with debug logging enabled, connect to the instance named "prod-web1" on its public ip address and run `uptime`:

```
$ essh -d -p prod-web1 -- uptime
DEBU[2020-03-19T16:42:38.920] aws region not set, trying AWS_DEFAULT_REGION environment variable
DEBU[2020-03-19T16:42:38.920] aws region found in AWS_DEFAULT_REGION environment variable: eu-west-1
DEBU[2020-03-19T16:42:38.920] using Name tag prod-web1 to find instance id
DEBU[2020-03-19T16:42:39.711] found instance id: i-0cc2be02456a7180c
DEBU[2020-03-19T16:42:39.770] looking up ip of: prod-web1
DEBU[2020-03-19T16:42:39.939] adding key to agent
DEBU[2020-03-19T16:42:39.940] host: 52.51.41.123
DEBU[2020-03-19T16:42:39.940] pushing public key to instance
INFO[2020-03-19T16:42:40.681] running command: ssh -l ec2-user 52.51.41.123 uptime
 16:42:42 up 16 min,  0 users,  load average: 0.13, 0.04, 0.01
```

## Build

```
go build
```

Put the resulting `essh` binary somewhere in your `$PATH`.


## TODO

- Exit with the ssh command exit code
- Create binary releases for popular platforms/architectures


## License

The project is open-source software licensed under the [MIT license](http://opensource.org/licenses/MIT).
