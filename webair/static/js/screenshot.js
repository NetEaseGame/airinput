/* javascript for screenshot */

// Example: fn = newRefreshFunc($("#screenshot"), "http://loalhost:5000", "time1")
// fn(false);

function newRefreshFunc($screen, imageAddr, cookieKey) {
    var images = [];
    var numOfImages = 3;
    var counter = 0;
    var lasttime = 0;

    for (var i = 0; i < numOfImages; i++) {
        images[i] = $('<img />');
    }

    var startnextimage = function(i, waitfordiff) {
        if (document.getElementById('auto').checked === true || waitfordiff === false) {
            var now = new Date();
            var t = document.getElementById('png').checked ? 'png' : 'jpg';
            var or = document.getElementById('orient').checked ? 'h' : 'v';
            // var lowr = document.getElementById('lowres').checked ? 'l' : 'n';
            // var f = document.getElementById('first').checked ? 'f' : 'n';
            // var fli = document.getElementById('flip').checked ? 'f' : 'n';
            var w = waitfordiff ? 'w' : 'n';
            var asfile = 'n';
            var addr = imageAddr + '/screen.' + t +
                '?cookiename=' + cookieKey + '&rand=' + or + w + now.getTime() + i;
            $(images[i])
                .unbind()
                .load(function() {
                    var ctime = getCookie("timestamp");
                    if (ctime >= lasttime) {
                        $screen.attr('src', addr);
                        lasttime = ctime;
                    }
                    var now = new Date();
                    counter++;
                    setTimeout(function() {
                        startnextimage(i, true);
                    }, 100);
                })
                .error(function() {
                    setTimeout(function() {
                        startnextimage(i, false);
                    }, 1000);
                })
                .attr('src', addr);
        }
    };
    return function(waitfordiff) {
        for (var i = 0; i < numOfImages; i++) {
            startnextimage(i, waitfordiff ? true : i !== 0);
        }
    };
}