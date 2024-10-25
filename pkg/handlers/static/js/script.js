function ready(fn) {
  if (document.readyState !== "loading") {
    fn();
  } else {
    document.addEventListener("DOMContentLoaded", fn);
  }
}

function displayError(error) {
  hideMessage();
  document.getElementById("error").innerHTML = error;
  document.getElementById("error").classList.remove("dn");
}

function hideError() {
  document.getElementById("error").innerHTML = "";
  document.getElementById("error").classList.add("dn");
}

function displayMessage(message) {
  hideError();
  document.getElementById("message").innerHTML = message;
  document.getElementById("message").classList.remove("dn");
}

function hideMessage() {
  document.getElementById("message").innerHTML = "";
  document.getElementById("message").classList.add("dn");
}

ready(function() {
  const errorDivId = "error";
  document.body.addEventListener("htmx:responseError", function(e) {
    document.getElementById(errorDivId).innerHTML = e.detail.xhr.response;
    document.getElementById(errorDivId).classList.remove("dn");
  });
  document.body.addEventListener("htmx:afterOnLoad", function(e) {
    if (e.detail.successful) {
      document.getElementById(errorDivId).innerHTML = "";
      document.getElementById(errorDivId).classList.add("dn");
    }
  });

  // history and navigation events
  document.body.addEventListener("htmx:beforeRequest", function(e) {
    document.getElementById("loader").classList.remove("dn");
  });
  document.body.addEventListener("htmx:afterRequest", function(e) {
    document.getElementById("loader").classList.add("dn");
  });
  document.body.addEventListener("htmx:historyRestore", function(e) {
    document.getElementById("loader").classList.add("dn");
  });

  document.body.addEventListener("htmx:wsError", function(e) {
    displayError("WebSocket Error: " + JSON.stringify(e));
  });
  document.body.addEventListener("htmx:wsConnecting", function(e) {
    displayMessage("WebSocket Connecting...");
  });
  document.body.addEventListener("htmx:wsOpen", function(e) {
    hideMessage();
    hideError();
  });
});
