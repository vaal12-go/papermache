// var GLOBAL_STATE = "qwe1";



//From: https://stackoverflow.com/questions/9719570/generate-random-password-string-with-requirements-in-javascript
//Specifically: https://stackoverflow.com/a/26528271
//TODO: review code below - no particular concerns, but may be to make it more secure
var Password = {
  _pattern: /[a-zA-Z0-9_\-\+\.\?\*]/,
  _getRandomByte: function () {
    // http://caniuse.com/#feat=getrandomvalues
    if (window.crypto && window.crypto.getRandomValues) {
      var result = new Uint8Array(1);
      window.crypto.getRandomValues(result);
      return result[0];
    } else if (window.msCrypto && window.msCrypto.getRandomValues) {
      var result = new Uint8Array(1);
      window.msCrypto.getRandomValues(result);
      return result[0];
    } else {
      return Math.floor(Math.random() * 256);
    }
  },

  generate: function (length) {
    return Array.apply(null, { length: length })
      .map(function () {
        var result;
        while (true) {
          result = String.fromCharCode(this._getRandomByte());
          if (this._pattern.test(result)) {
            return result;
          }
        }
      }, this)
      .join("");
  },
}; //var Password = {

//TODO: move to

function onSubmitClick(event) {
  event.preventDefault();
  text2Encode = document.getElementById("text2encrypt").value;
  key2Encode = document.getElementById("key").value;
  if (text2Encode.length == 0) {
    new Notify({
      status: "error",
      title: "Error",
      text: "Text to encrypt cannot be empty",
      effect: "fade",
      speed: 300,
      customClass: "",
      customIcon: "",
      showIcon: true,
      showCloseButton: true,
      autoclose: true,
      autotimeout: 5000,
      gap: 20,
      distance: 20,
      type: 1,
      position: "right top",
    });
    return;
  }//if (text2Encode.length == 0) {

  if (key2Encode.length == 0) {
    new Notify({
      status: "error",
      title: "Error",
      text: "Key cannot be empty",
      effect: "fade",
      speed: 300,
      customClass: "",
      customIcon: "",
      showIcon: true,
      showCloseButton: true,
      autoclose: true,
      autotimeout: 5000,
      gap: 20,
      distance: 20,
      type: 1,
      position: "right top",
    });
    return;
  }//if (key2Encode.length == 0) {

  encodedBody = JSON.stringify({
    Text2Encrypt: text2Encode,
    Key: key2Encode,
  });


  fetch(BASE_URL + "receiveData2Encrypt", {
    method: "POST",
    body: encodedBody,
    headers: {
      "Content-type": "application/json; charset=UTF-8",
    },
  })
    .then((response) => response.json())
    .then((json) => {
      if (json.ErrOccurred) {
        new Notify({
          status: "error",
          title: "Error",
          text: `Error received from qrCode generator:<p>${json.ErrDescription}</p>`,
          effect: "fade",
          speed: 300,
          customClass: "",
          customIcon: "",
          showIcon: true,
          showCloseButton: true,
          autoclose: true,
          autotimeout: 15000,
          gap: 20,
          distance: 20,
          type: 1,
          position: "right top",
        });
        if (json.ErrDescription.includes("Error generating QRCode.")) {
          new Notify({
            status: "warning",
            title: "Information",
            text: `This error most probably is due text to be encrypted is too long. Shorten the text to about 1200 bytes to fit.`,
            effect: "fade",
            speed: 300,
            customClass: "",
            customIcon: "",
            showIcon: true,
            showCloseButton: true,
            autoclose: true,
            autotimeout: 30000,
            gap: 20,
            distance: 20,
            type: 1,
            position: "right top",
          });
        }
        // console.log("Error creating QRcode:" + json.ErrDescription);
      } else {
        sessionStorage.setItem("EncryptedText", json.EncryptedText);
        sessionStorage.setItem("ImagePath", json.ImagePath);
        var comment = document.getElementById("openComment").value;
        sessionStorage.setItem("OpenComment", comment);
        window.location.href = BASE_URL + "/static/result.html";
      }
    });
} //function onSubmitClick(event) {

//TODO: write number of bytes and characters red when bytes are over 1200
function updateCharEnteredCounter(evt) {
  text_entered = document.getElementById("text2encrypt").value;
  numChars = text_entered.length;
  let utf8Encode = new TextEncoder();
  uintArr = utf8Encode.encode(text_entered);
  document.getElementById("char-entered-placeholder").innerHTML = numChars;
  document.getElementById("bytes-entered-placeholder").innerHTML =
    uintArr.length;
} //function updateCharEnteredCounter(evt) {

function generatePassword() {
  document.getElementById("key").value = Password.generate(32);
}

executeOnload(() => {
  document
    .getElementById("form-submit-btn")
    .addEventListener("click", onSubmitClick);

  document
    .getElementById("text2encrypt")
    .addEventListener("input", updateCharEnteredCounter);

  document
    .getElementById("generate-pass-btn")
    .addEventListener("click", generatePassword);
});//executeOnload(() => {
