import React from 'react';
import {fill, stroke} from './config'

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
        let text2 = <div/>
        let text = <text
            x="50%" y="50%"
            dominantBaseline="middle" textAnchor="middle"
            fill="white"
            >{this.props.text}</text>
        if (this.props.text2 !== undefined) {
            text = <text
                x="50%" y="40%"
                dominantBaseline="middle" textAnchor="middle"
                fill="white"
            >{this.props.text}</text>
            text2 = <text
                x="50%" y="60%"
                dominantBaseline="middle" textAnchor="middle"
                fill="white"
            >{this.props.text2}</text>

        }
        return (
            <svg height={80} width={120} x={x} y={y}>
                <polygon
                    points={"1,40 60,1 119,40 60,79 1,40"}
                    fill={this.props.color}
                    stroke={stroke} strokeWidth={2}
                />
                {text}
                {text2}
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
