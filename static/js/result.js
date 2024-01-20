window.onload = () => {
  // console.log("Result window: I am in.");
  // console.log("GLOBAL_STATE:", sessionStorage.getItem("key"));

  ImagePath = sessionStorage.getItem("ImagePath");
  EncryptedText = sessionStorage.getItem("EncryptedText");
  // console.log(EncryptedText);
  // console.log("EncrText len:", EncryptedText.length);

  if(EncryptedText.length > 900) {
    // console.log("Page breaker added");
    document.getElementById("page-breaker-para").classList.add("page-break");
  }

  document
    .getElementById("qr-code-placeholder")
    .setAttribute("src", ImagePath);
  document.getElementById("crypted-message-placeholder").innerHTML =
    EncryptedText;

    currTime = Date.now()
    formatted_date = (new Intl.DateTimeFormat('en-GB', {
      dateStyle: 'full',
      timeStyle: 'long',
      timeZone: 'UTC',
    })).format(currTime)

    // console.log("date:", formatted_date)
    document.getElementById("datetime-placeholder").innerHTML = formatted_date

    document.getElementById("show-cipher-text-checkbox").addEventListener("change",
      (evt)=> {
        // console.log("Checkbox changed", document.getElementById("show-cipher-text-checkbox").checked)
        if(document.getElementById("show-cipher-text-checkbox").checked) {
          document.getElementById("encrypted-message-div").classList.add("no-print")
          document.getElementById("encrypted-message-div").classList.add("no-show")
          
          // console.log(document.getElementById("crypted-message-placeholder").classList)
        } else {
          document.getElementById("encrypted-message-div").classList.remove("no-print")
          document.getElementById("encrypted-message-div").classList.remove("no-show")
        }

      }
    )

}; //window.onload = ()=> {



