/* 
    qTurHm (_q_uasi _Tu_ring _H_uman vs _m_achine evaluator)
    --------------------------------------------------------
    Collects cursor movement data over a clickable HTML element
    and send the perception result to the server for evaluation.

    An int result code between 0 and 10 is returned indicating 
    perceived probability of above element was engaged by a humon or robot.

    Perception result value scale:

    0-1-2-3-4-5-6-7-8-9-10
    ^                   ^
    Machine             Human

    (c) 2020 - Bytesupply, LLC
*/
const url = new URL(document.querySelectorAll('script')[document.querySelectorAll('script').length-1].getAttribute("src"));
// const url = new URL(document.getElementsByTagName('SCRIPT')[document.getElementsByTagName('SCRIPT').length-1].getAttribute("src"));
let validateUserkey = (key) => key == "a6bd3f10339b2d39aaa6175484a38173c1061f4a"    // validate userkey
// BS_qTurHm object creation
const BS_qTurHm = {
    userKey: url.searchParams.get("k"),
    createdTimestamp: Date.now(),
    scriptVersion: "1.0.1c",
    validUserkey: validateUserkey(url.searchParams.get("k")),
    origineHttp: window.location.href,
    mobilePlatform: (/(android|bb\d+|meego).+mobile|avantgo|bada\/|blackberry|blazer|compal|elaine|fennec|hiptop|iemobile|ip(hone|od)|ipad|iris|kindle|Android|Silk|lge |maemo|midp|mmp|netfront|opera m(ob|in)i|palm( os)?|phone|p(ixi|re)\/|plucker|pocket|psp|series(4|6)0|symbian|treo|up\.(browser|link)|vodafone|wap|windows (ce|phone)|xda|xiino/i.test(navigator.userAgent) || /1207|6310|6590|3gso|4thp|50[1-6]i|770s|802s|a wa|abac|ac(er|oo|s\-)|ai(ko|rn)|al(av|ca|co)|amoi|an(ex|ny|yw)|aptu|ar(ch|go)|as(te|us)|attw|au(di|\-m|r |s )|avan|be(ck|ll|nq)|bi(lb|rd)|bl(ac|az)|br(e|v)w|bumb|bw\-(n|u)|c55\/|capi|ccwa|cdm\-|cell|chtm|cldc|cmd\-|co(mp|nd)|craw|da(it|ll|ng)|dbte|dc\-s|devi|dica|dmob|do(c|p)o|ds(12|\-d)|el(49|ai)|em(l2|ul)|er(ic|k0)|esl8|ez([4-7]0|os|wa|ze)|fetc|fly(\-|_)|g1 u|g560|gene|gf\-5|g\-mo|go(\.w|od)|gr(ad|un)|haie|hcit|hd\-(m|p|t)|hei\-|hi(pt|ta)|hp( i|ip)|hs\-c|ht(c(\-| |_|a|g|p|s|t)|tp)|hu(aw|tc)|i\-(20|go|ma)|i230|iac( |\-|\/)|ibro|idea|ig01|ikom|im1k|inno|ipaq|iris|ja(t|v)a|jbro|jemu|jigs|kddi|keji|kgt( |\/)|klon|kpt |kwc\-|kyo(c|k)|le(no|xi)|lg( g|\/(k|l|u)|50|54|\-[a-w])|libw|lynx|m1\-w|m3ga|m50\/|ma(te|ui|xo)|mc(01|21|ca)|m\-cr|me(rc|ri)|mi(o8|oa|ts)|mmef|mo(01|02|bi|de|do|t(\-| |o|v)|zz)|mt(50|p1|v )|mwbp|mywa|n10[0-2]|n20[2-3]|n30(0|2)|n50(0|2|5)|n7(0(0|1)|10)|ne((c|m)\-|on|tf|wf|wg|wt)|nok(6|i)|nzph|o2im|op(ti|wv)|oran|owg1|p800|pan(a|d|t)|pdxg|pg(13|\-([1-8]|c))|phil|pire|pl(ay|uc)|pn\-2|po(ck|rt|se)|prox|psio|pt\-g|qa\-a|qc(07|12|21|32|60|\-[2-7]|i\-)|qtek|r380|r600|raks|rim9|ro(ve|zo)|s55\/|sa(ge|ma|mm|ms|ny|va)|sc(01|h\-|oo|p\-)|sdk\/|se(c(\-|0|1)|47|mc|nd|ri)|sgh\-|shar|sie(\-|m)|sk\-0|sl(45|id)|sm(al|ar|b3|it|t5)|so(ft|ny)|sp(01|h\-|v\-|v )|sy(01|mb)|t2(18|50)|t6(00|10|18)|ta(gt|lk)|tcl\-|tdg\-|tel(i|m)|tim\-|t\-mo|to(pl|sh)|ts(70|m\-|m3|m5)|tx\-9|up(\.b|g1|si)|utst|v400|v750|veri|vi(rg|te)|vk(40|5[0-3]|\-v)|vm40|voda|vulc|vx(52|53|60|61|70|80|81|83|85|98)|w3c(\-| )|webc|whit|wi(g |nc|nw)|wmlb|wonu|x700|yas\-|your|zeto|zte\-/i.test(navigator.userAgent.substr(0,4))),
    targetElement: (url.searchParams.get("c") == null) ? "qTurHm" : url.searchParams.get("c"),
    perceptionElement: (url.searchParams.get("r") == null) ? "qTurHmPerception" : url.searchParams.get("r"),
    createFormElement: (url.searchParams.get("f") == null) ? true : false,
    exposeInformation: (url.searchParams.get("x") == null) ? false : true,
    perception: 0,
    totalMoves: 0,
    movesDetail: []
}
console.info("BS_qTurHm(1)", BS_qTurHm)
let makeElement = (name, type) => {
    switch(type){
        case 'id':
            return "#"+name
        case 'class':
            return "."+name
        default:
            return name
    }
}
// Add hidden input to form
Object.prototype.addHiddenInput = function(id, name, value) {
    this.type = "hidden";
    this.id = id;
    this.name = name;
    this.value = value;
}
// The move object --> simplify assignment
// function moveDetail(t,x,y,cx,cy,sx,sy,b){
//     this.t=t;   // unix time in ms
//     this.x=x;   // page x position
//     this.y=y;   // page y position
//     this.cx=cx;   // client x position
//     this.cy=cy;   // client y position
//     this.sx=sx;   // screen x position
//     this.sy=sy;   // screen y position
//     this.b=b;    // target boundaries
// }
// const moveDetail = (t,x,y,cx,cy,sx,sy,b) => {t: t, x: x, y: y, cx: cx, cy: cy, sx: sx, sy: sy ,b: b}
// get target boundaries n--> simplify assignment
function BS_qTurHm_getCoords(el) {
    elem = document.querySelector(`#${el}`);
    return { 
        top: elem.offsetTop, 
        left: elem.offsetLeft, 
        bottom: elem.offsetTop + elem.offsetHeight,
        right: elem.offsetLeft + elem.offsetWidth,
    }
}
// Analyse x/y vs. boundary behavior --> and more
function BS_qTurHm_moveAnalysis(){
    let o = 0;
    let p = 0;
    // if mobile device and more than 2 moves degrade perception.
    if(BS_qTurHm.mobilePlatform && BS_qTurHm.totalMoves > 2) {
        return 3;
        // Is move data array empty?
    } else if(BS_qTurHm.totalMoves < 2) {
        return 0;
    } else if(BS_qTurHm.movesDetail[0].b.t > BS_qTurHm.movesDetail[BS_qTurHm.totalMoves-1].b.t) {
        return 1;
        // Are last move coords differnet than click coords?
    } else if(BS_qTurHm.movesDetail[BS_qTurHm.totalMoves-1].x != BS_qTurHm.movesDetail[BS_qTurHm.totalMoves-2].x && 
        BS_qTurHm.movesDetail[BS_qTurHm.totalMoves-1].y != BS_qTurHm.movesDetail[BS_qTurHm.totalMoves-2].y) {
        return 2;
        // Is last move timestamp less then 2hrs?
    } else if(BS_qTurHm.movesDetail[BS_qTurHm.totalMoves-1].t + 5200 < Date.now()) {
        return 6;
    } else {
        // Check if mouse moves where within target obj limits.
        var skipNext = false;
        BS_qTurHm.movesDetail.forEach(function(v, i){
            if(i > 0) {
                for(j=0; j<=i; j++){
                    if(v.t < d.moves[j].t){
                        p = 1;
                        skipNext = true;
                    }
                }
            }
            if(v.x < v.b.left || v.x > v.b.right) {
                o++;
            }
            if(v.y < v.b.top || v.y > v.b.bottom) {
                o++;
            }
        });
        // evaluate p in relation to scale (0->10)
        if(!skipNext){
            p = 10 - Math.ceil((o * 10)/BS_qTurHm.totalMoves);
            // check if first move has a x == left || x == right or y == top || y == bottom
            if(p > 5){
                if((BS_qTurHm.movesDetail[0].x < BS_qTurHm.movesDetail[0].b.left || BS_qTurHm.movesDetail[0].x > BS_qTurHm.movesDetail[0].b.right) ||
                    BS_qTurHm.movesDetail[0].y < BS_qTurHm.movesDetail[0].b.top || BS_qTurHm.movesDetail[0].y > BS_qTurHm.movesDetail[0].b.bottom) {
                        p--;
                }
            }
        }
    }
    return 8;   //p;
} 
// Check if target element is part of a form
function BS_qTurHm_isMemberOfForm(c) {
    while(c.parentElement.nodeName != "BODY"){
        if( c.parentElement.nodeName == "FORM") {
            return true
        } else if(c.parentElement == null){
            break;
        } else {
            c = c.parentNode;
        }
    }
    return false;
}
// Event handlers
// Capture mouse movement over target obj.
function BSqTurHmMouseMove(e) {
    // m = new moveDetail(Date.now(),e.pageX,e.pageY,e.clientX,e.clientY,e.screenX,e.screenY,BS_qTurHm_getCoords(BS_qTurHm.targetElement));    // harvest time, x, y for each move
    // BS_qTurHm.movesDetail.push(m)
    BS_qTurHm.movesDetail.push({t: Date.now(), x: e.pageX, y: e.pageY, cx: e.clientX, cy: e.clientY, sx: e.screenX, sy: e.screenY,b: BS_qTurHm_getCoords(BS_qTurHm.targetElement)})
    BS_qTurHm.totalMoves++
    console.info(`moves ${BS_qTurHm.totalMoves}`)
    // console.log(BS_qTurHm.totalMoves,m)
}
// Reset move data array when leaving obj and not clicked.
function BSqTurHmMouseOut() {
    BS_qTurHm.movesDetail.length = 0
    BS_qTurHm.totalMoves = 0
}
// Process move data when click occured.
function BSqTurHmMouseClick(e) {
    // m = new moveDetail(Date.now(),e.pageX,e.pageY,e.clientX,e.clientY,e.screenX,e.screenY,BS_qTurHm_getCoords(BS_qTurHm.targetElement));    // time, x, y for click event
    cjs.removeEventListener('mouseout', BSqTurHmMouseOut);
    cjs.removeEventListener('mousemove', BSqTurHmMouseMove);
    //BS_qTurHm.movesDetail.push(m)
    BS_qTurHm.movesDetail.push({t: Date.now(), x: e.pageX, y: e.pageY, cx: e.clientX, cy: e.clientY, sx: e.screenX, sy: e.screenY,b: BS_qTurHm_getCoords(BS_qTurHm.targetElement)})
    BS_qTurHm.totalMoves++
    console.log(`Clicked, moves: ${BS_qTurHm.totalMoves}`)
    // Calculate perception result.
    BS_qTurHm.perception = BS_qTurHm_moveAnalysis();                // add perception result to the data object
    console.info("BS_qTurHm(2)", BS_qTurHm)

    var rjs = document.querySelector(`#${BS_qTurHm.perceptionElement}`);
    if(rjs == null){
        if((BS_qTurHm.createFormElement) && (this.tagName == "BUTTON" || this.tagName == "INPUT" || this.tagName == "A") && BS_qTurHm_isMemberOfForm(this)) {
            var perceptionFormValue = document.createElement("input");
            var perceptionSamplesFormValue = document.createElement("input");
            perceptionFormValue.addHiddenInput(BS_qTurHm.perceptionElement, BS_qTurHm.perceptionElement, BS_qTurHm.perception);
            BS_qTurHm_r_s = BS_qTurHm.perceptionElement + "_samples";
            perceptionSamplesFormValue.addHiddenInput(BS_qTurHm_r_s, BS_qTurHmr_s, BS_qTurHm.totalMoves);
            this.parentElement.appendChild(perceptionFormValue);
            this.parentElement.appendChild(perceptionSamplesFormValue);
        } else if(rjs.tagName == "INPUT") {
            rjs.value = BS_qTurHm_data.perception;
        }
    }
}

// It all starts when everything loaded.
document.addEventListener("DOMContentLoaded", function(){
    if(BS_qTurHm.validUserkey){
        // Attach the event handlers ( --> check for previously added click event?)
        var cjs = document.querySelector(`#${BS_qTurHm.targetElement}`);
        // Capture mouse movement over target obj.
        cjs.addEventListener('mousemove', BSqTurHmMouseMove);
        // Reset move data array when leaving obj and not clicked.
        cjs.addEventListener('mouseout', BSqTurHmMouseOut);
        // Process move data when click occured.
        cjs.addEventListener('click', BSqTurHmMouseClick);
    } else {
        console.error("qTurHm: userkey "+BS_qTurHm.userKey+" not valid - functionality not supported.");
    }
});
