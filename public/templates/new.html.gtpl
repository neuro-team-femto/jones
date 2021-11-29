<!DOCTYPE html>
<html>
  <head>
    <title></title>
    <link rel="shortcut icon" href="data:image/x-icon;," type="image/x-icon">
    <link rel="stylesheet" href="{{.webPrefix}}/styles/main.css">
    <link rel="stylesheet" href="{{.webPrefix}}/styles/jspsych.css">
  </head>
  <body class="new-body jspsych-display-element" style="margin: 0px; height: 100%; width: 100%;">
    <div class="jspsych-content-wrapper">
      <div id="jspsych-content" class="jspsych-content">
        <div id="jspsych-survey-html-form-preamble" class="jspsych-survey-html-form-preamble">
          {{if .error}} <p class="strong">{{.error}}</p> {{end}}
          <p>{{.wording.instructions}}</p>
        </div>
        <form id="jspsych-survey-html-form" action="{{.webPrefix}}/xp/{{.experimentId}}/new" autocomplete="off" method="POST">
          <fieldset>
            <label>{{.wording.idLabel}}</label>
            <input name="id" type="text" maxlength="30" required value="{{.id}}">
          </fieldset>
          <fieldset>
            <label>{{.wording.passwordLabel}}</label>
            <input name="password" type="text" required>
          </fieldset>
          <input type="submit" id="jspsych-survey-html-form-next" class="jspsych-btn jspsych-survey-html-form" value="{{.wording.submitLabel}}">
        </form>
      </div>
    </div>
</body>
</html>