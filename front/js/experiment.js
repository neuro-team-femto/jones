import * as timelines from "./timelines";

const PREVENT_WS_CLOSE_MAX_IDLING = 10 * 60 * 1000; // x minutes in milliseconds
const PING_PERIOD = 55 * 1000; // x seconds in milliseconds

const state = {};

const updateActivity = () => {
  if (!state.intervalId) initPingInterval(); // start anti idling

  state.lastInteraction = new Date();
}

const isActive = () => {
  return (new Date() - state.lastInteraction) < PREVENT_WS_CLOSE_MAX_IDLING;
}

const initPingInterval = () => {
  state.intervalId = setInterval(() => {
    if (isActive()) { // prevent websocket close server-side for a bit
      state.ws.readyState === WebSocket.OPEN && state.ws.send( JSON.stringify({
        kind: "ping",
      }));
    } else { // stop preventing websocket close after PREVENT_WS_CLOSE_MAX_IDLING
      clearInterval(state.intervalId);
      delete(state.intervalId);
    }
  }, PING_PERIOD);
}

const shared = {
  inBlock: (done, trialsPerBlock) => Math.floor(done / trialsPerBlock),
  updateProgress: () => {
    document.getElementById("progress").innerHTML = `${
      state.position.trial + 1
    }/${state.totalLength}`;
  },
  showProgress: () => {
    if (state.settings.showProgress) {
      document.getElementById("progress").style = "display: block;";
    }
  },
  hideProgress: () => {
    document.getElementById("progress").style = "display: none;";
  },
};

const pairs = (arr) =>
  Array.from(new Array(Math.ceil(arr.length / 2)), (_, i) => {
    const pair = arr.slice(i * 2, i * 2 + 2);
    return { asset1: pair[0], asset2: pair[1] };
  });

// build and run experiment
export default (props, ws) => {
  // init
  const { settings, wording, participant } = props;

  // jspsych
  const jsPsych = initJsPsych({
    display_element: "jspsych-root",
    on_finish: function () {
      state.jsPsych.data.displayData();
    },
  });
  const blocks = settings.addRepeatBlock
    ? settings.blocksPerXp + 1
    : settings.blocksPerXp;

  const totalLength = blocks * settings.trialsPerBlock;
  let stimuli;
  let remainingLength = 0;
  if (!!participant.todo) {
    if (settings.nInterval === 1) {
      remainingLength = participant.todo.length;
      stimuli = participant.todo.map((t) => ({ asset: t}));
    } else {
      remainingLength = participant.todo.length / 2;
      stimuli = pairs(participant.todo, 2);
    }
  }
  const previouslyDoneLength = totalLength - remainingLength; // not 0 if user reconnects (page refresh for instance)
  const position = {
    trial: previouslyDoneLength,
    block: shared.inBlock(previouslyDoneLength, settings.trialsPerBlock),
  };
  const timeline = [];

  // ws
  ws.onerror = (event) => {
    console.error("[ws error] ", event);
    jsPsych.endExperiment(wording.connectionError || "Connexion perdue, veuillez rafraîchir la page");
  };
  
  ws.onclose = (event) => {
    console.log("[ws closed] ", event);
    jsPsych.endExperiment(wording.connectionError || "Connexion perdue, veuillez rafraîchir la page");
  };

  // shared state
  state.ws = ws;
  state.settings = settings;
  state.wording = wording;
  state.start = new Date().toISOString();
  state.stimuli = stimuli;
  state.jsPsych = jsPsych;
  state.position = position;
  state.totalLength = totalLength;
  state.previouslyDoneLength = previouslyDoneLength;
  updateActivity();

  // UX and timeline

  shared.updateProgress();

  // experiment has already be fully run by this participant
  if (remainingLength === 0) {
    // display closed message
    timeline.push({
      type: jsPsychHtmlKeyboardResponse,
      stimulus: `<h3>${wording.closed}</h3>`,
      choices: "NO_KEYS",
    });
  } else {
    // form to collect participant info
    if (
      settings.collectInfo &&
      settings.collectInfo.length > 0 &&
      !participant.infoCollected
    ) {
      timeline.push({
        type: jsPsychSurveyHtmlForm,
        preamble: `<p>${wording.collect}</p>`,
        html: () => {
          const fieldsets = settings.collectInfo
            .map((i) => {
              let pattern = !!i.pattern ? `pattern="${i.pattern}"` : "";
              let min = !!i.min ? `min="${i.min}"` : "";
              let max = !!i.max ? `max="${i.max}"` : "";
              let step = !!i.min && !!i.max ? `step="1"` : "";
              return ` <fieldset>
              <label>${i.label}</label>
              <input id="${i.key}" name="${i.key}" type="${i.inputType}" ${pattern} ${min} ${max} ${step} required />
            </fieldset>`;
            })
            .join("");
          return `<p>${fieldsets}</p>`;
        },
        autofocus: "age",
        button_label: wording.collectButton,
        on_finish: (data) => {
          ws.readyState === WebSocket.OPEN && ws.send(
            JSON.stringify({
              kind: "info",
              payload: JSON.stringify(data.response),
            })
          );
          updateActivity();
        },
      });
    }
    // display welcoming (or welcoming back) message
    const welcomingMessage =
      previouslyDoneLength == 0 ? wording.introduction : wording.resume;
  
    if ((typeof(wording.introduction) === "string") || (wording.introduction.length === 1)) {
      timeline.push({
        type: jsPsychHtmlKeyboardResponse,
        stimulus: `<p>${welcomingMessage}</p>`,
        prompt: `<p><span class='strong'>[${wording.space}]</span> ${wording.next}</p>`,
        choices: " ",
        on_finish: updateActivity
      })  
    // multi-page instructions
    // if introduction field in wording.run is list with multiple elements, show one-by-one
    } else {   
      wording.introduction.forEach((i) => {
        timeline.push({
          type: jsPsychHtmlKeyboardResponse,
          stimulus: `<p>${i}</p>`,
          prompt: `<p><span class='strong'>[${wording.space}]</span> ${wording.next}</p>`,
          choices: " ",
          on_finish: updateActivity
        })
      })
    }

    // choose main timeline
    let mainTimeline;
    if (settings.nInterval === 1) {
      mainTimeline =
        settings.kind === "video" ? timelines.videosN1 : settings.kind === "sound" ? timelines.soundsN1 : timelines.imagesN1;
    } else {
      mainTimeline =
        settings.kind === "sound" ? timelines.soundsN2 : timelines.imagesN2;
    }
    timeline.push(mainTimeline(state, shared, updateActivity));

    // display end message
    timeline.push({
      type: jsPsychHtmlKeyboardResponse,
      stimulus: `<h3>${wording.end}</h3>`,
      prompt: `<p>${wording.thanks}</p>`,
      choices: "NO_KEYS",
      on_start: function () {
        shared.hideProgress();
        console.log("The experiment is over");
      },
    });
  }

  state.jsPsych.run(timeline);
};
