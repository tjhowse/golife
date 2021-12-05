function update() {
    var source = '/world',
        timestamp = (new Date()).getTime(),
        newUrl = source + '?_=' + timestamp;
    document.getElementById("world").src = newUrl;
    var source = '/brain',
        timestamp = (new Date()).getTime(),
        newUrl = source + '?_=' + timestamp;
    document.getElementById("brain").src = newUrl;
    setTimeout(update, 100);
}

update();


function request(endpoint) {
    var xhr = new XMLHttpRequest();
    xhr.open('get', endpoint, true);
    xhr.send();
}

$(document).ready(function() {
    $("img").on("click", function(event) {
        var x = event.pageX - this.offsetLeft;
        var y = event.pageY - this.offsetTop;
        req = '/click?id='+this.id+'&x=' + x + '&y=' + y;
        request(req);
        // alert("X Coordinate: " + x + " Y Coordinate: " + y + " " + this.id);
    });
});