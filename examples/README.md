## Creation

To create a new experiment, follow these instructions:

1. Choose an experience ID using this character set: `a-zA-Z0-9-_`. From now on, let say you chose `experiment_3`.

If your revcor app is reachable at `https://example.com`, the root path of the experiment is `https://example.com/xp/experiment_3/`

2. Create the following folders and files in `data`:

* (folder) `experiment_3/` at the same location than this README.md
* (folder) `experiment_3/config/` to hold configuration files
* (folder) `experiment_3/assets/` to hold wav (or jpg/png for images) assets used during the experiment, plus *.txt files describing how the sounds/images have been generated (see more below)
* (folder) `experiment_3/results/` let empty at start, trials will populate it with data
* `experiment_3/config/settings.json` JSON file (see more below)
* `experiment_3/config/wording.new.json` and `experiment_3/config/wording.run.json` JSON file (see more below)
* `experiment_3/config/participants.txt` text file containing one participant ID per line (valid IDs are also made of `a-zA-Z0-9-_`)

For convenience you may in fact declare several `participants` files: any text file containing the string `participants` (in the `experiment_3/config/` folder) is considered a valid participants dictionnary.

If you happen want an experiment to be available in different languages, we currently recommend to create separate experiment folders (for instance `experiment_3_fr` and `experiment_3_en`) that will contain their own configuration and wording files.

3. The `settings.json` file configures the experiment and needs the following properties:

```json
{
  "adminPassword": "to_change",
  "allowCreate": true,
  "createPassword": "consent",
  "trialsPerBlock": 3,
  "blocksPerXp": 6,
  "addRepeatBlock": true,
  "assetsPerParticipant": true,
  "kind": "image",
  "fileExtension": "jpg",
  "forceWidth": "350px",
  "showProgress": true,
  "collectInfo": [
    {
      "key": "age",
      "label": "âge",
      "kind": "text",
      "pattern": "[0-9]*"
    },
    {
      "key": "sex",
      "label": "sexe",
      "kind": "text"
    }
  ]
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
* `kind` (optional): either `sound` (default) or `image`
* `fileExtension` (optional): if you want to set the assets extension (by default `wav` for sounds, `png` for images). Supported values are: `wav`, `png` and `jpg`
* `forceWidth` (optional and used only for `image` kind): force the width of displayed images using the provided value as a CSS property. If not set, images are displayed in their original size
* `showProgress` (optional and false by default): show the progress (in the form `trials done/total count`) at the right bottom corner
* `collectInfo` (optional, no HTML form if empty): add info definition to be collected before the trials in a HTML form, with the following subproperties:
  - `key` the column header corresponding to that piece of data in CSV results
  - `label` the displayed label in the HTML form
  - `kind` the HTML input type attribute
  - `pattern` the HTML input pattern attribute

4. The `wording.run.json` file configures messages displayed on screens and needs the following properties:

```json
{
    "title": "The smile of sounds",
    "collect": "Please fill in the following information:",
    "collectButton": "Continue",
    "introduction": "In this experiment, you will hear examples of pronunciations of the sound /a/, and we ask you to judge which one you think was pronounced with the most smile.",
    "pause": "Let's pause for a few seconds",
    "pauseOver": "The pause is over, you can resume the experiment",
    "resume": "Resuming",
    "end": "End of the experiment",
    "thanks": "Thanks for your participation",
    "closed": "Experiment already done",
    "playSounds": "listening to voices 1 & 2",
    "question": "Which pronunciation is the most smiling?",
    "next": "next",
    "label1": "voice 1",
    "label2": "voice 2",
    "space": "space",
    "choice1": "f",
    "choice2": "j"
}
```

**Caution**: the properties `"choice1` and `choice2` have to map actual keyboard letter keys, they are used as is to collect the participant decision.

An additional `wording.new.json` file has to be provided for the participant creation page (only available if `allowCreate` is true), please check `examples/sound_xp/config/wording.new.json`.

5. Put assets to be tested in `experiment_3/assets`. For each *asset* file (sound or image) there MUST be a *definition* file respecting the following conditions:
* for a given sound or image, the asset and definition files have an identical name but a different extension. For instance the asset `gomot_a.0001.eq.wav` is paired with its definition `gomot_a.0001.eq.txt` (definition extension are always `txt`)
* definition files are CSV formatted, with comma-separated headers on their first line
* in the same experiment, all definition files share the same headers (implying that all `experiment_3/assets/*.txt` start with the same first line, defining headers)

Here is a sample of a possible sound definition file:

```csv
filter_freq,filter_gain
0.00000000,-4.65473028
104.58767290,3.80355849
224.80189054,4.16050504
...
```

Note: headers are supposed to be identical for all sounds/images within a given experiment.

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

## Reset

If you want to reset an experience collected data, you should be aware that data relevant to participants is stored at two different places:

- `data/experiment_3/results/`: the CSV format output that is used for further analysis
- `data/experiment_3/state/`: which is an internal state folder used to store participant data during the experiment, particularly helpful if the participant takes a break or reloads the page

That's why you should delete files/folders at both places if you want to get rid of data related to given participants.