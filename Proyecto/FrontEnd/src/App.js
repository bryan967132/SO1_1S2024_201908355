// src/App.js
import React, { useState, useEffect } from 'react';
import Pie from './components/Pie';
import Line from './components/Line';
import Tree from './components/Tree'
import './App.css';

const HOSTBACK = 'http://localhost:8000'

export default function App() {
    const [chartPie, setChartPie] = useState(true);
    const [text, setText] = useState("Histórico");

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
        const response = await fetch(`${HOSTBACK}/cpuram`)
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
        const response = await fetch(`${HOSTBACK}/history`)
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
            await fetch(`${HOSTBACK}/inscpuram`, {
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

    const [dataT, setDataT] = useState({})

    const getPIDs = async () => {
        var response = await fetch(`${HOSTBACK}/pids`)
        if(!response.ok) {
            throw new Error('Network response was not ok');
        }
        var jsonData = await response.json()
        setOptions(jsonData.pids)
    }
    
    const getProc = async (option) => {
        var response = await fetch(`${HOSTBACK}/proc/${option}`)
        if(!response.ok) {
            throw new Error('Network response was not ok');
        }
        var jsonData = await response.json()

        var child = []
        var c = undefined

        for(var i = 0; i < jsonData.proc.child.length; i ++) {
            c = jsonData.proc.child[i]
            child.push({name: c.name, attributes: {pid: c.pid}})
        }

        setDataT({
            name: jsonData.proc.name,
            attributes: {pid: jsonData.proc.pid},
            children: child,
        })
    }

    const handleSelectChange = (event) => {
        setSelectedOption(event.target.value);
        getProc(event.target.value)
    }

    useEffect(() => {
        getPIDs()
    }, [])

    return (
        <div className="App">
            <h1>Monitoreo</h1>
            <button onClick={toggleCharts} className='custom-button'>{text}</button>
            <div>
                {chartPie ? (
                    <>
                        <div className="chart-container">
                            <Pie data={dataRAM} title={'RAM'}/>
                        </div>
                        <div className="chart-container">
                            <Pie data={dataCPU} title={'CPU'}/>
                        </div>
                    </>
                ) : (
                    <>
                        <div className="chart-row">
                            <div className="chart-container-line">
                                <Line data={historyRAM} title={'RAM'}/>
                            </div>
                            <div className="chart-container-line">
                                <Line data={historyCPU} title={'CPU'}/>
                            </div>
                        </div>
                    </>
                )}
            </div>
            <h1 className="titles2">Procesos</h1>
            <p>Seleccione un PID</p>
            <select value={selectedOption} onChange={handleSelectChange}>
                <option value="" disabled hidden>Seleccionar</option>
                {options.map((option, index) => (
                    <option key={index} value={option}>{option}</option>
                ))}
            </select>
            <div className="chart-row">
                <div className="chart-container-tree">
                    <Tree data={dataT}/>
                </div>
            </div>
        </div>
    );
}