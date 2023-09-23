const episodeInputEl = document.querySelector('input[name="episode"]');
const audioEl = document.getElementById('audio-node');
const btnStart = document.getElementById('btn-start');
const toolbar = document.getElementById('toolbar');

const clipParams = {
  from: 0,
  to: 0,
};

function onTimeMark(type) {
  clipParams[type] = Math.floor(audioEl.currentTime * 1000);
  updateTimeMarkers();
}

function updateTimeMarkers() {
  document.getElementById('time-label-from').innerText = formatTimestamp(clipParams.from);
  document.getElementById('time-label-to').innerText = formatTimestamp(clipParams.to);
}

function formatTimestamp(ms) {
  ms = Number(ms);
  let seconds = Math.floor(ms / 1000);
  let minutes = Math.floor(seconds / 60);
  let hours = Math.floor(minutes / 60);
  seconds -= minutes * 60;
  minutes -= hours * 60;

  return `${padTimeSegment(hours)}:${padTimeSegment(minutes)}:${padTimeSegment(seconds)}`;
}

function padTimeSegment(str) {
  return String(str).padStart(2, '0');
}

btnStart.addEventListener('click', () => {
  //
});

episodeInputEl.addEventListener('keydown', (event) => {
  if (event.key === "Enter") {
    loadPodcastAudio();
  } else if (event.key === "Escape") {
    episodeInputEl.blur();
  }
});

episodeInputEl.addEventListener('blur', () => {
  loadPodcastAudio();
});

audioEl.addEventListener("canplay", function() {
  toolbar.removeAttribute('hidden');
});


function loadPodcastAudio() {
  const episode = parseInt(episodeInputEl.value, 10);

  if (Number.isInteger(episode)) {
    const url = `http://cdn.radio-t.com/rt_podcast${episode}.mp3`;

    if (audioEl.src !== url) {
      audioEl.src = url;
    }

  }
}

updateTimeMarkers();
