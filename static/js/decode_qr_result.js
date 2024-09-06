// window.onload = ()=> {
executeOnload(() => {
    console.log('Execute on load works :>> ');
    decodedText = sessionStorage.getItem("DecodedQRImageText");
    document.getElementById("decodedText").innerHTML = decodedText;
});