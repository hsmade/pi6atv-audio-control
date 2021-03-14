import React from 'react';

export class Ellipse extends React.Component {
    render() {
        const x = parseInt(this.props.x)
        const y = parseInt(this.props.y) - 33
        return (
            <svg height={66} width={159} x={x} y={y}>
            <path
                fill="#4f88bb"
                stroke="#41719c"
                strokeWidth="1.67"
                d="M32 64h94.5c17.32 0 31.5-14.1 31.5-31.5s-14.18-31.5-31.5-31.5h-94.5c-17.3 0-31.48 14.1-31.48 31.5s14.17 31.5 31.5 31.5z"
            />
                <text x="50%" y="50%" dominantBaseline="middle" textAnchor="middle" fill="white">{this.props.text}</text>
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

export default class Diagram extends React.Component {
    clickme() {
        console.log("clicked")
    }
    render() {
        return (
            <svg
                width="1871"
                height="1323"
            >
                <g>
                    <Ellipse x={0} y={33} text={"DSP-1"}/>
                    <line x1="159" y1="33" x2="200" y2="33" style={{stroke:"#5b9bd5", strokeWidth:2.22}}/>
                    <Ellipse x={0} y={133} text={"DSP-2"}/>
                    <line x1="159" y1="133" x2="200" y2="133" style={{stroke:"#5b9bd5", strokeWidth:2.22}}/>
                    <Ellipse x={0} y={233} text={"DSP-3"}/>
                    <line x1="159" y1="233" x2="200" y2="233" style={{stroke:"#5b9bd5", strokeWidth:2.22}}/>
                    <Ellipse x={0} y={333} text={"DSP-4"}/>
                    <line x1="159" y1="333" x2="400" y2="333" style={{stroke:"#5b9bd5", strokeWidth:2.22}}/>
                    <Ellipse x={0} y={433} text={"DSP-5"}/>
                    <line x1="159" y1="433" x2="200" y2="433" style={{stroke:"#5b9bd5", strokeWidth:2.22}}/>
                    <line x1="200" y1="33" x2="200" y2="433" style={{stroke:"#5b9bd5", strokeWidth:2.22}}/>

                    <Ellipse x={0} y={633} text={"Program"}/>
                    <line x1="159" y1="633" x2="379" y2="633" style={{stroke:"#5b9bd5", strokeWidth:2.22}}/>
                    <Ellipse x={0} y={733} text={"Reset"}/>
                    <line x1="159" y1="733" x2="379" y2="733" style={{stroke:"#5b9bd5", strokeWidth:2.22}}/>

                    <ArrowRight x={282} y={333}/>
                    <Ruit x={300} y={333} text={"Besturing"}/>
                    <line x1="458" y1="333" x2="498" y2="333" style={{stroke:"#5b9bd5", strokeWidth:2.22}}/>
                    <ArrowUp x={372} y={386}/>
                    <line x1="379" y1="386" x2="379" y2="733" style={{stroke:"#5b9bd5", strokeWidth:2.22}}/>

                    <ArrowRight x={1182} y={33}/>
                    <Ellipse x={1200} y={33} text={"7.02"}/>
                    <line x1={498} y1={78} x2={582} y2={78} style={{stroke:"#5b9bd5", strokeWidth:2.22}}/>
                    <ArrowRight x={582} y={78}/>
                    <Disk x={600} y={78} text={"DSP-1"}/>
                    <ArrowRight x={882} y={78}/>
                    <Rectangle x={900} y={78} text={"Analog carrier"}/>
                    <ArrowRight x={1182} y={123}/>
                    <Ellipse x={1200} y={123} text={"7.20"}/>

                    <ArrowRight x={1182} y={213}/>
                    <Ellipse x={1200} y={213} text={"7.38"}/>
                    <line x1={498} y1={258} x2={582} y2={258} style={{stroke:"#5b9bd5", strokeWidth:2.22}}/>
                    <ArrowRight x={582} y={258}/>
                    <Disk x={600} y={258} text={"DSP-2"}/>
                    <ArrowRight x={882} y={258}/>
                    <Rectangle x={900} y={258} text={"Analog carrier"}/>
                    <ArrowRight x={1182} y={303}/>
                    <Ellipse x={1200} y={303} text={"7.56"}/>

                    <ArrowRight x={1182} y={393}/>
                    <Ellipse x={1200} y={393} text={"7.74"}/>
                    <line x1={498} y1={438} x2={582} y2={438} style={{stroke:"#5b9bd5", strokeWidth:2.22}}/>
                    <ArrowRight x={582} y={438}/>
                    <Disk x={600} y={438} text={"DSP-3"}/>
                    <ArrowRight x={882} y={438}/>
                    <Rectangle x={900} y={438} text={"Analog carrier"}/>
                    <ArrowRight x={1182} y={483}/>
                    <Ellipse x={1200} y={483} text={"7.92"}/>

                    <line x1={498} y1={633} x2={582} y2={633} style={{stroke:"#5b9bd5", strokeWidth:2.22}}/>
                    <ArrowRight x={582} y={633}/>
                    <Disk x={600} y={633} text={"DSP-4"}/>
                    <ArrowRight x={882} y={633}/>
                    <Rectangle x={900} y={633} text={"Nicam"}/>
                    <ArrowRight x={1182} y={633}/>
                    <Ellipse x={1200} y={633} text={"5.85"}/>

                    <line x1={498} y1={753} x2={582} y2={753} style={{stroke:"#5b9bd5", strokeWidth:2.22}}/>
                    <ArrowRight x={582} y={753}/>
                    <Disk x={600} y={753} text={"DSP-5"}/>
                    <ArrowRight x={882} y={753}/>
                    <Rectangle x={900} y={753} text={"Nicam"}/>
                    <ArrowRight x={1182} y={753}/>
                    <Ellipse x={1200} y={753} text={"6.552"}/>
                    <line x1={498} y1={78} x2={498} y2={753} style={{stroke:"#5b9bd5", strokeWidth:2.22}}/>
                </g>
            </svg>
        )
    }
}