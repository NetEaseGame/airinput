/* javascript */
$(function() {
    var $screen = $("#img-screen");


    refreshFuncs = [];
    /* refresh screen */
    $("img.screen").each(function(ind, value){
        var $screen = $(value);
        var refresh = newRefreshFunc($screen, $screen.attr("src"), "time"+ind);
        var serialno = $screen.attr("serialno");

        refresh(true);
        refreshFuncs.push(refresh);

        var pressed = false;
        var startPoint = null,
            endPoint = null;
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
                    endPoint = coord(event);
                    sendNow();
                    pressed = false;
                }
            });

        function sendNow(){
            if (document.getElementById('sync').checked){
                $("img.screen").each(function(ind, val){
                    var sno = $(val).attr("serialno");
                    send(sno, startPoint, endPoint);
                });
            } else {
                send(serialno, startPoint, endPoint);
            }
            
        }

        function send(serialno, start, end) {
            var url = '/api/'+serialno+'/touch';
            var data = start;
            var dist = distance(start, end);
            console.log("DIST=" + dist);
            if (dist > 5.0) {
                url = '/api/'+serialno+'/drag';
                data = {
                    start: start,
                    end: end
                };
            }
            $.ajax(url, {
                type: 'POST',
                data: {
                    data: JSON.stringify(data)
                },
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

        function distance(start, end) {
            if (start === undefined || end === undefined){
                return 0;
            }
            var dist = Math.sqrt(
                (start.x - end.x) * (start.x - end.x) + (start.y - end.y) * (start.y - end.y));
            console.log('distance=' + dist);
            return dist;
        }

        function coord(event) {
            var offset = $screen.offset();
            var scale = $screen[0].naturalHeight * 1.0 / $screen.height();
            var x = parseInt((event.pageX - offset.left) * scale, 10);
            var y = parseInt((event.pageY - offset.top) * scale, 10);
            return {
                x: x,
                y: y
            };
        }
        
    });

    refreshAll = function(){
        for (var x in refreshFuncs){
            refreshFuncs[x](true);
        }
    };

});