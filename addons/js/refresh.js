function updateClock() {
  var now = new Date();
      months = ['Jan', 'Feb', 'Mar', 'Apr', 'May', 'Jun', 'Jul', 'Aug', 'Sep', 'Oct', 'Nov', 'Dec'];
      weekday = ['Sun', 'Mon', 'Tue', 'Wed', 'Thu', 'Fri', 'Sat']
      time = ('0' + now.getHours()).slice(-2) + ':' + ('0' + now.getMinutes()).slice(-2);

      date = [weekday[now.getDay()], 
              months[now.getMonth()],
              now.getDate(),
              now.getFullYear()].join(' ');

  document.getElementById('timeCell').innerHTML = time;
  document.getElementById('dateCell').innerHTML = date;

  setTimeout(updateClock, 1000);
}

updateClock();
