import React from 'react';
import './diagram.css';
import {ArrowRight, ArrowUp, Diamond, Rectangle, RoundedRect, Tube} from './shapes';
import {api, fill, green, orange, red, resetDuration, resetWaitDuration} from './config'

export default class Diagram extends React.Component {
    constructor(props) {
        super(props);

        this.state = {
            carrierPorts: [],
            programmerSelected: 255,
            program: false,
            reset: false,
            resetButton: false,
            error: false,
        };
    }

    async fetchState() {
        fetch(api)
            .then((resp) => {
                // console.log("RESPONSE:", resp)
                if(!resp.ok) throw new Error(`backend returned: ${resp.status}`)
                    resp.json()
                        .then((data) => {
                            this.setState({ports: data['pca'], programmerSelected: data['tca'], error: false})
                        })
                        .catch((e) => {
                            console.log("failed to parse json, error:",e)
                            this.setState({error: true})
                        })
                })
            .catch((e) => {
                console.log("failed to do request:", api, e)
                this.setState({error: true})
            })
    }

    async setCarrierPort(port, state) {
        const url = api + port + "/" + state
        fetch(url, {method: "POST"})
            .then((resp) => {
                resp.json()
                    .then((data) => {
                        if (data === null || data['error'] !== "") {
                            console.log("set port error:", data)
                        }
                    })
            })
            .catch((e) => console.log("failed to do post request:", url, e))
        await this.fetchState()
    }

    async enableCarrier(port) {
        console.log("Enabling port", port)
        await this.setCarrierPort(port, true)
    }

    async disableCarrier(port) {
        console.log("Disabling port", port)
        await this.setCarrierPort(port, false)
    }

    // regularly update the ports
    componentDidMount() {
        try {
            setInterval(this.fetchState.bind(this), 1000);
        } catch(e) {
            console.log(e);
        }
    }

    dspButtonColor(port) {
        return this.state.programmerSelected===port?green:fill
    }

    dspColor(port) {
        let color = fill
        if (this.state.ports[port] ) {
            color = green
            if (this.state.program) {
                color = red
            }
            if (this.state.reset) {
                color = orange
            }
        }
        return color
    }

    carrierColor(port) {
        return this.state.ports[port]?green:red
    }

    hasPortSelected() {
        for (let i=0; i<=5; i++) {
            if (this.state.ports[i]) return true
        }
        return false
    }

    async toggleProgram() {
        if (this.state.program) {
            console.log("Disable program")
            await this.disableCarrier(5)
        } else {
            if (!this.hasPortSelected()) {
                console.log("Not enabling program as there is no active port")
                return
            }
            console.log("Enable program")
            await this.enableCarrier(5)
        }
        this.setState({program: !this.state.program})
    }

    async startReset() {
        if (!this.hasPortSelected()) {
            console.log("Not enabling reset as there is no active port")
            return
        }

        console.log("Enable reset")
        this.setState({resetButton: true})

        await this.disableCarrier(5) // disable program
        this.setState({program: false})

        setTimeout(async function (){
            await this.enableCarrier(6) // enable reset
            this.setState({reset: true})

            setTimeout(async function (){
                await this.disableCarrier(6) // disable reset
                this.setState({reset: false, resetButton: false})
            }.bind(this), resetDuration)

        }.bind(this), resetWaitDuration)
    }

    async setDSP(port) {
        if (this.state.program || this.state.resetButton) {
            console.log("Ignoring request to enable DSP", port, "as we're programming/resetting")
            return
        }

        // if port is already enabled, just disable it
        if (this.state.ports[port]) {
            await this.ProgrammerSetPort(255)
            return
        }

        console.log("Enable DSP",port)
        await this.ProgrammerSetPort(port)
    }

    async ProgrammerSetPort(port) {
        const url = api + "/programmer/" + port
        fetch(url, {method: "POST"})
            .then((resp) => {
                resp.json()
                    .then((data) => {
                        if (data === null || data['error'] !== "") {
                            console.log("set DSP error:", data)
                        }
                    })
            })
            .catch((e) => console.log("failed to do post request:", url, e))
        await this.fetchState()
    }

    async toggleCarrier(port) {
        if (this.state.ports[port]) {
            console.log("Disable carrier", port)
            await this.disableCarrier(port)
        } else {
            console.log("Enable carrier", port)
            await this.enableCarrier(port)
        }
    }

    drawDspButton(x, y, dsp) {
        return (
            <RoundedRect
                x={x} y={y}
                text={"DSP-" + dsp}
                color={this.dspButtonColor(dsp)}
                clickHandler={this.setDSP.bind(this)} clickValue={dsp}
            />
        )
    }

    drawCarrierButton(x, y, port, text) {
        const textContent = [
                            <tspan key={`${port}_name`} x={60} y={20}>{text}</tspan>,
                            <tspan key={`${port}_state`} x={60} y={40}>{this.carrierColor(port)===green?"Active":"Inactive"}</tspan>
        ]
        return (
            <RoundedRect
                x={x} y={y}
                text={textContent}
                color={this.carrierColor(port)}
                clickHandler={this.toggleCarrier.bind(this)} clickValue={port}
            />
        )
    }

    render() {
        let error = <div/>
        if (this.state.error) {
            error = <text x={200} y={30} style={{fill: red, fontSize: "2em", fontWeight: "bold"}}>
                Error connecting to controller or device
            </text>
        }
        return (
            <svg viewBox={"0 0 1012 763"}>
                {error}

                {this.drawDspButton(0, 30, 0)}
                <line x1={120} y1={30} x2={160} y2={30}/>

                {this.drawDspButton(0, 130, 1)}
                <line x1={120} y1={130} x2={160} y2={130}/>

                {this.drawDspButton(0, 230, 2)}
                <line x1={120} y1={230} x2={160} y2={230}/>

                {this.drawDspButton(0, 330, 3)}
                <line x1={120} y1={330} x2={160} y2={330}/>

                {this.drawDspButton(0, 430, 4)}
                <line x1={120} y1={430} x2={160} y2={430}/>

                {/*connect all DSPs*/}
                <line x1={160} y1={30} x2={160} y2={430}/>

                {/*Program*/}
                <RoundedRect
                    x={0} y={610}
                    text={"Program"}
                    color={this.state.program?red:fill}
                    clickHandler={this.toggleProgram.bind(this)}
                />
                <line x1={120} y1={610} x2={298} y2={610}/>

                {/*Reset*/}
                <RoundedRect
                    x={0} y={710}
                    text={"Reset"}
                    color={this.state.resetButton?red:fill}
                    clickHandler={this.startReset.bind(this)}
                />
                <line x1={120} y1={710} x2={298} y2={710}/>

                {/*Besturing*/}
                <line x1={160} y1={330} x2={220} y2={330}/>
                <ArrowRight x={220} y={330}/>
                <Diamond
                    x={238} y={330} text={"Besturing"}
                    color={this.state.program?red:fill}
                />
                <ArrowUp x={292} y={366}/>
                <line x1={298} y1={366} x2={298} y2={710}/>
                <line x1={358} y1={330} x2={398} y2={330}/>
                <line x1={398} y1={78} x2={398} y2={710}/>

                {/*DSP 1*/}
                <line x1={398} y1={78} x2={478} y2={78}/>
                <ArrowRight x={478} y={78}/>
                <Tube x={496} y={78} text={"DSP-1"} color={this.dspColor(0)}/>
                {/*Carrier 1*/}
                <line x1={616} y1={78} x2={656} y2={78}/>
                <ArrowRight x={656} y={78}/>
                <Rectangle x={674} y={78} text={"Analog carrier"}/>
                <line x1={794} y1={78} x2={834} y2={78}/>
                <line x1={834} y1={33} x2={834} y2={123}/>
                {/*7.02*/}
                <line x1={834} y1={33} x2={874} y2={33}/>
                <ArrowRight x={874} y={33}/>
                {this.drawCarrierButton(892, 33, 7, "7.02")}
                {/*7.20*/}
                <line x1={834} y1={123} x2={874} y2={123}/>
                <ArrowRight x={874} y={123}/>
                {this.drawCarrierButton(892, 123, 10, "7.20")}

                {/*DSP 2*/}
                <line x1={398} y1={258} x2={478} y2={258}/>
                <ArrowRight x={478} y={258}/>
                <Tube x={496} y={258} text={"DSP-2"} color={this.dspColor(1)}/>
                {/*Carrier 2*/}
                <line x1={616} y1={258} x2={656} y2={258}/>
                <ArrowRight x={656} y={258}/>
                <Rectangle x={674} y={258} text={"Analog carrier"}/>
                <line x1={794} y1={258} x2={834} y2={258}/>
                <line x1={834} y1={213} x2={834} y2={303}/>
                {/*7.38*/}
                <line x1={834} y1={213} x2={874} y2={213}/>
                <ArrowRight x={874} y={213}/>
                {this.drawCarrierButton(892, 213, 11, "7.38")}
                {/*7.56*/}
                <line x1={834} y1={303} x2={874} y2={303}/>
                <ArrowRight x={874} y={303}/>
                {this.drawCarrierButton(892, 303, 12, "7.56")}

                {/*DSP 3*/}
                <line x1={398} y1={438} x2={478} y2={438}/>
                <ArrowRight x={478} y={438}/>
                <Tube x={496} y={438} text={"DSP-3"} color={this.dspColor(2)}/>
                {/*Carrier 3*/}
                <line x1={616} y1={438} x2={656} y2={438}/>
                <ArrowRight x={656} y={438}/>
                <Rectangle x={674} y={438} text={"Analog carrier"}/>
                <line x1={794} y1={438} x2={834} y2={438}/>
                <line x1={834} y1={393} x2={834} y2={483}/>
                {/*7.74*/}
                <line x1={834} y1={393} x2={874} y2={393}/>
                <ArrowRight x={874} y={393}/>
                {this.drawCarrierButton(892, 393, 13, "7.74")}
                {/*7.92*/}
                <line x1={834} y1={483} x2={874} y2={483}/>
                <ArrowRight x={874} y={483}/>
                {this.drawCarrierButton(892, 483, 14, "7.92")}

                {/*DSP 4*/}
                <line x1={398} y1={610} x2={478} y2={610}/>
                <ArrowRight x={478} y={610}/>
                <Tube x={496} y={610} text={"DSP-4"} color={this.dspColor(3)}/>
                {/*Carrier 4*/}
                <line x1={616} y1={610} x2={656} y2={610}/>
                <ArrowRight x={656} y={610}/>
                <Rectangle x={674} y={610} text={"Nicam"}/>
                <line x1={794} y1={610} x2={834} y2={610}/>
                {/*5.85*/}
                <line x1={834} y1={610} x2={874} y2={610}/>
                <ArrowRight x={874} y={610}/>
                {this.drawCarrierButton(892, 610, 15, "5.85")}

                {/*DSP 5*/}
                <line x1={398} y1={710} x2={478} y2={710}/>
                <ArrowRight x={478} y={710}/>
                <Tube x={496} y={710} text={"DSP-5"} color={this.dspColor(4)}/>
                {/*Carrier 5*/}
                <line x1={616} y1={710} x2={656} y2={710}/>
                <ArrowRight x={656} y={710}/>
                <Rectangle x={674} y={710} text={"Nicam"}/>
                <line x1={794} y1={710} x2={834} y2={710}/>
                {/* 6.552*/}
                <line x1={834} y1={710} x2={874} y2={710}/>
                <ArrowRight x={874} y={710}/>
                {this.drawCarrierButton(892, 710, 16, "6.552")}
            </svg>
        )
    }
}