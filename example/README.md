# Example Greengrass lambda function

The gghelper utility will package and upload a zip file containing the Greengrass lambda function along with its dependencies.

This directory contains a "dist" directory containing the hello_vsphere.py function and the code to be uploaded. The Greengrass Python SDK should be unpacked into this directory to look like:

```
$ ls -1F dist
greengrass_common/
greengrass_ipc_python_sdk/
greengrasssdk/
hello_vsphere.py
```

This command will then upload the code as a zip file:

```
gghelper lambda -pinned -d dist -handler hello_vsphere.hello_vsphere_handler -name HellovSphere -role lambda-test-get -runtime python2.7
```
