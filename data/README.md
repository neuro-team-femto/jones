To create a new experiment, follow these instructions:

1. Choose an experience ID using this character set: [a-zA-Z0-9-_]

2. Let say you chose `experiment_3` in the previous step, create the following folders and files:

* an `experiment_3` folder (at the same location than this README.md)
* an `experiment_3/config` folder
* an `experiment_3/sounds` folder (to hold *.wav files used during the experiment)
* an `experiment_3/results` folder (let empty at start, trials will populate it with data)
* an `experiment_3/config/settings.json` JSON file (see more below)
* an `experiment_3/config/wording.json` JSON file (see more below)
* an `experiment_3/config/participants.txt` file containing one participant ID per line (valid IDs are also made of [a-zA-Z0-9-_])

For convenience you may in fact declare several `participants` files: any text file containing the string `participants` (in the `experiment_3/config/` folder) is considered a valid participants dictionnary.

3. The `settings.json` file configures the experiment and needs the following properties:

```json
{
    "blockCount": 4,
    "trialsPerBlock": 10,
}
```

4. The `wording.json` file configures messages displayed on screens and needs the following properties:

```json
{
    "title": "The smile of sounds",
    "collect": "Please fill in the following information:",
    "collectAge": "age",
    "collectSex": "sex (f/m)",
    "collectButton": "Continue",
    "introduction": "In this experiment, you will hear examples of pronunciations of the sound /a/, and we ask you to judge which one you think was pronounced with the most smile.",
    "pause": "Let's pause for a few seconds",
    "pauseOver": "The pause is over, you can resume the experiment",
    "resume": "Resuming",
    "end": "End of the experiment",
    "thanks": "Thanks for your participation",
    "closed": "Experiment already done",
    "stimuli": "listening to voices 1 & 2",
    "question": "Which pronunciation is the most smiling?",
    "next": "next",
    "sound1": "voice 1",
    "sound2": "voice 2",
    "space": "space"
}
```

5. Put wav files to be tested in `experiment_3/sounds` (with a `.wav` file extension)

6. Here is an example `participants` file, defining 4 participant IDs:
```text
b0c410eacc023237ca8d9cfea109ab70
d465f071d45d8a216b42d6411e865bcf
f003a58ffc73c3bd44f2c44662c98def
1de290f8d4e545f768851e4039770709
```

With this `participants` file, and if the webapp is hosted at `https://example.com/` you may share the following links to participants (don't forget the `xp` path prefix)

https://example.com/xp/experiment_3/b0c410eacc023237ca8d9cfea109ab70
https://example.com/xp/experiment_3/d465f071d45d8a216b42d6411e865bcf
https://example.com/xp/experiment_3/f003a58ffc73c3bd44f2c44662c98def
https://example.com/xp/experiment_3/1de290f8d4e545f768851e4039770709


