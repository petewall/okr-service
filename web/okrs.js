function insertOKRRow() {
  location.reload()
}

function updateOKRRow() {
  location.reload()
}

function getCurrentFiscalYear() {
  const year = new Date().getFullYear()
  const month = new Date().getMonth()

  // January is the only date where the fiscal year matches the calendar year
  if (month == 0) {
    return year
  }
  return year + 1
}

function getCurrentQuarter() {
  // Q1: February, March, April
  // Q2: May, June, July
  // Q3: August, September, October
  // Q4: November, December, January

  const month = new Date().getMonth()               // 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11
  const quarter = Math.floor((month-1)/3) % 4 + 1   // 4, 1, 1, 1, 2, 2, 2, 3, 3, 3,  4,  4
  return `Q${quarter}`
}

function showNewModal() {
  $('.ui.edit.modal>.header').text('New OKR')
  $('.ui.edit.modal select.quarter').val(getCurrentQuarter())
  $('.ui.edit.modal input.year').val(getCurrentFiscalYear())
  $('.ui.edit.modal input.category').val('')
  $('.ui.edit.modal input.description').val('')
  $('.ui.edit.modal select.type').val('boolean')
  $('.ui.edit.modal input.progress').val(0)
  $('.ui.edit.modal input.goal').val(1)
  $('.ui.edit.modal').modal({
    onApprove: function() {
      const newOKR = {
        quarter: `${$('.ui.edit.modal input.year').val()}${$('.ui.edit.modal select.quarter').val()}`,
        category: $('.ui.edit.modal input.category').val(),
        description: $('.ui.edit.modal input.description').val(),
        type: $('.ui.edit.modal select.type').val(),
        progress: parseFloat($('.ui.edit.modal input.progress').val()),
        goal: parseFloat($('.ui.edit.modal input.goal').val())
      }
      $.ajax({
        url: 'http://localhost:8080/api/okr',
        type: 'PUT',
        data: JSON.stringify(newOKR),
        contentType: 'application/json',
        success: function() {
          insertOKRRow(newOKR)
        }
      })
    }
  })
  $('.ui.edit.modal').modal('show')
}

function showEditModal(okr) {
  $('.ui.edit.modal>.header').text('Edit OKR')
  $('.ui.edit.modal select.quarter').val(okr.quarter.substring(4))
  $('.ui.edit.modal input.year').val(okr.quarter.substring(0,4))
  $('.ui.edit.modal input.category').val(okr.category)
  $('.ui.edit.modal input.description').val(okr.description)
  $('.ui.edit.modal select.type').val(okr.type)
  $('.ui.edit.modal input.progress').val(okr.progress)
  $('.ui.edit.modal input.goal').val(okr.goal)
  $('.ui.edit.modal').modal({
    onApprove: function() {
      const newOKR = {
        id: okr.id,
        quarter: `${$('.ui.edit.modal input.year').val()}${$('.ui.edit.modal select.quarter').val()}`,
        category: $('.ui.edit.modal input.category').val(),
        description: $('.ui.edit.modal input.description').val(),
        type: $('.ui.edit.modal select.type').val(),
        progress: parseFloat($('.ui.edit.modal input.progress').val()),
        goal: parseFloat($('.ui.edit.modal input.goal').val())
      }
      $.ajax({
        url: 'http://localhost:8080/api/okr',
        type: 'POST',
        data: JSON.stringify(newOKR),
        contentType: 'application/json',
        success: function() {
          updateOKRRow(newOKR)
        }
      })
    }
  })
  $('.ui.edit.modal').modal('show')
}

function showDeleteModal(okr) {
  $('.ui.delete.modal').modal({
    onApprove: function() {
      $.ajax({
        url: `http://localhost:8080/api/okr/${okr.id}`,
        type: 'DELETE',
        success: function() {
          $(`#${okr.id}`).remove()
        }
      })
    }
  })
  $('.ui.delete.modal').modal('show')
}

function makeProgress(type, progress, goal) {
  const percent = Math.round((progress / goal) * 100)

  if (type == 'boolean') {
    if (progress == goal) {
      return 'Done!'
    } else {
      return 'Not yet'
    }
  } else if (type == 'number') {
    return `${progress} / ${goal} (${percent}%)`
  } else if (type == 'percentage') {
    return `${percent}%`
  }
}

function setOKRs(okrs) {
  let structured = {}
  for (const okr of okrs) {
    if (!structured[okr.quarter]) {
      structured[okr.quarter] = {}
    }
    if (!structured[okr.quarter][okr.category]) {
      structured[okr.quarter][okr.category] = []
    }
    structured[okr.quarter][okr.category].push(okr)
  }

  $('#okr-list').empty()
  const quarters = Object.keys(structured).sort().reverse()
  for (const quarter of quarters) {
    const quarterSection = $('<div>', {id: quarter}).append(
      $('<a>', {href: `#${quarter}`})
    )
    $(quarterSection).append($('<h2>', {text: quarter}))
    const okrTable = $('<table>', {class: 'ui celled table'})

    const categories = Object.keys(structured[quarter]).sort()
    for (const category of categories) {
      okrTable.append($('<thead>').append($('<tr>').append(
        $('<th>', {colspan: 3, text: category}))
      ))

      const okrs = structured[quarter][category]
      for (const okr of okrs) {
        const myOKR = okr
        const editButton = $('<i>', {class: 'pencil icon'}).click(() => {showEditModal(myOKR)})
        const deleteButton = $('<i>', {class: 'trash icon'}).click(() => {showDeleteModal(myOKR)})

        okrTable.append($('<tr>', {
          id: okr.id
        }).append(
          $('<td>', {class: 'description', text: okr.description}),
          $('<td>', {class: 'collapsing progress'}).append(makeProgress(okr.type, okr.progress, okr.goal)),
          $('<td>', {class: 'collapsing tools'}).append(editButton, deleteButton)
        ))
      }

      $(quarterSection).append(okrTable)
    }
   
    $('#okr-list').append(quarterSection)
  }
}

$(document).ready(() => {
  $.get('http://localhost:8080/api/okrs', (data) => {
    const okrs = JSON.parse(data)
    setOKRs(okrs)

    $('.new.button').click(showNewModal)
  })
})
