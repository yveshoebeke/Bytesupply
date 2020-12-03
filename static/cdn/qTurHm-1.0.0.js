/* 
    qTurHm (_q_uasi _T_uring _H_uman vs _m_achine evaluator)
    --------------------------------------------------------
    Collects cursor movement data over a clickable HTML element
    and send the result to the server for evaluation.

    An int result code between 0 and 10 is returned indicating 
    perceived probability of above element was engaged by a humon or robot.

    Result value scale:

    0-1-2-3-4-5-6-7-8-9-10
    ^                   ^
    Machine             Human

    (c) 2020 - Bytesupply, LLC
*/
var url = "https://bytesupply.com/api/v1/qTurHm"; // Server url
var u = new URL($('script').last().attr('src'));    // Get this scripts url
var c = u.searchParams.get("c");                    // Move target class (default: 'qTurHm')
var k = u.searchParams.get("k");                    // User key (default: `sha-1 '21101956'`)
var r = u.searchParams.get("r");                    // Result element id (default: 'qTutrHm_Result')
// set defauts if data parameters are not given
// target move element class
if(c == null) { 
    c = "qTurHm";
}
// callback receiver element id default
if(r == null) { 
    r = "qTurHm_Result";
}
// key default
if (k == null){
    k = "1dc9b274eb754dfa1574984a56561b88214a1802";
}
// Cursor movement data storage
ms = new Array; 
// Move object
function mvd(t,x,y){
    this.t=t;   // unix time in ms
    this.x=x;   // x position
    this.y=y;   // y position
}

$(function() {
    c = "." + c;    // c = class of submit element
    $(c).mousemove(function(e){
        m = new mvd(Date.now(),e.pageX,e.pageY);    // harvest time, x, y for each move
        ms.push(m);                                 // ... and push it on the array
    }).click(function(){
        $(c).unbind("mousemove").unbind("click");   // disable click when clicked

        // Create JSON Object
        var t = {};                         // Cursor move target element dims
        t.top = ~~$(c).position().top;      // upper limit (min val on y-axis)
        t.left = ~~$(c).position().left;    // left limit (min val on x-axis)
        t.width = ~~$(c).width();           // width (+ left limit = max val on x-axis)
        t.height = ~~$(c).height();         // height (+ top = max val on y-axis)

        n = Date.now();
        rc = k + "_" + n.toString();

        var data = {};                      // data objet to be JSON-ized
        data.userkey = k;                   // user supplied key
        data.timestamp = n;                 // this object's creation date
        data.resultcontent = rc;            // id tag for server to attach to result
        data.origURL = window.location.href;// request coming from this URL
        data.subject = c;                   // elem class where moves were derived from
        data.target = t;                    // move target object (see above)
        data.receiver = r;                  // where to push result to for callback
        data.samples = ms.length;           // number of movements captured
        data.moves = ms;                    // movement data array
   
        // Create the JSON Object
        jsonData = JSON.stringify(data);

        // Send it to the server
        $.post(url, jsonData, function(jsonData, status){
            console.log("status is " + status);
        });

    });
});