# essh

SSH to an EC2 instance using an in memory, ephemeral ssh key and EC2 instance connect to push the new public key to the instance


## Description

`essh` does the following:

- Generates a one time RSA ssh keypair in memory
- Adds the private key to your ssh agent define by `SSH_AUTH_SOCK` (for a configurable number of seconds)
- Pushes the public key to the instance using [ec2-instance connect](https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/Connect-using-EC2-Instance-Connect.html)
- `ssh` to the instance using the private IP address (public IP can be used with `-p`), using user `ec2-user` by default


## Requirements

As `essh` uses AWS APIs, you will need you have valid credentials configured. If you're using this tool, then I'm presuming that you know how to do this, if not [see here](https://docs.aws.amazon.com/cli/latest/userguide/cli-chap-configure.html).

You should set the region with the `-r`/`--region` flag, or by setting the environment variable `AWS_DEFAULT_REGION` or `AWS_REGION`.


## Demo

[![asciicast](https://asciinema.org/a/318394.svg)](https://asciinema.org/a/318394?autoplay=1)

## Usage

```shell
$ essh --help
Usage of essh:
  -d, --debug             Enable debug logging
  -t, --key-ttl uint32    How long the private key will live in the ssh-agent in seconds (default 10)
  -r, --region string     AWS Region
  -p, --use-public-ip     Use the public ip instead of the private ip address
  -u, --username string   UNIX user name (default "ec2-user")
  -v, --version           Show version
```

Connect to an instance's private IP with ssh as user `ec2-user` (the default):

```shell
$ essh i-02fab0d7dd3ab737b
```

Connect to an instance's public IP with ssh as user `ec2-user`:

```shell
$ essh -p i-02fab0d7dd3ab737b
```

Connect as user `ubuntu` passing the flags `-A`, `-4` and the command `uptime` to the ssh command:

```shell
$ essh -u ubuntu i-02fab0d7dd3ab737b -- -A -4 uptime
```

Connect to an instance by its full name tag:

```shell
$ essh prod-web1
```

Display a menu of instances that match a partial tag:

```shell
$ essh gitlab
```

Display all running instances in a region:

```shell
$ essh
```

You can use `/` to search the list of instances.


## Changing the default UNIX user

If you use a different operating system that does not use the username `ec2-user`, you can set a different default username.

For example, if you use Ubuntu, you can set the environment variable:

```shell
$ export ESSH_DEFAULT_USER=ubuntu
```

From then on, you can just omit the `-u ubuntu` flag to log in as the `ubuntu` user:


## Usage Examples

Connect to an instance on its private IP:

```shell
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

```shell
$ essh -p prod-web1 -- uptime
running command: ssh -l ec2-user 52.51.41.123 uptime
 16:42:42 up 16 min,  0 users,  load average: 0.13, 0.04, 0.01
```

Connect to a host named "gitlab" if it exists and is running, or show a menu of instances with "gitlab" in their name:

```shell
$ essh gitlab
Use the arrow keys to navigate: ↓ ↑ → ←  and / toggles search
Select an instance:
  » fi-gitlab-runner-self-hosted-dev i-05e50f67e9dda4278 (10.100.7.92)
    fi-gitlab-runner-hosted-dev i-0191ea736eca6db2f (10.100.10.29)
```

If you don't know which instance to connect to, run without specifying a tag or instance id:

```shell
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

```shell
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

```shell
$ go build
```

Put the resulting `essh` binary somewhere in your `$PATH`.


## Releasing

To create a new release, just tag the repo and run goreleaser:

```shell
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
