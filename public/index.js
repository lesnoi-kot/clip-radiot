const episodeInputEl = document.querySelector('input[name="episode"]');
const audioEl = document.getElementById('audio-node');
const audioPreviewEl = document.getElementById('audio-preview');
const btnClip = document.getElementById('btn-clip');
const toolbar = document.getElementById('toolbar');
const errorContainer = document.getElementById('error-container');
const downloadSection = document.getElementById('download-section');
const downloadLink = document.getElementById('download-link');

let episodeChangeTimerId;

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
  } else {
    syncUI();
  }
}

function rewindAudio(delta) {
  audioEl.currentTime += delta;
}

function syncUI() {
  const { error, episodeLoaded, clipURL, isFetching } = appState;

  errorContainer.innerText = error ? error.message : '';
  setElementVisible(errorContainer, !!error);
  setElementVisible(toolbar, episodeLoaded && episodeInputEl.value !== '');
  setElementVisible(downloadSection, !!clipURL);

  if (isFetching) {
    btnClip.innerText = 'Загрузка...';
    btnClip.disabled = true;
  } else {
    btnClip.innerText = '✂ Кат!';
    btnClip.disabled = false;
  }

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
      downloadLink.href = audioURL;
      downloadLink.download = `Radio-T ${appState.episode} (from ${replaceColons(formatTimestamp(appState.from))} to ${replaceColons(formatTimestamp(appState.to))}).mp3`;
      audioPreviewEl.src = audioURL;
      appState.clipURL = audioURL;
    })
    .catch((error) => {
      appState.error = error;
    }).finally(() => {
      appState.isFetching = false;
      syncUI();
    });
});

function handleEpisodeInputValue() {
  if (episodeInputEl.value === '') {
    return;
  }

  const episode = parseInt(episodeInputEl.value, 10);

  if (!Number.isInteger(episode)) {
    appState.update({ error: new Error('Введен некорректный номер выпуска') });
    return;
  }

  if (episode !== appState.episode) {
    appState.update({ from: 0, to: 0, episode });
    audioEl.src = `https://cdn.radio-t.com/rt_podcast${episode}.mp3`;
  } else {
    appState.update({ error: null });
  }
}

episodeInputEl.addEventListener('keydown', (event) => {
  if (event.key === 'Enter' || event.key === 'Escape') {
    episodeInputEl.blur();
  }
});

episodeInputEl.oninput = () => {
  clearTimeout(episodeChangeTimerId);
  episodeChangeTimerId = setTimeout(handleEpisodeInputValue, 1000);
}

episodeInputEl.addEventListener('blur', handleEpisodeInputValue);

audioEl.oncanplay = () => {
  appState.update({ episodeLoaded: true });
};

audioEl.onerror = () => {
  appState.update({
    error: new Error('Выпуск не найден!'),
    episodeLoaded: false,
    clipURL: null,
  });
}

audioPreviewEl.oncanplay = () => {
  appState.update({ isFetching: false });
}

audioPreviewEl.onerror = () => {
  appState.update({ isFetching: false, error: new Error('Выпуск не найден!') });
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

function replaceColons(str, to) {
  return String(str).replaceAll(':', '-');
}
