# gghelper - AWS Greengrass helper

This is a commandline tool to assist in setting up [AWS Greengrass](https://aws.amazon.com/greengrass)
installations. 

## Workflow

A typical workflow is:

1. Create a Greengrass group

	```
	gghelper creategroup -name test
	```

1. The output of the creategroup will be a tar file containing the certificates and configuration for the Greengrass Core. This should be transferred onto the GGC system and unpacked:
	```
	tar xzf 5d7b82589d-setup.tar.gz -C /greengrass
	```

1. The GGC can be started with this configuration.
	```
	(cd /greengrass/ggc/core; sudo ./greengrassd start)
	```

1. Add in a lambda function
	```
	gghelper lambda -pinned -d dist -handler hello_vsphere.hello_vsphere_handler -name HellovSphere -role lambda-test-get -runtime python2.7
	```

1. Create a subscription between the function and cloud
	```
	gghelper createsub -source HellovSphere -target cloud -subject "hello/vsphere"
	```

1. Make a deployment to download config and code to the Greengrass core
   ```
   gghelper createdeployment
   ```

1. Going to the Greengrass Test page, create a subscription (using # will get all the messages) to see the function run periodically every 5 seconds.


## Credits
This project was sponsored by [VMware](http://www.vmware.com). And inspiration and some compatibility by the AWS Labs [aws-greengrass-group-setup](https://github.com/awslabs/aws-greengrass-group-setup).