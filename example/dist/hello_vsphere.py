import greengrasssdk
from threading import Timer


client = greengrasssdk.client('iot-data')


def hello_vsphere():
    client.publish(topic='hello/vsphere', payload='Hello vSphere!')
    Timer(5, hello_vsphere).start()


hello_vsphere()


def hello_vsphere_handler(event, context):
    return
