# Running Benchmarks
Forked from [hashrand-rs](https://github.com/akhilsb/hashrand-rs).

This document explains how to benchmark the codebase and read benchmarks' results. It also provides a step-by-step tutorial to run benchmarks on [Amazon Web Services (AWS)](https://aws.amazon.com) accross multiple data centers (WAN).

## Setup
We will use [Fabric](http://www.fabfile.org/) purely to start up the testbed. First, install the python dependencies:

```
$ pip install -r requirements.txt
```

You also need to install [tmux](https://linuxize.com/post/getting-started-with-tmux/#installing-tmux) (which runs all nodes and clients in the background). 

## AWS Benchmarks

### Step 1. Set up your AWS credentials
Set up your AWS credentials to enable programmatic access to your account from your local machine. These credentials will authorize your machine to create, delete, and edit instances on your AWS account programmatically. First of all, [find your 'access key id' and 'secret access key'](https://docs.aws.amazon.com/cli/latest/userguide/cli-configure-quickstart.html#cli-configure-quickstart-creds). Then, create a file `~/.aws/credentials` with the following content:
```
[default]
aws_access_key_id = YOUR_ACCESS_KEY_ID
aws_secret_access_key = YOUR_SECRET_ACCESS_KEY
```
Do not specify any AWS region in that file as the python scripts will allow you to handle multiple regions programmatically.

### Step 2. Add your SSH public key to your AWS account
You must now [add your SSH public key to your AWS account](https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/ec2-key-pairs.html). This operation is manual (AWS exposes little APIs to manipulate keys) and needs to be repeated for each AWS region that you plan to use. Upon importing your key, AWS requires you to choose a 'name' for your key; ensure you set the same name on all AWS regions. This SSH key will be used by the python scripts to execute commands and upload/download files to your AWS instances.
If you don't have an SSH key, you can create one using [ssh-keygen](https://www.ssh.com/ssh/keygen/):
```
$ ssh-keygen -f ~/.ssh/aws
```

### Step 3. Configure the testbed
The file settings.json contains all the configuration parameters of the testbed to deploy. Its content looks as follows:
```json
{
    "key": {
        "name": "aws",
        "path": "/absolute/key/path"
    },
    "port": 8500,
    "client_base_port": 9000,
    "client_run_port": 9500,
    "repo": {
        "name": "hashrand-rs",
        "url": "https://github.com/akhilsb/hashrand-rs.git",
        "branch": "master"
    },
    "instances": {
        "type": "t3a.medium",
        "regions": ["us-east-1","us-east-2","us-west-1","us-west-2","ca-central-1", "eu-west-1", "ap-southeast-1", "ap-northeast-1"]
    }
}
```
The first block (`key`) contains information regarding your SSH key:
```json
"key": {
    "name": "aws",
    "path": "/absolute/key/path"
},
```
Enter the name of your SSH key; this is the name you specified in the AWS web console in step 2. Also, enter the absolute path of your SSH private key (using a relative path won't work). 


The second block (`ports`) specifies the TCP ports to use:
```json
"port": 8500,
"client_base_port": 9000,
"client_run_port": 9500,
```
The artifact requires a number of TCP ports for communication between the processes. Note that the script will open a large port range (5000-10000) to the WAN on all your AWS instances. 

The the last block (`instances`) specifies the [AWS instance type](https://aws.amazon.com/ec2/instance-types) and the [AWS regions](https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/using-regions-availability-zones.html#concepts-available-regions) to use:
```json
"instances": {
    "type": "t3a.medium",
    "regions": ["us-east-1","us-east-2","us-west-1","us-west-2","ca-central-1", "eu-west-1", "ap-southeast-1", "ap-northeast-1"]
}
```
The instance type selects the hardware on which to deploy the testbed. For example, `t3a.medium` instances come with 2 vCPU (2 physical cores), and 4 GB of RAM. The python scripts will configure each instance with 300 GB of SSD hard drive. The `regions` field specifies the data centers to use. If you require more nodes than data centers, the python scripts will distribute the nodes as equally as possible amongst the data centers. All machines run a fresh install of Ubuntu Server 20.04.

### Step 4. Create a testbed
The AWS instances are orchestrated with [Fabric](http://www.fabfile.org) from the file [fabfile.py](https://github.com/akhil-sb/hashrand-rs/blob/master/benchmark/fabfile.py) (located in [hashrand-rs/benchmarks](https://github.com/akhilsb/hashrand-rs/blob/master/benchmark)); you can list all possible commands as follows:
```
$ cd hashrand-rs/benchmark
$ fab --list
```
The command `fab create` creates new AWS instances; open [fabfile.py](https://github.com/akhilsb/hashrand-rs/blob/master/benchmark/fabfile.py) and locate the `create` task:
```python
@task
def create(ctx, nodes=2):
    ...
```
The parameter `nodes` determines how many instances to create in *each* AWS region. That is, if you specified 5 AWS regions as in the example of step 3, setting `nodes=2` will creates a total of 16 machines:
```
$ fab create

Creating 16 instances |██████████████████████████████| 100.0% 
Waiting for all instances to boot...
Successfully created 16 new instances
```

The commands `fab stop` and `fab start` respectively stop and start the testbed without destroying it (it is good practice to stop the testbed when not in use as AWS can be quite expensive); and `fab destroy` terminates all instances and destroys the testbed. Note that, depending on the instance types, AWS instances may take up to several minutes to fully start or stop. The command `fab info` displays a nice summary of all available machines and information to manually connect to them (for debug).

### Step 5. Verify Testbed

After creating the testbed, you can check the status with this command:
```
$ fab info
```

### Step 6. Run Scripts

Run the scripts as described in the main README.

### Step 7. Destroy
After running the benchmarks, destroy the testbed with the following command. 
```
$ fab destroy
```
This command destroys the testbed and terminates all created AWS instances.
