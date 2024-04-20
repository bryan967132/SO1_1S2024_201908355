import React, { useState } from 'react';
import './App.css'; // Archivo de estilos CSS

function App() {

    const [showData, setShowData] = useState(false)
    const [carnet, setCarnet] = useState("")
    const [nombre, setNombre] = useState("")

    const handleButton = async () => {
        try {
            if(!showData) {
                const response = await fetch('http://127.0.0.1:8000/data'); // La URL '/data' asume que tu aplicación está configurada para manejar esta ruta correctamente
                if (!response.ok) {
                    alert('Error al obtener los datos');
                }
                const jsonData = await response.json();
                setCarnet(jsonData.carnet)
                setNombre(jsonData.nombre)
            }
            setShowData(!showData)
        } catch (error) {
            alert('Error:', error);
        }
    }

    return (
        <div className="container">
            <button className="centered-button" onClick={handleButton}>
                {showData ? 'Ocultar Data' : 'Mostrar Data'}
            </button>
            {showData && <p>{carnet}</p>}
            {showData && <p>{nombre}</p>}
        </div>
    );
}

export default App;