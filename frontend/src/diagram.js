import React from 'react';
import './diagram.css';
import {ArrowRight, ArrowUp, Diamond, Rectangle, RoundedRect, Tube} from './shapes';
import {api, fill, green, orange, red, resetDuration, resetWaitDuration} from './config'

export default class Diagram extends React.Component {
    constructor(props) {
        super(props);

        this.state = {
            expanderPorts: [],
            multiplexerPortSelected: 255,
            program: false,
            reset: false,
            resetButton: false,
            multiplexerError: false,
            multiplexerEnabled: true,
            expanderError: false,
        };
    }

    async sleep(milliseconds) {
        return new Promise(resolve => setTimeout(resolve, milliseconds))
    }

    async fetchState() {
        fetch(`${api}/io/`)
            .then((resp) => {
                // console.log("RESPONSE:", resp)
                if(!resp.ok) throw new Error(`[GET ${api}/io/] returned: ${resp.status}`)
                    resp.json()
                        .then((data) => {
                            this.setState({expanderPorts: data, expanderError: false})
                        })
                        .catch((e) => {
                            console.log(`[GET ${api}/io/] failed to parse json, error:`, e)
                            this.setState({expanderError: true})
                        })
                })
            .catch((e) => {
                console.log(`[GET ${api}/io/] failed to do request:`, e)
                this.setState({expanderError: true})
            })

        fetch(`${api}/mpx/`)
            .then((resp) => {
                // console.log("RESPONSE:", resp)
                if(!resp.ok) throw new Error(`[GET ${api}/mpx/] returned: ${resp.status}`)
                    resp.json()
                        .then((data) => {
                            this.setState({multiplexerPortSelected: data['Port'], multiplexerError: false})
                        })
                        .catch((e) => {
                            console.log(`[GET ${api}/mpx/] failed to parse json, error:`, e)
                            this.setState({multiplexerError: true})
                        })
                })
            .catch((e) => {
                console.log(`[GET ${api}/mpx/] failed to do request:`, e)
                this.setState({multiplexerError: true})
            })
    }

    async setExpanderPort(port, state) {
        const url = `${api}/io/${port}/${state}`
        fetch(url, {method: "POST"})
            .then((resp) => {
                resp.json()
                    .then((data) => {
                        if (data !== null && data['error'] !== undefined) {
                            console.log(`[POST ${url}] set port error:`, data)
                        }
                    })
                    .catch((e) => {
                        console.log(`[POST ${url}] failed to parse json, error:`, e)
                    })
            })
            .catch((e) => console.log(`[POST ${url}] failed to do post request:`, e))
        await this.fetchState()
    }

    async enableCarrier(port) {
        console.log("Enabling port", port)
        await this.setExpanderPort(port, true)
    }

    async disableCarrier(port) {
        console.log("Disabling port", port)
        await this.setExpanderPort(port, false)
    }

    // regularly update the ports
    componentDidMount() {
        try {
            setInterval(this.fetchState.bind(this), 1000);
        } catch(e) {
            console.log(e);
        }
    }

    // The color for the DSP buttons, on the left
    dspButtonColor(port) {
        return this.state.multiplexerPortSelected===port?green:fill
    }

    // The color for the DSPs, in the middle
    dspColor(port) {
        let color = fill
        if (this.state.multiplexerPortSelected===port) {
            color = green
            if (this.state.program) {
                color = red
            }
        }
        if (this.state.reset) {
            color = orange
        }
        return color
    }

    // the color for the carriers, on the right
    carrierColor(port) {
        if (this.state.expanderPorts===undefined) return '#606060'
        return this.state.expanderPorts[port]?green:red
    }

    async toggleProgram() {
        if (this.state.program) {
            console.log("Disable program")
            await this.disableCarrier(5)
        } else {
            if (this.state.multiplexerPortSelected === 255) {
                console.log("Not enabling program as there is no active port")
                return
            }
            console.log("Enable program")
            await this.enableCarrier(5)
        }
        this.setState({program: !this.state.program})
    }

    async startReset() {
        if (this.state.multiplexerPortSelected === 255) {
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

    async disableMultiplexer() {
        console.log("Disabling multiplexer")
        this.setState({multiplexerEnabled: false})
        await this.disableCarrier(0)
        await this.sleep(100)
        await this.disableCarrier(17)
        await this.sleep(100)
    }
    
    async enableMultiplexer() {
        console.log("Enabling multiplexer")
        await this.enableCarrier(0)
        await this.sleep(100)
        await this.enableCarrier(17)
        await this.sleep(300)
        this.setState({multiplexerEnabled: true})
    }
    
    async setDSP(port) {
        if (this.state.program || this.state.resetButton) {
            console.log("Ignoring request to enable DSP", port, "as we're programming/resetting")
            return
        }

        // if port is already enabled, just disable it
        if (this.state.multiplexerPortSelected===port) {
            await this.enableMultiplexer() // first make sure the i2c for the mpx is connected to the port
            await this.multiplexerSetPort(255)
            return
        }

        console.log("Enable DSP",port)
        await this.enableMultiplexer() // make sure it's on, before we try to switch
        await this.multiplexerSetPort(port)
        await this.disableMultiplexer() // disconnect the pi from i2c, so the programmer can take over
    }

    async multiplexerSetPort(port) {
        console.log(`Setting multiplexer to port ${port}`)
        const url = `${api}/mpx/${port}`
        fetch(url, {method: "POST"})
            .then((resp) => {
                resp.json()
                    .then((data) => {
                        if (data !== null && data['error'] !== undefined) {
                            console.log(`[POST ${url}] set DSP error`, data)
                        }
                    })
                    .catch((e) => {
                        console.log(`[POST ${url}] failed to parse json, error:`, e)
                    })
            })
            .catch((e) => console.log(`[POST ${url}] failed to do post request:`, e))
        await this.fetchState()
    }

    async toggleCarrier(port) {
        if (this.state.expanderPorts[port]) {
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
                color={this.dspButtonColor(dsp-1)}
                clickHandler={this.setDSP.bind(this)} clickValue={dsp-1}
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
        let multiplexerError = <div/>
        if (this.state.multiplexerError && this.state.multiplexerEnabled) {
            multiplexerError = <text x={200} y={40} style={{fill: red, fontSize: "1.5em", fontWeight: "bold"}}>
                Error connecting to multiplexer / DSP switch (toggle a DSP)
            </text>
        }
        let expanderError = <div/>
        if (this.state.expanderError) {
            expanderError = <text x={200} y={20} style={{fill: red, fontSize: "1.5em", fontWeight: "bold"}}>
                Error connecting to expander / Carrier switch
            </text>
        }

        const i2cDspDisabled = (this.state.expanderPorts !== undefined && (!this.state.expanderPorts[0] || !this.state.expanderPorts[17]))
        return (
            <svg viewBox={"0 0 1012 763"}>
                {multiplexerError}
                {expanderError}

                {this.drawDspButton(0, 30, 1)}
                <line x1={120} y1={30} x2={160} y2={30}/>

                {this.drawDspButton(0, 130, 2)}
                <line x1={120} y1={130} x2={160} y2={130}/>

                {this.drawDspButton(0, 230, 3)}
                <line x1={120} y1={230} x2={160} y2={230}/>

                {this.drawDspButton(0, 330, 4)}
                <line x1={120} y1={330} x2={160} y2={330}/>

                {this.drawDspButton(0, 430, 5)}
                <line x1={120} y1={430} x2={160} y2={430}/>

                {/*connect all DSPs*/}
                <line x1={160} y1={30} x2={160} y2={430}/>
                <line x1={160} y1={330} x2={220} y2={330}/>
                <ArrowRight x={220} y={330}/>

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
                    text={this.state.expanderPorts[6]?"Reset (on)":"Reset"}
                    color={this.state.resetButton?red:fill}
                    clickHandler={this.startReset.bind(this)}
                />
                <line x1={120} y1={710} x2={298} y2={710}/>

                {/*I2C DSP*/}
                <Diamond
                    x={238} y={215} text={"I2C-DSPs"}
                    text2={i2cDspDisabled?"Inactive":"Active"}
                    color={i2cDspDisabled?red:fill}
                />
                <ArrowUp x={292} y={270}/>
                <line x1={298} y1={330} x2={298} y2={282}/>

                {/*Besturing*/}
                <Diamond
                    x={238} y={320} text={"Besturing"}
                    color={this.state.program?red:fill}
                />
                <ArrowUp x={292} y={376}/>
                <line x1={298} y1={388} x2={298} y2={710}/>
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