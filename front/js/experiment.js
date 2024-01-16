import * as timelines from "./timelines";

const state = {};

const shared = {
  inBlock: (done, trialsPerBlock) => Math.floor(done / trialsPerBlock),
  updateProgress: () => {
    document.getElementById("progress").innerHTML = `${state.position.trial + 1}/${state.totalLength}`;
  },
  showProgress: () => {
    if (state.settings.showProgress) {
      document.getElementById("progress").style = "display: block;";
    }
  },
  hideProgress: () => {
    document.getElementById("progress").style = "display: none;";
  }
};

const pairs = (arr) =>
  Array.from(new Array(Math.ceil(arr.length / 2)), (_, i) => {
    const pair = arr.slice(i * 2, i * 2 + 2);
    return { s1: pair[0], s2: pair[1] };
  });

// build and run experiment
export default (props, ws) => {
  // init
  const { settings, wording, participant } = props;
  const jsPsych = initJsPsych({
    display_element: "jspsych-root",
    on_finish: function () {
      state.jsPsych.data.displayData();
    },
  });
  const todo = participant.todo.split(",");
  const blocks = settings.addRepeatBlock
    ? settings.blocksPerXp + 1
    : settings.blocksPerXp;

  const totalLength = blocks * settings.trialsPerBlock;
  const remainingLength = todo.length / 2; // 2 choices per trial
  const previouslyDoneLength = totalLength - remainingLength; // not 0 if user reconnects (page refresh for instance)
  const position = {
    trial: previouslyDoneLength,
    block: shared.inBlock(previouslyDoneLength, settings.trialsPerBlock),
  };
  const timeline = [];

  // shared state
  state.ws = ws;
  state.settings = settings;
  state.wording = wording;
  state.start = new Date().toISOString();
  state.stimuli = pairs(todo, 2);
  state.jsPsych = jsPsych;
  state.position = position;
  state.totalLength = totalLength;
  state.previouslyDoneLength = previouslyDoneLength;

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
    if (participant.age.length === 0 || participant.sex.length === 0) {
      timeline.push({
        type: jsPsychSurveyHtmlForm,
        preamble: `<p>${wording.collect}</p>`,
        html: `<p>
          <fieldset>
            <label>${wording.collectAge}</label>
            <input id="age" name="age" type="text" minlength="2" maxlength="3" pattern="[0-9]*" required />
          </fieldset>
          <fieldset>
            <label>${wording.collectSex}</label>
            <input name="sex" type="text" maxlength="16" required />
          </fieldset>
        </p>`,
        autofocus: "age",
        button_label: wording.collectButton,
        on_finish: (data) => {
          ws.send(
            JSON.stringify({
              kind: "info",
              payload: JSON.stringify(data.response),
            })
          );
        },
      });
    }
    // display welcoming (or welcoming back) message
    const welcomingMessage =
      previouslyDoneLength == 0 ? wording.introduction : wording.resume;
    timeline.push({
      type: jsPsychHtmlKeyboardResponse,
      stimulus: `<p>${welcomingMessage}</p>`,
      prompt: `<p><span class='strong'>[${wording.space}]</span> ${wording.next}</p>`,
      choices: " ",
    });


    // choose main timeline
    const mainTimeline =
      settings.kind === "sound" ? timelines.sounds : timelines.images;
    timeline.push(mainTimeline(state, shared));

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
