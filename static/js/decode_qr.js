import {
  Uppy,
  DragDrop,
  Dashboard,
  XHRUpload,
  // ScreenCapture,
  // Webcam,
} from "./uppy.min_3.25.2.min.js";

export function start() {
  // window.onload = () => {
    
  // }; //window.onload = ()=> {

    // console.log("Start started")
    const uppy = new Uppy({
      restrictions: {
        maxNumberOfFiles: 1,
      },
    });
    uppy
      .use(Dashboard, {
        target: "#file_upload_div",
        inline: true,
        hideUploadButton: true,
        singleFileFullScreen : true,
        hideRetryButton : true,
      })
      .use(XHRUpload, {
        endpoint: "../decodeQRCode",
        method: "POST",
        fieldName: "qr_file",
        bundle: true,
      });
    // .use(ScreenCapture, { target: Uppy.Dashboard })
    // .use(Webcam, { target: Uppy.Dashboard });

    document
      .getElementById("decode-qr-submit-btn")
      .addEventListener("click", (evt) => {

        document.getElementById("file_upload_div").addEventListener(
          "click", ()=> {
            document.getElementById("file_upload_div").focus()
            console.log('setting focus :>> ');
          }
        )
        if (uppy.getFiles().length > 0) {
          uppy.setMeta({
            key: document.getElementById("key").value,
          });
          uppy.upload();
        }
        uppy.on("upload-success", (file, response) => {
          console.log('answerJSON:>> ', response.body);
          if (response.body.ErrOccurred) {
            console.log("Error occurred")
            new Notify({
              status: "error",
              title: "Error",
              text: "Problem with QR Image uploaded:<br><br>"+
                  response.body.ErrDescription+
                  "<br><br> Check if QR image does not have any artifacts (e.g. mouse pointer)."              ,
              effect: "fade",
              speed: 300,
              customClass: "",
              customIcon: "",
              showIcon: true,
              showCloseButton: true,
              autoclose: false,
              autotimeout: 15000,
              gap: 20,
              distance: 20,
              type: 1,
              position: "right top",
            });
          } else {
            sessionStorage.setItem("DecodedQRImageText", response.body.answer);
            window.location.href = "decode_qr_result.html";
          }
          
        });
      });

    var els = document.getElementsByClassName("uppy-Dashboard-AddFiles-title");

}//export function start() {
