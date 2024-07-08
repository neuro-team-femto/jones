<!DOCTYPE html>
<html>
  <head>
    <title></title>
    <link rel="shortcut icon" href="data:image/x-icon;," type="image/x-icon">
    <script src="{{.webPrefix}}/scripts/jspsych/jspsych.js"></script>
    <script src="{{.webPrefix}}/scripts/jspsych/plugin-preload.js"></script>
    <script src="{{.webPrefix}}/scripts/jspsych/plugin-survey-html-form.js"></script>
    <script src="{{.webPrefix}}/scripts/jspsych/plugin-html-keyboard-response.js"></script>
    <script src="{{.webPrefix}}/scripts/jspsych/plugin-audio-keyboard-response.js"></script>
    <script src="{{.webPrefix}}/scripts/jspsych/plugin-video-keyboard-response.js"></script>
    <script src="{{.webPrefix}}/scripts/main.js"></script>
    <link rel="stylesheet" href="{{.webPrefix}}/styles/jspsych.css">
    <link rel="stylesheet" href="{{.webPrefix}}/styles/main.css">
  </head>
  <body>
    <div id="jspsych-root"></div>
    <div id="progress" style="display: none;"></div>
    <script>
      var stateIds = { experimentId: {{.experimentId}}, participantId: {{.participantId}}};
    </script>
  </body>
</html>