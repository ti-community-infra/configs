function rgbToHex(col) {
    if (col.charAt(0) == 'r') {
        col = col.replace('rgb(', '').replace(')', '').split(',');
        var r = parseInt(col[0], 10).toString(16);
        var g = parseInt(col[1], 10).toString(16);
        var b = parseInt(col[2], 10).toString(16);
        r = r.length == 1 ? '0' + r : r; g = g.length == 1 ? '0' + g : g; b = b.length == 1 ? '0' + b : b;
        var colHex = '#' + r + g + b;
        return colHex;
    }
}

var results = []
document.querySelectorAll("#label-select-menu > details-menu > div > div.SelectMenu-list.select-menu-list > div:nth-child(1) > a").forEach(n => {
    var labelNameNode = n.querySelector("div > div:nth-child(1)");
    var labelDesNode = n.querySelector("div > div:nth-child(2)")
    var labelColorNode = n.querySelector("span");
    if (labelNameNode) {
        var name = labelNameNode.textContent;
        var color = rgbToHex(labelColorNode.style.backgroundColor).replace('#', '');
        var e = { name, color }
        if (labelDesNode) {
            e.description = labelDesNode.textContent.replace(/[\n\t\r]/g, "").trim();
        }
        results.push(e)
    }
})
console.log(JSON.stringify(results, null, 2))