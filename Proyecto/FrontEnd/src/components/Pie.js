import React from 'react';
import { Chart as ChartJS, ArcElement, Tooltip, Legend } from 'chart.js';
import { Pie } from 'react-chartjs-2';

ChartJS.register(ArcElement, Tooltip, Legend);

export default function App({data, title}) {
    return (
        <div>
            <h2 style={{ textAlign: "center" }}>{title}</h2>
            <Pie data={data}/>
        </div>
    )
}