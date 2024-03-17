import React, { useEffect } from 'react';
import * as d3Graphviz from 'd3-graphviz';

export default function App({dotCode}) {
    useEffect(() => {
        d3Graphviz.graphviz("#graphStates")
            .width(document.getElementById('graphStates').clientWidth)
            .scale(1.5)
            .height(500)
            .renderDot(dotCode);
    }, [dotCode]);

    return (
        <div className="App">
            <div id="graphStates"></div>
        </div>
    );
}