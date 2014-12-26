/* javascript */
$(function() {
    var $screen = $("#img-screen");
    var $message = $("#message");
    var $codearea = $("#ipt-code");
    var $cropname = $("#ipt-cropname");

    var $connect = $("#btn-connect");

    var $dataX = $("#dataX"),
        $dataY = $("#dataY"),
        $dataWidth = $("#dataWidth"),
        $dataHeight = $("#dataHeight"),
        $dataMiddle = $("#dataMiddle"),
        $dataSiftCount = $("#dataSiftCount");

    /* refresh screen */
    var imageUrl = "http://10.242.116.68:21000";
    refreshImage = newRefreshFunc($screen, imageUrl, "time1");
    refreshImage(true);

    var pressed = false;
    var moveevent = [];
    var Mqueue = false;
    $screen
        .mousedown(function(event) {
            pressed = true;
            moveevent.push(coord(event, pressed)); //moveevent + '_' + coord(event) + '_1';
            moveLater();
            // send(event, pressed);
        })
        .mouseup(function(event) {
            pressed = false;
            moveevent.push(coord(event, pressed)); //moveevent + '_' + coord(event) + '_1';
            // moveevent = moveevent + '_' + coord(event) + '_0';
            // send(event, pressed);
            sendNow();
        })
        .mousemove(function(event) {
            if (pressed) {
                moveevent.push(coord(event, pressed)); //moveevent + '_' + coord(event) + '_1';
                // moveevent = moveevent + '_' + coord(event) + '_1';
                sendNow();
                // send(event, pressed);
            }
        })
        .mouseout(function(event) {
            if (pressed) {
                pressed = false;
                moveevent.push(coord(event, pressed)); //moveevent + '_' + coord(event) + '_1';
                // moveevent = moveevent + '_' + coord(event) + '_0';
                sendNow();
                // send(event, pressed);
            }
        });

    function sendNow() {
        if (moveevent !== '' && Mqueue === false) {
            Mqueue = true;
            $.ajax('/api/touch', {
                type: 'POST',
                data: {'event': JSON.stringify(moveevent)},
                cache: false,
                timeout: 10000,
                success: function() {
                    Mqueue = false;
                    sendNow();
                },
                error: function() {
                    MQueue = false;
                    sendNow();
                }
            });
            moveevent = '';
        }
    }

    function moveLater() {
        if (Mqueue === false) {
            Mqueue = true;
            setTimeout(function() {
                Mqueue = false;
                sendNow();
            }, 300);
        }
    }

    function send(event, pressed){
        var moveevent = coord(event, pressed);
        $.ajax('/api/touch', {
                type: 'POST',
                data: {'event': JSON.stringify(moveevent)},
                cache: false,
                timeout: 10000,
                success: function() {
                    console.log("send");
                }
            });
    }
    function coord(event, pressed) {
        var offset = $screen.offset();
        var scale = $screen[0].naturalHeight * 1.0 / $screen.height();
        var x = parseInt((event.pageX - offset.left) * scale, 10);
        var y = parseInt((event.pageY - offset.top) * scale, 10);
        return {x: x, y: y, pressed: pressed};

        //  var top = 0;
        //  var left = 0;
        //  var ob = document.images['screenshot'];
        //  do{
        //      left += ob.offsetLeft;
        //      top += ob.offsetTop;
        //      ob = ob.offsetParent;
        //  }while (ob);
        //  pos_x = event.offsetX?(event.offsetX):event.pageX-left;
        //  pos_y = event.offsetY?(event.offsetY):event.pageY-top;
        // var sc = $("#screenshot");
        // var pos_x = Math.floor(event.offsetX ? (event.offsetX) : event.pageX - sc.offset().left);
        // var pos_y = Math.floor(event.offsetY ? (event.offsetY) : event.pageY - sc.offset().top);
        // if (document.getElementById('halfsize').checked === true) {
        //     pos_x = pos_x * 2;
        //     pos_y = pos_y * 2;
        // }
        // if (document.getElementById('orient').checked === true)
        //     return 'h' + pos_x + '_' + pos_y;
        // else
        //     return 'v' + pos_x + '_' + pos_y;
    }
    // $screen.attr("src", "/screen.jpg?random="+new Date().getTime());
    // $("#refresh").click(function(){loadscreen()});



    var cropperData = {};
    var screenFilename = null;

    var updateSiftCount = function(data) {
        $.ajax({
            url: "/api/cropcheck",
            dataType: "json",
            data: $.extend(data, {
                screen: screenFilename
            }),
            success: function(e) {
                console.log(e);
                $dataSiftCount.text(e.siftcnt);
            },
            error: function(e) {
                console.log(e);
            }
        });
    };
    var cropperImage = function(source) {
        // $screen.cropper("reset");
        screenFilename = source;
        $screen.attr("src", source);
        $screen.cropper({
            resizable: true,
            preview: ".img-preview",
            autoCropArea: 0.3,
            data: cropperData,
            done: function(data) {
                var $preview = $(".img-preview");
                $preview.css("width", 150 * data.width / data.height + "px");
                $preview.css("height", 300 * data.height / data.width + "px");
                $dataX.text(data.x);
                $dataY.text(data.y);
                $dataWidth.text(data.width);
                $dataHeight.text(data.height);
                var middle = {
                    x: parseInt(data.x + data.width / 2, 10),
                    y: parseInt(data.y + data.height / 2, 10)
                };
                $dataMiddle.text('{0}, {1}'.format(middle.x, middle.y));
                cropperData = {
                    x: data.x,
                    y: data.y,
                    width: data.width,
                    height: data.height,
                    middle: middle
                };
            },
            dragend: function(data) {
                updateSiftCount(cropperData);
            }
        });
    };

    // cropperImage("/static/imgs/init.png");
    // cropperImage("http://10.242.116.68:21000/screen.jpg");

    $("#btn-refresh").click(function() {
        $message.text("taking snapshot");
        $.ajax({
            url: "/api/snapshot",
            dataType: "json",
            success: function(e) {
                $screen.cropper("destroy");
                cropperImage("/tmp/" + e.filename);
                $message.text("");
            },
            error: function(e) {
                $message.text("take snapshot failed for some reason");
            }
        });
    });

    $cropname.keydown(function(e) {
        if (e.keyCode == 13) { /* when press enter trigger crop */
            var ev = document.createEvent('MouseEvent');
            ev.initEvent('click', false, false);
            $("#btn-crop")[0].dispatchEvent(e);
        }
    });
    $("#btn-crop").click(function() {
        var filename = $.trim($cropname.val());
        if (filename === "") {
            console.log("canceled");
            $message.text("empty filename");
            return;
        }
        if (filename.indexOf(".") == -1) {
            filename = filename + ".png";
            $cropname.val(filename);
        }
        $.ajax({
            url: "/api/crop",
            dataType: "json",
            data: $.extend(cropperData, {
                screen: screenFilename,
                filename: filename
            }),
            success: function(e) {
                $message.text(e.message);
                $codearea.val("app.click(u'{0}')".format($cropname.val()));
            },
            error: function(e) {
                $message.text("crop image failed for some reason");
            }
        });
    });

    $("#btn-click-point").click(function() {
        var x = cropperData.middle.x;
        var y = cropperData.middle.y;
        $codearea.val("app.click(({0}, {1}))".format(x, y));
    });

    $("#btn-run").click(function() {
        $message.text("run in background ...");
        $.ajax({
            url: "/api/run",
            data: {
                code: $codearea.val()
            },
            dataType: "json",
            success: function(e) {
                $message.text(e.message);
            },
            error: function(e) {
                $message.text("run failed for some reason");
            }
        });
    });

    var $devno = $("#ipt-devno");
    var $device = $("#slt-device");
    $connect.click(function() {
        $message.text("connecting ...");
        $.ajax({
            url: "/api/connect",
            data: {
                devno: $devno.val(),
                device: $device.val()
            },
            dataType: "json",
            success: function(e) {
                $message.text(e.message);
                $connect.text("重新连接");
            },
            error: function(e) {
                $message.text("connect failed for some reason");
            }
        });
    });

});