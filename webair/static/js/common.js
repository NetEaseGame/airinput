/* javascript */
String.prototype.format = function() {
    var formatted = this;
    for (var arg in arguments) {
        formatted = formatted.replace("{" + arg + "}", arguments[arg]);
    }
    return formatted;
};

$(document).on("click", "input.click-select", function(e) {
    $(e.target).select();
});

function indof(x, arr) {
    if (arr.indexOf)
        return arr.indexOf(x);
    for (var i = 0; i < arr.length; i++) {
        if (arr[i] == x)
            return i;
    }
    return -1;
}

function indoff(x, arr, start) {
    if (arr.indexOf)
        return arr.indexOf(x, start);
    for (var i = start; i < arr.length; i++) {
        if (arr[i] == x)
            return i;
    }
    return -1;
}

function setCookie(c_name, value, expiredays) {
    var exdate = new Date();
    exdate.setDate(exdate.getDate() + expiredays);
    document.cookie = c_name + "=" + escape(value) +
        ((expiredays == null) ? "" : ";expires=" + exdate.toUTCString());
}

function getCookie(c_name) {
    if (document.cookie.length > 0) {
        c_start = indof(c_name + "=", document.cookie);
        if (c_start != -1) {
            c_start = c_start + c_name.length + 1;
            c_end = indoff(";", document.cookie, c_start);
            if (c_end == -1) c_end = document.cookie.length;
            return unescape(document.cookie.substring(c_start, c_end));
        }
    }
    return "";
}