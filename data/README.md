To create a new experiment, follow these instructions:

1. Choose an experience ID using this character set: `a-zA-Z0-9-_`. From now on, let say you choose `experiment_3`.

If your revcor app is reachable at `https://example.com`, the root path of the experiment is `https://example.com/xp/experiment_3/`

2. Create the following folders and files in `data`:

* (folder) `experiment_3/` at the same location than this README.md
* (folder) `experiment_3/config/` to hold configuration files
* (folder) `experiment_3/sounds/` to hold *.wav files used during the experiment, and *.txt files describing how the sounds have been generated
* (folder) `experiment_3/results/` let empty at start, trials will populate it with data
* `experiment_3/config/settings.json` JSON file (see more below)
* `experiment_3/config/wording.json` JSON file (see more below)
* `experiment_3/config/participants.txt` text file containing one participant ID per line (valid IDs are also made of `a-zA-Z0-9-_)

For convenience you may in fact declare several `participants` files: any text file containing the string `participants` (in the `experiment_3/config/` folder) is considered a valid participants dictionnary.

3. The `settings.json` file configures the experiment and needs the following properties:

```json
{
    "adminPassword": "unsecure",
    "allowCreate": true,
    "createPassword": "consent",
    "trialsPerBlock": 3,
    "blocksPerXp": 6,
    "addRepeatBlock": true
}
```

Let's review the meaning of each property:

* `adminPassword`: if you want to access CSV results go to `https://example.com/xp/experiment_3/results` and enter admin as the login and the value of `adminPassword` as the password
* `allowCreate`:
    * if `false`, only declared participants are allowed to access this experiment. Check step 6 to see how to declare them
    * if `true`, anyone going to `https://example.com/xp/experiment_3/new` will be able to access the experiment, if they know the `createPassword` 
* `createPassword`: password that protects the creation of participants
* `trialsPerBlock`: trials are grouped by blocks (with a pause between blocks), define how many trials make a block
* `blocksPerXp`: how many blocks are in the experiment
* `addRepeatBlock`: repeat last block (for consistency measurement). In the example above, the experiment is made of 7 (6 + 1) blocks.

4. The `wording.run.json` file configures messages displayed on screens and needs the following properties:

```json
{
    "title": "The smile of sounds",
    "collect": "Please fill in the following information:",
    "collectAge": "age",
    "collectSex": "sex",
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
    "space": "space",
    "choice1": "f",
    "choice2": "j"
}
```

**Caution**: the properties `"choice1` and `choice2` have to map actual keyboard letter keys, they are used as is to collect the participant decision.

An additional `wording.new.json` file has to be provided for the participant creation page (only available if `allowCreate` is true), please check `data/example/confid.wording.new.json`.

5. Put wav files to be tested in `experiment_3/sounds` and add a text file defining how the sound has been generated in a CSV format. The file names need to be identical, with only the file extension changing. For instance the `gomot_a.0001.eq.wav` sound file has to be paired with `gomot_a.0001.eq.txt`, whose contents looks like:

```csv
filter_freq,filter_gain
0.00000000,-4.65473028
104.58767290,3.80355849
224.80189054,4.16050504
...
```

Note: headers are supposed to be identical for all sounds within a given experiment.

This CSV definition is used when appending to the CSV result file. Here is an extract of a result file corresponding to the definition above (check the `filter_freq` and `filter_gain` from the definition above, and the added `param_index`):

```csv
subj,trial,block,sex,age,date,stim,stim_order,param_index,filter_freq,filter_gain,response,rt
100,0,0,f,33,2021-11-22T21:01:26.374Z,gomot_a.0291.eq.wav,0,0,0.00000000,-4.65473028,1,330
100,0,0,f,33,2021-11-22T21:01:26.374Z,gomot_a.0291.eq.wav,0,1,104.58767290,3.80355849,1,330
100,0,0,f,33,2021-11-22T21:01:26.374Z,gomot_a.0291.eq.wav,0,2,224.80189054,4.16050504,1,330
...
```

In this example the 3 lines refer to the same sound and the same trial result, the unfolded/multiline notation being intended to help with further analysis.

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


