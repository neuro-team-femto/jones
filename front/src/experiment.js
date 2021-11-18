const pairs = (arr) =>
    Array.from(
        new Array(Math.ceil(arr.length / 2)),
        (_, i) => {
            const pair = arr.slice(i * 2, i * 2 + 2);
            return { s1: pair[0], s2: pair[1]};
        }
    );

// build and run experiment
export default (state, ws) => {
    const { experiment: xp, participant: { sounds: soundsStr } } = state;
    const sounds = soundsStr.split(',');

    const todoLength = sounds.length / 2;
    const totalLength = xp.blockCount * xp.trialsPerBlock;
    const previouslyDoneLength = totalLength - todoLength;


    const stimuli = pairs(sounds, 2);

    const jsPsych = initJsPsych({
        on_finish: function() {
            jsPsych.data.displayData();
        }
    });

    const timeline = [];

    // experiment has already be fully run by this participant
    if(todoLength === 0) {
        timeline.push({
            type: jsPsychHtmlKeyboardResponse,
            stimulus: `<h3>Expérience déjà effectuée</h3>`,
            choices: "NO_KEYS"
        });
    } else {
        // does the participant start for the first time?
        if(previouslyDoneLength == 0) {
            // intro page
            timeline.push({
                type: jsPsychHtmlKeyboardResponse,
                stimulus: `<p>${xp.introduction}</p>`,
                prompt: "<p><span class='strong'>[espace]</span> pour continuer</p>",
                choices: " "
            });
        } else {
            timeline.push({
                type: jsPsychHtmlKeyboardResponse,
                stimulus: `<p>Reprise de l'expérience</p>`,
                prompt: "<p><span class='strong'>[espace]</span> pour continuer</p>",
                choices: " "
            });
        }

        const blockStop = {
            timeline: [
                {
                    type: jsPsychHtmlKeyboardResponse,
                    stimulus: "<p>Nous vous proposons une pause de quelques secondes.</p>",
                    prompt: "",
                    choices: "NO_KEYS",
                    trial_duration: 6000
                },
                {
                    type: jsPsychHtmlKeyboardResponse,
                    stimulus: "<p>La pause est terminée. Vous pouvez reprendre l'expérience.</p>",
                    prompt: "<p><span class='strong'>[espace]</span> pour continuer</p>",
                    choices: " ",
                }
            ],
            conditional_function: function(){
                const done = jsPsych.data.get().filter({answered: true}).count() + previouslyDoneLength;
                const blockEnd = (done % xp.trialsPerBlock) === 0;
                return blockEnd;
            }
        }

        timeline.push({
            prompt: `<p>[espace] <span style='font-weight:bold'>écoute ${xp.trialSoundLabel} 1 & 2</span></p>
            <p>${xp.trialQuestion}</p>
            <div class='choice'>
                <div>[f] ${xp.trialSoundLabel} 1</div>
                <div>${xp.trialSoundLabel} 2 [j]</div>
            </div>`,
            timeline: [
                {
                    type: jsPsychPreload,
                    audio: () => {
                        [`sounds/${jsPsych.timelineVariable('s1')}`, `sounds/${jsPsych.timelineVariable('s2')}`]
                    },
                    show_detailed_errors: true
                },
                {
                    type: jsPsychHtmlKeyboardResponse,
                    stimulus: '',
                    choices: " ",
                    prompt: `<p><span style='font-weight:bold'>[espace]</span> écoute ${xp.trialSoundLabel} 1 & 2</p>
                    <p>${xp.trialQuestion}</p>
                    <div class='choice'>
                        <div>[f] ${xp.trialSoundLabel} 1</div>
                        <div>${xp.trialSoundLabel} 2 [j]</div>
                    </div>`,
                },
                {
                    type: jsPsychAudioKeyboardResponse,
                    stimulus: () => `sounds/${jsPsych.timelineVariable('s1')}`,
                    choices: "NO_KEYS",
                    trial_ends_after_audio: true,
                    response_allowed_while_playing: false,
                },
                {
                    type: jsPsychHtmlKeyboardResponse,
                    stimulus: '',
                    choices: "NO_KEYS",
                    trial_duration: 500,
                },
                {
                    type: jsPsychAudioKeyboardResponse,
                    stimulus: () => `sounds/${jsPsych.timelineVariable('s2')}`,
                    choices: "NO_KEYS",
                    trial_ends_after_audio: true,
                    response_allowed_while_playing: false,
                },
                {
                    type: jsPsychHtmlKeyboardResponse,
                    stimulus: '',
                    choices: ["f", "j"],
                    prompt: `<p>[espace] écoute ${xp.trialSoundLabel} 1 & 2</p>
                    <p>${xp.trialQuestion}</p>
                    <div class='choice'>
                        <div><span class='strong'>[f]</span> ${xp.trialSoundLabel} 1</div>
                        <div>${xp.trialSoundLabel} 2 <span class='strong'>[j]</span></div>
                    </div>`,
                    data: {
                        answered: true
                    },
                    on_finish: (data) => {
                        const chosen = data.response === 'f' ? jsPsych.timelineVariable('s1') : jsPsych.timelineVariable('s2');
                        const dismissed = data.response === 'f' ? jsPsych.timelineVariable('s2') : jsPsych.timelineVariable('s1');
                        ws.send(
                            JSON.stringify({
                                kind: "result",
                                payload: JSON.stringify({ chosen, dismissed }),
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
            stimulus: `<h3>Fin de l'expérience</h3>`,
            prompt: "<p>Merci pour votre participation/p>",
            choices: "NO_KEYS",
            on_start: function() {
                console.log('The experiment is over');
            }
        });
    }
    
    jsPsych.run(timeline);
}