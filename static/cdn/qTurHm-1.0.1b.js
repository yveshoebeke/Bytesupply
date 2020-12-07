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
var u = new URL(document.getElementsByTagName('SCRIPT')[document.getElementsByTagName('SCRIPT').length-1].getAttribute("src"));
var k = u.searchParams.get("k");                    // User key (default: `sha-1 '21101956'`)
var c = u.searchParams.get("c");                    // Move target id (default: 'qTurHm')
var r = u.searchParams.get("r");                    // Result element id (default: 'qTutrHmPerception')
var f = u.searchParams.get("f");                    // Create hidden inputs if 'c' is part of a form (default: true)
var p = 5;                                          // Perception result
// set defaults if data parameters are not given
// move target id
if(c == null) { 
    c = "qTurHm";
}
// callback receiver element id default
if(r == null) {
    r = "qTurHmPerception";
}
// form input creation
if(f == null) {
    f = true;
} else {
    f = false;
}
// Add hidden input to form
Object.prototype.addHiddenInput = function(id, name, value) {
    this.type = "hidden";
    this.id = id;
    this.name = name;
    this.value = value;
}
// validate userkey
function validateUserkey(k){
    if(k == "a6bd3f10339b2d39aaa6175484a38173c1061f4a"){
        return true;
    } else {
        return false;
    }
}
// data object (to be JSON-ized if we need ajax -> rethink this?)
var data = {};
// Cursor movement data storage
var ms = new Array; 
// The move object
function mvd(t,x,y){
    this.t=t;   // unix time in ms
    this.x=x;   // x position
    this.y=y;   // y position
}
// Check if we are accessed from a mobile device
var isM = false; //initiate as false
// device detection
if(/(android|bb\d+|meego).+mobile|avantgo|bada\/|blackberry|blazer|compal|elaine|fennec|hiptop|iemobile|ip(hone|od)|ipad|iris|kindle|Android|Silk|lge |maemo|midp|mmp|netfront|opera m(ob|in)i|palm( os)?|phone|p(ixi|re)\/|plucker|pocket|psp|series(4|6)0|symbian|treo|up\.(browser|link)|vodafone|wap|windows (ce|phone)|xda|xiino/i.test(navigator.userAgent) 
|| /1207|6310|6590|3gso|4thp|50[1-6]i|770s|802s|a wa|abac|ac(er|oo|s\-)|ai(ko|rn)|al(av|ca|co)|amoi|an(ex|ny|yw)|aptu|ar(ch|go)|as(te|us)|attw|au(di|\-m|r |s )|avan|be(ck|ll|nq)|bi(lb|rd)|bl(ac|az)|br(e|v)w|bumb|bw\-(n|u)|c55\/|capi|ccwa|cdm\-|cell|chtm|cldc|cmd\-|co(mp|nd)|craw|da(it|ll|ng)|dbte|dc\-s|devi|dica|dmob|do(c|p)o|ds(12|\-d)|el(49|ai)|em(l2|ul)|er(ic|k0)|esl8|ez([4-7]0|os|wa|ze)|fetc|fly(\-|_)|g1 u|g560|gene|gf\-5|g\-mo|go(\.w|od)|gr(ad|un)|haie|hcit|hd\-(m|p|t)|hei\-|hi(pt|ta)|hp( i|ip)|hs\-c|ht(c(\-| |_|a|g|p|s|t)|tp)|hu(aw|tc)|i\-(20|go|ma)|i230|iac( |\-|\/)|ibro|idea|ig01|ikom|im1k|inno|ipaq|iris|ja(t|v)a|jbro|jemu|jigs|kddi|keji|kgt( |\/)|klon|kpt |kwc\-|kyo(c|k)|le(no|xi)|lg( g|\/(k|l|u)|50|54|\-[a-w])|libw|lynx|m1\-w|m3ga|m50\/|ma(te|ui|xo)|mc(01|21|ca)|m\-cr|me(rc|ri)|mi(o8|oa|ts)|mmef|mo(01|02|bi|de|do|t(\-| |o|v)|zz)|mt(50|p1|v )|mwbp|mywa|n10[0-2]|n20[2-3]|n30(0|2)|n50(0|2|5)|n7(0(0|1)|10)|ne((c|m)\-|on|tf|wf|wg|wt)|nok(6|i)|nzph|o2im|op(ti|wv)|oran|owg1|p800|pan(a|d|t)|pdxg|pg(13|\-([1-8]|c))|phil|pire|pl(ay|uc)|pn\-2|po(ck|rt|se)|prox|psio|pt\-g|qa\-a|qc(07|12|21|32|60|\-[2-7]|i\-)|qtek|r380|r600|raks|rim9|ro(ve|zo)|s55\/|sa(ge|ma|mm|ms|ny|va)|sc(01|h\-|oo|p\-)|sdk\/|se(c(\-|0|1)|47|mc|nd|ri)|sgh\-|shar|sie(\-|m)|sk\-0|sl(45|id)|sm(al|ar|b3|it|t5)|so(ft|ny)|sp(01|h\-|v\-|v )|sy(01|mb)|t2(18|50)|t6(00|10|18)|ta(gt|lk)|tcl\-|tdg\-|tel(i|m)|tim\-|t\-mo|to(pl|sh)|ts(70|m\-|m3|m5)|tx\-9|up(\.b|g1|si)|utst|v400|v750|veri|vi(rg|te)|vk(40|5[0-3]|\-v)|vm40|voda|vulc|vx(52|53|60|61|70|80|81|83|85|98)|w3c(\-| )|webc|whit|wi(g |nc|nw)|wmlb|wonu|x700|yas\-|your|zeto|zte\-/i.test(navigator.userAgent.substr(0,4))) { 
    isM = true;
}
// get target boundaries
function getCoords(el) {
    elem = document.getElementById(el);
    return { 
        top: elem.offsetTop, 
        left: elem.offsetLeft, 
        bottom: elem.offsetTop + elem.offsetHeight,
        right: elem.offsetLeft + elem.offsetWidth,
    }
}
// Analyse x/y vs. boundary behavior
function moveAnalysis(d){
    var o = 0;
    var p = 0;
    // Is move data array empty?
    if(d.samples < 2) {
        p = 0;
    } else if(d.moves[0].t > d.moves[d.samples-1].t) {
        p = 1;
        // Are last move coords differnet than click coords?
    } else if(d.moves[d.samples-1].x != d.moves[d.samples-2].x && 
        data.moves[d.samples-1].y != d.moves[d.samples-2].y) {
        p = 2;
        // Is last move timestamp less then 2hrs?
    } else if(d.moves[d.samples-1].t + 5200 < Date.now()) {
        p = 6;
    } else {
        // Check if mouse moves where within target obj limits.
        var skipNext = false;
        d.moves.forEach(function(v, i){
            if(i > 0) {
                for(j=0; j<=i; j++){
                    if(v.t < d.moves[j].t){
                        p = 1;
                        skipNext = true;
                    }
                }
            }
            if(v.x < d.target.left || v.x > d.target.right) {
                o++;
            }
            if(v.y < d.target.top || v.y > d.target.bottom) {
                o++;
            }
        });
        // evaluate p in relation to scale (0->10)
        if(!skipNext){
            p = 10 - Math.ceil((o * 10)/d.samples);
            // check if first move has a x == left || x == right or y == top || y == bottom
            if(p > 5){
                var firstX = d.moves[0].x;
                var firstY = d.moves[0].y;
                var l = 0
                if((firstX < d.target.left || firstX > d.target.right) ||
                    firstY < d.target.top || firstY > d.target.bottom) {
                    p--;
                }
            }
        }
    }
    // if mobile device and more than 2 moves degrade perception.
    if(d.isM && d.samples > 2) {
        p = Math.ceil(p/1.4);
    }

    return p;
} 
// Check if target element is part of a form
function isMemberOfForm(c) {
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
function qTurHmMousemove(e) {
    m = new mvd(Date.now(),e.pageX,e.pageY);    // harvest time, x, y for each move
    ms.push(m);                                 // ... and push it on the array
}
// Reset move data array when leaving obj and not clicked.
function qTurHmMouseout() {
    //console.log(e);
    ms.length = 0;
}
// Process move data when click occured.
function qTurHmClick(e) {
    m = new mvd(Date.now(),e.pageX,e.pageY);    // time, x, y for click event
    ms.push(m);                                 // ... and push it on the array
    this.removeEventListener('mousemove', qTurHmMousemove);
    this.removeEventListener('mouseout', qTurHmMouseout);
    // Cursor move target element dimensions
    var t = getCoords(c); 
    var n = Date.now();
    var rc = k + "_" + n.toString();
    // Create JSON Object
    data.userkey = k;                   // user supplied key
    data.timestamp = n;                 // this object's creation date
    data.resultcontent = rc;            // id tag for server to attach to result
    data.origURL = window.location.href;// request coming from this URL
    data.mobile = isM;                  // is this a mobile device?
    data.subject = c;                   // elem class where moves were derived from
    data.target = t;                    // move target object limits (see above)
    data.receiver = r;                  // where to push result to for callback
    data.samples = ms.length;           // number of movements captured
    data.moves = ms;                    // movement data array
    // Calculate perception result.
    p = moveAnalysis(data);
    data.perception = p;                // add perception result to the data object
    console.log("data",data);
    // Place perception result in appropriate element.
    // Check nature of target object clicked.
    // Add data store to target element with perception value.  and number of mousemove sample
    this.dataset.qturhm_perception = p;
    this.dataset.qturhm_samples = data.samples;
    // Check target tagname and see if it's part of a form, if so create the hidden inputs.
    var rjs = document.getElementById(r);
    if(rjs == null){
        if((f) && (this.tagName == "BUTTON" || this.tagName == "INPUT" || this.tagName == "A") && isMemberOfForm(this)) {
            var perceptionFormValue = document.createElement("input");
            var perceptionSamplesFormValue = document.createElement("input");
            perceptionFormValue.addHiddenInput(r, r, p);
            r = r + "_samples";
            perceptionSamplesFormValue.addHiddenInput(r, r, data.samples);
            this.parentElement.appendChild(perceptionFormValue);
            this.parentElement.appendChild(perceptionSamplesFormValue);
        }
    }
}

// It all starts when everything loaded.
document.addEventListener("DOMContentLoaded", function(){
    if(validateUserkey(k)){
        // Attach the event handlers ( --> check for previously added click event?)
        var cjs = document.getElementById(c);
        // Capture mouse movement over target obj.
        cjs.addEventListener('mousemove', qTurHmMousemove);
        // Reset move data array when leaving obj and not clicked.
        cjs.addEventListener('mouseout', qTurHmMouseout);
        // Process move data when click occured.
        cjs.addEventListener('click', qTurHmClick);
    } else {
        console.error("qTurHm: userkey not valid - functionality not supported.");
    }
});
