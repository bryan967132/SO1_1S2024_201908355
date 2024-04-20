import './App.css';
import {useState, useEffect} from 'react';
import {CatModule} from "../wailsjs/go/main/App";
import Chart from "chart.js/auto";
import { CategoryScale } from "chart.js";
import RingChart from "./Components/RingChart";

Chart.register(CategoryScale);

function App() {
    const [chartData, setChartData] = useState({
        labels: ['Libre', 'En Uso'],
        datasets: [
            {
                label: 'RAM',
                data: [50, 50],
                backgroundColor: [
                    'rgba(54, 162, 235, 0.6)',
                    'rgba(255, 99, 132, 0.6)'
                ],
                borderWidth: 1
            }
        ]
    });

    const [options, setOptions] = useState({
        plugins: {
            title: {
                display: true,
                text: "50 % En Uso"
            }
        }
    })

    const updateResultText = (result) => {
        var data = JSON.parse(result)
        var percent = (data.freeram / data.totalram * 100).toFixed(2)
        setChartData({
            labels: ['Libre', 'En Uso'],
            datasets: [
                {
                    label: 'RAM',
                    data: [percent, 100 - percent],
                    backgroundColor: [
                        'rgba(54, 162, 235, 0.6)',
                        'rgba(255, 99, 132, 0.6)'
                    ],
                    borderWidth: 1
                }
            ]
        })
        setOptions({
            plugins: {
                title: {
                    display: true,
                    text: 100 - percent + " % En Uso"
                }
            }
        })
    };

    function catModule() {
        CatModule().then(updateResultText);
    }

    useEffect(() => {
        const interval = setInterval(() => {
            catModule()
        }, 500);

        return () => clearInterval(interval);
    }, []);

    return (
        <div id="App">
            <RingChart chartData={chartData} options={options}/>
        </div>
    )
}

export default App
