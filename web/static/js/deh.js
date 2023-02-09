
function fetchApi(endpoint, id, dataset, docId) {
    if (document.getElementById(docId).innerHTML === '') {
        // so that the loading content will not flash.
        document.getElementById(docId).innerHTML = '<br><br><br><br><br><br><br><br><br><br>'
        const url = '/api/' + endpoint + '?id=' + id + '&dataset=' + dataset
        fetch(url)
            .then(response => response.text())
            .then((response) => {
                document.getElementById(docId).innerHTML = response
            })
    }
}

function getSelected(id) {
    selectList = document.getElementById(id)
    const text = selectList.options[selectList.selectedIndex].text
    return text
}

// see https://stackoverflow.com/questions/11338774/serialize-form-data-to-json
function getFormData($form) {
    const unindexed_array = $form.serializeArray()
    const indexed_array = {}

    $.map(unindexed_array, function(n, i) {
        indexed_array[n.name] = n.value
    })

    indexed_array.dehMinRating = $('#dehMinRating').val()
    indexed_array.dehMaxItems = $('#dehMaxItems').val()
    indexed_array.dehMaxDiff = $('#dehMaxDiff').val()

    let parent, name, query, divId
    const dataSourceDiv = document.getElementById('data-source-div')
    const divContents = dataSourceDiv.getElementsByClassName('buttonbar-content-div')[0]
    const divs = divContents.getElementsByTagName('div')
    for (const div of divs) {
        if (div.style.display == 'block') {
            divId = div.id
            break
        }
    }
    switch (divId) {
        case 'data-source-exists':
            indexed_array.sourceType = 'exists'
            name = getSelected('data-source-exists-select-name')
            if (name === '--None--') {
                name = ''
            }
            parent = ''
            query = ''
            break
        case 'data-source-new':
            indexed_array.sourceType = 'new'
            parent = getSelected('data-source-new-select-parent')
            if (parent === '--None--') {
                parent = ''
            }
            query = document.getElementById('data-source-new-query').value
            name = document.getElementById('data-source-new-name').value
            break
        case 'data-source-load':
            indexed_array.sourceType = 'load'
            name = document.getElementById('data-source-load-name').value
            query = document.getElementById('data-source-load-query').value
            parent = ''
            break
    }

    indexed_array.name = name
    indexed_array.parent = parent
    indexed_array.query = query
    return indexed_array
}

function handleCopyText(text) {
    const cb = navigator.clipboard
    cb.writeText(text)
}

function execQuery(query) {
    const params = new URLSearchParams('')
    params.set('query', query)
    const url = '/joda/query?' + params.toString()
    window.open(url)
}

function handleButtonBar(buttonId, buttonBarId, docId) {
    const button = document.getElementById(buttonId)
    const buttonBarDiv = document.getElementById(buttonBarId)
    const contents = buttonBarDiv.getElementsByClassName('content')
    for (i = 0; i < contents.length; i++) {
        contents[i].style.display = 'none'
    }
    const mainButtonBar = buttonBarDiv.getElementsByClassName('main-buttonbar')[0]
    const buttons = mainButtonBar.getElementsByTagName('button')
    for (i = 0; i < buttons.length; i++) {
        buttons[i].className = buttons[i].className.replace(' w3-green', '')
    }
    const currentContent = document.getElementById(docId)
    currentContent.style.display = 'block'
    button.className += ' w3-green'
}

function showStr(str) {
    if (str.length < 10) {
        return str
    }
    const array = str.match(/(.|[\r\n]){1,8}/g)
    return array.join('\n')
}

function newDoughnutChart() {
    return {
        type: 'doughnut',
        data: {
            labels: [],
            datasets: [{
                data: [],
                title: '',
                fill: false,
                borderWidth: 1,
                backgroundColor: [
                    'rgba(51, 160, 44, 1)',
                    'rgba(31, 120, 180, 1)',
                    'rgba(227, 26, 28, 1)',
                    'rgba(255, 127, 0, 1)',
                    'rgba(106, 61, 154, 1)',
                    'rgba(177, 89, 40, 1)',
                    'rgba(178, 223, 138, 1)',
                    'rgba(166, 206, 227, 1)',
                    'rgba(251, 154, 153, 1)',
                    'rgba(253, 191, 111, 1)'
                ]
            }]
        },
        plugins: [ChartDataLabels],
        options: {
            plugins: {
                legend: {
                    display: false
                },
                datalabels: {
                    labels: {
                        value: {
                            display: function(ctx) {
                                return ctx.dataset.data[ctx.dataIndex] / ctx.dataset.data[0] > 0.1
                            },
                            align: 'bottom',
                            color: 'black',
                            borderColor: 'white',
                            borderWidth: 2,
                            borderRadius: 4,
                            backgroundColor: 'white',
                            font: { size: 15 },
                            color: function(ctx) {
                                return ctx.dataset.backgroundColor
                            },
                            padding: 4,
                            font: {
                                weight: 'bold'
                            }
                        },
                        name: {
                            display: function(ctx) {
                                return ctx.chart.data.labels[ctx.dataIndex].length < 30 &&
                                    ctx.dataset.data[ctx.dataIndex] / ctx.dataset.data[0] > 0.1
                            },
                            align: 'top',
                            font: { size: 16 },
                            formatter: function(value, ctx) {
                                return showStr(ctx.chart.data.labels[ctx.dataIndex])
                            },
                            color: 'white',
                            font: {
                                weight: 'bold'
                            }
                        }
                    }

                }
            },
            title: {
                display: false
            },
            maintainAspectRatio: false,
            responsive: true
        }
    }
}

function newHorizontalBarChart() {
    return {
        type: 'bar',
        data: {
            labels: [],
            datasets: [{
                barPercentage: 0.9,
                minBarLength: 2,
                data: [],
                fill: false,
                backgroundColor: '#4CAF50',
                borderWidth: 1
            }]
        },
        plugins: [ChartDataLabels],
        options: {
            onClick: chartClickEvent,
            plugins: {
                legend: {
                    display: false
                },
                datalabels: {
                    color: '#808080',
                    align: 'end',
                    offset: 10,
                    anchor: 'end',
                    font: {
                        weight: 'bold'
                    }
                }
            },
            scales: {
                xAxes: {
                    position: 'top',
                    grid: {
                        display: false
                    }
                },
                yAxes: {
                    grid: {
                        display: false
                    }
                }
            },
            indexAxis: 'y',
            title: {
                display: false
            },
            maintainAspectRatio: false,
            responsive: true
        }
    }
}

function genChartClickFunction(dataset, id, endpoint, obj) {
    return function(evt) {
        const points = obj.getElementsAtEventForMode(evt, 'nearest', { intersect: true }, true)
        if (points.length) {
            const firstPoint = points[0]
            const label = obj.data.labels[firstPoint.index]
            const value = obj.data.datasets[firstPoint.datasetIndex].data[firstPoint.index]
        }

        const url = '/api/' + endpoint + '?id=' + id + '&dataset=' + dataset + '&path=' + label
        const response = fetch(url)
        const res = response.text()
        execQuery(res)
    }
}

function genNewChart(type, labels, data, obj, dataset, id) {
    let chartData, newChart
    switch (type) {
        case 'doughnut':
            chartData = newDoughnutChart()
            break
        case 'horizontalBar':
            chartData = newHorizontalBarChart()
            break
    }
    chartData.data.labels = labels
    chartData.data.datasets[0].data = data
    newChart = new Chart(obj, chartData)
    return chartData
}

function destoryCharts() {
    const charts = document.querySelectorAll('.coverage-doughnets')
    charts.forEach(chart => chart.destory())
}

function chartClickEvent(e) {
    const activePoints = myChart.getElementsAtEvent(e)
    const selectedIndex = activePoints[0]._index
    alert(this.data.datasets[0].data[selectedIndex])
}

function filteredListLength(list, minCount, length) {
    if (list[Math.min(length - 1, list.length - 1)] >= minCount) {
        return Math.min(length, list.length)
    }
    let startIdx = Math.min(Math.round(length / 2), list.length - 1)
    while (list[startIdx] < minCount && startIdx > 1) {
        startIdx = Math.round(startIdx / 2)
    }
    let res = startIdx
    for (let i = startIdx; i < length; i++) {
        if (list[i] >= minCount) {
            res = i
        } else {
            break
        }
    }
    return res + 1
}

async function loadDistinctValues(dataset, path, number, id) {
    const url = '/api/distinctValues' + '?number=' + number + '&dataset=' + dataset + '&path=' + path
    const response = await fetch(url)
    const distinctValues = await response.json()

    const coverageDoughnets = document.getElementById('coverage-doughnets_' + id)
    if (!coverageDoughnets.classList.contains('loaded')) {
        len = filteredListLength(distinctValues.Counts, 3, 10)
        genNewChart('doughnut', distinctValues.Values.slice(0, len), distinctValues.Counts.slice(0, len), coverageDoughnets, dataset, id)
        coverageDoughnets.style.height = '75vh'
        coverageDoughnets.classList.add('loaded')
    }

    const coverageBar = document.getElementById('coverage-bar_' + id)
    if (!coverageBar.classList.contains('loaded')) {
        len = filteredListLength(distinctValues.Counts, 3, 50)
        genNewChart('horizontalBar', distinctValues.Values.slice(0, len), distinctValues.Counts.slice(0, len), coverageBar, dataset, id)
        coverageBar.style.height = Math.max(75, len * 25) + 'px'
        coverageBar.classList.add('loaded')
    }
}

async function loadMemberFreq(dataset, path, number, id) {
    let url = '/api/allDistinctMembersCount' + '?dataset=' + dataset + '&path=' + path
    let response = await fetch(url)
    const count = await response.json()
    document.getElementById('totalNumber_' + id).innerText = count

    url = '/api/memberFreq' + '?number=' + number + '&dataset=' + dataset + '&path=' + path
    response = await fetch(url)
    memberFreq = await response.json()

    const coverageDoughnets = document.getElementById('coverage-doughnets_' + id)
    if (!coverageDoughnets.classList.contains('loaded')) {
        // len = filteredListLength(memberFreq.Counts, 3, 10)
        len = 10
        genNewChart('doughnut', memberFreq.Values.slice(0, len), memberFreq.Counts.slice(0, len), coverageDoughnets, dataset, id)
        coverageDoughnets.style.height = '75vh'
        coverageDoughnets.classList.add('loaded')
    }

    const coverageBar = document.getElementById('coverage-bar_' + id)
    if (!coverageBar.classList.contains('loaded')) {
        // len = filteredListLength(memberFreq.Counts, 3, 50)
        len = 50
        genNewChart('horizontalBar', memberFreq.Values.slice(0, len), memberFreq.Counts.slice(0, len), coverageBar, dataset, id)
        coverageBar.style.height = Math.max(75, len * 25) + 'px'
        coverageBar.classList.add('loaded')
    }
}

function handleQueryGenerator(id, dataset) {
    fetchApi('QueryGenerator', id, dataset, 'QueryGenerator_' + id)
}

function handleCopyQuery(queryContentId) {
    const q = document.getElementById(queryContentId).innerText
    handleCopyText(q)
}

function handleExecQuery(queryContentId) {
    const q = document.getElementById(queryContentId).innerText
    execQuery(q)
}

async function handleExploreQuery(queryContentId, parent, id) {
    handleButtonBar('data-source-new-btn', 'data-source-div', 'data-source-new')
    const query = document.getElementById(queryContentId).innerText
    changeSelectOpt('data-source-new-select-parent', parent)
    const queryContent = parserLoad(query, 'remove')
    $('#data-source-new-query').val(queryContent)
    const response = await fetch('/api/newChild?dataset=' + parent)
    const childName = await response.text()
    $('#data-source-new-name').val(childName)
    scrollAndBlink('data-source-container')
}

// see https://stackoverflow.com/questions/610406/javascript-equivalent-to-printf-string-format/32202320#32202320
String.prototype.format = function() {
    return [...arguments].reduce((p, c) => p.replace(/%s/, c), this)
}

function escapeQuery(q) {
    return q.replaceAll('\"', '\\\"')
}

function selectQuery(tmpl, queryContentId, selectListId, initValue) {
    const selectList = document.getElementById(selectListId)
    const selected = selectList.options[selectList.selectedIndex].text
    document.getElementById(queryContentId).innerText =
        selected.startsWith('--choose ') ? initValue : tmpl.format(escapeQuery(selected))
}

// function selectSource(id) {
//     const newDataSource = document.getElementById(id)
//     const inputs = newDataSource.getElementsByTagName('input')
//     if (isSelectedNone()) {
//         for (const input of inputs) {
//             input.disabled = false
//         }
//     } else {
//         for (const input of inputs) {
//             input.disabled = true
//         }
//     }
// }

function isSelectedNone() {
    const selectList = document.getElementById('select-data-source')
    const selected = selectList.options[selectList.selectedIndex].text
    return selected === '--None--'
}

async function handleHistory() {
    const url = '/api/history'
    fetch(url)
        .then(response => response.text())
        .then((response) => {
            document.getElementById('misc-history-list').innerHTML = response
        })
}

async function updateDatasetDesc() {
    await handleHistory()
    const dataset = document.getElementById('dataset-details-name').innerText
    handleDatasetDesc(dataset)
}

async function handleDatasetDesc(dataset) {
    const url = '/api/datasetDesc?dataset=' + dataset
    fetch(url)
        .then(response => response.text())
        .then((response) => {
            document.getElementById('misc-history-dataset-desc').innerHTML = response
        })
}

async function updateDatasetsOpt(id) {
    const selected = getSelected(id)
    const response = await fetch('/api/filteredSources')
    const datasets = await response.json()
    const selectList = document.getElementById(id)
    for (i = selectList.options.length - 1; i >= 1; i--) {
        selectList.remove(i)
    }
    for (const dataset of datasets) {
        const opt = document.createElement('option')
        opt.text = dataset
        selectList.add(opt)
    }
    changeSelectOpt(id, selected)
}

function changeSelectOpt(id, newOpt) {
    const selectList = document.getElementById(id)
    const opts = selectList.options
    for (var opt, j = 0; opt = opts[j]; j++) {
        if (opt.value == newOpt) {
            selectList.selectedIndex = j
            break
        }
    }
}

function parserLoad(query, cmd) {
    const regex = /LOAD (?<dataset>[a-zA-Z0-9_]+) /
    switch (cmd) {
        case 'dataset':
            const found = query.match(regex)
            return found.groups.dataset
        case 'remove':
            return query.replace(regex, '')
    }
}

function updateAll() {
    handleHistory()
    updateDatasetsOpt('data-source-exists-select-name')
    updateDatasetsOpt('data-source-new-select-parent')
}

function scrollAndBlink(id) {
    $('html').animate({
        scrollTop: $('#' + id).offset().top - 200
    }, 200, function() {
        $('#' + id).effect('highlight', {}, 500)
    })
}

function hoverButtons(id, button) {
    $('#list li').hover(function() {
        $(this).append('<button id="but1" >button1</button><button id="but2">button2</button>')
    }, function() {
        $(this).empty()
    })
}

function deleteAllBtns(elem) {
    Array.prototype.slice.call(elem.getElementsByTagName('button')).forEach(
        function(item) {
            item.remove()
        })
}

function addDatasetBtns(elem, dataset) {
    const removeBtn = document.createElement('button')
    removeBtn.textContent = 'details'
    removeBtn.classList.add('w3-ripple', 'w3-hover-green', 'inline-btn', 'no-default-btn')
    removeBtn.addEventListener('click', function() {
        handleDatasetDesc(dataset)
    })
    elem.append(removeBtn)

    const detailsBtn = document.createElement('button')
    detailsBtn.textContent = 'delete var'
    detailsBtn.classList.add('w3-ripple', 'w3-hover-green', 'inline-btn', 'no-default-btn')
    detailsBtn.addEventListener('click', function() {

    })
    elem.append(detailsBtn)
}

function addOptBtns(elem, idx) {
    const loadBtn = document.createElement('button')
    loadBtn.textContent = 'load'
    loadBtn.classList.add('w3-ripple', 'w3-hover-green', 'grid_opt-btn', 'no-default-btn', 'float-btn')
    loadBtn.addEventListener('click', function() {
        loadOption(idx)
    })
    elem.append(loadBtn)
}

// const backends = ['Value Type', 'Structure Difference', 'Distinct Values']

async function loadOption(idx) {
    handleButtonBar('data-source-exists-btn', 'data-source-div', 'data-source-exists')
    const dataset = document.getElementById('dataset-details-name').innerText
    const response = await fetch('/api/searchOpt?dataset=' + dataset + '&id=' + idx)
    const opt = await response.json()
    changeSelectOpt('data-source-exists-select-name', opt.Dataset)
    $('#dehMinRating').val(opt.Filter.MinRating)
    $('#dehMaxItems').val(opt.Filter.MaxNumber)
    opt.EnabledBackends.forEach(function(enabled, i) {
        $('#backend_' + i).prop('checked', enabled)
    })
    scrollAndBlink('data-source-container')
}
