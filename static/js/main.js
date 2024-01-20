var GLOBAL_STATE = "qwe1";

//From: https://stackoverflow.com/questions/9719570/generate-random-password-string-with-requirements-in-javascript
//Specifically: https://stackoverflow.com/a/26528271
//TODO: review code below - no particular concerns, but may be to make it more secure
var Password = {
  _pattern: /[a-zA-Z0-9_\-\+\.\?\*]/,

  _getRandomByte: function () {
    // http://caniuse.com/#feat=getrandomvalues
    if (window.crypto && window.crypto.getRandomValues) {
      var result = new Uint8Array(1);
      // console.log("Have window.crypto");
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
  // console.log("Submit click event");
  event.preventDefault();

  text2Encode = document.getElementById("text2encrypt").value;
  // console.log("Text:", text2Encode);
  key2Encode = document.getElementById("key").value;

  if (text2Encode.length == 0) {
    // alert("Text to encrypt cannot be empty.");
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
  }

  if (key2Encode.length == 0) {
    // alert("Key cannot be empty.");
    // let demoErrElem = SnackBar({
    //   message: "Key cannot be empty",
    //   status: "danger",
    // });
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
  }

  encodedBody = JSON.stringify({
    Text2Encrypt: text2Encode,
    Key: key2Encode,
  });

  // console.log("Encoded body:", encodedBody);

  fetch("../receiveData2Encrypt", {
    method: "POST",
    body: encodedBody,
    headers: {
      "Content-type": "application/json; charset=UTF-8",
    },
  })
    .then((response) => response.json())
    .then((json) => {
      // console.log("2nd Then:", json);
      // console.log("IMG path:", json.ImagePath);

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
        if(json.ErrDescription.includes("Error generating QRCode.")) {
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
        //
        console.log("Error creating QRcode:" + json.ErrDescription);
      } else {
        // GLOBAL_STATE = "qwe2"
        // window.GLOBAL_STATE = "qwe3"
        sessionStorage.setItem("EncryptedText", json.EncryptedText);
        sessionStorage.setItem("ImagePath", json.ImagePath);
        // console.log("Global state:", window.GLOBAL_STATE)
        window.location.href = "result.html";
      }
    });
};//function onSubmitClick(event) {


//TODO: write number of bytes and characters red when bytes are over 1200
function updateCharEnteredCounter(evt) {
  text_entered = document.getElementById("text2encrypt").value;
  numChars = text_entered.length;
  // console.log("numChars:", numChars);
  if (numChars == 0) {
    document.getElementById("char-counters-para").classList.add("no-show");
  } else {
    document.getElementById("char-counters-para").classList.remove("no-show");
    let utf8Encode = new TextEncoder();
    uintArr = utf8Encode.encode(text_entered);
    // console.log("UINT Array:", uintArr)
    // console.log("Len of uint array:", uintArr.length)

    document.getElementById("char-entered-placeholder").innerHTML = numChars;
    document.getElementById("bytes-entered-placeholder").innerHTML =
      uintArr.length;
  }
} //function updateCharEnteredCounter(evt) {

function generatePassword() {
  document.getElementById("key").value = Password.generate(32);
}

window.onload = () => {
  // console.log("I am in.");

  document
    .getElementById("form-submit-btn")
    .addEventListener("click", onSubmitClick);

  document
    .getElementById("text2encrypt")
    .addEventListener("input", updateCharEnteredCounter);

  document
    .getElementById("generate-pass-btn")
    .addEventListener("click", generatePassword);
}; //window.onload = ()=> {
