var j = 0
var words = [
  'ding1','ding2','ding3',
  '顶1','顶2','顶3'
]

function run() {
  setTimeout(function () {
    let word = words[j % word.length]
    j++
    $("#ueditor_replace>p").text(word)
    $(".poster_submit").click()
    run()
  }, parseInt(Math.random() * 1000) * 1000)
}

run()
