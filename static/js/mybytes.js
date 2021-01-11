// function byteShow(id, brand, graph) {
function byteShow(){
    const id = "#mybox"
    const brand = "#brand"
    const graph = "#bytegraph"
    const maxVal = 255
    const iterations = 800
    const maxX = document.querySelector(graph).getAttribute("width")
    const maxY = document.querySelector(graph).getAttribute("height")

    let interval = offset = 10
    let result = ""
    let byterun = 0
    
    let canvas = document.querySelector(graph)
    let ctx = canvas.getContext("2d")

    var gX = gY = sum = 0
    var coords = []
    var xCoef = maxX / iterations
    var yCoef = maxY / maxVal
    
    //           ^ ^
    //          (o O)
    // _______oOO(.)OOo________
    // ________________________
    // ________________________
    // 3) Dim and remove the canvasses prior to deleting them. Remove event listener.   
    // document.querySelector(id).fadeBrand(60, (interval * iterations) + 1500)
    // document.querySelector(brand).fadeBrand(40, (interval * iterations) + 1400)
    // document.querySelector(graph).fadeBrand(40, (interval * iterations) + 2000)
    // See notes below at EOF why Object.prototype.fadeBrand is not compatible.
    fadeBrand(id, 60, (interval * iterations) + 1500)
    fadeBrand(brand, 40, (interval * iterations) + 1400)
    fadeBrand(graph, 40, (interval * iterations) + 2000)
    setTimeout(document.removeEventListener("DOMContentLoaded", byteShow, true), (interval * iterations) + 2100)
    
    // 2) Accent brandname in blue
    setTimeout(function(){ document.querySelector(brand).style.color = "#0000FF" }, (interval * iterations) - 600)   

    // 1) Heart of the beast:
    //  a. Generate an int 0->255
    //  b. Display binary value
    //  c. Graph the value (yellow-ish)
    //  d. Graph the mean average (red)
    do {
        setTimeout(function(){
            // Generate the number
            // Original was: let mybyte = (Math.floor((Math.random() * maxVal) + 1))
            let mybyte = (Math.floor((Math.random() * Math.abs(maxVal-(gX / 1.4)))) + 1)
            // Make it a string with padding as neccessary
            let byteOut = convertToBinaryString(mybyte)
            // Append it to existing string value
            result += byteOut
            // Output result to id canvas div
            document.querySelector(id).innerHTML = '<span class="mybits">' + result + '</span>'
            
            // Graph values and averages
            // First the value itself
            gY = mybyte
            ctx.translate(0, 0)
            ctx.beginPath()
            ctx.moveTo(Math.round(gX * xCoef), maxY)
            ctx.lineTo(Math.round(gX * xCoef), maxY - Math.round(gY * yCoef))
            ctx.strokeStyle = "#BFBF00"
            ctx.stroke()
            ctx.closePath()

            // Calculate mean average and draw it
            sum = 0
            coords.forEach(coord => sum += coord.y)

            if(coords.length > 0) {
                pAvg = sum / coords.length
                pX = gX - 1
            } else {
                pAvg = mybyte
                pX = 0
            }

            coords.push({x: gX, y: gY})
            sum = 0
            coords.forEach(coord => sum += coord.y)

            ctx.beginPath()
            ctx.moveTo(pX * xCoef, maxY - Math.round(pAvg * yCoef))
            ctx.lineTo(gX * xCoef, maxY - Math.round((sum / coords.length) * yCoef))
            ctx.strokeStyle = "#FF0000"
            ctx.stroke()
            ctx.closePath()

            gX++

        }, interval)
        interval += offset
    } while(byterun++ <= iterations)
}

// Start when DOM is loaded
document.addEventListener("DOMContentLoaded", byteShow, true)

// fade Obj given fade speed (ms) and execution delay (ms)
// ** Note: Object prototyping not compatible with bootstrap  :( **
// Object.prototype.fadeBrand = function(speed, delay) {
//     let dimmer = 1.0        
//     for(i = speed; i <= (speed * 10) + speed; i += speed){
//         var elem = this
//         setTimeout(function() {
//             elem.style.opacity = dimmer
//             dimmer -= 0.1
//         }, (delay + i))
//         setTimeout(function() {
//             elem.style.display = "none"
//             elem.remove()
//         }, delay + (speed * 10))
//     }
// }

function fadeBrand(id, speed, delay) {
    let dimmer = 1.0        
    var elem = document.querySelector(id)
    for(i = speed; i <= (speed * 10) + speed; i += speed){
        setTimeout(function() {
            elem.style.opacity = dimmer
            dimmer -= 0.1
        }, (delay + i))
        setTimeout(function() {
            elem.style.display = "none"
            elem.remove()
        }, delay + (speed * 10))
    }
}

// Helper func - convert int to string binary representation with leading padding as needed
function convertToBinaryString(num) {
    let padding = "00000000"
    let strNum = num.toString(2)
    let bytelength = strNum.length
    if(bytelength < 8) {
        return padding.substr(0, 8 - bytelength) + strNum + " "
    }
    return strNum.substr(0, 8) + " "
}
