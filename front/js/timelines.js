// Internals

const ASSET_PREFIX = "../assets/";

const genBlockStop = (state, shared) => {
  const { jsPsych, previouslyDoneLength, settings, totalLength, wording} = state;
  const { hideProgress} = shared;
  return {
    timeline: [
      {
        type: jsPsychHtmlKeyboardResponse,
        stimulus: `<p>${wording.pause}</p>`,
        prompt: "",
        choices: "NO_KEYS",
        trial_duration: 4000,
      },
      {
        type: jsPsychHtmlKeyboardResponse,
        stimulus: `<p>${wording.pauseOver}</p>`,
        prompt: `<p><span class='strong'>[${wording.space}]</span> ${wording.next}</p>`,
        choices: " ",
      },
    ],
    on_start: hideProgress,
    conditional_function: function () {
      const done =
        jsPsych.data.get().filter({ answered: true }).count() +
        previouslyDoneLength;
      const blockEnd = done % settings.trialsPerBlock === 0;
      // display if end of block and if not last block
      return blockEnd && done !== totalLength;
    },
  }
}

// higher-order function 
const trialSubmitTwoIntervals = (state, shared) => {
  const { jsPsych, position, settings, start, wording, ws } = state;
  const { inBlock, updateProgress } = shared;
  
  return (data) => {
    const result1 = {
      trial: position.trial.toString(),
      block: position.block.toString(),
      stimulus: jsPsych.timelineVariable("s1"),
      order: "0",
      response: data.response === wording.choice1 ? "True" : "False",
      rt: data.rt.toString(),
      date: start,
    };
    const result2 = {
      trial: position.trial.toString(),
      block: position.block.toString(),
      stimulus: jsPsych.timelineVariable("s2"),
      order: "1",
      response: data.response === wording.choice2 ? "True" : "False",
      rt: data.rt.toString(),
      date: start,
    };
    position.trial++;
    position.block = inBlock(position.trial, settings.trialsPerBlock);
    updateProgress();
  
    ws.send(
      JSON.stringify({
        kind: "trial",
        payload: JSON.stringify({ result1, result2 }),
      })
    );
  }
};


// API

export const sounds = (state, shared) => {
  const { jsPsych, wording, stimuli } = state;
  const { showProgress } = shared;

  const blockStop = genBlockStop(state, shared);
  return {
    // default value preventing text flickering when loading assets
    prompt: `<p>[${wording.space}] <span style='font-weight:bold'> ${wording.playSounds}</span></p>
    <p>${wording.question}</p>
    <div class='sound-choice'>
      <div>[${wording.choice1}] ${wording.label1}</div>
      <div>${wording.label2} [${wording.choice2}]</div>
    </div>`,
    // actual timeline
    timeline: [
      {
        type: jsPsychPreload,
        audio: () => {
          return [
            `${ASSET_PREFIX}${jsPsych.timelineVariable("s1")}`,
            `${ASSET_PREFIX}${jsPsych.timelineVariable("s2")}`,
          ];
        },
        show_progress_bar: false,
        post_trial_gap: 200,
        on_start: showProgress,
      },
      {
        type: jsPsychHtmlKeyboardResponse,
        stimulus: "",
        choices: " ",
        prompt: `<p><span style='font-weight:bold'>[${wording.space}]</span> ${wording.playSounds}</p>
        <p>${wording.question}</p>
        <div class='sound-choice'>
          <div>[${wording.choice1}] ${wording.label1}</div>
          <div>${wording.label2} [${wording.choice2}]</div>
        </div>`,
      },
      {
        type: jsPsychAudioKeyboardResponse,
        stimulus: () => `${ASSET_PREFIX}${jsPsych.timelineVariable("s1")}`,
        choices: "NO_KEYS",
        trial_ends_after_audio: true,
        response_allowed_while_playing: false,
      },
      {
        type: jsPsychHtmlKeyboardResponse,
        stimulus: "",
        choices: "NO_KEYS",
        trial_duration: 500,
      },
      {
        type: jsPsychAudioKeyboardResponse,
        stimulus: () => `${ASSET_PREFIX}${jsPsych.timelineVariable("s2")}`,
        choices: "NO_KEYS",
        trial_ends_after_audio: true,
        response_allowed_while_playing: false,
      },
      {
        type: jsPsychHtmlKeyboardResponse,
        stimulus: "",
        choices: [wording.choice1, wording.choice2],
        prompt: `<p>[${wording.space}] ${wording.playSounds}</p>
        <p>${wording.question}</p>
        <div class='sound-choice'>
          <div><span class='strong'>[${wording.choice1}]</span> ${wording.label1}</div>
          <div>${wording.label2} <span class='strong'>[${wording.choice2}]</span></div>
        </div>`,
        data: {
          answered: true,
        },
        on_finish: trialSubmitTwoIntervals(state, shared),
      },
      blockStop,
    ],
    timeline_variables: stimuli,
  };
};

export const images = (state, shared) => {
  const { jsPsych, settings, wording, stimuli } = state;
  const { showProgress } = shared;

  const blockStop = genBlockStop(state, shared);
  return {
    // default values
    css_classes: ["image"],
    // actual timeline
    timeline: [
      {
        type: jsPsychPreload,
        images: () => {
          return [
            `${ASSET_PREFIX}${jsPsych.timelineVariable("s1")}`,
            `${ASSET_PREFIX}${jsPsych.timelineVariable("s2")}`,
          ];
        },
        on_start: showProgress,
        show_progress_bar: false,
        post_trial_gap: 200,
      },
      {
        type: jsPsychHtmlKeyboardResponse,
        stimulus: "",
        choices: [wording.choice1, wording.choice2],
        prompt: () => {
          const imgWidth =
            settings.forceWidth.length == 0 ? "auto" : settings.forceWidth;
          return `<p>${wording.question}</p>
          <div class='image-choice'>
            <div>
              <img style="width:${imgWidth};" src="${ASSET_PREFIX}${jsPsych.timelineVariable(
              "s1"
            )}">
              <div><span class='strong'>[${wording.choice1}]</span> ${
              wording.label1
            }</div>
            </div>
            <div>
              <img style="width:${imgWidth};" src="${ASSET_PREFIX}${jsPsych.timelineVariable(
              "s2"
            )}">
              <div>${wording.label2} <span class='strong'>[${
              wording.choice2
            }]</span></div>
            </div>
          </div>`;
        },
        data: {
          answered: true,
        },
        on_finish: trialSubmitTwoIntervals(state, shared),
      },
      blockStop,
    ],
    timeline_variables: stimuli,
  };
};

