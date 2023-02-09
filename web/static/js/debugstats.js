
var cy_debug = cytoscape({
  container: document.getElementById('pipeline-canvas-debug'), // container to render in
  elements: translatePipeline(debug_pipeline_data),
  autoungrabify: true,
  layout: { 
    name: 'klay',
    nodeDimensionsIncludeLabels: true,
    klay: {
      edgeSpacingFactor: 1,
      direction: 'RIGHT'
    }

  },
  style: [ // the stylesheet for the graph
    {
      selector: 'node',
      style: {
        'background-color': '#666',
        'label': 'data(label)',
        'width': '200px',
        'height': '60px',
        'shape': 'roundrectangle',
        "text-valign": "center",
        "text-halign": "center",
        "text-wrap": "wrap"
      }
    },
    {
      selector: 'edge',
      style: {
        'label': 'data(label)',
        "text-valign": "top",
        "text-halign": "center",
      }
    },
    {
      "selector": ".single-threaded",
      "style": {
        'background-color': '#46a046',
      }
    },
    {
      "selector": ".multi-threaded",
      "style": {
        'background-color': '#80a643',
        'height': '120px',
      }
    },
    {
      "selector": ".synchronous",
      "style": {
        'background-color': '#01673a',
      }
    },
  ]
});
cy_debug.minZoom(cy_debug.zoom()/2);
