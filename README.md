# PI6ATV audio control
This package installs the software to:
 * read the different sensors on the audio control board
 * control the programming of the DSPs
 * expose metrics on the sensors

## Installation
Install the latest debian package from the releases page.

## Development
* React sources [here](frontend/src/App.js)
* backend [here](backend/)

## Testing
### SNMP integration
Configure snmpd with the following:

    pass_persist .1.3.6.1.4.1.8072.2.255 /opt/repeater-audio-control/venv/bin/python -u /opt/repeater-audio-control/snmp-passthrough.py

### backend
Connects to the PCA and servers out the API and metrics

    /opt/repeater-audio-control/backend -config /opt/repeater-audio-control/config.json -verbose

### Website
Install nginx and copy (`/etc/nginx/sites-enabled/default`) [default](build/etc/nginx/sites-enabled/default)

#### building the website
When changes have been done, they need to be compiled:

    cd frontend
    npm start

Copy the build directory to /opt/repeater-audio-control/web

### SNMP
[example snmpd.conf](snmpd.conf)

## Architecture
This software controls the audio for PI6ATV. It runs on a raspberry pi.
It switches the carrier boards on/off by controlling 4066s (a Quad Bilateral Switch) by means of a PCA9671 i2c IO expander.
For programming the DSPs, a USB programmer is used. To control which DSP the programmer can access,
a tca9548a i2c multiplexer is used to select which DSP has its ic2 connected to the programmer.
The pi and the programmer both share the same i2c bus. 
This is because the pi needs to control the i2c multiplexer to be able to connect the right DSP to the bus,
and the programmer then needs to control the bus to program the DSP.
We can't use multi-master on this i2c bus, so we instead disconnect the pi from the bus, when the programmer needs control.
We do this by using two ports on the IO expander.

### pca9671 i2c IO expander
| port | connected to |
| ---- | ------------ |
| 0 | SDA between the pi i2c-2 and the tca9548a |
| 1 | |
| 2 | |
| 3 | |
| 4 | |
| 5 | program enable (DSP) |
| 6 | reset (DSP) |
| 7 | carrier 7.02 |
| 10 | carrier 7.20 |
| 11 | carrier 7.38 |
| 12 | carrier 7.56 |
| 13 | carrier 7.74 |
| 14 | carrier 7.92 |
| 15 | carrier 5.85 |
| 16 | carrier 6.552 |
| 17 | SCL between the pi i2c-2 and the tca9548a |

### tca9548a i2c multiplexer
| port | connected to |
| ---- | ------------ |
| 0 | DSP 1 |
| 1 | DSP 2 |
| 2 | DSP 3 |
| 3 | DSP 4 |
| 4 | DSP 5 |

### logic rules
#### Programming
When a DSP has been selected, the `program` button can be pressed.
This will enable the program pin, and then disconnect the i2c multiplexer from the pi.
The programmer can now be used over USB to program the selected DSP.

If `program` is pressed again, the program pin is disabled, and the i2c bus to the multiplexer is reconnected.
This cancels the program mode.

#### Reset
`Reset` is used to reset the DSP after programming.
When `reset` is pressed:
* the program pin is disabled
* a delay is started
* the reset pin is enabled
* a delay is started
* the reset pin is disabled
* a delay is started
* the i2c bus to the multiplexer is reconnected

## TODO
* tests for the multiplexer