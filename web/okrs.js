function setOKRs(okrs) {
  let structured = {}
  for (okr of okrs) {
    if (!structured[okr.quarter]) {
      structured[okr.quarter] = {}
    }
    if (!structured[okr.quarter][okr.category]) {
      structured[okr.quarter][okr.category] = []
    }
    structured[okr.quarter][okr.category].push(okr)
  }

  $("#okr-list").empty()
  const quarters = Object.keys(structured).sort().reverse()
  for (quarter of quarters) {
    const quarterSection = $("<div>", {id: quarter}).append(
      $("<a>", {href: `#${quarter}`})
    )
    $(quarterSection).append($("<h2>", {text: quarter}))
    
    const categories = Object.keys(structured[quarter]).sort()
    for (category of categories) {
      const categorySection = $("<div>", {id: category})
      const okrTable = $("<table>", {class: "okrList"})

      const okrs = structured[quarter][category]
      for (okr of okrs) {
        okrTable.append($("<tr>").append(
          $("<td>", {class: "description", text: okr.description}),
          $("<td>", {class: "progress", text: `${okr.progress} / ${okr.goal}`}),
          $("<td>", {class: "tools", text: "X"})
        ))
      }

      categorySection.append(
        $("<h3>", {text: category}),
        okrTable
      )
      $(quarterSection).append(categorySection)
    }
   
    $("#okr-list").append(quarterSection)
  }
}

$(document).ready(() => {
  $.get("http://localhost:8080/api/okrs", (data) => {
    const okrs = JSON.parse(data)
    setOKRs(okrs)
  })
})

