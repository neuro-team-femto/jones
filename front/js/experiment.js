const pairs = (arr) =>
  Array.from(new Array(Math.ceil(arr.length / 2)), (_, i) => {
    const pair = arr.slice(i * 2, i * 2 + 2);
    return { s1: pair[0], s2: pair[1] };
  });

const inBlock = (done, trialsPerBlock) => Math.floor(done / trialsPerBlock);

const ASSET_PREFIX = "../assets/";

const soundTimeline = (ws, jsPsych, start, settings, wording, stimuli, position, blockStop) => ({
  prompt: `<p>[${wording.space}] <span style='font-weight:bold'> ${wording.playSounds}</span></p>
  <p>${wording.question}</p>
  <div class='choice'>
    <div>[${wording.choice1}] ${wording.label1}</div>
    <div>${wording.label2} [${wording.choice2}]</div>
  </div>`,
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
    },
    {
      type: jsPsychHtmlKeyboardResponse,
      stimulus: "",
      choices: " ",
      prompt: `<p><span style='font-weight:bold'>[${wording.space}]</span> ${wording.playSounds}</p>
      <p>${wording.question}</p>
      <div class='choice'>
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
      <div class='choice'>
        <div><span class='strong'>[${wording.choice1}]</span> ${wording.label1}</div>
        <div>${wording.label2} <span class='strong'>[${wording.choice2}]</span></div>
      </div>`,
      data: {
        answered: true,
      },
      on_finish: (data) => {
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
        ws.send(
          JSON.stringify({
            kind: "trial",
            payload: JSON.stringify({ result1, result2 }),
          })
        );
      },
    },
    blockStop,
  ],
  timeline_variables: stimuli,
});

const imageTimeline = (ws, jsPsych, start, settings, wording, stimuli, position, blockStop) => ({
  prompt: "",
  timeline: [
    {
      type: jsPsychPreload,
      images: () => {
        return [
          `${ASSET_PREFIX}${jsPsych.timelineVariable("s1")}`,
          `${ASSET_PREFIX}${jsPsych.timelineVariable("s2")}`,
        ];
      },
      show_progress_bar: false,
      post_trial_gap: 200,
    },
    {
      type: jsPsychHtmlKeyboardResponse,
      stimulus: "",
      choices: [wording.choice1, wording.choice2],
      prompt: () => {
        return `<p>${wording.question}</p>
        <div class='choice'>
          <img src="${ASSET_PREFIX}${jsPsych.timelineVariable("s1")}">
          <img src="${ASSET_PREFIX}${jsPsych.timelineVariable("s2")}">
          <div><span class='strong'>[${wording.choice1}]</span> ${wording.label1}</div>
          <div>${wording.label2} <span class='strong'>[${wording.choice2}]</span></div>
        </div>`;
      },
      data: {
        answered: true,
      },
      on_finish: (data) => {
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
        ws.send(
          JSON.stringify({
            kind: "trial",
            payload: JSON.stringify({ result1, result2 }),
          })
        );
      },
    },
    blockStop,
  ],
  timeline_variables: stimuli,
});

// build and run experiment
export default (state, ws) => {
  const start = new Date().toISOString();
  const { settings, wording, participant } = state;
  const todo = participant.todo.split(",");
  const stimuli = pairs(todo, 2);
  const blocks = settings.addRepeatBlock
    ? settings.blocksPerXp + 1
    : settings.blocksPerXp;

  const totalLength = blocks * settings.trialsPerBlock;
  const remainingLength = todo.length / 2; // 2 choices per trial
  const previouslyDoneLength = totalLength - remainingLength; // not 0 if user reconnects (page refresh for instance)

  const timeline = [];
  const position = {
    trial: previouslyDoneLength,
    block: inBlock(previouslyDoneLength, settings.trialsPerBlock),
  };

  const jsPsych = initJsPsych({
    on_finish: function () {
      jsPsych.data.displayData();
    },
  });

  // experiment has already be fully run by this participant
  if (remainingLength === 0) {
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
    // does the participant start for the first time?
    const firstStimulus =
      previouslyDoneLength == 0 ? wording.introduction : wording.resume;
    timeline.push({
      type: jsPsychHtmlKeyboardResponse,
      stimulus: `<p>${firstStimulus}</p>`,
      prompt: `<p><span class='strong'>[${wording.space}]</span> ${wording.next}</p>`,
      choices: " ",
    });

    const blockStop = {
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
      conditional_function: function () {
        const done =
          jsPsych.data.get().filter({ answered: true }).count() +
          previouslyDoneLength;
        const blockEnd = done % settings.trialsPerBlock === 0;
        // display if end of block and if not last block
        return blockEnd && done !== totalLength;
      },
    };

    const assetTimeline = settings.kind === "sound" ? soundTimeline : imageTimeline;
    timeline.push(assetTimeline(ws, jsPsych, start, settings, wording, stimuli, position, blockStop));

    timeline.push({
      type: jsPsychHtmlKeyboardResponse,
      stimulus: `<h3>${wording.end}</h3>`,
      prompt: `<p>${wording.thanks}</p>`,
      choices: "NO_KEYS",
      on_start: function () {
        console.log("The experiment is over");
      },
    });
  }

  jsPsych.run(timeline);
};
