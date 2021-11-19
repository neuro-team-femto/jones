var jsPsychExtensionMouseTracking = (function () {
  'use strict';

  class MouseTrackingExtension {
      constructor(jsPsych) {
          this.jsPsych = jsPsych;
          this.initialize = ({ minimum_sample_time = 0 }) => {
              this.domObserver = new MutationObserver(this.mutationObserverCallback);
              this.minimumSampleTime = minimum_sample_time;
              return new Promise((resolve, reject) => {
                  resolve();
              });
          };
          this.on_start = (params) => {
              this.currentTrialData = [];
              this.currentTrialTargets = new Map();
              this.currentTrialSelectors = typeof params !== "undefined" ? params.targets : [];
              this.lastSampleTime = null;
              this.domObserver.observe(this.jsPsych.getDisplayElement(), { childList: true });
          };
          this.on_load = () => {
              // set current trial start time
              this.currentTrialStartTime = performance.now();
              // start data collection
              window.addEventListener("mousemove", this.mouseEventHandler);
          };
          this.on_finish = () => {
              this.domObserver.disconnect();
              window.removeEventListener("mousemove", this.mouseEventHandler);
              return {
                  mouse_tracking_data: this.currentTrialData,
                  mouse_tracking_targets: this.currentTrialTargets,
              };
          };
          this.mouseEventHandler = (e) => {
              const x = e.clientX;
              const y = e.clientY;
              const event_time = performance.now();
              const t = Math.round(event_time - this.currentTrialStartTime);
              if (this.lastSampleTime === null ||
                  event_time - this.lastSampleTime >= this.minimumSampleTime) {
                  this.lastSampleTime = event_time;
                  this.currentTrialData.push({ x, y, t });
              }
          };
          this.mutationObserverCallback = (mutationsList, observer) => {
              for (const selector of this.currentTrialSelectors) {
                  if (!this.currentTrialTargets[selector]) {
                      if (this.jsPsych.getDisplayElement().querySelector(selector)) {
                          var coords = this.jsPsych
                              .getDisplayElement()
                              .querySelector(selector)
                              .getBoundingClientRect();
                          this.currentTrialTargets[selector] = coords;
                      }
                  }
              }
          };
      }
  }
  MouseTrackingExtension.info = {
      name: "mouse-tracking",
  };

  return MouseTrackingExtension;

})();
