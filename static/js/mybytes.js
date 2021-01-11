// function byteShow(id, brand, graph) {
function byteShow(){
    const id = "#mybox"
    const brand = "#brand"
    const graph = "#bytegraph"
    const maxVal = 255
    const iterations = 800
    const maxX = document.querySelector(graph).getAttribute("width")
    const maxY = document.querySelector(graph).getAttribute("height")

    let isMobile = (/(android|bb\d+|meego).+mobile|avantgo|bada\/|blackberry|blazer|compal|elaine|fennec|hiptop|iemobile|ip(hone|od)|ipad|iris|kindle|Android|Silk|lge |maemo|midp|mmp|netfront|opera m(ob|in)i|palm( os)?|phone|p(ixi|re)\/|plucker|pocket|psp|series(4|6)0|symbian|treo|up\.(browser|link)|vodafone|wap|windows (ce|phone)|xda|xiino/i.test(navigator.userAgent) || /1207|6310|6590|3gso|4thp|50[1-6]i|770s|802s|a wa|abac|ac(er|oo|s\-)|ai(ko|rn)|al(av|ca|co)|amoi|an(ex|ny|yw)|aptu|ar(ch|go)|as(te|us)|attw|au(di|\-m|r |s )|avan|be(ck|ll|nq)|bi(lb|rd)|bl(ac|az)|br(e|v)w|bumb|bw\-(n|u)|c55\/|capi|ccwa|cdm\-|cell|chtm|cldc|cmd\-|co(mp|nd)|craw|da(it|ll|ng)|dbte|dc\-s|devi|dica|dmob|do(c|p)o|ds(12|\-d)|el(49|ai)|em(l2|ul)|er(ic|k0)|esl8|ez([4-7]0|os|wa|ze)|fetc|fly(\-|_)|g1 u|g560|gene|gf\-5|g\-mo|go(\.w|od)|gr(ad|un)|haie|hcit|hd\-(m|p|t)|hei\-|hi(pt|ta)|hp( i|ip)|hs\-c|ht(c(\-| |_|a|g|p|s|t)|tp)|hu(aw|tc)|i\-(20|go|ma)|i230|iac( |\-|\/)|ibro|idea|ig01|ikom|im1k|inno|ipaq|iris|ja(t|v)a|jbro|jemu|jigs|kddi|keji|kgt( |\/)|klon|kpt |kwc\-|kyo(c|k)|le(no|xi)|lg( g|\/(k|l|u)|50|54|\-[a-w])|libw|lynx|m1\-w|m3ga|m50\/|ma(te|ui|xo)|mc(01|21|ca)|m\-cr|me(rc|ri)|mi(o8|oa|ts)|mmef|mo(01|02|bi|de|do|t(\-| |o|v)|zz)|mt(50|p1|v )|mwbp|mywa|n10[0-2]|n20[2-3]|n30(0|2)|n50(0|2|5)|n7(0(0|1)|10)|ne((c|m)\-|on|tf|wf|wg|wt)|nok(6|i)|nzph|o2im|op(ti|wv)|oran|owg1|p800|pan(a|d|t)|pdxg|pg(13|\-([1-8]|c))|phil|pire|pl(ay|uc)|pn\-2|po(ck|rt|se)|prox|psio|pt\-g|qa\-a|qc(07|12|21|32|60|\-[2-7]|i\-)|qtek|r380|r600|raks|rim9|ro(ve|zo)|s55\/|sa(ge|ma|mm|ms|ny|va)|sc(01|h\-|oo|p\-)|sdk\/|se(c(\-|0|1)|47|mc|nd|ri)|sgh\-|shar|sie(\-|m)|sk\-0|sl(45|id)|sm(al|ar|b3|it|t5)|so(ft|ny)|sp(01|h\-|v\-|v )|sy(01|mb)|t2(18|50)|t6(00|10|18)|ta(gt|lk)|tcl\-|tdg\-|tel(i|m)|tim\-|t\-mo|to(pl|sh)|ts(70|m\-|m3|m5)|tx\-9|up(\.b|g1|si)|utst|v400|v750|veri|vi(rg|te)|vk(40|5[0-3]|\-v)|vm40|voda|vulc|vx(52|53|60|61|70|80|81|83|85|98)|w3c(\-| )|webc|whit|wi(g |nc|nw)|wmlb|wonu|x700|yas\-|your|zeto|zte\-/i.test(navigator.userAgent.substr(0,4)))
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
    if(isMobile) {
        document.querySelector(brand).style.fontSize = "250%"
    }

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
