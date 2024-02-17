// src/components/PieChart.js
import React from "react";
import { Doughnut } from "react-chartjs-2";

function RingChart({ chartData, options }) {
    return (
        <div className="chart-container">
        <h2 style={{ textAlign: "center" }}>RAM</h2>
        <Doughnut data={chartData} options={options}/>
        </div>
    );
}
export default RingChart;