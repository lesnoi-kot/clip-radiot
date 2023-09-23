const episodeInputEl = document.querySelector('input[name="episode"]');
const audioEl = document.getElementById('audio-node');
const btnClip = document.getElementById('btn-clip');
const toolbar = document.getElementById('toolbar');
const errorContainer = document.getElementById('error-container');
const downloadSection = document.getElementById('download-section');

const appState = {
  episode: 0,
  from: 0,
  to: 0,

  error: null,
  episodeLoaded: false,
  clipURL: null,
  isFetching: false,

  update(diff) {
    Object.assign(this, diff);
    syncUI();
  }
};

function onTimeMark(type) {
  appState[type] = Math.floor(audioEl.currentTime * 1000);

  if (type === "to" && appState.to < appState.from) {
    appState.update({
      from: appState.to,
      to: appState.from,
    });
  }
}

function rewindAudio(delta) {
  audioEl.currentTime += delta;
}

function syncUI() {
  const { error, episodeLoaded, clipURL } = appState;

  errorContainer.innerText = error ? error.message : '';
  setElementVisible(errorContainer, !!error);
  setElementVisible(toolbar, episodeLoaded);
  setElementVisible(downloadSection, !!clipURL);

  updateTimeMarkers();
}

function updateTimeMarkers() {
  document.getElementById('time-label-from').innerText =
    formatTimestamp(appState.from);
  document.getElementById('time-label-to').innerText =
    formatTimestamp(appState.to);
}

btnClip.addEventListener('click', () => {
  const endpoint = new URL('/api/cut', document.location.origin);
  endpoint.searchParams.append('from', appState.from);
  endpoint.searchParams.append('to', appState.to);
  endpoint.searchParams.append('episode', appState.episode);

  appState.update({ error: null, isFetching: true });

  fetch(endpoint)
    .then(response => {
      if (response.ok) {
        return response.blob()
      }

      return response.json().then(body => {
        throw new Error(body.message);
      });
    })
    .then(audioBlob => {
      const audioURL = URL.createObjectURL(audioBlob);
      document.getElementById('download-link').href = audioURL;
      document.getElementById('audio-preview').src = audioURL;
      appState.clipURL = audioURL;
    })
    .catch((error) => {
      appState.error = error;
    }).finally(() => {
      syncUI();
    });
});

episodeInputEl.addEventListener('keydown', (event) => {
  if (event.key === 'Enter' || event.key === 'Escape') {
    episodeInputEl.blur();
  }
});

episodeInputEl.addEventListener('blur', () => {
  const episode = parseInt(episodeInputEl.value, 10);

  if (Number.isInteger(episode) && episode !== appState.episode) {
    appState.update({ from: 0, to: 0, episode });
    audioEl.src = `https://cdn.radio-t.com/rt_podcast${episode}.mp3`;
  }
});

audioEl.addEventListener('canplay', function () {
  appState.update({ episodeLoaded: true });
});

audioEl.onerror = () => {
  appState.update({
    error: new Error('Выпуск не найден!'),
    episodeLoaded: false,
    clipURL: null,
  });
}

syncUI();

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

function setElementVisible(el, visible) {
  if (visible) {
    el.classList.remove('hidden');
  } else {
    el.classList.add('hidden');
  }
}
