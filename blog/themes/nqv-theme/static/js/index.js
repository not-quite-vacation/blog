function requestInterval(fn, delay) {
    var start = new Date().getTime(),
        handle = new (Object);
    function loop() {
        var current = new Date().getTime(),
            delta = current - start;
        if (delta >= delay) {
            fn.call();
            start = new Date().getTime();
        }
        handle.value = requestAnimationFrame(loop);
    };
    handle.value = requestAnimationFrame(loop);
    return handle;
}
function clearRequestInterval(handle) {
    cancelAnimationFrame(handle.value);
}
function getTimeRemaining(endtime) {
    var t = Date.parse(new Date()) - Date.parse(endtime);
    var seconds = Math.floor((t / 1000) % 60);
    var minutes = Math.floor((t / 1000 / 60) % 60);
    var hours = Math.floor((t / (1000 * 60 * 60)) % 24);
    var days = Math.floor(t / (1000 * 60 * 60 * 24));
    return {
        'total': t,
        'days': days,
        'hours': hours,
        'minutes': minutes,
        'seconds': seconds
    };
}

function initializeClock(id, endtime) {
    var clock = document.getElementById(id);
    if (clock === null) {
        return;
    }
    var daysSpan = clock.querySelector('.days');
    var hoursSpan = clock.querySelector('.hours');
    var minutesSpan = clock.querySelector('.minutes');
    var secondsSpan = clock.querySelector('.seconds');

    function updateClock() {
        var t = getTimeRemaining(endtime);

        daysSpan.innerHTML = t.days;
        hoursSpan.innerHTML = ('0' + t.hours).slice(-2);
        minutesSpan.innerHTML = ('0' + t.minutes).slice(-2);
        secondsSpan.innerHTML = ('0' + t.seconds).slice(-2);

        if (t.total <= 0) {
            clearRequestInterval(timeinterval);
        }
    }

    updateClock();
    var timeinterval = requestInterval(updateClock, 1000);
}

var deadline = '2017-05-17';
initializeClock('clockdiv', deadline);

var toggler = (selector) => {
    var nav = document.querySelector(selector);
    return () => {
        nav.classList.toggle("open");
    }
};
var navToggle = toggler(".header > nav");