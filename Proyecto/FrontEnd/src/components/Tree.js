import React from "react";
import Tree from 'react-d3-tree';
import './custom-tree.css'

export default function App({data}) {
    return (
        <div id="treeWrapper" style={{ height: '40em' }}>
            <Tree
                data={data}
                separation={{
                    siblings: 2, // Separación entre nodos hermanos
                    nonSiblings: 2 // Separación entre nodos no hermanos
                }}
                rootNodeClassName="node__root"
                branchNodeClassName="node__branch"
                leafNodeClassName="node__leaf"
                orientation="vertical"
            />
        </div>
    )
}