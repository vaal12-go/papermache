executeOnload(() => {
  initResultWindow();

});

function initResultWindow() {
  ImagePath = sessionStorage.getItem("ImagePath");
  EncryptedText = sessionStorage.getItem("EncryptedText");
  OpenComment = sessionStorage.getItem("OpenComment");
  document.getElementById("openComment").innerHTML =
    "Comment: "+OpenComment;
// 
  // console.log('OpenComment :>> ', OpenComment);
  document.getElementById("page1-number-placeholder").innerHTML = "Page <strong>1</strong> of 1";
  if(EncryptedText.length > 900) {
    document.getElementById("page-breaker-para").classList.add("page-break");
    document.getElementById("page1-number-placeholder").innerHTML = "Page <strong>1</strong> of 2";
    document.getElementById("page2-number-placeholder").innerHTML = "Page <strong>2</strong> of 2";
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
    document.getElementById("datetime-placeholder").innerHTML = "Generated:"+formatted_date
    document.getElementById("show-cipher-text-checkbox").addEventListener("change",
      (evt)=> {
        if(document.getElementById("show-cipher-text-checkbox").checked) {
          document.getElementById("encrypted-message-div").classList.add("no-print")
          document.getElementById("encrypted-message-div").classList.add("no-show")
        } else {
          document.getElementById("encrypted-message-div").classList.remove("no-print")
          document.getElementById("encrypted-message-div").classList.remove("no-show")
        }
      }
    )

    
}; //initResultWindow() {



