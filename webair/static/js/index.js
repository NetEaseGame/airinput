/* javascript */
$(function() {
    var $screen = $("#img-screen");


    /* refresh screen */
    var imageUrl = "http://10.242.134.91:21000";
    refreshImage = newRefreshFunc($screen, imageUrl, "time1");
    refreshImage(true);

    var pressed = false;
    var startPoint = null, endPoint = null;
    $screen
        .mousedown(function(event) {
            pressed = true;
            startPoint = coord(event);
        })
        .mouseup(function(event) {
            pressed = false;
            endPoint = coord(event);
            sendNow();
        })
        .mousemove(function(event) {
            if (pressed) {
                // just ignore
            }
        })
        .mouseout(function(event) {
            if (pressed) {
                endPoint= coord(event);
                sendNow();
                pressed = false;
            }
        });

    function sendNow() {
        var url = '/api/touch';
        var data = startPoint;
        var dist = distance(startPoint, endPoint);
        console.log("DIST="+dist);
        if (dist > 5.0){
            url = '/api/drag';
            data = {start: startPoint, end: endPoint};
        }
        $.ajax(url, {
            type: 'POST',
            data: {data: JSON.stringify(data)},
            cache: false,
            timeout: 10000,
            success: function() {
                console.log("good");
            },
            error: function() {
                console.log("bad");
            }
        });
    }
    function distance(start, end){
        var dist = Math.sqrt(
            (start.x-end.x)*(start.x-end.x) + (start.y-end.y)*(start.y-end.y));
        console.log('distance='+dist);
        return dist;
    }
    function coord(event) {
        var offset = $screen.offset();
        var scale = $screen[0].naturalHeight * 1.0 / $screen.height();
        var x = parseInt((event.pageX - offset.left) * scale, 10);
        var y = parseInt((event.pageY - offset.top) * scale, 10);
        return {x: x, y: y};
    }

});