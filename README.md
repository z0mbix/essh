# essh

SSH to an EC2 instance using an in memory, ephemeral ssh key and EC2 instance connect to push the new public key to the instance


## Description

`essh` does the following:

- Generates a one time RSA ssh keypair in memory
- Adds the private key to you ssh agent (for a few seconds)
- Pushes the public key to the instance using [ec2-instance connect](https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/Connect-using-EC2-Instance-Connect.html)
- `ssh` to the instance using the private IP address (public IP can be used with `-p`), using user `ec2-user` by default


## Requirements

As `essh` uses AWS APIs, you will need you have valid credentials configured. If you're using this tool, then I'm presuming that you know how to do this, if not [see here](https://docs.aws.amazon.com/cli/latest/userguide/cli-chap-configure.html).

You should set the region with the `-r`/`--region` flag, or by setting the environment variable `AWS_DEFAULT_REGION` or `AWS_REGION`.


## Usage

```
$ essh --help
Usage of essh:
  -d, --debug             Enable debug logging
  -r, --region string     AWS Region
  -p, --use-public-ip     Use the public ip instead of the private ip address
  -u, --username string   UNIX user name (default "ec2-user")
  -v, --version           Show version
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

Connect to an instance by it's full name tag:

```
$ essh prod-web1
```

Display a menu of instances that match a partial tag:

```
$ essh gitlab
```

Display all running instances in a region:

```
$ essh
```


You can use `/` to search the list of instances.


## Examples

Connect to an instance on it's private IP:

```
$ essh i-03faf0d7dd3ab737a
running command: ssh -l ec2-user 10.200.3.25
Last login: Mon Mar 16 22:49:14 2020 from ip-10-200-42-219.eu-west-1.compute.internal

       __|  __|_  )
       _|  (     /   Amazon Linux 2 AMI
      ___|\___|___|

https://aws.amazon.com/amazon-linux-2/
No packages needed for security; 6 packages available
Run "sudo yum update" to apply all updates.
[ec2-user@ip-10-200-3-25 ~]$
```

Connect to the instance named "prod-web1" on its public ip address and run `uptime`:

```
$ essh -p prod-web1 -- uptime
running command: ssh -l ec2-user 52.51.41.123 uptime
 16:42:42 up 16 min,  0 users,  load average: 0.13, 0.04, 0.01
```

Connect to a host named "gitlab" if it exists and is running, or show a menu of instances with "gitlab" in their name:

```
$ essh gitlab
Use the arrow keys to navigate: ↓ ↑ → ←  and / toggles search
Select an instance:
  » fi-gitlab-runner-self-hosted-dev i-05e50f67e9dda4278 (10.100.7.92)
    fi-gitlab-runner-hosted-dev i-0191ea736eca6db2f (10.100.10.29)
```

If you don't know which instance to connect to, run without specifing a tag or instance id:

```
$ essh
Use the arrow keys to navigate: ↓ ↑ → ←  and / toggles search
Select an instance:
    bastion i-06a049e3dbbdc37ae (10.100.12.213)
    eks i-02dbc94c2efe19e68 (10.100.0.67)
    eks i-0907b9bb45af5b43e (10.100.2.252)
    eks i-0d7344c185041ba14 (10.100.4.10)
    fi-gitlab-runner-self-hosted-dev i-05e50f67e9dda4278 (10.100.7.92)
  » eks i-03623ba03fc2dab6f (10.100.9.64)
    eks i-07f1430c8a05d00a7 (10.100.11.246)
    fi-gitlab-runner-hosted-dev i-0191ea736eca6db2f (10.100.10.29)
```

Run with debug logging enabled:

```
$ essh -d -p i-0cc2be02456a7180c
DEBUG Setting region from AWS_DEFAULT_REGION env: eu-west-1
DEBUG All cmd line args passed in
DEBUG flag_pos: 0, flag: i-0cc2be02456a7180c
DEBUG
DEBUG host: 34.245.6.105
DEBUG adding key to agent
DEBUG pushing public key to instance
running command: ssh -l ec2-user 34.245.6.105
Last login: Fri Apr  3 21:40:35 2020 from 90.199.173.2

       __|  __|_  )
       _|  (     /   Amazon Linux 2 AMI
      ___|\___|___|

https://aws.amazon.com/amazon-linux-2/
5 package(s) needed for security, out of 5 available
Run "sudo yum update" to apply all updates.
[ec2-user@ip-172-30-0-254 ~]$
```


## Build

```
$ go build
```

Put the resulting `essh` binary somewhere in your `$PATH`.


## Releasing

To create a new release, just tag the repo and run goreleaser:

```
$ git tag -a [tag] -m "Release message"
$ git push origin [tag]
$ goreleaser --rm-dist
```

## TODO

- Exit with the ssh command exit code
- Add support for setting the default user as an environment variable for shops that use ubuntu etc.
- Add tests


## License

The project is open-source software licensed under the [MIT license](http://opensource.org/licenses/MIT).
