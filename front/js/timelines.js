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
        choices: " ", // press space to next step
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
const trialSubmitN1 = (state, shared) => {
  const { jsPsych, position, settings, start, wording, ws } = state;
  const { inBlock, updateProgress } = shared;
  
  return (data) => {
    const result = {
      trial: position.trial.toString(),
      block: position.block.toString(),
      stimulus: jsPsych.timelineVariable("asset"),
      order: "0",
      response: data.response === wording.keyAlt1 ? wording.labelAlt1 : wording.labelAlt2,
      rt: data.rt.toString(),
      date: start,
    };
    position.trial++;
    position.block = inBlock(position.trial, settings.trialsPerBlock);
    updateProgress();
  
    ws.send(
      JSON.stringify({
        kind: "trial",
        payload: JSON.stringify(result),
      })
    );
  }
};

const trialSubmitN2 = (state, shared) => {
  const { jsPsych, position, settings, start, wording, ws } = state;
  const { inBlock, updateProgress } = shared;
  
  return (data) => {
    const result1 = {
      trial: position.trial.toString(),
      block: position.block.toString(),
      stimulus: jsPsych.timelineVariable("asset1"),
      order: "0",
      response: data.response === wording.keyAlt1 ? "True" : "False",
      rt: data.rt.toString(),
      date: start,
    };
    const result2 = {
      trial: position.trial.toString(),
      block: position.block.toString(),
      stimulus: jsPsych.timelineVariable("asset2"),
      order: "1",
      response: data.response === wording.keyAlt2 ? "True" : "False",
      rt: data.rt.toString(),
      date: start,
    };
    position.trial++;
    position.block = inBlock(position.trial, settings.trialsPerBlock);
    updateProgress();

    ws.send(
      JSON.stringify({
        kind: "trial",
        payload: JSON.stringify(result1),
      })
    );
    ws.send(
      JSON.stringify({
        kind: "trial",
        payload: JSON.stringify(result2),
      })
    );
  }
};


// API

export const soundsN1 = (state, shared) => {
  const { jsPsych, wording, stimuli } = state;
  const { showProgress } = shared;

  const blockStop = genBlockStop(state, shared);
  return {
    // default value preventing text flickering when loading assets
    prompt: `<p>[${wording.space}] <span style='font-weight:bold'> ${wording.playSounds}</span></p>
    <p>${wording.question}</p>
    <div class='sound-choice'>
      <div>[${wording.keyAlt1}] ${wording.labelAlt1}</div>
      <div>${wording.labelAlt2} [${wording.keyAlt2}]</div>
    </div>`,
    // actual timeline
    timeline: [
      {
        type: jsPsychPreload,
        audio: () => {
          return [
            `${ASSET_PREFIX}${jsPsych.timelineVariable("asset")}`,
          ];
        },
        show_progress_bar: false,
        post_trial_gap: 200,
        on_start: showProgress,
      },
      {
        type: jsPsychHtmlKeyboardResponse,
        stimulus: "",
        choices: " ", // press space to next step
        prompt: `<p><span style='font-weight:bold'>[${wording.space}]</span> ${wording.playSounds}</p>
        <p>${wording.question}</p>
        <div class='sound-choice'>
          <div>[${wording.keyAlt1}] ${wording.labelAlt1}</div>
          <div>${wording.labelAlt2} [${wording.keyAlt2}]</div>
        </div>`,
      },
      {
        type: jsPsychAudioKeyboardResponse,
        stimulus: () => `${ASSET_PREFIX}${jsPsych.timelineVariable("asset")}`,
        choices: "NO_KEYS",
        trial_ends_after_audio: true,
        response_allowed_while_playing: false,
      },
      {
        type: jsPsychHtmlKeyboardResponse,
        stimulus: "",
        choices: [wording.keyAlt1, wording.keyAlt2],
        prompt: `<p>[${wording.space}] ${wording.playSounds}</p>
        <p>${wording.question}</p>
        <div class='sound-choice'>
          <div><span class='strong'>[${wording.keyAlt1}]</span> ${wording.labelAlt1}</div>
          <div>${wording.labelAlt2} <span class='strong'>[${wording.keyAlt2}]</span></div>
        </div>`,
        data: {
          answered: true,
        },
        on_finish: trialSubmitN1(state, shared),
      },
      blockStop,
    ],
    timeline_variables: stimuli,
  };
}

export const soundsN2 = (state, shared) => {
  const { jsPsych, wording, stimuli } = state;
  const { showProgress } = shared;

  const blockStop = genBlockStop(state, shared);
  return {
    // default value preventing text flickering when loading assets
    prompt: `<p>[${wording.space}] <span style='font-weight:bold'> ${wording.playSounds}</span></p>
    <p>${wording.question}</p>
    <div class='sound-choice'>
      <div>[${wording.keyAlt1}] ${wording.labelAlt1}</div>
      <div>${wording.labelAlt2} [${wording.keyAlt2}]</div>
    </div>`,
    // actual timeline
    timeline: [
      {
        type: jsPsychPreload,
        audio: () => {
          return [
            `${ASSET_PREFIX}${jsPsych.timelineVariable("asset1")}`,
            `${ASSET_PREFIX}${jsPsych.timelineVariable("asset2")}`,
          ];
        },
        show_progress_bar: false,
        post_trial_gap: 200,
        on_start: showProgress,
      },
      {
        type: jsPsychHtmlKeyboardResponse,
        stimulus: "",
        choices: " ", // press space to next step
        prompt: `<p><span style='font-weight:bold'>[${wording.space}]</span> ${wording.playSounds}</p>
        <p>${wording.question}</p>
        <div class='sound-choice'>
          <div>[${wording.keyAlt1}] ${wording.labelAlt1}</div>
          <div>${wording.labelAlt2} [${wording.keyAlt2}]</div>
        </div>`,
      },
      {
        type: jsPsychAudioKeyboardResponse,
        stimulus: () => `${ASSET_PREFIX}${jsPsych.timelineVariable("asset1")}`,
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
        stimulus: () => `${ASSET_PREFIX}${jsPsych.timelineVariable("asset2")}`,
        choices: "NO_KEYS",
        trial_ends_after_audio: true,
        response_allowed_while_playing: false,
      },
      {
        type: jsPsychHtmlKeyboardResponse,
        stimulus: "",
        choices: [wording.keyAlt1, wording.keyAlt2],
        prompt: `<p>[${wording.space}] ${wording.playSounds}</p>
        <p>${wording.question}</p>
        <div class='sound-choice'>
          <div><span class='strong'>[${wording.keyAlt1}]</span> ${wording.labelAlt1}</div>
          <div>${wording.labelAlt2} <span class='strong'>[${wording.keyAlt2}]</span></div>
        </div>`,
        data: {
          answered: true,
        },
        on_finish: trialSubmitN2(state, shared),
      },
      blockStop,
    ],
    timeline_variables: stimuli,
  };
};

export const imagesN1 = (state, shared) => {
  const { jsPsych, settings, wording, stimuli } = state;
  const { showProgress } = shared;

  console.log(settings);

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
            `${ASSET_PREFIX}${jsPsych.timelineVariable("asset")}`,
          ];
        },
        on_start: showProgress,
        show_progress_bar: false,
        post_trial_gap: 200,
      },
      {
        type: jsPsychHtmlKeyboardResponse,
        stimulus: "",
        choices: [wording.keyAlt1, wording.keyAlt2],
        prompt: () => {
          const imgWidth =
            settings.forceWidth.length == 0 ? "auto" : settings.forceWidth;
          return `<div>
            <img style="width:${imgWidth};" src="${ASSET_PREFIX}${jsPsych.timelineVariable(
            "asset")}">
          </div>
          <p>${wording.question}</p>
          <div class='image-choice'>
            <div>
              <span class='strong'>[${wording.keyAlt1}]</span> ${wording.labelAlt1}
            </div>
            <div>
              ${wording.labelAlt2} <span class='strong'>[${
              wording.keyAlt2}]</span>
            </div>
          </div>`;
        },
        data: {
          answered: true,
        },
        on_finish: trialSubmitN1(state, shared),
      },
      blockStop,
    ],
    timeline_variables: stimuli,
  };
}

export const imagesN2 = (state, shared) => {
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
            `${ASSET_PREFIX}${jsPsych.timelineVariable("asset1")}`,
            `${ASSET_PREFIX}${jsPsych.timelineVariable("asset2")}`,
          ];
        },
        on_start: showProgress,
        show_progress_bar: false,
        post_trial_gap: 200,
      },
      {
        type: jsPsychHtmlKeyboardResponse,
        stimulus: "",
        choices: [wording.keyAlt1, wording.keyAlt2],
        prompt: () => {
          const imgWidth =
            settings.forceWidth.length == 0 ? "auto" : settings.forceWidth;
          return `<p>${wording.question}</p>
          <div class='image-choice'>
            <div>
              <img style="width:${imgWidth};" src="${ASSET_PREFIX}${jsPsych.timelineVariable(
              "asset1"
            )}">
              <div><span class='strong'>[${wording.keyAlt1}]</span> ${
              wording.labelAlt1
            }</div>
            </div>
            <div>
              <img style="width:${imgWidth};" src="${ASSET_PREFIX}${jsPsych.timelineVariable(
              "asset2"
            )}">
              <div>${wording.labelAlt2} <span class='strong'>[${
              wording.keyAlt2
            }]</span></div>
            </div>
          </div>`;
        },
        data: {
          answered: true,
        },
        on_finish: trialSubmitN2(state, shared),
      },
      blockStop,
    ],
    timeline_variables: stimuli,
  };
};

