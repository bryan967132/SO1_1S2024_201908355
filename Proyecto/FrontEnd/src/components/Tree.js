import React, { useEffect } from 'react';
import * as d3Graphviz from 'd3-graphviz';

export default function App({dotCode}) {
    useEffect(() => {
        d3Graphviz.graphviz("#graphProcess")
            .scale(0.5)
            .width(document.getElementById('graphProcess').clientWidth)
            .height(1000)
            .renderDot(dotCode)
    }, [dotCode]);

    return (
        <div>
            <div id="graphProcess"></div>
        </div>
    );
}