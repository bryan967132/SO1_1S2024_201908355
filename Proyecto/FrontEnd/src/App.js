// src/App.js
import React, { useState, useEffect } from 'react';
import Pie from './components/Pie';
import Line from './components/Line';
import Tree from './components/Tree'
import State from './components/States'
import './App.css';

export default function App() {
    const [chartPie, setChartPie] = useState(true);
    const [text, setText] = useState("Históricos");

    const toggleCharts = () => {
        setChartPie(!chartPie);
        if(chartPie) {
            getHistory()
        }
        setText(text === 'Histórico' ? 'Uso' : 'Histórico')
    };

    const [dataRAM, setDataRAM] = useState({
        labels: ['Usado: 50.00 %', 'Disponible: 50.00 %'],
        datasets: [
            {
                label: 'Porcentaje',
                data: [50, 50],
                backgroundColor: [
                    'rgba(255, 99, 132, 0.2)',
                    'rgba(75, 192, 192, 0.2)',
                ],
                borderColor: [
                    'rgba(255, 99, 132, 1)',
                    'rgba(75, 192, 192, 1)',
                ],
                borderWidth: 1,
            },
        ],
    })

    const [dataCPU, setDataCPU] = useState({
        labels: ['Usado: 50.00 %', 'Disponible: 50.00 %'],
        datasets: [
            {
                label: 'Porcentaje',
                data: [50, 50],
                backgroundColor: [
                    'rgba(255, 99, 132, 0.2)',
                    'rgba(75, 192, 192, 0.2)',
                ],
                borderColor: [
                    'rgba(255, 99, 132, 1)',
                    'rgba(75, 192, 192, 1)',
                ],
                borderWidth: 1,
            },
        ],
    })

    const [historyRAM, setHistoryRAM] = useState({
        labels: [],
        datasets: [
            {
                label: 'Disponible',
                data: [],
                borderColor: 'rgb(75, 192, 192)',
                backgroundColor: 'rgba(75, 192, 192, 0.5)',
            },
            {
                label: 'En Uso',
                data: [],
                borderColor: 'rgb(255, 99, 132)',
                backgroundColor: 'rgba(255, 99, 132, 0.5)',
            },
        ],
    })


    const [historyCPU, setHistoryCPU] = useState({
        labels: [],
        datasets: [
            {
                label: 'Disponible',
                data: [],
                borderColor: 'rgb(75, 192, 192)',
                backgroundColor: 'rgba(75, 192, 192, 0.5)',
            },
            {
                label: 'En Uso',
                data: [],
                borderColor: 'rgb(255, 99, 132)',
                backgroundColor: 'rgba(255, 99, 132, 0.5)',
            },
        ],
    })

    const getData = async () => {
        const response = await fetch(`/back/cpuram`)
        if(!response.ok) {
            throw new Error('Network response was not ok');
        }
        const jsonData = await response.json();
        setDataRAM({
            labels: [`Usado: ${jsonData.percentusedRAM.toFixed(2)} %`, `Disponible: ${(100 - jsonData.percentusedRAM).toFixed(2)} %`],
            datasets: [
                {
                    label: 'Porcentaje',
                    data: [jsonData.percentusedRAM, 100 - jsonData.percentusedRAM],
                    backgroundColor: [
                        'rgba(255, 99, 132, 0.2)',
                        'rgba(75, 192, 192, 0.2)',
                    ],
                    borderColor: [
                        'rgba(255, 99, 132, 1)',
                        'rgba(75, 192, 192, 1)',
                    ],
                    borderWidth: 1,
                },
            ],
        })
        setDataCPU({
            labels: [`Usado: ${jsonData.percentusedCPU.toFixed(2)} %`, `Disponible: ${(100 - jsonData.percentusedCPU).toFixed(2)} %`],
            datasets: [
                {
                    label: 'Porcentaje',
                    data: [jsonData.percentusedCPU, 100 - jsonData.percentusedCPU],
                    backgroundColor: [
                        'rgba(255, 99, 132, 0.2)',
                        'rgba(75, 192, 192, 0.2)',
                    ],
                    borderColor: [
                        'rgba(255, 99, 132, 1)',
                        'rgba(75, 192, 192, 1)',
                    ],
                    borderWidth: 1,
                },
            ],
        })
    }

    const getHistory = async() => {
        const response = await fetch(`/back/history`)
        if(!response.ok) {
            throw new Error('Network response was not ok');
        }
        const jsonData = await response.json();
        var labels = []
        var datap = undefined
        var dataDisponible = []
        var dataUsado = []
        var i = 0
        for(i = 0; i < jsonData.ram.length; i ++) {
            datap = jsonData.ram[i]
            labels.push(datap.tiempo.split(' ')[1])
            dataDisponible.push(datap.disponible)
            dataUsado.push(datap.usado)
        }

        setHistoryRAM({
            labels: labels,
            datasets: [
                {
                    label: 'Disponible',
                    data: dataDisponible,
                    borderColor: 'rgb(75, 192, 192)',
                    backgroundColor: 'rgba(75, 192, 192, 0.5)',
                },
                {
                    label: 'En Uso',
                    data: dataUsado,
                    borderColor: 'rgb(255, 99, 132)',
                    backgroundColor: 'rgba(255, 99, 132, 0.5)',
                },
            ],
        })

        labels = []
        datap = undefined
        dataDisponible = []
        dataUsado = []
        for(i = 0; i < jsonData.cpu.length; i ++) {
            datap = jsonData.cpu[i]
            labels.push(datap.tiempo.split(' ')[1])
            dataDisponible.push(datap.disponible)
            dataUsado.push(datap.usado)
        }

        setHistoryCPU({
            labels: labels,
            datasets: [
                {
                    label: 'Disponible',
                    data: dataDisponible,
                    borderColor: 'rgb(75, 192, 192)',
                    backgroundColor: 'rgba(75, 192, 192, 0.5)',
                },
                {
                    label: 'En Uso',
                    data: dataUsado,
                    borderColor: 'rgb(255, 99, 132)',
                    backgroundColor: 'rgba(255, 99, 132, 0.5)',
                },
            ],
        })
    }

    useEffect(() => {
        const intervalUso = setInterval(async () => {
            if(chartPie) {
                getData()
            }
        }, 500)

        const setHistory = async () => {
            const body = {
                usadoram: parseFloat(dataRAM.datasets[0].data[0].toFixed(2)),
                disponibleram: parseFloat(dataRAM.datasets[0].data[1].toFixed(2)),
                usadocpu: parseFloat(dataCPU.datasets[0].data[0].toFixed(2)),
                disponiblecpu: parseFloat(dataCPU.datasets[0].data[1].toFixed(2)),
            }
            await fetch(`/back/inscpuram`, {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json'
                },
                body: JSON.stringify(body)
            })
        }

        const intervalHistory = setInterval(async () => {
            setHistory()
        }, 1000 * 240)

        return () => {
            clearInterval(intervalUso)
            clearInterval(intervalHistory)
        }
    }, [chartPie, dataRAM, dataCPU])

    const [selectedOption, setSelectedOption] = useState('Seleccionar');
    const [options, setOptions] = useState([]);
    const [dotCode, setDotCode] = useState('digraph G {}')

    const getPIDs = async () => {
        var response = await fetch(`/back/pids`)
        if(!response.ok) {
            throw new Error('Network response was not ok');
        }
        var jsonData = await response.json()
        setOptions(jsonData.pids)
    }

    const getProc = async (option) => {
        var response = await fetch(`/back/proc/${option}`)
        if(!response.ok) {
            throw new Error('Network response was not ok');
        }
        var jsonData = await response.json()

        var c = undefined

        var dot = `\n\tn_root[label = "${jsonData.proc.name}\\nPID = ${jsonData.proc.pid}"];`

        for(var i = 0; i < jsonData.proc.child.length; i ++) {
            c = jsonData.proc.child[i]
            dot += `\n\tn_${i}[label = "${c.name}\\nPID = ${c.pid}"];`
            dot += `\n\tn_root -> n_${i};`
        }

        setDotCode(`digraph G {\n\tlayout = circo;\n\tedge[color="#1975ff"];\n\tnode[fontcolor=white, color="#163c75", style=filled, fillcolor="#3a78d6"];${dot}\n}`)
    }

    const [visibleTree, setVisibleTree] = useState(false)

    const handleSelectChange = (event) => {
        setSelectedOption(event.target.value);
        getProc(event.target.value)
        setVisibleTree(true)
    }

    const [dotCodeStates, setDotCodeStates] = useState('digraph G {}')

    const generateGraphState = (data) => {
        const ready = data.ready
        const running = data.running
        const terminated = data.terminated

        var isTerminated = false

        var dot = '\n\tn_new[label = "New"];'
        dot += `\n\tn_ready[label = "Ready"${ready.status === 'Current' ? ' color="#0d3b18" fillcolor="#46b860"' : ''}];`
        dot += `\n\tn_new -> n_ready;`
        dot += `\n\tn_running[label = "Running"${running.status === 'Current' ? ' color="#0d3b18" fillcolor="#46b860"' : ''}];`
        dot += `\n\tn_ready -> n_running${running.status === 'Current' ? ' [color="#46b860" dir=front]' : ''};`
        if(ready.to.includes('terminated')) {
            if(!isTerminated) {
                dot += `\n\tn_terminated[label = "Terminated"${terminated.status === 'Current' ? ' color="#0d3b18" fillcolor="#46b860"' : ''}];`
                isTerminated = true
            }
            dot += `\n\tn_ready -> n_terminated${ready.status === 'Current' ? ' [color="#46b860" dir=front]' : ''};`
        }
        if(running.to.length) {
            if(running.to.includes('ready')) {
                dot += `\n\tn_running -> n_ready${ready.status === 'Current' ? ' [color="#46b860" dir=front]' : ''};`
            }
            if(running.to.includes('terminated')) {
                if(!isTerminated) {
                    dot += `\n\tn_terminated[label = "Terminated"${terminated.status === 'Current' ? ' color="#0d3b18" fillcolor="#46b860"' : ''}];`
                    isTerminated = true
                }
                dot += `\n\tn_running -> n_terminated${ready.status === 'Current' ? ' [color="#46b860" dir=front]' : ''};`
            }
        }

        return `digraph G {\n\tedge[dir=none];\n\tnode[fontcolor=white color="#163c75" style=filled fillcolor="#5e9cff"];\n\trankdir=LR;${dot}\n}`
    }

    const [pidCurrent, setPidCurrent] = useState(-1)

    // eslint-disable-next-line
    const getPIDCurrent = async () => {
        const response = await fetch(`/back/thereisproc`)
        if(!response.ok) {
            throw new Error('Network response was not ok');
        }
        const jsonData = await response.json()
        if(jsonData.status) {
            setPidCurrent(jsonData.PID)
            setDotCodeStates(generateGraphState(jsonData.graph))
        }
    }

    const newProcess = async () => {
        const response = await fetch(`/back/start`)
        if(!response.ok) {
            throw new Error('Network response was not ok');
        }
    }

    const stopProcess = async () => {
        const response = await fetch(`/back/stop/${pidCurrent}`)
        if(!response.ok) {
            throw new Error('Network response was not ok');
        }
    }

    const resumeProcess = async () => {
        const response = await fetch(`/back/resume/${pidCurrent}`)
        if(!response.ok) {
            throw new Error('Network response was not ok');
        }
    }

    const killProcess = async () => {
        const response = await fetch(`/back/kill/${pidCurrent}`)
        if(!response.ok) {
            throw new Error('Network response was not ok');
        }
        const jsonData = await response.json()
        if(jsonData.status === 'Success') {
            setDotCodeStates(generateGraphState(jsonData.graph))
        }
    }

    useEffect(() => {
        getPIDs()
    }, [])

    useEffect(() => {
        getPIDCurrent()
    }, [getPIDCurrent])

    return (
        <div className="App">
            <div>
                <h1>Monitoreo</h1>
                <button onClick={toggleCharts} className='custom-button'>{text}</button>
                <div className="monitor">
                    {chartPie ? (
                        <>
                            <div className="chart-container-pie">
                                <Pie data={dataRAM} title={'RAM'}/>
                            </div>
                            <div className="chart-container-pie">
                                <Pie data={dataCPU} title={'CPU'}/>
                            </div>
                        </>
                    ) : (
                        <>
                            <div className="chart-container-line">
                                <Line data={historyRAM} title={'RAM'}/>
                            </div>
                            <div className="chart-container-line">
                                <Line data={historyCPU} title={'CPU'}/>
                            </div>
                        </>
                    )}
                </div>
            </div>
            <div>
                <h1 className="titles2">Procesos</h1>
                <p>Seleccione un PID</p>
                <select value={selectedOption} onChange={handleSelectChange}>
                    <option value="" disabled hidden>Seleccionar</option>
                    {options.map((option, index) => (
                        <option key={index} value={option}>{option}</option>
                    ))}
                </select>
                <div className="chart-row-1">
                    <div className="chart-container-tree">
                        { visibleTree ? <Tree dotCode={dotCode}/> : <></> }
                    </div>
                </div>
            </div>
            <div className="last">
                <h1>Diagrama de Estados</h1>
                <div className="button-container">
                    <h4 className='pid-area'>PID: { pidCurrent > 0 ? pidCurrent : 'Nulo' }</h4>
                    <button className="state-button color1" onClick={newProcess}>New</button>
                    <button className="state-button color2" onClick={stopProcess}>Stop</button>
                    <button className="state-button color3" onClick={resumeProcess}>Ready</button>
                    <button className="state-button color4" onClick={killProcess}>Kill</button>
                </div>
                <div className="chart-row-1">
                    <div className="chart-container-tree">
                        <State dotCode={dotCodeStates}/>
                    </div>
                </div>
            </div>
        </div>
    );
}