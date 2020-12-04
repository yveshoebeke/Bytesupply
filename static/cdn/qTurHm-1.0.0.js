/* 
    qTurHm (_q_uasi _T_uring _H_uman vs _m_achine evaluator)
    --------------------------------------------------------
    Collects cursor movement data over a clickable HTML element
    and send the result to the server for evaluation.

    An int result code between 0 and 10 is returned indicating 
    perceived probability of above element was engaged by a humon or robot.

    Perception result value scale:

    0-1-2-3-4-5-6-7-8-9-10
    ^                   ^
    Machine             Human

    (c) 2020 - Bytesupply, LLC
*/
var url = "https://bytesupply.com/api/v1/qTurHm";   // Server url
var u = new URL($('script').last().attr('src'));    // Get this scripts url
var k = u.searchParams.get("k");                    // User key (default: `sha-1 '21101956'`)
var c = u.searchParams.get("c");                    // Move target class (default: 'qTurHm')
var r = u.searchParams.get("r");                    // Result element id (default: 'qTutrHmPerception')
var p = 5;                                          // Perception result
// set defauts if data parameters are not given
// target move element class
if(c == null) { 
    c = ".qTurHm";
} else {
    c = "." + c;    // c = class of submit element
}
// callback receiver element id default
if(r == null) {
    r = "#qTurHmPerception";
} else {
    r = "#" + r;    // r = id of result receiving element
}
// key default
if (k == null){
    k = "1dc9b274eb754dfa1574984a56561b88214a1802";
}
// Check if we are accessed from a mobile device
var isM = false; //initiate as false
// device detection
if(/(android|bb\d+|meego).+mobile|avantgo|bada\/|blackberry|blazer|compal|elaine|fennec|hiptop|iemobile|ip(hone|od)|ipad|iris|kindle|Android|Silk|lge |maemo|midp|mmp|netfront|opera m(ob|in)i|palm( os)?|phone|p(ixi|re)\/|plucker|pocket|psp|series(4|6)0|symbian|treo|up\.(browser|link)|vodafone|wap|windows (ce|phone)|xda|xiino/i.test(navigator.userAgent) 
    || /1207|6310|6590|3gso|4thp|50[1-6]i|770s|802s|a wa|abac|ac(er|oo|s\-)|ai(ko|rn)|al(av|ca|co)|amoi|an(ex|ny|yw)|aptu|ar(ch|go)|as(te|us)|attw|au(di|\-m|r |s )|avan|be(ck|ll|nq)|bi(lb|rd)|bl(ac|az)|br(e|v)w|bumb|bw\-(n|u)|c55\/|capi|ccwa|cdm\-|cell|chtm|cldc|cmd\-|co(mp|nd)|craw|da(it|ll|ng)|dbte|dc\-s|devi|dica|dmob|do(c|p)o|ds(12|\-d)|el(49|ai)|em(l2|ul)|er(ic|k0)|esl8|ez([4-7]0|os|wa|ze)|fetc|fly(\-|_)|g1 u|g560|gene|gf\-5|g\-mo|go(\.w|od)|gr(ad|un)|haie|hcit|hd\-(m|p|t)|hei\-|hi(pt|ta)|hp( i|ip)|hs\-c|ht(c(\-| |_|a|g|p|s|t)|tp)|hu(aw|tc)|i\-(20|go|ma)|i230|iac( |\-|\/)|ibro|idea|ig01|ikom|im1k|inno|ipaq|iris|ja(t|v)a|jbro|jemu|jigs|kddi|keji|kgt( |\/)|klon|kpt |kwc\-|kyo(c|k)|le(no|xi)|lg( g|\/(k|l|u)|50|54|\-[a-w])|libw|lynx|m1\-w|m3ga|m50\/|ma(te|ui|xo)|mc(01|21|ca)|m\-cr|me(rc|ri)|mi(o8|oa|ts)|mmef|mo(01|02|bi|de|do|t(\-| |o|v)|zz)|mt(50|p1|v )|mwbp|mywa|n10[0-2]|n20[2-3]|n30(0|2)|n50(0|2|5)|n7(0(0|1)|10)|ne((c|m)\-|on|tf|wf|wg|wt)|nok(6|i)|nzph|o2im|op(ti|wv)|oran|owg1|p800|pan(a|d|t)|pdxg|pg(13|\-([1-8]|c))|phil|pire|pl(ay|uc)|pn\-2|po(ck|rt|se)|prox|psio|pt\-g|qa\-a|qc(07|12|21|32|60|\-[2-7]|i\-)|qtek|r380|r600|raks|rim9|ro(ve|zo)|s55\/|sa(ge|ma|mm|ms|ny|va)|sc(01|h\-|oo|p\-)|sdk\/|se(c(\-|0|1)|47|mc|nd|ri)|sgh\-|shar|sie(\-|m)|sk\-0|sl(45|id)|sm(al|ar|b3|it|t5)|so(ft|ny)|sp(01|h\-|v\-|v )|sy(01|mb)|t2(18|50)|t6(00|10|18)|ta(gt|lk)|tcl\-|tdg\-|tel(i|m)|tim\-|t\-mo|to(pl|sh)|ts(70|m\-|m3|m5)|tx\-9|up(\.b|g1|si)|utst|v400|v750|veri|vi(rg|te)|vk(40|5[0-3]|\-v)|vm40|voda|vulc|vx(52|53|60|61|70|80|81|83|85|98)|w3c(\-| )|webc|whit|wi(g |nc|nw)|wmlb|wonu|x700|yas\-|your|zeto|zte\-/i.test(navigator.userAgent.substr(0,4))) { 
    isM = true;
}

function moveAnalysis(mvs, area){
    console.log(area);
    $.each(mvs, function(i, v){
        console.log(i,v.t,v.x,v.y);
    });

    return 10;
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
    $(c).mousemove(function(e){
        m = new mvd(Date.now(),e.pageX,e.pageY);    // harvest time, x, y for each move
        ms.push(m);                                 // ... and push it on the array
    }).click(function(e){
        m = new mvd(Date.now(),e.pageX,e.pageY);    // time, x, y for click event
        ms.push(m);                                 // ... and push it on the array
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
        data.mobile = isM;                  // is this a mobile device?
        data.subject = c;                   // elem class where moves were derived from
        data.target = t;                    // move target object (see above)
        data.receiver = r;                  // where to push result to for callback
        data.samples = ms.length;           // number of movements captured
        data.moves = ms;                    // movement data array
   
        // Create the JSON Object
        jsonData = JSON.stringify(data);

        // Send it to the server
        $.post(url, jsonData, function(jsonData, status){
            //console.log("POST status is " + status);
        });

        // Get evaluation result back and push it in designated element --> to be revised.
        /*
        $.get(url, function(result) {
            if (result == 'ON') {
                alert('ON');
            } else if (result == 'OFF') {
                alert('OFF');
            } else {
                alert(result);
            }
        });
        */
       
       // Calculate perception result.
       // Is last move timestamp greater then "now"? -> 1
       // Is move data array empty? -> 0
       if(data.samples < 2) {
            p = 0;
       } else if(data.moves[0].t > data.moves[data.samples-1].t) {
            p = 1;
            // Are last move coords differnet than click coords? -> 2
       } else if(data.moves[data.samples-1].x != data.moves[data.samples-2].x && 
                data.moves[data.samples-1].y != data.moves[data.samples-2].y) {
            p = 2;
            // Is last move timestamp less then 2hrs? -> 3
        } else if(data.moves[data.samples-1].t + 5200 < Date.now()) {
            p = 6;
        } else {
            // Are there move coords outside target object dimentions? -> 4
            p = moveAnalysis(data.moves, data.target);
        }
       
       // Place perception result in appropriate element.
       // Check nature of target object clicked.
       // Add data store to target element with perception value.
       $(c).data(r.slice(1), p);
       // if <input type ?> get parent <form> and append <input type hidden> with value.
        if($(c).parent().get(0).tagName == "FORM") {
            $(c).parent().append("<input id=\""+r+"\" type=\"hidden\" value=\""+p.toString()+"\" />");
        };

    });
});