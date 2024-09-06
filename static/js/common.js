var prev_on_load = window.onload;

const VERSION_STRING = "ver 0.2.1 06Sep2024"

var onload_func_arr = [];
function executeOnload(fn2Exec) {
  onload_func_arr.push(fn2Exec);
}

executeOnload(() => {
  // console.log("Hello I should execute :>> ");

  ver_spanEl = document.getElementById("version_span")
  if (ver_spanEl)
      ver_spanEl.innerHTML = VERSION_STRING;
  document.title = "Papier-mâché "+VERSION_STRING;

  stopSrvrEl = document.getElementById("stop-server")
  if(stopSrvrEl)
    stopSrvrEl.addEventListener("click", (evt) => {
    // console.log("Stop server clicked :>> ");
    //https://getbootstrap.com/docs/4.0/components/modal/
    fetch(BASE_URL + "stopServer", {
      method: "GET",
      headers: {
        "Content-type": "application/json; charset=UTF-8",
      },
    })
      .then((response) => response.json())
      .then((json) => {
        if (json.result == "ServerShutdown success")
            $('#server-stopped-modal').on('hide.bs.modal', function(e) {
                e.preventDefault();
            });
          $("#server-stopped-modal").modal("show", { show: true, backdrop: "static", keyboard: false });
      });
  });
}); //executeOnload(() => {

window.onload = () => {
  // console.log("I am loaded2 :>> ");
  // console.log("onload_func_arr :>> ", onload_func_arr);
  for (idx in onload_func_arr) {
    // console.log("fn :>> ", onload_func_arr[idx]);
    onload_func_arr[idx]();
  }
}; //window.onload = ()=> {
