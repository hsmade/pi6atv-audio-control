import React from 'react';

export class Ellipse extends React.Component {
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
        const y = parseInt(this.props.y) - 33
        return (
            <svg height={66} width={159} x={x} y={y}>
            <path
                fill={this.props.color}
                stroke="#41719c"
                strokeWidth="1.67"
                d="M32 64h94.5c17.32 0 31.5-14.1 31.5-31.5s-14.18-31.5-31.5-31.5h-94.5c-17.3 0-31.48 14.1-31.48 31.5s14.17 31.5 31.5 31.5z"
                onClick={this.handleClick}
            />
                <text x="50%" y="50%" dominantBaseline="middle" textAnchor="middle" fill="white" onClick={this.handleClick}>{this.props.text}</text>
            </svg>
        )
    }
}

export class Rectangle extends React.Component {
    render() {
        const x = parseInt(this.props.x)
        const y = parseInt(this.props.y) - 47
        return (
            <svg height={94} width={158} x={x} y={y}>
                <path
                    fill="#5b9bd5"
                    stroke="#c8c8c8"
                    strokeWidth="0.56"
                    d="M0 94c0 .17.12.3.28.3h156.92c.15 0 .28-.13.28-.3v-93.9c0-.17-.13-.3-.28-.3h-156.92c-.16 0-.28.13-.28.3z"
                />
                <text x="50%" y="50%" dominantBaseline="middle" textAnchor="middle" fill="white">{this.props.text}</text>
            </svg>
        )
    }
}

export class Disk extends React.Component {
    render() {
        const x = parseInt(this.props.x)
        const y = parseInt(this.props.y) - 48
        return (
            <svg height={96} width={173} x={x} y={y}>
                <path
                    fill="#5b9bd5"
                    stroke="#c8c8c8"
                    strokeWidth="0.56"
                    d="M154 0c9.4 0 17.02 21.15 17.02 47.24 0 26.1-7.6 47.25-17 47.25H18c-9.4 0-17-21.15-17-47.25s7.6-47.24 17-47.24z"
                />
                <path
                    fill="none"
                    stroke="#c8c8c8"
                    strokeWidth="0.56"
                    d="M154 94.49c-9.38 0-17-21.15-17-47.25s7.62-47.24 17-47.24"
                />
                <text x="50%" y="50%" dominantBaseline="middle" textAnchor="middle" fill="white">{this.props.text}</text>
            </svg>
        )
    }
}

export class Ruit extends React.Component {
    render() {
        const x = parseInt(this.props.x)
        const y = parseInt(this.props.y) - 48
        return (
            <svg height={96} width={158} x={x} y={y}>
                <path
                    fill="#4f88bb"
                    stroke="#41719c"
                    strokeWidth="1.67"
                    d="M1 48c-.4-.24-.4-.62 0-.86l77.3-46.38c.4-.23 1.04-.23 1.43 0l77.32 46.4c.4.23.4.6 0 .85l-77.32 46.38c-.4.24-1.03.24-1.42 0z"
                />
                <text x="50%" y="50%" dominantBaseline="middle" textAnchor="middle" fill="white">{this.props.text}</text>
            </svg>
        )
    }
}

export class ArrowRight extends React.Component {
    render() {
        const x = parseInt(this.props.x)
        const y = parseInt(this.props.y) - 6
        return (
            <svg height={12} width={18} x={x} y={y}>
                <path
                    fill="#5b9bd5"
                    stroke="#5b9bd5"
                    strokeWidth="2.22"
                    d="M15 6l-15 4.22V2z"
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
            <svg height={18} width={14} x={x} y={y}>
                <path
                    fill="#5b9bd5"
                    stroke="#5b9bd5"
                    strokeWidth="2.22"
                    d="M7.11 2.22 L2.22 16.78 L12 16.78 Z"
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
            reset: false
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

    dspColor(port) {
        let color = "#5b9bd5"
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
        let color = "#5b9bd5"
        if (this.state.ports[port] ) {
            color = "#5bb59b"
        } else {
            color = "#d55b5b"
        }
        return color
    }

    startProgram() {
        console.log("Enable program")
        enablePort(5)
        this.setState({program: true})
    }

    startReset() {
        console.log("Enable reset")
        disablePort(5)
        this.setState({program: false})
        setTimeout(function (){
            enablePort(6)
            this.setState({reset: true})
            setTimeout(function (){
                disablePort(6)
                this.setState({reset: false})
            }.bind(this), 1000)
        }.bind(this), 1000)
    }

    render() {
        function enableDSP(port) {
            console.log("Enable DSP",port)
            for (let i=0; i<=5; i++) {
                disablePort(i)
            }
            enablePort(port)
        }

        function toggleCarrier(port) {
            if (this.state.ports[port]) {
                console.log("Disable carrier", port)
                disablePort(port)
            } else {
                console.log("Enable carrier", port)
                enablePort(port)
            }
        }

        return (
            <svg
                width="1315"
                height="800"
            >
                <g>
                    <Ellipse x={0} y={33} text={"DSP-1"} color={this.dspColor(0)} clickHandler={enableDSP} clickValue={0}/>
                    <line x1="159" y1="33" x2="200" y2="33" style={{stroke:"#5b9bd5", strokeWidth:2.22}}/>
                    <Ellipse x={0} y={133} text={"DSP-2"} color={this.dspColor(1)} clickHandler={enableDSP} clickValue={1}/>
                    <line x1="159" y1="133" x2="200" y2="133" style={{stroke:"#5b9bd5", strokeWidth:2.22}}/>
                    <Ellipse x={0} y={233} text={"DSP-3"} color={this.dspColor(2)} clickHandler={enableDSP} clickValue={2}/>
                    <line x1="159" y1="233" x2="200" y2="233" style={{stroke:"#5b9bd5", strokeWidth:2.22}}/>
                    <Ellipse x={0} y={333} text={"DSP-4"} color={this.dspColor(3)} clickHandler={enableDSP} clickValue={3}/>
                    <line x1="159" y1="333" x2="400" y2="333" style={{stroke:"#5b9bd5", strokeWidth:2.22}}/>
                    <Ellipse x={0} y={433} text={"DSP-5"} color={this.dspColor(4)} clickHandler={enableDSP} clickValue={4}/>
                    <line x1="159" y1="433" x2="200" y2="433" style={{stroke:"#5b9bd5", strokeWidth:2.22}}/>
                    <line x1="200" y1="33" x2="200" y2="433" style={{stroke:"#5b9bd5", strokeWidth:2.22}}/>

                    <Ellipse x={0} y={633} text={"Program"} color={this.state.program?"red":"blue"} clickHandler={this.startProgram.bind(this)}/>
                    <line x1="159" y1="633" x2="379" y2="633" style={{stroke:"#5b9bd5", strokeWidth:2.22}}/>
                    <Ellipse x={0} y={733} text={"Reset"} color={this.state.reset?"red":"blue"} clickHandler={this.startReset.bind(this)}/>
                    <line x1="159" y1="733" x2="379" y2="733" style={{stroke:"#5b9bd5", strokeWidth:2.22}}/>

                    <ArrowRight x={282} y={333}/>
                    <Ruit x={300} y={333} text={"Besturing"}/>
                    <line x1="458" y1="333" x2="498" y2="333" style={{stroke:"#5b9bd5", strokeWidth:2.22}}/>
                    <ArrowUp x={372} y={386}/>
                    <line x1="379" y1="386" x2="379" y2="733" style={{stroke:"#5b9bd5", strokeWidth:2.22}}/>

                    <line x1={498} y1={78} x2={582} y2={78} style={{stroke:"#5b9bd5", strokeWidth:2.22}}/>
                    <ArrowRight x={582} y={78}/>
                    <Disk x={600} y={78} text={"DSP-1"}/>
                    <line x1={771} y1={78} x2={882} y2={78} style={{stroke:"#5b9bd5", strokeWidth:2.22}}/>
                    <ArrowRight x={882} y={78}/>
                    <Rectangle x={900} y={78} text={"Analog carrier"}/>
                    <line x1={1058} y1={78} x2={1098} y2={78} style={{stroke:"#5b9bd5", strokeWidth:2.22}}/>
                    <line x1={1098} y1={33} x2={1098} y2={123} style={{stroke:"#5b9bd5", strokeWidth:2.22}}/>

                    <line x1={1098} y1={33} x2={1138} y2={33} style={{stroke:"#5b9bd5", strokeWidth:2.22}}/>
                    <ArrowRight x={1138} y={33}/>
                    <Ellipse x={1156} y={33} text={"7.02"} color={this.carrierColor(7)} clickHandler={toggleCarrier.bind(this)} clickValue={7}/>

                    <line x1={1098} y1={123} x2={1138} y2={123} style={{stroke:"#5b9bd5", strokeWidth:2.22}}/>
                    <ArrowRight x={1138} y={123}/>
                    <Ellipse x={1156} y={123} text={"7.20"} color={this.carrierColor(10)} clickHandler={toggleCarrier.bind(this)} clickValue={10}/>

                    <line x1={498} y1={258} x2={582} y2={258} style={{stroke:"#5b9bd5", strokeWidth:2.22}}/>
                    <ArrowRight x={582} y={258}/>
                    <Disk x={600} y={258} text={"DSP-2"}/>
                    <line x1={771} y1={258} x2={882} y2={258} style={{stroke:"#5b9bd5", strokeWidth:2.22}}/>
                    <ArrowRight x={882} y={258}/>
                    <Rectangle x={900} y={258} text={"Analog carrier"}/>
                    <line x1={1058} y1={258} x2={1098} y2={258} style={{stroke:"#5b9bd5", strokeWidth:2.22}}/>
                    <line x1={1098} y1={213} x2={1098} y2={303} style={{stroke:"#5b9bd5", strokeWidth:2.22}}/>

                    <line x1={1098} y1={213} x2={1138} y2={213} style={{stroke:"#5b9bd5", strokeWidth:2.22}}/>
                    <ArrowRight x={1138} y={213}/>
                    <Ellipse x={1156} y={213} text={"7.38"} color={this.carrierColor(11)} clickHandler={toggleCarrier.bind(this)} clickValue={11}/>

                    <line x1={1098} y1={303} x2={1138} y2={303} style={{stroke:"#5b9bd5", strokeWidth:2.22}}/>
                    <ArrowRight x={1138} y={303}/>
                    <Ellipse x={1156} y={303} text={"7.56"} color={this.carrierColor(12)} clickHandler={toggleCarrier.bind(this)} clickValue={12}/>

                    <line x1={498} y1={438} x2={582} y2={438} style={{stroke:"#5b9bd5", strokeWidth:2.22}}/>
                    <ArrowRight x={582} y={438}/>
                    <Disk x={600} y={438} text={"DSP-3"}/>
                    <line x1={771} y1={438} x2={882} y2={438} style={{stroke:"#5b9bd5", strokeWidth:2.22}}/>
                    <ArrowRight x={882} y={438}/>
                    <Rectangle x={900} y={438} text={"Analog carrier"}/>
                    <line x1={1058} y1={438} x2={1098} y2={438} style={{stroke:"#5b9bd5", strokeWidth:2.22}}/>
                    <line x1={1098} y1={393} x2={1098} y2={483} style={{stroke:"#5b9bd5", strokeWidth:2.22}}/>

                    <line x1={1098} y1={393} x2={1138} y2={393} style={{stroke:"#5b9bd5", strokeWidth:2.22}}/>
                    <ArrowRight x={1138} y={393}/>
                    <Ellipse x={1156} y={393} text={"7.74"} color={this.carrierColor(13)} clickHandler={toggleCarrier.bind(this)} clickValue={13}/>

                    <line x1={1098} y1={483} x2={1138} y2={483} style={{stroke:"#5b9bd5", strokeWidth:2.22}}/>
                    <ArrowRight x={1138} y={483}/>
                    <Ellipse x={1156} y={483} text={"7.92"} color={this.carrierColor(14)} clickHandler={toggleCarrier.bind(this)} clickValue={14}/>

                    <line x1={498} y1={633} x2={582} y2={633} style={{stroke:"#5b9bd5", strokeWidth:2.22}}/>
                    <ArrowRight x={582} y={633}/>
                    <Disk x={600} y={633} text={"DSP-4"}/>
                    <line x1={771} y1={633} x2={882} y2={633} style={{stroke:"#5b9bd5", strokeWidth:2.22}}/>
                    <ArrowRight x={882} y={633}/>
                    <Rectangle x={900} y={633} text={"Nicam"}/>
                    <line x1={1058} y1={633} x2={1200} y2={633} style={{stroke:"#5b9bd5", strokeWidth:2.22}}/>
                    <ArrowRight x={1138} y={633}/>
                    <Ellipse x={1156} y={633} text={"5.85"} color={this.carrierColor(15)} clickHandler={toggleCarrier.bind(this)} clickValue={15}/>

                    <line x1={498} y1={753} x2={582} y2={753} style={{stroke:"#5b9bd5", strokeWidth:2.22}}/>
                    <ArrowRight x={582} y={753}/>
                    <Disk x={600} y={753} text={"DSP-5"}/>
                    <line x1={771} y1={753} x2={882} y2={753} style={{stroke:"#5b9bd5", strokeWidth:2.22}}/>
                    <ArrowRight x={882} y={753}/>
                    <Rectangle x={900} y={753} text={"Nicam"}/>
                    <line x1={1058} y1={753} x2={1200} y2={753} style={{stroke:"#5b9bd5", strokeWidth:2.22}}/>
                    <ArrowRight x={1138} y={753}/>
                    <line x1={498} y1={78} x2={498} y2={753} style={{stroke:"#5b9bd5", strokeWidth:2.22}}/>
                    <Ellipse x={1156} y={753} text={"6.552"} color={this.carrierColor(16)} clickHandler={toggleCarrier.bind(this)} clickValue={16}/>
                </g>
            </svg>
        )
    }
}