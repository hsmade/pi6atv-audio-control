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

### scraper
This script reads out all the sensors and writes the result into a json file

    ./venv/bin/python scraper.py

### Website
Install nginx and copy (`/etc/nginx/sites-enabled/default`) [default](build/etc/nginx/sites-enabled/default)

#### building the website
When changes have been done, they need to be compiled:

    cd frontend
    yarn build

Copy the build directory to /opt/repeater-audio-control/web
