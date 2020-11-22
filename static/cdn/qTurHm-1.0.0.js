var url = new URL($('script').last().attr('src')); 
var c = url.searchParams.get("c"); 
var k = url.searchParams.get("k");
var r = url.searchParams.get("r");
console.log(c)
console.log(k)
console.log(r)
//var k = "21101956";
//var r = "qTurHm_Result";
ms = new Array;

function mvd(t,x,y){
    this.t=t;
    this.x=x;
    this.y=y;
}

$(function() {
    c = "." + c;
    $(c).data("r", 0);
    $(c).mousemove(function(e){
        m = new mvd(Date.now(),e.pageX,e.pageY);
        ms.push(m);
    }).click(function(){
        $(c).unbind("mousemove").unbind("click");
        ms.forEach(function(d){
            $("#moves").append("<tr><td>"+d.t+"</td><td>"+d.x+"</td><td>"+d.y+"</td></tr>")
        });

        // Create JSON
        var t = {};
        t.top = ~~$(c).position().top;
        t.left = ~~$(c).position().left;
        t.width = ~~$(c).width();
        t.height = ~~$(c).height();

        var data = {};
        data.userkey = k;
        data.timestamp = Date.now();
        data.origURL = window.location.href;
        data.target = t;
        data.reciever = r;
        data.samples = ms.length;
        data.moves = ms;
   
        jsonData = JSON.stringify(data);

        $.post("https://bytesupply.com/api/v1/qTurHm", jsonData, function(jsonData, status){
            console.log(`${jsonData} and status is ${status}`);
        });
    });
});