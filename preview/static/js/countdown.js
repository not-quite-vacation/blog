function getTimeRemaining(endtime) {
  var t = Date.parse(endtime) - Date.parse(new Date());
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
    
    var countdownComponents = ["days", "hours", "minutes", "seconds"];

    var ccElements = countdownComponents.map(function(name) {
        var container = document.createElement("div");
        var timeSpan = document.createElement("span");
        var textDiv = document.createElement("div");
        textDiv.className = "smalltext";
        textDiv.innerText = name;
        timeSpan.className = name;
        container.appendChild(timeSpan);
        container.appendChild(textDiv);
        return container;
    });

    ccElements.forEach(function(el) {
        clock.appendChild(el);
    });

    function updateClock() {
        var t = getTimeRemaining(endtime);
        countdownComponents.forEach(function(name, i) {
            ccElements[i].firstChild.textContent = t[name]
        });

        if (t.total <= 0) {
            clearInterval(timeinterval);
        }
    }

    updateClock();
    var timeinterval = setInterval(updateClock, 1000);
}

var deadline = '2017-05-12';
initializeClock('countdown', deadline);