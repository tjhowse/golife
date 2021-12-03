function update() {
    var source = '/img',
        timestamp = (new Date()).getTime(),
        newUrl = source + '?_=' + timestamp;
    document.getElementById("world").src = newUrl;
    setTimeout(update, 100);
}

update();