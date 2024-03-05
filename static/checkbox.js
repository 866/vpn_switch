async function sendRequest(upVPN) {
    const response = await fetch("/vpn", {
      method: "POST", // *GET, POST, PUT, DELETE, etc.
      mode: "cors", // no-cors, *cors, same-origin
      cache: "no-cache", // *default, no-cache, reload, force-cache, only-if-cached
      credentials: "same-origin", // include, *same-origin, omit
      headers: {
        "Content-Type": "application/json",
        // 'Content-Type': 'application/x-www-form-urlencoded',
      },
      redirect: "follow", // manual, *follow, error
      referrerPolicy: "no-referrer", // no-referrer, *no-referrer-when-downgrade, origin, origin-when-cross-origin, same-origin, strict-origin, strict-origin-when-cross-origin, unsafe-url
      body: JSON.stringify({"vpn": upVPN}), // body data type must match "Content-Type" header
    });
    const post = await response.json();
    console.log(post)
  }

async function usbOff() {
  const response = await fetch("/usboff");
  const get = await response.json();
  console.log(get)
}
  
const switchEl = document.getElementById('switch');
  
switchEl.addEventListener('change', (event) => {
  sendRequest(event.currentTarget.checked)
});
  