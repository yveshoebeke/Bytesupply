var url = new URL($('script').last().attr('src')), c = url.searchParams.get("c"), ms=new Array;
k = "75cab53138db56817e12500dcdb06e6c61500ea6f0085ca54946b87ec1c4d409";
r = "qTurHm_Result";

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
        t.top = $(c).position().top;
        t.left = $(c).position().left;
        t.width = $(c).width();
        t.height = $(c).height();

        var data = {};
        data.userkey = k;
        data.timestamp = Date();
        data.origURL = window.location.href;
        data.target = t;
        data.reciever = r;
        data.samples = ms.length;
        data.moves = ms;
   
        userJson = JSON.stringify(data);
        alert(userJson);

        $.post("https://bytesupply.com/api/v1/qTurHm", function(data, status){
            alert("Data: " + data + "\nStatus: " + status);
        });
    });
});