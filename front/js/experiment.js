const pairs = (arr) =>
    Array.from(
        new Array(Math.ceil(arr.length / 2)),
        (_, i) => {
            const pair = arr.slice(i * 2, i * 2 + 2);
            return { s1: pair[0], s2: pair[1]};
        }
    );

const SELECT_1_KEY = "f";
const SELECT_2_KEY = "j";

const inBlock = (done, trialsPerBlock) => Math.floor(done / trialsPerBlock);

// build and run experiment
export default (state, ws) => {
    const start = (new Date()).toISOString();
    const { settings, wording, participant } = state;
    const todo = participant.todo.split(",");
    const stimuli = pairs(todo, 2);

    const totalLength = settings.blockCount * settings.trialsPerBlock;
    const remainingLength = todo.length / 2; // 2 sounds per trial
    const previouslyDoneLength = totalLength - remainingLength; // not 0 if user reconnects (page refresh for instance) 

    const timeline = [];
    const position = {
        trial: previouslyDoneLength,
        block: inBlock(previouslyDoneLength, settings.trialsPerBlock)
    };

    // experiment has already be fully run by this participant
    if(remainingLength === 0) {
        timeline.push({
            type: jsPsychHtmlKeyboardResponse,
            stimulus: `<h3>${wording.closed}</h3>`,
            choices: "NO_KEYS"
        });
    } else {
        // form to collect participant info
        if(participant.age.length === 0 || participant.sex.length === 0) {
            timeline.push({
                type: jsPsychSurveyHtmlForm,
                preamble: `<p>${wording.collect}</p>`,
                html: `<p>
                    <fieldset>
                        <label>${wording.collectAge}</label>
                        <input id="age" name="age" type="text" minlength="2" maxlength="3" required />
                    </fieldset>
                    <fieldset>
                        <label>${wording.collectSex}</label>
                        <input name="sex" type="text" minlength="1" maxlength="1" required />
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
                }
            });
        }
        // does the participant start for the first time?
        if(previouslyDoneLength == 0) {
            timeline.push({
                type: jsPsychHtmlKeyboardResponse,
                stimulus: `<p>${wording.introduction}</p>`,
                prompt: `<p><span class='strong'>[${wording.space}]</span> ${wording.next}</p>`,
                choices: " "
            });
        } else {
            timeline.push({
                type: jsPsychHtmlKeyboardResponse,
                stimulus: `<p>${wording.resume}</p>`,
                prompt: `<p><span class='strong'>[${wording.space}]</span> ${wording.next}</p>`,
                choices: " "
            });
        }

        const blockStop = {
            timeline: [
                {
                    type: jsPsychHtmlKeyboardResponse,
                    stimulus: `<p>${wording.pause}</p>`,
                    prompt: "",
                    choices: "NO_KEYS",
                    trial_duration: 4000
                },
                {
                    type: jsPsychHtmlKeyboardResponse,
                    stimulus: `<p>${wording.pauseOver}</p>`,
                    prompt: `<p><span class='strong'>[${wording.space}]</span> ${wording.next}</p>`,
                    choices: " ",
                }
            ],
            conditional_function: function(){
                const done = jsPsych.data.get().filter({answered: true}).count() + previouslyDoneLength;
                const blockEnd = (done % settings.trialsPerBlock) === 0;
                // display if end of block and if not last block
                return blockEnd && done !== totalLength;
            }
        }

        timeline.push({
            prompt: `<p>[${wording.space}] <span style='font-weight:bold'> ${wording.stimuli}</span></p>
            <p>${wording.question}</p>
            <div class='choice'>
                <div>[f] ${wording.sound1}</div>
                <div>${wording.sound2} [j]</div>
            </div>`,
            timeline: [
                {
                    type: jsPsychPreload,
                    audio: () => {
                        return [`sounds/${jsPsych.timelineVariable("s1")}`, `sounds/${jsPsych.timelineVariable("s2")}`]
                    },
                    show_progress_bar: false,
                    post_trial_gap: 200
                },
                {
                    type: jsPsychHtmlKeyboardResponse,
                    stimulus: "",
                    choices: " ",
                    prompt: `<p><span style='font-weight:bold'>[${wording.space}]</span> ${wording.stimuli}</p>
                    <p>${wording.question}</p>
                    <div class='choice'>
                        <div>[f] ${wording.sound1}</div>
                        <div>${wording.sound2} [j]</div>
                    </div>`,
                },
                {
                    type: jsPsychAudioKeyboardResponse,
                    stimulus: () => `sounds/${jsPsych.timelineVariable("s1")}`,
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
                    stimulus: () => `sounds/${jsPsych.timelineVariable("s2")}`,
                    choices: "NO_KEYS",
                    trial_ends_after_audio: true,
                    response_allowed_while_playing: false,
                },
                {
                    type: jsPsychHtmlKeyboardResponse,
                    stimulus: "",
                    choices: [SELECT_1_KEY, SELECT_2_KEY],
                    prompt: `<p>[${wording.space}] ${wording.stimuli}</p>
                    <p>${wording.question}</p>
                    <div class='choice'>
                        <div><span class='strong'>[f]</span> ${wording.sound1}</div>
                        <div>${wording.sound2} <span class='strong'>[j]</span></div>
                    </div>`,
                    data: {
                        answered: true
                    },
                    on_finish: (data) => {
                        const result1 = {
                            trial: position.trial.toString(),                         
                            block: position.block.toString(),                         
                            stimulus: jsPsych.timelineVariable("s1"),
                            order: "0",
                            response: data.response === SELECT_1_KEY ? "True" : "False",
                            rt: data.rt.toString(),
                            date: start
                        };
                        const result2 = {
                            trial: position.trial.toString(),                         
                            block: position.block.toString(),   
                            stimulus: jsPsych.timelineVariable("s2"),
                            order: "1",
                            response: data.response === SELECT_2_KEY ? "True" : "False",
                            rt: data.rt.toString(),
                            date: start
                        }
                        position.trial++;
                        position.block = inBlock(position.trial, settings.trialsPerBlock);
                        ws.send(
                            JSON.stringify({
                                kind: "trial",
                                payload: JSON.stringify({ result1, result2 }),
                            })
                        );
                    }
                },
                blockStop,
            ],
            timeline_variables: stimuli
        });

        timeline.push({
            type: jsPsychHtmlKeyboardResponse,
            stimulus: `<h3>${wording.end}</h3>`,
            prompt: `<p>${wording.thanks}</p>`,
            choices: "NO_KEYS",
            on_start: function() {
                console.log("The experiment is over");
            }
        });
    }
    
    const jsPsych = initJsPsych({
        on_finish: function() {
            jsPsych.data.displayData();
        }
    });
    jsPsych.run(timeline);
}