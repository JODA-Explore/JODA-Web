function translatePipeline(pipeline) {
  var elts = [];

  pipeline.Tasks.forEach(function (task) {
    var async_class = "";
    switch (task.Async) {
      case "SingleThreaded":
        async_class = "single-threaded";
        break;
      case "MultiThreaded":
        async_class = "multi-threaded";
        break;
      case "Synchronous":
        async_class = "synchronous";
        break;
    }
    elts.push({
      group: 'nodes',
      data: { id: task.Num, name: task.Name, async: task.Async, runtime: task.Runtime, taskcount: task.TaskCount, desc: task.Description, label: '(' + task.Num + ')\n' + task.Name},
      classes: async_class,
    });
  })

  pipeline.Connections.forEach(function (con) {
    con.From.forEach(function (from) {
      con.To.forEach(function (to) {
        elts.push({
          group: 'edges',
          data: {
            id: from + '-' + to,
            source: from,
            target: to,
            label: con.Throughput,
          }
        });
      })
    })
  });

  return elts;
}

var currentPopper = null;
var updatePopper = function () {
  if (currentPopper != null) {
    currentPopper.update();
  }

};

function removePipelinePopup() {
  if (currentPopper != null) {
    var pop = document.getElementById("pipeline-popup");
    pop.remove();
    currentPopper = null;
  }
}

function createPopupCard(node) {
  var card = document.createElement('div');
  card.classList.add('w3-card-4');

  var header = document.createElement('header');
  header.classList.add('w3-container');
  header.classList.add('w3-dark-grey');
  header.innerHTML = '<h3>' + node.id() + ' - ' + node.data('name') + '</h3>';
  card.appendChild(header);

  var body = document.createElement('div');
  body.classList.add('w3-container');
  body.classList.add('w3-light-grey');

  var desc = "";
  if(node.data('desc') != null) {
    desc = '<hr class="darkhr"> ' + node.data('desc');
  }
  body.innerHTML = '<p><b>Runtime:</b> ' + node.data('runtime') + '</p>' + '<p><b>#Threads:</b> ' + node.data('taskcount') + '</p>' + desc;
  card.appendChild(body);

  var footer = document.createElement('footer');
  footer.classList.add('w3-container');
  footer.classList.add('w3-dark-grey');
  footer.innerHTML = '<h4>' + node.data('async') + '</h4>'
  card.appendChild(footer);

  return card;
}

var createPopupFromNode = function (node) {
  removePipelinePopup();
  let div = document.createElement('div');
  div.setAttribute('id', 'pipeline-popup');
  div.classList.add('pipeline-popup-div');

  div.appendChild(createPopupCard(node));

  document.body.appendChild(div);
  return div;
}



var cy = cytoscape({
  container: document.getElementById('pipeline-canvas'), // container to render in
  elements: translatePipeline(pipeline_data),
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
        "text-wrap": "wrap",
        "text-max-width": "190px"
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
cy.minZoom(cy.zoom()/2);


cy.on('tap', 'node', function(evt){
  var node = evt.target;
  currentPopper = node.popper({
      content: function(){ return createPopupFromNode(node); },
      popper: {} // options
  });
  updatePopper();
});

cy.on('tap', function(evt){
  if( evt.target === cy ){
    removePipelinePopup();
  }
});

cy.on('pan zoom resize', updatePopper);