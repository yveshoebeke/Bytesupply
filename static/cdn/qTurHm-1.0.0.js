var url = new URL($('script').last().attr('src')), c = url.searchParams.get("c"), ms=new Array;

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
        var r=8;
        $(c).data("r", r);
        console.log("Result: "+$(c).data("r"));
        console.log("Samples: "+ms.length);
        console.log("t: " + $(c).position().top);
        console.log("l: " + $(c).position().left);
        console.log("w: " + $(c).width());
        console.log("h: " + $(c).height());

        console.log("Stored data:");
        ms.forEach(function(d){
            $("#moves").append("<tr><td>"+d.t+"</td><td>"+d.x+"</td><td>"+d.y+"</td></tr>")
            console.log(d.t+" -> "+d.x+" - "+d.y);
        });
    });
});