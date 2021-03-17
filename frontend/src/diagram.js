import React from 'react';
import './diagram.css';

const stroke = "#41719c"
const fill = "#5b9bd5"

export class RoundedRect extends React.Component {
    // run the click handler, and if we have click values, use them as arguments
    handleClick = () => {
        if (this.props.clickHandler !== undefined) {
            if (this.props.clickValue !== undefined) {
                this.props.clickHandler(this.props.clickValue);
            } else {
                this.props.clickHandler();
            }
        }
    }

    render() {
        const x = parseInt(this.props.x)
        const y = parseInt(this.props.y) - 30
        return (
            <svg width={120} height={60} x={x} y={y}>
                <rect
                    x={1} y={1}
                    width={118} height={58}
                    rx={30}
                    fill={this.props.color}
                    stroke={stroke} strokeWidth={2}
                    onClick={this.handleClick}
                />
                <text
                    x={"50%"} y={"50%"}
                    dominantBaseline={"middle"} textAnchor={"middle"}
                    fill={"white"}
                    onClick={this.handleClick}
                >{this.props.text}</text>
            </svg>
        )
    }
}

export class Rectangle extends React.Component {
    render() {
        const x = parseInt(this.props.x)
        const y = parseInt(this.props.y) - 30
        return (
            <svg x={x} y={y} width={120} height={60}>
                <rect
                    x={1} y={1}
                    width={118} height={58}
                    fill={fill}
                    stroke={stroke} strokeWidth={2}
                />
                <text
                    x={"50%"} y={"50%"}
                    dominantBaseline={"middle"} textAnchor={"middle"}
                    fill={"white"}
                >{this.props.text}</text>
            </svg>
        )
    }
}

export class Tube extends React.Component {
    render() {
        const x = parseInt(this.props.x)
        const y = parseInt(this.props.y) - 30
        return (
            <svg x={x} y={y} width={120} height={60}>
                <ellipse
                    cx={12} cy={30}
                    ry={29} rx={10}
                    fill={this.props.color}
                    stroke={stroke}
                    strokeWidth={2}
                />
                <rect
                    x={10} y={2}
                    width={100} height={56}
                    fill={this.props.color}
                />
                <line x1={11} y1={1}  x2={109} y2={1}  style={{stroke:stroke}}/>
                <line x1={11} y1={59} x2={109} y2={59} style={{stroke:stroke}}/>
                <ellipse
                    cx={108} cy={30}
                    ry={29} rx={10}
                    fill={this.props.color}
                    stroke={stroke} strokeWidth={2}
                />

                <text
                    x={"50%"} y={"50%"}
                    dominantBaseline={"middle"} textAnchor={"middle"}
                    fill={"white"}
                >{this.props.text}</text>
            </svg>
        )
    }
}

export class Diamond extends React.Component {
    render() {
        const x = parseInt(this.props.x)
        const y = parseInt(this.props.y) - 30
        return (
            <svg height={60} width={120} x={x} y={y}>
                <polygon
                    points={"1,30 60,1 119,30 60,59 1,30"}
                    fill={fill}
                    stroke={stroke} strokeWidth={2}
                />
                <text
                    x="50%" y="50%"
                    dominantBaseline="middle" textAnchor="middle"
                    fill="white"
                >{this.props.text}</text>
            </svg>
        )
    }
}

export class ArrowRight extends React.Component {
    render() {
        const x = parseInt(this.props.x)
        const y = parseInt(this.props.y) - 6
        return (
            <svg width={18} height={12} x={x} y={y}>
                <polygon
                    points={"1,1 1,11 17,6 1,1"}
                    fill={fill}
                    stroke={stroke}
                    strokeWidth="2"
                />
            </svg>
        )
    }
}

export class ArrowUp extends React.Component {
    render() {
        const x = parseInt(this.props.x)
        const y = parseInt(this.props.y) - 6
        return (
            <svg width={12} height={18} x={x} y={y}>
                <polygon
                    points={"1,17 11,17 6,1 1,17"}
                    fill={fill}
                    stroke={stroke}
                    strokeWidth="2"
                />
            </svg>
        )
    }
}

function enablePort(port) {
    console.log("Enabling port", port)
}

function disablePort(port) {
    console.log("Disabling port", port)
}

export default class Diagram extends React.Component {
    constructor(props) {
        super(props);

        this.state = {
            ports: [],
            program: false,
            reset: false,
            resetButton: false,
        };
    }

    componentDidMount() {
        try {
            setInterval(async () => {
                const res = await fetch('control-example.json');
                const data = await res.json();
                this.setState({ports: data})
            }, 1000);
        } catch(e) {
            console.log(e);
        }
    }

    dspButtonColor(port) {
        let color = fill
        if (this.state.ports[port] ) {
            color = "#5bb59b"
        }
        return color
    }

    dspColor(port) {
        let color = fill
        if (this.state.ports[port] ) {
            color = "#5bb59b"
            if (this.state.program) {
                color = "#d55b5b"
            }
            if (this.state.reset) {
                color = "#d59b5b"
            }
        }
        return color
    }

    carrierColor(port) {
        let color
        if (this.state.ports[port] ) {
            color = "#5bb59b"
        } else {
            color = "#d55b5b"
        }
        return color
    }

    toggleProgram() {
        if (this.state.program) {
            console.log("Disable program")
            disablePort(5)
        } else {
            console.log("Enable program")
            enablePort(5)
        }
        this.setState({program: !this.state.program})
    }

    startReset() {
        console.log("Enable reset")
        this.setState({resetButton: true})
        disablePort(5) // disable program
        this.setState({program: false})
        setTimeout(function (){
            enablePort(6) // enable reset
            this.setState({reset: true})
            setTimeout(function (){
                disablePort(6) // disable reset
                this.setState({reset: false, resetButton: false})
            }.bind(this), 1000) // FIXME: config
        }.bind(this), 1000) // FIXME: config
    }

    enableDSP(port) {
        console.log("Enable DSP",port)
        for (let i=0; i<=5; i++) {
            disablePort(i)
        }
        enablePort(port)
    }

    toggleCarrier(port) {
        if (this.state.ports[port]) {
            console.log("Disable carrier", port)
            disablePort(port)
        } else {
            console.log("Enable carrier", port)
            enablePort(port)
        }
    }

    drawDspButton(x, y, dsp) {
        return (
            <RoundedRect
                x={x} y={y}
                text={"DSP-" + dsp}
                color={this.dspButtonColor(dsp-1)}
                clickHandler={this.enableDSP} clickValue={dsp-1}
            />
        )
    }

    drawCarrierButton(x, y, port, text) {
        return (
            <RoundedRect
                x={x} y={y}
                text={text}
                color={this.carrierColor(port)}
                clickHandler={this.toggleCarrier.bind(this)} clickValue={port}
            />
        )
    }

    render() {

        return (
            <svg width="1012" height="763">
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

                // connect all DSPs
                <line x1={160} y1={30} x2={160} y2={430}/>

                // Program
                <RoundedRect
                    x={0} y={610}
                    text={"Program"}
                    color={this.state.program?"#d55b5b":fill}
                    clickHandler={this.toggleProgram.bind(this)}
                />
                <line x1={120} y1={610} x2={298} y2={610}/>

                // Reset
                <RoundedRect
                    x={0} y={710}
                    text={"Reset"}
                    color={this.state.resetButton?"#d55b5b":fill}
                    clickHandler={this.startReset.bind(this)}
                />
                <line x1={120} y1={710} x2={298} y2={710}/>

                // Besturing
                <line x1={160} y1={330} x2={220} y2={330}/>
                <ArrowRight x={220} y={330}/>
                <Diamond x={238} y={330} text={"Besturing"}/>
                <ArrowUp x={292} y={366}/>
                <line x1={298} y1={366} x2={298} y2={710}/>
                <line x1={358} y1={330} x2={398} y2={330}/>
                <line x1={398} y1={78} x2={398} y2={710}/>

                // DSP 1
                <line x1={398} y1={78} x2={478} y2={78}/>
                <ArrowRight x={478} y={78}/>
                <Tube x={496} y={78} text={"DSP-1"} color={this.dspColor(0)}/>
                // Carrier 1
                <line x1={616} y1={78} x2={656} y2={78}/>
                <ArrowRight x={656} y={78}/>
                <Rectangle x={674} y={78} text={"Analog carrier"}/>
                <line x1={794} y1={78} x2={834} y2={78}/>
                <line x1={834} y1={33} x2={834} y2={123}/>
                // 7.02
                <line x1={834} y1={33} x2={874} y2={33}/>
                <ArrowRight x={874} y={33}/>
                {this.drawCarrierButton(892, 33, 7, "7.02")}
                // 7.20
                <line x1={834} y1={123} x2={874} y2={123}/>
                <ArrowRight x={874} y={123}/>
                {this.drawCarrierButton(892, 123, 10, "7.20")}

                // DSP 2
                <line x1={398} y1={258} x2={478} y2={258}/>
                <ArrowRight x={478} y={258}/>
                <Tube x={496} y={258} text={"DSP-2"} color={this.dspColor(1)}/>
                // Carrier 2
                <line x1={616} y1={258} x2={656} y2={258}/>
                <ArrowRight x={656} y={258}/>
                <Rectangle x={674} y={258} text={"Analog carrier"}/>
                <line x1={794} y1={258} x2={834} y2={258}/>
                <line x1={834} y1={213} x2={834} y2={303}/>
                // 7.38
                <line x1={834} y1={213} x2={874} y2={213}/>
                <ArrowRight x={874} y={213}/>
                {this.drawCarrierButton(892, 213, 11, "7.38")}
                // 7.56
                <line x1={834} y1={303} x2={874} y2={303}/>
                <ArrowRight x={874} y={303}/>
                {this.drawCarrierButton(892, 303, 12, "7.56")}

                // DSP 3
                <line x1={398} y1={438} x2={478} y2={438}/>
                <ArrowRight x={478} y={438}/>
                <Tube x={496} y={438} text={"DSP-3"} color={this.dspColor(2)}/>
                // Carrier 3
                <line x1={616} y1={438} x2={656} y2={438}/>
                <ArrowRight x={656} y={438}/>
                <Rectangle x={674} y={438} text={"Analog carrier"}/>
                <line x1={794} y1={438} x2={834} y2={438}/>
                <line x1={834} y1={393} x2={834} y2={483}/>
                // 7.74
                <line x1={834} y1={393} x2={874} y2={393}/>
                <ArrowRight x={874} y={393}/>
                {this.drawCarrierButton(892, 393, 13, "7.74")}
                // 7.92
                <line x1={834} y1={483} x2={874} y2={483}/>
                <ArrowRight x={874} y={483}/>
                {this.drawCarrierButton(892, 483, 14, "7.92")}

                // DSP 4
                <line x1={398} y1={610} x2={478} y2={610}/>
                <ArrowRight x={478} y={610}/>
                <Tube x={496} y={610} text={"DSP-4"} color={this.dspColor(3)}/>
                // Carrier 4
                <line x1={616} y1={610} x2={656} y2={610}/>
                <ArrowRight x={656} y={610}/>
                <Rectangle x={674} y={610} text={"Nicam"}/>
                <line x1={794} y1={610} x2={834} y2={610}/>
                // 5.85
                <line x1={834} y1={610} x2={874} y2={610}/>
                <ArrowRight x={874} y={610}/>
                {this.drawCarrierButton(892, 610, 15, "5.85")}

                // DSP 5
                <line x1={398} y1={710} x2={478} y2={710}/>
                <ArrowRight x={478} y={710}/>
                <Tube x={496} y={710} text={"DSP-5"} color={this.dspColor(4)}/>
                // Carrier 5
                <line x1={616} y1={710} x2={656} y2={710}/>
                <ArrowRight x={656} y={710}/>
                <Rectangle x={674} y={710} text={"Nicam"}/>
                <line x1={794} y1={710} x2={834} y2={710}/>
                // 6.552
                <line x1={834} y1={710} x2={874} y2={710}/>
                <ArrowRight x={874} y={710}/>
                {this.drawCarrierButton(892, 710, 16, "6.552")}
            </svg>
        )
    }
}