# milesightct
iot microservices for the milesight ct(x)

## Device submodules

Devices each have their own repo - 
however, they are organized here under the devices directory and their module path is based on this structure

## Process
The process microservice takes our simpleMessage format and interprets its message.

The gives the pipeline microservices access to the MilesightCTReading type.
Process handles a cache/store lookup of a device and appends the information to the devices message.

This allows downstream events and stores to operate without direct access to databases etc...
